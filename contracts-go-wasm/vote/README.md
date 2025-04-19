This is the vote contract for test net demo.

## The description of methods are below:
## 1. issueProject
### args: 
#### key: "projectInfo"
#### value: json
#### example:
```json
{"projectInfo":{"Id":"projectId1","PicUrl":"www.sina.com","Title":"wonderful","StartTime":"1664450291","EndTime":"1665314291","Desc":"the 1","Items":[{"Id":"item1","PicUrl":"www.baidu.com","Desc":"beautiful", "Url":"www.qq.com"},{"Id":"item2","PicUrl":"www.baidu.com","Desc":"beautiful", "Url":"www.qq.com"}]}}
```
### event:
#### topic: issue project
#### data: same as args

## 2. vote
### args:
#### key1: "projectId"
#### value1: string
#### key2: "itemId"
#### value2: string
#### example:
```json
{"projectId":"projectId1","itemId":"item1"}
```
### event:
#### topic: vote
#### data: projectId, itemId, voter
#### example:
```json
["projectId1","itemId1","441224c1757ec1cc67f4b7b3ac29c76cf5799ee5"]
```

## 3. queryProjectVoters
### args:
#### key1: "projectId"
#### value1: string
#### example:
```json
{"projectId":"projectId1"}
```
#### resp exampl:
```json
{"ProjectId":"projectId1","ItemVotes":[{"ItemId":"item1","VotesCount":1,"Voters":["441224c1757ec1cc67f4b7b3ac29c76cf5799ee5"]},{"ItemId":"item2","VotesCount":1,"Voters":["441224c1757ec1cc67f4b7b3ac29c76cf5799ee5"]}]}
```

## 4. queryProjectItemVoters
### args:
#### key1: "projectId"
#### value1: string
#### key2: "itemId"
#### value2: string
#### example:
```json
{"projectId":"projectId1","itemId":"item1"}
```
#### resp example:
```json
{"ItemId":"item1","VotesCount":1,"Voters":["441224c1757ec1cc67f4b7b3ac29c76cf5799ee5"]}
```

## Test

### 安装合约
```sh

./cmc client contract user create \
--contract-name=vote \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/vote-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"


```

### 发布投票项目
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=issueProject \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"projectInfo\":{\"Id\":\"projectId2\",\"PicUrl\":\"www.sina.com\",\"Title\":\"wonderful\",\"StartTime\":\"1664450291\",\"EndTime\":\"1665314291\",\"Desc\":\"the 1\",\"Items\":[{\"Id\":\"item1\",\"PicUrl\":\"www.baidu.com\",\"Desc\":\"beautiful\", \"Url\":\"www.qq.com\"},{\"Id\":\"item2\",\"PicUrl\":\"www.baidu.com\",\"Desc\":\"beautiful\", \"Url\":\"www.qq.com\"}]}}"

期望输出
{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "vote",
        "contract_version": "1.0",
        "event_data": [
          "{\"Desc\":\"the 1\",\"EndTime\":\"1665314291\",\"Id\":\"projectId2\",\"Items\":[{\"Desc\":\"beautiful\",\"Id\":\"item1\",\"PicUrl\":\"www.baidu.com\",\"Url\":\"www.qq.com\"},{\"Desc\":\"beautiful\",\"Id\":\"item2\",\"PicUrl\":\"www.baidu.com\",\"Url\":\"www.qq.com\"}],\"PicUrl\":\"www.sina.com\",\"StartTime\":\"1664450291\",\"Title\":\"wonderful\"}"
        ],
        "topic": "issue_project",
        "tx_id": "1834a1154d540061ca567b9165537ed65871ebe7fef643088656878eb136569a"
      }
    ],
    "gas_used": 4162116,
    "result": "issue project success"
  },
  "tx_block_height": 3,
  "tx_id": "1834a1154d540061ca567b9165537ed65871ebe7fef643088656878eb136569a",
  "tx_timestamp": 1744196068
}
```

### 投票
### 测试Case1: 
投票的projectId或itemId不存在时应报错
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=vote \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true \
--params="{\"projectId\":\"projectId1\",\"itemId\":\"item1\"}"

期望输出
{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 298600,
    "message": "contract message:get project info from store failed",
    "result": ""
  },
  "tx_block_height": 5,
  "tx_id": "18362e4f47a49e4dca1620ab4088635adab933b1e9784ddfa1e90b6ac18f82cc",
  "tx_timestamp": 1744632823
}

./cmc client contract user invoke \
--contract-name=vote \
--method=vote \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true \
--params="{\"projectId\":\"projectId2\",\"itemId\":\"item1\"}"

期望输出
{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "vote",
        "contract_version": "1.0",
        "event_data": [
          "projectId2",
          "item1",
          "b0b9c232e14b79e8ab2208fdccea4e1fc64837db21c1c489caa563d217460e33"
        ],
        "topic": "vote",
        "tx_id": "18362fb5b0326058ca8aadadc250b71771cbbd3bbc554f1093a4324a8877c728"
      }
    ],
    "gas_used": 2312060,
    "result": "vote success"
  },
  "tx_block_height": 5,
  "tx_id": "18362fb5b0326058ca8aadadc250b71771cbbd3bbc554f1093a4324a8877c728",
  "tx_timestamp": 1744634363
}

```

### 查询项目投票人员
### 测试Case1:
查询投票项目的projectId不存在时应报错
### 测试Case2:
查询投票项目的projectId存在时应返回正确的投票人员集合
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=queryProjectVoters \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true \
--params="{\"projectId\":\"projectId1\"}"

{
  "contract_result": {
    "gas_used": 3351324,
    "result": "{\"ProjectId\":\"projectId1\",\"ItemVotes\":[]}"
  },
  "tx_block_height": 5,
  "tx_id": "1834a122f3d90222ca58acbc5c036556bdd4e8aa48f4452daa52d1396f777d2e",
  "tx_timestamp": 1744196127
}
```

### 查询投票项投票人员
### 测试Case1:
查询投票项目的projectId或itemId不存在时应报错
### 测试Case2:
查询投票项目的projectId和itemId存在时应返回正确的投票人员集合
```sh
./cmc client contract user invoke \
--contract-name=vote \
--method=queryProjectItemVoters \
--sync-result=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true \
--params="{\"projectId\":\"projectId1\",\"itemId\":\"item1\"}"

{
  "contract_result": {
    "gas_used": 3154061,
    "result": "{\"ItemId\":\"item1\",\"VotesCount\":0,\"Voters\":[]}"
  },
  "tx_block_height": 6,
  "tx_id": "1834a12c4456f186cacf98e03e2143728db8ce8e00d74a8e814c00e0602e7dae",
  "tx_timestamp": 1744196167
}

```

