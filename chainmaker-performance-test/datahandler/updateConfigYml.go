/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package datahandler

import (
	logger "chain-performance-test/log"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var (
	//文件夹名
	configFileName    string = "config"         //配置文件所在文件夹名
	configExampleName string = "config_example" //配置文件样例所在文件夹名

	//样例文件名
	sdkConfigExample    string = "sdk_config_ca_example"
	sdkConfigPkExample  string = "sdk_config_pk_example"  // pk模式下
	sdkConfigPWKExample string = "sdk_config_pwk_example" // pwk模式下
	constConfigExample  string = "const_config_example"
	clientsExample      string = "clients_example"

	//写入配置文件名
	sdkConfigPkFile   string = "sdk_config_pk" //pk模式下
	sdkConfigFile     string = "sdk_config"
	sdkConfigPWKFile  string = "sdk_config_pwk" //pwk模式下
	constConfigFile   string = "const_config"
	clientsConfigFile string = "clients"
)

// ChainNodesYaml cert模式下配置yml文件中ChainNodes对应结构体定义
type ChainNodesYaml struct {
	ConnCnt        int      `yaml:"conn_cnt"`
	EnableTls      bool     `yaml:"enable_tls"`
	NodeAddr       string   `yaml:"node_addr"`
	TlsHostName    string   `yaml:"tls_host_name"`
	TrustRootPaths []string `yaml:"trust_root_paths"`
}

// ChainNodesYamlForPWK PWK模式下配置yml文件中 ChainNodes对应结构体定义
type ChainNodesYamlForPWK struct {
	ConnCnt  int    `yaml:"conn_cnt"`
	NodeAddr string `yaml:"node_addr"`
}

// ChainNodesYamlForPK PK模式下配置yml文件中 ChainNodes对应结构体定义
type ChainNodesYamlForPK struct {
	ConnCnt  int    `yaml:"conn_cnt"`
	NodeAddr string `yaml:"node_addr"`
}

// GetConfigFilePathByCaller 获取要写入的配置文件路径
func GetConfigFilePathByCaller(configFileName string) string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	abPath2 := filepath.Join(abPath, "../"+configFileName)
	return strings.Replace(abPath2, OldSep, NewSep, -1)
}

// UpdateConstConfigYaml 更新ConstConfig配置文件yml文件
func UpdateConstConfigYaml() {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	// 连接yaml文件
	ConstConfig := viper.New()
	ConstConfig.AddConfigPath(examplePath)
	ConstConfig.SetConfigName(constConfigExample)
	ConstConfig.SetConfigType("yaml")
	if err := ConstConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}

	// 修改viper实例内容
	ConstConfig.Set("chain_information.model", ParametersList.ChainInformationList.Model)
	ConstConfig.Set("contract_configurable_parameters.contract_method", ParametersList.ContractParametersList.ContractFunction)
	ConstConfig.Set("contract_configurable_parameters.contract_name", ParametersList.ContractParametersList.ContractName)

	// 根据智能合约函数名判断合约类型
	saveContain := strings.Contains(strings.ToLower(ParametersList.ContractParametersList.ContractFunction), strings.ToLower("save"))
	findContain := strings.Contains(strings.ToLower(ParametersList.ContractParametersList.ContractFunction), strings.ToLower("find"))
	if saveContain {
		ConstConfig.Set("contract_configurable_parameters.contract_type", "Claim")
	} else if findContain {
		ConstConfig.Set("contract_configurable_parameters.contract_type", "Query")
	}

	ConstConfig.Set("contract_configurable_parameters.runtime_type", ParametersList.ContractParametersList.RunTimeType)
	ConstConfig.Set("pressure_configurable_parameters.climb_time", ParametersList.PressureParametersList.ClimbTime)
	ConstConfig.Set("pressure_configurable_parameters.loop_num", ParametersList.PressureParametersList.LoopNum)
	ConstConfig.Set("pressure_configurable_parameters.sleep_time", ParametersList.PressureParametersList.SleepTime)
	ConstConfig.Set("pressure_configurable_parameters.thread_num", ParametersList.PressureParametersList.ThreadNum)
	ConstConfig.Set("contract_configurable_parameters.params", ParametersList.ContractParametersList.ContractFunctionParameters)

	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	ConstSavePath := newPath + NewSep + constConfigFile + ".yml"
	errWrite := ConstConfig.WriteConfigAs(ConstSavePath)
	if errWrite != nil {
		return
	}
}

