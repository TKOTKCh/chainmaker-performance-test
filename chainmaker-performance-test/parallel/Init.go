/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package parallel

import (
	"chain-performance-test/chainclient"
	"chain-performance-test/datahandler"
	logger "chain-performance-test/log"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"encoding/hex"
	"fmt"
	"github.com/go-yaml/yaml"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	UserNameOrg1Admin1    = "org1admin1" // admin节点名
	UserNameOrg2Admin1    = "org2admin1"
	UserNameOrg3Admin1    = "org3admin1"
	UserNameOrg4Admin1    = "org4admin1"
	createContractTimeout = int64(50)
	claimVersion          = "2.0.0"
)

var (
	Client                *sdk.ChainClient       // 被压测的链
	Clients1              []*sdk.ChainClient     // 被压测的链数组
	sdkConfigPaths        []string               // 被压测的cert模式链配置文件数组
	sdkPKConfigPaths      []string               // 被压测的public模式链配置文件数组
	SdkPWKConfigPaths     []string               // 被压测的permissionedWithKey模式链配置文件数组
	Clientslen            int                    // 被压测的链个数
	ContractByteCodePath  string                 // 存证类合约智能合约路径
	txCounts              uint32             = 0 // 成功上链交易数
	timeStart             time.Time              // 压测开始时间
	timeEnd               time.Time              // 压测开始时间
	BlockStart            int64                  // 压测开始区块
	Wg                    = sync.WaitGroup{}
	ConfigurableParameter *ConfigurableParameters // 可配置参数
	ClientsParameter      *ClientsParameters      // 多加点加压个数参数
	Model                 string                  // 身份认证模型
	ContractName          string                  // 智能合约名称
	ContractType          string                  // 智能合约类型
	RuntimeTypeString     string                  // 智能合约语言类型, string类型
	RuntimeType           common.RuntimeType      // 智能合约语言类型, common.RuntimeType类型
	contractMethod        string                  // 智能合约被压测方法
	Params                string                  // 智能合约被压测方法参数
	ThreadNum             int                     // 单次并发进程数
	LoopNum               int                     // 压测并发次数
	SleepTime             int                     // 并发间隔,单位ms
	ClimbTime             int                     // 爬坡时间,单位s
	AddOption             string
	RandomSeed            int64 // 生成tokenid的随机种子
	lastTokenId           = new(big.Int)
	mu                    sync.Mutex
)

// ConfigurableParameters 可配置参数定义
type ConfigurableParameters struct {
	ChainParameters    *ChainParameters    `yaml:"chain_information"`                // 长安链可配置参数
	ContractParameters *ContractParameters `yaml:"contract_configurable_parameters"` // 智能合约可配置参数
	PressureParameters *PressureParameters `yaml:"pressure_configurable_parameters"` // 并发压测可配置参数
}

// ChainParameters 智能合约可配置参数定义
type ChainParameters struct {
	Model string `json:"model"` //身份权限管理
}

// ContractParameters 智能合约可配置参数定义
type ContractParameters struct {
	ContractName   string `yaml:"contract_name"`   // 智能合约名称参数
	ContractType   string `yaml:"contract_type"`   // 智能合约类型参数
	RuntimeType    string `yaml:"runtime_type"`    // 智能合约语言类型参数
	ContractMethod string `yaml:"contract_method"` // 智能合约被压测方法参数
	Params         string `yaml:"params"`          // 智能合约被压测方法参数列表参数
}

// PressureParameters 并发压测可配置参数定义
type PressureParameters struct {
	ThreadNum int    `yaml:"thread_num"` // 单次并发进程数参数
	LoopNum   int    `yaml:"loop_num"`   // 压测并发次数参数
	SleepTime int    `yaml:"sleep_time"` // 并发间隔参数,单位ms
	AddOption string `yaml:"add_option"` // 用于找loopNum
	ClimbTime int    `yaml:"climb_time"` // 并发间隔参数,单位ms
}

