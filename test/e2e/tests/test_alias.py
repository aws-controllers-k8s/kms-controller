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

import logging
import pytest

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_kms_resource
from e2e.replacement_values import REPLACEMENT_VALUES

DELETE_WAIT_AFTER_SECONDS = 30

KEY_RESOURCE_PLURAL = "aliases"

@pytest.fixture
def simple_key(kms_client):
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
        CRD_GROUP, CRD_VERSION, KEY_RESOURCE_PLURAL,
        alias_name, namespace="default",
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
class TestAlias:
    def test_create_delete_alias(self, kms_client, simple_alias):
        (ref, cr) = simple_alias

        assert 'arn' in cr['status']['ackResourceMetadata']
        alias_arn = cr['status']['ackResourceMetadata']['arn']

        _, deleted = k8s.delete_custom_resource(ref, 3, 5)
        assert deleted

        for alias in kms_client.list_aliases(KeyId=cr['spec']['targetKeyID'])['Aliases']:
            if alias['AliasArn'] == alias_arn:
                # The alias was not deleted correctly
                pytest.fail("Alias was not deleted from KMS")
        