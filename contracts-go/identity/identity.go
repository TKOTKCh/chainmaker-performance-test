/*
 Copyright (C) BABEC. All rights reserved.
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	utils "chainmaker.org/chainmaker/contract-utils/address"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

const (
	paramAdminAddress = "adminAddress"
	paramAddress      = "address"
	keyAdminAddress   = "adminAddress"
)

type identity interface {
	// 安装合约
	initContract(adminAddresses []string) protogo.Response
	// 升级合约
	upgradeContract() protogo.Response
	// 添加白名单列表
	addWriteList(addresses []string) protogo.Response
	// 移除白名单列表
	removeWriteList(addresses []string) protogo.Response
	// 修改管理员
	alterAdminAddress(adminAddress []string) protogo.Response
	// 是否在白名单： "true" 在 "false" 不在
	isApprovedUser(addresses []string) protogo.Response
	// 返回自己的地址
	address() protogo.Response
}

var _ identity = (*IdentityContract)(nil)

// IdentityContract contract
type IdentityContract struct {
}

// InitContract install contract func
func (ic *IdentityContract) InitContract() protogo.Response {
	args := sdk.Instance.GetArgs()
	adminAddress := args[paramAdminAddress]
	var adminAddressStr string
	if len(adminAddress) == 0 {
		adminAddressStr, _ = sdk.Instance.Origin()
	} else {
		adminAddressStr = string(adminAddress)
	}
	adminAddresses := strings.Split(adminAddressStr, ",")
	return ic.initContract(adminAddresses)
}

func (ic *IdentityContract) initContract(adminAddresses []string) protogo.Response {
	identityInfo, err := getIdentityMap()
	if err != nil {
		return sdk.Error(fmt.Sprintf("new storeMap of erc20Info failed, err:%s", err))
	}

	adminAddressByte, _ := json.Marshal(adminAddresses)
	err = identityInfo.Set([]string{keyAdminAddress}, adminAddressByte)
	if err != nil {
		return sdk.Error("set admin address of identityInfo failed")
	}
	err = identityInfo.Set([]string{"userCount"}, []byte("0"))
	if err != nil {
		return sdk.Error("set user count of identityInfo failed")
	}
	sdk.Instance.EmitEvent("alterAdminAddress", adminAddresses)
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrade contract func
func (ic *IdentityContract) UpgradeContract() protogo.Response {
	return ic.upgradeContract()
}
func (ic *IdentityContract) upgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// InvokeContract the entry func of invoke contract func
func (ic *IdentityContract) InvokeContract(method string) protogo.Response {
	args := sdk.Instance.GetArgs()
	if len(method) == 0 {
		return sdk.Error("method of param should not be empty")
	}

	switch method {
	case "addWriteList":
		address := args[paramAddress]
		var addresses []string
		if len(address) != 0 {
			addresses = strings.Split(string(address), ",")
		}
		return ic.addWriteList(addresses)
	case "removeWriteList":
		address := args[paramAddress]
		var addresses []string
		if len(address) != 0 {
			addresses = strings.Split(string(address), ",")
		}
		return ic.removeWriteList(addresses)
	case "alterAdminAddress":
		address := args[paramAdminAddress]
		var addresses []string
		if len(address) != 0 {
			addresses = strings.Split(string(address), ",")
		}
		return ic.alterAdminAddress(addresses)
	case "isApprovedUser":
		address := args[paramAddress]
		var addresses []string
		if len(address) != 0 {
			addresses = strings.Split(string(address), ",")
		}
		return ic.isApprovedUser(addresses)
	case "address":
		return ic.address()
	case "callerAddress":
		return ic.callerAddress()
	default:
		return sdk.Error("Invalid method" + method)
	}
}

func (ic *IdentityContract) addWriteList(addresses []string) protogo.Response {
	if len(addresses) == 0 {
		return sdk.Error("address of param should not be empty")
	}
	//if !ic.senderIsAdmin() {
	//	return sdk.Error("sender is not admin")
	//}
	identityInfo, err := getIdentityWriteMap()
	if err != nil {
		return sdk.Error(fmt.Sprintf("new storeMap of identityInfo failed, err:%s", err))
	}
	for _, address := range addresses {
		if !utils.IsValidAddress(address) {
			return sdk.Error(fmt.Sprintf("addWriteList address[%s,%d] format error", address, len(address)))
		}
		_ = identityInfo.Set([]string{address}, []byte("1"))
	}
	//sdk.Instance.EmitEvent("addWriteList", addresses)
	return sdk.Success([]byte("add write list success"))
}

func (ic *IdentityContract) removeWriteList(addresses []string) protogo.Response {
	if len(addresses) == 0 {
		return sdk.Error("address of param should not be empty")
	}
	//if !ic.senderIsAdmin() {
	//	return sdk.Error("sender is not admin")
	//}
	identityInfo, err := getIdentityWriteMap()
	if err != nil {
		return sdk.Error(fmt.Sprintf("new storeMap of identityInfo failed, err:%s", err))
	}
	for _, address := range addresses {
		_ = identityInfo.Del([]string{address})
	}
	//sdk.Instance.EmitEvent("removeWriteList", addresses)
	return sdk.Success([]byte("remove write list success"))
}

func (ic *IdentityContract) address() protogo.Response {
	addr, err := sdk.Instance.Origin()
	if err != nil {
		return sdk.Error(err.Error())
	}
	if len(addr) == 0 {
		return sdk.Error("addr is empty")
	}
	//sdk.Instance.Infof("sender is %s, len is %d", addr, len(addr))
	return sdk.Success([]byte(addr))
}

func (ic *IdentityContract) callerAddress() protogo.Response {
	var param = make(map[string][]byte)
	resp := sdk.Instance.CallContract("identity", "address", param)
	return resp
}
func (ic *IdentityContract) isApprovedUser(addresses []string) protogo.Response {
	if len(addresses) == 0 {
		//sdk.Instance.Warnf("address is empty")
		return sdk.Success([]byte("address is empty"))
	}

	identityInfo, err := getIdentityWriteMap()
	if err != nil {
		//sdk.Instance.Warnf("new storeMap of identityInfo failed, err:%s", err)
		return sdk.Error(fmt.Sprintf("new storeMap of identityInfo failed, err:%s", err))
	}
	flag := true
	for _, addr := range addresses {
		val, err := identityInfo.Get([]string{addr})
		if len(val) == 0 || err != nil {
			flag = false
		}

	}
	if flag {
		return sdk.Success([]byte("true"))
	} else {
		return sdk.Success([]byte("false"))
	}
}
func (ic *IdentityContract) alterAdminAddress(adminAddress []string) protogo.Response {
	if len(adminAddress) == 0 {
		return sdk.Error("adminAddress of param should not be empty")
	}
	if !ic.senderIsAdmin() {
		return sdk.Error("sender is not admin")
	}

	identityInfo, err := getIdentityMap()
	if err != nil {
		return sdk.Error(fmt.Sprintf("new storeMap of identityInfo failed, err:%s", err))
	}
	adminAddressByte, _ := json.Marshal(adminAddress)
	err = identityInfo.Set([]string{keyAdminAddress}, adminAddressByte)
	if err != nil {
		return sdk.Error("alter admin address of identityInfo failed")
	}
	sdk.Instance.EmitEvent("alterAdminAddress", adminAddress)
	return sdk.Success([]byte("OK"))
}

func (ic *IdentityContract) senderIsAdmin() bool {
	sender, _ := sdk.Instance.Origin()
	identityInfo, err := getIdentityMap()
	if err != nil {
		sdk.Instance.Warnf("new storeMap of allowanceInfo failed, err:%s", err)
		return false
	}
	adminAddressByte, err := identityInfo.Get([]string{keyAdminAddress})
	if len(adminAddressByte) == 0 || err != nil {
		sdk.Instance.Warnf("Get totalSupply failed, err:%s", err)
		return false
	}
	var adminAddress []string
	_ = json.Unmarshal(adminAddressByte, &adminAddress)
	for _, addr := range adminAddress {
		if addr == sender {
			return true
		}
	}
	return false
}

func getIdentityMap() (*sdk.StoreMap, error) {
	return sdk.NewStoreMap("identity", 1, crypto.HASH_TYPE_SHA256)
}
func getIdentityWriteMap() (*sdk.StoreMap, error) {
	return sdk.NewStoreMap("identityWriteList", 1, crypto.HASH_TYPE_SHA256)
}

func main() {
	err := sandbox.Start(new(IdentityContract))
	if err != nil {
		log.Fatal(err)
	}
}
