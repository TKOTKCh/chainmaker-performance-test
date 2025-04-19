/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chainmaker.org/chainmaker/contract-utils/safemath"
	"encoding/json"
	"errors"
	"fmt"

	"strings"

	"chainmaker.org/chainmaker/contract-utils/address"
	"chainmaker.org/chainmaker/contract-utils/standard"
)

const (
	paramAdmin        = "admin"
	paramCategoryName = "categoryName"
	paramCategoryURI  = "categoryURI"
	paramCategory     = "category"
	paramFrom         = "from"
	paramTo           = "to"
	paramTokenId      = "tokenId"
	paramTokenIds     = "tokenIds"
	paramMetadata     = "metadata"
	paramTokens       = "tokens"
	paramIsApproval   = "isApproval"
	paramStandardName = "standardName"

	categoryMapName                = "categoryMap"
	categoryTotalSupplyMapName     = "categoryTotalSupplyMap"
	balanceInfoMapName             = "balanceInfoMap"
	accountMapName                 = "accountInfoMap"
	tokenOwnerMapName              = "tokenOwnerMap"
	tokenInfoMapName               = "tokenInfoMap"
	tokenApprovalMapName           = "tokenApprovalMap"
	tokenApprovalForAllMapName     = "tokenApprovalForAllMap"
	tokenApprovalByCategoryMapName = "tokenApprovalByCategoryMap"

	adminStoreKey        = "admin"
	metadataStoreKey     = "metadata"
	categoryNameStoreKey = "categoryName"
	totalSupplyStoreKey  = "TotalSupply"

	nfaStandardName = standard.ContractStandardNameCMNFA
)

//go:wasmexport init_contract
func InitContract() {
	ctx := NewSimContext()
	err1 := updateNFAInfo(true)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("Init contract err: %s", err1.Error()))
		return
	}
	ctx.SuccessResult("Init contract success")
}

//go:wasmexport upgrade
func Upgrade() {
	ctx := NewSimContext()
	err1 := updateNFAInfo(false)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("Upgrade contract err: %s", err1.Error()))
		return
	}
	ctx.SuccessResult("Upgrade contract success")
}

//go:wasmexport Standards
func standardsCore() {
	ctx := NewSimContext()
	standards := Standards()
	standardsBytes, err1 := json.Marshal(standards)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	ctx.SuccessResult(string(standardsBytes))
}

//go:wasmexport SupportStandard
func supportStandardCore() {
	ctx := NewSimContext()
	standardName, _ := ctx.ArgString(paramStandardName)
	isSupport := SupportStandard(standardName)
	if isSupport {
		ctx.SuccessResult(standard.TrueString)
		return
	} else {
		ctx.SuccessResult(standard.FalseString)
		return
	}
}

//go:wasmexport Mint
func mintCore() {
	ctx := NewSimContext()
	to, _ := ctx.ArgString(paramTo)
	tokenId, _ := ctx.ArgString(paramTokenId)
	categoryName, _ := ctx.ArgString(paramCategoryName)
	metadata, _ := ctx.Arg(paramMetadata)
	err1 := verifyToken(to, tokenId, categoryName)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	adminByte, err := ctx.GetStateFromKey(adminStoreKey)
	if err != SUCCESS {
		ctx.ErrorResult(fmt.Sprintf("get admin failed, err:%s", err))
		return
	}
	admin := string(adminByte)
	if sender != admin {
		ctx.ErrorResult("only admin can mint tokens")
		return
	}
	err1 = Mint(to, tokenId, categoryName, metadata)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	//ctx.EmitEvent("Mint", strings.Join([]string{to, tokenId, categoryName, string(metadata)}, ",")
	ctx.SuccessResult("Mint success")
}

//go:wasmexport MintBatch
func mintBatchCore() {
	ctx := NewSimContext()
	tokens, _ := ctx.Arg(paramTokens)
	if len(tokens) == 0 {
		ctx.ErrorResult("invalid tokens parameter")
		return
	}
	nfas := make([]standard.NFA, 0)
	err1 := json.Unmarshal(tokens, &nfas)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("unmarshal tokens failed, err:%s", err1))
		return
	}
	for _, nfa := range nfas {
		err1 = verifyToken(nfa.To, nfa.TokenId, nfa.CategoryName)
		if err1 != nil {
			ctx.ErrorResult(err1.Error())
			return
		}
	}
	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	adminByte, err := ctx.GetStateFromKey(adminStoreKey)
	if err != SUCCESS {
		ctx.ErrorResult(fmt.Sprintf("get admin failed, err:%s", err))
		return
	}
	admin := string(adminByte)
	if sender != admin {
		ctx.ErrorResult("only admin can mint tokens")
		return
	}
	err1 = MintBatch(nfas)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	ctx.SuccessResult("MintBatch success")
}

