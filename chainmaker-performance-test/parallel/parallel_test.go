/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package parallel

import (
	"chain-performance-test/datahandler"
	logger "chain-performance-test/log"
	"chain-performance-test/mock"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	// 改为你的sdk配置文件路径
	sdkPwkConfigPath  = "../config_example/sdk_config_pwk_example.yml"
	sdkPKConfigPath   = "../config_example/sdk_config_pk_example.yml"
	sdkCertConfigPath = "../config_example/sdk_config_ca_example.yml"
)

func TestBeginLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestBeginLog"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.BeginLog()
		})
	}
}

// 安装Claim智能合约测试
func TestInstallContractInstance(t *testing.T) {
	cc := &sdk.ChainClient{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClientCreator := mock.NewMockClientCreator(ctrl)
	mockContractClaimCreator := mock.NewMockContractClaimCreator(ctrl)

	mockClientCreator.EXPECT().CreateClientWithConfig(gomock.Any()).Return(cc, nil).AnyTimes()
	mockContractClaimCreator.EXPECT().UserContractCreate(cc, true, gomock.Any()).Return(nil).AnyTimes()

	tests := []struct {
		name  string
		model string
	}{
		{
			name:  "Test_PermissionWithKey",
			model: "PermissionWithKey",
		},
		{
			name:  "Test_Public",
			model: "Public",
		},
		{
			name:  "Test_PermissionWithCert",
			model: "PermissionWithCert",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Model = tt.model
			Clientslen = 1
			sdkConfigPaths = make([]string, Clientslen)
			sdkPKConfigPaths = make([]string, Clientslen)
			SdkPWKConfigPaths = make([]string, Clientslen)
			SdkPWKConfigPaths[0] = sdkPwkConfigPath
			sdkPKConfigPaths[0] = sdkPKConfigPath
			sdkConfigPaths[0] = sdkCertConfigPath
			err := InstallContractInstance(mockClientCreator, mockContractClaimCreator)
			require.Nil(t, err)
		})
	}
}

//func TestInitBeforeTest(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//		{
//			name: "TestInitBeforeTest",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			datahandler.ParametersList.ChainInformationList.HostList = []string{""}
//			ContractType = "Claim"
//			Model = "PermissionWithKey"
//			Clientslen = 1
//			SdkPWKConfigPaths = make([]string, Clientslen)
//			SdkPWKConfigPaths[0] = sdkPwkConfigPath
//			InitBeforeTest()
//		})
//	}
//}

//func TestContinuousClaimPressureTest(t *testing.T) {
//	type args struct {
//		client *sdk.ChainClient
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "TestContinuousClaimPressureTest",
//			args: args{client: nil},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			Wg.Add(10)
//			// 被测试的链
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			LoopNum = 10
//			// 10个并发,每个并发10次交易测试
//			for i := 0; i < 10; i++ {
//				err = ContinuousClaimPressureTest(cc)
//			}
//			Wg.Wait()
//			require.Nil(t, err)
//		})
//	}
//}

//func TestContinuousQueryPressureTest(t *testing.T) {
//	type args struct {
//		client *sdk.ChainClient
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "TestContinuousQueryPressureTest",
//			args: args{client: nil},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			Wg.Add(10)
//			// 被测试的链
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			LoopNum = 10
//			// 10个并发,每个并发10次交易测试
//			for i := 0; i < 10; i++ {
//				err = ContinuousQueryPressureTest(cc)
//			}
//			Wg.Wait()
//			require.Nil(t, err)
//		})
//	}
//}
//
//func TestCreateClientWithConfig(t *testing.T) {
//	type args struct {
//		sdkConfPath string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "TestCreateClientWithConfig",
//			args: args{
//				sdkConfPath: sdkConfigPath,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			_, err := CreateClientWithConfig(tt.args.sdkConfPath)
//			require.Nil(t, err)
//		})
//	}
//}

func TestGenerateClintConfig(t *testing.T) {
	type args struct {
		nodeNum int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_GenerateConstConfig",
			args: args{
				nodeNum: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datahandler.GenerateClintConfig(tt.args.nodeNum)
		})
	}
}

