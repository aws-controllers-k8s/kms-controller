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

import logging
import pytest

from datetime import datetime, timedelta

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_kms_resource
from e2e.replacement_values import REPLACEMENT_VALUES

DELETE_WAIT_AFTER_SECONDS = 30

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
        _, deleted = k8s.delete_custom_resource(ref, 3, 10)
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
        _, deleted = k8s.delete_custom_resource(ref, 3, 10)
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

        _, deleted = k8s.delete_custom_resource(ref, 3, 5)
        assert deleted

        # Should still exist, and have a deleted timestamp
        key = kms_client.describe_key(KeyId=key_id)
        self._assert_key_deleted(key)

    def test_delete_annotated_key(self, kms_client, delete_annotated_key):
        (ref, cr) = delete_annotated_key

        assert 'keyID' in cr['status']
        key_id = cr['status']['keyID']

        key = kms_client.describe_key(KeyId=key_id)
        assert key['KeyMetadata']['KeyId'] == key_id
        assert 'DeletionDate' not in key['KeyMetadata']

        _, deleted = k8s.delete_custom_resource(ref, 3, 5)
        assert deleted

        # Should still exist, and have a deleted timestamp
        key = kms_client.describe_key(KeyId=key_id)
        self._assert_key_deleted(key, PENDING_WINDOW_IN_DAYS)