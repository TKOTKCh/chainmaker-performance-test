package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/contract-utils/address"
	"chainmaker.org/chainmaker/contract-utils/safemath"
)

// Address 地址格式字符串
type Address string

const (
	// use for approve true
	trueString = "1"
	// use for approve false
	falseString = "0"
)
const (
	keyAdminAddress   = "adminAddress"
	keyUri            = "uri"
	paramAdminAddress = "adminAddress"
	paramAmount       = "amount"
	paramAmounts      = "amounts"
	paramApproved     = "approved"
	paramData         = "data"
	paramFrom         = "from"
	paramTo           = "to"
	paramId           = "id"
	paramIds          = "ids"
	paramIdStart      = "idStart"
	paramOperator     = "operator"
	paramOwner        = "owner"
	paramUri          = "uri"

	Erc1155Map = "fxtoonErc1155"
	ApproveMap = "approve"
	BalanceMap = "balance"
	TokenMap   = "token"
)

//go:wasmexport init_contract
func InitContract() {
	ctx := NewSimContext()

	// 处理管理员地址
	adminAddress, _ := ctx.Arg(paramAdminAddress)
	uri, _ := ctx.ArgString(paramUri)

	var adminAddressStr string
	if len(adminAddress) == 0 {
		adminAddressStr, _ = ctx.GetSenderPk()
		if len(adminAddressStr) > 40 {
			adminAddressStr = adminAddressStr[:40]
		}
	} else {
		adminAddressStr = string(adminAddress)

	}
	adminAddresses := strings.Split(adminAddressStr, ",")
	resp := initContract(adminAddresses, uri)
	if !resp {
		ctx.ErrorResult("Init contract failed")
		return
	} else {
		ctx.SuccessResult("Init contract success")
		return
	}
}

//go:wasmexport upgrade
func Upgrade() {
	ctx := NewSimContext()
	ctx.SuccessResult("Upgrade contract success")
}

//go:wasmexport Uri
func GetUri() {
	ctx := NewSimContext()
	idStr, _ := ctx.ArgString(paramId)
	tokenId, _ := safemath.ParseSafeUint256(idStr)
	val, _ := ctx.GetStateByte(Erc1155Map, keyUri)
	url := strings.ReplaceAll(string(val), "{tokenId}", tokenId.ToString())
	ctx.SuccessResult(url)
}

// SetUri 设置uri
//
//go:wasmexport SetUri
func SetUri() {
	ctx := NewSimContext()
	data, _ := ctx.ArgString(paramUri)
	if !senderIsAdmin() {
		ctx.ErrorResult("No permission")
		return
	}
	err := ctx.PutStateByte(Erc1155Map, keyUri, []byte(data))
	if err != SUCCESS {
		//ctx.Warnf("Set uri failed, err:%s", err)
		ctx.ErrorResult("Set uri failed")
		return
	}
	successData(data)
}

// AlterAdminAddress 修改管理员
//
//go:wasmexport AlterAdminAddress
func AlterAdminAddress() {
	ctx := NewSimContext()
	adminAddrBytes, _ := ctx.Arg(paramAdminAddress)
	var adminAddress []string
	if len(adminAddrBytes) != 0 {
		adminAddress = strings.Split(string(adminAddrBytes), ",")
	}
	if len(adminAddress) == 0 {
		ctx.ErrorResult("adminAddress of param should not be empty")
		return
	}
	if !senderIsAdmin() {
		ctx.ErrorResult("sender is not admin")
		return
	}
	if !address.IsValidAddress(adminAddress...) {
		ctx.ErrorResult("address format error")
		return
	}
	adminAddressByte, _ := json.Marshal(adminAddress)

	err := ctx.PutStateByte(Erc1155Map, keyAdminAddress, adminAddressByte)
	if err != SUCCESS {
		ctx.ErrorResult("alter admin address of identityInfo failed.")
		return
	}
	ctx.EmitEvent("AlterAdminAddress", adminAddress...)
	successNormal()
}

