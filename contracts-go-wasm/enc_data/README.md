ERC721 Token Standard:
https://eips.ethereum.org/EIPS/eip-721

## The description of methods are below:
## 1. InitContract
### args:
#### key1: name(optional)
#### value1: string
#### key2: symbol(optional)
#### value2: string
#### key3: tokenURI(optional)
#### value3: string
#### example:
```json
{"tokenURI":"http://chainmaker.org.cn"}
```

## 2. name
### args: no args
### response example: "erc721"

## 3. symbol
### args: no args
### response example: "erc721X"

## 4. balanceOf
### args:
#### key1: "account"
#### value1: string
#### example:
```json
{"account":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f"}
```
### response example: "0"

## 5. ownerOf
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
### response example: "ec47ae0f0d6a0e952c240383d70ab43b19997a9f"

## 6. mint
### args:
#### key1: "to"
#### value1: string
#### key2: "tokenId"
#### value2: string
#### key3: "metadata"
#### value3: bytes
#### example:
```json
{"to":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f", "tokenId":"111111111111111111111112", "metadata": "url:https://chainmaker.org.cn"}
```
#### resp exampl: "mint success"

## 7. tokenURI
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
#### resp exampl: "http://chainmaker.org.cn/111111111111111111111112"

## 8. tokenMetadata
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
#### resp exampl: "url:http://chainmaker.org.cn/111111111111111111111112"

## 9. tokenLatestTxInfo
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
#### resp exampl: 
```json
{"TxId":"17262429164a0e82ca17c10d4d4bc2b11be6c7c1e9cd4d6db287a8a4f3f2e2e5","BlockHeight":79,"From":"0000000000000000000000000000000000000000","To":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f","Timestamp":"1668060470"}
```

## 10. accountTokens
### args:
#### key1: "account"
#### value1: string
#### example:
```json
{"account":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f"}
```
#### resp exampl:
```json
{"Account":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f","Tokens":["111111111111111111111112","111111111111111111111113"]}
```

## 11. approve
### args:
#### key1: "to"
#### value1: string
#### key2: "tokenId"
#### value2: string
#### example:
```json
{"to":"a04f7895de24f61807a729be230f03da8c0eef42", "tokenId":"111111111111111111111112"}
```
#### resp exampl: "approve success"
### event:
#### topic: approve
#### data: owner, to, tokenId
#### example:
```json
["ec47ae0f0d6a0e952c240383d70ab43b19997a9f","a04f7895de24f61807a729be230f03da8c0eef42","111111111111111111111112"]
```

## 12. getApprove
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
#### resp exampl: "ec47ae0f0d6a0e952c240383d70ab43b19997a9f"

## 13. transferFrom
### args:
#### key1: "from"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "tokenId"
#### value2: string
#### example:
```json
{"from":"ec47ae0f0d6a0e952c240383d70ab43b19997a9f", "to":"a04f7895de24f61807a729be230f03da8c0eef42", "tokenId":"111111111111111111111112"}
```
#### resp exampl: "transfer success"
### event:
#### topic: transfer
#### data: from, to, tokenId
#### example:
```json
["ec47ae0f0d6a0e952c240383d70ab43b19997a9f","a04f7895de24f61807a729be230f03da8c0eef42","111111111111111111111112"]
```

## Test

### 部署合约
```sh
./cmc client contract user create \
--contract-name=erc721 \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/erc721-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{\"name\":\"huanletoken\", \"symbol\":\"hlt\", \"tokenURI\":\"https://chainmaker.org.cn\"}"


./cmc client contract user create \
--contract-name=erc721 \
--runtime-type=DOCKER_GO \
--byte-code-path=./testdata/go-wasm-demo/erc721.7z \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{\"name\":\"huanletoken\", \"symbol\":\"hlt\", \"tokenURI\":\"https://chainmaker.org.cn\"}"


```

### 查询name
#### 验证Case1：
部署合约时如果没有指定erc721的name，默认的name为空，需要验证name为空
#### 验证Case2：
部署合约时如果指定了name参数，在这儿获取时验证name是否和部署合约时指定的name一致
```sh
./cmc client contract user invoke \
--contract-name=erc721 \
--method=name \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}"

{
  "contract_result": {
    "gas_used": 2478706,
    "result": "[name] name=huanletoken"
  },
  "tx_block_height": 4,
  "tx_id": "1834a1fdc579d2c4ca360806b60e8b337430c6761a85463fb2e77cdd94c8cf76",
  "tx_timestamp": 1744197067
}

```

