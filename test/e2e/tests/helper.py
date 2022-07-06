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

"""Helper functions for kms tests
"""


class KMSValidator:

    def __init__(self, kms_client):
        self.kms_client = kms_client

    def assert_alias_exists(self, arn, target_key_id):
        assert self.get_alias(arn, target_key_id) is not None,\
            f"alias {arn} for key {target_key_id} is not present"

    def assert_alias_deleted(self, arn, target_key_id):
        assert self.get_alias(arn, target_key_id) is None,\
            f"Alias {arn} for key {target_key_id} is not deleted"

    def get_alias(self, arn, target_key_id):
        for alias in self.kms_client.list_aliases(KeyId=target_key_id)['Aliases']:
            if alias['AliasArn'] == arn:
                return alias
        return None