//go:wasmexport SetApproval
func setApprovalCore() {
	ctx := NewSimContext()
	tokenId, _ := ctx.ArgString(paramTokenId)
	to, _ := ctx.ArgString(paramTo)
	isApproval, _ := ctx.ArgString(paramIsApproval)

	if len(tokenId) == 0 {
		ctx.ErrorResult("invalid tokenId")
		return
	}
	if !address.IsValidAddress(to) || address.IsZeroAddress(to) {
		ctx.ErrorResult("invalid to address")
		return
	}
	if isApproval != standard.TrueString && isApproval != standard.FalseString {
		ctx.ErrorResult("isApproval only can be true or false")
		return
	}

	// Get token owner
	owner, err1 := OwnerOf(tokenId)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	// Check sender is owner
	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	if sender != owner {
		ctx.ErrorResult("only owner can set approval")
		return
	}

	err1 = SetApproval(owner, to, tokenId, isApproval == standard.TrueString)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	ctx.SuccessResult("SetApproval success")
}

//go:wasmexport SetApprovalForAll
func setApprovalForAllCore() {
	ctx := NewSimContext()
	to, _ := ctx.ArgString(paramTo)
	isApproval, _ := ctx.ArgString(paramIsApproval)

	if !address.IsValidAddress(to) || address.IsZeroAddress(to) {
		ctx.ErrorResult("invalid to address")
		return
	}
	if isApproval != standard.TrueString && isApproval != standard.FalseString {
		ctx.ErrorResult("invalid isApproval value")
		return
	}

	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	err1 := SetApprovalForAll(sender, to, isApproval == standard.TrueString)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("set approve for all failed, err:%s", err1))
		return
	}
	ctx.SuccessResult("SetApprovalForAll success")
}

//go:wasmexport TransferFrom
func transferFromCore() {
	ctx := NewSimContext()
	from, _ := ctx.ArgString(paramFrom)
	to, _ := ctx.ArgString(paramTo)
	tokenId, _ := ctx.ArgString(paramTokenId)

	// Validate addresses
	if !address.IsValidAddress(from) || address.IsZeroAddress(from) {
		ctx.ErrorResult("invalid from address")
		return
	}
	if !address.IsValidAddress(to) || address.IsZeroAddress(to) {
		ctx.ErrorResult("invalid to address")
		return
	}
	if len(tokenId) == 0 {
		ctx.ErrorResult("invalid tokenId")
		return
	}
	// Check sender permission
	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	IsApprovedOrOwner, err1 := isApprovedOrOwner(sender, tokenId)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("check isApprovedOrOwner failed, err:%s", err1))
		return
	}
	if !IsApprovedOrOwner {
		ctx.ErrorResult("sender is not token owner or approved")
		return
	}

	err1 = TransferFrom(from, to, tokenId)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	ctx.SuccessResult("TransferFrom success")
}

// Batch transfer functions
//
//go:wasmexport TransferFromBatch
func transferFromBatchCore() {
	ctx := NewSimContext()
	from, _ := ctx.ArgString(paramFrom)
	to, _ := ctx.ArgString(paramTo)
	tokenIdsBytes, _ := ctx.Arg(paramTokenIds)
	if address.IsZeroAddress(from) || !address.IsValidAddress(from) {
		ctx.ErrorResult("invalid from address")
		return
	}
	if address.IsZeroAddress(to) || !address.IsValidAddress(to) {
		ctx.ErrorResult("invalid to address")
		return
	}
	if len(tokenIdsBytes) == 0 {
		ctx.ErrorResult("invalid tokenIdsBytes")
		return
	}
	tokenIds := make([]string, 0)
	err1 := json.Unmarshal(tokenIdsBytes, &tokenIds)
	if err1 != nil {
		ctx.ErrorResult("invalid tokenIds unmarshal")
		return
	}

	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	if len(tokenIds) == 0 {
		ctx.ErrorResult("invalid tokenIds")
		return
	}
	for _, tokenId := range tokenIds {
		if len(tokenId) == 0 {
			ctx.ErrorResult("invalid tokenId")
			return
		}

		IsApprovedOrOwner, err1 := isApprovedOrOwner(sender, tokenId)
		if err1 != nil {
			ctx.ErrorResult(err1.Error())
			return
		}
		if !IsApprovedOrOwner {
			ctx.ErrorResult("only owner or approved account can transfer token")
			return
		}
	}
	err1 = TransferFromBatch(from, to, tokenIds)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	ctx.SuccessResult("TransferFromBatch success")
}

//go:wasmexport OwnerOf
func ownerOfCore() {
	ctx := NewSimContext()
	tokenId, _ := ctx.ArgString(paramTokenId)
	if len(tokenId) == 0 {
		ctx.ErrorResult("invalid tokenId")
		return
	}
	owner, err1 := OwnerOf(tokenId)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	ctx.SuccessResult(owner)
}

//go:wasmexport TokenURI
func tokenURICore() {
	ctx := NewSimContext()
	tokenId, _ := ctx.ArgString(paramTokenId)

	if len(tokenId) == 0 {
		ctx.ErrorResult("invalid tokenId")
		return
	}
	uri, err1 := TokenURI(tokenId)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	ctx.SuccessResult(uri)
}