// SafeTransferFrom 交易
//
//go:wasmexport SafeTransferFrom
func SafeTransferFrom() { // nolint
	ctx := NewSimContext()
	fromStr, _ := ctx.ArgString(paramFrom)
	toStr, _ := ctx.ArgString(paramTo)
	idStr, _ := ctx.ArgString(paramId)
	amountStr, _ := ctx.ArgString(paramAmount)

	fromStr = address.GetCleanAddr(fromStr)
	toStr = address.GetCleanAddr(toStr)
	id, ok := safemath.ParseSafeUint256(idStr)
	if !ok {
		ctx.ErrorResult("parse id failed")
		return
	}
	value, ok := safemath.ParseSafeUint256(amountStr)
	if !ok {
		ctx.ErrorResult("parse amount failed")
		return
	}
	from := Address(fromStr)
	to := Address(toStr)

	if !address.IsValidAddress(string(from), string(to)) {
		ctx.ErrorResult("address format error")
		return
	}
	err := transferCore(from, to, id, value)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}

	// emit event
	ctx.EmitEvent("SafeTransferFrom", string(from), string(to), id.ToString(), value.ToString())
	successNormal()
}

// SafeBatchTransferFrom 批量交易
//
//go:wasmexport SafeBatchTransferFrom
func SafeBatchTransferFrom() { // nolint
	ctx := NewSimContext()
	fromStr, _ := ctx.ArgString(paramFrom)
	toStr, _ := ctx.ArgString(paramTo)
	idsStr, _ := ctx.ArgString(paramIds)
	amountsStr, _ := ctx.ArgString(paramAmounts)
	//data, _ := ctx.ArgString(paramData)

	fromStr = address.GetCleanAddr(fromStr)
	toStr = address.GetCleanAddr(toStr)
	from := Address(fromStr)
	to := Address(toStr)
	idArr := strings.Split(idsStr, ",")
	amountArr := strings.Split(amountsStr, ",")
	if len(amountArr) != len(idArr) || len(amountArr) == 0 {
		ctx.ErrorResult("param error")
		return
	}

	var ids = make([]*safemath.SafeUint256, 0)
	var values = make([]*safemath.SafeUint256, 0)
	for i, id := range idArr {
		tokenId, err := safemath.ParseSafeUint256(id)
		if !err {
			ctx.ErrorResult("parse id failed")
			return
		}
		amount, err := safemath.ParseSafeUint256(amountArr[i])
		if !err {
			ctx.ErrorResult("parse amount failed")
			return
		}
		ids = append(ids, tokenId)
		values = append(values, amount)
	}
	if len(ids) != len(values) {
		ctx.ErrorResult("param len error")
		return
	}
	if !address.IsValidAddress(string(from), string(to)) {
		ctx.ErrorResult("address format error")
		return
	}
	for i, id := range ids {
		value := values[i]
		err := transferCore(from, to, id, value)
		if err != nil {
			ctx.ErrorResult(err.Error())
			return
		}
	}
	// emit event
	idsByte, _ := json.Marshal(ids)
	valuesStr, _ := json.Marshal(values)
	ctx.EmitEvent("SafeTransferFrom", string(from), string(to), string(idsByte), string(valuesStr))
	successNormal()
}

// BalanceOf 余额
//
//go:wasmexport BalanceOf
func BalanceOf() {
	ctx := NewSimContext()
	owner, _ := ctx.ArgString(paramOwner)
	idStr, _ := ctx.ArgString(paramId)
	owner = address.GetCleanAddr(owner)

	id, err1 := safemath.ParseSafeUint256(idStr)
	if !err1 {
		ctx.ErrorResult("parse id failed")
		return
	}
	balanceByte, err := ctx.GetStateByte(BalanceMap, string(owner)+"_"+id.ToString())
	if err != SUCCESS {
		ctx.ErrorResult("get balance error.")
		return
	}
	balance, _ := safemath.ParseSafeUint256(string(balanceByte))
	ctx.SuccessResult(balance.ToString())
}

