/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package datahandler

import (
	logger "chain-performance-test/log"
	"chain-performance-test/testdata"
	"chain-performance-test/traffic"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
	"strings"
)

var (
	ipSep       string = "||" //多个ip间分隔符
	ParaListSep string = "||" //多个参数列表间分隔符
	ParaSep     string = ":"  //参数列表内参数间分隔符

	//路径分隔符
	OldSep string = "/"
	//使用系统分隔符
	NewSep string = testdata.NewSep

	// ParametersList json结构体
	ParametersList Parameters

	// TestResult 结果结构体
	TestResult Result
)

// Result 返回结果结构体定义
type Result struct {
	TxNum  int // 总交易数
	CTxNum int // 总交易上链数
	TPS    float64
	CTPS   float64
	QTPS   float64
}

// 链信息结构体
type ChainInformation struct {
	Host       string `json:"host"`       //虚拟机IP
	Grpc1      string `json:"grpc1"`      //虚拟机端口号1
	Grpc2      string `json:"grpc2"`      //虚拟机端口号2
	Model      string `json:"model"`      //身份权限管理
	ConfigFile string `json:"configFile"` //链和节点的配置文件路径
	HostList   []string
	Grpc1List  []string
}

// 合约文件信息结构体
type ContractConfigurableParameters struct {
	RunTimeType                   string `json:"contractType"`     //智能合约类型
	ContractFunction              string `json:"contractFunction"` //智能合约函数
	ContractName                  string `json:"contractName"`     //智能合约名称
	ContractFunctionParameters    string `json:"allParams"`
	ContractFunctionParametersMap map[string]string
}

// 配置参数结构体
type PressureConfigurableParameters struct {
	LoopNum   int `json:"loopNum"`   //压测并发次数
	ThreadNum int `json:"threadNum"` //并发线程数
	ClimbTime int `json:"climbTime"` //线程爬坡时间
	SleepTime int `json:"sleepTime"` //线程发送间隔
}

// //流量分析参数结构体
type TrafficAnalysisParameters struct {
	SourceIp   string `json:"sourceIp"`   //源IP
	SourcePort string `json:"sourcePort"` //源端口号
	Protocol   string `json:"protocol"`   //协议
}

type Parameters struct {
	ChainInformationList          ChainInformation               //压测链信息列表
	ContractParametersList        ContractConfigurableParameters //合约参数列表
	PressureParametersList        PressureConfigurableParameters //压测参数列表
	TrafficAnalysisParametersList TrafficAnalysisParameters      // 流量分析参数列表
}

// 连接redis
func ConnectRedis(addr string, password string) *redis.Client {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // 密码
		DB:       0,        // 数据库
		PoolSize: 20,       // 连接池大小
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Logger.Printf("连接redis出错，错误信息: %v\n", err)
	} else {
		logger.Logger.Println("连接redis成功！")
	}
	return rdb
}

// GetKey 获取redis中缓存的指定 key
func GetKey(rdb *redis.Client, keys string) string {
	var ctx = context.Background()
	// 匹配key
	getKeys, err := rdb.Keys(ctx, keys).Result()
	if err != nil {
		logger.Logger.Println(err)
	}
	//获取匹配到的所有key中第一个key
	var getKey string
	if len(getKeys) != 0 {
		getKey = getKeys[0]
		logger.Logger.Println("本次处理的redis缓存中的key：", getKey)
	} else {
		getKey = ""
	}
	logger.Logger.Println("============================================================")
	return getKey
}

// ReadRedis 根据指定key从redis中获取对应内容
func ReadRedis(rdb *redis.Client, getKey string) string {
	var ctx = context.Background()
	// 获取对应key的内存  taskQueue  go_redis_test
	val, err := rdb.Get(ctx, getKey).Result()
	if err != nil {
		logger.Logger.Println(err)
	}
	content := strings.Replace(val, "\\", "", -1)
	content = content[1 : len(content)-1] //调整content成可解析为json的字符串
	logger.Logger.Printf("key：%v\n", content)
	return content

	//return ParametersList
}

// 根据指定key删除redis缓存
func DeleteRedisCache(rdb *redis.Client, deleteKey string) {
	var ctx = context.Background()
	//删除缓存项
	_, err := rdb.Del(ctx, deleteKey).Result()
	if err != nil {
		logger.Logger.Println(err)
	}
}

