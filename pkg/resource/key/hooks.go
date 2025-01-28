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
	"strings"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/kms"
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
	if immutableFieldChanges := rm.getImmutableFieldChanges(delta); len(immutableFieldChanges) > 0 {
		msg := fmt.Sprintf("Immutable Spec fields have been modified: %s",
			strings.Join(immutableFieldChanges, ","))
		return nil, ackerr.NewTerminalError(fmt.Errorf(msg))
	}
	updatedRes := rm.concreteResource(desired.DeepCopy())
	updatedRes.SetStatus(latest)

	if delta.DifferentAt("Spec.Policy") {
		if updatedRes.ko.Spec.Policy != nil && *updatedRes.ko.Spec.Policy != "" {
			if err = rm.updatePolicy(ctx, updatedRes); err != nil {
				return updatedRes, err
			}
		}
	}
	if delta.DifferentAt("Spec.Tags") {
		err = rm.updateTags(ctx, updatedRes)
		if err != nil {
			return updatedRes, err
		}
	}
	if delta.DifferentAt("Spec.EnableKeyRotation") {
		err = rm.updateKeyRotation(ctx, updatedRes)
		if err != nil {
			return updatedRes, err
		}
	}
	rm.setStatusDefaults(updatedRes.ko)
	return updatedRes, nil
}

func (rm *resourceManager) updateKeyRotation(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateKeyRotation")

	defer func() {
		exit(err)
	}()

	keyRotationStatus, err := rm.getKeyRotationStatus(ctx, r)
	if err != nil {
		return err
	}

	// sanity check for response from key rotation status
	if keyRotationStatus == nil {
		return nil
	}

	if r.ko.Spec.EnableKeyRotation == nil {
		// check if current status of key is enabled
		if keyRotationStatus.KeyRotationEnabled {
			return rm.disableKeyRotation(ctx, r.ko.Status.KeyID)
		}
		return nil
	}

	// Check if desired state and actual state has any difference
	if keyRotationStatus.KeyRotationEnabled != *r.ko.Spec.EnableKeyRotation {
		switch *r.ko.Spec.EnableKeyRotation {
		case true:
			return rm.enableKeyRotation(ctx, r.ko.Status.KeyID)
		case false:
			return rm.disableKeyRotation(ctx, r.ko.Status.KeyID)
		}
	}

	return nil
}

// get key rotation status at the kms key
func (rm *resourceManager) getKeyRotationStatus(ctx context.Context, r *resource) (*svcsdk.GetKeyRotationStatusOutput, error) {
	keyRotationInput := svcsdk.GetKeyRotationStatusInput{
		KeyId: r.ko.Status.KeyID,
	}
	resp, err := rm.sdkapi.GetKeyRotationStatus(ctx, &keyRotationInput)
	rm.metrics.RecordAPICall("GET", "GetKeyRotationStatus", err)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// enable key rotation on the kms key
func (rm *resourceManager) enableKeyRotation(ctx context.Context, keyId *string) error {
	enableKeyRotationInput := svcsdk.EnableKeyRotationInput{
		KeyId: keyId,
	}
	_, err := rm.sdkapi.EnableKeyRotation(ctx, &enableKeyRotationInput)
	rm.metrics.RecordAPICall("UPDATE", "EnableKeyRotation", err)
	if err != nil {
		return err
	}
	return nil
}

// disable key rotation on the kms key
func (rm *resourceManager) disableKeyRotation(ctx context.Context, keyId *string) error {
	disableKeyRotationInput := svcsdk.DisableKeyRotationInput{
		KeyId: keyId,
	}
	_, err := rm.sdkapi.DisableKeyRotation(ctx, &disableKeyRotationInput)
	rm.metrics.RecordAPICall("UPDATE", "DisableKeyRotation", err)
	if err != nil {
		return err
	}
	return nil
}