// BalanceOfBatch 批量查询余额
//
//go:wasmexport BalanceOfBatch
func BalanceOfBatch() {
	ctx := NewSimContext()
	ownersStr, _ := ctx.ArgString(paramOwner)
	idStrs, _ := ctx.ArgString(paramId)

	o := strings.Split(ownersStr, ",")
	tempids := strings.Split(idStrs, ",")

	if len(tempids) != len(o) {
		ctx.ErrorResult("param error")
		return
	}
	var ids = make([]*safemath.SafeUint256, 0)
	var owners = make([]Address, 0)
	for i, id := range tempids {
		tokenId, err := safemath.ParseSafeUint256(id)
		if !err {
			ctx.ErrorResult("parse id failed")
			return
		}
		ids = append(ids, tokenId)
		owners = append(owners, Address(o[i]))
	}
	if len(owners) != len(ids) {
		ctx.ErrorResult("param len error")
		return
	}
	var balances []string
	for i, owner := range owners {
		id := ids[i]
		balanceByte, err := ctx.GetStateByte(BalanceMap, string(owner)+"_"+id.ToString())
		if err != SUCCESS {
			ctx.ErrorResult("get balance error.")
		}
		balance, _ := safemath.ParseSafeUint256(string(balanceByte))
		balances = append(balances, balance.ToString())
	}
	// 1,2,0,3,3
	balancesJson, _ := json.Marshal(balances)
	ctx.SuccessResultByte(balancesJson)
}

// SetApprovalForAll 授权
//
//go:wasmexport SetApprovalForAll
func SetApprovalForAll() {
	ctx := NewSimContext()
	papprovedStr, _ := ctx.ArgString(paramApproved)
	operatorStr, _ := ctx.ArgString(paramOperator)
	senderStr, _ := ctx.ArgString("sender")
	operatorStr = address.GetCleanAddr(operatorStr)

	approved, err1 := strconv.ParseBool(papprovedStr)
	if err1 != nil {
		ctx.ErrorResult("parse approved failed" + err1.Error())
		return
	}
	var sender Address
	if len(senderStr) == 0 {
		senderStr, _ = ctx.GetSenderPk()
		if len(senderStr) > 40 {
			senderStr = senderStr[:40]
		}
	}
	sender = Address(senderStr)
	operator := Address(operatorStr)

	// check param
	if !address.IsValidAddress(string(operator)) {
		ctx.ErrorResult("operator address invalid")
		return
	}
	if sender == operator {
		ctx.ErrorResult("approve to caller")
		return
	}

	var approvedStr string
	if approved {
		approvedStr = trueString
	} else {
		approvedStr = falseString
	}
	// save approval
	err := ctx.PutStateByte(ApproveMap, string(sender)+"_"+string(operator), []byte(approvedStr))
	if err != SUCCESS {
		ctx.ErrorResult("set ApprovalForAll error, ")
		return
	}
	ctx.EmitEvent("ApprovalForAll", string(sender), string(operator), approvedStr)
	successNormal()
}

// IsApprovedForAll 查询权限
//
//go:wasmexport IsApprovedForAll
func IsApprovedForAll() bool {
	ctx := NewSimContext()
	owner, _ := ctx.ArgString(paramOwner)
	operator, _ := ctx.ArgString(paramOperator)

	owner = address.GetCleanAddr(owner)
	operator = address.GetCleanAddr(operator)
	val, err := ctx.GetState(ApproveMap, string(owner)+"_"+string(operator))
	if err != SUCCESS {
		ctx.ErrorResult(fmt.Sprintf("get approved val from approve info failed, err:%s", err))
		return false
	}
	if string(val) == trueString {
		ctx.SuccessResult(trueString)
		return true
	}
	ctx.SuccessResult(falseString)
	return false
}
func isApprovedForAll(owner, operator Address) bool {
	ctx := NewSimContext()
	val, err := ctx.GetState(ApproveMap, string(owner)+"_"+string(operator))
	if err != SUCCESS {
		ctx.ErrorResult(fmt.Sprintf("get approved val from approve info failed, err:%s", err))
		return false
	}
	ctx.EmitEvent("isApprovedForAll", string(owner), string(operator), val)
	if string(val) == trueString {
		ctx.SuccessResult(trueString)
		return true
	}
	ctx.SuccessResult(falseString)
	return false
}

