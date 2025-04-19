/*
Copyright (C) BABEC. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	//"log"
	"strconv"

	bccrypto "chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/hash"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
)

const (
	enc_data_func_name        = "enc_data"
	enc_get_data_func_name    = "get_enc_data"
	enc_auth_func_name        = "enc_auth"
	enc_get_auth_func_name    = "get_enc_auth"
	enc_update_auth_func_name = "update_enc_auth"
)

const (
	// ENC_KEY parameter of the data encryption contract -- enc_key
	ENC_KEY = "enc_key"

	// ENC_AUTHED_PERSON parameter of the data encryption contract -- authorized_person
	ENC_AUTHED_PERSON = "authorized_person"

	// ENC_AUTHOR parameter of the data encryption contract -- authorizer
	ENC_AUTHOR = "authorizer"

	// ENC_AUTH_SIGN parameter of the data encryption contract -- auth_sign
	ENC_AUTH_SIGN = "auth_sign"

	// ENC_AUTH_LEVEL parameter of the data encryption contract -- auth_level
	ENC_AUTH_LEVEL = "auth_level"

	// DATA_KEY parameter of the data encryption contract -- data_key
	DATA_KEY = "data_key"

	// DATA_VALUE parameter of the data encryption contract -- data_value
	DATA_VALUE = "data_value"

	enc_data_filed = "enc_data_contract_filed"
)

// EncAuth the encdata contract struct
type EncAuth struct {

	// DataKey 加密数据的key
	DataKey []byte `json:"dataKey"`

	// AuthorizedPerson 被授权人 （证书或者公钥）
	AuthorizedPerson []byte `json:"authorizedPerson"`

	// EncKey 加密后的 KEY
	EncKey []byte `json:"encAESKey"`

	// Authorizer 授权人（证书或者公钥）
	Authorizer []byte `json:"authorizer"`

	// AuthSignature 授权人签名
	AuthSignature []byte `json:"authSignature"`

	// AuthLevel 授权等级
	AuthLevel AuthLevel `json:"authLevel"`
}

// AuthMsg information of the authorized person
type AuthMsg struct {
	// AuthorizedPerson the identity of the authorized person
	AuthorizedPerson []byte `json:"authorizedPerson"`

	// AuthLevel the level of the authorizer
	AuthLevel []byte `json:"authLevel"`
}

// AuthLevel the Auth Level
type AuthLevel int32

const (
	// ROOT 合约创建人的等级，能授权ADMIN权限
	ROOT AuthLevel = iota + 1
	// ADMIN 可以向下授权COMMON
	ADMIN
	// COMMON 无法继续授权
	COMMON
)

// 安装合约时会执行此方法，必须。ChainMaker不允许用户直接调用该方法。
//
//go:wasmexport init_contract
func initContract() {
	ctx := NewSimContext()
	ctx.SuccessResult("Init contract success")
}

// 升级合约时会执行此方法，必须。ChainMaker不允许用户直接调用该方法。
//
//go:wasmexport upgrade
func upgrade() {
	ctx := NewSimContext()
	ctx.SuccessResult("Upgrade contract success")
}

func hashHex(data []byte) (string, error) {
	hashBytes, err := hash.GetByStrType("SHA256", data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hashBytes), nil
}

//go:wasmexport enc_data
func enc_data() {
	ctx := NewSimContext()
	var err error
	dataKey, ok := ctx.Arg(DATA_KEY)
	if ok != SUCCESS {
		err = fmt.Errorf("failed to store encrypted data, err: the parameter [%s] does not exist", DATA_KEY)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	dataValue, ok := ctx.Arg(DATA_VALUE)
	if ok != SUCCESS {
		err = fmt.Errorf("failed to store encrypted data, err: the parameter [%s] does not exist", DATA_VALUE)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	encKey, ok := ctx.Arg(ENC_KEY)
	if ok != SUCCESS {
		err = fmt.Errorf("failed to store encrypted data, err: the parameter [%s] does not exist", ENC_KEY)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	authorizedPerson, ok := ctx.Arg(ENC_AUTHED_PERSON)
	if ok != SUCCESS {
		err = fmt.Errorf("failed to store encrypted data, err: the parameter [%s] does not exist", ENC_AUTHED_PERSON)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	result, err1 := ctx.GetStateByte(string(dataKey), enc_data_filed)
	if err1 != SUCCESS {
		err = fmt.Errorf("failed to store encrypted data")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if len(result) != 0 {
		err = fmt.Errorf("fail to store encrypted data, err: the data key already exist")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	err1 = ctx.PutStateByte(string(dataKey), enc_data_filed, dataValue)
	if err1 != SUCCESS {
		err = fmt.Errorf("failed to store encrypted data")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	encAuth := &EncAuth{
		DataKey:          dataKey,
		AuthorizedPerson: authorizedPerson,
		EncKey:           encKey,
		Authorizer:       authorizedPerson,
		AuthLevel:        ROOT,
	}

	encAuthBytes, err := json.Marshal(encAuth)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}
	apHashHex, err := hashHex(authorizedPerson)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 存储数据
	err1 = ctx.PutStateByte(apHashHex, string(dataKey), encAuthBytes)
	if err1 != SUCCESS {
		err = fmt.Errorf("fail to init contract")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	ctx.SuccessResult("store encrypted data successfully")
}

//go:wasmexport get_enc_data
func get_enc_data() {
	ctx := NewSimContext()
	var err error
	dataKey, ok := ctx.Arg(DATA_KEY)
	if ok != SUCCESS {
		err = fmt.Errorf("failed to get encrypted data, err: the parameter [%s] does not exist", DATA_KEY)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 查询结果
	result, err1 := ctx.GetStateByte(string(dataKey), enc_data_filed)
	if err1 != SUCCESS {
		err = fmt.Errorf("failed to get encrypted data")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	ctx.SuccessResultByte(result)
}

// nolint:gocyclo,revive
//
//go:wasmexport enc_auth
func enc_auth() {
	ctx := NewSimContext()
	var err error
	authPerson, ok := ctx.Arg(ENC_AUTHED_PERSON)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTHED_PERSON)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	authorizer, ok := ctx.Arg(ENC_AUTHOR)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTHOR)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	authSign, ok := ctx.Arg(ENC_AUTH_SIGN)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTH_SIGN)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	authLevel, ok := ctx.Arg(ENC_AUTH_LEVEL)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTH_LEVEL)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	encKey, ok := ctx.Arg(ENC_KEY)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_KEY)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	dataKey, ok := ctx.Arg(DATA_KEY)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", DATA_KEY)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if bytes.Equal(authPerson, authorizer) {
		err = fmt.Errorf("fail to authorize, err: can not authorize for self")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 验证授权者是否在链上
	authorizerHashHex, err := hashHex(authorizer)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 查询结果
	result, err1 := ctx.GetStateByte(authorizerHashHex, string(dataKey))
	if err1 != SUCCESS {
		err = fmt.Errorf("fail to authorize")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if len(result) == 0 {
		err = fmt.Errorf("fail to authorize, err: the authorizer does not exist")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 验证被授权者是否已经存在
	authPersonHashHex, err := hashHex(authPerson)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 查询结果
	result2, err1 := ctx.GetStateByte(authPersonHashHex, string(dataKey))
	if err1 != SUCCESS {
		err = fmt.Errorf("fail to authorize")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if len(result2) != 0 {
		err = fmt.Errorf("fail to authorize, err: the authorized person already exist")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	var authorizerInfo EncAuth
	err = json.Unmarshal(result, &authorizerInfo)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 构建MSG，验签
	msg := &AuthMsg{
		AuthorizedPerson: authPerson,
		AuthLevel:        authLevel,
	}

	var authorizerCert *bcx509.Certificate
	certBlock, rest := pem.Decode(authorizer)
	if certBlock == nil {
		authorizerCert, err = bcx509.ParseCertificate(rest)
		if err != nil {
			err = fmt.Errorf("fail to authorize, err: %s", err.Error())
			ctx.Log(err.Error())
			ctx.ErrorResult(err.Error())
			return
		}
	} else {
		authorizerCert, err = bcx509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			err = fmt.Errorf("fail to authorize, err: %s", err.Error())
			ctx.Log(err.Error())
			ctx.ErrorResult(err.Error())
			return
		}
	}

	// msg 序列化
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 验签
	ok1, err := authorizerCert.PublicKey.VerifyWithOpts(msgBytes, authSign, &bccrypto.SignOpts{
		Hash: bccrypto.HASH_TYPE_SM3,
		UID:  bccrypto.CRYPTO_DEFAULT_UID,
	})

	if !ok1 || err != nil {
		err = fmt.Errorf("fail to authorize, err: invalid signature")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	intAuthLevel, err := strconv.Atoi(string(authLevel))
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	switch AuthLevel(intAuthLevel) {
	case ADMIN:
		if authorizerInfo.AuthLevel != ROOT {
			err = fmt.Errorf("fail to authorize, err: the authorizer has no power to authorize")
			ctx.Log(err.Error())
			ctx.ErrorResult(err.Error())
			return
		}
	case COMMON:
		if authorizerInfo.AuthLevel != ROOT && authorizerInfo.AuthLevel != ADMIN {
			err = fmt.Errorf("fail to authorize, err: the authorizer has no power to authorize")
			ctx.Log(err.Error())
			ctx.ErrorResult(err.Error())
			return
		}
	default:
		err = fmt.Errorf("fail to authorize, err: invalid auth level")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 授权存储上链
	encAuth := &EncAuth{
		AuthorizedPerson: authPerson,
		EncKey:           encKey,
		Authorizer:       authorizer,
		AuthLevel:        AuthLevel(intAuthLevel),
		DataKey:          dataKey,
		AuthSignature:    authSign,
	}

	encAuthBytes, err := json.Marshal(encAuth)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}
	apHashHex, err := hashHex(authPerson)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 存储数据
	err1 = ctx.PutStateByte(apHashHex, string(dataKey), encAuthBytes)
	if err1 != SUCCESS {
		err = fmt.Errorf("fail to init contract")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	ctx.SuccessResult("store auth info successfully")
}

//go:wasmexport get_enc_auth
func get_enc_auth() {
	ctx := NewSimContext()
	var err error
	authorizer, ok := ctx.Arg(ENC_AUTHOR)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to get auth info, err: the parameter [%s] does not exist", ENC_AUTHOR)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	dataKey, ok := ctx.Arg(DATA_KEY)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to get auth info, err: the parameter [%s] does not exist", DATA_KEY)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	authorizerHashHex, err := hashHex(authorizer)
	if err != nil {
		err = fmt.Errorf("fail to get auth info, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 查询结果
	result, err1 := ctx.GetStateByte(authorizerHashHex, string(dataKey))
	if err1 != SUCCESS {
		err = fmt.Errorf("fail to get auth info")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if result == nil {
		err = fmt.Errorf("fail to get auth info, err: the authorized information does not exist")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	var authInfo EncAuth

	err = json.Unmarshal(result, &authInfo)
	if err != nil {
		err = fmt.Errorf("fail to get auth info, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if len(authInfo.EncKey) == 0 {
		ctx.ErrorResult("the encrypted key is empty")
		return
	}

	ctx.SuccessResultByte(authInfo.EncKey)
}

// nolint:gocyclo,revive
//
//go:wasmexport update_enc_auth
func update_enc_auth() {
	ctx := NewSimContext()
	var err error
	authPerson, ok := ctx.Arg(ENC_AUTHED_PERSON)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTHED_PERSON)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	authorizer, ok := ctx.Arg(ENC_AUTHOR)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTHOR)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	authSign, ok := ctx.Arg(ENC_AUTH_SIGN)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTH_SIGN)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	authLevel, ok := ctx.Arg(ENC_AUTH_LEVEL)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", ENC_AUTH_LEVEL)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	dataKey, ok := ctx.Arg(DATA_KEY)
	if ok != SUCCESS {
		err = fmt.Errorf("fail to authorize, err: the parameter [%s] does not exist", DATA_KEY)
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if bytes.Equal(authPerson, authorizer) {
		err = fmt.Errorf("fail to authorize, err: can not authorize for self")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 验证授权者是否在链上
	authorizerHashHex, err := hashHex(authorizer)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 查询结果
	result, err1 := ctx.GetStateByte(authorizerHashHex, string(dataKey))
	if err1 != SUCCESS {
		err = fmt.Errorf("fail to authorize")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if len(result) == 0 {
		err = fmt.Errorf("fail to authorize, err: the authorizer does not exist")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 验证被授权者是否已经存在
	authPersonHashHex, err := hashHex(authPerson)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 查询结果
	result2, err1 := ctx.GetStateByte(authPersonHashHex, string(dataKey))
	if err1 != SUCCESS {
		err = fmt.Errorf("fail to authorize")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if len(result2) == 0 {
		err = fmt.Errorf("fail to authorize, err: the authorized person does not exist")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	var authorizerInfo EncAuth
	err = json.Unmarshal(result, &authorizerInfo)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	var authPersonInfo EncAuth
	err = json.Unmarshal(result2, &authPersonInfo)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 构建MSG，验签
	msg := &AuthMsg{
		AuthorizedPerson: authPerson,
		AuthLevel:        authLevel,
	}

	var authorizerCert *bcx509.Certificate
	certBlock, rest := pem.Decode(authorizer)
	if certBlock == nil {
		authorizerCert, err = bcx509.ParseCertificate(rest)
		if err != nil {
			err = fmt.Errorf("fail to authorize, err: %s", err.Error())
			ctx.Log(err.Error())
			ctx.ErrorResult(err.Error())
			return
		}
	} else {
		authorizerCert, err = bcx509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			err = fmt.Errorf("fail to authorize, err: %s", err.Error())
			ctx.Log(err.Error())
			ctx.ErrorResult(err.Error())
			return
		}
	}

	// msg 序列化
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 验签
	ok1, err := authorizerCert.PublicKey.VerifyWithOpts(msgBytes, authSign, &bccrypto.SignOpts{
		Hash: bccrypto.HASH_TYPE_SM3,
		UID:  bccrypto.CRYPTO_DEFAULT_UID,
	})

	if !ok1 || err != nil {
		err = fmt.Errorf("fail to authorize, err: invalid signature")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	intAuthLevel, err := strconv.Atoi(string(authLevel))
	if err != nil {
		err = fmt.Errorf("fail to authorize, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if authorizerInfo.AuthLevel != ROOT {
		err = fmt.Errorf("fail to authorize, only the root can change the user permissions")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if authPersonInfo.AuthLevel == ROOT {
		err = fmt.Errorf("fail to authorize, can't update the user permissions who is the root")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	if AuthLevel(intAuthLevel) != ADMIN && AuthLevel(intAuthLevel) != COMMON {
		err = fmt.Errorf("fail to authorize, err: invalid auth level")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 授权存储上链
	encAuth := &EncAuth{
		AuthorizedPerson: authPerson,
		EncKey:           authPersonInfo.EncKey,
		Authorizer:       authorizer,
		AuthLevel:        AuthLevel(intAuthLevel),
		DataKey:          dataKey,
		AuthSignature:    authSign,
	}

	encAuthBytes, err := json.Marshal(encAuth)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}
	apHashHex, err := hashHex(authPerson)
	if err != nil {
		err = fmt.Errorf("fail to init contract, err: %s", err.Error())
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	// 存储数据
	err1 = ctx.PutStateByte(apHashHex, string(dataKey), encAuthBytes)
	if err1 != SUCCESS {
		err = fmt.Errorf("fail to init contract")
		ctx.Log(err.Error())
		ctx.ErrorResult(err.Error())
		return
	}

	ctx.SuccessResult("store auth info successfully")
}

func main() {}
