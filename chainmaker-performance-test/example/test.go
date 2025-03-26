/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	logger "chain-performance-test/log"
	"chain-performance-test/parallel"
	"fmt"
	"os/exec"
	"time"
)

type TestResult struct {
	TxNum  int
	CTxNum int
	TPS    float64
	CTPS   float64
}

// 示例代码
func main() {
	fmt.Println("import finish")
	// 移动上次生成的日志文件,可能会覆盖，请注意保存log下日志
	parallel.ProcessLog()
	// 创建一个日志文件
	logger.BeginLog()
	fmt.Println("日志创建成功")

	// 处理敏捷测试任务
	var totalResult TestResult
	testCount := 5
	scriptDir := "/home/chenhang/WorkSpace/chainmaker-go/scripts"
	for i := 0; i < testCount; i++ {
		fmt.Printf("开始第 %d 轮测试...\n", i+1)
		result := parallel.ProcessPressureTest(i)

		// 累加各项数值
		totalResult.TxNum += result.TxNum
		totalResult.CTxNum += result.CTxNum
		totalResult.TPS += result.TPS
		totalResult.CTPS += result.CTPS
		fmt.Printf("第 %d 轮测试结果: %+v\n", i+1, result)
		// 先停止集群
		cmdStop := exec.Command("./cluster_quick_stop.sh", "clean")
		cmdStop.Dir = scriptDir // 设置工作目录
		cmdStop.Stdout = nil
		cmdStop.Stderr = nil
		err := cmdStop.Run()
		if err != nil {
			fmt.Println("执行 cluster_quick_stop.sh 失败:", err)
		}

		// 再启动集群
		cmdStart := exec.Command("./cluster_quick_start.sh", "normal")
		cmdStart.Dir = scriptDir // 设置工作目录
		cmdStart.Stdout = nil
		cmdStart.Stderr = nil
		err = cmdStart.Run()
		if err != nil {
			fmt.Println("执行 cluster_quick_start.sh 失败:", err)
		}
		if i < testCount-1 {
			fmt.Println("暂停 10 秒后继续下一轮测试...")
			// 等待 10 秒
			time.Sleep(15 * time.Second)
		}
	}

	// 计算平均值
	avgResult := TestResult{
		TxNum:  totalResult.TxNum / testCount,
		CTxNum: totalResult.CTxNum / testCount,
		TPS:    totalResult.TPS / float64(testCount),
		CTPS:   totalResult.CTPS / float64(testCount),
	}

	fmt.Println("===== 平均测试结果 =====")
	fmt.Printf("交易数: %d,TPS: %.2f, CTPS: %.2f\n", avgResult.TxNum, avgResult.TPS, avgResult.CTPS)

}
