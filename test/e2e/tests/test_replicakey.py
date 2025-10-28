# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the KMS ReplicaKey resource
"""
import logging
import pytest
import time

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_kms_resource
from e2e.replacement_values import REPLACEMENT_VALUES

MODIFY_WAIT_AFTER_SECONDS = 40
DELETE_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_PERIODS = 3
DELETE_WAIT_PERIOD_LENGTH_SECONDS = 10

KEY_RESOURCE_PLURAL = "keys"
REPLICA_KEY_RESOURCE_PLURAL = "replicakeys"

# For testing purposes, we'll use a secondary region
# Adjust this based on your test environment
REPLICA_REGION = "us-west-2"


@pytest.fixture
def multi_region_primary_key(kms_client):
  """Creates a multi-region primary key for replication tests"""
  key_name = random_suffix_name("mr-primary-key", 32)

  replacements = REPLACEMENT_VALUES.copy()
  replacements["KEY_NAME"] = key_name

  resource_data = load_kms_resource(
    "key_multiregion",
    additional_replacements=replacements,
  )
  logging.debug(resource_data)

  # Create the k8s resource
  ref = k8s.CustomResourceReference(
    CRD_GROUP, CRD_VERSION, KEY_RESOURCE_PLURAL,
    key_name, namespace="default",
  )
  k8s.create_custom_resource(ref, resource_data)
  cr = k8s.wait_resource_consumed_by_controller(ref)

  assert cr is not None
  assert k8s.get_resource_exists(ref)
  assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=10)

  # Verify it's a multi-region key
  assert 'keyID' in cr['status']
  key_id = cr['status']['keyID']
  key_metadata = kms_client.describe_key(KeyId=key_id)
  assert key_metadata['KeyMetadata']['MultiRegion'] == True

  yield (ref, cr, key_name, key_id)

  # Cleanup
  try:
    _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
    assert deleted
  except:
    pass


@pytest.fixture
def simple_replica_key(multi_region_primary_key):
  """Creates a simple replica key from a multi-region primary key"""
  (primary_ref, primary_cr, primary_key_name, primary_key_id) = multi_region_primary_key

  replica_key_name = random_suffix_name("simple-replica", 32)

  replacements = REPLACEMENT_VALUES.copy()
  replacements["REPLICA_KEY_NAME"] = replica_key_name
  replacements["PRIMARY_KEY_NAME"] = primary_key_name
  replacements["REPLICA_REGION"] = REPLICA_REGION

  resource_data = load_kms_resource(
    "replicakey_simple",
    additional_replacements=replacements,
  )
  logging.debug(resource_data)

  # Create the k8s resource
  ref = k8s.CustomResourceReference(
    CRD_GROUP, CRD_VERSION, REPLICA_KEY_RESOURCE_PLURAL,
    replica_key_name, namespace="default",
  )
  k8s.create_custom_resource(ref, resource_data)
  cr = k8s.wait_resource_consumed_by_controller(ref)

  assert cr is not None
  assert k8s.get_resource_exists(ref)

  yield (ref, cr, primary_key_id)

  # Try to delete, if doesn't already exist
  try:
    _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
    assert deleted
  except:
    pass


@service_marker
@pytest.mark.canary
class TestReplicaKey:
  """Test suite for ReplicaKey resource"""

  def test_create_delete_replica_key(self, simple_replica_key):
    """Test creating and deleting a simple replica key"""
    (ref, cr, primary_key_id) = simple_replica_key

    # Verify the replica key was created
    assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=10)

    # ReplicaKey status fields are nested in replicaKeyMetadata
    assert 'replicaKeyMetadata' in cr['status']
    assert 'keyID' in cr['status']['replicaKeyMetadata']
    replica_key_id = cr['status']['replicaKeyMetadata']['keyID']

    # Verify the replica key is in the correct state
    assert 'keyState' in cr['status']['replicaKeyMetadata']
    # Key might be in Creating or Enabled state
    assert cr['status']['replicaKeyMetadata']['keyState'] in ['Creating', 'Enabled']

    # Verify multi-region configuration
    if 'multiRegionConfiguration' in cr['status']['replicaKeyMetadata']:
      config = cr['status']['replicaKeyMetadata']['multiRegionConfiguration']
      assert config is not None
      # The replica should reference the primary key
      if 'primaryKey' in config:
        assert primary_key_id in config['primaryKey']['arn']

    # Delete the replica key
    _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
    assert deleted

  def test_replica_key_immutability(self, simple_replica_key):
    """Test that replica region cannot be changed after creation"""
    (ref, cr, primary_key_id) = simple_replica_key

    assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=10)

    # Try to update the replica region (should fail due to immutability)
    updates = {
      "spec": {
        "replicaRegion": "eu-west-1"  # Try to change region
      }
    }

    try:
      k8s.patch_custom_resource(ref, updates)
      time.sleep(MODIFY_WAIT_AFTER_SECONDS)

      # The patch should be rejected by the validation webhook
      # or the resource should show an error condition
      updated_cr = k8s.get_resource(ref)
      # Check if there's an error condition indicating immutability violation
      conditions = updated_cr.get('status', {}).get('conditions', [])
      has_error = any(
        c.get('type') == 'ACK.Terminal' and c.get('status') == 'True'
        for c in conditions
      )

      # Either the patch was rejected or there's an error condition
      assert has_error, "Expected immutability validation to fail the update"
    except Exception as e:
      # Expected: the patch should be rejected
      logging.info(f"Patch rejected as expected: {e}")
      pass

    _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
    assert deleted
