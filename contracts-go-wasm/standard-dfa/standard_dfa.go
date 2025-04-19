package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	//"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	//"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	//"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"chainmaker.org/chainmaker/contract-utils/address"
	"chainmaker.org/chainmaker/contract-utils/safemath"
	"chainmaker.org/chainmaker/contract-utils/standard"
)

const (
	pName          = "name"
	pSymbol        = "symbol"
	pDecimals      = "decimals"
	pAccount       = "account"
	pAmount        = "amount"
	pTotalSupply   = "totalSupply"
	pFrom          = "from"
	pTo            = "to"
	pSpender       = "spender"
	pOwner         = "owner"
	balanceKey     = "b"
	allowanceKey   = "a"
	totalSupplyKey = "totalSupplyKey"
	nameKey        = "name"
	symbolKey      = "symbol"
	decimalKey     = "decimal"
	adminKey       = "admin"
)

var (
	defaultName        = "TestToken"
	defaultSymbol      = "TT"
	defaultDecimals    = 18
	defaultTotalSupply = safemath.SafeUintZero
)

// 安装合约时会执行此方法，必须
//
//go:wasmexport init_contract
func InitContract() {
	ctx := NewSimContext()
	err1 := updateErc20Info()
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	ctx.SuccessResult("Init contract success")

}

// 升级合约时会执行此方法，必须
//
//go:wasmexport upgrade
func Upgrade() {
	ctx := NewSimContext()
	ctx.SuccessResult("Upgrade contract success")

}

// UpgradeContract upgrade contract func
func updateErc20Info() error {
	ctx := NewSimContext()
	// name, symbol and decimal are optional
	name, _ := ctx.ArgString(pName)
	symbol, _ := ctx.ArgString(pSymbol)
	decimalsStr, _ := ctx.ArgString(pDecimals)
	totalSupply, _ := ctx.ArgString(pTotalSupply)

	if len(name) > 0 {
		defaultName = string(name)
	}
	err := ctx.PutState(nameKey, "", defaultName)
	if err != SUCCESS {
		return fmt.Errorf("set name of erc20Info failed")
	}
	if len(symbol) > 0 {
		defaultSymbol = string(symbol)
	}
	err = ctx.PutState(symbolKey, "", defaultSymbol)
	if err != SUCCESS {
		return fmt.Errorf("set symbol of erc20Info failed")
	}
	//decimals default to 18
	if len(decimalsStr) > 0 {
		num, err1 := strconv.Atoi(string(decimalsStr))
		if err1 != nil {
			return fmt.Errorf("param decimals err")
		}
		defaultDecimals = num
	}
	err = ctx.PutState(decimalKey, "", strconv.Itoa(defaultDecimals))
	if err != SUCCESS {
		return fmt.Errorf("set decimal of erc20Info failed")
	}
	//set admin
	admin, err := ctx.GetSenderPk()
	if len(admin) > 0 {
		admin = admin[:40]
	}
	if err != SUCCESS {
		return fmt.Errorf("get sender failed")
	}
	err = ctx.PutState(adminKey, "", admin)
	if err != SUCCESS {
		return fmt.Errorf("set admin of erc20Info failed")
	}
	//set totalSupply
	if len(totalSupply) > 0 {
		totalSupplyNum, ok := safemath.ParseSafeUint256(string(totalSupply))
		if !ok {
			return errors.New("invalid totalSupply number")
		}
		defaultTotalSupply = totalSupplyNum
	}
	//Mint
	if defaultTotalSupply.ToString() != "0" {
		return baseMint(admin, defaultTotalSupply)
	}
	return nil
}

////////////////////////////////Helper//////////////////////////////////

// ReturnUint256 封装返回SafeUint256类型为Response，如果有error则忽略num，封装error
// @param num
// @param err
// @return Response
func ReturnUint256(num *safemath.SafeUint256, err error) {
	ctx := NewSimContext()
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult((num.ToString()))
}

