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

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GrantSpec defines the desired state of Grant.
type GrantSpec struct {

	// Specifies a grant constraint.
	//
	// Do not include confidential or sensitive information in this field. This
	// field may be displayed in plaintext in CloudTrail logs and other output.
	//
	// KMS supports the EncryptionContextEquals and EncryptionContextSubset grant
	// constraints, which allow the permissions in the grant only when the encryption
	// context in the request matches (EncryptionContextEquals) or includes (EncryptionContextSubset)
	// the encryption context specified in the constraint.
	//
	// The encryption context grant constraints are supported only on grant operations
	// (https://docs.aws.amazon.com/kms/latest/developerguide/grants.html#terms-grant-operations)
	// that include an EncryptionContext parameter, such as cryptographic operations
	// on symmetric encryption KMS keys. Grants with grant constraints can include
	// the DescribeKey and RetireGrant operations, but the constraint doesn't apply
	// to these operations. If a grant with a grant constraint includes the CreateGrant
	// operation, the constraint requires that any grants created with the CreateGrant
	// permission have an equally strict or stricter encryption context constraint.
	//
	// You cannot use an encryption context grant constraint for cryptographic operations
	// with asymmetric KMS keys or HMAC KMS keys. Operations with these keys don't
	// support an encryption context.
	//
	// Each constraint value can include up to 8 encryption context pairs. The encryption
	// context value in each constraint cannot exceed 384 characters. For information
	// about grant constraints, see Using grant constraints (https://docs.aws.amazon.com/kms/latest/developerguide/create-grant-overview.html#grant-constraints)
	// in the Key Management Service Developer Guide. For more information about
	// encryption context, see Encryption context (https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#encrypt_context)
	// in the Key Management Service Developer Guide .
	Constraints *GrantConstraints `json:"constraints,omitempty"`
	// A list of grant tokens.
	//
	// Use a grant token when your permission to call this operation comes from
	// a new grant that has not yet achieved eventual consistency. For more information,
	// see Grant token (https://docs.aws.amazon.com/kms/latest/developerguide/grants.html#grant_token)
	// and Using a grant token (https://docs.aws.amazon.com/kms/latest/developerguide/grant-manage.html#using-grant-token)
	// in the Key Management Service Developer Guide.
	GrantTokens []*string `json:"grantTokens,omitempty"`
	// The identity that gets the permissions specified in the grant.
	//
	// To specify the grantee principal, use the Amazon Resource Name (ARN) of an
	// Amazon Web Services principal. Valid principals include Amazon Web Services
	// accounts, IAM users, IAM roles, federated users, and assumed role users.
	// For help with the ARN syntax for a principal, see IAM ARNs (https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-arns)
	// in the Identity and Access Management User Guide .
	//
	// Regex Pattern: `^[\w+=,.@:/-]+$`
	// +kubebuilder:validation:Required
	GranteePrincipal *string `json:"granteePrincipal"`
	// Identifies the KMS key for the grant. The grant gives principals permission
	// to use this KMS key.
	//
	// Specify the key ID or key ARN of the KMS key. To specify a KMS key in adifferent
	// Amazon Web Services account, you must use the key ARN.
	//
	// For example:
	//
	//   - Key ID: 1234abcd-12ab-34cd-56ef-1234567890ab
	//
	//   - Key ARN: arn:aws:kms:us-east-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab
	//
	// To get the key ID and key ARN for a KMS key, use ListKeys or DescribeKey.
	KeyID  *string                                  `json:"keyID,omitempty"`
	KeyRef *ackv1alpha1.AWSResourceReferenceWrapper `json:"keyRef,omitempty"`
	// A friendly name for the grant. Use this value to prevent the unintended creation
	// of duplicate grants when retrying this request.
	//
	// Do not include confidential or sensitive information in this field. This
	// field may be displayed in plaintext in CloudTrail logs and other output.
	//
	// When this value is absent, all CreateGrant requests result in a new grant
	// with a unique GrantId even if all the supplied parameters are identical.
	// This can result in unintended duplicates when you retry the CreateGrant request.
	//
	// When this value is present, you can retry a CreateGrant request with identical
	// parameters; if the grant already exists, the original GrantId is returned
	// without creating a new grant. Note that the returned grant token is unique
	// with every CreateGrant request, even when a duplicate GrantId is returned.
	// All grant tokens for the same grant ID can be used interchangeably.
	//
	// Regex Pattern: `^[a-zA-Z0-9:/_-]+$`
	Name *string `json:"name,omitempty"`
	// A list of operations that the grant permits.
	//
	// This list must include only operations that are permitted in a grant. Also,
	// the operation must be supported on the KMS key. For example, you cannot create
	// a grant for a symmetric encryption KMS key that allows the Sign operation,
	// or a grant for an asymmetric KMS key that allows the GenerateDataKey operation.
	// If you try, KMS returns a ValidationError exception. For details, see Grant
	// operations (https://docs.aws.amazon.com/kms/latest/developerguide/grants.html#terms-grant-operations)
	// in the Key Management Service Developer Guide.
	// +kubebuilder:validation:Required
	Operations []*string `json:"operations"`
	// The principal that has permission to use the RetireGrant operation to retire
	// the grant.
	//
	// To specify the principal, use the Amazon Resource Name (ARN) (https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html)
	// of an Amazon Web Services principal. Valid principals include Amazon Web
	// Services accounts, IAM users, IAM roles, federated users, and assumed role
	// users. For help with the ARN syntax for a principal, see IAM ARNs (https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-arns)
	// in the Identity and Access Management User Guide .
	//
	// The grant determines the retiring principal. Other principals might have
	// permission to retire the grant or revoke the grant. For details, see RevokeGrant
	// and Retiring and revoking grants (https://docs.aws.amazon.com/kms/latest/developerguide/grant-manage.html#grant-delete)
	// in the Key Management Service Developer Guide.
	//
	// Regex Pattern: `^[\w+=,.@:/-]+$`
	RetiringPrincipal *string `json:"retiringPrincipal,omitempty"`
}

// GrantStatus defines the observed state of Grant
type GrantStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRs managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// The unique identifier for the grant.
	//
	// You can use the GrantId in a ListGrants, RetireGrant, or RevokeGrant operation.
	// +kubebuilder:validation:Optional
	GrantID *string `json:"grantID,omitempty"`
	// The grant token.
	//
	// Use a grant token when your permission to call this operation comes from
	// a new grant that has not yet achieved eventual consistency. For more information,
	// see Grant token (https://docs.aws.amazon.com/kms/latest/developerguide/grants.html#grant_token)
	// and Using a grant token (https://docs.aws.amazon.com/kms/latest/developerguide/grant-manage.html#using-grant-token)
	// in the Key Management Service Developer Guide.
	// +kubebuilder:validation:Optional
	GrantToken *string `json:"grantToken,omitempty"`
}

// Grant is the Schema for the Grants API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Grant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GrantSpec   `json:"spec,omitempty"`
	Status            GrantStatus `json:"status,omitempty"`
}

// GrantList contains a list of Grant
// +kubebuilder:object:root=true
type GrantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Grant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Grant{}, &GrantList{})
}