// ParseJson redis 中 key 对应 value 参数解析
func ParseJson(content string) error {
	//解析 redisData  []byte(val1)
	ParametersList = Parameters{} //初始化
	//解析到对应参数结构体
	errs := json.Unmarshal([]byte(content), &ParametersList.ChainInformationList)
	if errs != nil {
		logger.Logger.Println("json unmarshal error:", errs)
		return errs
	}

	errs = json.Unmarshal([]byte(content), &ParametersList.ContractParametersList)
	if errs != nil {
		logger.Logger.Println("json unmarshal error:", errs)
		return errs
	}
	errs = json.Unmarshal([]byte(content), &ParametersList.PressureParametersList)
	if errs != nil {
		logger.Logger.Println("json unmarshal error:", errs)
		return errs
	}

	errs = json.Unmarshal([]byte(content), &ParametersList.TrafficAnalysisParametersList)
	if errs != nil {
		logger.Logger.Println("json unmarshal error:", errs)
		return errs
	}

	//解析IP
	HostListTemp := strings.Split(ParametersList.ChainInformationList.Host, ipSep)
	ParametersList.ChainInformationList.HostList = make([]string, len(HostListTemp)-1)
	for i := 0; i < len(HostListTemp)-1; i++ {
		ParametersList.ChainInformationList.HostList[i] = HostListTemp[i]
	}
	Grpc1ListTemp := strings.Split(ParametersList.ChainInformationList.Grpc1, ipSep)
	ParametersList.ChainInformationList.Grpc1List = make([]string, len(Grpc1ListTemp)-1)
	for i := 0; i < len(Grpc1ListTemp)-1; i++ {
		ParametersList.ChainInformationList.Grpc1List[i] = Grpc1ListTemp[i]
	}

	//解析参数map
	ParametersList.ContractParametersList.ContractFunctionParametersMap = make(map[string]string)
	FunctionParametersList := strings.Split(ParametersList.ContractParametersList.ContractFunctionParameters, ParaListSep)
	for i := 0; i < len(FunctionParametersList); i++ {
		FunctionParameters := strings.Split(FunctionParametersList[i], ParaSep)
		ParametersList.ContractParametersList.ContractFunctionParametersMap[FunctionParameters[0]] = FunctionParameters[1]
	}

	//修改DOCKER-GO下划线
	ParametersList.ContractParametersList.RunTimeType = strings.Replace(ParametersList.ContractParametersList.RunTimeType, "-", "_", -1)

	logger.Logger.Println("key参数处理结果：", ParametersList)
	return nil
}

// PushRedis 结果字段推入redis
func PushRedis(getKey string, rdb *redis.Client, trafficInfo map[traffic.SourceInfo][]traffic.DestInfo) {
	var ctx = context.Background()
	result := make(map[string]string)
	result["TxNum"] = strconv.Itoa(TestResult.TxNum)
	result["CTxNum"] = strconv.Itoa(TestResult.CTxNum)
	result["TPS"] = strconv.Itoa(int(TestResult.TPS))
	result["CTPS"] = strconv.Itoa(int(TestResult.CTPS))
	result["QTPS"] = strconv.Itoa(int(TestResult.QTPS))
	for source, dest := range trafficInfo {
		//获取结果
		sourceInfo, err := json.Marshal(source)
		logger.Logger.Println("err = ", err)
		destInfo, err := json.Marshal(dest)
		if err != nil {
			logger.Logger.Println("err = ", err)
			return
		}
		result["source"] = string(sourceInfo)
		result["port"] = string(destInfo)
	}

	redisData, err := json.Marshal(result)
	if err != nil {
		logger.Logger.Println("err = ", err)
		return
	}
	logger.Logger.Println("redisData = ", string(redisData))
	//生成结果key值
	pre := getKey[0:21]
	next := string(getKey[21:])
	num := fmt.Sprintf("%08v", rand.Intn(99999999))
	resKey := pre + "_" + "res" + "_" + next + "_" + num
	logger.Logger.Println("结果返回：", resKey)
	err = rdb.Set(ctx, resKey, string(redisData), 0).Err()
	if err != nil {
		logger.Logger.Println(err)
	}
}
