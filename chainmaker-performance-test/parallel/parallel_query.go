/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package parallel

import (
	"chain-performance-test/chainclient"
	"chain-performance-test/datahandler"
	logger "chain-performance-test/log"
	sub "chain-performance-test/subservice"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"context"
	"strconv"
	"time"
)

// ContinuousQueryPressureTester 创建持续查询压测的接口
type ContinuousQueryPressureTester interface {
	QueryPressureTest(*sdk.ChainClient) error
}

type QueryPressureTester struct{}

// QueryTPSCTPSWithSub 计算TPS、CTPS 边压测边订阅
func QueryTPSCTPSWithSub(pressureTester ContinuousQueryPressureTester,
	clientCreator chainclient.ClientCreator, blockSubscriber sub.BlockSubscriber) error {

	logger.Logger.Println("============ 开始查询压测 ============")

	countClimb := ThreadNum / 10
	if countClimb == 0 {
		countClimb = 1
	}
	// 设置线程启动时间
	interval := time.Duration(ClimbTime/countClimb) * time.Second

	// 开始压测计时
	timeStart = time.Now()
	dateStart := timeStart.Format("2006-01-02 15:04:05")

	ctxForTest, cancelForTest := context.WithCancel(context.Background())

	for index := 0; index < ThreadNum; {
		for j := 0; j < 10; j++ {
			// 多节点加压
			for i := 0; i < Clientslen; i++ {
				// 持续的执行压测方法
				// 设置需要等待的协程数量
				Wg.Add(1)
				k := i
				go func(ctx context.Context) {
					defer Wg.Done()
					select {
					case <-ctx.Done():
						// 如果收到取消信号，直接返回不执行任务
						return
					default:
						defer func() {
							if err := recover(); err != nil {
								cancelForTest()
								logger.Logger.Println("Error:", err)
							}
						}()
						err := pressureTester.QueryPressureTest(Clients1[k])
						if err != nil {
							cancelForTest()
							logger.Logger.Panicln(err)
						}
					}
				}(ctxForTest)
			}
			index++
			if index >= ThreadNum {
				break
			}
		}
		time.Sleep(interval)
	}

	Wg.Wait()
	cancelForTest()

	// 结束压测计时
	timeEnd = time.Now()
	dateEnd := timeEnd.Format("2006-01-02 15:04:05")

	// txNum交易数
	txNum := LoopNum * ThreadNum * Clientslen
	datahandler.TestResult.TxNum = txNum
	count := float64(txNum)

	timeResult := float64((timeEnd.UnixNano()-timeStart.UnixNano())/1e6) / 1000.0
	logger.Logger.Printf("并发数: %d\n", ThreadNum)
	logger.Logger.Println("timeStart:", dateStart, "timeEnd:", dateEnd)
	logger.Logger.Println("txNum:", txNum, "Duration:", strconv.FormatFloat(timeResult, 'g', 30, 32)+" s",
		"TPS:", count/timeResult)

	datahandler.TestResult.TPS = count / timeResult
	sub.BlockInfos = []sub.BlockInfo{}

	// 消息订阅上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	// 建立连接
	if Model == "PermissionWithCert" {
		Client, err = clientCreator.CreateClientWithConfig(sdkConfigPaths[0])

	} else if Model == "Public" {
		Client, err = clientCreator.CreateClientWithConfig(sdkPKConfigPaths[0])

	} else if Model == "PermissionWithKey" {
		Client, err = clientCreator.CreateClientWithConfig(SdkPWKConfigPaths[0])
	}
	if err != nil {
		logger.Logger.Println(err)
		return err
	}

	logger.Logger.Println("------- start a subscribe service to monitor chain block information -------")
	// 创建一个消息订阅服务
	blockSubscriber.SetClient(Client)
	err = blockSubscriber.Run(ctx, BlockStart, -1)
	if err != nil {
		logger.Logger.Println(err)
	}
	blockSubscriber.Close()

	// QTPS计算
	txCounts = 0

	// 实际上链数计算
	for num, Info := range sub.BlockInfos {
		if num == 0 {
			timeStart = Info.Timestamp
		}
		txCounts = txCounts + Info.TxCount
		timeEnd = Info.Timestamp
	}

	datahandler.TestResult.CTxNum = int(txCounts)
	count = float64(txCounts)
	logger.Logger.Printf("成功交易数：%f\n", count)
	timeResultQtps := float64((timeEnd.UnixNano()-timeStart.UnixNano())/1e6) / 1000.0
	logger.Logger.Println("Duration:", strconv.FormatFloat(timeResultQtps, 'g', 30, 32)+" s",
		"QTPS:", (count)/timeResultQtps)
	datahandler.TestResult.QTPS = (count) / timeResultQtps
	return nil
}

// QueryPressureTest 持续的执行压测方法
func (q QueryPressureTester) QueryPressureTest(client *sdk.ChainClient) error {

	for i := 0; i < LoopNum; i++ {
		err := UserContractQueryInvoke(client, contractMethod, false)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(SleepTime) * time.Millisecond)
	}
	return nil
}

// UserContractQueryInvoke 调用查询合约
func UserContractQueryInvoke(client *sdk.ChainClient, method string, withSyncResult bool) error {

	params := RandParams(datahandler.ParametersList.ContractParametersList.ContractFunctionParametersMap)
	_, err := chainclient.InvokeUserContract(client, ContractName, method, "", params, withSyncResult)
	if err != nil {
		return err
	}
	return nil
}
