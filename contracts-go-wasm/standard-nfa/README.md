https://git.chainmaker.org.cn/contracts/standard/-/blob/master/living/CM-CS-221221-NFA.md

## The description of methods are below:
## InitContract
### args:
#### key1: categoryName(optional)
#### value1: string
#### key2: categoryURI(optional)
#### value2: string
#### key3: admin(optional)
#### value3: string
#### example:
```json
{
  "categoryName": "chainmaker",
  "categoryURI": "http://chainmaker.org.cn",
  "admin": "huanle"
}
```

## UpgradeContract
### args:
#### key1: categoryName(optional)
#### value1: string
#### key2: categoryURI(optional)
#### value2: string
#### key3: admin(optional)
#### value3: string
#### example:
```json
{
  "categoryName": "chainmaker",
  "categoryURI": "http://chainmaker.org.cn",
  "admin": "huanle"
}
```

## Standards
### response example: {"standards":["CMNFA"]}

## Mint
### args:
#### key1: "to"
#### value1: string
#### key2: "tokenId"
#### value2: string
#### key3: "categoryName"
#### value2: string
#### key4: "metadata"(optional)
#### value3: bytes
#### example:
```json
{
  "to":"b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "tokenId":"111111111111111111111112",
  "categoryName":"huanle",
  "metadata": "url:https://chainmaker.org.cn"
}
```
#### resp exampl: "Mint success"
### event:
#### topic: Mint
#### data: ZeroAddr, to, tokenId, categoryName, metadata
#### example:
```json
[
  "00000000000000000000",
  "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "111111111111111111111112",
  "huanle",
  "url:https://chainmaker.org.cn"
]
```

## MintBatch
### args:
#### key1: "tokens"
#### value1: json
#### example:
```json
{
  "tokens": [{
    "tokenId": "xxxxx",
    "categoryName": "111",
    "to": "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
    "metadata": "aaa"
  }, {
    "tokenId": "xxxxx1",
    "categoryName": "111",
    "to": "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
    "metadata": "aaa"
  }]
}
```
#### resp exampl: "MintBatch success"
### event:
#### topic: Mint
#### data: tokens
#### example:
```json
[
  "00000000000000000000",
  "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "xxxxx",
  "111",
  "url:https://chainmaker.org.cn"
]
```

## SetApproval
### args:
#### key1: "owner"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "tokenId"
#### value3: string
#### key4: "isApproval"
#### value4: string
#### example:
```json
{
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "tokenId":"111111111111111111111112",
  "isApproval": "true"
}
```
#### resp exampl: "SetApproval success"
### event:
#### topic: SetApproval
#### data: owner, to, tokenId, isApproval
#### example:
```json
[
  "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "111111111111111111111112",
  "true"
]
```

## SetApprovalForAll
### args:
#### key1: "owner"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "isApproval"
#### value3: string
#### example:
```json
{
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "isApproval": "true"
}
```
#### resp exampl: "SetApprovalForAll success"
### event:
#### topic: SetApprovalForAll
#### data: owner, to, isApproval
#### example:
```json
[
  "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "true"
]
```

## TransferFrom
### args:
#### key1: "from"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "tokenId"
#### value3: string
#### example:
```json
{
  "from":"b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "tokenId": "111111111111111111111112"
}
```
#### resp exampl: "TransferFrom success"
### event:
#### topic: TransferFrom
#### data: from, to, tokenId
#### example:
```json
[
  "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "111111111111111111111112"
]
```

## TransferFromBatch
### args:
#### key1: "from"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "tokenIds"
#### value3: array of string
#### example:
```json
{
  "from":"b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "tokenIds": [
    "111111111111111111111111",
    "111111111111111111111112"
  ]
}
```
#### resp exampl: "TransferFromBatch success"
### event:
#### topic: TransferFrom
#### data: from, to, tokenId
#### example:
```json
[
  "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "111111111111111111111112"
]
```

## OwnerOf
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
### response example: "b0831c4aebde4eb5a97b0eb5b1310a746763e752"

## TokenURI
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
#### resp exampl: "http://chainmaker.org.cn/111111111111111111111112"

