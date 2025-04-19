/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
参照CMID合约标准实现：
https://git.chainmaker.org.cn/contracts/standard/-/blob/master/draft/CM-CS-221221-Identity.md
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//"log"
	"strconv"

	"chainmaker.org/chainmaker/contract-utils/standard"

	//"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	//"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	//"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/address"
	"chainmaker.org/chainmaker/contract-utils/str"
)

const (
	paramAdminAddress = "adminAddress"
	paramAddress      = "address"
	paramIdentities   = "identities"
	paramLevel        = "level"
	paramPkPem        = "pkPem"
	paramMetadata     = "metadata"
	paramStandardName = "standardName"

	methodPreIdentities        = "[identities]"
	methodPreSetIdentity       = "[setIdentity]"
	methodPreSetIdentityBatch  = "[setIdentityBatch]"
	methodPreIdentityOf        = "[identityOf]"
	methodPreLevelOf           = "[levelOf]"
	methodPrePkPemOf           = "[pkPemOf]"
	methodPreAlterAdminAddress = "[alterAdminAddress]"
	methodPreStandards         = "[standards]"
	methodPreSupportsStandard  = "[supportsStandard]"

	keyAdminAddress    = "a"
	keyAddressLevel    = "l"
	keyAddressIdentity = "i"
	keyAddressPkPem    = "p"
)
const (
	// LevelIllegalAccount level 权限编号示例
	LevelIllegalAccount    int = iota // nolint 0 非法
	LevelPlatformAccount              // nolint 1 平台的用户
	LevelPhoneAccount                 // nolint 2 个人手机注册用户
	LevelPersonalAccount              // nolint 3 个人实名用户
	LevelEnterpriseAccount            // nolint 4 企业实名用户
	LevelMax
)

// 安装合约时会执行此方法，必须。ChainMaker不允许用户直接调用该方法。
//
//go:wasmexport init_contract
func initContract() {
	// 记录异常结果日志
	ctx := NewSimContext()
	sender, err := ctx.GetSenderPk()
	if err != SUCCESS {
		ctx.ErrorResult("get sender fail")
		return
	}

	err1 := AlterAdminAddress(sender)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("Alter admin address fail,%s", err1.Error()))
		return
	}

	ctx.SuccessResult("Init contract success")
}

// 升级合约时会执行此方法，必须。ChainMaker不允许用户直接调用该方法。
//
//go:wasmexport upgrade
func upgrade() {
	ctx := NewSimContext()
	ctx.SuccessResult("Upgrade contract success")
}

//go:wasmexport Identities
func identitiesCore() {
	ctx := NewSimContext()
	identityMetas := Identities()
	data, err := json.Marshal(identityMetas)
	if err != nil {
		func_error(methodPreIdentities, err.Error())
		return
	}
	ctx.SuccessResultByte(data)
}

//go:wasmexport SetIdentity
func setIdentityCore() {
	ctx := NewSimContext()

	// 获取参数
	addressParam, _ := ctx.ArgString(paramAddress)
	levelParam, _ := ctx.ArgString(paramLevel)
	pkPemParam, _ := ctx.ArgString(paramPkPem)
	metadataParam, _ := ctx.ArgString(paramMetadata)

	// 检查参数
	if !address.IsValidAddress(addressParam) {
		func_error(methodPreSetIdentity, "address format error")
		return
	}
	if str.IsAnyBlank(addressParam, levelParam) {
		func_error(methodPreSetIdentity, "address or level of param is empty")
		return
	}
	intLevel, err := strconv.Atoi(levelParam)
	if err != nil {
		func_error(methodPreSetIdentity, err.Error())
		return
	}
	if intLevel < 0 || intLevel >= LevelMax {
		func_error(methodPreSetIdentity, "level of param is illegal, level="+levelParam)
		return
	}

	// 执行逻辑
	err = SetIdentity(addressParam, pkPemParam, intLevel, metadataParam)

	// 返回响应
	if err != nil {
		func_error(methodPreSetIdentity, err.Error())
		return
	}
	ctx.SuccessResult("ok")
}

//go:wasmexport SetIdentityBatch
func setIdentityBatchCore() {
	ctx := NewSimContext()

	// 获取参数
	identitiesParam, _ := ctx.Arg(paramIdentities)
	if str.IsAnyBlank(identitiesParam) {
		func_error(methodPreSetIdentityBatch, "identities of param  is empty")
		return
	}

	// 解析参数
	identities := make([]standard.Identity, 0)
	err := json.Unmarshal(identitiesParam, &identities)
	if err != nil {
		func_error(methodPreSetIdentityBatch, err.Error())
		return
	}

	// 执行逻辑
	err = SetIdentityBatch(identities)

	// 返回响应
	if err != nil {
		func_error(methodPreSetIdentityBatch, err.Error())
		return
	}
	ctx.SuccessResult("ok")
}

