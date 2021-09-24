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

package key

import (
	"strconv"

	svcapitypes "github.com/aws-controllers-k8s/kms-controller/apis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DefaultDeletePendingWindowInDays = int64(7)
)

// GetDeletePendingWindowInDays returns the pending window (in days) as
// determined by the annotation on the object, or the default value otherwise.
func GetDeletePendingWindowInDays(
	m *metav1.ObjectMeta,
) int64 {
	resAnnotations := m.GetAnnotations()
	pendingWindow, ok := resAnnotations[svcapitypes.AnnotationDeletePendingWindow]
	if !ok {
		return DefaultDeletePendingWindowInDays
	}

	pendingWindowInt, err := strconv.Atoi(pendingWindow)
	if err != nil {
		return DefaultDeletePendingWindowInDays
	}

	return int64(pendingWindowInt)
}
