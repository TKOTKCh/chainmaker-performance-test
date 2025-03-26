# 性能测评（长安链）

### [相关文档]()

### 相关接口说明

- 长安链流量抓取

  参数说明
  - host：服务器地址
  - port：服务器端口号
  - username：服务器登陆用户名
  - password：登陆服务器密码
  - catch：服务端抓取流量端口
  
  结果说明
  - 抓取流量是否成功

```shell
    CatchTraffic(host string, port string, username string, password string, 
     trafficPort string) bool
```

- 长安链流量分析

  参数说明
    - path：流量包所在地址

  结果说明
    - 通信源节点IP、端口号，节点ID，目的节点IP、端口号，节点ID，节点间的通信协议和加密方式
```shell
    Analyse(pcapPath string) map[SourceInfo][]DestInfo
```

- 长安链性能压测

  参数说明
    - 无，从配置文件中读取压测任务所需参数

  结果说明
    - 压测结果，包括交易数、成功上链数、TPS、CTPS、QTPS
```shell
    processPressureTest() Result
```

- 长安链配置文件下载

  参数说明
    - host：服务器地址
    - port：服务器端口号
    - username：服务器登陆用户名
    - password：登陆服务器密码
    - configFile：需要下载的配置文件名
    - newSep：新分隔符，将配置文件中原分隔符替换为新分隔符，windows为“\\”，linux为“/”

  结果说明
    - 下载的配置文件在压测引擎中的路径

```shell
    DownConfigAndLoad(host string, port int64, username, password, configFile, newSep string) string
```

- 长安链SDK配置文件下载

  参数说明
    - model：长安链的配置链模式
    - hostList：需要加压的主机IP地址列表
    - grpcList：需要加压的主机列表对应的端口列表

  结果说明
    - 生成的SDK配置文件在压测引擎中的路径

```shell
    GenerateSdkConfig(hostList, grpc1List []string, model string) string
```

- 长安链压测参数配置文件生成

  参数说明
    - model：长安链的配置链模式
    - contractFunction：智能合约函数名称
    - contractName：智能合约名称
    - runTimeType：智能合约语言类型
    - allParams：运行智能合约时所需参数
    - climbTime：线程爬坡时间
    - loopNum：循环次数
    - sleepTime：睡眠时间
    - threadNum：并发线程数

  结果说明
    - 生成的压测参数配置文件在压测引擎中的路径

```shell
    GenerateConstConfig(model, contractFunction, contractName, runTimeType, allParams string,
     climbTime, loopNum, sleepTime, threadNum int) string
```