//go:wasmexport LevelOf
func levelOfCore() {
	ctx := NewSimContext()

	// 获取、校验参数
	addressParam, _ := ctx.ArgString(paramAddress)
	if str.IsAnyBlank(addressParam) {
		func_error(methodPreLevelOf, "address of param is empty")
		return
	}

	// 查询level
	level, err := LevelOf(addressParam)
	if err != nil {
		func_error(methodPreLevelOf, err.Error())
		return
	}

	ctx.SuccessResult(strconv.Itoa(level))
}

//go:wasmexport IdentityOf
func identityOfCore() {
	ctx := NewSimContext()

	// 获取、校验参数
	addressStr, _ := ctx.ArgString(paramAddress)
	if str.IsAnyBlank(addressStr) {
		func_error(methodPreIdentityOf, "address of param is empty")
		return
	}

	// 查询level
	identity, err := IdentityOf(addressStr)
	if err != nil {
		func_error(methodPreIdentityOf, err.Error())
		return
	}

	identityBytes, err := json.Marshal(identity)
	if err != nil {
		func_error(methodPreIdentityOf, err.Error())
		return
	}
	ctx.SuccessResult(string(identityBytes))
}

//go:wasmexport PkPemOf
func pkPemOfCore() {
	ctx := NewSimContext()

	// 获取、校验参数
	addressParam, _ := ctx.ArgString(paramAddress)
	if str.IsAnyBlank(addressParam) {
		func_error(methodPrePkPemOf, "address of param is empty")
		return
	}

	// 查询level
	pkPem, err := PkPemOf(addressParam)
	if err != nil {
		func_error(methodPrePkPemOf, err.Error())
		return
	}

	ctx.SuccessResult(pkPem)
}

//go:wasmexport AlterAdminAddress
func alterAdminAddressCore() {
	ctx := NewSimContext()

	// 获取参数
	adminAddressParam, _ := ctx.ArgString(paramAdminAddress)

	// 解析参数
	if str.IsAnyBlank(adminAddressParam) {
		func_error(methodPreAlterAdminAddress, "adminAddress of param is empty")
		return
	}
	if !address.IsValidAddress(adminAddressParam) {
		func_error(methodPreAlterAdminAddress, "address format error")
		return
	}

	// 保存admin地址
	err := AlterAdminAddress(adminAddressParam)

	// 返回响应
	if err != nil {
		func_error(methodPreAlterAdminAddress, err.Error())
		return
	}
	ctx.SuccessResult("ok")
}

//go:wasmexport Standards
func standardsCore() {
	ctx := NewSimContext()
	data, err := json.Marshal(Standards())
	if err != nil {
		func_error(methodPreStandards, err.Error())
		return
	}
	ctx.SuccessResult(string(data))
}

//go:wasmexport SupportStandard
func supportStandardCore() {
	ctx := NewSimContext()

	standardName, _ := ctx.ArgString(paramStandardName)
	if str.IsAnyBlank(standardName) {
		func_error(methodPreSupportsStandard, "standardName is empty")
		return
	}

	if SupportStandard(standardName) {
		ctx.SuccessResult(standard.TrueString)
		return
	}
	ctx.SuccessResult(standard.FalseString)
}

// Identities 获取该合约支持的所有认证类型
func Identities() (metas []standard.IdentityMeta) {
	metas = make([]standard.IdentityMeta, 5)
	metas[0] = standard.IdentityMeta{Level: LevelIllegalAccount, Description: "未认证"}
	metas[1] = standard.IdentityMeta{Level: LevelPhoneAccount, Description: "个人手机号注册用户"}
	metas[2] = standard.IdentityMeta{Level: LevelPersonalAccount, Description: "个人实名用户"}
	metas[3] = standard.IdentityMeta{Level: LevelPlatformAccount, Description: "企业的用户"}
	metas[4] = standard.IdentityMeta{Level: LevelEnterpriseAccount, Description: "企业实名用户"}
	return metas
}

