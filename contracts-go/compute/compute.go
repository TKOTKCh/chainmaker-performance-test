package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

// ComputeContract for complexity testing
type ComputeContract struct{}

// InitContract installs the contract
func (c *ComputeContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

// UpgradeContract upgrades the contract
func (c *ComputeContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

// InvokeContract entry point
func (c *ComputeContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "normalCal":
		return c.normalCal()
	case "hashCal":
		return c.hashCal()
	case "bigNumCal":
		return c.bigNumCal()
	default:
		return sdk.Error("invalid method")
	}
}

// normalCal performs a simple arithmetic calculation
func (c *ComputeContract) normalCal() protogo.Response {
	result := 0
	for i := 0; i < 1000000; i++ {
		result += i
	}
	return sdk.Success([]byte(fmt.Sprintf("success normalCal: %d", result)))
}

// hashCal performs SHA-256 hashing multiple times
func (c *ComputeContract) hashCal() protogo.Response {
	hashInput := "ChainMaker Performance Test"
	var hashResult [32]byte
	for i := 0; i < 100000; i++ {
		hashResult = sha256.Sum256([]byte(hashInput))
	}
	return sdk.Success([]byte(fmt.Sprintf("success hashCal", hashResult)))
	//hashInput := "ChainMaker Performance Test"
	//for i := 0; i < 100000; i++ {
	//	sha256.Sum256([]byte(hashInput))
	//}
	//return sdk.Success([]byte(fmt.Sprintf("success hashCal")))

}

// bigNumCal performs large number exponentiation
func (c *ComputeContract) bigNumCal() protogo.Response {
	a := big.NewInt(2)
	exp := big.NewInt(100000)
	mod := big.NewInt(1000000007)
	var result *big.Int

	// 进行 10000 次大数运算
	for i := 0; i < 10000; i++ {
		result = new(big.Int).Exp(a, exp, mod)
	}

	// 返回最终一次计算的结果
	return sdk.Success([]byte(fmt.Sprintf("success bigNumCal: %s", result.String())))

}

func main() {
	err := sandbox.Start(new(ComputeContract))
	if err != nil {
		log.Fatal(err)
	}
}
