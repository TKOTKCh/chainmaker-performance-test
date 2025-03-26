/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package testdata

// traffic 相关配置变量
const (
	// 远程服务器主机号
	ServerHost = "127.0.0.1"
	// 远程服务器端口
	ServerPort = "22"
	// 远程服务器用户名
	ServerUser = "root"
	// 远程服务器密码
	ServerPasswd = "123456"
	// 监听流量端口号
	ServerListenPort = "11301"
	// 流量包名
	PcapPath = "test1.pcap"
	// 本地chainmaker.yml配置文件路径
	ConfigPath = "config_example"
)

// datahandler 相关配置变量
const (
	// DownFileHost 下载文件服务器ip
	DownFileHost = "127.0.0.1"
	// DownFilePort 端口
	DownFilePort = 22
	// DownFileUser 用户名
	DownFileUser = "root"
	// DownFilePasswd 密码
	DownFilePasswd = "123456"

	// ConnectRedisAddr 连接redis地址
	ConnectRedisAddr = "127.0.0.1:6379"
	// ConnectRedisPasswd 密码
	ConnectRedisPasswd = ""
	// ReadRedisKey 读取redis的key值
	ReadRedisKey = "cm_stressTest_default[0-9]*"
)

// 通用
const (
	// NewSep 操作系统 window \\ , mac和linux /
	NewSep = "/"
)
