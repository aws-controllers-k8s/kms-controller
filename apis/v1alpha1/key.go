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

// KeySpec defines the desired state of Key.
type KeySpec struct {

	// Skips ("bypasses") the key policy lockout safety check. The default value
	// is false.
	//
	// Setting this value to true increases the risk that the KMS key becomes unmanageable.
	// Do not set this value to true indiscriminately.
	//
	// For more information, see Default key policy (https://docs.aws.amazon.com/kms/latest/developerguide/key-policy-default.html#prevent-unmanageable-key)
	// in the Key Management Service Developer Guide.
	//
	// Use this parameter only when you intend to prevent the principal that is
	// making the request from making a subsequent PutKeyPolicy (https://docs.aws.amazon.com/kms/latest/APIReference/API_PutKeyPolicy.html)
	// request on the KMS key.
	BypassPolicyLockoutSafetyCheck *bool `json:"bypassPolicyLockoutSafetyCheck,omitempty"`
	// Creates the KMS key in the specified custom key store (https://docs.aws.amazon.com/kms/latest/developerguide/custom-key-store-overview.html).
	// The ConnectionState of the custom key store must be CONNECTED. To find the
	// CustomKeyStoreID and ConnectionState use the DescribeCustomKeyStores operation.
	//
	// This parameter is valid only for symmetric encryption KMS keys in a single
	// Region. You cannot create any other type of KMS key in a custom key store.
	//
	// When you create a KMS key in an CloudHSM key store, KMS generates a non-exportable
	// 256-bit symmetric key in its associated CloudHSM cluster and associates it
	// with the KMS key. When you create a KMS key in an external key store, you
	// must use the XksKeyId parameter to specify an external key that serves as
	// key material for the KMS key.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	CustomKeyStoreID *string `json:"customKeyStoreID,omitempty"`
	// A description of the KMS key. Use a description that helps you decide whether
	// the KMS key is appropriate for a task. The default value is an empty string
	// (no description).
	//
	// Do not include confidential or sensitive information in this field. This
	// field may be displayed in plaintext in CloudTrail logs and other output.
	//
	// To set or change the description after the key is created, use UpdateKeyDescription.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	Description       *string `json:"description,omitempty"`
	EnableKeyRotation *bool   `json:"enableKeyRotation,omitempty"`
	// Specifies the type of KMS key to create. The default value, SYMMETRIC_DEFAULT,
	// creates a KMS key with a 256-bit AES-GCM key that is used for encryption
	// and decryption, except in China Regions, where it creates a 128-bit symmetric
	// key that uses SM4 encryption. For help choosing a key spec for your KMS key,
	// see Choosing a KMS key type (https://docs.aws.amazon.com/kms/latest/developerguide/key-types.html#symm-asymm-choose)
	// in the Key Management Service Developer Guide .
	//
	// The KeySpec determines whether the KMS key contains a symmetric key or an
	// asymmetric key pair. It also determines the algorithms that the KMS key supports.
	// You can't change the KeySpec after the KMS key is created. To further restrict
	// the algorithms that can be used with the KMS key, use a condition key in
	// its key policy or IAM policy. For more information, see kms:EncryptionAlgorithm
	// (https://docs.aws.amazon.com/kms/latest/developerguide/policy-conditions.html#conditions-kms-encryption-algorithm),
	// kms:MacAlgorithm (https://docs.aws.amazon.com/kms/latest/developerguide/policy-conditions.html#conditions-kms-mac-algorithm)
	// or kms:Signing Algorithm (https://docs.aws.amazon.com/kms/latest/developerguide/policy-conditions.html#conditions-kms-signing-algorithm)
	// in the Key Management Service Developer Guide .
	//
	// Amazon Web Services services that are integrated with KMS (http://aws.amazon.com/kms/features/#AWS_Service_Integration)
	// use symmetric encryption KMS keys to protect your data. These services do
	// not support asymmetric KMS keys or HMAC KMS keys.
	//
	// KMS supports the following key specs for KMS keys:
	//
	//   - Symmetric encryption key (default) SYMMETRIC_DEFAULT
	//
	//   - HMAC keys (symmetric) HMAC_224 HMAC_256 HMAC_384 HMAC_512
	//
	//   - Asymmetric RSA key pairs (encryption and decryption -or- signing and
	//     verification) RSA_2048 RSA_3072 RSA_4096
	//
	//   - Asymmetric NIST-recommended elliptic curve key pairs (signing and verification
	//     -or- deriving shared secrets) ECC_NIST_P256 (secp256r1) ECC_NIST_P384
	//     (secp384r1) ECC_NIST_P521 (secp521r1)
	//
	//   - Other asymmetric elliptic curve key pairs (signing and verification)
	//     ECC_SECG_P256K1 (secp256k1), commonly used for cryptocurrencies.
	//
	//   - SM2 key pairs (encryption and decryption -or- signing and verification
	//     -or- deriving shared secrets) SM2 (China Regions only)
	//
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	KeySpec *string `json:"keySpec,omitempty"`
	// Determines the cryptographic operations (https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#cryptographic-operations)
	// for which you can use the KMS key. The default value is ENCRYPT_DECRYPT.
	// This parameter is optional when you are creating a symmetric encryption KMS
	// key; otherwise, it is required. You can't change the KeyUsage value after
	// the KMS key is created.
	//
	// Select only one valid value.
	//
	//   - For symmetric encryption KMS keys, omit the parameter or specify ENCRYPT_DECRYPT.
	//
	//   - For HMAC KMS keys (symmetric), specify GENERATE_VERIFY_MAC.
	//
	//   - For asymmetric KMS keys with RSA key pairs, specify ENCRYPT_DECRYPT
	//     or SIGN_VERIFY.
	//
	//   - For asymmetric KMS keys with NIST-recommended elliptic curve key pairs,
	//     specify SIGN_VERIFY or KEY_AGREEMENT.
	//
	//   - For asymmetric KMS keys with ECC_SECG_P256K1 key pairs specify SIGN_VERIFY.
	//
	//   - For asymmetric KMS keys with SM2 key pairs (China Regions only), specify
	//     ENCRYPT_DECRYPT, SIGN_VERIFY, or KEY_AGREEMENT.
	//
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	KeyUsage *string `json:"keyUsage,omitempty"`
	// Creates a multi-Region primary key that you can replicate into other Amazon
	// Web Services Regions. You cannot change this value after you create the KMS
	// key.
	//
	// For a multi-Region key, set this parameter to True. For a single-Region KMS
	// key, omit this parameter or set it to False. The default value is False.
	//
	// This operation supports multi-Region keys, an KMS feature that lets you create
	// multiple interoperable KMS keys in different Amazon Web Services Regions.
	// Because these KMS keys have the same key ID, key material, and other metadata,
	// you can use them interchangeably to encrypt data in one Amazon Web Services
	// Region and decrypt it in a different Amazon Web Services Region without re-encrypting
	// the data or making a cross-Region call. For more information about multi-Region
	// keys, see Multi-Region keys in KMS (https://docs.aws.amazon.com/kms/latest/developerguide/multi-region-keys-overview.html)
	// in the Key Management Service Developer Guide.
	//
	// This value creates a primary key, not a replica. To create a replica key,
	// use the ReplicateKey operation.
	//
	// You can create a symmetric or asymmetric multi-Region key, and you can create
	// a multi-Region key with imported key material. However, you cannot create
	// a multi-Region key in a custom key store.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	MultiRegion *bool `json:"multiRegion,omitempty"`
	// The source of the key material for the KMS key. You cannot change the origin
	// after you create the KMS key. The default is AWS_KMS, which means that KMS
	// creates the key material.
	//
	// To create a KMS key with no key material (https://docs.aws.amazon.com/kms/latest/developerguide/importing-keys-create-cmk.html)
	// (for imported key material), set this value to EXTERNAL. For more information
	// about importing key material into KMS, see Importing Key Material (https://docs.aws.amazon.com/kms/latest/developerguide/importing-keys.html)
	// in the Key Management Service Developer Guide. The EXTERNAL origin value
	// is valid only for symmetric KMS keys.
	//
	// To create a KMS key in an CloudHSM key store (https://docs.aws.amazon.com/kms/latest/developerguide/create-cmk-keystore.html)
	// and create its key material in the associated CloudHSM cluster, set this
	// value to AWS_CLOUDHSM. You must also use the CustomKeyStoreId parameter to
	// identify the CloudHSM key store. The KeySpec value must be SYMMETRIC_DEFAULT.
	//
	// To create a KMS key in an external key store (https://docs.aws.amazon.com/kms/latest/developerguide/create-xks-keys.html),
	// set this value to EXTERNAL_KEY_STORE. You must also use the CustomKeyStoreId
	// parameter to identify the external key store and the XksKeyId parameter to
	// identify the associated external key. The KeySpec value must be SYMMETRIC_DEFAULT.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	Origin *string `json:"origin,omitempty"`
	// The key policy to attach to the KMS key.
	//
	// If you provide a key policy, it must meet the following criteria:
	//
	//   - The key policy must allow the calling principal to make a subsequent
	//     PutKeyPolicy request on the KMS key. This reduces the risk that the KMS
	//     key becomes unmanageable. For more information, see Default key policy
	//     (https://docs.aws.amazon.com/kms/latest/developerguide/key-policy-default.html#prevent-unmanageable-key)
	//     in the Key Management Service Developer Guide. (To omit this condition,
	//     set BypassPolicyLockoutSafetyCheck to true.)
	//
	//   - Each statement in the key policy must contain one or more principals.
	//     The principals in the key policy must exist and be visible to KMS. When
	//     you create a new Amazon Web Services principal, you might need to enforce
	//     a delay before including the new principal in a key policy because the
	//     new principal might not be immediately visible to KMS. For more information,
	//     see Changes that I make are not always immediately visible (https://docs.aws.amazon.com/IAM/latest/UserGuide/troubleshoot_general.html#troubleshoot_general_eventual-consistency)
	//     in the Amazon Web Services Identity and Access Management User Guide.
	//
	// If you do not provide a key policy, KMS attaches a default key policy to
	// the KMS key. For more information, see Default key policy (https://docs.aws.amazon.com/kms/latest/developerguide/key-policies.html#key-policy-default)
	// in the Key Management Service Developer Guide.
	//
	// The key policy size quota is 32 kilobytes (32768 bytes).
	//
	// For help writing and formatting a JSON policy document, see the IAM JSON
	// Policy Reference (https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies.html)
	// in the Identity and Access Management User Guide .
	//
	// Regex Pattern: `^[\u0009\u000A\u000D\u0020-\u00FF]+$`
	Policy *string `json:"policy,omitempty"`
	// Assigns one or more tags to the KMS key. Use this parameter to tag the KMS
	// key when it is created. To tag an existing KMS key, use the TagResource operation.
	//
	// Do not include confidential or sensitive information in this field. This
	// field may be displayed in plaintext in CloudTrail logs and other output.
	//
	// Tagging or untagging a KMS key can allow or deny permission to the KMS key.
	// For details, see ABAC for KMS (https://docs.aws.amazon.com/kms/latest/developerguide/abac.html)
	// in the Key Management Service Developer Guide.
	//
	// To use this parameter, you must have kms:TagResource (https://docs.aws.amazon.com/kms/latest/developerguide/kms-api-permissions-reference.html)
	// permission in an IAM policy.
	//
	// Each tag consists of a tag key and a tag value. Both the tag key and the
	// tag value are required, but the tag value can be an empty (null) string.
	// You cannot have more than one tag on a KMS key with the same tag key. If
	// you specify an existing tag key with a different tag value, KMS replaces
	// the current tag value with the specified one.
	//
	// When you add tags to an Amazon Web Services resource, Amazon Web Services
	// generates a cost allocation report with usage and costs aggregated by tags.
	// Tags can also be used to control access to a KMS key. For details, see Tagging
	// Keys (https://docs.aws.amazon.com/kms/latest/developerguide/tagging-keys.html).
	Tags []*Tag `json:"tags,omitempty"`
}