## SetApprovalByCategory
### args:
#### key1: "owner"
#### value1: string
#### key2: "to"
#### value2: string
#### key3: "categoryName"
#### value3: string
#### key4: "isApproval"
#### value4: string
#### example:
```json
{
  "to":"a04f7895de24f61807a729be230f03da8c0eef42", 
  "categoryName": "1111",
  "isApproval": "true"
}
```
#### resp exampl: "SetApprovalByCategory success"
### event:
#### topic: SetApprovalByCategory
#### data: owner, to, categoryName, isApproval
#### example:
```json
[
  "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
  "a04f7895de24f61807a729be230f03da8c0eef42",
  "1111",
  "true"
]
```

## CreateOrSetCategory
### args:
#### key1: "category"
#### value1: string
#### example:
```json
{
  "category": {
    "categoryName":"1111",
    "categoryURI":"http://www.chainmaker.org.cn/"
  }
}
```
#### resp exampl: "CreateOrSetCategory success"
### event:
#### topic: CreateOrSetCategory
#### data: categoryName, categoryURI
#### example:
```json
[
  "1111",
  "http://www.chainmaker.org.cn/"
]
```

## Burn
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{
  "tokenId": "111111111111111111111112"
}
```
#### resp exampl: "Burn success"
### event:
#### topic: Burn
#### data: from, to, tokenId
#### example:
```json
[
  "111111111111111111111112"
]
```

## GetCategoryByName
### args:
#### key1: "categoryName"
#### value1: string
#### example:
```json
{
  "categoryName": "1111"
}
```
#### resp exampl:
```json
{
"categoryName":"1111",
"categoryURI":"http://www.chainmaker.org.cn/"
}
```

## GetCategoryByTokenId
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{
  "tokenId": "111111111111111111111112"
}
```
#### resp exampl:
```json
{
"categoryName":"1111",
"categoryURI":"http://www.chainmaker.org.cn/"
}
```

## TotalSupply
### args:
#### resp exampl: "0"

## TotalSupplyOfCategory
### args:
#### key1: "categoryName"
#### value1: string
#### example:
```json
{
  "categoryName": "111111111111111111111112"
}
```
#### resp exampl: "0"

## BalanceOf
### args:
#### key1: "account"
#### value1: string
#### example:
```json
{"account":"b0831c4aebde4eb5a97b0eb5b1310a746763e752"}
```
### response example: "0"

## AccountTokens
### args:
#### key1: "account"
#### value1: string
#### example:
```json
{"account":"b0831c4aebde4eb5a97b0eb5b1310a746763e752"}
```
#### resp exampl:
```json
{"account":"b0831c4aebde4eb5a97b0eb5b1310a746763e752","tokens":["111111111111111111111112","111111111111111111111113"]}
```

## TokenMetadata
### args:
#### key1: "tokenId"
#### value1: string
#### example:
```json
{"tokenId":"111111111111111111111112"}
```
#### resp exampl: "url:http://chainmaker.org.cn/111111111111111111111112"

## Test

### 部署合约
```sh

./cmc client contract user create \
--contract-name=CMNFA \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/standard_nfa-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{\"categoryName\":\"1111\", \"categoryURI\":\"chainmaker.org.cn/\", \"admin\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}"

{
  "contract_result": {
    "gas_used": 434515,
    "result": {
      "address": "dc2e430c7b155738dc0bf07c69aa2a9f55f81abf",
      "creator": {
        "member_id": "client1.sign.wx-org1.chainmaker.org",
        "member_info": "Y9PdAUogrWJ0Z1LsrlgbL7ltvZIi1+MNRPVgholFo9c=",
        "member_type": 1,
        "org_id": "wx-org1.chainmaker.org",
        "role": "CLIENT",
        "uid": "b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674"
      },
      "name": "CMNFA",
      "runtime_type": 2,
      "version": "1.0"
    }
  },
  "tx_block_height": 2,
  "tx_id": "1837c2e5d4ce3acaca4345c98a79cc9496492c9532e749a6898ae781869605ea",
  "tx_timestamp": 1745077673
}
```

### 升级合约
```sh
./cmc client contract user upgrade \
--contract-name=CMNFA \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/standard_nfa-go.wasm \
--version=2.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{\"categoryName\":\"1111\", \"categoryURI\":\"chainmaker.org.cn/\", \"admin\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}"


```

