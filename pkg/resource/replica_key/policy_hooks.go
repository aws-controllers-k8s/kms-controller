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

	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/kms"
)

var (
	// PolicyName is the only allowed value for KMS Key's PolicyName
	// https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/kms.html#KMS.Client.put_key_policy
	PolicyName = "default"
)

// updatePolicy performs the PutKeyPolicy operation after reading the Policy
// and BypassPolicyLockoutSafetyCheck from resource spec
func (rm *resourceManager) updatePolicy(ctx context.Context, r *resource) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updatePolicy")
	defer func() {
		exit(err)
	}()

	input := &svcsdk.PutKeyPolicyInput{
		BypassPolicyLockoutSafetyCheck: r.ko.Spec.BypassPolicyLockoutSafetyCheck != nil && *r.ko.Spec.BypassPolicyLockoutSafetyCheck,
		KeyId:                          r.ko.Status.ReplicaKeyMetadata.KeyID,
		Policy:                         r.ko.Spec.Policy,
		PolicyName:                     &PolicyName,
	}

	_, err = rm.sdkapi.PutKeyPolicy(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "PutKeyPolicy", err)
	return err
}

// getPolicy performs the GetKeyPolicy API call and returns the key policy
func (rm *resourceManager) getPolicy(ctx context.Context, r *resource) (policy *string, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getPolicy")
	defer func() {
		exit(err)
	}()
	input := &svcsdk.GetKeyPolicyInput{
		KeyId:      r.ko.Status.ReplicaKeyMetadata.KeyID,
		PolicyName: &PolicyName,
	}
	resp, err := rm.sdkapi.GetKeyPolicy(ctx, input)
	rm.metrics.RecordAPICall("GET", "GetKeyPolicy", err)
	if err != nil {
		return nil, err
	}
	return resp.Policy, nil
}