func TestGenerateConstConfig(t *testing.T) {
	type args struct {
		model            string
		contractFunction string
		contractName     string
		runTimeType      string
		allParams        string
		climbTime        int
		loopNum          int
		sleepTime        int
		threadNum        int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_GenerateConstConfig",
			args: args{
				model:            "PermissionWithCert",
				contractFunction: "save",
				contractName:     "claim001",
				runTimeType:      "WASMER",
				allParams:        "fileName:aaaaaaaaaaaaa||fileHash:bbbbbbbbbbbbb||fileTime:cccccccccccccc",
				climbTime:        5,
				loopNum:          100,
				sleepTime:        100,
				threadNum:        100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Chdir("..") // 进入上一级目录
			if err != nil {
				fmt.Println("Failed to change directory:", err)
				return
			}
			datahandler.GenerateConstConfig(tt.args.model, tt.args.contractFunction, tt.args.contractName, tt.args.runTimeType, tt.args.allParams, tt.args.climbTime, tt.args.loopNum, tt.args.sleepTime, tt.args.threadNum)
			TestInitConfig(t)
		})
	}
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name                   string
		wantContractParameters *ContractParameters
		wantPressureParameters *PressureParameters
	}{
		{
			name: "TestInitConfig",
			wantContractParameters: &ContractParameters{
				ContractName:   "claim001",
				ContractType:   "Claim",
				RuntimeType:    "WASMER",
				ContractMethod: "save",
				Params:         "fileName:aaaaaaaaaaaaa||fileHash:bbbbbbbbbbbbb||fileTime:cccccccccccccc",
			},
			wantPressureParameters: &PressureParameters{
				ThreadNum: 100,
				LoopNum:   100,
				SleepTime: 100,
				ClimbTime: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datahandler.ParametersList.ChainInformationList.HostList = []string{""}
			InitConfig()
			require.Equal(t, tt.wantContractParameters, ConfigurableParameter.ContractParameters)
			require.Equal(t, tt.wantPressureParameters, ConfigurableParameter.PressureParameters)
		})
	}
}

func TestProcessLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestProcessLog",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ProcessLog()
			files, err := os.ReadDir(".")
			if err != nil {
				log.Fatal(err)
			}
			count := 0
			for _, file := range files {
				if strings.HasPrefix(file.Name(), "sdk.log.") && !file.IsDir() {
					count++
				}
			}
			require.Equal(t, 0, count)
		})
	}
}

func TestClaimTPSCTPSWithSub(t *testing.T) {
	cc := &sdk.ChainClient{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ThreadNum = 0
	Clientslen = 1
	Client = cc
	mockPressureTester := mock.NewMockContinuousClaimPressureTester(ctrl)
	mockClientCreator := mock.NewMockClientCreator(ctrl)
	mockBlockSubscriber := mock.NewMockBlockSubscriber(ctrl)

	// 设置模拟对象的期望行为
	mockPressureTester.EXPECT().ClaimPressureTest(gomock.Any()).Return(nil).AnyTimes()
	mockClientCreator.EXPECT().CreateClientWithConfig(gomock.Any()).Return(cc, nil).AnyTimes()
	mockBlockSubscriber.EXPECT().SetClient(gomock.Any()).AnyTimes()
	mockBlockSubscriber.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockBlockSubscriber.EXPECT().Close().Return(nil).AnyTimes()

	// 使用模拟对象调用被测函数
	err := ClaimTPSCTPSWithSub(mockPressureTester, mockClientCreator, mockBlockSubscriber)

	// 断言函数返回的错误是否符合预期
	require.Nil(t, err)
}

func TestQueryTPSCTPSWithSub(t *testing.T) {
	cc := &sdk.ChainClient{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ThreadNum = 0
	Clientslen = 1
	Client = cc
	mockPressureTester := mock.NewMockContinuousQueryPressureTester(ctrl)
	mockClientCreator := mock.NewMockClientCreator(ctrl)
	mockBlockSubscriber := mock.NewMockBlockSubscriber(ctrl)

	// 设置模拟对象的期望行为
	mockPressureTester.EXPECT().QueryPressureTest(gomock.Any()).Return(nil).AnyTimes()
	mockClientCreator.EXPECT().CreateClientWithConfig(gomock.Any()).Return(cc, nil).AnyTimes()
	mockBlockSubscriber.EXPECT().SetClient(gomock.Any()).AnyTimes()
	mockBlockSubscriber.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockBlockSubscriber.EXPECT().Close().Return(nil).AnyTimes()

	// 使用模拟对象调用被测函数
	err := QueryTPSCTPSWithSub(mockPressureTester, mockClientCreator, mockBlockSubscriber)
	// 断言函数返回的错误是否符合预期
	require.Nil(t, err)
}

//func TestClaimTPSCTPSWithSub(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//		{
//			name: "TestClaimTPSCTPSWithSub",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ThreadNum = 2
//			Clientslen = 1
//			ClimbTime = 5
//			LoopNum = 2
//			// 被测试的链
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			TestUserContractCreate(t)
//			Clients1 = make([]*sdk.ChainClient, Clientslen)
//			Clients1[0] = cc
//			Client = cc
//			ParametersList.ContractParametersList.ContractFunctionParametersMap = make(map[string]string, 3)
//			ParametersList.ContractParametersList.ContractFunctionParametersMap["fileName"] = "value1"
//			ParametersList.ContractParametersList.ContractFunctionParametersMap["fileHash"] = "value2"
//			ParametersList.ContractParametersList.ContractFunctionParametersMap["fileTime"] = "value3"
//			err = ClaimTPSCTPSWithSub()
//			require.Nil(t, err)
//		})
//	}
//}

//func TestQueryTPSCTPSWithSub(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//		{
//			name: "TestQueryTPSCTPSWithSub",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ThreadNum = 2
//			Clientslen = 1
//			ClimbTime = 5
//			LoopNum = 2
//			// 被测试的链
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			TestUserContractCreate(t)
//			Clients1 = make([]*sdk.ChainClient, Clientslen)
//			Clients1[0] = cc
//			Client = cc
//			ParametersList.ContractParametersList.ContractFunctionParametersMap = make(map[string]string, 3)
//			ParametersList.ContractParametersList.ContractFunctionParametersMap["fileName"] = "value1"
//			ParametersList.ContractParametersList.ContractFunctionParametersMap["fileHash"] = "value2"
//			err = QueryTPSCTPSWithSub()
//			require.Nil(t, err)
//		})
//	}
//}

func TestUnzip(t *testing.T) {
	type args struct {
		zipPath string
		dstDir  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestUnzip",
			args: args{
				zipPath: "",
				dstDir:  "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := datahandler.Unzip(tt.args.zipPath, tt.args.dstDir); (err != nil) == tt.wantErr {
				t.Errorf("Unzip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateBuildFileName(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "TestUpdateBuildFileName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datahandler.UpdateBuildFileName("/")
		})
	}
}

func TestUpdateConstConfigYaml(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "TestUpdateConstConfigYaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datahandler.UpdateConstConfigYaml()
		})
	}
}

