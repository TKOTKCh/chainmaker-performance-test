/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 安装合约时会执行此方法，必须
//
//go:wasmexport init_contract
func InitContract() {
	//ctx := NewSimContext()
	//ctx.SuccessResult("Init contract success")

}

// 升级合约时会执行此方法，必须
//
//go:wasmexport upgrade
func Upgrade() {
	//ctx := NewSimContext()
	//ctx.SuccessResult("Upgrade contract success")

}

//go:wasmexport manualInit
func manualInit() {
	ctx := NewSimContext()
	ctx.SuccessResult("Init contract success")

}

// tinygo的编译逻辑是我这里export了，他会继续找发现ctx.SuccessResult会调用syscall和logmessage
// 而这两个应该是由wasmer提供实现，所以函数体为空，这时编译出来的wasm文件会有import env syscall,如果函数体不为空就不会有
//
//go:wasmexport normalCal
func normalCal() {
	// 获取上下文
	ctx := NewSimContext()

	result := 0
	for i := 0; i < 1000000; i++ {
		result += i
	}

	// 返回结果
	ctx.SuccessResult(fmt.Sprintf("success normalCal: %d", result))
}

//go:wasmexport hashCal
func hashCal() {
	// 获取上下文
	ctx := NewSimContext()

	hashInput := "ChainMaker Performance Test"
	var hashResult [32]byte
	for i := 0; i < 100000; i++ {
		hashResult = sha256.Sum256([]byte(hashInput))
	}
	// 返回结果
	ctx.SuccessResult(fmt.Sprintf("success hashCal %x", hashResult))
}

//go:wasmexport bigNumCal
func bigNumCal() {
	ctx := NewSimContext()
	a := big.NewInt(2)
	exp := big.NewInt(100000)
	mod := big.NewInt(1000000007)
	var result *big.Int

	// 进行 10000 次大数运算
	for i := 0; i < 10000; i++ {
		result = new(big.Int).Exp(a, exp, mod)
	}
	ctx.SuccessResult(fmt.Sprintf("success bigNumCal: %s", result.String()))
}
func main() {

}