// MintBatchNormal 批量发行
//
//go:wasmexport MintBatchNormal
func MintBatchNormal() {
	ctx := NewSimContext()
	toStr, _ := ctx.ArgString(paramTo)
	ids, _ := ctx.ArgString(paramIds)
	amounts, _ := ctx.ArgString(paramAmounts)
	data, _ := ctx.ArgString(paramData)
	toStr = address.GetCleanAddr(toStr)
	to := Address(toStr)
	if len(ids) == 0 || len(amounts) == 0 {
		ctx.ErrorResult("param empty")
		return
	}
	idStrArr := strings.Split(ids, ",")
	amountStrArr := strings.Split(amounts, ",")
	if len(idStrArr) != len(amountStrArr) {
		ctx.ErrorResult("param err")
		return
	}
	var idsArr []*safemath.SafeUint256
	var amountArr []*safemath.SafeUint256
	for i, idStr := range idStrArr {
		amountStr := amountStrArr[i]
		id, flag := safemath.ParseSafeUint256(idStr)
		if !flag {
			ctx.ErrorResult("param id error")
			return
		}
		amount, flag := safemath.ParseSafeUint256(amountStr)
		if !flag {
			ctx.ErrorResult("param amount error")
			return
		}
		idsArr = append(idsArr, id)
		amountArr = append(amountArr, amount)
	}
	MintBatch(to, idsArr, amountArr, data)
}

// MintBatchNft 批量发行nft
//
//go:wasmexport MintBatchNft
func MintBatchNft() { // nolint
	ctx := NewSimContext()
	toStr, _ := ctx.ArgString(paramTo)
	idStartStr, _ := ctx.ArgString(paramIdStart)
	data, _ := ctx.ArgString(paramData)
	amountStr, _ := ctx.ArgString(paramAmount)
	toStr = address.GetCleanAddr(toStr)
	to := Address(toStr)
	idStart, flag1 := safemath.ParseSafeUint256(idStartStr)
	if !flag1 {
		ctx.ErrorResult("param idStart error")
		return
	}
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		ctx.ErrorResult("param amount error")
		return
	}
	// check param
	if amount <= 0 {
		ctx.ErrorResult("nft amount mast large than 0")
		return
	}
	if !senderIsAdmin() {
		ctx.ErrorResult("No permission")
		return
	}
	if !address.IsValidAddress(string(to)) {
		ctx.ErrorResult("param error")
		return
	}
	if !idStart.GTE(safemath.SafeUintOne) {
		ctx.ErrorResult("param idStart error")
		return
	}
	if len(ownerOf(idStart)) > 0 {
		ctx.ErrorResult("exists id " + idStart.ToString())
		return
	}

	// make param
	var tokenIdArr []*safemath.SafeUint256
	var amountArr []*safemath.SafeUint256
	var tokenIdStart = idStart

	for i := 0; i < amount; i++ {
		tokenIdArr = append(tokenIdArr, tokenIdStart)
		amountArr = append(amountArr, safemath.SafeUintOne)
		tokenIdStart, _ = safemath.SafeAdd(tokenIdStart, safemath.SafeUintOne)
	}

	MintBatch(to, tokenIdArr, amountArr, data)
}