#### 目录结构
```text
chainmaker-performance-test
├─ LICENSE
├─ README.md
├─ build
│    └─ config    # 长安链的配置文件，不同的链模式此文件结构不同
│           ├─ node1
│           │    ├─ admin
│           │    │    ├─ admin1
│           │    │    │    ├─ admin1.key
│           │    │    │    └─ admin1.pem
│           │    │    ├─ admin2
│           │    │    │    ├─ admin2.key
│           │    │    │    └─ admin2.pem
│           │    │    ├─ admin3
│           │    │    │    ├─ admin3.key
│           │    │    │    └─ admin3.pem
│           │    │    ├─ admin4
│           │    │    │    ├─ admin4.key
│           │    │    │    └─ admin4.pem
│           │    │    └─ admin5
│           │    │           ├─ admin5.key
│           │    │           └─ admin5.pem
│           │    ├─ node1.key
│           │    ├─ node1.nodeid
│           │    ├─ node1.pem
│           │    └─ user
│           │           └─ client1
│           │                  ├─ client1.addr
│           │                  ├─ client1.key
│           │                  └─ client1.pem
│           ├─ node2
│           │    ├─ admin
│           │    │    ├─ admin1
│           │    │    │    ├─ admin1.key
│           │    │    │    └─ admin1.pem
│           │    │    ├─ admin2
│           │    │    │    ├─ admin2.key
│           │    │    │    └─ admin2.pem
│           │    │    ├─ admin3
│           │    │    │    ├─ admin3.key
│           │    │    │    └─ admin3.pem
│           │    │    ├─ admin4
│           │    │    │    ├─ admin4.key
│           │    │    │    └─ admin4.pem
│           │    │    └─ admin5
│           │    │           ├─ admin5.key
│           │    │           └─ admin5.pem
│           │    ├─ node2.key
│           │    ├─ node2.nodeid
│           │    ├─ node2.pem
│           │    └─ user
│           │           └─ client1
│           │                  ├─ client1.addr
│           │                  ├─ client1.key
│           │                  └─ client1.pem
│           ├─ node3
│           │    ├─ admin
│           │    │    ├─ admin1
│           │    │    │    ├─ admin1.key
│           │    │    │    └─ admin1.pem
│           │    │    ├─ admin2
│           │    │    │    ├─ admin2.key
│           │    │    │    └─ admin2.pem
│           │    │    ├─ admin3
│           │    │    │    ├─ admin3.key
│           │    │    │    └─ admin3.pem
│           │    │    ├─ admin4
│           │    │    │    ├─ admin4.key
│           │    │    │    └─ admin4.pem
│           │    │    └─ admin5
│           │    │           ├─ admin5.key
│           │    │           └─ admin5.pem
│           │    ├─ node3.key
│           │    ├─ node3.nodeid
│           │    ├─ node3.pem
│           │    └─ user
│           │           └─ client1
│           │                  ├─ client1.addr
│           │                  ├─ client1.key
│           │                  └─ client1.pem
│           └─ node4
│                  ├─ admin
│                  │    ├─ admin1
│                  │    │    ├─ admin1.key
│                  │    │    └─ admin1.pem
│                  │    ├─ admin2
│                  │    │    ├─ admin2.key
│                  │    │    └─ admin2.pem
│                  │    ├─ admin3
│                  │    │    ├─ admin3.key
│                  │    │    └─ admin3.pem
│                  │    ├─ admin4
│                  │    │    ├─ admin4.key
│                  │    │    └─ admin4.pem
│                  │    └─ admin5
│                  │           ├─ admin5.key
│                  │           └─ admin5.pem
│                  ├─ node4.key
│                  ├─ node4.nodeid
│                  ├─ node4.pem
│                  └─ user
│                         └─ client1
│                                ├─ client1.addr
│                                ├─ client1.key
│                                └─ client1.pem
├─ chainclient
│    └─ client.go
├─ config                                         # 本次压测的配置文件夹
│    ├─ clients.yml
│    └─ const_config.yml
├─ config_example                                 # 示例配置文件夹
│    ├─ bc1.yml                                   # chainmaker配置文件
│    ├─ chainmaker.yml                            # chainmaker配置文件
│    ├─ clients_example.yml                       # 多节点加压示例配置文件
│    ├─ const_config_example.yml                  # 压测参数示例配置文件
│    ├─ sdk_config_ca_example.yml
│    ├─ sdk_config_pk_example.yml                 # pk链模式 长安链连接示例配置文件
│    └─ sdk_config_pwk_example.yml                # pwk链模式 长安链连接示例配置文件
├─ contract
│    ├─ asset_demo
│    │    └─ rust-asset-2.0.0.wasm
│    └─ claim_demo
│           ├─ dockerFact230.7z
│           └─ rustFact.wasm
├─ datahandler
│    ├─ downLoadFile.go                           # 下载并上传长安链配置文件
│    ├─ get_data.go                               # 数据处理部分，负责连接redis，从redis接受任务数据并将结果上传到redis
│    └─ updateConfigYml.go                        # 根据任务数据更新压测的配置文件
├─ example
│    └─ main.go                                  # 前端交互示例引擎代码
├─ go.mod
├─ go.sum
├─ log                                           # 日志存档文件
│    └─ log.go
├─ mock
│    ├─ claim_mock.go                            # 单元测试mock存证任务压测接口
│    ├─ client_mock.go                           # 单元测试mock与长安连连接接口
│    ├─ query_mock.go                            # 单元测试mock查询任务压测接口
│    └─ subservice_mock.go                       # 单元测试mock消息订阅接口
├── parallel
│    ├── Init.go                                 # 根据配置文件更新本次压测参数
│    ├── parallel_claim.go                       # 存证任务压测
│    ├── parallel_query.go                       # 查询任务压测
│    ├── parallel_test.go                        # 单侧相关
│    ├── process_log.go                          # 处理客户端生成的日志文件
│    └─ process_pressure.go
├─ subservice
│    └─ subService.go                            # 消息订阅封装
├─ testdata
│    ├─ config.go
│    └─ test.pcap
└─ traffic
       ├─ analysis_traffic.go                    # 长安链流量分析相关操作封装
       ├─ catch_traffic.go                       # 长安链流量分析相关操作封装
       ├─ config.go                              # 配置文件对应结构定义
       ├─ ssh.go                                 # ssh连接相关操作封装
       └─ traffic_test.go                        # 单测相关
```

