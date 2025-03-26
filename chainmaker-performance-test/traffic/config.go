/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package traffic

type Config struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type FileResult struct {
	Status         bool   `yaml:"status" json:"status"`
	LocalFileName  string `yaml:"local_file_name" json:"local_file_name"`
	LocalFileSize  uint64 `yaml:"local_file_size" json:"local_file_size"`
	RemoteFileName string `yaml:"remote_file_name" json:"remote_file_name"`
	RemoteFileSize uint64 `yaml:"remote_file_size" json:"remote_file_size"`
}
