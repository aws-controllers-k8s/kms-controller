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
	"errors"
	"fmt"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/smithy-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/kms-controller/apis/v1alpha1"
)

// sdkFind returns SDK-specific information about a supplied resource by using
// the DescribeKey API operation. This is implemented in hook.go because the
// code generator cannot auto-generate a read operation for ReplicateKey.
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()

	// Check if we have a KeyID to look up
	// The KeyID should be in Status.ReplicaKeyMetadata.KeyID after creation
	var keyID *string
	if r.ko.Status.ReplicaKeyMetadata != nil && r.ko.Status.ReplicaKeyMetadata.KeyID != nil {
		keyID = r.ko.Status.ReplicaKeyMetadata.KeyID
	} else if r.ko.Spec.KeyID != nil {
		// Fallback to Spec.KeyID if status doesn't have it yet
		keyID = r.ko.Spec.KeyID
	}

	if keyID == nil {
		return nil, ackerr.NotFound
	}

	input := &svcsdk.DescribeKeyInput{
		KeyId: keyID,
	}

	var resp *svcsdk.DescribeKeyOutput
	resp, err = rm.sdkapi.DescribeKey(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeKey", err)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "NotFoundException" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge the information from the DescribeKey response into a copy of
	// the Kubernetes resource
	ko := r.ko.DeepCopy()

	if resp.KeyMetadata != nil {
		f0 := &svcapitypes.KeyMetadata{}
		if resp.KeyMetadata.AWSAccountId != nil {
			f0.AWSAccountID = resp.KeyMetadata.AWSAccountId
		}
		if resp.KeyMetadata.Arn != nil {
			f0.ARN = resp.KeyMetadata.Arn
		}
		if resp.KeyMetadata.CloudHsmClusterId != nil {
			f0.CloudHsmClusterID = resp.KeyMetadata.CloudHsmClusterId
		}
		if resp.KeyMetadata.CreationDate != nil {
			f0.CreationDate = &metav1.Time{*resp.KeyMetadata.CreationDate}
		}
		if resp.KeyMetadata.CustomKeyStoreId != nil {
			f0.CustomKeyStoreID = resp.KeyMetadata.CustomKeyStoreId
		}
		if resp.KeyMetadata.DeletionDate != nil {
			f0.DeletionDate = &metav1.Time{*resp.KeyMetadata.DeletionDate}
		}
		if resp.KeyMetadata.Description != nil {
			f0.Description = resp.KeyMetadata.Description
		}
		f0.Enabled = &resp.KeyMetadata.Enabled
		if resp.KeyMetadata.EncryptionAlgorithms != nil {
			f0f8 := []*string{}
			for _, f0f8iter := range resp.KeyMetadata.EncryptionAlgorithms {
				var f0f8elem *string
				f0f8elem = aws.String(string(f0f8iter))
				f0f8 = append(f0f8, f0f8elem)
			}
			f0.EncryptionAlgorithms = f0f8
		}
		if resp.KeyMetadata.ExpirationModel != "" {
			f0.ExpirationModel = aws.String(string(resp.KeyMetadata.ExpirationModel))
		}
		if resp.KeyMetadata.KeyId != nil {
			f0.KeyID = resp.KeyMetadata.KeyId
		}
		if resp.KeyMetadata.KeyManager != "" {
			f0.KeyManager = aws.String(string(resp.KeyMetadata.KeyManager))
		}
		if resp.KeyMetadata.KeySpec != "" {
			f0.KeySpec = aws.String(string(resp.KeyMetadata.KeySpec))
		}
		if resp.KeyMetadata.KeyState != "" {
			f0.KeyState = aws.String(string(resp.KeyMetadata.KeyState))
		}
		if resp.KeyMetadata.KeyUsage != "" {
			f0.KeyUsage = aws.String(string(resp.KeyMetadata.KeyUsage))
		}
		if resp.KeyMetadata.MacAlgorithms != nil {
			f0f15 := []*string{}
			for _, f0f15iter := range resp.KeyMetadata.MacAlgorithms {
				var f0f15elem *string
				f0f15elem = aws.String(string(f0f15iter))
				f0f15 = append(f0f15, f0f15elem)
			}
			f0.MacAlgorithms = f0f15
		}
		if resp.KeyMetadata.MultiRegion != nil {
			f0.MultiRegion = resp.KeyMetadata.MultiRegion
		}
		if resp.KeyMetadata.MultiRegionConfiguration != nil {
			f0f17 := &svcapitypes.MultiRegionConfiguration{}
			if resp.KeyMetadata.MultiRegionConfiguration.MultiRegionKeyType != "" {
				f0f17.MultiRegionKeyType = aws.String(string(resp.KeyMetadata.MultiRegionConfiguration.MultiRegionKeyType))
			}
			if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey != nil {
				f0f17f1 := &svcapitypes.MultiRegionKey{}
				if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Arn != nil {
					f0f17f1.ARN = resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Arn
				}
				if resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Region != nil {
					f0f17f1.Region = resp.KeyMetadata.MultiRegionConfiguration.PrimaryKey.Region
				}
				f0f17.PrimaryKey = f0f17f1
			}
			if resp.KeyMetadata.MultiRegionConfiguration.ReplicaKeys != nil {
				f0f17f2 := []*svcapitypes.MultiRegionKey{}
				for _, f0f17f2iter := range resp.KeyMetadata.MultiRegionConfiguration.ReplicaKeys {
					f0f17f2elem := &svcapitypes.MultiRegionKey{}
					if f0f17f2iter.Arn != nil {
						f0f17f2elem.ARN = f0f17f2iter.Arn
					}
					if f0f17f2iter.Region != nil {
						f0f17f2elem.Region = f0f17f2iter.Region
					}
					f0f17f2 = append(f0f17f2, f0f17f2elem)
				}
				f0f17.ReplicaKeys = f0f17f2
			}
			f0.MultiRegionConfiguration = f0f17
		}
		if resp.KeyMetadata.Origin != "" {
			f0.Origin = aws.String(string(resp.KeyMetadata.Origin))
		}
		if resp.KeyMetadata.PendingDeletionWindowInDays != nil {
			pendingDeletionWindowInDaysCopy := int64(*resp.KeyMetadata.PendingDeletionWindowInDays)
			f0.PendingDeletionWindowInDays = &pendingDeletionWindowInDaysCopy
		}
		if resp.KeyMetadata.SigningAlgorithms != nil {
			f0f20 := []*string{}
			for _, f0f20iter := range resp.KeyMetadata.SigningAlgorithms {
				var f0f20elem *string
				f0f20elem = aws.String(string(f0f20iter))
				f0f20 = append(f0f20, f0f20elem)
			}
			f0.SigningAlgorithms = f0f20
		}
		if resp.KeyMetadata.ValidTo != nil {
			f0.ValidTo = &metav1.Time{*resp.KeyMetadata.ValidTo}
		}
		ko.Status.ReplicaKeyMetadata = f0
	} else {
		ko.Status.ReplicaKeyMetadata = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

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