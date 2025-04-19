## Test
### 安装合约
```sh
./cmc client contract user create \
--contract-name=fact \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/fact-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"
{
  "contract_result": {
    "gas_used": 12300,
    "result": {
      "address": "c867f29fbc90b619206c7aef0ba4a2efb30ab9f6",
      "creator": {
        "member_id": "client1.sign.wx-org1.chainmaker.org",
        "member_info": "Y9PdAUogrWJ0Z1LsrlgbL7ltvZIi1+MNRPVgholFo9c=",
        "member_type": 1,
        "org_id": "wx-org1.chainmaker.org",
        "role": "CLIENT",
        "uid": "b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674"
      },
      "name": "fact",
      "runtime_type": 2,
      "version": "1.0"
    }
  },
  "tx_block_height": 6,
  "tx_id": "1837c646955b5758cacf4e03486067733f5238b7c61f4398a7570e8b0381e10a",
  "tx_timestamp": 1745081387
}
```



### 存入存证数据
#### 测试Case1:
存入存证数据后需要验证数据是否能够查询得到
```sh
./cmc client contract user invoke \
--contract-name=fact \
--method=save \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true \
--params="{\"file_hash\":\"005521f27d745a04999c6d09f559764f9c44376a\",\"file_name\":\"aoteman.jpg\",\"time\":\"16456254\"}" \
--sync-result=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "fact",
        "contract_version": "1.0",
        "event_data": [
          "005521f27d745a04999c6d09f559764f9c44376a",
          "aoteman.jpg"
        ],
        "topic": "topic_vx",
        "tx_id": "1837c65bb762fb0dca5a9c1e13efb28f61829af2196148d7a98e917e6853c9a4"
      }
    ],
    "gas_used": 624725,
    "result": "aoteman.jpg005521f27d745a04999c6d09f559764f9c44376a"
  },
  "tx_block_height": 7,
  "tx_id": "1837c65bb762fb0dca5a9c1e13efb28f61829af2196148d7a98e917e6853c9a4",
  "tx_timestamp": 1745081477
}
```


### 查询数据
#### 测试Case1:
数据不存在时，查到数据应为空
```sh
./cmc client contract user invoke \
--contract-name=fact \
--method=findByFileHash \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true \
--params="{\"file_hash\":\"005521f27d745a04999c6d09f559764f9c44376a\"}" \
--sync-result=true

{
  "contract_result": {
    "gas_used": 693384,
    "result": "{\"fileHash\":\"005521f27d745a04999c6d09f559764f9c44376a\",\"fileName\":\"aoteman.jpg\",\"time\":16456254}"
  },
  "tx_block_height": 4,
  "tx_id": "1837c67d604fa6accaf636e9af84671ba0ae3ce6f3e944a98c8b8c42223e2660",
  "tx_timestamp": 1745081622
}
```