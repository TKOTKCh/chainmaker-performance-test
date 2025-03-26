/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package parallel

import (
	"chain-performance-test/chainclient"
	"chain-performance-test/datahandler"
	sub "chain-performance-test/subservice"
	"chain-performance-test/traffic"
	"fmt"
)

// 通信流量信息
var TrafficInfo map[traffic.SourceInfo][]traffic.DestInfo

// ProcessPressureTest 处理敏捷测试任务
func ProcessPressureTest(loopi int) datahandler.Result {
	// 压测前准备
	InitBeforeTest(loopi)
	fmt.Println("压测准备完成")
	//go func() {
	//	// 抓取通信流量
	//	pcapPath := traffic.CatchTraffic(testdata.ServerHost, testdata.ServerPort, testdata.ServerUser, testdata.ServerPasswd, testdata.ServerListenPort)
	//	// 分析通信流量
	//	TrafficInfo = traffic.Analyse(pcapPath)
	//	logger.Logger.Println("trafficInfo1:", TrafficInfo)
	//}()

	// 压测并计算相应的指标
	var pressureTester1 ClaimPressureTester
	//var pressureTester2 QueryPressureTester
	var clientCreator chainclient.ClientCreate
	var blockSubscriber sub.SubscribeBlock
	err := ClaimTPSCTPSWithSub(pressureTester1, clientCreator, &blockSubscriber)

	if err != nil {
		return datahandler.Result{}
	}
	//switch ContractType {
	//case "Claim":
	//	err := ClaimTPSCTPSWithSub(pressureTester1, clientCreator, &blockSubscriber)
	//	if err != nil {
	//		return datahandler.Result{}
	//	}
	//case "Query":
	//	err := QueryTPSCTPSWithSub(pressureTester2, clientCreator, &blockSubscriber)
	//	if err != nil {
	//		return datahandler.Result{}
	//	}
	//case "Asset":
	//	// todo
	//}
	return datahandler.TestResult
}