// pwk模式下修改编号i的SdkConfig配置文件
func UpdateYamlForIForPWK(i int) {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	// 连接yaml文件
	SdkConfig := viper.New()
	SdkConfig.AddConfigPath(examplePath)
	SdkConfig.SetConfigName(sdkConfigPWKExample)
	SdkConfig.SetConfigType("yaml")
	if err := SdkConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}

	// 修改viper实例内容
	hostIp := ParametersList.ChainInformationList.HostList[i]
	grpc := ParametersList.ChainInformationList.Grpc1List[i]
	i = i + 1
	// fmt.Println(host + grpc)
	messages := SdkConfig.Get("chain_client.nodes")
	chainNodesList := make([]ChainNodesYamlForPWK, len(messages.([]interface{})))
	for j, curNode := range messages.([]interface{}) {
		node := curNode.(map[interface{}]interface{})["conn_cnt"]
		chainNodesList[j].ConnCnt = node.(int)

		//node = curNode.(map[interface{}]interface{})["node_addr"]
		chainNodesList[j].NodeAddr = hostIp + ":" + grpc
	}
	SdkConfig.Set("chain_client.nodes", chainNodesList)

	orgId := SdkConfig.Get("chain_client.org_id")
	orgId = strings.Replace(orgId.(string), "org1", "org"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.org_id", orgId)

	userKeyFilePath := SdkConfig.Get("chain_client.user_sign_key_file_path")
	userKeyFilePath = strings.Replace(userKeyFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_sign_key_file_path", userKeyFilePath)

	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	SdkSavePath := newPath + NewSep + sdkConfigPWKFile + strconv.Itoa(i) + ".yml"
	errWrite := SdkConfig.WriteConfigAs(SdkSavePath)
	if errWrite != nil {
		return
	}
}

// UpdateYamlForI 修改编号i的SdkConfig配置文件
func UpdateYamlForI(i int) {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	// 连接yaml文件
	SdkConfig := viper.New()
	SdkConfig.AddConfigPath(examplePath)
	SdkConfig.SetConfigName(sdkConfigExample)
	SdkConfig.SetConfigType("yaml")
	if err := SdkConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}

	// 修改viper实例内容
	hostIp := ParametersList.ChainInformationList.HostList[i]
	grpc := ParametersList.ChainInformationList.Grpc1List[i]
	i = i + 1
	messages := SdkConfig.Get("chain_client.nodes")
	chainNodesList := make([]ChainNodesYaml, len(messages.([]interface{})))
	for j, curNode := range messages.([]interface{}) {
		node := curNode.(map[interface{}]interface{})["conn_cnt"]
		chainNodesList[j].ConnCnt = node.(int)
		node = curNode.(map[interface{}]interface{})["enable_tls"]
		chainNodesList[j].EnableTls = node.(bool)
		chainNodesList[j].NodeAddr = hostIp + ":" + grpc

		node = curNode.(map[interface{}]interface{})["tls_host_name"]
		chainNodesList[j].TlsHostName = node.(string)
		node = curNode.(map[interface{}]interface{})["trust_root_paths"]
		chainNodesList[j].TrustRootPaths = make([]string, len(node.([]interface{})))
		for k, trustRootPath := range node.([]interface{}) {
			newPath := trustRootPath.(string)
			newPath = strings.Replace(newPath, "node1", "node"+strconv.Itoa(i), -1)
			newPath = strings.Replace(newPath, "org1", "org"+strconv.Itoa(i), -1)
			chainNodesList[j].TrustRootPaths[k] = newPath
		}
	}
	SdkConfig.Set("chain_client.nodes", chainNodesList)
	orgId := SdkConfig.Get("chain_client.org_id")
	orgId = strings.Replace(orgId.(string), "org1", "org"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.org_id", orgId)

	userCrtFilePath := SdkConfig.Get("chain_client.user_crt_file_path")
	userCrtFilePath = strings.Replace(userCrtFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_crt_file_path", userCrtFilePath)

	userKeyFilePath := SdkConfig.Get("chain_client.user_key_file_path")
	userKeyFilePath = strings.Replace(userKeyFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_key_file_path", userKeyFilePath)

	userSignCrtFilePath := SdkConfig.Get("chain_client.user_sign_crt_file_path")
	userSignCrtFilePath = strings.Replace(userSignCrtFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_sign_crt_file_path", userSignCrtFilePath)

	userSignKeyFilePath := SdkConfig.Get("chain_client.user_sign_key_file_path")
	userSignKeyFilePath = strings.Replace(userSignKeyFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_sign_key_file_path", userSignKeyFilePath)

	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	SdkSavePath := newPath + NewSep + sdkConfigFile + strconv.Itoa(i) + ".yml"
	errWrite := SdkConfig.WriteConfigAs(SdkSavePath)
	if errWrite != nil {
		return
	}
}

// UpdateYamlForIForPK PK模式下修改编号i的SdkConfig配置文件
func UpdateYamlForIForPK(i int) {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	//fmt.Printf(examplePath)
	// 连接yaml文件
	SdkConfig := viper.New()
	SdkConfig.AddConfigPath(examplePath)
	SdkConfig.SetConfigName(sdkConfigPkExample)
	SdkConfig.SetConfigType("yaml")
	if err := SdkConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}

	// 修改viper实例内容
	hostIp := ParametersList.ChainInformationList.HostList[i]
	grpc := ParametersList.ChainInformationList.Grpc1List[i]
	i = i + 1
	messages := SdkConfig.Get("chain_client.nodes")
	chainNodesList := make([]ChainNodesYamlForPK, len(messages.([]interface{})))
	for j, curNode := range messages.([]interface{}) {
		node := curNode.(map[interface{}]interface{})["conn_cnt"]
		chainNodesList[j].ConnCnt = node.(int)
		chainNodesList[j].NodeAddr = hostIp + ":" + grpc
	}
	SdkConfig.Set("chain_client.nodes", chainNodesList)

	userKeyFilePath := SdkConfig.Get("chain_client.user_sign_key_file_path")
	userKeyFilePath = strings.Replace(userKeyFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	userKeyFilePath = strings.Replace(userKeyFilePath.(string), "admin1", "admin"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_sign_key_file_path", userKeyFilePath)

	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	SdkSavePath := newPath + NewSep + sdkConfigPkFile + strconv.Itoa(i) + ".yml"
	errWrite := SdkConfig.WriteConfigAs(SdkSavePath)
	if errWrite != nil {
		return
	}
}

// UpdateYaml 更新配置文件总函数
func UpdateYaml() {

	dir, _ := os.ReadDir(GetConfigFilePathByCaller(configFileName))
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{configFileName, d.Name()}...))
	}

	UpdateConstConfigYaml()

	examplePath := GetConfigFilePathByCaller(configExampleName)
	//fmt.Printf(examplePath)
	// 连接yaml文件
	ClientsConfig := viper.New()
	ClientsConfig.AddConfigPath(examplePath)
	ClientsConfig.SetConfigName(clientsExample)
	ClientsConfig.SetConfigType("yaml")
	if err := ClientsConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}
	ClientsConfig.Set("node_num", len(ParametersList.ChainInformationList.HostList))
	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	ClientsSavePath := newPath + NewSep + clientsConfigFile + ".yml"
	errWrite := ClientsConfig.WriteConfigAs(ClientsSavePath)
	if errWrite != nil {
		return
	}

	if ParametersList.ChainInformationList.Model == "PermissionWithCert" {
		for i := 0; i < len(ParametersList.ChainInformationList.HostList); i++ {
			UpdateYamlForI(i)
		}
	}
	if ParametersList.ChainInformationList.Model == "PermissionWithKey" {
		for i := 0; i < len(ParametersList.ChainInformationList.HostList); i++ {
			UpdateYamlForIForPWK(i)
		}
	}
	if ParametersList.ChainInformationList.Model == "Public" {
		for i := 0; i < len(ParametersList.ChainInformationList.HostList); i++ {
			UpdateYamlForIForPK(i)
		}
	}
}

// GenerateConstConfig 更新ConstConfig配置文件yml文件
func GenerateConstConfig(model, contractFunction, contractName, runTimeType, allParams string, climbTime, loopNum, sleepTime, threadNum int) string {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	// 连接yaml文件
	ConstConfig := viper.New()
	ConstConfig.AddConfigPath(examplePath)
	ConstConfig.SetConfigName(constConfigExample)
	ConstConfig.SetConfigType("yaml")
	if err := ConstConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}

	// 修改viper实例内容
	ConstConfig.Set("chain_information.model", model)
	ConstConfig.Set("contract_configurable_parameters.contract_method", contractFunction)
	ConstConfig.Set("contract_configurable_parameters.contract_name", contractName)

	// 根据智能合约函数名判断合约类型
	saveContain := strings.Contains(strings.ToLower(contractFunction), strings.ToLower("save"))
	findContain := strings.Contains(strings.ToLower(contractFunction), strings.ToLower("find"))
	if saveContain {
		ConstConfig.Set("contract_configurable_parameters.contract_type", "Claim")
	} else if findContain {
		ConstConfig.Set("contract_configurable_parameters.contract_type", "Query")
	}

	ConstConfig.Set("contract_configurable_parameters.runtime_type", runTimeType)
	ConstConfig.Set("pressure_configurable_parameters.climb_time", climbTime)
	ConstConfig.Set("pressure_configurable_parameters.loop_num", loopNum)
	ConstConfig.Set("pressure_configurable_parameters.sleep_time", sleepTime)
	ConstConfig.Set("pressure_configurable_parameters.thread_num", threadNum)
	ConstConfig.Set("contract_configurable_parameters.params", allParams)

	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	// 获取操作系统
	osKind := runtime.GOOS
	// 根据操作系统生成新分隔符
	var newSep string = "/"
	// windows系统
	if osKind == "windows" {
		newSep = "\\"
	}
	// linux系统
	if osKind == "linux" {
		newSep = "/"
	}
	ConstSavePath := newPath + newSep + constConfigFile + ".yml"
	errWrite := ConstConfig.WriteConfigAs(ConstSavePath)
	if errWrite != nil {
		return errWrite.Error()
	}
	return newPath
}

// UpdateYamlForIForPWKAPI pwk模式下修改编号i的SdkConfig配置文件
func UpdateYamlForIForPWKAPI(hostList, grpc1List []string, i, connCnt int) {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	// 连接yaml文件
	SdkConfig := viper.New()
	SdkConfig.AddConfigPath(examplePath)
	SdkConfig.SetConfigName(sdkConfigPWKExample)
	SdkConfig.SetConfigType("yaml")
	if err := SdkConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}

	// 修改viper实例内容
	hostIp := hostList[i]
	grpc := grpc1List[i]
	i = i + 1
	// fmt.Println(host + grpc)
	messages := SdkConfig.Get("chain_client.nodes")
	chainNodesList := make([]ChainNodesYamlForPWK, len(messages.([]interface{})))
	for j := range messages.([]interface{}) {
		chainNodesList[j].ConnCnt = connCnt

		//node = curNode.(map[interface{}]interface{})["node_addr"]
		chainNodesList[j].NodeAddr = hostIp + ":" + grpc
	}
	SdkConfig.Set("chain_client.nodes", chainNodesList)

	orgId := SdkConfig.Get("chain_client.org_id")
	orgId = strings.Replace(orgId.(string), "org1", "org"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.org_id", orgId)

	userKeyFilePath := SdkConfig.Get("chain_client.user_sign_key_file_path")
	userKeyFilePath = strings.Replace(userKeyFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_sign_key_file_path", userKeyFilePath)

	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	// 获取操作系统
	osKind := runtime.GOOS
	// 根据操作系统生成新分隔符
	var newSep string
	// windows系统
	if osKind == "windows" {
		newSep = "\\"
	}
	// linux系统
	if osKind == "linux" {
		newSep = "/"
	}
	SdkSavePath := newPath + newSep + sdkConfigPWKFile + strconv.Itoa(i) + ".yml"
	errWrite := SdkConfig.WriteConfigAs(SdkSavePath)
	if errWrite != nil {
		return
	}
}

// UpdateYamlForIAPI 修改编号i的SdkConfig配置文件
func UpdateYamlForIAPI(hostList, grpc1List []string, i, connCnt int) {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	// 连接yaml文件
	SdkConfig := viper.New()
	SdkConfig.AddConfigPath(examplePath)
	SdkConfig.SetConfigName(sdkConfigExample)
	SdkConfig.SetConfigType("yaml")
	if err := SdkConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}

	// 修改viper实例内容
	hostIp := hostList[i]
	grpc := grpc1List[i]
	i = i + 1
	messages := SdkConfig.Get("chain_client.nodes")
	chainNodesList := make([]ChainNodesYaml, len(messages.([]interface{})))
	for j, curNode := range messages.([]interface{}) {
		chainNodesList[j].ConnCnt = connCnt
		node := curNode.(map[interface{}]interface{})["enable_tls"]
		chainNodesList[j].EnableTls = node.(bool)
		chainNodesList[j].NodeAddr = hostIp + ":" + grpc

		node = curNode.(map[interface{}]interface{})["tls_host_name"]
		chainNodesList[j].TlsHostName = node.(string)
		node = curNode.(map[interface{}]interface{})["trust_root_paths"]
		chainNodesList[j].TrustRootPaths = make([]string, len(node.([]interface{})))
		for k, trustRootPath := range node.([]interface{}) {
			newPath := trustRootPath.(string)
			newPath = strings.Replace(newPath, "node1", "node"+strconv.Itoa(i), -1)
			newPath = strings.Replace(newPath, "org1", "org"+strconv.Itoa(i), -1)
			chainNodesList[j].TrustRootPaths[k] = newPath
		}
	}
	SdkConfig.Set("chain_client.nodes", chainNodesList)
	orgId := SdkConfig.Get("chain_client.org_id")
	orgId = strings.Replace(orgId.(string), "org1", "org"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.org_id", orgId)

	userCrtFilePath := SdkConfig.Get("chain_client.user_crt_file_path")
	userCrtFilePath = strings.Replace(userCrtFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_crt_file_path", userCrtFilePath)

	userKeyFilePath := SdkConfig.Get("chain_client.user_key_file_path")
	userKeyFilePath = strings.Replace(userKeyFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_key_file_path", userKeyFilePath)

	userSignCrtFilePath := SdkConfig.Get("chain_client.user_sign_crt_file_path")
	userSignCrtFilePath = strings.Replace(userSignCrtFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_sign_crt_file_path", userSignCrtFilePath)

	userSignKeyFilePath := SdkConfig.Get("chain_client.user_sign_key_file_path")
	userSignKeyFilePath = strings.Replace(userSignKeyFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_sign_key_file_path", userSignKeyFilePath)

	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	// 获取操作系统
	osKind := runtime.GOOS
	// 根据操作系统生成新分隔符
	var newSep string
	// windows系统
	if osKind == "windows" {
		newSep = "\\"
	}
	// linux系统
	if osKind == "linux" {
		newSep = "/"
	}
	SdkSavePath := newPath + newSep + sdkConfigFile + strconv.Itoa(i) + ".yml"
	errWrite := SdkConfig.WriteConfigAs(SdkSavePath)
	if errWrite != nil {
		return
	}
}

// UpdateYamlForIForPK PK模式下修改编号i的SdkConfig配置文件
func UpdateYamlForIForPKAPI(hostList, grpc1List []string, i, connCnt int) {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	//fmt.Printf(examplePath)
	// 连接yaml文件
	SdkConfig := viper.New()
	SdkConfig.AddConfigPath(examplePath)
	SdkConfig.SetConfigName(sdkConfigPkExample)
	SdkConfig.SetConfigType("yaml")
	if err := SdkConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}

	// 修改viper实例内容
	hostIp := hostList[i]
	grpc := grpc1List[i]
	i = i + 1
	messages := SdkConfig.Get("chain_client.nodes")
	chainNodesList := make([]ChainNodesYamlForPK, len(messages.([]interface{})))
	for j := range messages.([]interface{}) {
		chainNodesList[j].ConnCnt = connCnt
		chainNodesList[j].NodeAddr = hostIp + ":" + grpc
	}
	SdkConfig.Set("chain_client.nodes", chainNodesList)

	userKeyFilePath := SdkConfig.Get("chain_client.user_sign_key_file_path")
	userKeyFilePath = strings.Replace(userKeyFilePath.(string), "node1", "node"+strconv.Itoa(i), -1)
	userKeyFilePath = strings.Replace(userKeyFilePath.(string), "admin1", "admin"+strconv.Itoa(i), -1)
	SdkConfig.Set("chain_client.user_sign_key_file_path", userKeyFilePath)

	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	// 获取操作系统
	osKind := runtime.GOOS
	// 根据操作系统生成新分隔符
	var newSep string
	// windows系统
	if osKind == "windows" {
		newSep = "\\"
	}
	// linux系统
	if osKind == "linux" {
		newSep = "/"
	}
	SdkSavePath := newPath + newSep + sdkConfigPkFile + strconv.Itoa(i) + ".yml"
	errWrite := SdkConfig.WriteConfigAs(SdkSavePath)
	if errWrite != nil {
		return
	}
}

func GenerateClintConfig(nodeNum int) string {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	//fmt.Printf(examplePath)
	// 连接yaml文件
	ClientsConfig := viper.New()
	ClientsConfig.AddConfigPath(examplePath)
	ClientsConfig.SetConfigName(clientsExample)
	ClientsConfig.SetConfigType("yaml")
	if err := ClientsConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}
	ClientsConfig.Set("node_num", nodeNum)
	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	// 获取操作系统
	osKind := runtime.GOOS
	// 根据操作系统生成新分隔符
	var newSep string = "/"
	// windows系统
	if osKind == "windows" {
		newSep = "\\"
	}
	// linux系统
	if osKind == "linux" {
		newSep = "/"
	}
	ClientsSavePath := newPath + newSep + clientsConfigFile + ".yml"
	errWrite := ClientsConfig.WriteConfigAs(ClientsSavePath)
	if errWrite != nil {
		return ""
	}
	return ClientsSavePath
}

func GenerateSdkConfig(hostList, grpc1List []string, model string, connCnt int) string {
	examplePath := GetConfigFilePathByCaller(configExampleName)
	//fmt.Printf(examplePath)
	// 连接yaml文件
	ClientsConfig := viper.New()
	ClientsConfig.AddConfigPath(examplePath)
	ClientsConfig.SetConfigName(clientsExample)
	ClientsConfig.SetConfigType("yaml")
	if err := ClientsConfig.ReadInConfig(); err != nil {
		logger.Logger.Println(err)
	}
	ClientsConfig.Set("node_num", len(hostList))
	//生成yml文件
	newPath := GetConfigFilePathByCaller(configFileName)
	// 获取操作系统
	osKind := runtime.GOOS
	// 根据操作系统生成新分隔符
	var newSep string
	// windows系统
	if osKind == "windows" {
		newSep = "\\"
	}
	// linux系统
	if osKind == "linux" {
		newSep = "/"
	}
	ClientsSavePath := newPath + newSep + clientsConfigFile + ".yml"
	errWrite := ClientsConfig.WriteConfigAs(ClientsSavePath)
	if errWrite != nil {
		return ""
	}

	if model == "PermissionWithCert" {
		for i := 0; i < len(hostList); i++ {
			UpdateYamlForIAPI(hostList, grpc1List, i, connCnt)
		}
	}
	if model == "PermissionWithKey" {
		for i := 0; i < len(hostList); i++ {
			UpdateYamlForIForPWKAPI(hostList, grpc1List, i, connCnt)
		}
	}
	if model == "Public" {
		for i := 0; i < len(hostList); i++ {
			UpdateYamlForIForPKAPI(hostList, grpc1List, i, connCnt)
		}
	}
	return GetConfigFilePathByCaller(configFileName)
}

// UpdateConfigNew 更新长安链配置文件、压测任务参数配置文件
func UpdateConfigNew(rdb *redis.Client, getKey string) {
	content := ReadRedis(rdb, getKey)
	err := ParseJson(content)
	if err != nil {
		return
	}
	DownLoadRun()
	UpdateYaml()
}