//go:wasmexport SetApprovalByCategory
func setApprovalByCategoryCore() {
	ctx := NewSimContext()
	to, _ := ctx.ArgString(paramTo)
	categoryName, _ := ctx.ArgString(paramCategoryName)
	isApproval, _ := ctx.ArgString(paramIsApproval)

	// Validate inputs
	if !address.IsValidAddress(to) || address.IsZeroAddress(to) {
		ctx.ErrorResult("invalid to address")
		return
	}
	if isApproval != standard.TrueString && isApproval != standard.FalseString {
		ctx.ErrorResult("isApprove only can be true or false")
		return
	}

	category, err1 := GetCategoryByName(categoryName)
	if err1 != nil || category == nil {
		ctx.ErrorResult("invalid categoryName")
		return
	}

	// Get sender and set approval
	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	err1 = SetApprovalByCategory(sender, to, categoryName, isApproval == standard.TrueString)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	ctx.SuccessResult("setApprovalByCategoryCore success")
}

//go:wasmexport CreateOrSetCategory
func createOrSetCategoryCore() {
	ctx := NewSimContext()
	categoryBytes, _ := ctx.Arg(paramCategory)
	category := &standard.Category{}
	err1 := json.Unmarshal(categoryBytes, &category)
	if err1 != nil {
		ctx.ErrorResult(fmt.Sprintf("invalid category format %s", err1.Error()))
		return
	}
	err1 = CreateOrSetCategory(category)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}

	ctx.SuccessResult("CreateOrSetCategory success")
}

//go:wasmexport Burn
func burnCore() {
	ctx := NewSimContext()
	tokenId, _ := ctx.ArgString(paramTokenId)
	if len(tokenId) == 0 {
		ctx.ErrorResult("invalid tokenId")
		return
	}
	// Check ownership
	sender, _ := ctx.GetSenderPk()
	if len(sender) > 40 {
		sender = sender[:40]
	}
	approved, err := isApprovedOrOwner(sender, tokenId)
	if err != nil {
		ctx.ErrorResult(fmt.Sprintf("check approved failed, err:%s", err))
	}
	if !approved {
		ctx.ErrorResult("only owner or approved user can Burn the token")
	}
	err = Burn(tokenId)
	if err != nil {
		ctx.ErrorResult(err.Error())
	}
	ctx.SuccessResult("Burn success")
}

//go:wasmexport GetCategoryByName
func getCategoryByNameCore() {
	ctx := NewSimContext()
	categoryName, _ := ctx.ArgString(paramCategoryName)
	if len(categoryName) == 0 {
		ctx.ErrorResult("invalid categoryName")
		return
	}
	category, err := GetCategoryByName(categoryName)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	if category == nil {
		ctx.ErrorResult("category not exist")
		return
	}
	categoryBytes, err := json.Marshal(category)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult(string(categoryBytes))
}

//go:wasmexport GetCategoryByTokenId
func getCategoryByTokenIdCore() {
	ctx := NewSimContext()
	tokenId, _ := ctx.ArgString(paramTokenId)
	if len(tokenId) == 0 {
		ctx.ErrorResult("invalid categoryName")
		return
	}
	category, err := GetCategoryByTokenId(tokenId)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}

	categoryBytes, err := json.Marshal(category)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult(string(categoryBytes))
}

//go:wasmexport TotalSupply
func totalSupplyCore() {
	ctx := NewSimContext()
	totalSupply, err := TotalSupply()
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult(totalSupply.ToString())
}

//go:wasmexport TotalSupplyOfCategory
func totalSupplyOfCategoryCore() {
	ctx := NewSimContext()
	categoryName, _ := ctx.ArgString(paramCategoryName)

	if len(categoryName) == 0 {
		ctx.ErrorResult("invalid categoryName")
		return
	}
	totalSupply, err := TotalSupplyOfCategory(categoryName)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult(totalSupply.ToString())
}

//go:wasmexport BalanceOf
func balanceOfCore() {
	ctx := NewSimContext()
	account, _ := ctx.ArgString("account")
	if len(account) == 0 {
		ctx.ErrorResult("Param account should not be empty")
		return
	}
	if !address.IsValidAddress(account) {
		ctx.ErrorResult("invalid account address")
		return
	}
	if address.IsZeroAddress(account) {
		ctx.ErrorResult("address zero is not a valid owner")
		return
	}
	balance, err := BalanceOf(account)
	if err != nil {
		ctx.ErrorResult(err.Error())
	}

	ctx.SuccessResult(balance.ToString())
}

//go:wasmexport AccountTokens
func accountTokensCore() {
	ctx := NewSimContext()
	account, _ := ctx.ArgString("account")

	if !address.IsValidAddress(account) || address.IsZeroAddress(account) {
		ctx.ErrorResult("invalid account")
		return
	}
	tokens, err := AccountTokens(account)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ats := &standard.AccountTokens{
		Account: account,
		Tokens:  tokens,
	}
	var atsBytes []byte
	atsBytes, err = json.Marshal(ats)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	ctx.SuccessResult(string(atsBytes))
}