func TestUpdateYaml(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "TestUpdateYaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datahandler.UpdateYaml()
		})
	}
}

//func TestUserContractCreate(t *testing.T) {
//	type args struct {
//		client         *sdk.ChainClient
//		withSyncResult bool
//		usernames      []string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "TestUserContractCreate",
//			args: args{
//				client:         nil,
//				withSyncResult: true,
//				usernames: []string{UserNameOrg1Admin1, UserNameOrg2Admin1,
//					UserNameOrg3Admin1, UserNameOrg4Admin1},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			err = UserContractCreate(cc, tt.args.withSyncResult, tt.args.usernames...)
//			require.Nil(t, err)
//		})
//	}
//}

func TestgenerateByteCodePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test_generateByteCodePath",
			want: "./contract/claim_demo/rustFact.wasm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RuntimeType = common.RuntimeType_WASMER
			generateByteCodePath()
			require.Equal(t, tt.want, ContractByteCodePath)
		})
	}
}

//func TestCheckProposalRequestResp(t *testing.T) {
//	type args struct {
//		resp               *common.TxResponse
//		needContractResult bool
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "Test_checkProposalRequestResp",
//			args: args{
//				resp: &common.TxResponse{Code: common.TxStatusCode_SUCCESS, ContractResult: &common.ContractResult{
//					Code: 0,
//				}},
//				needContractResult: true,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			err := CheckProposalRequestResp(tt.args.resp, tt.args.needContractResult)
//			require.Nil(t, err)
//		})
//	}
//}

//func TestCreateUserContract(t *testing.T) {
//	type args struct {
//		client         *sdk.ChainClient
//		version        string
//		byteCodePath   string
//		kvs            []*common.KeyValuePair
//		withSyncResult bool
//		usernames      []string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "Test_createUserContract",
//			args: args{
//				client:         nil,
//				version:        "3.0.0",
//				byteCodePath:   "./contract/claim_demo/rustFact.wasm",
//				kvs:            nil,
//				withSyncResult: true,
//				usernames: []string{UserNameOrg1Admin1, UserNameOrg2Admin1,
//					UserNameOrg3Admin1, UserNameOrg4Admin1},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			ContractName = "test"
//			RuntimeType = common.RuntimeType_WASMER
//			ContractByteCodePath = "./contract/claim_demo/rustFact.wasm"
//			var kvs []*common.KeyValuePair
//			_, err = CreateUserContract(cc, tt.args.version, tt.args.byteCodePath, kvs, tt.args.withSyncResult, tt.args.usernames...)
//			require.Nil(t, err)
//		})
//	}
//}