type ClientsParameters struct {
	NodeNum int `yaml:"node_num"` // 长安链多节点加压，加压节点个数
}

type ContractCreator struct{}

// InitBeforeTest 并发压测初始化
func InitBeforeTest(loopi int) {

	//1.解析可配置参数
	InitConfig(loopi)
	//2.安装智能合约
	var clientCreator chainclient.ClientCreate
	var contractClaimCreator ContractCreator
	//对于exchange要跨合约调用的函数，需要先安装identity，erc721合约
	if ContractName == "exchange" {
		contractName := "identity"
		err := InstallContractInstance(clientCreator, contractName, generateByteCodePath(contractName), contractClaimCreator)
		if err != nil {
			fmt.Println("安装先导合约", contractName, "失败")
		}
		contractName = "erc721"
		err = InstallContractInstance(clientCreator, contractName, generateByteCodePath(contractName), contractClaimCreator)
		if err != nil {
			fmt.Println("安装先导合约", contractName, "失败")
		}
	}
	err := InstallContractInstance(clientCreator, ContractName, ContractByteCodePath, contractClaimCreator)
	if err != nil {
		//logger.Logger.Panic("安装智能合约失败:", err)
		logger.Logger.Println("安装智能合约失败：", err)
	}

}

//生成随机长度地址
func randomHexString(length int) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// 生成指定长度的随机数字字符串（首位不为0），用于erc721的tokenid生成
func RandomNumberString(length int) string {
	//加锁这一步对性能影响不大
	mu.Lock()         // 加锁确保并发安全
	defer mu.Unlock() // 解锁
	lastTokenId.Add(lastTokenId, big.NewInt(1))

	return lastTokenId.String()
}

//处理参数
func handleParams(paramString string) map[string]string {
	FunctionParametersMap := make(map[string]string)
	initParametersList := strings.Split(paramString, "||")
	for i := 0; i < len(initParametersList); i++ {
		index := strings.Index(initParametersList[i], ":")
		if index != -1 {
			firstPart := initParametersList[i][:index]
			secondPart := initParametersList[i][index+1:]
			FunctionParametersMap[firstPart] = secondPart
		}
	}
	return FunctionParametersMap
}

//生成调用函数参数
func getMethodParams(ContractName, contractMethod string) string {
	if ContractName == "identity" {
		if contractMethod == "callerAddress" || contractMethod == "address" {
			Params = ""
		} else {
			Params = "address:"
			for i := 0; i < 100; i++ {
				s, _ := randomHexString(40)
				Params += s
				if i != 99 {
					Params += ","
				}
			}
		}
	}
	if ContractName == "erc721" {
		if contractMethod == "tokenURI" || contractMethod == "ownerOf" || contractMethod == "tokenMetadata" || contractMethod == "tokenLatestTxInfo" || contractMethod == "getApprove" {
			Params = "tokenId:111111111111111111111112"
		}
		if contractMethod == "balanceOf" || contractMethod == "accountTokens" {
			Params = "account:c0d8e4ce07a48081eff14a3016699b1c839c4375"
		}
		if contractMethod == "mint" {
			Params = "to:8acfaca5eeec9f6f7c23c4ffac969b86f27799b0||tokenId:111111111111111111111111||metadata:http://chainmaker.org.cn/"
		}
		//if contractMethod == "approve" {
		//	Params = "to:818fac1ac51525aeedf619a9a339b95854930159||tokenId:111111111111111111111111"
		//}
		if contractMethod == "setApprovalForAll2" {
			Params = "approvalFrom:8acfaca5eeec9f6f7c23c4ffac969b86f27799b0"
		}
		if contractMethod == "transferFrom" {
			Params = "from:8acfaca5eeec9f6f7c23c4ffac969b86f27799b0||to:818fac1ac51525aeedf619a9a339b95854930159||tokenId:11111111111111111111111||metadata:http://chainmaker.org.cn/"
		}
	}
	if ContractName == "compute" {
		Params = ""
	}
	if ContractName == "exchange" {
		if contractMethod == "buyNow" {
			Params = "from:8acfaca5eeec9f6f7c23c4ffac969b86f27799b0||to:818fac1ac51525aeedf619a9a339b95854930159||tokenId:11111111111111111111111||metadata:http://chainmaker.org.cn/"
		}
	}
	//fmt.Println(Params)
	return Params
}

