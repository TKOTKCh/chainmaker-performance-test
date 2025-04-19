/*
  Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

  SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	projectInfoArgKey       = "projectInfo"
	projectIdArgKey         = "projectId"
	projectItemIdArgKey     = "itemId"
	projectInfoStoreKey     = "project"
	projectItemsStoreMapKey = "projectItems"
	projectVotesStoreMapKey = "projectVotes"
	trueString              = "1"
)

type ProjectInfo struct {
	Id        string      `json:"Id"`
	PicUrl    string      `json:"PicUrl"`
	Title     string      `json:"Title"`
	StartTime string      `json:"StartTime"`
	EndTime   string      `json:"EndTime"`
	Desc      string      `json:"Desc"`
	Items     []*ItemInfo `json:"Items"`
}

type ItemInfo struct {
	Id     string `json:"Id"`
	PicUrl string `json:"PicUrl"`
	Desc   string `json:"Desc"`
	Url    string `json:"Url"`
}

type ProjectVotesInfo struct {
	ProjectId string           `json:"ProjectId"`
	ItemVotes []*ItemVotesInfo `json:"ItemVotes"`
}

type ItemVotesInfo struct {
	ItemId     string   `json:"ItemId"`
	VotesCount int      `json:"VotesCount"`
	Voters     []string `json:"Voters"`
}

///////////////////////// 合约入口函数 /////////////////////////

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

///////////////////////// 核心业务函数 /////////////////////////

//go:wasmexport issueProject
func issueProject() {
	ctx := NewSimContext()

	// 获取并解析参数
	projectInfoBytes, _ := ctx.Arg(projectInfoArgKey)

	if len(projectInfoBytes) == 0 {
		ctx.ErrorResult("issueProject should have arg of " + projectInfoArgKey)
		return
	}

	var pi ProjectInfo
	if err1 := json.Unmarshal(projectInfoBytes, &pi); err1 != nil {
		ctx.ErrorResult("unmarshal project info failed")
		return
	}

	// 检查项目是否存在
	piBytes, err := ctx.GetStateByte(projectInfoStoreKey, pi.Id)
	if err != SUCCESS {
		ctx.ErrorResult("get project info failed")
		return
	}
	if len(piBytes) > 0 {
		ctx.ErrorResult("project already exists")
		return
	}

	// 存储项目信息
	if err := ctx.PutStateByte(projectInfoStoreKey, pi.Id, projectInfoBytes); err != SUCCESS {
		ctx.ErrorResult("store project info failed")
		return
	}

	// 存储项目选项
	for _, item := range pi.Items {
		key := fmt.Sprintf("%s_%s", pi.Id, item.Id)
		if err := ctx.PutState(projectItemsStoreMapKey, key, trueString); err != SUCCESS {
			ctx.ErrorResult("set project item to storemap failed")
			return
		}
	}

	ctx.EmitEvent("issue_project", string(projectInfoBytes))
	ctx.SuccessResult("issue project success")
	return
}

//go:wasmexport vote
func vote() {
	ctx := NewSimContext()

	// 获取参数
	projectId, _ := ctx.ArgString(projectIdArgKey)
	projectItemId, _ := ctx.ArgString(projectItemIdArgKey)
	if len(projectId) == 0 {
		ctx.ErrorResult("invalid project id")
		return
	}
	if len(projectItemId) == 0 {
		ctx.ErrorResult("invalid project item id")
		return
	}

	// 获取项目信息
	projectInfoBytes, err := ctx.GetStateByte(projectInfoStoreKey, projectId)
	if err != SUCCESS || len(projectInfoBytes) == 0 {
		ctx.ErrorResult("get project info from store failed")
		return
	}

	var pi ProjectInfo
	if err1 := json.Unmarshal(projectInfoBytes, &pi); err1 != nil {
		ctx.ErrorResult("unmarshal project info from store failed")
		return
	}

	// 检查是否重复投票
	sender, _ := ctx.GetSenderPk()
	if isVoted, _ := ctx.GetState(projectId, sender); isVoted == trueString {
		ctx.ErrorResult(fmt.Sprintf("get is voted from store failed"))
		return
	}

	// 记录投票
	if err = ctx.PutState(projectId, sender, trueString); err != SUCCESS {
		ctx.ErrorResult("store voted info failed")
		return
	}

	// 验证投票选项
	itemKey := fmt.Sprintf("%s_%s", projectId, projectItemId)
	if val, _ := ctx.GetState(projectItemsStoreMapKey, itemKey); val != trueString {
		ctx.ErrorResult("can't found project or item")
		return
	}

	// 记录投票详情
	voteKey := fmt.Sprintf("%s_%s_%s", projectId, projectItemId, sender)
	if err = ctx.PutState(projectVotesStoreMapKey, voteKey, trueString); err != SUCCESS {
		ctx.ErrorResult("save vote detail failed")
		return
	}

	ctx.EmitEvent("vote", projectId, projectItemId, sender)
	ctx.SuccessResult("vote success")
	return
}

//go:wasmexport queryProjectVoters
func queryProjectVoters() {
	ctx := NewSimContext()

	pi, _ := ctx.ArgString(projectIdArgKey)
	if len(pi) == 0 {
		ctx.ErrorResult("invalid project id")
		return
	}

	rs, err := ctx.NewIteratorPrefixWithKeyField(projectItemsStoreMapKey, pi)
	if err != SUCCESS {
		ctx.ErrorResult(fmt.Sprintf("new store map iterator of project info failed, err: %s", err))
		return
	}
	projectVotesInfo := newProjectVotesInfo(pi)
	for {
		if !rs.HasNext() {
			break
		}
		var item string
		item, _, _, err = rs.Next()
		if err != SUCCESS {
			ctx.ErrorResult(fmt.Sprintf("iterator next failed, err: %s", err))
			return
		}
		itemId := strings.TrimPrefix(strings.TrimPrefix(item, projectItemsStoreMapKey), pi)
		if len(itemId) == 0 {
			ctx.ErrorResult("invalid itemId")
			return
		}
		var itemVotesInfo *ItemVotesInfo
		itemVotesInfo, err1 := queryItemVoters(pi, itemId)
		if err1 != nil {
			ctx.ErrorResult(err1.Error())
			return
		}
		projectVotesInfo.ItemVotes = append(projectVotesInfo.ItemVotes, itemVotesInfo)
	}

	var projectVotesInfoBytes []byte
	projectVotesInfoBytes, err1 := json.Marshal(projectVotesInfo)
	if err1 != nil {
		ctx.ErrorResult(err1.Error())
		return
	}
	ctx.SuccessResult(string(projectVotesInfoBytes))
	return
}

//go:wasmexport queryProjectItemVoters
func queryProjectItemVoters() {
	ctx := NewSimContext()

	pi, _ := ctx.ArgString(projectIdArgKey)
	if len(pi) == 0 {
		ctx.ErrorResult("invalid project id")
		return
	}
	projectItemId, _ := ctx.ArgString(projectItemIdArgKey)
	if len(projectItemId) == 0 {
		ctx.ErrorResult("invalid project item id")
		return
	}
	itemVotesInfo, err := queryItemVoters(pi, projectItemId)
	if err != nil {
		ctx.ErrorResult(err.Error())
		return
	}
	itemVotesInfoBytes, err := json.Marshal(itemVotesInfo)
	if err != nil {
		ctx.ErrorResult("marshal itemVotesInfo failed")
		return
	}
	ctx.SuccessResult(string(itemVotesInfoBytes))
	return
}

func queryItemVoters(projectId, itemId string) (*ItemVotesInfo, error) {

	ctx := NewSimContext()
	rs, err := ctx.NewIteratorPrefixWithKeyField(projectVotesStoreMapKey, projectId+"_"+itemId)
	if err != SUCCESS {
		return nil, fmt.Errorf("new store map iterator of project info failed, err: %s", err)
	}
	itemVotesInfo := newItemVotesInfo(itemId)

	for {
		if !rs.HasNext() {
			break
		}
		item, _, _, err := rs.Next()
		if err != SUCCESS {
			return nil, fmt.Errorf("iterator next failed, err: %s", err)
		}
		voter := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(item, projectVotesStoreMapKey), projectId), itemId)
		if len(voter) == 0 {
			return nil, fmt.Errorf("found invalid voter")
		}
		itemVotesInfo.Voters = append(itemVotesInfo.Voters, voter)
		itemVotesInfo.VotesCount++
	}
	return itemVotesInfo, nil
}
func newItemVotesInfo(itemId string) *ItemVotesInfo {
	return &ItemVotesInfo{
		ItemId: itemId,
		Voters: make([]string, 0, 5),
	}
}

func newProjectVotesInfo(projectId string) *ProjectVotesInfo {
	return &ProjectVotesInfo{
		ProjectId: projectId,
		ItemVotes: make([]*ItemVotesInfo, 0, 5),
	}
}

func main() {}
