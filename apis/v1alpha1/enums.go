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

type AlgorithmSpec string

const (
	AlgorithmSpec_RSAES_PKCS1_V1_5   AlgorithmSpec = "RSAES_PKCS1_V1_5"
	AlgorithmSpec_RSAES_OAEP_SHA_1   AlgorithmSpec = "RSAES_OAEP_SHA_1"
	AlgorithmSpec_RSAES_OAEP_SHA_256 AlgorithmSpec = "RSAES_OAEP_SHA_256"
)

type ConnectionErrorCodeType string

const (
	ConnectionErrorCodeType_INVALID_CREDENTIALS        ConnectionErrorCodeType = "INVALID_CREDENTIALS"
	ConnectionErrorCodeType_CLUSTER_NOT_FOUND          ConnectionErrorCodeType = "CLUSTER_NOT_FOUND"
	ConnectionErrorCodeType_NETWORK_ERRORS             ConnectionErrorCodeType = "NETWORK_ERRORS"
	ConnectionErrorCodeType_INTERNAL_ERROR             ConnectionErrorCodeType = "INTERNAL_ERROR"
	ConnectionErrorCodeType_INSUFFICIENT_CLOUDHSM_HSMS ConnectionErrorCodeType = "INSUFFICIENT_CLOUDHSM_HSMS"
	ConnectionErrorCodeType_USER_LOCKED_OUT            ConnectionErrorCodeType = "USER_LOCKED_OUT"
	ConnectionErrorCodeType_USER_NOT_FOUND             ConnectionErrorCodeType = "USER_NOT_FOUND"
	ConnectionErrorCodeType_USER_LOGGED_IN             ConnectionErrorCodeType = "USER_LOGGED_IN"
	ConnectionErrorCodeType_SUBNET_NOT_FOUND           ConnectionErrorCodeType = "SUBNET_NOT_FOUND"
)

type ConnectionStateType string

const (
	ConnectionStateType_CONNECTED     ConnectionStateType = "CONNECTED"
	ConnectionStateType_CONNECTING    ConnectionStateType = "CONNECTING"
	ConnectionStateType_FAILED        ConnectionStateType = "FAILED"
	ConnectionStateType_DISCONNECTED  ConnectionStateType = "DISCONNECTED"
	ConnectionStateType_DISCONNECTING ConnectionStateType = "DISCONNECTING"
)

type CustomerMasterKeySpec string

const (
	CustomerMasterKeySpec_RSA_2048          CustomerMasterKeySpec = "RSA_2048"
	CustomerMasterKeySpec_RSA_3072          CustomerMasterKeySpec = "RSA_3072"
	CustomerMasterKeySpec_RSA_4096          CustomerMasterKeySpec = "RSA_4096"
	CustomerMasterKeySpec_ECC_NIST_P256     CustomerMasterKeySpec = "ECC_NIST_P256"
	CustomerMasterKeySpec_ECC_NIST_P384     CustomerMasterKeySpec = "ECC_NIST_P384"
	CustomerMasterKeySpec_ECC_NIST_P521     CustomerMasterKeySpec = "ECC_NIST_P521"
	CustomerMasterKeySpec_ECC_SECG_P256K1   CustomerMasterKeySpec = "ECC_SECG_P256K1"
	CustomerMasterKeySpec_SYMMETRIC_DEFAULT CustomerMasterKeySpec = "SYMMETRIC_DEFAULT"
)

type DataKeyPairSpec string

const (
	DataKeyPairSpec_RSA_2048        DataKeyPairSpec = "RSA_2048"
	DataKeyPairSpec_RSA_3072        DataKeyPairSpec = "RSA_3072"
	DataKeyPairSpec_RSA_4096        DataKeyPairSpec = "RSA_4096"
	DataKeyPairSpec_ECC_NIST_P256   DataKeyPairSpec = "ECC_NIST_P256"
	DataKeyPairSpec_ECC_NIST_P384   DataKeyPairSpec = "ECC_NIST_P384"
	DataKeyPairSpec_ECC_NIST_P521   DataKeyPairSpec = "ECC_NIST_P521"
	DataKeyPairSpec_ECC_SECG_P256K1 DataKeyPairSpec = "ECC_SECG_P256K1"
)

type DataKeySpec string

const (
	DataKeySpec_AES_256 DataKeySpec = "AES_256"
	DataKeySpec_AES_128 DataKeySpec = "AES_128"
)

type EncryptionAlgorithmSpec string

const (
	EncryptionAlgorithmSpec_SYMMETRIC_DEFAULT  EncryptionAlgorithmSpec = "SYMMETRIC_DEFAULT"
	EncryptionAlgorithmSpec_RSAES_OAEP_SHA_1   EncryptionAlgorithmSpec = "RSAES_OAEP_SHA_1"
	EncryptionAlgorithmSpec_RSAES_OAEP_SHA_256 EncryptionAlgorithmSpec = "RSAES_OAEP_SHA_256"
)

type ExpirationModelType string

const (
	ExpirationModelType_KEY_MATERIAL_EXPIRES         ExpirationModelType = "KEY_MATERIAL_EXPIRES"
	ExpirationModelType_KEY_MATERIAL_DOES_NOT_EXPIRE ExpirationModelType = "KEY_MATERIAL_DOES_NOT_EXPIRE"
)

type GrantOperation string

