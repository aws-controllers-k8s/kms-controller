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

"""Integration tests for the KMS Alias resource
"""

import boto3
import logging
import pytest
import time

from acktest.k8s import resource as k8s
from acktest.k8s.condition import CONDITION_TYPE_RESOURCE_SYNCED, CONDITION_TYPE_REFERENCES_RESOLVED
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_kms_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.tests.helper import KMSValidator
from e2e.tests.test_key import KEY_RESOURCE_PLURAL

CREATE_WAIT_AFTER_SECONDS = 30
MODIFY_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_PERIODS = 3
DELETE_WAIT_PERIOD_LENGTH_SECONDS = 10

ALIAS_RESOURCE_PLURAL = "aliases"

kms_validator = KMSValidator(boto3.client('kms'))

@pytest.fixture
def simple_key(kms_client):
    key = kms_client.create_key()

    yield key

    kms_client.schedule_key_deletion(KeyId=key['KeyMetadata']['KeyId'])

@pytest.fixture
def another_key(kms_client):
    key = kms_client.create_key()

    yield key

    kms_client.schedule_key_deletion(KeyId=key['KeyMetadata']['KeyId'])

@pytest.fixture
def simple_alias(simple_key):
    alias_name = random_suffix_name("simple-alias", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["ALIAS_NAME"] = alias_name
    replacements["TARGET_KEY_ID"] = simple_key['KeyMetadata']['KeyId']

    resource_data = load_kms_resource(
        "alias_simple",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, ALIAS_RESOURCE_PLURAL,
        alias_name, namespace="default",
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
class TestAlias:
    def test_create_update_delete_alias(self, simple_alias, another_key):
        (ref, cr) = simple_alias
        another_key_id = another_key['KeyMetadata']['KeyId']

        assert k8s.wait_on_condition(ref, CONDITION_TYPE_RESOURCE_SYNCED, "True", wait_periods=10)
        assert 'arn' in cr['status']['ackResourceMetadata']
        alias_arn = cr['status']['ackResourceMetadata']['arn']

        # Update
        updates = {
            "spec": {
                "targetKeyID": another_key_id
            }
        }

        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(ref, CONDITION_TYPE_RESOURCE_SYNCED, "True", wait_periods=10)
        alias = kms_validator.get_alias(arn=alias_arn, target_key_id=another_key_id)
        assert alias is not None, f"Alias should not be None for key id {another_key_id}"

        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_PERIODS, DELETE_WAIT_PERIOD_LENGTH_SECONDS)
        assert deleted
        time.sleep(DELETE_WAIT_AFTER_SECONDS)
        kms_validator.assert_alias_deleted(arn=alias_arn, target_key_id=another_key_id)

    def test_create_alias_ref(self):
        key_name = random_suffix_name("ref-key", 15)
        alias_name = random_suffix_name("ref-alias", 20)
        replacements = REPLACEMENT_VALUES.copy()
        replacements["REF_KEY_NAME"] = key_name
        replacements["REF_ALIAS_NAME"] = alias_name

        key_res_data = load_kms_resource(
            "key_ref",
            additional_replacements=replacements,
        )

        key_ref = k8s.CustomResourceReference(
            CRD_GROUP, CRD_VERSION, KEY_RESOURCE_PLURAL,
            key_name, namespace="default",
        )

        logging.debug(f"key resource. name: {key_name}, data: {key_res_data}")

        alias_res_data = load_kms_resource(
                "alias_ref",
                additional_replacements=replacements,
            )
        logging.debug(alias_res_data)

        # Create the k8s resource
        alias_ref = k8s.CustomResourceReference(
            CRD_GROUP, CRD_VERSION, ALIAS_RESOURCE_PLURAL,
            alias_name, namespace="default",
        )

        logging.debug(f"alias resource. name: {alias_name}, data: {alias_res_data}")

        # create the resources in order that initially the reference resolution fails and
        # then when the referenced resource gets created, then all resolutions eventually
        # pass and resources get synced.

        # Create Alias. Needs Key reference
        k8s.create_custom_resource(alias_ref, alias_res_data)

        # Create Key. Needs no reference
        k8s.create_custom_resource(key_ref, key_res_data)

        time.sleep(CREATE_WAIT_AFTER_SECONDS)

        assert k8s.wait_on_condition(key_ref, CONDITION_TYPE_RESOURCE_SYNCED, "True", wait_periods=10)
        assert k8s.wait_on_condition(alias_ref, CONDITION_TYPE_RESOURCE_SYNCED, "True", wait_periods=10)

        assert k8s.wait_on_condition(alias_ref, CONDITION_TYPE_REFERENCES_RESOLVED, "True", wait_periods=10)

        alias_cr = k8s.get_resource(alias_ref)
        alias_arn = alias_cr['status']['ackResourceMetadata']['arn']

        key_cr = k8s.get_resource(key_ref)
        key_id = key_cr['status']['keyID']

        kms_validator.assert_alias_exists(arn=alias_arn, target_key_id=key_id)

        # DELETE
        k8s.delete_custom_resource(alias_ref)
        time.sleep(DELETE_WAIT_AFTER_SECONDS)

        k8s.delete_custom_resource(key_ref)
        time.sleep(DELETE_WAIT_AFTER_SECONDS)

        assert not k8s.get_resource_exists(alias_ref)
        assert not k8s.get_resource_exists(key_ref)

        # check that the alias resources does not exist in AWS KMS
        kms_validator.assert_alias_deleted(arn=alias_arn, target_key_id=key_id)