### 发行NFA
#### 验证Case1：
发行后可以使用BalanceOf、OwnerOf以及AccountTokens进行验证是否正确发行
```sh

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=Mint \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"to\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\", \"tokenId\":\"111111111111111111111112\", \"categoryName\":\"1111\", \"metadata\":\"aaa\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "0000000000000000000000000000000000000000",
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "111111111111111111111112",
          "1111",
          "aaa"
        ],
        "topic": "Mint",
        "tx_id": "1837c3d1c19e4f20ca0680ac619de3670f54a432bc24476388bc9bb295bf957a"
      }
    ],
    "gas_used": 1445812,
    "result": "Mint success"
  },
  "tx_block_height": 3,
  "tx_id": "1837c3d1c19e4f20ca0680ac619de3670f54a432bc24476388bc9bb295bf957a",
  "tx_timestamp": 1745078686
}
```

### 批量发行NFA
#### 验证Case1：
发行后可以使用BalanceOf、OwnerOf以及AccountTokens进行验证是否正确发行
```sh
./cmc client contract user invoke \
--contract-name=CMNFA \
--method=MintBatch \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{
  \"tokens\": [{
    \"tokenId\": \"xxxxx\",
    \"categoryName\": \"1111\",
    \"to\": \"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",
    \"metadata\": \"YWFh\"
  }, {
    \"tokenId\": \"xxxxx1\",
    \"categoryName\": \"1111\",
    \"to\": \"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",
    \"metadata\": \"YWFh\"
  }]
}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "0000000000000000000000000000000000000000",
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "xxxxx",
          "1111",
          "aaa"
        ],
        "topic": "Mint",
        "tx_id": "1837c3e4857a3e14ca79141b874d29ebb6e5d78f5eac4ed9b591ff3fcea79fd9"
      },
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "0000000000000000000000000000000000000000",
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "xxxxx1",
          "1111",
          "aaa"
        ],
        "topic": "Mint",
        "tx_id": "1837c3e4857a3e14ca79141b874d29ebb6e5d78f5eac4ed9b591ff3fcea79fd9"
      }
    ],
    "gas_used": 3663530,
    "result": "MintBatch success"
  },
  "tx_block_height": 4,
  "tx_id": "1837c3e4857a3e14ca79141b874d29ebb6e5d78f5eac4ed9b591ff3fcea79fd9",
  "tx_timestamp": 1745078766
}
```

### 授权
#### 验证Case1：
授权后需要验证授权信息是否正确
```sh

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=SetApproval \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"111111111111111111111112\", \"categoryName\":\"1111\", \"metadata\":\"aaa\",\"isApproval\":\"true\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "a04f7895de24f61807a729be230f03da8c0eef42",
          "111111111111111111111112",
          "true"
        ],
        "topic": "SetApproval",
        "tx_id": "1837c407e1a67105ca5fac137fd9dea49884e4d7dbd447ce8b20007c4263e968"
      }
    ],
    "gas_used": 611928,
    "result": "SetApproval success"
  },
  "tx_block_height": 6,
  "tx_id": "1837c407e1a67105ca5fac137fd9dea49884e4d7dbd447ce8b20007c4263e968",
  "tx_timestamp": 1745078918
}
```

### 全部授权
#### 验证Case1：
授权后需要验证授权信息是否正确
```sh

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=SetApprovalForAll \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"isApproval\":\"true\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "a04f7895de24f61807a729be230f03da8c0eef42",
          "true"
        ],
        "topic": "SetApprovalForAll",
        "tx_id": "1837c41568668807cad86578e453a672d4d37db17629473190f11815fc120da7"
      }
    ],
    "gas_used": 445239,
    "result": "SetApprovalForAll success"
  },
  "tx_block_height": 7,
  "tx_id": "1837c41568668807cad86578e453a672d4d37db17629473190f11815fc120da7",
  "tx_timestamp": 1745078976
}
```

### 转账
#### 验证Case1：
转账后需要验证Owner是否发生了变化
```sh

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=TransferFrom \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"from\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\", \"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"tokenId\":\"111111111111111111111112\"}"
{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "a04f7895de24f61807a729be230f03da8c0eef42",
          "111111111111111111111112"
        ],
        "topic": "TransferFrom",
        "tx_id": "1837c429a62f5e9dca125cc354aa755e579d2482ee4f46a6b4048df76cbadc55"
      }
    ],
    "gas_used": 1081066,
    "result": "TransferFrom success"
  },
  "tx_block_height": 8,
  "tx_id": "1837c429a62f5e9dca125cc354aa755e579d2482ee4f46a6b4048df76cbadc55",
  "tx_timestamp": 1745079063
}
```