// Mint 发行
//
//go:wasmexport Mint
func Mint() {
	ctx := NewSimContext()
	to, _ := ctx.ArgString(paramTo)
	idStr, _ := ctx.ArgString(paramId)
	amountStr, _ := ctx.ArgString(paramAmount)
	to = address.GetCleanAddr(to)

	id, flag1 := safemath.ParseSafeUint256(idStr)
	if !flag1 {
		ctx.ErrorResult("param error")
		return
	}
	amount, flag2 := safemath.ParseSafeUint256(amountStr)
	if !flag2 {
		ctx.ErrorResult("param error")
		return
	}
	if !address.IsValidAddress(string(to)) {
		ctx.ErrorResult("param to error")
		return
	}
	if !senderIsAdmin() {
		ctx.ErrorResult("No permission")
		return
	}
	if !id.GTE(safemath.SafeUintZero) {
		ctx.ErrorResult("param id error")
		return
	}
	//if !address.IsValidAddress(string(to)) {
	//	ctx.ErrorResult("param error")
	//	return
	//}
	balanceByte, err := ctx.GetStateByte(BalanceMap, string(to)+"_"+id.ToString())
	if err != SUCCESS {
		ctx.ErrorResult("Mint get balance error")
		return
	}
	balance, _ := safemath.ParseSafeUint256(string(balanceByte))
	balance, _ = safemath.SafeAdd(balance, amount)

	// set balance
	_ = ctx.PutStateByte(BalanceMap, string(to)+"_"+id.ToString(), []byte(balance.ToString()))
	// set owner only nft valid
	_ = ctx.PutState("owner", id.ToString(), string(to))
	// emit event
	ctx.EmitEvent("Mint", string(to), id.ToString(), amount.ToString())
	// set the number of token transactions
	_ = ctx.PutStateByte(TokenMap, id.ToString(), []byte("0"))

	successNormal()
}

// MintBatch 批量发行实现
func MintBatch(to Address, ids, amounts []*safemath.SafeUint256, data string) {
	ctx := NewSimContext()
	// check param
	idStrArr := make([]string, 0)
	amountStrArr := make([]string, 0)
	for i, id := range ids {
		amount := amounts[i]
		// query balance
		balanceByte, err := ctx.GetStateByte(BalanceMap, string(to)+"_"+id.ToString())
		if err != SUCCESS {
			ctx.ErrorResult("MintBatch get balance error")
			return
		}
		balance, _ := safemath.ParseSafeUint256(string(balanceByte))
		balance, _ = safemath.SafeAdd(balance, amount)

		// set balance
		_ = ctx.PutStateByte(BalanceMap, string(to)+"_"+id.ToString(), []byte(balance.ToString()))
		// set owner only nft valid
		_ = ctx.PutState("owner", id.ToString(), string(to))
		// set the number of token transactions
		_ = ctx.PutStateByte(TokenMap, id.ToString(), []byte("0"))

		idStrArr = append(idStrArr, id.ToString())
		amountStrArr = append(amountStrArr, amount.ToString())
	}
	// emit event
	i, _ := json.Marshal(idStrArr)
	a, _ := json.Marshal(amountStrArr)
	ctx.EmitEvent("MintBatch", string(to), string(i), string(a))
	successNormal()
}
func initContract(adminAddresses []string, uri string) bool {
	ctx := NewSimContext()

	if !address.IsValidAddress(adminAddresses...) {
		ctx.ErrorResult("address error")
	}
	ctx.PutStateByte(Erc1155Map, keyUri, []byte(uri))
	adminAddressByte, _ := json.Marshal(adminAddresses)
	err := ctx.PutStateByte(Erc1155Map, keyAdminAddress, adminAddressByte)
	if err != SUCCESS {
		ctx.ErrorResult("set admin address of contractInfo failed")
		return false
	}
	err = ctx.PutStateByte(Erc1155Map, "userCount", []byte("0"))
	if err != SUCCESS {
		ctx.ErrorResult("set user count of contractInfo failed")
		return false
	}
	ctx.EmitEvent("AlterAdminAddress", adminAddresses...)
	return true
}