//go:wasmexport TokenMetadata
func tokenMetadataCore() {
	ctx := NewSimContext()
	tokenId, _ := ctx.ArgString(paramTokenId)

	if len(tokenId) == 0 {
		ctx.ErrorResult("invalid tokenId")
		return
	}
	metadata, err := TokenMetadata(tokenId)
	if err != nil {
		ctx.ErrorResult(err.Error())
	}
	ctx.SuccessResult(string(metadata))
}

func updateNFAInfo(isInitContract bool) error {
	ctx := NewSimContext()
	categoryName, _ := ctx.Arg(paramCategoryName)
	categoryURI, _ := ctx.ArgString(paramCategoryURI)
	admin, _ := ctx.ArgString(paramAdmin)

	if len(categoryName) > 0 {
		err := ctx.PutState(categoryMapName, string(categoryName), categoryURI)
		if err != SUCCESS {
			return fmt.Errorf("categoryMap set default category failed, err:%s", err)
		}
	}

	if len(admin) > 0 {
		err := ctx.PutStateFromKey(adminStoreKey, admin)
		if err != SUCCESS {
			return fmt.Errorf("save admin failed, err:%s", err)
		}
	} else if isInitContract {
		origin, err := ctx.GetSenderPk()
		if len(origin) > 40 {
			origin = origin[:40]
		}
		if err != SUCCESS {
			return errors.New("get sender failed")
		}
		err = ctx.PutStateFromKey(adminStoreKey, origin)
		if err != SUCCESS {
			return fmt.Errorf("save admin failed, err:%s", err)
		}
	}
	return nil
}

// Standards returns standard strings which contain "CMNFA"
func Standards() (standards []string) {
	return []string{nfaStandardName}
}

