/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package datahandler

import (
	"archive/zip"
	logger "chain-performance-test/log"
	"chain-performance-test/testdata"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	//连接后端服务器配置
	host     string = testdata.DownFileHost
	port     int64  = testdata.DownFilePort
	username string = testdata.DownFileUser
	Password string = testdata.DownFilePasswd
	Network  string = "tcp"

	//文件夹名
	buildFileName       string = "build"  //下载文件所在文件夹名
	buildConfigFileName string = "config" //解压文件所在文件夹名

	//chainmaker.yml文件名
	chainmaker string = "chainmaker"
	dir        string = "config_example"
)

// 建立连接结构体
type ClientConfig struct {
	Host       string       //ip
	Port       int64        // 端口
	Username   string       //用户名
	Password   string       //密码
	SshClient  *ssh.Client  //ssh client
	SftpClient *sftp.Client //sftp client
	LastResult string       //最近一次运行的结果
}

// GetCurrentBuildPathByCaller 获取工作目录路径
func GetCurrentBuildPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	abPath3 := filepath.Join(abPath, "../"+buildFileName)
	return strings.Replace(abPath3, OldSep, NewSep, -1)
}

// DelOldBuild 删除旧的Build文件夹下配置文件
func DelOldBuild(configPath string) {
	dir, _ := os.ReadDir(configPath)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{buildFileName, d.Name()}...))
	}
}

// DownLoadRun 下载长安链 config文件并放到 build目录下
func DownLoadRun() {
	// 创建链接
	cliConf := new(ClientConfig)
	cliConf.CreateClient(host, port, username, Password)
	// 获取build目录
	buildPath := GetCurrentBuildPathByCaller()
	// 删除build目录下旧文件
	DelOldBuild(buildPath)
	// 下载新文件
	cliConf.Download(ParametersList.ChainInformationList.ConfigFile, buildPath, NewSep+buildConfigFileName+".zip")
	// 更改解压后文件名
	UpdateBuildFileName(NewSep)
}

// DownConfigAndLoad 下载长安链 config文件并放到 build目录下
func DownConfigAndLoad(host string, port int64, username, password, configFile string) string {
	// 创建链接
	cliConf := new(ClientConfig)
	cliConf.CreateClient(host, port, username, password)
	// 获取build目录
	buildPath := GetCurrentBuildPathByCaller()
	logger.Logger.Println(buildPath)
	// 删除build目录下旧文件
	DelOldBuild(buildPath)
	osKind := runtime.GOOS
	var newSep string
	if osKind == "windows" {
		newSep = "\\"
	}
	if osKind == "linux" {
		newSep = "/"
	}
	// 下载新文件
	cliConf.Download(configFile, buildPath, newSep+buildConfigFileName+".zip")
	// 更改解压后文件名
	UpdateBuildFileName(newSep)
	return buildPath + newSep + "config"
}

// DownTrafficYml 下载长安链 chainmaker.yml文件并放到 config_example目录下
func DownTrafficYml(host string, port int64, username, password, filePath string) string {
	// 创建链接
	cliConf := new(ClientConfig)
	cliConf.CreateClient(host, port, username, password)
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	logger.Logger.Println(filename)
	if ok {
		abPath = path.Dir(filename)
	}
	osKind := runtime.GOOS
	var sep string
	if osKind == "windows" {
		sep = "\\"
	}
	if osKind == "linux" {
		sep = "/"
	}
	downFilePath := filepath.Join(abPath, "../"+dir+sep+chainmaker+".yml")
	logger.Logger.Println(downFilePath)
	// 下载新文件
	srcFile, err := cliConf.SftpClient.Open(filePath) //远程
	logger.Logger.Println(srcFile)
	if err != nil {
		logger.Logger.Println(err)
	}
	logger.Logger.Println(downFilePath)
	dstFile, _ := os.Create(downFilePath) //本地
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		logger.Logger.Println("error occurred", err)
	}
	logger.Logger.Println("YML文件下载完毕")

	return downFilePath
}

// CreateClient 创建客户端链接
func (cliConf *ClientConfig) CreateClient(host string, port int64, username, password string) {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		err        error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username
	cliConf.Password = password

	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	if sshClient, err = ssh.Dial(Network, addr, &config); err != nil {
		logger.Logger.Println("error occurred:", err)
	}
	cliConf.SshClient = sshClient

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		logger.Logger.Println("error occurred:", err)
	}
	cliConf.SftpClient = sftpClient
}

// Download 从服务器下载文件
func (cliConf *ClientConfig) Download(srcPath, dstPath, dstFileName string) {
	zipPath := dstPath + dstFileName
	srcFile, err := cliConf.SftpClient.Open(srcPath) //远程
	if err != nil {
		logger.Logger.Panic("从服务器下载文件错误：", err)
	}

	dstFile, _ := os.Create(zipPath) //本地
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		logger.Logger.Println("error occurred", err)
	}
	logger.Logger.Println("从服务器下载文件完毕")
	err = Unzip(zipPath, dstPath) //解压
	if err != nil {
		logger.Logger.Panic("解压发生错误", err)
	}
}

// Unzip 找到zip路径然后打开文件
func Unzip(zipPath, dstDir string) error {
	// open zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	// 将打开文件解压到build目录下
	for _, file := range reader.File {
		if err = UnzipFile(file, dstDir); err != nil {
			return err
		}
	}
	return nil
}

// UnzipFile 把打开的文件解压到对应路径
func UnzipFile(file *zip.File, dstDir string) error {
	// create the directory of file
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// open the file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// save the decompressed file content
	_, err = io.Copy(w, rc)
	return err
}

// UpdateBuildFileName 更新配置文件目录名称
func UpdateBuildFileName(newSep string) {
	buildPath := GetCurrentBuildPathByCaller()
	//获取文件或目录相关信息
	fileInfoList, err := os.ReadDir(buildPath)
	if err != nil {
		logger.Logger.Println(err)
	}
	for i := range fileInfoList { //打印当前文件或目录下的文件或目录名
		// 判断是否为目录，是则更改目录名
		if fileInfoList[i].IsDir() {
			oldPath := buildPath + newSep + fileInfoList[i].Name()
			newPath := buildPath + newSep + "config"
			err = os.Rename(oldPath, newPath)
			if err != nil {
				logger.Logger.Println(err)
			} // 改名
			fileInfoListOfConfig, err1 := os.ReadDir(newPath)
			if err1 != nil {
				logger.Logger.Println(err)
			}
			// 修改该目录下的子级目录名
			for j := range fileInfoListOfConfig {
				if fileInfoListOfConfig[j].IsDir() {
					oldPathOfConfig := newPath + newSep + fileInfoListOfConfig[j].Name()
					newPathOfConfig := newPath + newSep + "node" + strconv.Itoa(j+1)
					err = os.Rename(oldPathOfConfig, newPathOfConfig)
					if err != nil {
						logger.Logger.Println(err)
					}
				}
			}
		}
	}
}
