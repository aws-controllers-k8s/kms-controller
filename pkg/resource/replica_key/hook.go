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

package replica_key

import (
	"context"
	"fmt"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
)

// updateNotSupported returns a terminal error because KMS replica key
// resource does not support updates.
func (rm *resourceManager) updateNotSupported(ctx context.Context,
	desired *resource,
	latest *resource,
	diffReporter *ackcompare.Delta,
) (*resource, error) {
	return nil, ackerr.NewTerminalError(
		fmt.Errorf("replica key resource does not support updates"),
	)
}