// ReturnString 封装返回string类型为Response，如果有error则忽略str，封装error
// @param str
// @param err
// @return Response
func ReturnString(str string, err error) {
	ctx := NewSimContext()
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult((str))
}

// ReturnJson 封装返回interface类型为json string Response
// @param data
// @return Response
func ReturnJson(data interface{}) {
	ctx := NewSimContext()
	standardsBytes, err := json.Marshal(data)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult(string(standardsBytes))
}

// Return 封装返回Bool类型为Response，如果有error则忽略bool，封装error
// @param err
// @return Response
func Return(err error) {
	ctx := NewSimContext()
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult("OK")
}

// ReturnUint8 封装返回uint8类型为Response，如果有error则忽略num，封装error
// @param num
// @param err
// @return Response
func ReturnUint8(num uint8, err error) {
	ctx := NewSimContext()
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult(strconv.Itoa(int(num)))
}

// RequireString 必须要有参数 string类型
// @param key
// @return string
// @return error
func RequireString(key string) (string, error) {
	ctx := NewSimContext()
	b, ok := ctx.Arg(key)
	if ok != SUCCESS || len(b) == 0 {
		return "", fmt.Errorf("CMDFA: require parameter:'%s'", key)
	}
	return string(b), nil
}

// RequireUint256 必须要有参数 Uint256类型
// @param key
// @return *safemath.SafeUint256
// @return error
func RequireUint256(key string) (*safemath.SafeUint256, error) {
	ctx := NewSimContext()
	b, ok := ctx.Arg(key)
	if ok != SUCCESS || len(b) == 0 {
		return nil, fmt.Errorf("CMDFA: require parameter:'%s'", key)
	}
	num, err1 := safemath.ParseSafeUint256(string(b))
	if !err1 {
		return nil, fmt.Errorf("CMDFA: parameter:'%s' not a valid uint256", key)
	}
	return num, nil
}

// Standards 合约支持的标准
// @return []string
//
//go:wasmexport Standards
func Standards() {
	ReturnJson([]string{standard.ContractStandardNameCMBC, standard.ContractStandardNameCMDFA})
}

// SupportStandard 查询本合约是否支持某合约标准
// @param standardName
// @return bool
func SupportStandard(standardName string) bool {
	return standardName == standard.ContractStandardNameCMDFA || standardName == standard.ContractStandardNameCMBC
}

// TotalSupply 发行总量
// @return *safemath.SafeUint256
// @return error
//
//go:wasmexport TotalSupply
func TotalSupply() {
	ReturnUint256(GetTotalSupply())
}

// BalanceOf 账户余额查询
// @param account
// @return *safemath.SafeUint256
// @return error
//
//go:wasmexport BalanceOf
func BalanceOf() {
	ctx := NewSimContext()
	account, err1 := RequireString(pAccount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	ReturnUint256(GetBalance(account))
}

// Transfer 转账操作
// @param to
// @param amount
// @return error
//
//go:wasmexport Transfer
func Transfer() {
	ctx := NewSimContext()
	to, err1 := RequireString(pTo)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		Return(err1)
		return
	}
	amount, err1 := RequireUint256(pAmount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		Return(err1)
		return
	}

	from, err := ctx.GetSenderPk()
	if len(from) > 40 {
		from = from[:40]
	}
	if err != SUCCESS {
		Return(fmt.Errorf("Get sender address failed"))
		return
	}
	err1 = baseTransfer(from, to, amount)
	if err1 != nil {
		Return(err1)
	}
	Return(nil)
}