### 批量转账
#### 验证Case1：
转账后需要Owner是否发生了变化
```sh


./cmc client contract user invoke \
--contract-name=CMNFA \
--method=TransferFromBatch \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{
  \"from\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",
  \"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", 
  \"tokenIds\": [
    \"111111111111111111111111\",
    \"111111111111111111111112\"
  ]
}"

111111111111111111111111没有发行过，111111111111111111111112已经转移了所以这个正常会报错
{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 634285,
    "message": "contract message:only owner or approved account can transfer token",
    "result": ""
  },
  "tx_block_height": 9,
  "tx_id": "1837c4395cd4b19bca656bc710fb9704ed2152f01a824cc3afdbf0f9863718bc",
  "tx_timestamp": 1745079131
}
```

### 查询OwnerOf
#### 验证Case1：
验证返回的owner是否为安装合约时指定的tokenURI+'/'+tokenId
```sh

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=OwnerOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "gas_used": 306683,
    "result": "a04f7895de24f61807a729be230f03da8c0eef42"
  },
  "tx_block_height": 10,
  "tx_id": "1837c45cd0cf0ca1ca1505c47016885d3b97052493e741a1a0732a0116ffa32d",
  "tx_timestamp": 1745079283
}
```

### 查询tokenURI
#### 验证Case1：
验证返回的tokenURI是否为安装合约时指定的tokenURI+'/'+tokenId
```sh

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=TokenURI \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "gas_used": 433529,
    "result": "chainmaker.org.cn//111111111111111111111112"
  },
  "tx_block_height": 11,
  "tx_id": "1837c4667497bfd4ca27de8f8d32d44acfd9dcc5116a435f82472ac0a0be4045",
  "tx_timestamp": 1745079324
}
```

### 按分类授权
#### 验证Case1：
授权后需要验证授权信息是否正确
```sh
./cmc client contract user invoke \
--contract-name=CMNFA \
--method=SetApprovalByCategory \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"to\":\"a04f7895de24f61807a729be230f03da8c0eef42\", \"categoryName\":\"1111\", \"isApproval\":\"true\"}"


{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "a04f7895de24f61807a729be230f03da8c0eef42",
          "1111",
          "true"
        ],
        "topic": "SetApprovalByCategory",
        "tx_id": "1837c4735cea4702ca6ac77096f8dcc3d276527af8f549bb8b01f898cf036981"
      }
    ],
    "gas_used": 532243,
    "result": "setApprovalByCategoryCore success"
  },
  "tx_block_height": 12,
  "tx_id": "1837c4735cea4702ca6ac77096f8dcc3d276527af8f549bb8b01f898cf036981",
  "tx_timestamp": 1745079380
}
```

### 创建或设置分类信息
#### 验证Case1：
创建后需要验证分类信息是否正确
```sh
./cmc client contract user invoke \
--contract-name=CMNFA \
--method=CreateOrSetCategory \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{
  \"category\": {
    \"categoryName\":\"1111\",
    \"categoryURI\":\"http://www.chainmaker.org.cn/\"
  }
}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "1111",
          "http://www.chainmaker.org.cn/"
        ],
        "topic": "CreateOrSetCategory",
        "tx_id": "1837c485a4b8f972ca5db04307a59a312bf2a0d17818458f8cdb991f31e0d36b"
      }
    ],
    "gas_used": 791865,
    "result": "CreateOrSetCategory success"
  },
  "tx_block_height": 13,
  "tx_id": "1837c485a4b8f972ca5db04307a59a312bf2a0d17818458f8cdb991f31e0d36b",
  "tx_timestamp": 1745079458
}
```

### 销毁NFA
#### 验证Case1：
销毁后需要验证token是否还存在Owner
```sh

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=Burn \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "code": 4,
  "contract_result": {
    "code": 1,
    "contract_event": [
      {
        "contract_name": "CMNFA",
        "contract_version": "1.0",
        "event_data": [
          "111111111111111111111112"
        ],
        "topic": "Burn",
        "tx_id": "1837c492f4ad2c3cca6b0c932d2065bc86c2b1cd6e3b4367bb745a76bf4d4770"
      }
    ],
    "gas_used": 1404955,
    "message": "contract message:only owner or approved user can Burn the token",
    "result": ""
  },
  "tx_block_height": 14,
  "tx_id": "1837c492f4ad2c3cca6b0c932d2065bc86c2b1cd6e3b4367bb745a76bf4d4770",
  "tx_timestamp": 1745079516
}
```

