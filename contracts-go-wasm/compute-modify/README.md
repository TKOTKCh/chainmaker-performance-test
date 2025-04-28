# compute合约
该合约为性能测试用自己编写，包含简单计算、hash计算、大数计算

## cmc测试示例

命令行工具测试示例

```sh
"安装合约 compute"
./cmc client contract user create \
--contract-name=compute \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/compute-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"

./cmc client contract user invoke \
--contract-name=compute \
--method=manualInit \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}" \
--sync-result=true \
--result-to-string=true
```
```sh
"执行合约函数 normalCal，循环累加1000000次"
./cmc client contract user invoke \
--contract-name=compute \
--method=normalCal \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}" \
--sync-result=true \
--result-to-string=true

"期望输出："
{
  "contract_result": {
    "gas_used": 37435737,
    "result": "success normalCal: 499999500000"
  },
  "tx_block_height": 3,
  "tx_id": "183467c460d09b30ca545be9b6d9e7eabb782a5ffdf843f09c9a454b473c0958",
  "tx_timestamp": 1744133048
}
```


```sh
"执行合约函数 hashCal，sha256计算100000次"
./cmc client contract user invoke \
--contract-name=compute \
--method=hashCal \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}" \
--sync-result=true \
--result-to-string=true

"期望输出："
{
  "contract_result": {
    "gas_used": 2353580788,
    "result": "success hashCal 2f49fd75f950f2f513de583f0105ece419c34e09263db87013e62724eea8b44d"
  },
  "tx_block_height": 3,
  "tx_id": "18346868577be1bdcae3fb455fee94d40c7d2ff0ba3c4d92afedc92926b63dac",
  "tx_timestamp": 1744133753
}
```

```sh
"执行合约函数 bigNumCal，10000次求幂取余大数运算"
./cmc client contract user invoke \
--contract-name=compute \
--method=bigNumCal \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}" \
--sync-result=true \
--result-to-string=true

"期望输出："
{
  "contract_result": {
    "gas_used": 654526981,
    "result": "success bigNumCal: 607723520"
  },
  "tx_block_height": 5,
  "tx_id": "183468877178e8a8ca46ad39629fba474e605089585c4015af6de3de1a45b032",
  "tx_timestamp": 1744133886
}
```