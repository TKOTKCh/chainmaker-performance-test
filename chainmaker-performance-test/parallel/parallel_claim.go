/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package parallel

import (
	"bufio"
	"chain-performance-test/chainclient"
	"chain-performance-test/datahandler"
	logger "chain-performance-test/log"
	sub "chain-performance-test/subservice"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ContinuousClaimPressureTester 创建持续压测的接口
type ContinuousClaimPressureTester interface {
	ClaimPressureTest(*sdk.ChainClient, bool) (error, []string)
}

func addMaps(m1, m2 map[int64]int64) map[int64]int64 {
	for key, value := range m2 {
		m1[key] += value // 如果键存在，累加值；如果键不存在，初始化为 value
	}
	return m1
}

type ClaimPressureTester struct{}

// 定义交易数据结构
type TransactionRunningInfo struct {
	contractName          string
	contractMethod        string
	RuntimeType           string
	RuntimeContractResult int
	StartTime             int64
	EndTime               int64
	ExecutionTime         float64
}

// 处理日志文件
func getLastTimestamp(logFilePath string) (map[string]TransactionRunningInfo, int64, int64) {
	// 正则表达式匹配交易信息
	logPattern := regexp.MustCompile(`tx id:(?P<txid>[a-f0-9]+),\s*` +
		`contractName:(?P<contractName>[^,]+),\s*` +
		`contractMethod:(?P<contractMethod>[^,]+),\s*` +
		`runtime type:(?P<runtime_type>\w+),\s*` +
		`runtimeContractResult:(?P<runtimeContractResult>\d+)\s*,\s*` +
		`startTime:(?P<startTime>\d+),\s*` +
		`endTime:(?P<endTime>\d+),\s*` +
		`executionTime:(?P<executionTime>[\d.]+) s`)

	txData := make(map[string]TransactionRunningInfo)
	var minStartTime int64 = 1<<63 - 1 // 设为 int64 最大值
	var maxEndTime int64 = 0

	// 打开日志文件
	file, err := os.Open(logFilePath)
	if err != nil {
		fmt.Println("无法打开文件:", logFilePath, "错误:", err)
		return nil, 0, 0
	}
	defer file.Close()
	cnt := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "vm_factory.go:417") { // 仅匹配包含特定标识的行
			matches := logPattern.FindStringSubmatch(line)
			if matches != nil {
				//第一个匹配的是安装合约执行的结果，跳过
				if cnt == 0 {
					cnt++
					continue
				}
				// 提取匹配组
				txid := matches[1]
				thisContractName := matches[2]
				thisContractMethod := matches[3]
				runtimeType := matches[4]
				runtimeContractResult, _ := strconv.Atoi(matches[5])
				startTime, _ := strconv.ParseInt(matches[6], 10, 64)
				endTime, _ := strconv.ParseInt(matches[7], 10, 64)
				executionTime, _ := strconv.ParseFloat(matches[8], 64)
				//如果合约执行结果失败跳过
				if runtimeContractResult != 0 || thisContractName != ContractName || thisContractMethod != contractMethod {
					continue
				}
				// 存入字典
				txData[txid] = TransactionRunningInfo{
					contractName:          ContractName,
					contractMethod:        contractMethod,
					RuntimeType:           runtimeType,
					RuntimeContractResult: runtimeContractResult,
					StartTime:             startTime,
					EndTime:               endTime,
					ExecutionTime:         executionTime,
				}

				// 更新最小开始时间 & 最大结束时间
				if startTime < minStartTime {
					minStartTime = startTime
				}
				if endTime > maxEndTime {
					maxEndTime = endTime
				}
			}
		}
	}

	if len(txData) > 0 && maxEndTime > minStartTime {

	} else {
		fmt.Printf("文件: %s, 没有有效交易数据\n", logFilePath)
	}

	return txData, minStartTime, maxEndTime
}

