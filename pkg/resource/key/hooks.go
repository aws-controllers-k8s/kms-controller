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
	"context"
	"fmt"
	"strconv"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/kms-controller/apis/v1alpha1"
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

// customUpdate is the implementation of update operation for KMS Key resource.
// This operation will return a Terminal error if the update is not for 'Policy'
// , 'Tags' or 'BypassPolicyLockoutSafetyCheck' because KMS Key only supports
// updating those three fields.
func (rm *resourceManager) customUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdate")
	defer func() {
		exit(err)
	}()
	updatedRes := rm.concreteResource(desired.DeepCopy())
	updatedRes.SetStatus(latest)

	if delta.DifferentExcept("Spec.Policy", "Spec.Tags", "Spec.BypassPolicyLockoutSafetyCheck") {
		return updatedRes, ackerr.NewTerminalError(
			fmt.Errorf("KMS Key only supports update for Policy and Tags"),
		)
	}
	if delta.DifferentAt("Spec.Policy") {
		if updatedRes.ko.Spec.Policy != nil && *updatedRes.ko.Spec.Policy != "" {
			if err = rm.updatePolicy(ctx, updatedRes); err != nil {
				return updatedRes, err
			}
		}
	}
	if delta.DifferentAt("Spec.Tags") {
		err = rm.removeOldTags(ctx, updatedRes)
		if err != nil {
			return updatedRes, err
		}
		err = rm.updateTags(ctx, updatedRes)
		if err != nil {
			return updatedRes, err
		}
	}
	rm.setStatusDefaults(updatedRes.ko)
	return updatedRes, nil
}