func transferCore(from Address, to Address, id *safemath.SafeUint256, value *safemath.SafeUint256) error { // nolint
	ctx := NewSimContext()
	if !value.GTE(safemath.SafeUintOne) {
		return errors.New("error value " + value.ToString())
	}
	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	isApprovedOrOwner, err := isApprovedOrOwner(from, Address(sender), id)
	if err != nil {
		return fmt.Errorf("check isApprovedOrOwner failed, err:%s", err)
	}
	if !isApprovedOrOwner {
		return errors.New("caller is not token owner or approved")
	}
	balanceFromByte, _ := ctx.GetStateByte(BalanceMap, string(from)+"_"+id.ToString())
	balanceToByte, _ := ctx.GetStateByte(BalanceMap, string(to)+"_"+id.ToString())

	balanceFrom, _ := safemath.ParseSafeUint256(string(balanceFromByte))
	balanceTo, _ := safemath.ParseSafeUint256(string(balanceToByte))
	if !balanceFrom.GTE(value) {
		return errors.New("insufficient balance from:" + balanceFrom.ToString())
	}
	balanceFrom, _ = safemath.SafeSub(balanceFrom, value)
	balanceTo, _ = safemath.SafeAdd(balanceTo, value)

	// set balance
	_ = ctx.PutStateByte(BalanceMap, string(from)+"_"+id.ToString(), []byte(balanceFrom.ToString()))
	_ = ctx.PutStateByte(BalanceMap, string(to)+"_"+id.ToString(), []byte(balanceTo.ToString()))
	// set owner only nft valid
	_ = ctx.PutState("owner", id.ToString(), string(to))

	// set tx number
	tokenTxCountByte, _ := ctx.GetStateByte(TokenMap, id.ToString())
	tokenTxCount, _ := safemath.ParseSafeUint256(string(tokenTxCountByte))
	tokenTxCount, _ = safemath.SafeAdd(tokenTxCount, safemath.SafeUintOne)
	_ = ctx.PutStateByte(TokenMap, id.ToString(), []byte(tokenTxCount.ToString()))
	return nil
}

func isApprovedOrOwner(owner, sender Address, tokenId *safemath.SafeUint256) (bool, error) {
	// check owner
	if owner == sender {
		return true, nil
	}

	// check operatorApprove
	resp := isApprovedForAll(owner, sender)
	if !resp {
		return false, errors.New("Is not Approved For All")
	} else {
		return true, nil
	}
	return false, nil
}

// OwnerOf 查询owner
//
//go:wasmexport OwnerOf
func OwnerOf() {
	ctx := NewSimContext()
	idStr, _ := ctx.ArgString(paramId)
	tokenId, ok := safemath.ParseSafeUint256(idStr)
	if !ok {
		ctx.ErrorResult("param error")
	}
	val, _ := ctx.GetState("owner", tokenId.ToString())
	addr := Address(val)
	ctx.SuccessResult(string(addr))
}

func ownerOf(tokenId *safemath.SafeUint256) Address {
	ctx := NewSimContext()
	val, _ := ctx.GetState("owner", tokenId.ToString())
	addr := Address(val)
	ctx.SuccessResult(string(addr))
	return addr
}

func senderIsAdmin() bool {
	ctx := NewSimContext()
	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	adminAddressByte, err := ctx.GetStateByte(Erc1155Map, keyAdminAddress)
	if len(adminAddressByte) == 0 || err != SUCCESS {
		ctx.ErrorResult("Get adminAddress failed")
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

func successNormal() {
	ctx := NewSimContext()
	ctx.SuccessResult("ok")
}
func successData(data string) {
	ctx := NewSimContext()
	ctx.SuccessResult(data)
}

func main() {

}