func copyFile(srcPath, distDir string) error {
	// 获取当前时间，并格式化为 "20060102150405" (YYYYMMDDHHmmss)
	timeStr := time.Now().Format("20060102150405")
	filename := "system.log." + timeStr
	distPath := filepath.Join(distDir, filename)

	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("无法打开源文件 %s: %v", srcPath, err)
	}
	defer srcFile.Close()

	// 创建目标文件
	distFile, err := os.Create(distPath)
	if err != nil {
		return fmt.Errorf("无法创建目标文件 %s: %v", distPath, err)
	}
	defer distFile.Close()

	// 复制文件内容
	_, err = io.Copy(distFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件失败: %v", err)
	}

	// 确保数据写入磁盘
	err = distFile.Sync()
	if err != nil {
		return fmt.Errorf("写入磁盘失败: %v", err)
	}

	fmt.Printf("文件%s 成功复制到 %s\n", srcPath, distPath)
	logger.Logger.Printf("文件%s 成功复制到 %s\n", srcPath, distPath)
	return nil
}

func ClaimTPSCTPSWithSub(pressureTester ContinuousClaimPressureTester,
	clientCreator chainclient.ClientCreator, blockSubscriber sub.BlockSubscriber) error {

	logger.Logger.Println("============ 开始压测 ============")

	countClimb := ThreadNum / 10
	if countClimb == 0 {
		countClimb = 1
	}
	// 设置线程启动时间
	interval := time.Duration(ClimbTime/countClimb) * time.Second

	// 开始压测计时
	//pTimeStart := time.Now()
	//pDateStart := pTimeStart.Format("2006-01-02 15:04:05")

	ctxForTest, cancelForTest := context.WithCancel(context.Background())

	//txIdMap := make(map[int][]string)
	//var totalTxIds []string
	for index := 0; index < ThreadNum; {
		for j := 0; j < 10; j++ {
			// 多节点加压
			for i := 0; i < Clientslen; i++ {
				// 持续的执行压测方法
				// 设置需要等待的协程数量
				k := i
				Wg.Add(1)
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
						err, _ := pressureTester.ClaimPressureTest(Clients1[k], false)

						//err, txIds := pressureTester.ClaimPressureTest(Clients1[k], false)
						//totalTxIds = append(totalTxIds, txIds...)
						//if existingTxIds, exists := txIdMap[k]; exists {
						//	txIdMap[k] = append(existingTxIds, txIds...) // 追加新的 txIds
						//} else {
						//	// 如果 k 不存在，则直接存入 txIds
						//	txIdMap[k] = txIds
						//}

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
	if AddOption == "no" {
		time.Sleep(40 * time.Second)
	} else {
		time.Sleep(40 * time.Second)
	}

	// 构建文件名
	originSrcDir := "/home/chenhang/WorkSpace/chainmaker-go/build/release/chainmaker-v2.3.0-wx-org1.chainmaker.org/log/"
	now := time.Now()
	// 格式化时间为文件名中的格式（例如 YYYYMMDDHH）
	timeStr := now.Format("2006010215")
	filename := "system.log." + timeStr

	minStartTime := int64(math.MaxInt64)
	maxEndTime := int64(0)
	lastSrcPath := ""
	txNum := 0
	orgNum := Clientslen
	for i := 1; i <= orgNum; i++ {
		orgNum := "org" + strconv.Itoa(i)                          // 生成 orgX
		srcDir := strings.Replace(originSrcDir, "org1", orgNum, 1) // 替换 org1 为 orgX
		srcPath := filepath.Join(srcDir, filename)

		currentTxData, currentMinStartTime, currentMaxEndTime := getLastTimestamp(srcPath)
		currentTxNum := len(currentTxData)
		fmt.Println(currentMinStartTime, currentMaxEndTime, currentTxNum, float64(currentTxNum)/(float64(currentMaxEndTime-currentMinStartTime)/1e9))
		// 找到最大的 maxEndTime
		if currentMaxEndTime > maxEndTime {
			maxEndTime = currentMaxEndTime
			lastSrcPath = srcPath
		}
		if currentMinStartTime < minStartTime {
			minStartTime = currentMinStartTime
		}
		txNum = txNum + len(currentTxData)
	}
	distDir := "/home/chenhang/WorkSpace/chainmaker-test-toolkit-master/chainmaker-performance-test/log/experimentlog"
	copyFile(lastSrcPath, distDir)
	//timeStart := pTimeStart
	timeStart := time.Unix(0, minStartTime)
	dateStart := timeStart.Format("2006-01-02 15:04:05")
	timeEnd := time.Unix(0, maxEndTime)
	dateEnd := timeEnd.Format("2006-01-02 15:04:05")
	//fmt.Println(timeStart.UnixNano(), timeEnd.UnixNano())
	// txNum交易数
	//txNum := LoopNum * ThreadNum * Clientslen
	// txNum为实际虚拟机执行成功的交易数,如果全都正常完成的话等于loop_num*thread_num*Clientslen
	count := float64(txNum)
	timeResult := float64((timeEnd.UnixNano()-timeStart.UnixNano())/1e6) / 1000.0

	datahandler.TestResult.TxNum = txNum
	datahandler.TestResult.TPS = count / timeResult

	logger.Logger.Printf("并发数: %d\n", ThreadNum)
	logger.Logger.Println("timeStart:", dateStart, timeStart.UnixNano(), "timeEnd:", dateEnd, timeEnd.UnixNano())
	logger.Logger.Println("txNum:", txNum, "Duration:", strconv.FormatFloat(timeResult, 'g', 30, 32)+" s",
		"TPS:", count/timeResult)
	//totalTimeMap := make(map[int64]int64)
	//for key, _ := range txIdMap {
	//	addMaps(totalTimeMap, checkTime(Clients1[key], totalTxIds))
	//}
	//var minKey, maxKey int64
	//first := true
	//
	//// 遍历 map
	//for k := range totalTimeMap {
	//	if first {
	//		minKey, maxKey = k, k
	//		first = false
	//	} else {
	//		if k < minKey {
	//			minKey = k
	//		}
	//		if k > maxKey {
	//			maxKey = k
	//		}
	//	}
	//}
	//fmt.Println(minKey, maxKey, float64(txNum)/(float64(maxKey-minKey)), "timeStart:", dateStart, timeStart.UnixNano(), "timeEnd:", dateEnd, timeEnd.UnixNano())
	fmt.Println("timeStart:", dateStart, timeStart.UnixNano(), "timeEnd:", dateEnd, timeEnd.UnixNano())
	return nil
}
func (p ClaimPressureTester) ClaimPressureTest(client *sdk.ChainClient, withSyncResult bool) (error, []string) {
	txIds := []string{}
	for i := 0; i < LoopNum; i++ {
		err, resp := UserContractClaimInvoke(client, contractMethod, withSyncResult, datahandler.ParametersList.ContractParametersList.ContractFunctionParametersMap)
		if err != nil {
			return err, txIds
		}
		txIds = append(txIds, resp.TxId)
		time.Sleep(time.Duration(SleepTime) * time.Millisecond)
	}
	return nil, txIds
}

// UserContractClaimInvoke 调用合约
func UserContractClaimInvoke(client *sdk.ChainClient, method string, withSyncResult bool, FunctionParametersMap map[string]string) (error, *common.TxResponse) {

	params := RandParams(FunctionParametersMap)
	resp, err := chainclient.InvokeUserContract(client, ContractName, method, "", params, withSyncResult)
	if err != nil {
		return err, resp
	}
	return nil, resp
}

func checkTime(client *sdk.ChainClient, txIds []string) map[int64]int64 {
	timeMap := map[int64]int64{}
	resultMap := map[string]int64{}
	wrongMap := map[string]int64{}
	for _, value := range txIds {
		transactionInfo, err := client.GetTxByTxId(value)

		if err != nil {
			wrongMap["err"] += 1
		} else {

			resultMap[string(transactionInfo.Transaction.Result.ContractResult.Result)] += 1
			tempTime := transactionInfo.BlockTimestamp
			if tempvalue, exist := timeMap[tempTime]; exist {
				timeMap[tempTime] = tempvalue + 1
			} else {
				timeMap[tempTime] = 1
			}

		}
	}
	fmt.Println(resultMap)
	return timeMap
}