### 查询symbol
#### 验证Case1：
部署合约时如果没有指定erc721的symbol，默认的symbol为空，这儿需要验证symbol为空
#### 验证Case2：
部署合约时如果指定了symbol参数，在这儿获取时验证symbol是否和部署合约时指定的symbol一致
```sh

./cmc client contract user invoke \
--contract-name=erc721 \
--method=symbol \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}"

{
  "contract_result": {
    "gas_used": 2475694,
    "result": "[symbol] symbol=hlt"
  },
  "tx_block_height": 5,
  "tx_id": "1834a20acb22f3cfca69756c336abe6c7985111043bd4f4fa46272803867169b",
  "tx_timestamp": 1744197122
}
```

### 查询tokenURI
#### 验证Case1：
验证返回的tokenURI是否为安装合约时指定的tokenURI+'/'+tokenId
```sh
./cmc client contract user invoke \
--contract-name=erc721 \
--method=tokenURI \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}"

{
  "contract_result": {
    "gas_used": 2501544,
    "result": "https://chainmaker.org.cn/0"
  },
  "tx_block_height": 6,
  "tx_id": "1834a216445afd5ccae9a71a90f8990c86b76f04799f4b96861e887e632061a1",
  "tx_timestamp": 1744197172
}
```

### 查询账户nft数量
#### 验证Case1：
部署合约后所有账户默认的nft数量为0，这儿需要验证账户默认nft数量是否为0
```sh
./cmc client contract user invoke \
--contract-name=erc721 \
--method=balanceOf \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"account\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"

{
  "contract_result": {
    "gas_used": 2465838,
    "result": "0"
  },
  "tx_block_height": 7,
  "tx_id": "1834a22ab28d0563ca959a653718d35034e797dcf3f34b9bb7c1e41c97193d85",
  "tx_timestamp": 1744197259
}
```

### 查询nft所属账户
#### 验证Case1：
部署合约后如果nft不存在，查询nft所属账户会报错，这儿需要验证nft不存在时的错误情况
```sh

./cmc client contract user invoke \
--contract-name=erc721 \
--method=ownerOf \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 2493988,
    "message": "contract message:get owner failed",
    "result": ""
  },
  "tx_block_height": 4,
  "tx_id": "1834a288b1539ac2caaba6485172b2ddb53d2d0eda834acf88aaa373ffbd977e",
  "tx_timestamp": 1744197663
}
```

### 发行nft
#### 验证Case1：
发行nft后需要验证账户nft数量是否增加1
#### 验证Case2：
发行nft后需要验证nft所属账户是否正确
```sh

./cmc client contract user invoke \
--contract-name=erc721 \
--method=mint \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"to\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"tokenId\":\"111111111111111111111112\", \"metadata\":\"url:http://chainmaker.org.cn/\"}"

{
  "contract_result": {
    "gas_used": 1684842,
    "result": "mint success"
  },
  "tx_block_height": 3,
  "tx_id": "18363075bb0796d0cae3c1fc88c43d10a8b5af755f8b47398c58f72edda442ad",
  "tx_timestamp": 1744635187
}
```

### 查询token metadata信息
#### 验证Case1：
这儿验证查询到的metadata是否和mint时传递的一致
```sh

./cmc client contract user invoke \
--contract-name=erc721 \
--method=tokenMetadata \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "gas_used": 2583152,
    "result": "tokenMetadata is url:http://chainmaker.org.cn/"
  },
  "tx_block_height": 5,
  "tx_id": "1834a2f56ce4e6f8ca435325e6c4f200cc2a5534c11c41ae9b83718647d73019",
  "tx_timestamp": 1744198130
}
```

### 查询account tokens信息
#### 验证Case1：
验证账户下是否包含了所有发行的nft
```sh
./cmc client contract user invoke \
--contract-name=erc721 \
--method=accountTokens \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"account\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"

{
  "contract_result": {
    "gas_used": 1156636,
    "result": "{\"Account\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\",\"Tokens\":[\"111111111111111111111112\"]}"
  },
  "tx_block_height": 4,
  "tx_id": "18363364be791f3eca93faa317f01233f0e3cd46651347b2bf4ccbd86134ba46",
  "tx_timestamp": 1744638413
}
```

