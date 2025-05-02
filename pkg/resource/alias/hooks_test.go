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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureAliasName(t *testing.T) {
	tests := []struct {
		name     string
		provided *string
		expected *string
	}{
		{
			name:     "alias without prefix",
			provided: stringPtr("foo"),
			expected: stringPtr("alias/foo"),
		},
		{
			name:     "alias with prefix",
			provided: stringPtr("alias/foo"),
			expected: stringPtr("alias/foo"),
		},
		{
			// This one should not be
			// expected either since
			// aliasName is required
			// field
			name:     "nil alias",
			provided: nil,
			expected: nil,
		},
		{
			// This one should be an error,
			// but expected to be handled by
			// the AWS API
			name:     "just the prefix",
			provided: stringPtr("alias/"),
			expected: stringPtr("alias/"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ensureAliasName(tt.provided)
			if actual != nil {
				assert.NotNil(t, tt.expected)
				assert.Equal(t, *tt.expected, *actual)
			} else {
				assert.Nil(t, tt.expected)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
