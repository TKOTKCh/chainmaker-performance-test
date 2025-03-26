/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"fmt"
	"chain-performance-test/datahandler"
	logger "chain-performance-test/log"
	"chain-performance-test/parallel"
	"chain-performance-test/testdata"
	"time"
)

// 示例代码
func main() {
	fmt.Println("import finish")
	// 移动上次生成的日志文件,可能会覆盖，请注意保存log下日志
	parallel.ProcessLog()
	// 创建一个日志文件
	logger.BeginLog()
	rdb := datahandler.ConnectRedis(testdata.ConnectRedisAddr, testdata.ConnectRedisPasswd)
	// 持续监听 redis中 key并处理
	for {
		func() {
			getKey := datahandler.GetKey(rdb, testdata.ReadRedisKey)
			if getKey != "" {
				// 捕获本次任务的错误
				defer func() {
					if err := recover(); err != nil {
						logger.Logger.Println("Error:", err)
						// 结果全为0推回 redis，删除错误的 key
						datahandler.PushRedis(getKey, rdb, parallel.TrafficInfo)
						datahandler.DeleteRedisCache(rdb, getKey)
					}
				}()
				// 初始化结果值
				datahandler.TestResult = datahandler.Result{}
				datahandler.UpdateConfigNew(rdb, getKey)

				// 处理敏捷测试任务
				parallel.ProcessPressureTest()

				// 结果推回 redis，删除该 key
				datahandler.PushRedis(getKey, rdb, parallel.TrafficInfo)
				datahandler.DeleteRedisCache(rdb, getKey)
			}
		}()
		time.Sleep(20 * time.Second)
	}
}
