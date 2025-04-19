# identity合约
该合约用于身份验证

## cmc测试示例

命令行工具测试示例

```sh
"安装合约 identity"
./cmc client contract user create \
--contract-name=identity \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/identity-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"

```


调用合约函数 查询调用者地址
```sh
callerAddress跨合约调address，callerAddress本质和address没有区别
./cmc client contract user invoke \
--contract-name=identity \
--method=callerAddress \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 7434707,
    "result": "call contract success :sender is 7a83769df9cdfe9c96bf8e01c623e9686a7dc1e796ce12c25ef327d7fd1871ee, len is 64"
  },
  "tx_block_height": 6,
  "tx_id": "18349700dae3b3bbcaf4e430a22b8b6c80cd3553dbe743f9bfd6fe91252b208b",
  "tx_timestamp": 1744184985
}

./cmc client contract user invoke \
--contract-name=identity \
--method=address \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 2370624,
    "result": "sender is 7a83769df9cdfe9c96bf8e01c623e9686a7dc1e796ce12c25ef327d7fd1871ee, len is 64"
  },
  "tx_block_height": 4,
  "tx_id": "183496e49c93cd4cca1dd96a2dc6fb6040f65d10eb6541b586baec9d82296bd1",
  "tx_timestamp": 1744184864
}
```
添加白名单
```sh
./cmc client contract user invoke \
--contract-name=identity \
--method=addWriteList \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"005521f27d745a04999c6d09f559764f9c44376a,335338acd9f5a757a12cf4a35bf327d6051e7068,2a230c7ea2110446e6320a44091089f111cb5028\"}" \
--sync-result=true \
--result-to-string=true
期望输出
{
  "contract_result": {
    "gas_used": 2564677,
    "result": "add write list"
  },
  "tx_block_height": 7,
  "tx_id": "1834970f991a9af0ca48e084d93c737f3b11bc7b4d204a84ab9df7b3606e66d8",
  "tx_timestamp": 1744185048
}
期望输出

```
查询是否在白名单
```sh
./cmc client contract user invoke \
--contract-name=identity \
--method=isApprovedUser \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"005521f27d745a04999c6d09f559764f9c44376a\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 2498716,
    "result": "true"
  },
  "tx_block_height": 5,
  "tx_id": "1834974354460c49cab0f589b1258586493027b14b7543db8c88bbc9365618f5",
  "tx_timestamp": 1744185271
}

```

移除白名单
```sh
./cmc client contract user invoke \
--contract-name=identity \
--method=removeWriteList \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"005521f27d745a04999c6d09f559764f9c44376a\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 2421820,
    "result": "remove write list success"
  },
  "tx_block_height": 6,
  "tx_id": "1834974de03ddf7fca36765c0be2ed0ef9e5aa60d6404038bc794827c6f54db3",
  "tx_timestamp": 1744185316
}
```