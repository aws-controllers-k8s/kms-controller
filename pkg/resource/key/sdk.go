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

// Code generated by ack-generate. DO NOT EDIT.

package key

import (
	"context"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/kms"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/kms-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.KMS{}
	_ = &svcapitypes.Key{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer exit(err)
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.DescribeKeyOutput
	resp, err = rm.sdkapi.DescribeKeyWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeKey", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "UNKNOWN" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.KeyMetadata.AWSAccountId != nil {
		ko.Status.AWSAccountID = resp.KeyMetadata.AWSAccountId
	} else {
		ko.Status.AWSAccountID = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.KeyMetadata.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.KeyMetadata.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.KeyMetadata.CloudHsmClusterId != nil {
		ko.Status.CloudHsmClusterID = resp.KeyMetadata.CloudHsmClusterId
	} else {
		ko.Status.CloudHsmClusterID = nil
	}
	if resp.KeyMetadata.CreationDate != nil {
		ko.Status.CreationDate = &metav1.Time{*resp.KeyMetadata.CreationDate}
	} else {
		ko.Status.CreationDate = nil
	}
	if resp.KeyMetadata.CustomKeyStoreId != nil {
		ko.Spec.CustomKeyStoreID = resp.KeyMetadata.CustomKeyStoreId
	} else {
		ko.Spec.CustomKeyStoreID = nil
	}
	if resp.KeyMetadata.DeletionDate != nil {
		ko.Status.DeletionDate = &metav1.Time{*resp.KeyMetadata.DeletionDate}
	} else {
		ko.Status.DeletionDate = nil
	}
	if resp.KeyMetadata.Description != nil {
		ko.Spec.Description = resp.KeyMetadata.Description
	} else {
		ko.Spec.Description = nil
	}
	if resp.KeyMetadata.Enabled != nil {
		ko.Status.Enabled = resp.KeyMetadata.Enabled
	} else {
		ko.Status.Enabled = nil
	}
	if resp.KeyMetadata.EncryptionAlgorithms != nil {
		f8 := []*string{}
		for _, f8iter := range resp.KeyMetadata.EncryptionAlgorithms {
			var f8elem string
			f8elem = *f8iter
			f8 = append(f8, &f8elem)
		}
		ko.Status.EncryptionAlgorithms = f8
	} else {
		ko.Status.EncryptionAlgorithms = nil
	}
	if resp.KeyMetadata.ExpirationModel != nil {
		ko.Status.ExpirationModel = resp.KeyMetadata.ExpirationModel
	} else {
		ko.Status.ExpirationModel = nil
	}
	if resp.KeyMetadata.KeyId != nil {
		ko.Status.KeyID = resp.KeyMetadata.KeyId
	} else {
		ko.Status.KeyID = nil
	}
	if resp.KeyMetadata.KeyManager != nil {
		ko.Status.KeyManager = resp.KeyMetadata.KeyManager
	} else {
		ko.Status.KeyManager = nil
	}
	if resp.KeyMetadata.KeySpec != nil {
		ko.Spec.KeySpec = resp.KeyMetadata.KeySpec
	} else {
		ko.Spec.KeySpec = nil
	}
	if resp.KeyMetadata.KeyState != nil {
		ko.Status.KeyState = resp.KeyMetadata.KeyState
	} else {
		ko.Status.KeyState = nil
	}
	if resp.KeyMetadata.KeyUsage != nil {
		ko.Spec.KeyUsage = resp.KeyMetadata.KeyUsage
	} else {
		ko.Spec.KeyUsage = nil
	}
	if resp.KeyMetadata.MultiRegion != nil {
		ko.Spec.MultiRegion = resp.KeyMetadata.MultiRegion
	} else {
		ko.Spec.MultiRegion = nil
	}
	if resp.KeyMetadata.MultiRegionConfiguration != nil {
		f16 := &svcapitypes.MultiRegionConfiguration{}
		if resp.KeyMetadata.MultiRegionConfiguration.MultiRegionKeyType != nil {
			f16.MultiRegionKeyType = resp.KeyMetadata.MultiRegionConfiguration.MultiRegionKeyType
		}
		if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey != nil {
			f16f1 := &svcapitypes.MultiRegionKey{}
			if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Arn != nil {
				f16f1.ARN = resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Arn
			}
			if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Region != nil {
				f16f1.Region = resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Region
			}
			f16.PrimaryKey = f16f1
		}
		if resp.KeyMetadata.MultiRegionConfiguration.ReplicaKeys != nil {
			f16f2 := []*svcapitypes.MultiRegionKey{}
			for _, f16f2iter := range resp.KeyMetadata.MultiRegionConfiguration.ReplicaKeys {
				f16f2elem := &svcapitypes.MultiRegionKey{}
				if f16f2iter.Arn != nil {
					f16f2elem.ARN = f16f2iter.Arn
				}
				if f16f2iter.Region != nil {
					f16f2elem.Region = f16f2iter.Region
				}
				f16f2 = append(f16f2, f16f2elem)
			}
			f16.ReplicaKeys = f16f2
		}
		ko.Status.MultiRegionConfiguration = f16
	} else {
		ko.Status.MultiRegionConfiguration = nil
	}
	if resp.KeyMetadata.Origin != nil {
		ko.Spec.Origin = resp.KeyMetadata.Origin
	} else {
		ko.Spec.Origin = nil
	}
	if resp.KeyMetadata.PendingDeletionWindowInDays != nil {
		ko.Status.PendingDeletionWindowInDays = resp.KeyMetadata.PendingDeletionWindowInDays
	} else {
		ko.Status.PendingDeletionWindowInDays = nil
	}
	if resp.KeyMetadata.SigningAlgorithms != nil {
		f19 := []*string{}
		for _, f19iter := range resp.KeyMetadata.SigningAlgorithms {
			var f19elem string
			f19elem = *f19iter
			f19 = append(f19, &f19elem)
		}
		ko.Status.SigningAlgorithms = f19
	} else {
		ko.Status.SigningAlgorithms = nil
	}
	if resp.KeyMetadata.ValidTo != nil {
		ko.Status.ValidTo = &metav1.Time{*resp.KeyMetadata.ValidTo}
	} else {
		ko.Status.ValidTo = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Status.KeyID == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeKeyInput, error) {
	res := &svcsdk.DescribeKeyInput{}

	if r.ko.Status.KeyID != nil {
		res.SetKeyId(*r.ko.Status.KeyID)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer exit(err)
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateKeyOutput
	_ = resp
	resp, err = rm.sdkapi.CreateKeyWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateKey", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.KeyMetadata.AWSAccountId != nil {
		ko.Status.AWSAccountID = resp.KeyMetadata.AWSAccountId
	} else {
		ko.Status.AWSAccountID = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.KeyMetadata.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.KeyMetadata.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.KeyMetadata.CloudHsmClusterId != nil {
		ko.Status.CloudHsmClusterID = resp.KeyMetadata.CloudHsmClusterId
	} else {
		ko.Status.CloudHsmClusterID = nil
	}
	if resp.KeyMetadata.CreationDate != nil {
		ko.Status.CreationDate = &metav1.Time{*resp.KeyMetadata.CreationDate}
	} else {
		ko.Status.CreationDate = nil
	}
	if resp.KeyMetadata.CustomKeyStoreId != nil {
		ko.Spec.CustomKeyStoreID = resp.KeyMetadata.CustomKeyStoreId
	} else {
		ko.Spec.CustomKeyStoreID = nil
	}
	if resp.KeyMetadata.DeletionDate != nil {
		ko.Status.DeletionDate = &metav1.Time{*resp.KeyMetadata.DeletionDate}
	} else {
		ko.Status.DeletionDate = nil
	}
	if resp.KeyMetadata.Description != nil {
		ko.Spec.Description = resp.KeyMetadata.Description
	} else {
		ko.Spec.Description = nil
	}
	if resp.KeyMetadata.Enabled != nil {
		ko.Status.Enabled = resp.KeyMetadata.Enabled
	} else {
		ko.Status.Enabled = nil
	}
	if resp.KeyMetadata.EncryptionAlgorithms != nil {
		f8 := []*string{}
		for _, f8iter := range resp.KeyMetadata.EncryptionAlgorithms {
			var f8elem string
			f8elem = *f8iter
			f8 = append(f8, &f8elem)
		}
		ko.Status.EncryptionAlgorithms = f8
	} else {
		ko.Status.EncryptionAlgorithms = nil
	}
	if resp.KeyMetadata.ExpirationModel != nil {
		ko.Status.ExpirationModel = resp.KeyMetadata.ExpirationModel
	} else {
		ko.Status.ExpirationModel = nil
	}
	if resp.KeyMetadata.KeyId != nil {
		ko.Status.KeyID = resp.KeyMetadata.KeyId
	} else {
		ko.Status.KeyID = nil
	}
	if resp.KeyMetadata.KeyManager != nil {
		ko.Status.KeyManager = resp.KeyMetadata.KeyManager
	} else {
		ko.Status.KeyManager = nil
	}
	if resp.KeyMetadata.KeySpec != nil {
		ko.Spec.KeySpec = resp.KeyMetadata.KeySpec
	} else {
		ko.Spec.KeySpec = nil
	}
	if resp.KeyMetadata.KeyState != nil {
		ko.Status.KeyState = resp.KeyMetadata.KeyState
	} else {
		ko.Status.KeyState = nil
	}
	if resp.KeyMetadata.KeyUsage != nil {
		ko.Spec.KeyUsage = resp.KeyMetadata.KeyUsage
	} else {
		ko.Spec.KeyUsage = nil
	}
	if resp.KeyMetadata.MultiRegion != nil {
		ko.Spec.MultiRegion = resp.KeyMetadata.MultiRegion
	} else {
		ko.Spec.MultiRegion = nil
	}
	if resp.KeyMetadata.MultiRegionConfiguration != nil {
		f16 := &svcapitypes.MultiRegionConfiguration{}
		if resp.KeyMetadata.MultiRegionConfiguration.MultiRegionKeyType != nil {
			f16.MultiRegionKeyType = resp.KeyMetadata.MultiRegionConfiguration.MultiRegionKeyType
		}
		if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey != nil {
			f16f1 := &svcapitypes.MultiRegionKey{}
			if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Arn != nil {
				f16f1.ARN = resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Arn
			}
			if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Region != nil {
				f16f1.Region = resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Region
			}
			f16.PrimaryKey = f16f1
		}
		if resp.KeyMetadata.MultiRegionConfiguration.ReplicaKeys != nil {
			f16f2 := []*svcapitypes.MultiRegionKey{}
			for _, f16f2iter := range resp.KeyMetadata.MultiRegionConfiguration.ReplicaKeys {
				f16f2elem := &svcapitypes.MultiRegionKey{}
				if f16f2iter.Arn != nil {
					f16f2elem.ARN = f16f2iter.Arn
				}
				if f16f2iter.Region != nil {
					f16f2elem.Region = f16f2iter.Region
				}
				f16f2 = append(f16f2, f16f2elem)
			}
			f16.ReplicaKeys = f16f2
		}
		ko.Status.MultiRegionConfiguration = f16
	} else {
		ko.Status.MultiRegionConfiguration = nil
	}
	if resp.KeyMetadata.Origin != nil {
		ko.Spec.Origin = resp.KeyMetadata.Origin
	} else {
		ko.Spec.Origin = nil
	}
	if resp.KeyMetadata.PendingDeletionWindowInDays != nil {
		ko.Status.PendingDeletionWindowInDays = resp.KeyMetadata.PendingDeletionWindowInDays
	} else {
		ko.Status.PendingDeletionWindowInDays = nil
	}
	if resp.KeyMetadata.SigningAlgorithms != nil {
		f19 := []*string{}
		for _, f19iter := range resp.KeyMetadata.SigningAlgorithms {
			var f19elem string
			f19elem = *f19iter
			f19 = append(f19, &f19elem)
		}
		ko.Status.SigningAlgorithms = f19
	} else {
		ko.Status.SigningAlgorithms = nil
	}
	if resp.KeyMetadata.ValidTo != nil {
		ko.Status.ValidTo = &metav1.Time{*resp.KeyMetadata.ValidTo}
	} else {
		ko.Status.ValidTo = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateKeyInput, error) {
	res := &svcsdk.CreateKeyInput{}

	if r.ko.Spec.BypassPolicyLockoutSafetyCheck != nil {
		res.SetBypassPolicyLockoutSafetyCheck(*r.ko.Spec.BypassPolicyLockoutSafetyCheck)
	}
	if r.ko.Spec.CustomKeyStoreID != nil {
		res.SetCustomKeyStoreId(*r.ko.Spec.CustomKeyStoreID)
	}
	if r.ko.Spec.Description != nil {
		res.SetDescription(*r.ko.Spec.Description)
	}
	if r.ko.Spec.KeySpec != nil {
		res.SetKeySpec(*r.ko.Spec.KeySpec)
	}
	if r.ko.Spec.KeyUsage != nil {
		res.SetKeyUsage(*r.ko.Spec.KeyUsage)
	}
	if r.ko.Spec.MultiRegion != nil {
		res.SetMultiRegion(*r.ko.Spec.MultiRegion)
	}
	if r.ko.Spec.Origin != nil {
		res.SetOrigin(*r.ko.Spec.Origin)
	}
	if r.ko.Spec.Policy != nil {
		res.SetPolicy(*r.ko.Spec.Policy)
	}
	if r.ko.Spec.Tags != nil {
		f8 := []*svcsdk.Tag{}
		for _, f8iter := range r.ko.Spec.Tags {
			f8elem := &svcsdk.Tag{}
			if f8iter.TagKey != nil {
				f8elem.SetTagKey(*f8iter.TagKey)
			}
			if f8iter.TagValue != nil {
				f8elem.SetTagValue(*f8iter.TagValue)
			}
			f8 = append(f8, f8elem)
		}
		res.SetTags(f8)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	// TODO(jaypipes): Figure this out...
	return nil, ackerr.NotImplemented
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer exit(err)
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	input.SetPendingWindowInDays(GetDeletePendingWindowInDays(&r.ko.ObjectMeta))
	var resp *svcsdk.ScheduleKeyDeletionOutput
	_ = resp
	resp, err = rm.sdkapi.ScheduleKeyDeletionWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "ScheduleKeyDeletion", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.ScheduleKeyDeletionInput, error) {
	res := &svcsdk.ScheduleKeyDeletionInput{}

	if r.ko.Status.KeyID != nil {
		res.SetKeyId(*r.ko.Status.KeyID)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Key,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}

	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	// No terminal_errors specified for this resource in generator config
	return false
}
