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

package key_test

import (
	"testing"

	svcapitypes "github.com/aws-controllers-k8s/kms-controller/apis/v1alpha1"
	key "github.com/aws-controllers-k8s/kms-controller/pkg/resource/key"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_GetDeletePendingWindowInDays(t *testing.T) {
	assert := assert.New(t)

	noAnnotation := metav1.ObjectMeta{
		Annotations: map[string]string{},
	}
	badAnnotation := metav1.ObjectMeta{
		Annotations: map[string]string{
			svcapitypes.AnnotationDeletePendingWindow: "not-an-int",
		},
	}
	validAnnotation := metav1.ObjectMeta{
		Annotations: map[string]string{
			svcapitypes.AnnotationDeletePendingWindow: "25",
		},
	}

	assert.Equal(key.GetDeletePendingWindowInDays(&noAnnotation), key.DefaultDeletePendingWindowInDays)
	assert.Equal(key.GetDeletePendingWindowInDays(&badAnnotation), key.DefaultDeletePendingWindowInDays)
	assert.Equal(key.GetDeletePendingWindowInDays(&validAnnotation), int64(25))
}
