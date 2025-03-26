/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package parallel

import (
	logger "chain-performance-test/log"
	"chain-performance-test/testdata"
	"os"
	"path"
	"path/filepath"
	. "regexp"
	"runtime"
	"strings"
)

// ProcessLog 处理日志
func ProcessLog() {
	logDir := GetCurrentAbPathConfigByCaller()
	fileName := FindLogName(logDir)
	MoveLog(fileName, logDir)
}

// GetCurrentAbPathConfigByCaller 获取当前执行程序所在的上一级目录
func GetCurrentAbPathConfigByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	abPath2 := filepath.Join(abPath, "..")
	return strings.Replace(abPath2, "/", testdata.NewSep, -1)
}

// FindLogName 匹配当前目录下日志文件
func FindLogName(dir string) []string {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var logName []string
	reg := MustCompile(".log\\.(.*?)")
	for _, fileInfo := range fileInfos {
		currentFileName := fileInfo.Name()
		txidStr := reg.FindAllString(currentFileName, -1)
		if len(txidStr) > 0 {
			logName = append(logName, currentFileName)
		}
	}
	return logName
}

// MoveLog 移动日志文件
func MoveLog(fileNames []string, dirPath string) {
	oldLogDir := dirPath + "/log/"
	// 确保 log 文件夹存在，如果不存在则创建
	err := os.MkdirAll(oldLogDir, os.ModePerm)
	if err != nil {
		logger.Logger.Println(err)
	}

	for _, fileName := range fileNames {
		err = os.Rename(dirPath+testdata.NewSep+fileName, oldLogDir+fileName)
		if err != nil {
			logger.Logger.Println(err)
		}
	}
}