// InitConfig 初始化可配置参数
func InitConfig(loopi int) {

	yamlFile1, err := os.ReadFile("/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/config/clients.yml")
	if err != nil {
		logger.Logger.Panic("打开文件失败:", err)
	}
	err = yaml.Unmarshal(yamlFile1, &ClientsParameter)
	if err != nil {
		logger.Logger.Panic(err)
	}
	Clientslen = ClientsParameter.NodeNum // 节点个数

	yamlFile, err := os.ReadFile("/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/config/const_config.yml")
	if err != nil {
		logger.Logger.Panic("打开文件失败:", err)
	}
	err = yaml.Unmarshal(yamlFile, &ConfigurableParameter)
	if err != nil {
		logger.Logger.Panic(err)
	}
	//num := fmt.Sprintf("%08v", rand.Intn(99999999))
	Model = ConfigurableParameter.ChainParameters.Model // 身份认证模式
	//ContractName = ConfigurableParameter.ContractParameters.ContractName + num  // 智能合约名称
	ContractName = ConfigurableParameter.ContractParameters.ContractName     // 智能合约名称
	ContractType = ConfigurableParameter.ContractParameters.ContractType     // 智能合约类型
	RuntimeTypeString = ConfigurableParameter.ContractParameters.RuntimeType // 智能合约语言类型
	contractMethod = ConfigurableParameter.ContractParameters.ContractMethod // 智能合约被压测方法

	//if judgeAdmin(ContractName, contractMethod) {
	//	Clientslen = 1
	//}

	ThreadNum = ConfigurableParameter.PressureParameters.ThreadNum / Clientslen // 单次并发进程数,总并发进程/节点数
	AddOption = ConfigurableParameter.PressureParameters.AddOption
	// 压测并发次数
	if AddOption == "yes" {
		LoopNum = ConfigurableParameter.PressureParameters.LoopNum + 10*loopi
	} else {
		LoopNum = ConfigurableParameter.PressureParameters.LoopNum
	}

	SleepTime = ConfigurableParameter.PressureParameters.SleepTime // 并发间隔,单位ms
	ClimbTime = ConfigurableParameter.PressureParameters.ClimbTime // 爬坡时间,单位s
	Params = getMethodParams(ContractName, contractMethod)         // 智能合约被压测方法参数
	RandomSeed = 0
	lastTokenId.SetString("111111111111111111111111", 10)
	fmt.Println(ContractType, ContractName, contractMethod, RuntimeTypeString, Params)
	// runtimeType 参数类型修改
	RuntimeTypeModification()
	// 链配置信息生成
	if Model == "PermissionWithCert" {
		SdkConfigPathsMake()
	} else if Model == "Public" {
		SdkPKConfigPathsMake()
	} else if Model == "PermissionWithKey" {
		SdkPWKConfigPathsMake()
	}

	ContractByteCodePath = generateByteCodePath(ContractName)

	// 参数列表解析
	datahandler.ParametersList = datahandler.Parameters{} //初始化
	datahandler.ParametersList.ContractParametersList.ContractFunctionParametersMap = handleParams(Params)

	logger.Logger.Println("====================== 可配置参数 ======================")
	logger.Logger.Printf("config.ContractParameters: %#v\n", ConfigurableParameter.ContractParameters)
	logger.Logger.Printf("config.PressureParameters: %#v\n", ConfigurableParameter.PressureParameters)
}