### 查询token 最近一笔交易信息
#### 验证Case1：
验证token最近一笔的交易信息是否正确
```sh
./cmc client contract user invoke \
--contract-name=erc721 \
--method=tokenLatestTxInfo \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "gas_used": 453783,
    "result": "tokenLatestTxInfo is {\"TxId\":\"1836336446e805cfcaf0f3d4b94c1e2b34963e034af5486c93962681880c85ee\",\"BlockHeight\":3,\"From\":\"0000000000000000000000000000000000000000\",\"To\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
  },
  "tx_block_height": 5,
  "tx_id": "1836337ed9e89ca3ca4f9c8f75d043d64a97c3cb77e54d59a277eb4e8bcfb641",
  "tx_timestamp": 1744638525
}
```

### 获取授权信息
#### 验证Case1：
如果nft没有进行过授权，查询到的授权信息应为空
```sh
./cmc client contract user invoke \
--contract-name=erc721 \
--method=getApproved \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "gas_used": 2498257,
    "result": ""
  },
  "tx_block_height": 7,
  "tx_id": "1834a332d8e4910fca034e9c1e20580ecc922190b7db462fade1d4bf5089881e",
  "tx_timestamp": 1744198394
}
```

### 授权
#### 验证Case1：
授权后需要验证授权信息是否正确
```sh
./cmc client contract user invoke \
--contract-name=erc721 \
--method=setApprovalForAll2 \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"approvalFrom\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc721",
        "contract_version": "1.0",
        "event_data": [
          "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
          "7a83769df9cdfe9c96bf8e01c623e9686a7dc1e796ce12c25ef327d7fd1871ee",
          "1"
        ],
        "topic": "ApprovalForAll2",
        "tx_id": "1834a3cf3acbb6c0cae46184038eedba573c83a1ffd3424faaf3b1d13bdcf3ae"
      }
    ],
    "gas_used": 2598591,
    "result": "setApprovalForAll2 success"
  },
  "tx_block_height": 6,
  "tx_id": "1834a3cf3acbb6c0cae46184038eedba573c83a1ffd3424faaf3b1d13bdcf3ae",
  "tx_timestamp": 1744199066
}

./cmc client contract user invoke --contract-name=erc721 --method=approve --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"111111111111111111111112\"}"
./cmc client contract user invoke \
--contract-name=erc721 \
--method=approve \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc721",
        "contract_version": "1.0",
        "event_data": [
          "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
          "a04f7895de24f61807a729be230f03da8c0eef42",
          "111111111111111111111112"
        ],
        "topic": "approve",
        "tx_id": "1834a3de67d121ccca28712356623bff74bede510827497783d8763a1e5a7995"
      }
    ],
    "gas_used": 2972485,
    "result": "approve success"
  },
  "tx_block_height": 7,
  "tx_id": "1834a3de67d121ccca28712356623bff74bede510827497783d8763a1e5a7995",
  "tx_timestamp": 1744199131
}
```

### 根据授权转账
#### 验证Case1：
转账后需要验证授权信息是否发生了变化
```sh
./cmc client contract user invoke --contract-name=erc721 --method=transferFrom --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml --params="{\"from\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"111111111111111111111112\"}"

./cmc client contract user invoke \
--contract-name=erc721 \
--method=setApprovalForAll2 \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"approvalFrom\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\"}"
{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc721",
        "contract_version": "1.0",
        "event_data": [
          "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
          "b0b9c232e14b79e8ab2208fdccea4e1fc64837db21c1c489caa563d217460e33",
          "1"
        ],
        "topic": "ApprovalForAll2",
        "tx_id": "1836347cbf6e3a34caf041b182a3366063ba599e16cd45ad8fe612c43211e1a6"
      }
    ],
    "gas_used": 440789,
    "result": "setApprovalForAll2 success"
  },
  "tx_block_height": 4,
  "tx_id": "1836347cbf6e3a34caf041b182a3366063ba599e16cd45ad8fe612c43211e1a6",
  "tx_timestamp": 1744639616
}

./cmc client contract user invoke \
--contract-name=erc721 \
--method=transferFrom \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"from\":\"ec47ae0f0d6a0e952c240383d70ab43b19997a9f\", \"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc721",
        "contract_version": "1.0",
        "event_data": [
          "ec47ae0f0d6a0e952c240383d70ab43b19997a9f",
          "a04f7895de24f61807a729be230f03da8c0eef42",
          "111111111111111111111112"
        ],
        "topic": "transfer",
        "tx_id": "1836347e7c98aebdca04e20e1398eee7e8bc9c9e33934d15b36758c65fce0f7b"
      }
    ],
    "gas_used": 2325650,
    "result": "transfer success"
  },
  "tx_block_height": 5,
  "tx_id": "1836347e7c98aebdca04e20e1398eee7e8bc9c9e33934d15b36758c65fce0f7b",
  "tx_timestamp": 1744639623
}
```
