/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package chainclient

import (
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
	"errors"
	"fmt"
)

const (
	OrgId1 = "wx-org1.chainmaker.org"
	OrgId2 = "wx-org2.chainmaker.org"
	OrgId3 = "wx-org3.chainmaker.org"
	OrgId4 = "wx-org4.chainmaker.org"
)

var Users = map[string]*User{
	"org1client1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/client1/client1.tls.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/client1/client1.tls.crt",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/client1/client1.sign.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/client1/client1.sign.crt",
	},
	"org2client1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/client1/client1.tls.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/client1/client1.tls.crt",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/client1/client1.sign.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/client1/client1.sign.crt",
	},
	"org1admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/admin1/admin1.tls.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/admin1/admin1.tls.crt",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/admin1/admin1.sign.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/admin1/admin1.sign.crt",
	},
	"org2admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/admin1/admin1.tls.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/admin1/admin1.tls.crt",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/admin1/admin1.sign.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/admin1/admin1.sign.crt",
	},
	"org3admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node3/user/admin1/admin1.tls.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node3/user/admin1/admin1.tls.crt",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node3/user/admin1/admin1.sign.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node3/user/admin1/admin1.sign.crt",
	},
	"org4admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node4/user/admin1/admin1.tls.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node4/user/admin1/admin1.tls.crt",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node4/user/admin1/admin1.sign.key",
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node4/user/admin1/admin1.sign.crt",
	},
}

var PermissionedPkUserss = map[string]*PermissionedPkUsers{
	"org1client1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/client1/client1.key",
		OrgId1,
	},
	"org2client1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/client1/client1.key",
		OrgId2,
	},
	"org1admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/admin/admin.key",
		OrgId1,
	},
	"org2admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/admin/admin.key",
		OrgId2,
	},
	"org3admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node3/admin/admin.key",
		OrgId3,
	},
	"org4admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node4/admin/admin.key",
		OrgId4,
	},
}

var PkUserss = map[string]*PkUsers{
	"org1client1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/user/client1/client1.key",
	},
	"org2client1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/user/client1/client1.key",
	},
	"org1admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node1/admin/admin1/admin1.key",
	},
	"org2admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node2/admin/admin2/admin2.key",
	},
	"org3admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node3/admin/admin3/admin3.key",
	},
	"org4admin1": {
		"/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/build/config/node4/admin/admin4/admin4.key",
	},
}

// PkUsers Pk用户定义
type PkUsers struct {
	SignKeyPath string
}

// PermissionedPkUsers PWk用户定义
type PermissionedPkUsers struct {
	SignKeyPath string
	OrgId       string
}

// User 用户节点证书路径定义
type User struct {
	TlsKeyPath, TlsCrtPath   string
	SignKeyPath, SignCrtPath string
}

// ClientCreator 创建连接的接口
type ClientCreator interface {
	CreateClientWithConfig(string) (*sdk.ChainClient, error)
}

type ClientCreate struct{}

// ContractClaimCreator 创建合约声明的接口
type ContractClaimCreator interface {
	UserContractCreate(*sdk.ChainClient, string, string, bool, ...string) error
}

// CreateClientWithConfig 用配置文件的方式创建链
func (c ClientCreate) CreateClientWithConfig(sdkConfPath string) (*sdk.ChainClient, error) {
	chainClient, err := sdk.NewChainClient(sdk.WithConfPath(sdkConfPath))
	if err != nil {
		return nil, err
	}
	return chainClient, nil
}

// InvokeUserContract 调用智能合约
func InvokeUserContract(client *sdk.ChainClient, contractName, method, txId string,
	params []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

	resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)

	if err != nil {
		return resp, err
	}
	if resp.Code != common.TxStatusCode_SUCCESS {
		return resp, fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}
	return resp, nil
}

// GetEndorsersWithAuthType 各组织Admin权限用户签名
func GetEndorsersWithAuthType(hashType crypto.HashType, authType sdk.AuthType, payload *common.Payload,
	usernames ...string) ([]*common.EndorsementEntry, error) {
	var endorsers []*common.EndorsementEntry

	for _, name := range usernames {
		var entry *common.EndorsementEntry
		var err error
		switch authType {
		case sdk.PermissionedWithCert:
			u, ok := Users[name]
			if !ok {
				return nil, errors.New("user not found")
			}
			entry, err = sdkutils.MakeEndorserWithPath(u.SignKeyPath, u.SignCrtPath, payload)
			if err != nil {
				return nil, err
			}
		case sdk.PermissionedWithKey:
			u, ok := PermissionedPkUserss[name]
			if !ok {
				return nil, errors.New("user not found")
			}
			entry, err = sdkutils.MakePkEndorserWithPath(u.SignKeyPath, hashType, u.OrgId, payload)
			if err != nil {
				return nil, err
			}
		case sdk.Public:
			u, ok := PkUserss[name]
			if !ok {
				return nil, errors.New("user not found")
			}
			entry, err = sdkutils.MakePkEndorserWithPath(u.SignKeyPath, hashType, "", payload)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("invalid authType")
		}
		endorsers = append(endorsers, entry)
	}

	return endorsers, nil
}

// CheckProposalRequestResp 检查交易是否成功
func CheckProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		if resp.Message == "" {
			resp.Message = resp.Code.String()
		}
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != 0 {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}
