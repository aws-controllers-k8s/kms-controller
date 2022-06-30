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

"""Utilities for working with tags
"""

ACK_SYSTEM_NAMESPACE_TAG_KEY = "services.k8s.aws/namespace"
ACK_SYSTEM_CONTROLLER_VERSION_TAG_KEY = "services.k8s.aws/controller-version"


def assert_ack_system_tags(tag_map):
    """
    assert_ack_system_tags verifies that ACK system tags are present inside
    tag_map parameter
    TODO(vijtrip2): move this functionality to test-infra and reuse in other
    controller tests.
    """
    assert (ACK_SYSTEM_CONTROLLER_VERSION_TAG_KEY in tag_map) & (ACK_SYSTEM_NAMESPACE_TAG_KEY in tag_map),\
        "Expected both ACK system 'namespace' and 'controller-version' tags to be present"


def convert_to_map(tags):
    """
    convert_to_map converts the tags into a map of string to string
    TODO(vijtrip2): move this functionality to test-infra and reuse in other
    controller tests
    """
    tag_map = {}
    for t in tags:
        tag_map[t['TagKey']] = t['TagValue']
    return tag_map