//func TestInvokeUserContract(t *testing.T) {
//	type args struct {
//		client         *sdk.ChainClient
//		contractName   string
//		method         string
//		txId           string
//		params         []*common.KeyValuePair
//		withSyncResult bool
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "test_invokeUserContract",
//			args: args{
//				client:         nil,
//				contractName:   "T",
//				method:         "P",
//				txId:           "",
//				params:         nil,
//				withSyncResult: false,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			_, err = InvokeUserContract(cc, tt.args.contractName, tt.args.method, tt.args.txId, tt.args.params, tt.args.withSyncResult)
//			require.Nil(t, err)
//		})
//	}
//}

func TestMoveLog(t *testing.T) {
	type args struct {
		fileNames []string
		dirPath   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test_moveLog",
			args: args{
				// 你需要先在当前目录下创建一个sdk.log.2023052412文件
				fileNames: []string{"sdk.log.2023052412"},
				dirPath:   "./",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MoveLog(tt.args.fileNames, tt.args.dirPath)
		})
	}
}

func TestRandParams(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test_randParams",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datahandler.ParametersList.ContractParametersList.ContractFunctionParametersMap = make(map[string]string, 2)
			datahandler.ParametersList.ContractParametersList.ContractFunctionParametersMap["test"] = "value1"
			datahandler.ParametersList.ContractParametersList.ContractFunctionParametersMap["test2"] = "value2"
			parms := RandParams()
			fmt.Println(parms)
		})
	}
}

func TestRuntimeTypeModification(t *testing.T) {
	type args struct {
		runtimeType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test_INVALID",
			args: args{runtimeType: "INVALID"},
		},
		{
			name: "test_NATIVE",
			args: args{runtimeType: "NATIVE"},
		},
		{
			name: "test_WASMER",
			args: args{runtimeType: "WASMER"},
		},
		{
			name: "test_WXVM",
			args: args{runtimeType: "WXVM"},
		},
		{
			name: "test_GASM",
			args: args{runtimeType: "GASM"},
		},
		{
			name: "test_EVM",
			args: args{runtimeType: "EVM"},
		},
		{
			name: "test_DOCKER_GO",
			args: args{runtimeType: "DOCKER_GO"},
		},
		{
			name: "test_JAVA",
			args: args{runtimeType: "JAVA"},
		},
		{
			name: "test_GO",
			args: args{runtimeType: "GO"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RuntimeTypeString = tt.args.runtimeType
			RuntimeTypeModification()
		})
	}
}

func TestSdkConfigPathsMake(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test_sdkConfigPathsMake",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Clientslen = 4
			SdkConfigPathsMake()
		})
	}
}

func TestSdkPKConfigPathsMake(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test_sdkPKConfigPathsMake",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Clientslen = 4
			SdkPKConfigPathsMake()
		})
	}
}

func TestSdkPWKConfigPathsMake(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test_sdkPWKConfigPathsMake",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Clientslen = 4
			SdkPWKConfigPathsMake()
		})
	}
}

//func TestUserContractClaimInvoke(t *testing.T) {
//	type args struct {
//		client         *sdk.ChainClient
//		method         string
//		withSyncResult bool
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "Test_UserContractClaimInvoke",
//			args: args{client: nil, method: "save", withSyncResult: false},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			err = UserContractClaimInvoke(cc, tt.args.method, tt.args.withSyncResult)
//			require.Nil(t, err)
//		})
//	}
//}

// 调用查询合约测试
//func TestUserContractQueryInvoke(t *testing.T) {
//	type args struct {
//		client         *sdk.ChainClient
//		method         string
//		withSyncResult bool
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "Test_UserContractQueryInvoke",
//			args: args{client: nil, method: "query", withSyncResult: false},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cc, err := CreateClientWithConfig(sdkConfigPath)
//			require.Nil(t, err)
//			defer func(cc *sdk.ChainClient) {
//				err = cc.Stop()
//				if err != nil {
//					fmt.Println(err)
//				}
//			}(cc)
//			err = UserContractQueryInvoke(cc, tt.args.method, tt.args.withSyncResult)
//			require.Nil(t, err)
//		})
//	}
//}