// KeyStatus defines the observed state of Key
type KeyStatus struct {
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
	// The twelve-digit account ID of the Amazon Web Services account that owns
	// the KMS key.
	// +kubebuilder:validation:Optional
	AWSAccountID *string `json:"awsAccountID,omitempty"`
	// The cluster ID of the CloudHSM cluster that contains the key material for
	// the KMS key. When you create a KMS key in an CloudHSM custom key store (https://docs.aws.amazon.com/kms/latest/developerguide/custom-key-store-overview.html),
	// KMS creates the key material for the KMS key in the associated CloudHSM cluster.
	// This field is present only when the KMS key is created in an CloudHSM key
	// store.
	//
	// Regex Pattern: `^cluster-[2-7a-zA-Z]{11,16}$`
	// +kubebuilder:validation:Optional
	CloudHsmClusterID *string `json:"cloudHsmClusterID,omitempty"`
	// The date and time when the KMS key was created.
	// +kubebuilder:validation:Optional
	CreationDate *metav1.Time `json:"creationDate,omitempty"`
	// The date and time after which KMS deletes this KMS key. This value is present
	// only when the KMS key is scheduled for deletion, that is, when its KeyState
	// is PendingDeletion.
	//
	// When the primary key in a multi-Region key is scheduled for deletion but
	// still has replica keys, its key state is PendingReplicaDeletion and the length
	// of its waiting period is displayed in the PendingDeletionWindowInDays field.
	// +kubebuilder:validation:Optional
	DeletionDate *metav1.Time `json:"deletionDate,omitempty"`
	// Specifies whether the KMS key is enabled. When KeyState is Enabled this value
	// is true, otherwise it is false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`
	// The encryption algorithms that the KMS key supports. You cannot use the KMS
	// key with other encryption algorithms within KMS.
	//
	// This value is present only when the KeyUsage of the KMS key is ENCRYPT_DECRYPT.
	// +kubebuilder:validation:Optional
	EncryptionAlgorithms []*string `json:"encryptionAlgorithms,omitempty"`
	// Specifies whether the KMS key's key material expires. This value is present
	// only when Origin is EXTERNAL, otherwise this value is omitted.
	// +kubebuilder:validation:Optional
	ExpirationModel *string `json:"expirationModel,omitempty"`
	// The globally unique identifier for the KMS key.
	// +kubebuilder:validation:Optional
	KeyID *string `json:"keyID,omitempty"`
	// The manager of the KMS key. KMS keys in your Amazon Web Services account
	// are either customer managed or Amazon Web Services managed. For more information
	// about the difference, see KMS keys (https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#kms_keys)
	// in the Key Management Service Developer Guide.
	// +kubebuilder:validation:Optional
	KeyManager *string `json:"keyManager,omitempty"`
	// The current status of the KMS key.
	//
	// For more information about how key state affects the use of a KMS key, see
	// Key states of KMS keys (https://docs.aws.amazon.com/kms/latest/developerguide/key-state.html)
	// in the Key Management Service Developer Guide.
	// +kubebuilder:validation:Optional
	KeyState *string `json:"keyState,omitempty"`
	// The message authentication code (MAC) algorithm that the HMAC KMS key supports.
	//
	// This value is present only when the KeyUsage of the KMS key is GENERATE_VERIFY_MAC.
	// +kubebuilder:validation:Optional
	MacAlgorithms []*string `json:"macAlgorithms,omitempty"`
	// Lists the primary and replica keys in same multi-Region key. This field is
	// present only when the value of the MultiRegion field is True.
	//
	// For more information about any listed KMS key, use the DescribeKey operation.
	//
	//    * MultiRegionKeyType indicates whether the KMS key is a PRIMARY or REPLICA
	//    key.
	//
	//    * PrimaryKey displays the key ARN and Region of the primary key. This
	//    field displays the current KMS key if it is the primary key.
	//
	//    * ReplicaKeys displays the key ARNs and Regions of all replica keys. This
	//    field includes the current KMS key if it is a replica key.
	// +kubebuilder:validation:Optional
	MultiRegionConfiguration *MultiRegionConfiguration `json:"multiRegionConfiguration,omitempty"`
	// The waiting period before the primary key in a multi-Region key is deleted.
	// This waiting period begins when the last of its replica keys is deleted.
	// This value is present only when the KeyState of the KMS key is PendingReplicaDeletion.
	// That indicates that the KMS key is the primary key in a multi-Region key,
	// it is scheduled for deletion, and it still has existing replica keys.
	//
	// When a single-Region KMS key or a multi-Region replica key is scheduled for
	// deletion, its deletion date is displayed in the DeletionDate field. However,
	// when the primary key in a multi-Region key is scheduled for deletion, its
	// waiting period doesn't begin until all of its replica keys are deleted. This
	// value displays that waiting period. When the last replica key in the multi-Region
	// key is deleted, the KeyState of the scheduled primary key changes from PendingReplicaDeletion
	// to PendingDeletion and the deletion date appears in the DeletionDate field.
	// +kubebuilder:validation:Optional
	PendingDeletionWindowInDays *int64 `json:"pendingDeletionWindowInDays,omitempty"`
	// The signing algorithms that the KMS key supports. You cannot use the KMS
	// key with other signing algorithms within KMS.
	//
	// This field appears only when the KeyUsage of the KMS key is SIGN_VERIFY.
	// +kubebuilder:validation:Optional
	SigningAlgorithms []*string `json:"signingAlgorithms,omitempty"`
	// The time at which the imported key material expires. When the key material
	// expires, KMS deletes the key material and the KMS key becomes unusable. This
	// value is present only for KMS keys whose Origin is EXTERNAL and whose ExpirationModel
	// is KEY_MATERIAL_EXPIRES, otherwise this value is omitted.
	// +kubebuilder:validation:Optional
	ValidTo *metav1.Time `json:"validTo,omitempty"`
}

// Key is the Schema for the Keys API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Key struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KeySpec   `json:"spec,omitempty"`
	Status            KeyStatus `json:"status,omitempty"`
}

// KeyList contains a list of Key
// +kubebuilder:object:root=true
type KeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Key `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Key{}, &KeyList{})
}
