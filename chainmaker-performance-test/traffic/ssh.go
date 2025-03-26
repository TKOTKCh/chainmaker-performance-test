/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package traffic

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_ssh/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

// SSH SSH连接对象
type SSH struct {
	Client *ssh.Client // ssh客户端
	Config *Config
}

func NewWithConfig(config *Config) *SSH {
	s := &SSH{}

	// 配置
	s.Config = config

	// 返回
	return s
}

// Connect 创建连接
func (s *SSH) Connect() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		err          error
	)

	// 获取权限方法
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(s.Config.Password))
	clientConfig = &ssh.ClientConfig{
		User:            s.Config.Username,
		Auth:            auth,
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接SSH
	addr = fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
	if s.Client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}

	// 返回
	return nil
}

// Sudo 执行sudo命令
func (s *SSH) Sudo(command string) (string, error) {

	// 创建连接
	if s.Client == nil {
		err := s.Connect()
		if err != nil {
			return "", err
		}
	}

	// 创建session会话
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// 执行命令
	command = fmt.Sprintf("echo %s | sudo -S %s", s.Config.Password, command)
	buf, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	// 返回指令执行结果
	return string(buf), nil
}

// GetSshConfig 获取SSH链接配置
func (s *SSH) GetSshConfig() *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: s.Config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         10 * time.Second,
	}
	return sshConfig
}

// DownloadFile 下载文件
func (s *SSH) DownloadFile(remoteFileName, localFileName string) FileResult {
	result := FileResult{
		LocalFileName:  localFileName,
		RemoteFileName: remoteFileName,
	}

	// 创建客户端
	sshConfig := s.GetSshConfig()

	// 建立一个SSH服务器连接
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), sshConfig)
	if err != nil {
		return result
	}
	defer sshClient.Close()

	// 获取SFTP客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return result
	}
	defer sftpClient.Close()

	//下载文件
	remoteFile, err := sftpClient.Open(remoteFileName)
	if err != nil {
		return result
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localFileName)
	if err != nil {
		return result
	}
	defer localFile.Close()

	// 将远程文件复制到本地文件流
	n, err := io.Copy(localFile, remoteFile)
	if err != nil {
		return result
	}

	// 获取远程文件大小
	remoteFileInfo, err := sftpClient.Stat(remoteFileName)
	if err != nil {
		return result
	}

	// 下载结果
	result.Status = true
	result.LocalFileSize = uint64(n)
	result.RemoteFileSize = uint64(remoteFileInfo.Size())

	return result
}