const (
	GrantOperation_Decrypt                             GrantOperation = "Decrypt"
	GrantOperation_Encrypt                             GrantOperation = "Encrypt"
	GrantOperation_GenerateDataKey                     GrantOperation = "GenerateDataKey"
	GrantOperation_GenerateDataKeyWithoutPlaintext     GrantOperation = "GenerateDataKeyWithoutPlaintext"
	GrantOperation_ReEncryptFrom                       GrantOperation = "ReEncryptFrom"
	GrantOperation_ReEncryptTo                         GrantOperation = "ReEncryptTo"
	GrantOperation_Sign                                GrantOperation = "Sign"
	GrantOperation_Verify                              GrantOperation = "Verify"
	GrantOperation_GetPublicKey                        GrantOperation = "GetPublicKey"
	GrantOperation_CreateGrant                         GrantOperation = "CreateGrant"
	GrantOperation_RetireGrant                         GrantOperation = "RetireGrant"
	GrantOperation_DescribeKey                         GrantOperation = "DescribeKey"
	GrantOperation_GenerateDataKeyPair                 GrantOperation = "GenerateDataKeyPair"
	GrantOperation_GenerateDataKeyPairWithoutPlaintext GrantOperation = "GenerateDataKeyPairWithoutPlaintext"
)

type KeyManagerType string

const (
	KeyManagerType_AWS      KeyManagerType = "AWS"
	KeyManagerType_CUSTOMER KeyManagerType = "CUSTOMER"
)

type KeySpec_SDK string

const (
	KeySpec_SDK_RSA_2048          KeySpec_SDK = "RSA_2048"
	KeySpec_SDK_RSA_3072          KeySpec_SDK = "RSA_3072"
	KeySpec_SDK_RSA_4096          KeySpec_SDK = "RSA_4096"
	KeySpec_SDK_ECC_NIST_P256     KeySpec_SDK = "ECC_NIST_P256"
	KeySpec_SDK_ECC_NIST_P384     KeySpec_SDK = "ECC_NIST_P384"
	KeySpec_SDK_ECC_NIST_P521     KeySpec_SDK = "ECC_NIST_P521"
	KeySpec_SDK_ECC_SECG_P256K1   KeySpec_SDK = "ECC_SECG_P256K1"
	KeySpec_SDK_SYMMETRIC_DEFAULT KeySpec_SDK = "SYMMETRIC_DEFAULT"
)

type KeyState string

const (
	KeyState_Creating               KeyState = "Creating"
	KeyState_Enabled                KeyState = "Enabled"
	KeyState_Disabled               KeyState = "Disabled"
	KeyState_PendingDeletion        KeyState = "PendingDeletion"
	KeyState_PendingImport          KeyState = "PendingImport"
	KeyState_PendingReplicaDeletion KeyState = "PendingReplicaDeletion"
	KeyState_Unavailable            KeyState = "Unavailable"
	KeyState_Updating               KeyState = "Updating"
)

type KeyUsageType string

const (
	KeyUsageType_SIGN_VERIFY     KeyUsageType = "SIGN_VERIFY"
	KeyUsageType_ENCRYPT_DECRYPT KeyUsageType = "ENCRYPT_DECRYPT"
)

type MessageType string

const (
	MessageType_RAW    MessageType = "RAW"
	MessageType_DIGEST MessageType = "DIGEST"
)

type MultiRegionKeyType string

const (
	MultiRegionKeyType_PRIMARY MultiRegionKeyType = "PRIMARY"
	MultiRegionKeyType_REPLICA MultiRegionKeyType = "REPLICA"
)

type OriginType string

const (
	OriginType_AWS_KMS      OriginType = "AWS_KMS"
	OriginType_EXTERNAL     OriginType = "EXTERNAL"
	OriginType_AWS_CLOUDHSM OriginType = "AWS_CLOUDHSM"
)

type SigningAlgorithmSpec string

const (
	SigningAlgorithmSpec_RSASSA_PSS_SHA_256        SigningAlgorithmSpec = "RSASSA_PSS_SHA_256"
	SigningAlgorithmSpec_RSASSA_PSS_SHA_384        SigningAlgorithmSpec = "RSASSA_PSS_SHA_384"
	SigningAlgorithmSpec_RSASSA_PSS_SHA_512        SigningAlgorithmSpec = "RSASSA_PSS_SHA_512"
	SigningAlgorithmSpec_RSASSA_PKCS1_V1_5_SHA_256 SigningAlgorithmSpec = "RSASSA_PKCS1_V1_5_SHA_256"
	SigningAlgorithmSpec_RSASSA_PKCS1_V1_5_SHA_384 SigningAlgorithmSpec = "RSASSA_PKCS1_V1_5_SHA_384"
	SigningAlgorithmSpec_RSASSA_PKCS1_V1_5_SHA_512 SigningAlgorithmSpec = "RSASSA_PKCS1_V1_5_SHA_512"
	SigningAlgorithmSpec_ECDSA_SHA_256             SigningAlgorithmSpec = "ECDSA_SHA_256"
	SigningAlgorithmSpec_ECDSA_SHA_384             SigningAlgorithmSpec = "ECDSA_SHA_384"
	SigningAlgorithmSpec_ECDSA_SHA_512             SigningAlgorithmSpec = "ECDSA_SHA_512"
)

type WrappingKeySpec string

const (
	WrappingKeySpec_RSA_2048 WrappingKeySpec = "RSA_2048"
)
