/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	utils "chainmaker.org/chainmaker/contract-utils/address"
)

const (
	paramAdminAddress = "adminAddress"
	paramAddress      = "address"
	keyAdminAddress   = "adminAddress"
)

// 安装合约时会执行此方法
//
//go:wasmexport init_contract
func InitContract() {
	ctx := NewSimContext()
	// 获取参数
	adminAddress, err := ctx.ArgString(paramAdminAddress)
	var adminAddressStr string
	if err != SUCCESS || len(adminAddress) == 0 {
		adminAddressStr, _ = ctx.GetSenderPk()
	}
	adminAddresses := strings.Split(adminAddressStr, ",")
	initContract(adminAddresses)
}

// 升级合约时会执行此方法
//
//go:wasmexport upgrade
func Upgrade() {
	ctx := NewSimContext()
	ctx.SuccessResult("Upgrade contract success")
}

// 这个函数不暴露给wasmer，InitContract暴露
func initContract(adminAddresses []string) {
	ctx := NewSimContext()
	adminAddressBytes, _ := json.Marshal(adminAddresses)
	err1 := ctx.PutStateByte("identity", keyAdminAddress, adminAddressBytes)
	if err1 != SUCCESS {
		ctx.ErrorResult("set admin address of identityInfo failed")
		return
	}
	err2 := ctx.PutState("identity", "userCount", "0") //这个userCount没什么作用，只是为了和go的identity合约保持一致
	if err2 != SUCCESS {
		ctx.ErrorResult("set user count of identityInfo failed")
		return
	}
	ctx.EmitEvent("alterAdminAddress", "")
	ctx.SuccessResult("Init contract success")
	return
}

// 由于wasmer虚拟机的写法是直接调用对应函数，不像dockergo那样从一个总的invokeContract中调
// 这里与rust保持一致，不写invokecontract而是直接调，
//
//go:wasmexport addWriteList
func addWriteList() {
	ctx := NewSimContext()
	paramAddress, _ := ctx.ArgString(paramAddress)
	var addresses []string
	if len(paramAddress) != 0 {
		addresses = strings.Split(paramAddress, ",")
	}
	if len(addresses) == 0 {
		ctx.ErrorResult("address of param should not be empty")
		return
	}
	for _, address := range addresses {
		if !utils.IsValidAddress(address) {
			ctx.ErrorResult(fmt.Sprintf("addWriteList address[%s,%d] format error", address, len(address)))
			return
		}
		//其实ctx可以直接存string，但dockergo的sdk存的byte
		ctx.PutState("identity", address, "1")
	}
	ctx.SuccessResult("add write list success")
	return
}

//go:wasmexport alterAdminAddress
func alterAdminAddress() {
	ctx := NewSimContext()
	paramAddress, _ := ctx.ArgString(paramAddress)
	var adminAddress []string
	if len(paramAddress) != 0 {
		adminAddress = strings.Split(paramAddress, ",")
	}
	if len(adminAddress) == 0 {
		ctx.ErrorResult("adminAddress of param should not be empty")
		return
	}

	if !senderIsAdmin() {
		ctx.ErrorResult("sender is not admin")
		return
	}
	adminAddressByte, _ := json.Marshal(adminAddress)
	err := ctx.PutStateByte("identity", keyAdminAddress, adminAddressByte)
	if err != SUCCESS {
		ctx.ErrorResult("alter admin address of identityInfo failed")
		return
	}
	ctx.EmitEvent("alterAdminAddress", "")
	ctx.SuccessResult("OK")
	return
}

func senderIsAdmin() bool {
	ctx := NewSimContext()
	sender, _ := ctx.GetSenderPk()
	adminAddressByte, err := ctx.GetStateByte("identity", keyAdminAddress)
	if err != SUCCESS || len(adminAddressByte) == 0 {
		ctx.Log(fmt.Sprintf("Get adminAddressList failed, err:%s", err))
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

//go:wasmexport removeWriteList
func removeWriteList() {
	ctx := NewSimContext()
	paramAddress, _ := ctx.ArgString(paramAddress)
	var addresses []string
	if len(paramAddress) != 0 {
		addresses = strings.Split(paramAddress, ",")
	}
	if len(addresses) == 0 {
		ctx.ErrorResult("address of param should not be empty")
		return
	}
	for _, address := range addresses {
		ctx.DeleteState("identity", address)
	}
	ctx.SuccessResult("remove write list success")
	return
}

//go:wasmexport isApprovedUser
func isApprovedUser() {
	ctx := NewSimContext()
	paramAddress, _ := ctx.ArgString(paramAddress)
	var addresses []string
	if len(paramAddress) != 0 {
		addresses = strings.Split(paramAddress, ",")
	}
	if len(addresses) == 0 {
		ctx.ErrorResult("address of param should not be empty")
		return
	}

	flag := true
	for _, addr := range addresses {
		val, err := ctx.GetState("identity", addr)
		if err != SUCCESS || len(val) == 0 {
			flag = false
		}
	}
	if flag {
		ctx.SuccessResult("true")
		return
	} else {
		ctx.SuccessResult("false")
		return
	}
}

//go:wasmexport address
func address() {
	ctx := NewSimContext()
	addr, err := ctx.GetSenderPk()
	if err != SUCCESS {
		ctx.ErrorResult(fmt.Sprintf("get address of sender failed, err:%s", err))
		return
	}
	if len(addr) == 0 {
		ctx.ErrorResult("addr is empty")
		return
	}
	//sdk.Instance.Infof("sender is %s, len is %d", addr, len(addr))
	//ctx.SuccessResult("address success")
	return

	ctx.SuccessResult(fmt.Sprintf("sender is %s, len is %d", addr, len(addr)))
	return
}

//go:wasmexport callerAddress
func callerAddress() {
	var param = make(map[string][]byte)
	ctx := NewSimContext()
	resp, resultCode := ctx.CallContract("identity", "address", param)
	if resultCode != SUCCESS {
		ctx.ErrorResult(fmt.Sprintf("call contract failed, err:%s", resp))
		return
	} else {
		ctx.SuccessResult(fmt.Sprintf("call contract success :%s", resp))
		return
	}
}

func main() {

}
