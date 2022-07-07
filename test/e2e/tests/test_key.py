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

"""Integration tests for the KMS Key resource
"""
import json
import logging
import pytest
import time

from datetime import datetime, timedelta

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest.aws.identity import get_region, get_account_id
from acktest import tags
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_kms_resource
from e2e.replacement_values import REPLACEMENT_VALUES

MODIFY_WAIT_AFTER_SECONDS = 40
DELETE_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_PERIODS = 3
DELETE_WAIT_PERIOD_LENGTH_SECONDS = 10

KEY_RESOURCE_PLURAL = "keys"

PENDING_WINDOW_IN_DAYS = 8

@pytest.fixture
def simple_key():
    key_name = random_suffix_name("simple-key", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["KEY_NAME"] = key_name

    resource_data = load_kms_resource(
        "key_simple",
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

    yield (ref, cr)

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted
    except:
        pass

@pytest.fixture
def delete_annotated_key():
    key_name = random_suffix_name("annotated-key", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["KEY_NAME"] = key_name
    replacements["PENDING_WINDOW"] = str(PENDING_WINDOW_IN_DAYS)

    resource_data = load_kms_resource(
        "key_annotated",
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

    yield (ref, cr)

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted
    except:
        pass

@pytest.fixture
def key_with_policy():
    key_name = random_suffix_name("key-w-policy", 32)
    key_description = 'Used for ACK key_with_policy testing'
    account_id = get_account_id()
    key_policy = {
        "Version": "2012-10-17",
        "Id": "ack-key-with-policy",
        "Statement": [
            {
                "Sid": "Enable IAM User Permissions",
                "Effect": "Allow",
                "Principal": {
                    "AWS": f'arn:aws:iam::{account_id}:root'
                },
                "Action": "kms:*",
                "Resource": "*"
            }
        ]
    }

    replacements = REPLACEMENT_VALUES.copy()
    replacements["KEY_NAME"] = key_name
    replacements["DESCRIPTION"] = key_description
    replacements["KEY_POLICY"] = json.dumps(key_policy)

    resource_data = load_kms_resource(
        "key_with_policy",
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

    yield (ref, cr, key_policy)

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted
    except:
        pass


@service_marker
@pytest.mark.canary
class TestKey:
    def _assert_key_alive(self, key):
        assert 'DeletionDate' not in key['KeyMetadata']

    def _assert_key_deleted(self, key, pending_window: int = None):
        assert 'DeletionDate' in key['KeyMetadata']
        assert 'KeyState' in key['KeyMetadata']
        assert key['KeyMetadata']['KeyState'] == 'PendingDeletion'

        if pending_window is None:
            return

        expected_end_date = datetime.today() + timedelta(days=pending_window)
        deletion_date = key['KeyMetadata']['DeletionDate']
        assert deletion_date.date() == expected_end_date.date()

    def test_create_delete_key(self, kms_client, simple_key):
        (ref, cr) = simple_key

        assert 'keyID' in cr['status']
        key_id = cr['status']['keyID']

        key = kms_client.describe_key(KeyId=key_id)
        assert key['KeyMetadata']['KeyId'] == key_id
        self._assert_key_alive(key)

        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted

    def test_update_key_policy(self, kms_client, key_with_policy):
        (ref, cr, input_policy) = key_with_policy

        assert 'keyID' in cr['status']
        key_id = cr['status']['keyID']

        key_policy = kms_client.get_key_policy(KeyId=key_id, PolicyName='default')
        assert 'ack-key-with-policy' in key_policy['Policy']

        input_policy['Id'] = 'updated-key-policy'
        updates = {
            "spec": {
                "policy": json.dumps(input_policy)
            }
        }

        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=10)

        key_policy = kms_client.get_key_policy(KeyId=key_id, PolicyName='default')
        assert 'updated-key-policy' in key_policy['Policy']

        # updating description should set terminal condition on the resource
        # because only tags and policy related fields can be updated on the Key
        # resource
        updates = {
            "spec": {
                "description": 'only policy and tags update are supported'
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(ref, "ACK.Terminal", "True", wait_periods=10)

    def test_update_tags(self, kms_client, simple_key):
        (ref, cr) = simple_key
        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=10)

        assert 'keyID' in cr['status']
        key_id = cr['status']['keyID']

        key_tags = kms_client.list_resource_tags(KeyId=key_id)['Tags']
        tags.assert_ack_system_tags(
            tags=key_tags,
            key_member_name='TagKey',
            value_member_name='TagValue'
        )

        # add new tags
        updates = {
                    "spec": {
                        "tags": [
                            {
                                "tagKey": "key1",
                                "tagValue": "value1"
                            }
                        ]
                    }
                }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=10)

        key_tags = kms_client.list_resource_tags(KeyId=key_id)['Tags']
        tags.assert_ack_system_tags(
            tags=key_tags,
            key_member_name='TagKey',
            value_member_name='TagValue'
        )
        tags.assert_equal_without_ack_tags(
            expected={"key1": "value1"},
            actual=key_tags,
            key_member_name='TagKey',
            value_member_name='TagValue'
        )

        # update existing tag
        updates = {
            "spec": {
                "tags": [
                    {
                        "tagKey": "key1",
                        "tagValue": "newValue"
                    },
                    {
                        "tagKey": "key2",
                        "tagValue": "value2"
                    }
                ]
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=10)

        key_tags = kms_client.list_resource_tags(KeyId=key_id)['Tags']
        tags.assert_ack_system_tags(
            tags=key_tags,
            key_member_name='TagKey',
            value_member_name='TagValue'
        )
        tags.assert_equal_without_ack_tags(
            expected={"key1": "newValue", "key2": "value2"},
            actual=key_tags,
            key_member_name='TagKey',
            value_member_name='TagValue'
        )

        # remove existing tag
        updates = {
            "spec": {
                "tags": [
                    {
                        "tagKey": "key3",
                        "tagValue": "value3"
                    },
                    {
                        "tagKey": "key2",
                        "tagValue": "value2"
                    }
                ]
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=10)

        key_tags = kms_client.list_resource_tags(KeyId=key_id)['Tags']
        tags.assert_ack_system_tags(
            tags=key_tags,
            key_member_name='TagKey',
            value_member_name='TagValue'
        )
        tags.assert_equal_without_ack_tags(
            expected={"key3": "value3", "key2": "value2"},
            actual=key_tags,
            key_member_name='TagKey',
            value_member_name='TagValue'
        )

        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted

    def test_delete_annotated_key(self, kms_client, delete_annotated_key):
        (ref, cr) = delete_annotated_key

        assert 'keyID' in cr['status']
        key_id = cr['status']['keyID']

        key = kms_client.describe_key(KeyId=key_id)
        assert key['KeyMetadata']['KeyId'] == key_id
        self._assert_key_alive(key)

        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted

        # Should still exist, and have a deleted timestamp
        key = kms_client.describe_key(KeyId=key_id)
        self._assert_key_deleted(key, PENDING_WINDOW_IN_DAYS)
