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

"""Integration tests for the KMS Grant resource
"""

import boto3
import logging
import pytest
import time

from acktest.k8s import resource as k8s
from acktest.k8s.condition import CONDITION_TYPE_READY, CONDITION_TYPE_TERMINAL
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_kms_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.tests.helper import KMSValidator

CREATE_WAIT_AFTER_SECONDS = 30
MODIFY_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_PERIODS = 3
DELETE_WAIT_PERIOD_LENGTH_SECONDS = 10

GRANT_RESOURCE_PLURAL = "grants"

kms_validator = KMSValidator(boto3.client('kms'))

@pytest.fixture
def simple_key(kms_client):
    key = kms_client.create_key()

    yield key

    kms_client.schedule_key_deletion(KeyId=key['KeyMetadata']['KeyId'])

@pytest.fixture
def simple_grant(simple_key):
    grant_name = random_suffix_name("simple-grant", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["GRANT_NAME"] = grant_name
    replacements["KEY_ID"] = simple_key['KeyMetadata']['KeyId']

    resource_data = load_kms_resource(
        "grant_simple",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, GRANT_RESOURCE_PLURAL,
        grant_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted
    except:
        pass


@service_marker
@pytest.mark.canary
class TestGrant:
    def test_create_update_delete_grant(self, simple_key, simple_grant):
        (ref, cr) = simple_grant
        key_id = simple_key['KeyMetadata']['KeyId']

        assert k8s.wait_on_condition(ref, CONDITION_TYPE_READY, "True", wait_periods=10)
        assert 'grantID' in cr['status']
        grant_id = cr['status']['grantID']

        # validate that grant exists
        kms_validator.assert_grant_exists(grant_id, key_id)

        # Update
        updates = {
            "spec": {
                "operations": ["Decrypt"]
            }
        }

        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        # Grant resource does not support update operation, so the terminal condition should be set
        # on the resource
        assert k8s.wait_on_condition(ref, CONDITION_TYPE_TERMINAL, "True", wait_periods=10)

        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted
        time.sleep(DELETE_WAIT_AFTER_SECONDS)
        kms_validator.assert_grant_deleted(grant_id=grant_id, key_id=key_id)
