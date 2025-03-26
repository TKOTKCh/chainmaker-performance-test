/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package traffic

import (
	logger "chain-performance-test/log"
	"chain-performance-test/testdata"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"strconv"
)

// Connection 创建连接
type Connection struct {
	*ssh.Client
}

// Exists 确定文件是否存在
func Exists(path string) bool {
	// Get file information
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// CatchTraffic 实时捕获流量
func CatchTraffic(host string, port string, username string, password string, trafficPort string) string {
	// 创建远程连接
	num, err := strconv.Atoi(port)
	if err != nil {
		logger.Logger.Println("字符串转换为整数失败:", err)
		return ""
	}
	s := NewWithConfig(&Config{
		Host:     host,
		Port:     num,
		Username: username,
		Password: password,
	})
	var command string
	// 向远程服务器发送命令以获取流量
	if trafficPort == "" {
		command = "sudo tcpdump -i any -c 200 -n -w ./test1.pcap"
	} else {
		command = fmt.Sprintf("sudo tcpdump -i any -c 200 -n src port %s -w %s", trafficPort, testdata.PcapPath)
	}
	output, err := s.Sudo(command)
	if err != nil {
		logger.Logger.Printf("远程执行sudo命令失败 %s: %v\n", output, err)
		return ""
	}
	// 将文件从服务器下载到客户端
	s.DownloadFile(testdata.PcapPath, testdata.PcapPath)

	// 删除服务器上的流量包
	output, err = s.Sudo(fmt.Sprintf("sudo rm %s", testdata.PcapPath))
	if err != nil {
		logger.Logger.Printf("删除服务器流量包失败 %s :%v\n", output, err)
		return ""
	}
	return testdata.PcapPath
}