// SupportStandard returns true if standardName equals "CMNFA"
func SupportStandard(standardName string) bool {
	return standardName == nfaStandardName
}
func minted(tokenId string) (bool, error) {
	owner, err := OwnerOf(tokenId)
	if err != nil {
		return false, err
	}
	return address.IsValidAddress(owner) && !address.IsZeroAddress(owner), nil
}
func Mint(to, tokenId, categoryName string, metadata []byte) error {
	err1 := increaseBalanceByOne(to)
	if err1 != nil {
		return err1
	}
	err1 = increaseTotalSupplyByOne()
	if err1 != nil {
		return err1
	}
	err1 = increaseTotalSupplyOfCategoryByOne(categoryName)
	if err1 != nil {
		return err1
	}
	err1 = setTokenOwner(to, tokenId)
	if err1 != nil {
		return err1
	}

	err1 = setAccountToken(address.ZeroAddr, to, tokenId)
	if err1 != nil {
		return err1
	}

	err1 = setTokenInfo(tokenId, categoryName, metadata)
	if err1 != nil {
		return err1
	}
	EmitMintEvent(to, tokenId, categoryName, string(metadata))
	return nil
}
func MintBatch(tokens []standard.NFA) error {
	for _, token := range tokens {
		err1 := Mint(token.To, token.TokenId, token.CategoryName, token.Metadata)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func verifyToken(to, tokenId, categoryName string) error {
	ctx := NewSimContext()
	if !address.IsValidAddress(to) || address.IsZeroAddress(to) {
		return fmt.Errorf("mint to invalid address")
	}
	if len(tokenId) == 0 {
		return fmt.Errorf("invalid tokenId")
	}
	minted, err1 := minted(tokenId)
	if err1 != nil {
		return err1
	}
	if minted {
		return fmt.Errorf("duplicated token")
	}
	if len(categoryName) == 0 {
		return fmt.Errorf("invalid category name")
	}
	exist, err := ctx.GetStateByte(categoryMapName, string(categoryName))
	if err != SUCCESS || len(exist) == 0 {
		return fmt.Errorf("category not exist")
	}
	return nil
}
func SetApproval(owner, to, tokenId string, isApproval bool) error {
	ctx := NewSimContext()
	var err ResultCode
	if isApproval {
		err = ctx.PutStateByte(tokenApprovalMapName, tokenId, []byte(to))
	} else {
		err = ctx.DeleteState(tokenApprovalMapName, tokenId)
	}
	if err != SUCCESS {
		return fmt.Errorf("set approve failed, err:%s", err)
	}

	EmitSetApprovalEvent(owner, to, tokenId, isApproval)

	return nil
}
func SetApprovalForAll(owner, to string, isApproval bool) error {
	ctx := NewSimContext()
	var err ResultCode
	if isApproval {
		err = ctx.PutStateByte(tokenApprovalForAllMapName, owner+"_"+to, []byte(standard.TrueString))
	} else {
		err = ctx.DeleteState(tokenApprovalForAllMapName, owner+"_"+to)
	}
	if err != SUCCESS {
		return fmt.Errorf("set operator approve failed, err:%s", err)
	}

	EmitSetApprovalForAllEvent(owner, to, isApproval)

	return nil
}
func getApproved(tokenId string) (bool, error) {
	ctx := NewSimContext()
	approveTo, err := ctx.GetStateByte(tokenApprovalMapName, tokenId)
	if err != SUCCESS {
		return false, errors.New("get approve info failed")
	}

	return !address.IsZeroAddress(string(approveTo)) && address.IsValidAddress(string(approveTo)), nil
}

func isApprovedForAll(owner, sender string) (bool, error) {
	ctx := NewSimContext()
	val, err := ctx.GetStateByte(tokenApprovalForAllMapName, owner+"_"+sender)
	if err != SUCCESS {
		return false, fmt.Errorf("get approved val from approve info failed, err:%s", err)
	}
	return string(val) == standard.TrueString, nil
}

func isApprovedByCategory(owner, sender, categoryName string) (bool, error) {
	ctx := NewSimContext()
	val, err := ctx.GetStateByte(tokenApprovalByCategoryMapName, owner+"_"+sender+"_"+categoryName)
	if err != SUCCESS {
		return false, fmt.Errorf("get approved val from approve info failed, err:%s", err)
	}
	return string(val) == standard.TrueString, nil
}
func OwnerOf(tokenId string) (account string, err1 error) {
	ctx := NewSimContext()
	owner, _ := ctx.GetState(tokenOwnerMapName, tokenId)
	if len(owner) == 0 {
		return "", nil
	}
	return owner, nil
}
func isApprovedOrOwner(sender string, tokenId string) (bool, error) {
	// check owner
	owner, err1 := OwnerOf(tokenId)
	if err1 != nil {
		return false, err1
	}
	if owner == sender {
		return true, nil
	}

	// check tokenApprove
	approved, err1 := getApproved(tokenId)
	if err1 != nil {
		return false, err1
	}
	if approved {
		return true, nil
	}
	approved, err1 = isApprovedForAll(owner, sender)
	if err1 != nil {
		return false, err1
	}
	if approved {
		return true, nil
	}
	categoryName, err1 := getCategoryNameByTokenId(tokenId)
	if err1 != nil {
		return false, err1
	}
	approved, err1 = isApprovedByCategory(owner, sender, categoryName)
	if err1 != nil {
		return false, err1
	}

	return approved, nil
}

// 以下是辅助函数
func increaseBalanceByOne(account string) error {
	ctx := NewSimContext()
	originTokenCount, err1 := getBalance(account)
	if err1 != nil {
		return fmt.Errorf("get token count failed, err:%s", err1)
	}
	newTokenCount, ok := safemath.SafeAdd(originTokenCount, safemath.SafeUintOne)
	if !ok {
		return fmt.Errorf("balance count of from is overflow")
	}
	err := ctx.PutStateByte(balanceInfoMapName, account, []byte(newTokenCount.ToString()))
	if err != SUCCESS {
		return fmt.Errorf("put balance info failed, err:%s", err)
	}
	return nil
}
func increaseTotalSupplyByOne() error {
	ctx := NewSimContext()
	currentTotalSupplyBytes, err := ctx.GetStateFromKey(totalSupplyStoreKey)
	currentTotalSupplyStr := string(currentTotalSupplyBytes)
	if err != SUCCESS {
		return errors.New("get current total supply failed")
	}
	currentTotalSupply, ok := safemath.ParseSafeUint256(currentTotalSupplyStr)
	if !ok {
		return fmt.Errorf("parse current total supply failed")
	}
	newTotalSupply, ok := safemath.SafeAdd(currentTotalSupply, safemath.SafeUintOne)
	if !ok {
		return fmt.Errorf("total supply too big")
	}
	err = ctx.PutStateFromKey(totalSupplyStoreKey, newTotalSupply.ToString())
	if err != SUCCESS {
		return errors.New("put total supply failed")
	}
	return nil
}

func increaseTotalSupplyOfCategoryByOne(categoryName string) error {
	ctx := NewSimContext()
	currentTotalSupplyOfCategoryBytes, err := ctx.GetState(categoryTotalSupplyMapName, categoryName)
	if err != SUCCESS {
		return errors.New("get total supply of category failed")
	}
	currentTotalSupplyOfCategory, ok := safemath.ParseSafeUint256(string(currentTotalSupplyOfCategoryBytes))
	if !ok {
		return fmt.Errorf("invalid TotalSupply of category in store")
	}
	newTotalSupplyOfCategory, ok := safemath.SafeAdd(currentTotalSupplyOfCategory, safemath.SafeUintOne)
	if !ok {
		return fmt.Errorf("TotalSupply of category too big")
	}
	err = ctx.PutState(categoryTotalSupplyMapName, categoryName, newTotalSupplyOfCategory.ToString())
	if err != SUCCESS {
		return errors.New("put total supply of category failed")
	}
	return nil
}

func decreaseBalanceByOne(account string) error {
	ctx := NewSimContext()
	originTokenCount, err1 := getBalance(account)
	if err1 != nil {
		return fmt.Errorf("get token count failed, err:%s", err1)
	}
	newTokenCount, ok := safemath.SafeSub(originTokenCount, safemath.SafeUintOne)
	if !ok {
		return fmt.Errorf("balance count of from is overflow")
	}
	err := ctx.PutStateByte(balanceInfoMapName, account, []byte(newTokenCount.ToString()))
	if err != SUCCESS {
		return fmt.Errorf("put balance info failed, err:%s", err)
	}
	return nil
}

func setTokenOwner(to string, tokenId string) error {
	ctx := NewSimContext()
	err := ctx.PutState(tokenOwnerMapName, tokenId, to)
	if err != SUCCESS {
		return fmt.Errorf("set token owner failed, err:%s", err)
	}
	return nil
}
func EmitMintEvent(to, tokenId, categoryName, metadata string) {
	ctx := NewSimContext()
	ctx.EmitEvent("Mint", address.ZeroAddr, to, tokenId, categoryName, metadata)
}

// EmitSetApprovalEvent emits SetApproval event
func EmitSetApprovalEvent(owner, to, tokenId string, isApproval bool) {
	ctx := NewSimContext()
	if isApproval {
		ctx.EmitEvent("SetApproval", owner, to, tokenId, standard.TrueString)
	} else {
		ctx.EmitEvent("SetApproval", owner, to, tokenId, standard.FalseString)
	}
}

// EmitBurnEvent emits Burn event
func EmitBurnEvent(tokenId string) {
	ctx := NewSimContext()
	ctx.EmitEvent("Burn", tokenId)
}

// EmitSetApprovalByCategoryEvent emits SetApprovalByCategory event
func EmitSetApprovalByCategoryEvent(owner, to, categoryName string, isApproval bool) {
	ctx := NewSimContext()
	if isApproval {
		ctx.EmitEvent("SetApprovalByCategory", owner, to, categoryName, standard.TrueString)
	} else {
		ctx.EmitEvent("SetApprovalByCategory", owner, to, categoryName, standard.FalseString)
	}
}

// EmitCreateOrSetCategoryEvent emits CreateOrSetCategory event
func EmitCreateOrSetCategoryEvent(categoryName, categoryURI string) {
	ctx := NewSimContext()
	ctx.EmitEvent("CreateOrSetCategory", categoryName, categoryURI)
}

// EmitSetApprovalForAllEvent emits SetApprovalForAll event
func EmitSetApprovalForAllEvent(owner, to string, isApproval bool) {
	ctx := NewSimContext()
	if isApproval {
		ctx.EmitEvent("SetApprovalForAll", owner, to, standard.TrueString)
	} else {
		ctx.EmitEvent("SetApprovalForAll", owner, to, standard.FalseString)
	}
}

// EmitTransferFromEvent emits TransferFrom event
func EmitTransferFromEvent(from, to, tokenId string) {
	ctx := NewSimContext()
	ctx.EmitEvent("TransferFrom", from, to, tokenId)
}
func setAccountToken(from, to string, tokenId string) error {
	ctx := NewSimContext()
	err := ctx.PutStateByte(accountMapName, to+"_"+tokenId, []byte(standard.TrueString))
	if err != SUCCESS {
		return fmt.Errorf("setAccountToken failed, err:%s", err)
	}
	if address.IsZeroAddress(from) {
		return nil
	}
	err = ctx.DeleteState(accountMapName, to+"_"+tokenId)
	if err != SUCCESS {
		return errors.New("delete account token failed")
	}
	return nil
}

func setTokenInfo(tokenId, categoryName string, metadata []byte) error {
	ctx := NewSimContext()
	err := ctx.PutStateByte(tokenInfoMapName, tokenId+"_"+categoryNameStoreKey, []byte(categoryName))
	if err != SUCCESS {
		return fmt.Errorf("set category name of token info failed, err:%s", err)
	}
	if len(metadata) > 0 {
		err = ctx.PutStateByte(tokenInfoMapName, tokenId+"_"+metadataStoreKey, metadata)
		if err != SUCCESS {
			return fmt.Errorf("set metadata of token info failed, err:%s", err)
		}
	}

	return nil
}
func getBalance(account string) (balance *safemath.SafeUint256, err1 error) {
	ctx := NewSimContext()
	balanceBytes, err := ctx.GetStateByte(balanceInfoMapName, account)
	if err != SUCCESS {
		return nil, fmt.Errorf("get balance failed, err:%s", err)
	}
	balance, ok := safemath.ParseSafeUint256(string(balanceBytes))
	if !ok {
		return nil, fmt.Errorf("balance bytes invalid")
	}

	return balance, nil
}

// SetApprovalByCategory approve or cancel approve tokens of category to 'to' account. Optional.
// @param to, destination address approve to. Obligatory.
// @categoryName, the category of tokens. Obligatory.
// @isApproval, to approve or to cancel approve. Obligatory.
// @return error, the error msg if some error occur.
// @event, topic: 'SetApprovalByCategory'; data: to, categoryName, isApproval
func SetApprovalByCategory(owner, to, categoryName string, isApproval bool) error {
	ctx := NewSimContext()
	var err ResultCode
	if isApproval {
		err = ctx.PutStateByte(tokenApprovalByCategoryMapName, owner+"_"+to+"_"+categoryName, []byte(standard.TrueString))
	} else {
		err = ctx.DeleteState(tokenApprovalByCategoryMapName, owner+"_"+to+"_"+categoryName)
	}

	if err != SUCCESS {
		return errors.New("set approval ByCategory info failed")
	}

	EmitSetApprovalByCategoryEvent(owner, to, categoryName, isApproval)

	return nil
}

// CreateOrSetCategory create ore reset a category. Optional.
// @param categoryName, the category name. Obligatory.
// @param categoryURI, the category uri. Obligatory.
// @return error, the error msg if some error occur.
// @event, topic: 'CreateOrSetCategory'; data: category
func CreateOrSetCategory(category *standard.Category) error {
	ctx := NewSimContext()

	err := ctx.PutStateByte(categoryMapName, category.CategoryName, []byte(category.CategoryURI))
	EmitCreateOrSetCategoryEvent(category.CategoryName, category.CategoryURI)
	if err != SUCCESS {
		return errors.New("CreateOrSetCategory failed")
	}
	return nil
}

// Burn burn a token
// @param tokenId
// @event, topic: 'Burn'; data: tokenId
func Burn(tokenId string) error {
	ctx := NewSimContext()
	err := ctx.DeleteState(tokenApprovalMapName, tokenId)
	if err != SUCCESS {
		return fmt.Errorf("delete token approve failed, err:%s", err)
	}

	// update "owner" balance count
	owner, err1 := OwnerOf(tokenId)
	if err1 != nil {
		return err1
	}
	err1 = decreaseBalanceByOne(owner)
	if err1 != nil {
		return err1
	}

	// update token owner
	err1 = setTokenOwner(address.ZeroAddr, tokenId)
	if err1 != nil {
		return err1
	}

	err1 = setAccountToken(owner, address.ZeroAddr, tokenId)
	if err1 != nil {
		return err1
	}

	EmitBurnEvent(tokenId)

	return nil
}

// GetCategoryByName get specific category by name. Optional.
// @param categoryName, the name of the category. Obligatory.
// @return category, the category returned.
// @return err, the error msg if some error occur.
func GetCategoryByName(categoryName string) (*standard.Category, error) {
	ctx := NewSimContext()
	uri, err := ctx.GetState(categoryMapName, categoryName)
	if err != SUCCESS {
		return nil, errors.New("GetCategoryByName failed")
	}
	if len(uri) == 0 {
		return nil, nil
	}
	return &standard.Category{
		CategoryName: categoryName,
		CategoryURI:  string(uri),
	}, nil
}

// GetCategoryByTokenId get a specific category by tokenId. Optional.
// @param tokenId, the names of category to be queried. Obligatory.
// @return category, the result queried.
// @return err, the error msg if some error occur.
func GetCategoryByTokenId(tokenId string) (category *standard.Category, err error) {
	categoryName, err1 := getCategoryNameByTokenId(tokenId)
	if err1 != nil {
		return nil, err1
	}
	category, err1 = GetCategoryByName(categoryName)
	if err1 != nil || category == nil {
		return nil, err1
	}
	return category, nil
}
func getCategoryNameByTokenId(tokenId string) (string, error) {
	ctx := NewSimContext()
	categoryName, err := ctx.GetState(tokenInfoMapName, tokenId+"_"+categoryNameStoreKey)
	if err != SUCCESS {
		return "", errors.New("getCategoryNameByTokenId failed")
	}
	return string(categoryName), nil
}

// TotalSupply get total token supply of this contract. Obligatory.
// @return TotalSupply, the total token supply value returned.
// @return err, the error msg if some error occur.
func TotalSupply() (*safemath.SafeUint256, error) {
	ctx := NewSimContext()
	totalSupplyByte, err := ctx.GetStateFromKey(totalSupplyStoreKey)
	totalSupplyStr := string(totalSupplyByte)
	if err != SUCCESS {
		return nil, errors.New("TotalSupply failed")
	}
	totalSupply, ok := safemath.ParseSafeUint256(totalSupplyStr)
	if !ok {
		return nil, fmt.Errorf("the total supply in store is invalid")
	}
	return totalSupply, nil
}

// TotalSupplyOfCategory get total token supply of the category. Obligatory.
// @param category, the category of tokens. Obligatory.
// @return TotalSupply, the total token supply value returned.
// @return err, the error msg if some error occur.
func TotalSupplyOfCategory(categoryName string) (*safemath.SafeUint256, error) {
	ctx := NewSimContext()
	totalSupplyOfCategoryStr, err := ctx.GetState(categoryTotalSupplyMapName, categoryName)
	if err != SUCCESS {
		return nil, errors.New("TotalSupplyOfCategory failed")
	}
	totalSupplyOfCategory, ok := safemath.ParseSafeUint256(string(totalSupplyOfCategoryStr))
	if !ok {
		return nil, fmt.Errorf("invalid TotalSupply of category in store")
	}
	return totalSupplyOfCategory, nil
}

// BalanceOf get total token number of the account. Optional
// @param account, the account which will be queried. Obligatory.
// @return balance, the token number of the account.
// @return err, the error msg if some error occur.
func BalanceOf(account string) (*safemath.SafeUint256, error) {
	balanceCount, err1 := getBalance(account)
	if err1 != nil {
		return nil, fmt.Errorf("Get balance failed, err:%s", err1)
	}
	return balanceCount, nil
}

// AccountTokens get the token list of the account. Optional
// @param account, the account which will be queried. Obligatory.
// @return tokenId, the list of tokenId.
// @return err, the error msg if some error occur.
func AccountTokens(account string) ([]string, error) {
	ctx := NewSimContext()
	rs, err := ctx.NewIteratorPrefixWithKeyField(accountMapName, account)
	if err != SUCCESS {
		return nil, fmt.Errorf("new store map iterator of project info failed, err: %s", err)
	}

	tokens := make([]string, 0)
	for {
		if !rs.HasNext() {
			break
		}
		_, field, _, err := rs.Next()
		if err != SUCCESS {
			ctx.ErrorResult(fmt.Sprintf("iterator next failed, err: %s", err))
			return nil, fmt.Errorf("iterator next failed, err: %s", err)
		}
		itemId := strings.Split(field, "_")[1]
		if len(itemId) == 0 {
			return nil, fmt.Errorf("invalid itemId")
		}
		tokens = append(tokens, itemId)
	}
	return tokens, nil
}

// TokenMetadata get the metadata of a token. Optional.
// @param tokenId, tokenId which will be queried.
// @return metadata, the metadata of the token.
// @return err, the error msg if some error occur.
func TokenMetadata(tokenId string) ([]byte, error) {
	ctx := NewSimContext()

	metadata, err := ctx.GetStateByte(tokenInfoMapName, tokenId+"_"+metadataStoreKey)
	if err != SUCCESS {
		return nil, fmt.Errorf("set metadata of erc721Info failed, err:%s", err)
	}

	return metadata, nil
}

// TransferFrom transfer single token after approve. Obligatory.
// @param from, owner account of token. Obligatory.
// @param to, destination account transferred to. Obligatory.
// @param tokenId, the token being transferred. Obligatory.
// @return error, the error msg if some error occur.
// @event, topic: 'TransferFrom'; data: from, to, tokenId
func TransferFrom(from, to, tokenId string) error {
	// delete token approve
	ctx := NewSimContext()
	err := ctx.DeleteState(tokenApprovalMapName, tokenId)
	if err != SUCCESS {
		return fmt.Errorf("delete token approve failed, err:%s", err)
	}

	// update "from" balance count
	err1 := decreaseBalanceByOne(from)
	if err1 != nil {
		return err1
	}

	// update "to" balance count
	err1 = increaseBalanceByOne(to)
	if err1 != nil {
		return err1
	}

	// update token owner
	err1 = setTokenOwner(to, tokenId)
	if err1 != nil {
		return err1
	}

	err1 = setAccountToken(from, to, tokenId)
	if err1 != nil {
		return err1
	}

	EmitTransferFromEvent(from, to, tokenId)

	return nil
}

// TransferFromBatch transfer tokens after approve. Obligatory.
// @param from, owner account of token. Obligatory.
// @param to, destination account transferred to. Obligatory.
// @param tokenIds, the tokens being transferred. Obligatory.
// @return error, the error msg if some error occur.
// @event, topic: 'TransferFromBatch'; data: from, to, tokenIds
func TransferFromBatch(from, to string, tokenIds []string) error {
	for _, tokenId := range tokenIds {
		err1 := TransferFrom(from, to, tokenId)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

// TokenURI get the URI of the token. a token's uri consists of CategoryURI and tokenId. Obligatory.
// @param tokenId, tokenId be queried. Obligatory.
// @return uri, the uri of the token.
// @return err, the error msg if some error occur.
func TokenURI(tokenId string) (uri string, err error) {
	category, err1 := GetCategoryByTokenId(tokenId)
	if err1 != nil {
		return "", err1
	}
	if category == nil {
		return "", errors.New("category not found")
	}
	//todo:format
	//http://chainmaker.org.cn/token/%s.json
	//http://chainmaker.org.cn/token/1000.json
	//http://www.chainmaker.org, 1000
	//http://www.chainmaker.org/1000
	//return fmt.Sprintf(category.CategoryURI, tokenId), nil
	return category.CategoryURI + "/" + tokenId, nil
}

func main() {

}