// TransferFrom 代为转账操作
// @param from
// @param to
// @param amount
// @return error
//
//go:wasmexport TransferFrom
func TransferFrom() {
	ctx := NewSimContext()
	from, err1 := RequireString(pFrom)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	to, err1 := RequireString(pTo)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	amount, err1 := RequireUint256(pAmount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	sender, err := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	if err != SUCCESS {
		Return(fmt.Errorf("Get sender address failed"))
	}
	err1 = baseSpendAllowance(from, sender, amount)
	if err1 != nil {
		Return(fmt.Errorf("spend allowance failed, err:%s", err1))
	}
	err1 = baseTransfer(from, to, amount)
	if err1 != nil {
		Return(err1)
	}
	Return(nil)
}

// Approve 授权额度
// @param spender
// @param amount
// @return error
//
//go:wasmexport Approve
func Approve() {
	ctx := NewSimContext()
	spender, err1 := RequireString(pSpender)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	amount, err1 := RequireUint256(pAmount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	sender, err := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	if err != SUCCESS {
		Return(fmt.Errorf("Get sender address failed"))
	}
	err1 = baseApprove(sender, spender, amount)
	if err1 != nil {
		Return(err1)
	}
	Return(nil)
}

// Allowance 查询授权额度
// @param owner
// @param spender
// @return *safemath.SafeUint256
// @return error
//
//go:wasmexport Allowance
func Allowance() {
	ctx := NewSimContext()
	spender, err1 := RequireString(pSpender)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	owner, err1 := RequireString(pOwner)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	ReturnUint256(GetAllowance(owner, spender))
}

// Name Token的名字
// @return string
// @return error
//
//go:wasmexport Name
func Name() {
	ReturnString(GetName())
}

// Symbol Token的符号
// @return string
// @return error
//
//go:wasmexport Symbol
func Symbol() {
	ReturnString(GetSymbol())
}

// Decimals Token的小数位数
// @return uint8
// @return error
//
//go:wasmexport Decimals
func Decimals() {
	ReturnUint8(GetDecimals())
}

// Mint 铸造Token
// @param account
// @param amount
// @return error
//
//go:wasmexport Mint
func Mint() {
	//check is admin
	ctx := NewSimContext()
	account, err1 := RequireString(pAccount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	amount, err1 := RequireUint256(pAmount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	sender, err := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	if err != SUCCESS {
		Return(fmt.Errorf("Get sender address failed"))
	}
	admin, err1 := GetAdmin()
	if err1 != nil {
		Return(err1)
	}
	if sender != admin {
		Return(errors.New("only admin can mint tokens"))
	}
	//call base mint
	err1 = baseMint(account, amount)
	if err1 != nil {
		Return(err1)
	}
	Return(nil)
}

// Burn 销毁Token
// @param amount
// @return error
//
//go:wasmexport Burn
func Burn() {
	ctx := NewSimContext()
	amount, err1 := RequireUint256(pAmount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		Return(err1)
	}
	spender, err := ctx.GetSenderPk()
	if len(spender) > 40 {
		spender = spender[:40]
	}
	if err != SUCCESS {
		Return(fmt.Errorf("Get sender address failed"))
	}
	//call base burn
	err1 = baseBurn(spender, amount)
	if err1 != nil {
		Return(err1)
	}
	Return(nil)
}

// BurnFrom 代为销毁Token
// @param account
// @param amount
// @return error
//
//go:wasmexport BurnFrom
func BurnFrom() {
	ctx := NewSimContext()
	account, err1 := RequireString(pAccount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	amount, err1 := RequireUint256(pAmount)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	spender, err := ctx.GetSenderPk()
	if len(spender) > 40 {
		spender = spender[:40]
	}
	if err != SUCCESS {
		Return(fmt.Errorf("Get sender address failed"))
	}
	err1 = baseSpendAllowance(account, spender, amount)
	if err1 != nil {
		Return(err1)
	}
	//call base burn
	err1 = baseBurn(account, amount)
	if err1 != nil {
		Return(err1)
	}
	Return(nil)
}

// EmitTransferEvent 触发转账事件
// @param spender
// @param to
// @param amount
func EmitTransferEvent(spender, to string, amount *safemath.SafeUint256) {
	ctx := NewSimContext()
	ctx.EmitEvent(standard.TopicTransfer, spender, to, amount.ToString())
}

// EmitApproveEvent 触发授权事件
// @param owner
// @param spender
// @param amount
func EmitApproveEvent(owner, spender string, amount *safemath.SafeUint256) {
	ctx := NewSimContext()
	ctx.EmitEvent(standard.TopicApprove, owner, spender, amount.ToString())
}

// EmitMintEvent 触发铸造事件
// @param account
// @param amount
func EmitMintEvent(account string, amount *safemath.SafeUint256) {
	ctx := NewSimContext()
	ctx.EmitEvent(standard.TopicMint, account, amount.ToString())
}

// EmitBurnEvent 触发销毁事件
// @param spender
// @param amount
func EmitBurnEvent(spender string, amount *safemath.SafeUint256) {
	ctx := NewSimContext()
	ctx.EmitEvent(standard.TopicBurn, spender, amount.ToString())
}

/////////////////////////Data Access Layer/////////////////////////////////

func createCompositeKey(prefix string, data ...string) string {
	return prefix + "_" + strings.Join(data, "_")
}

// GetUint256 获得DB中的SafeUint256
// @param key
// @param field
// @return *safemath.SafeUint256
// @return error
func GetUint256(key, field string) (*safemath.SafeUint256, error) {
	ctx := NewSimContext()
	fromBalStr, err := ctx.GetState(key, field)
	if err != SUCCESS {
		return nil, errors.New("func GetUint256 GetState failed")
	}

	fromBalance, pass := safemath.ParseSafeUint256(string(fromBalStr))
	if !pass {
		return nil, errors.New("invalid uint256 data")
	}
	//ctx.EmitEvent("GetUint256", key, field, fromBalance.ToString())
	return fromBalance, nil
}

// GetBalance 获得DB中的账户余额
// @param account
// @return *safemath.SafeUint256
// @return error
func GetBalance(account string) (*safemath.SafeUint256, error) {
	return GetUint256(balanceKey, account)
}

// SetBalance 设置DB中的账户余额
// @param account
// @param amount
// @return error
func SetBalance(account string, amount *safemath.SafeUint256) error {
	ctx := NewSimContext()
	err := ctx.PutState(balanceKey, account, amount.ToString())
	if err != SUCCESS {
		return errors.New("SetBalance PutState failed")
	}
	//ctx.EmitEvent("SetBalance", account, amount.ToString())
	return nil
}

// SetAllowance 设置DB中的授权额度
// @param owner
// @param spender
// @param amount
// @return error
func SetAllowance(owner string, spender string, amount *safemath.SafeUint256) error {
	key := createCompositeKey(allowanceKey, owner, spender)
	ctx := NewSimContext()
	err := ctx.PutState(key, "", amount.ToString())
	if err != SUCCESS {
		return errors.New("SetAllowance PutState failed")
	}
	//ctx.EmitEvent("SetAllowance", owner, spender, amount.ToString())
	return nil
}

// GetAllowance 获得DB中的授权额度
// @param owner
// @param spender
// @return *safemath.SafeUint256
// @return error
func GetAllowance(owner string, spender string) (*safemath.SafeUint256, error) {
	key := createCompositeKey(allowanceKey, owner, spender)
	return GetUint256(key, "")
}

// GetTotalSupply 获得发行总量
// @return *safemath.SafeUint256
// @return error
func GetTotalSupply() (*safemath.SafeUint256, error) {
	return GetUint256(totalSupplyKey, "")
}

// SetTotalSupply 设置发行总量
// @param amount
// @return error
func SetTotalSupply(amount *safemath.SafeUint256) error {
	ctx := NewSimContext()
	err := ctx.PutState(totalSupplyKey, "", amount.ToString())

	if err != SUCCESS {
		return errors.New("GetTotalSupply PutState failed")
	}
	//ctx.EmitEvent("SetTotalSupply", amount.ToString())
	return nil
}

// GetName 获得DB中的Name
// @return string
// @return error
func GetName() (string, error) {
	ctx := NewSimContext()
	name, err := ctx.GetState(nameKey, "")
	if err != SUCCESS {
		return "", errors.New("GetName GetState failed")
	}
	return name, nil
}

// SetName 设置DB中的Name
// @param name
// @return error
func SetName(name string) error {
	ctx := NewSimContext()
	err := ctx.PutState(nameKey, "", name)
	if err != SUCCESS {
		return errors.New("SetName PutState failed")
	}
	return nil
}

// GetSymbol 获得DB中的符号
// @return string
// @return error
func GetSymbol() (string, error) {
	ctx := NewSimContext()
	symbol, err := ctx.GetState(symbolKey, "")
	if err != SUCCESS {
		return "", errors.New("GetSymbol GetState failed")
	}
	return symbol, nil
}

// SetSymbol 设置DB中的符号
// @param symbol
// @return error
func SetSymbol(symbol string) error {
	ctx := NewSimContext()
	err := ctx.PutState(symbolKey, "", symbol)
	if err != SUCCESS {
		return errors.New("SetSymbol PutState failed")
	}
	return nil
}

// GetDecimals 获得DB中的小数位数
// @return uint8
// @return error
func GetDecimals() (uint8, error) {
	ctx := NewSimContext()
	d, err1 := ctx.GetState(decimalKey, "")
	if err1 != SUCCESS {
		return 0, errors.New("GetDecimals GetState failed")
	}
	decimal, err := strconv.ParseUint(d, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(decimal), nil
}

// SetDecimals 设置DB中的小数位数
// @param decimal
// @return error
func SetDecimals(decimal uint8) error {
	ctx := NewSimContext()
	err := ctx.PutState(decimalKey, "", strconv.Itoa(int(decimal)))
	if err != SUCCESS {
		return errors.New("SetDecimals PutState failed")
	}
	return nil
}

// GetAdmin 获得DB中的Admin
// @return string
// @return error
func GetAdmin() (string, error) {
	ctx := NewSimContext()
	admin, err := ctx.GetState(adminKey, "")
	if err != SUCCESS {
		return "", errors.New("GetAdmin GetState failed")
	}
	return admin, nil
}

// SetAdmin 设置DB中的Admin
// @param admin
// @return error
func SetAdmin(admin string) error {
	ctx := NewSimContext()
	err := ctx.PutState(adminKey, "", admin)
	if err != SUCCESS {
		return errors.New("SetAdmin PutState failed")
	}
	return nil
}

////////////////////////////CMDFA Core/////////////////////////////

// baseTransfer
// @param from
// @param to
// @param amount
// @return error
func baseTransfer(from string, to string, amount *safemath.SafeUint256) error {
	//检查from和to的合法性
	if !address.IsValidAddress(from) {
		return errors.New("CMDFA: transfer from the invalid address")
	}
	if address.IsZeroAddress(from) {
		return errors.New("CMDFA: transfer from the zero address")
	}
	if !address.IsValidAddress(to) {
		return errors.New("CMDFA: transfer to the invalid address")
	}
	if address.IsZeroAddress(to) {
		return errors.New("CMDFA: transfer to the zero address")
	}

	//检查from余额充足
	fromBalance, err := GetBalance(from)
	if err != nil {
		return err
	}
	if !fromBalance.GTE(amount) {
		return errors.New("CMDFA: transfer amount exceeds balance")
	}
	//更新from和to的余额
	fromNewBalance, _ := safemath.SafeSub(fromBalance, amount)
	err = SetBalance(from, fromNewBalance)
	if err != nil {
		return err
	}
	toBalance, err := GetBalance(to)
	if err != nil {
		return err
	}
	toNewBalance, ok := safemath.SafeAdd(toBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = SetBalance(to, toNewBalance)
	if err != nil {
		return err
	}
	//触发事件
	EmitTransferEvent(from, to, amount)

	return nil
}

func baseApprove(owner string, spender string, amount *safemath.SafeUint256) error {
	//检查from和to的合法性
	if !address.IsValidAddress(owner) {
		return errors.New("CMDFA: approve from the invalid address")
	}
	if address.IsZeroAddress(owner) {
		return errors.New("CMDFA: approve from the zero address")
	}
	if !address.IsValidAddress(spender) {
		return errors.New("CMDFA: approve to the invalid address")
	}
	if address.IsZeroAddress(spender) {
		return errors.New("CMDFA: approve to the zero address")
	}
	//设置Allowance
	err := SetAllowance(owner, spender, amount)
	if err != nil {
		return err
	}
	//触发事件Approval
	EmitApproveEvent(owner, spender, amount)
	return nil
}

func baseSpendAllowance(owner string, spender string, amount *safemath.SafeUint256) error {
	//获得授权的额度
	currentAllowance, err := GetAllowance(owner, spender)
	if err != nil {
		return err
	}
	// Does not update the allowance amount in case of infinite allowance.
	if currentAllowance.IsMaxSafeUint256() {
		return nil
	}
	//计算额度是否够用
	if !currentAllowance.GTE(amount) {
		return errors.New("CMDFA: insufficient allowance")
	}
	//扣减授权额度
	newCurrentAllowance, ok := safemath.SafeSub(currentAllowance, amount)
	if !ok {
		return errors.New("spend allowance error")
	}
	return baseApprove(owner, spender, newCurrentAllowance)
}
func baseMint(account string, amount *safemath.SafeUint256) error {
	//检查account的合法性
	if !address.IsValidAddress(account) {
		return errors.New("CMDFA: mint to the invalid address")
	}
	if address.IsZeroAddress(account) {
		return errors.New("CMDFA: mint to the zero address")
	}
	//ctx := NewSimContext()
	//更新TotalSupply
	totalSupply, err := GetTotalSupply()
	if err != nil {
		return err
	}
	//ctx.EmitEvent("GetTotalSupply Success", totalSupply.ToString())

	newTotal, ok := safemath.SafeAdd(totalSupply, amount)
	if !ok {
		return errors.New("calculate totalSupply failed")
	}
	err = SetTotalSupply(newTotal)
	if err != nil {
		return err
	}
	//ctx.EmitEvent("SetTotalSupply Success", newTotal.ToString())

	//更新余额
	toBalance, err := GetBalance(account)
	if err != nil {
		return err
	}
	//ctx.EmitEvent("GetBalance Success", toBalance.ToString())
	toNewBalance, ok := safemath.SafeAdd(toBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = SetBalance(account, toNewBalance)
	if err != nil {
		return err
	}
	//ctx.EmitEvent("SetBalance Success", toNewBalance.ToString())

	//触发事件
	EmitMintEvent(account, amount)
	return nil
}

func baseBurn(account string, amount *safemath.SafeUint256) error {
	//检查account的合法性
	if !address.IsValidAddress(account) {
		return errors.New("CMDFA: burn from the invalid address")
	}
	if address.IsZeroAddress(account) {
		return errors.New("CMDFA: burn from the zero address")
	}
	//检查用户余额充足
	fromBalance, err := GetBalance(account)
	if err != nil {
		return err
	}
	if !fromBalance.GTE(amount) {
		return errors.New("CMDFA: burn amount exceeds balance")
	}
	//更新TotalSupply
	totalSupply, err := GetTotalSupply()
	if err != nil {
		return err
	}
	newTotal, ok := safemath.SafeSub(totalSupply, amount)
	if !ok {
		return errors.New("calculate totalSupply failed")
	}
	err = SetTotalSupply(newTotal)
	if err != nil {
		return err
	}
	//更新余额
	fromNewBalance, ok := safemath.SafeSub(fromBalance, amount)
	if !ok {
		return errors.New("calculate new to balance error")
	}
	err = SetBalance(account, fromNewBalance)
	if err != nil {
		return err
	}
	//触发事件
	EmitBurnEvent(account, amount)

	return nil
}

func main() {

}