// SetIdentity 为地址设置认证类型，管理员可调用
func SetIdentity(address, pkPem string, level int, metadata string) error {
	ctx := NewSimContext()
	if !senderIsAdmin() {
		sender, _ := ctx.GetSenderPk()
		return errors.New("sender is not admin " + sender)
	}

	identity := standard.Identity{
		Address:  address,
		PkPem:    pkPem,
		Level:    level,
		Metadata: metadata,
	}
	identityBytes, err := json.Marshal(identity)
	if err != nil {
		return err
	}

	err1 := ctx.PutStateByte(keyAddressIdentity, address, identityBytes)
	if err1 != SUCCESS {
		return errors.New("func SetIdentity putstate keyAddressIdentity failed")
	}
	err1 = ctx.PutState(keyAddressLevel, address, strconv.Itoa(level))
	if err1 != SUCCESS {
		return errors.New("func SetIdentity putstate keyAddressLevel failed")
	}
	if len(pkPem) > 0 {
		err1 = ctx.PutState(keyAddressPkPem, address, pkPem)
		if err1 != SUCCESS {
			return errors.New("func SetIdentity putstate keyAddressPkPem failed")
		}
	}

	// 你可以使用 metadata 在下方执行自身业务逻辑
	EmitSetIdentityEvent(address, pkPem, level)

	return nil
}

// EmitSetIdentityEvent 发送设置认证类型事件
func EmitSetIdentityEvent(address, pkPem string, level int) {
	ctx := NewSimContext()
	ctx.EmitEvent("setIdentity", address, strconv.Itoa(level), pkPem)
}

// SetIdentityBatch 设置多个认证类型，管理员可调用
func SetIdentityBatch(identities []standard.Identity) error {
	for j := range identities {
		id := identities[j]
		if !address.IsValidAddress(id.Address) {
			return errors.New("address format error")
		}
		if str.IsAnyBlank(id.PkPem) {
			return errors.New("address or level of param is empty")
		}
		if id.Level < 0 || id.Level >= LevelMax {
			return fmt.Errorf("level of param is illegal, level=%d", id.Level)
		}

		err := SetIdentity(id.Address, id.PkPem, id.Level, id.Metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

// IdentityOf 获取认证信息
func IdentityOf(address string) (identity standard.Identity, err error) {
	ctx := NewSimContext()
	identityByte, err1 := ctx.GetStateByte(keyAddressIdentity, address)
	if err1 != SUCCESS {
		return identity, errors.New("func IdentityOf GetStateByte keyAddressIdentity failed")
	}
	if str.IsAnyBlank(identityByte) {
		return identity, errors.New("not found")
	}
	err = json.Unmarshal(identityByte, &identity)
	return identity, err
}

// LevelOf 获取认证编号
func LevelOf(address string) (int, error) {
	ctx := NewSimContext()
	level, err1 := ctx.GetState(keyAddressLevel, address)
	if err1 != SUCCESS {
		return 0, errors.New("func LevelOf GetState keyAddressLevel failed")
	}
	if str.IsAnyBlank(level) {
		return 0, errors.New("not found")
	}
	return strconv.Atoi(level)
}

// PkPemOf 获取公钥
func PkPemOf(address string) (string, error) {
	ctx := NewSimContext()
	pkPem, err1 := ctx.GetState(keyAddressPkPem, address)
	if err1 != SUCCESS {
		return "", errors.New("func PkPemOf GetState keyAddressPkPem failed")
	}

	if str.IsAnyBlank(pkPem) {
		return "0", errors.New("not found")
	}
	return pkPem, nil
}

// AlterAdminAddress 修改管理员，管理员可调用
func AlterAdminAddress(adminAddress string) error {
	ctx := NewSimContext()
	if !senderIsAdmin() {
		sender, _ := ctx.GetSenderPk()
		return errors.New("sender is not admin " + sender)
	}

	adminAddressByte, err := json.Marshal([]string{adminAddress})
	if err != nil {
		return err
	}

	err1 := ctx.PutStateFromKeyByte(keyAdminAddress, adminAddressByte)
	if err1 != SUCCESS {
		return errors.New("alter admin address of identityInfo failed")
	}

	ctx.EmitEvent("alterAdminAddress", adminAddress)

	return nil
}
func senderIsAdmin() bool {
	ctx := NewSimContext()
	sender, _ := ctx.GetSenderPk()
	adminAddressByte, err1 := ctx.GetStateFromKey(keyAdminAddress)

	if err1 != SUCCESS {
		ctx.Log("Get totalSupply failed")
		return false
	}

	if len(adminAddressByte) == 0 {
		return true
	}
	var adminAddress []string
	_ = json.Unmarshal(adminAddressByte, &adminAddress)
	for j := range adminAddress {
		if adminAddress[j] == sender {
			return true
		}
	}
	return false
}

// Standards  获取当前合约支持的标准协议列表
func Standards() []string {
	return []string{standard.ContractStandardNameCMID, standard.ContractStandardNameCMBC}
}

// SupportStandard  获取当前合约是否支持某合约标准协议
func SupportStandard(standardName string) bool {
	return standardName == standard.ContractStandardNameCMID || standardName == standard.ContractStandardNameCMBC
}

func func_error(methodPre, error string) {
	ctx := NewSimContext()
	ctx.ErrorResult(methodPre + "method invoke fail, error: " + error)
}

func main() {}
