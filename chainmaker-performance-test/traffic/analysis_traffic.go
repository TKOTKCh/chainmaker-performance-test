/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package traffic

import (
	logger "chain-performance-test/log"
	"chain-performance-test/testdata"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SourceInfo 节点源信息
type SourceInfo struct {
	Ip     string `json:"Ip"`
	Port   string `json:"Port"`
	Nodeid string `json:"Nodeid"`
}

// DestInfo 目标节点信息
type DestInfo struct {
	Ip            string `json:"Ip"`
	Port          string `json:"Port"`
	Protocol      string `json:"Protocol"`
	Cryptoversion string `json:"Cryptoversion"`
	Nodeid        string `json:"Nodeid"`
}

// ResolveConfig 解析配置文件
func ResolveConfig(path string) map[string]string {
	// 读取指定路径的配置文件
	config := viper.New()
	config.AddConfigPath(path)
	config.SetConfigName("chainmaker")
	config.SetConfigType("yml")
	if err := config.ReadInConfig(); err != nil {
		return map[string]string{}
	}
	re := regexp.MustCompile(`/(\d+)/`)
	// 读取YML文件中的指定字段信息
	messages := map[string]string{}
	nodes := config.GetStringSlice("net.seeds")
	for _, node := range nodes {
		info := strings.Split(node, "/")
		messages[re.FindStringSubmatch(node)[1]] = info[len(info)-1]
	}

	// 返回
	return messages
}

// Analyse 从给定的 pcap 文件中读取数据包并返回包含流量信息map
func Analyse(pcapPath string) map[SourceInfo][]DestInfo {
	// 创建一个存储流量信息的map
	trafficInfo := make(map[SourceInfo][]DestInfo)

	// 打开pcap文件
	handle, err := pcap.OpenOffline(pcapPath)
	if err != nil {
		return map[SourceInfo][]DestInfo{}
	}

	// 创建新的数据包源以从PCAP文件中读取数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// 通过数据包源循环读取数据包
	for packet := range packetSource.Packets() {
		// 创建源和目标信息的实例
		var source SourceInfo
		var dest DestInfo

		// 解析 TLS 协议
		if packet.ApplicationLayer() != nil {
			var tls layers.TLS
			var decoded []gopacket.LayerType

			// 创建一个解码器，从数据包中解码 TLS 层并存储解码结果
			parser := gopacket.NewDecodingLayerParser(layers.LayerTypeTLS, &tls)
			err := parser.DecodeLayers(packet.ApplicationLayer().LayerContents(), &decoded)
			if err != nil {
				continue
			}

			// 遍历解码结果并提取加密版本信息
			for _, layerType := range decoded {
				switch layerType {
				case layers.LayerTypeTLS:

					dest.Cryptoversion = tls.AppData[0].Version.String()
				}
			}
		}

		// 如果未解析加密版本信息，则跳过当前数据包
		if len(dest.Cryptoversion) == 0 {
			continue
		}

		// 解析IP层
		ip4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ip4Layer != nil {
			// 对于 IPv4 数据包，将提取源和目标 IP 地址
			if ipv4, ok := ip4Layer.(*layers.IPv4); ok {
				// 解码相关信息、源和目标 IP 地址
				source.Ip = ipv4.SrcIP.String()
				dest.Ip = ipv4.DstIP.String()
				dest.Protocol = ipv4.Protocol.String()
			}
		}

		// 解析传输层
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			// 提取 TCP 数据包的源端口和目标端口
			if tcp, ok := tcpLayer.(*layers.TCP); ok {
				source.Port = tcp.SrcPort.String()
				dest.Port = tcp.DstPort.String()
			}

			// 提取 UDP 数据包的源端口和目标端口
			if udp, ok := tcpLayer.(*layers.UDP); ok {
				source.Port = udp.SrcPort.String()
				dest.Port = udp.DstPort.String()
			}
		}

		// 获取指定配置文件中端口对应的节点 ID
		currentDir, err := os.Getwd()
		if err != nil {
			logger.Logger.Println(err)
		}
		//获取项目的根目录
		projectRoot := filepath.Dir(currentDir)
		file, err := findDir(projectRoot, testdata.ConfigPath)
		if err != nil {
			logger.Logger.Println(err)
		}
		result := ResolveConfig(file)
		dest.Nodeid = result[dest.Port]
		source.Nodeid = result[source.Port]
		isRecord := false

		// 确定消息是否被记录
		for index := range trafficInfo[source] {
			if (trafficInfo[source][index].Ip == dest.Ip) && (trafficInfo[source][index].Port == dest.Port) {
				isRecord = true
			}
		}

		// 如果未记录消息，则记录消息
		if !isRecord {
			trafficInfo[source] = append(trafficInfo[source], dest)
		}
	}

	// 确定文件是否存在
	isExist := Exists(testdata.PcapPath)
	if isExist {
		if err = os.Remove("./test1.pcap"); err != nil {
			logger.Logger.Println("test1.pcap流量包不存在", zap.Error(err))
		}
	}

	// 返回
	return trafficInfo
}

func findDir(rootDir, targetDir string) (string, error) {
	var result string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == targetDir {
			// 找到目标目录，将其绝对路径保存到 result 变量中
			result = path
			return filepath.SkipDir // 停止进一步遍历当前目录
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if result == "" {
		return "", fmt.Errorf("target directory %s not found in %s", targetDir, rootDir)
	}

	return result, nil
}
