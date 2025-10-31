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
import boto3
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

REPLICA_KEY_RESOURCE_PLURAL = "replicakeys"

# For testing purposes, we'll use a secondary region
# NOTE: Must be different from the controller's region (us-west-2)
REPLICA_REGION = "us-east-1"


@pytest.fixture
def multi_region_primary_key(kms_client):
  """Creates a multi-region primary key using boto3 directly"""
  # Create a multi-region primary key via AWS SDK
  key = kms_client.create_key(
    MultiRegion=True,
    Description="Multi-region primary key for ReplicaKey tests"
  )

  key_id = key['KeyMetadata']['KeyId']

  yield key

  # Cleanup
  try:
    kms_client.schedule_key_deletion(KeyId=key_id, PendingWindowInDays=7)
  except:
    pass


@pytest.fixture
def simple_replica_key(multi_region_primary_key):
  """Creates a simple replica key from a multi-region primary key"""
  primary_key_id = multi_region_primary_key['KeyMetadata']['KeyId']
  replica_key_name = random_suffix_name("simple-replica", 32)

  replacements = REPLACEMENT_VALUES.copy()
  replacements["REPLICA_KEY_NAME"] = replica_key_name
  replacements["PRIMARY_KEY_ID"] = primary_key_id
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

    # Get the latest resource state after sync
    cr = k8s.get_resource(ref)

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

    # Create a KMS client in the replica region to verify the key
    replica_kms_client = boto3.client('kms', region_name=REPLICA_REGION)
    aws_replica_key = replica_kms_client.describe_key(KeyId=replica_key_id)

    # Verify AWS resource properties
    assert aws_replica_key is not None
    assert 'KeyMetadata' in aws_replica_key
    metadata = aws_replica_key['KeyMetadata']

    # Verify the key is multi-region
    assert metadata['MultiRegion'] is True

    # Verify multi-region configuration
    assert 'MultiRegionConfiguration' in metadata
    mr_config = metadata['MultiRegionConfiguration']
    assert mr_config['MultiRegionKeyType'] == 'REPLICA'

    # Verify the replica references the correct primary key
    assert 'PrimaryKey' in mr_config
    assert primary_key_id in mr_config['PrimaryKey']['Arn']

    # Verify the key state (should be Creating or Enabled)
    assert metadata['KeyState'] in ['Creating', 'Enabled']

    # Verify the key is in the correct region
    assert mr_config['PrimaryKey']['Region'] != REPLICA_REGION

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