// InstallContractInstance 安装智能合约
func InstallContractInstance(clientCreator chainclient.ClientCreator, contractName string, contractByteCodePath string, contractClaimCreator chainclient.ContractClaimCreator) error {
	var err error
	Clients1 = make([]*sdk.ChainClient, Clientslen)
	if Model == "PermissionWithCert" {
		for i := 0; i < Clientslen; i++ {
			Clients1[i], err = clientCreator.CreateClientWithConfig(sdkConfigPaths[i])
		}
	} else if Model == "Public" {
		for i := 0; i < Clientslen; i++ {
			Clients1[i], err = clientCreator.CreateClientWithConfig(sdkPKConfigPaths[i])
		}
	} else if Model == "PermissionWithKey" {
		for i := 0; i < Clientslen; i++ {
			Clients1[i], err = clientCreator.CreateClientWithConfig(SdkPWKConfigPaths[i])
		}
	}
	if err != nil {
		return err
	}

	logger.Logger.Println("====================== 安装合约 ======================")
	// Admin name
	var usernames []string
	if Model == "PermissionWithCert" {
		usernames = []string{UserNameOrg1Admin1, UserNameOrg2Admin1, UserNameOrg3Admin1, UserNameOrg4Admin1}
	} else if Model == "Public" {
		usernames = []string{UserNameOrg1Admin1}
	} else if Model == "PermissionWithKey" {
		usernames = []string{UserNameOrg1Admin1, UserNameOrg2Admin1, UserNameOrg3Admin1, UserNameOrg4Admin1}
	}
	err = contractClaimCreator.UserContractCreate(Clients1[0], contractName, contractByteCodePath, true, usernames...)
	if err != nil {
		return err
	}
	return nil
}

// RuntimeTypeModification runtimeType参数类型修改
func RuntimeTypeModification() {
	switch RuntimeTypeString {
	case "INVALID":
		RuntimeType = common.RuntimeType_INVALID
	case "NATIVE":
		RuntimeType = common.RuntimeType_NATIVE
	case "WASMER":
		RuntimeType = common.RuntimeType_WASMER
	case "WXVM":
		RuntimeType = common.RuntimeType_WXVM
	case "GASM":
		RuntimeType = common.RuntimeType_GASM
	case "EVM":
		RuntimeType = common.RuntimeType_EVM
	case "DOCKER_GO":
		RuntimeType = common.RuntimeType_DOCKER_GO
	case "JAVA":
		RuntimeType = common.RuntimeType_JAVA
	case "GO":
		RuntimeType = common.RuntimeType_GO
	}
}

// SdkConfigPathsMake 生成cert模式链配置文件数组
func SdkConfigPathsMake() {
	sdkConfigPaths = make([]string, Clientslen)
	for i := 0; i < Clientslen; i++ {
		sdkConfigPaths[i] = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/config/sdk_config" + strconv.Itoa(i+1) + ".yml"
	}
}

// SdkPKConfigPathsMake 生成pubilc模式链配置文件数组
func SdkPKConfigPathsMake() {
	sdkPKConfigPaths = make([]string, Clientslen)
	for i := 0; i < Clientslen; i++ {
		sdkPKConfigPaths[i] = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/config/sdk_config_pk" + strconv.Itoa(i+1) + ".yml"
	}
}

// SdkPWKConfigPathsMake 生成permissionedWithKey模式链配置文件数组
func SdkPWKConfigPathsMake() {
	SdkPWKConfigPaths = make([]string, Clientslen)
	for i := 0; i < Clientslen; i++ {
		SdkPWKConfigPaths[i] = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/config/sdk_config_pwk" + strconv.Itoa(i+1) + ".yml"
	}
}

