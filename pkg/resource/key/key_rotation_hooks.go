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

	if r.ko.Spec.EnableKeyRotation == nil {
		return nil
	}

	keyRotationStatus, err := rm.getKeyRotationStatus(r.ko.Status.KeyID)
	if err != nil {
		return err
	}

	// Check if desired state and actual state has any difference
	if keyRotationStatus != nil &&
		keyRotationStatus.KeyRotationEnabled != nil &&
		*keyRotationStatus.KeyRotationEnabled != *r.ko.Spec.EnableKeyRotation {
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