### 根据名称查询分类信息
```sh


./cmc client contract user invoke \
--contract-name=CMNFA \
--method=GetCategoryByName \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"categoryName\":\"1111\"}"

{
  "contract_result": {
    "gas_used": 835029,
    "result": "{\"categoryName\":\"1111\",\"categoryURI\":\"http://www.chainmaker.org.cn/\"}"
  },
  "tx_block_height": 15,
  "tx_id": "1837c49fd2f822a8cad141a0fa8d340b6723076ef95345a5a1635470ab64eb55",
  "tx_timestamp": 1745079571
}
```

### 根据TokenId查询所属分类信息
```sh


./cmc client contract user invoke \
--contract-name=CMNFA \
--method=GetCategoryByTokenId \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "gas_used": 964714,
    "result": "{\"categoryName\":\"1111\",\"categoryURI\":\"http://www.chainmaker.org.cn/\"}"
  },
  "tx_block_height": 16,
  "tx_id": "1837c4abf0e99182ca749b4e6d52b31be48589a9cc0a43719f6f286168bb3e35",
  "tx_timestamp": 1745079623
}
```

### 查询已发行NFA总数量
```sh
./cmc client contract user invoke \
--contract-name=CMNFA \
--method=TotalSupply \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{}"

{
  "contract_result": {
    "gas_used": 311468,
    "result": "3"
  },
  "tx_block_height": 17,
  "tx_id": "1837c4b924d06327ca2415318b5c6f65bf92ebe7b5594ef78b2d3593b9776212",
  "tx_timestamp": 1745079680
}
```

### 查询某个分类下已发行NFA总数量
```sh

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=TotalSupplyOfCategory \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"categoryName\":\"1111\"}"

{
  "contract_result": {
    "gas_used": 314683,
    "result": "3"
  },
  "tx_block_height": 18,
  "tx_id": "1837c4c6cbddcec0cacf4551713a8c2b282314c05837445f90ab0dfc19d6bab2",
  "tx_timestamp": 1745079738
}
```

### 查询账户NFA数量
```sh


./cmc client contract user invoke \
--contract-name=CMNFA \
--method=BalanceOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"account\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}"

{
  "contract_result": {
    "gas_used": 318454,
    "result": "2"
  },
  "tx_block_height": 19,
  "tx_id": "1837c4d403626d17cacbf2c976fc2365bc15c239c04a4de0b5bd74fd567e57be",
  "tx_timestamp": 1745079795
}
```

### 查询账户下的所有NFA
#### 验证Case1：
验证账户下是否包含了所有发行的NFA
```sh
 ./cmc client contract user invoke \
--contract-name=CMNFA \
--method=AccountTokens \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"account\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}"

{
  "contract_result": {
    "gas_used": 1533111,
    "result": "{\"account\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"tokens\":[\"111111111111111111111112\",\"xxxxx\",\"xxxxx1\"]}"
  },
  "tx_block_height": 14,
  "tx_id": "1837c54f8d077e11cab96259eca009df5423a2e9668748f99d5d1b9626f15373",
  "tx_timestamp": 1745080326
}
```

### 查询token metadata信息
#### 验证Case1：
这儿验证查询到的metadata是否和mint时传递的一致
```sh
./cmc client contract user invoke --contract-name=CMNFA --method=TokenMetadata --sync-result=true --sdk-conf-path=./testdata/sdk_config_solo.yml \
--params="{\"tokenId\":\"111111111111111111111112\"}"

./cmc client contract user invoke \
--contract-name=CMNFA \
--method=TokenMetadata \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"tokenId\":\"111111111111111111111112\"}"

{
  "contract_result": {
    "gas_used": 298514,
    "result": "aaa"
  },
  "tx_block_height": 21,
  "tx_id": "1837c4fb4e2840e5caf4e5ed6b89fdb6aa6b9bb6cf8c4fb3963fc094d5e5dae2",
  "tx_timestamp": 1745079964
}
```
