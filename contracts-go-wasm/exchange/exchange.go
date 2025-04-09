package main

const (
	paramToken    = "tokenId"
	paramFrom     = "from"
	paramTo       = "to"
	paramMetaData = "metadata"
	paramAmount   = "amount"
	trueString    = "true"
)

//go:wasmexport init_contract
func InitContract() {
	ctx := NewSimContext()
	ctx.SuccessResult("Init contract success")
}

//go:wasmexport upgrade
func Upgrade() {
	ctx := NewSimContext()
	ctx.SuccessResult("Upgrade contract success")
}

//go:wasmexport buyNow
func buyNow() {
	ctx := NewSimContext()

	tokenId, _ := ctx.ArgString(paramToken)
	from, _ := ctx.ArgString(paramFrom)
	to, _ := ctx.ArgString(paramTo)
	metadata, _ := ctx.ArgString(paramMetaData)

	if tokenId == "" || from == "" || to == "" {
		ctx.ErrorResult("tokenId/from/to should not be empty")
	}

	// 添加白名单
	args := map[string][]byte{
		"address": []byte(from + "," + to),
	}
	resp, code := ctx.CallContract("identity", "addWriteList", args)
	if code != SUCCESS || string(resp) != "add write list success" {
		ctx.Log("[buyNow] addWriteList failed: " + string(resp))
		ctx.ErrorResult("[buyNow] addWriteList failed")
	}

	// 检查from是否在白名单
	args = map[string][]byte{"address": []byte(from)}
	resp, code = ctx.CallContract("identity", "isApprovedUser", args)
	if code != SUCCESS || string(resp) != trueString {
		ctx.ErrorResult("address[" + from + "] not registered")
	}

	// 检查to是否在白名单
	args = map[string][]byte{"address": []byte(to)}
	resp, code = ctx.CallContract("identity", "isApprovedUser", args)
	if code != SUCCESS || string(resp) != trueString {
		ctx.ErrorResult("address[" + to + "] not registered")
	}

	// 向from铸造 NFT
	args = map[string][]byte{
		"to":       []byte(from),
		"tokenId":  []byte(tokenId),
		"metadata": []byte(metadata),
	}
	resp, code = ctx.CallContract("erc721", "mint", args)
	if code != SUCCESS || string(resp) != "mint success" {
		ctx.Log("[buyNow] mint failed: " + string(resp))
		ctx.ErrorResult("[buyNow] mint failed")
	}

	// 设置 from 的资产授权
	args = map[string][]byte{
		"approvalFrom": []byte(from),
	}
	resp, code = ctx.CallContract("erc721", "setApprovalForAll2", args)
	if code != SUCCESS || string(resp) != "setApprovalForAll2 success" {
		ctx.Log("[buyNow] setApprovalForAll2 failed: " + string(resp))
		ctx.ErrorResult("[buyNow] setApprovalForAll2 failed")
	}

	// erc721 NFT 转移 from -> to
	args = map[string][]byte{
		"from":    []byte(from),
		"to":      []byte(to),
		"tokenId": []byte(tokenId),
	}
	resp, code = ctx.CallContract("erc721", "transferFrom", args)
	if code != SUCCESS || string(resp) != "transfer success" {
		ctx.ErrorResult("erc721 transferFrom error: " + string(resp))
	}
	spender, _ := ctx.GetSenderPk()
	ctx.EmitEvent("buyNow", from+to+spender)
	ctx.SuccessResult("buyNow success")
}

func main() {}
