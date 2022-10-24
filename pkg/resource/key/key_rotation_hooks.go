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

	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/kms"
)

func (rm *resourceManager) updateKeyRotation(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateKeyRotation")

	defer func() {
		exit(err)
	}()

	keyRotationStatus, err := rm.getKeyRotationStatus(r.ko.Status.KeyID)
	if err != nil {
		return err
	}

	// sanity check for response from key rotation status
	if keyRotationStatus == nil || keyRotationStatus.KeyRotationEnabled == nil {
		return nil
	}

	if r.ko.Spec.EnableKeyRotation == nil {
		// check if current status of key is enabled
		// if yes, then only disable, or skip and return
		if *keyRotationStatus.KeyRotationEnabled {
			return rm.disableKeyRotation(r.ko.Status.KeyID)
		}
		return nil
	}

	// Check if desired state and actual state has any difference
	if *keyRotationStatus.KeyRotationEnabled != *r.ko.Spec.EnableKeyRotation {
		switch *r.ko.Spec.EnableKeyRotation {
		case true:
			return rm.enableKeyRotation(r.ko.Status.KeyID)
		case false:
			return rm.disableKeyRotation(r.ko.Status.KeyID)
		}
	}

	return nil
}

// get key rotation status at the kms key
func (rm *resourceManager) getKeyRotationStatus(keyId *string) (*svcsdk.GetKeyRotationStatusOutput, error) {
	keyRotationInput := svcsdk.GetKeyRotationStatusInput{
		KeyId: keyId,
	}
	resp, err := rm.sdkapi.GetKeyRotationStatus(&keyRotationInput)
	rm.metrics.RecordAPICall("GET", "GetKeyRotationStatus", err)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// enable key rotation on the kms key
func (rm *resourceManager) enableKeyRotation(keyId *string) error {
	enableKeyRotationInput := svcsdk.EnableKeyRotationInput{
		KeyId: keyId,
	}
	_, err := rm.sdkapi.EnableKeyRotation(&enableKeyRotationInput)
	rm.metrics.RecordAPICall("POST", "EnableKeyRotation", err)
	if err != nil {
		return err
	}

	return nil
}

// disable key rotation on the kms key
func (rm *resourceManager) disableKeyRotation(keyId *string) error {
	disableKeyRotationInput := svcsdk.DisableKeyRotationInput{
		KeyId: keyId,
	}
	_, err := rm.sdkapi.DisableKeyRotation(&disableKeyRotationInput)
	rm.metrics.RecordAPICall("POST", "DisableKeyRotation", err)
	if err != nil {
		return err
	}

	return nil
}
