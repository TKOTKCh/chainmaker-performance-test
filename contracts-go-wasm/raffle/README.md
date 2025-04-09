This is the raffle contract for test net demo.

## Test

### 安装合约
```sh
./cmc client contract user create \
--contract-name=raffle \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/raffle-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"

./cmc client contract user invoke \
--contract-name=raffle \
--method=manualInit \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 2312031,
    "result": "Init contract success"
  },
  "tx_block_height": 8,
  "tx_id": "183497699544f5f8cac918a491021df4cdb4f07d90dc41ab8a064fc2b24f12b4",
  "tx_timestamp": 1744185435
}
```

### 注册奖池名单
```sh
./cmc client contract user invoke \
--contract-name=raffle \
--method=registerAll \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"peoples\":{\"peoples\":[{\"num\":1,\"name\":\"Chris\"},{\"num\":2,\"name\":\"Linus\"}]}}"
期望输出
{
  "contract_result": {
    "gas_used": 3131255,
    "result": "ok"
  },
  "tx_block_height": 9,
  "tx_id": "1834976e2b90c0d9cab1a04ecb00e6a414b038991ad04b72a9d931b7be41f29b",
  "tx_timestamp": 1744185455
}
```

### 抽奖
### 测试Case1: 
参数level或timestamp没有时应报错
### 测试Case2:
参数level和timestamp正确时应返回从奖池中抽出中奖的人员
### 测试Case3:
多次抽奖，中奖人员不应重复
### 测试Case4:
抽奖后查询奖池人员，应不包含已中奖人员
```sh
./cmc client contract user invoke \
--contract-name=raffle \
--method=raffle \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"level\":\"1\",\"timestamp\":\"13235432\"}"

期望输出
{
  "contract_result": {
    "gas_used": 3358668,
    "result": "num: 2, name: Linus, level: 1"
  },
  "tx_block_height": 10,
  "tx_id": "18349773c91159b3ca3eb9edc85f1a734177fab90f2c4c56acc754300913bcb7",
  "tx_timestamp": 1744185479
}
```

### 查询奖池名单
### 测试Case1:
应返回正确的奖池人员名单
```sh
./cmc client contract user invoke \
--contract-name=raffle \
--method=query \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml 

期望输出
{
  "contract_result": {
    "gas_used": 2453969,
    "result": "{\"peoples\":[{\"num\":1,\"name\":\"Chris\"}]}"
  },
  "tx_block_height": 11,
  "tx_id": "18349781781f4f67ca0814d55a863dba650c83b2855b41d49e4951cf10df1440",
  "tx_timestamp": 1744185538
}
```
