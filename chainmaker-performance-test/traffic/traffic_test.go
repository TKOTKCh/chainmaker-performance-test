/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package traffic

import (
	logger "chain-performance-test/log"
	"golang.org/x/crypto/ssh"
	"reflect"
	"strings"
	"testing"
)

func TestBeginLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestBeginLog"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.BeginLog()
		})
	}
}

func TestAnalyse(t *testing.T) {
	type args struct {
		pcapPath string
	}
	tests := []struct {
		name string
		args args
		want map[SourceInfo][]DestInfo
	}{
		{
			name: "Analyse success",
			args: args{
				pcapPath: "../testdata/test.pcap",
			},
			want: map[SourceInfo][]DestInfo{
				SourceInfo{
					Ip:     "127.0.0.1",
					Port:   "11301",
					Nodeid: "QmNfJGVrR63dCucjNxYfTausxGiWuDJpaWsMksRa3bp2RK",
				}: {
					DestInfo{
						Ip:            "127.0.0.1",
						Port:          "11302",
						Protocol:      "TCP",
						Cryptoversion: "TLS 1.2",
						Nodeid:        "QmRmxah2L1jEN6XB38XmTbvQdfEZ2srG5GXJgGA9vLu5bc",
					},
					DestInfo{
						Ip:            "127.0.0.1",
						Port:          "11304",
						Protocol:      "TCP",
						Cryptoversion: "TLS 1.2",
						Nodeid:        "QmaWoafqQgMzC9NXrsYmk6iLzg2nwUEfekG9NiuSr5niDW",
					},
					DestInfo{
						Ip:            "127.0.0.1",
						Port:          "11303",
						Protocol:      "TCP",
						Cryptoversion: "TLS 1.2",
						Nodeid:        "QmcVot7DR2hMadzbUEuo5j8vZo4fxUqgGZSb5Q6CjYm6G8",
					},
				},
			},
		},
		{
			name: "File not exist",
			args: args{
				pcapPath: "../testdata/test2.pcap",
			},
			want: map[SourceInfo][]DestInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Analyse(tt.args.pcapPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Analyse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resolveConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Resolve success",
			args: args{
				path: "../config_example",
			},
			want: map[string]string{
				"11301": "QmNfJGVrR63dCucjNxYfTausxGiWuDJpaWsMksRa3bp2RK",
				"11302": "QmRmxah2L1jEN6XB38XmTbvQdfEZ2srG5GXJgGA9vLu5bc",
				"11303": "QmcVot7DR2hMadzbUEuo5j8vZo4fxUqgGZSb5Q6CjYm6G8",
				"11304": "QmaWoafqQgMzC9NXrsYmk6iLzg2nwUEfekG9NiuSr5niDW",
			},
		},
		{
			name: "File not exist",
			args: args{
				path: "../build/config/node5",
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveConfig(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "File exists",
			args: args{
				path: "../go.mod",
			},
			want: true,
		},
		{
			name: "File not exists",
			args: args{
				path: "../not_main.go",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.args.path); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSSH_Connect(t *testing.T) {
	type fields struct {
		Client *ssh.Client
		Config *Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "connect fail",
			fields: fields{
				Client: nil,
				Config: &Config{
					Host:     "10.112.231.102",
					Port:     22,
					Username: "error",
					Password: "error",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SSH{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			if err := s.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSSH_DownloadFile(t *testing.T) {
	type fields struct {
		Client *ssh.Client
		Config *Config
	}
	type args struct {
		remoteFileName string
		localFileName  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   FileResult
	}{
		{
			name: "file not exist",
			fields: fields{
				Client: nil,
				Config: &Config{
					Host:     "10.112.231.102",
					Port:     22,
					Username: "stack",
					Password: "StackStack123",
				},
			},
			args: args{
				remoteFileName: "test111.pcap",
				localFileName:  "test111.pcap",
			},
			want: FileResult{
				Status:         false,
				LocalFileName:  "test111.pcap",
				LocalFileSize:  0,
				RemoteFileName: "test111.pcap",
				RemoteFileSize: 0,
			},
		},
		{
			name: "file not exist",
			fields: fields{
				Client: nil,
				Config: &Config{
					Host:     "10.112.231.102",
					Port:     22,
					Username: "error",
					Password: "error",
				},
			},
			args: args{
				remoteFileName: "test111.pcap",
				localFileName:  "test111.pcap",
			},
			want: FileResult{
				Status:         false,
				LocalFileName:  "test111.pcap",
				LocalFileSize:  0,
				RemoteFileName: "test111.pcap",
				RemoteFileSize: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SSH{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			if got := s.DownloadFile(tt.args.remoteFileName, tt.args.localFileName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DownloadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSSH_Sudo(t *testing.T) {
	type fields struct {
		Client *ssh.Client
		Config *Config
	}
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "connect fail",
			fields: fields{
				Client: nil,
				Config: &Config{
					Host:     "10.112.231.102",
					Port:     22,
					Username: "error",
					Password: "error",
				},
			},
			args: args{
				command: "sudo echo OK",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SSH{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := s.Sudo(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sudo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("Sudo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCatchTraffic(t *testing.T) {
	type args struct {
		host        string
		port        string
		username    string
		password    string
		trafficPort string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "server error",
			args: args{
				host:        "10.112.211.11",
				port:        "22",
				username:    "error",
				password:    "error",
				trafficPort: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CatchTraffic(tt.args.host, tt.args.port, tt.args.username, tt.args.password, tt.args.trafficPort)
		})
	}
}
