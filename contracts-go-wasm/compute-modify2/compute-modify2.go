/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"math/big"
	"unsafe"
)

// 安装合约时会执行此方法，必须
//
//go:wasmexport init_contract
func InitContract() {
	ctx := NewSimContext()
	ctx.SuccessResult("Init contract success")
	//fmt.Println("init contract test test")
}

// 升级合约时会执行此方法，必须
//
//go:wasmexport upgrade
func Upgrade() {
	ctx := NewSimContext()
	ctx.SuccessResult("Upgrade contract success")

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

	hashResult := make([]byte, 32)
	hashResultPtr := int32(uintptr(unsafe.Pointer(&hashResult[0])))

	for i := 0; i < 10000; i++ {
		nativeSha256(hashInput, hashResultPtr)
	}
	// 返回结果
	ctx.SuccessResult(fmt.Sprintf("success hashCal %x", hashResult))
}

//go:wasmimport env native_sha
func nativeSha256(hashInput string, hashResultPtr int32) int32

//go:wasmimport env native_BigExp
func nativeBigExp(num, exp, mod int64, resultPtr int32) int32

//go:wasmexport bigNumCal
func bigNumCal() {
	//startTime := time.Now().UnixNano()
	ctx := NewSimContext()
	a := int64(2)
	exp := int64(100000)
	mod := int64(1000000007)
	var resultLen int32
	var result *big.Int
	// 假设结果不会超过 256 字节（2048bit）
	resultBuf := make([]byte, 256)
	resultPtr := int32(uintptr(unsafe.Pointer(&resultBuf[0])))

	// 进行 10000 次大数运算
	for i := 0; i < 10000; i++ {

		resultLen = nativeBigExp(a, exp, mod, resultPtr)
		result = new(big.Int).SetBytes(resultBuf[:resultLen])
	}
	//endTime := time.Now().UnixNano()
	//ctx.SuccessResult(fmt.Sprintf("success bigNumCal: %s,startTime %d,endTime %d,executionTime %d", result.String(), startTime, endTime, endTime-startTime))
	ctx.SuccessResult(fmt.Sprintf("success bigNumCal: %s", result.String()))
}
func main() {

}