// generateByteCodePath 示例智能合约路径修改
func generateByteCodePath(contractName string) string {
	codePathPrefix := "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/"
	contractByteCodePath := ""
	if RuntimeType == common.RuntimeType_WASMER {
		contractByteCodePath = codePathPrefix + contractName + ".wasm"
		//contractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/tiny-go-chainmaker-contract-go.wasm"
		//ContractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/rustFact.wasm"
		//ContractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/tiny-go-chainmaker-contract-go.wasm"
	}
	if RuntimeType == common.RuntimeType_DOCKER_GO {
		contractByteCodePath = codePathPrefix + contractName + ".7z"
		//ContractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/dockerFact230.7z"
		//ContractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/factgo.7z"
		//ContractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/erc20.7z"

		//ContractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/factgomodify.7z"
	}
	if RuntimeType == common.RuntimeType_GASM {
		//tiny-go-chainmaker-contract-go.wasm是我用tinygo转成wasm
		contractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/tiny-go-chainmaker-contract-go.wasm"
		//ContractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/chainmaker-contract-go-modify.wasm"
		//ContractByteCodePath = "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/contract/claim_demo/factgo.wasm"

	}
	return contractByteCodePath
}

// UserContractCreate 创建用户智能合约
func (con ContractCreator) UserContractCreate(client *sdk.ChainClient, contractName string, contractByteCodePath string, withSyncResult bool, usernames ...string) error {

	//var kvs []*common.KeyValuePair
	kvs := RandCreateParams(contractName)
	fmt.Println(contractByteCodePath)
	// 安装用户智能合约
	resp, err := CreateUserContract(client, contractName, claimVersion, contractByteCodePath, kvs, withSyncResult, usernames...)
	// 记录压测前初始块高，用于消息订阅
	BlockStart = int64(resp.TxBlockHeight) + 1
	logger.Logger.Printf("BlockStart:%d\n", BlockStart)
	logger.Logger.Println("resp.ContractResult.Code:", resp.ContractResult.Code)
	logger.Logger.Println("resp.ContractResult.Message:", resp.ContractResult.Message)
	//chenhang修改
	//logger.Logger.Println("resp.Code:", resp.Code)
	//logger.Logger.Println("resp.Message:", resp.Message)
	if err != nil {
		return err
	}
	return nil
}

// CreateUserContract 安装用户智能合约
func CreateUserContract(client *sdk.ChainClient, contractName string, version,
	byteCodePath string, kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {
	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, RuntimeType, kvs)
	if err != nil {
		logger.Logger.Panic(err)
	}
	// 各组织Admin权限用户签名
	endorsers, err := chainclient.GetEndorsersWithAuthType(client.GetHashType(), client.GetAuthType(), payload, usernames...)
	if err != nil {
		logger.Logger.Panic(err)
	}

	// 发送请求
	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return resp, err
	}

	// 检查交易是否成功
	err = chainclient.CheckProposalRequestResp(resp, true)
	return resp, err
}

func RandCreateParams(contractName string) []*common.KeyValuePair {

	params := []*common.KeyValuePair{}
	//定义一个不需要参数的合约类型列表
	noParamContracts := []string{"identity", "save", "compute", "exchange"}
	// 检查 contractName 是否在 noParamContracts 中
	for _, t := range noParamContracts {
		if contractName == t {
			return params
		}
	}
	//下面是需要初始化参数的合约
	var initParams string = ""
	if contractName == "erc721" {
		initParams = "name:huanletoken||symbol:hlt||tokenURI:https://chainmaker.org.cn"
	}
	FunctionParametersMap := handleParams(initParams)
	fmt.Println(FunctionParametersMap)
	params = RandParams(FunctionParametersMap)
	return params
}

// RandParams 随机化参数
func RandParams(FunctionParametersMap map[string]string) []*common.KeyValuePair {
	normalParams := make(map[string]string)
	params := []*common.KeyValuePair{}
	// 获取带压测方法的参数列表
	set := FunctionParametersMap
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	for k, v := range set {
		param := new(common.KeyValuePair)
		param.Key = k
		value := v
		if strings.Contains(strings.ToLower(k), strings.ToLower("time")) {
			value = curTime
		}
		if strings.Contains(strings.ToLower(k), strings.ToLower("tokenId")) {
			value = RandomNumberString(len(value))
		}
		normalParams[k] = value
		param.Value = []byte(value)
		params = append(params, param)
	}

	return params
}
