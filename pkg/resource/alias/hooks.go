// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package alias

import (
	"strings"
)

const AliasPrefix = "alias/"

// ensureAliasName accepts the name of an Alias, and
// ensures it has the alias prefix. If it does not
// it returns a name with the prefix
func ensureAliasName(name *string) *string {
	if name == nil {
		return nil
	}

	nameVal := *name
	// alias/ should be the prefix of the name
	if !strings.HasPrefix(nameVal, AliasPrefix) {
		nameVal = AliasPrefix + nameVal
	}
	return &nameVal
}
