# evidence合约

存证合约，提供存证、查验功能。

## 主要合约接口

参考： 长安链CMEVI（CM-CS-221221-Evidence）存证合约标准实现：
https://git.chainmaker.org.cn/contracts/standard/-/blob/master/living/CM-CS-221221-Evidence.md

```go

// CMEVI 长安链存证合约go接口
// https://git.chainmaker.org.cn/contracts/wx-organdard/-/blob/master/draft/CM-CS-221221-Evidence.md


// Evidence 存证
// @param id 必填，流水号
// @param hash 必填，上链哈希值
// @param metadata 可选，其他信息；比如：哈希的类型（文字，文件）、文字描述的json格式字符串，具体参考下方 Metadata 对象。
// @return error 返回错误信息
Evidence(id string, hash string, metadata string) error

// ExistsOfHash 哈希是否存在
// @param hash 必填，上链的哈希值
// @return exist 存在：true，"true"；不存在：false，"false"
// @return err 错误信息
ExistsOfHash(hash string) (exist bool, er error)

// ExistsOfId 哈希是否存在
// @param id 必填，上链的ID值
// @return exist 存在：true，"true"；不存在：false，"false"
// @return err 错误信息
ExistsOfId(id string) (exist bool, er error)

// FindByHash 根据哈希查找
// @param hash 必填，上链哈希值
// @return evidence 上链时传入的evidence信息
// @return err 返回错误信息
FindByHash(hash string) (evidence *Evidence, err error)

// FindById 根据id查找
// @param id 必填，流水号
// @return evidence 上链时传入的evidence信息
// @return err 返回错误信息
FindById(id string) (evidence *Evidence, err error)

// EvidenceBatch 批量存证
// @param evidences 必填，存证信息
// @return error 返回错误信息
EvidenceBatch(evidences []Evidence) error

// UpdateEvidence 根据ID更新存证哈希和metadata
// @param id 必填，已经在链上存证的流水号。 如果是新流水号返回错误信息不存在
// @param hash 必填，上链哈希值。必须与链上已经存储的hash不同
// @param metadata 可选，其他信息；具体参考下方 Metadata 对象。
// @return error 返回错误信息
// @desc 该方法由长安链社区志愿者@sunhuiyuan提供建议，感谢支持
UpdateEvidence(id string, hash string, metadata string) error
```

## cmc使用示例

命令行工具使用示例

```sh
部署合约
./cmc client contract user create \
--contract-name=evidence \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/standard-evidence-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"

./cmc client contract user invoke \
--contract-name=evidence \
--method=manualInit \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 2348964,
    "result": "Init contract success"
  },
  "tx_block_height": 3,
  "tx_id": "183499630b2070d2caaa867996c79793dea6fa8c332b4c469669ab6052b61e49",
  "tx_timestamp": 1744187606
}

升级合约
./cmc client contract user upgrade \
--contract-name=evidence \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/standard-evidence-go.wasm \
--version=2.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"

```
```sh
echo "invoke evidence.Evidence(id, hash, metadata string)  hash:a2 id:02"
./cmc client contract user invoke \
--contract-name=evidence \
--method=Evidence \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"id\":\"02\",\"hash\":\"a2\",\"metadata\":\"{\\\"hashType\\\":\\\"file\\\",\\\"hashAlgorithm\\\":\\\"sha256\\\",\\\"username\\\":\\\"taifu\\\",\\\"timestamp\\\":\\\"1672048892\\\",\\\"proveTimestamp\\\":\\\"\\\"}\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 3506712,
    "result": "evidence saved"
  },
  "tx_block_height": 4,
  "tx_id": "1834996eca8a6d7dcad5b93ff911c9dee2a4fb540f3448a59e78bd1edbc90bed",
  "tx_timestamp": 1744187656
}


```
```sh
echo "invoke evidence.EvidenceBatch(evidences []Evidence)"
./cmc client contract user invoke \
--contract-name=evidence \
--method=EvidenceBatch \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"evidences\":\"[{\\\"id\\\":\\\"id1\\\",\\\"hash\\\":\\\"hash1\\\",\\\"txId\\\":\\\"\\\",\\\"blockHeight\\\":0,\\\"timestamp\\\":\\\"\\\",\\\"metadata\\\":\\\"11\\\"},{\\\"id\\\":\\\"id2\\\",\\\"hash\\\":\\\"hash2\\\",\\\"txId\\\":\\\"\\\",\\\"blockHeight\\\":0,\\\"timestamp\\\":\\\"\\\",\\\"metadata\\\":\\\"11\\\"}]\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 4112388,
    "result": "ok"
  },
  "tx_block_height": 5,
  "tx_id": "1834997b9e9e2d01ca8f49f6a231b7da5ffc854381204fa492b1bffe74dc8790",
  "tx_timestamp": 1744187711
}


```
```sh
echo "query evidence.ExistsOfHash(hash string) hash:a2"
./cmc client contract user get \
--contract-name=evidence \
--method=ExistsOfHash \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"hash\":\"a2\"}" \
--result-to-string=true

期望输出
{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 2513923,
    "message": "contract message:[existsOfHash]hash isn't exists"
  },
  "message": "txStatusCode:4, resultCode:1, contractName[evidence] method[ExistsOfHash] txType[QUERY_CONTRACT], contract message:[existsOfHash]hash isn't exists",
  "tx_id": "18349a5be8f19572ca2586447b5d954bb2cb58720de14d7089d3914bcad70356"
}

```
```sh
echo "query evidence.ExistsOfId(id string) id:02"
./cmc client contract user get \
--contract-name=evidence \
--method=ExistsOfId \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"id\":\"02\"}" \
--result-to-string=true
期望输出
{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 2497679,
    "message": "contract message:[existsOfId]id isn't exists"
  },
  "message": "txStatusCode:4, resultCode:1, contractName[evidence] method[ExistsOfId] txType[QUERY_CONTRACT], contract message:[existsOfId]id isn't exists",
  "tx_id": "18349a6ac5148affca5b9952fdd0c435ba456ad006664900883916ca4ec741c2"
}

```
```sh
echo "query evidence.FindByHash(hash string) hash:a2"
./cmc client contract user get \
--contract-name=evidence \
--method=FindByHash \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"hash\":\"a2\"}" \
--result-to-string=true
期望输出
{
  "contract_result": {
    "gas_used": 3587882,
    "result": "{\"id\":\"02\",\"hash\":\"a2\",\"txId\":\"18349a38fee8c350ca1cbb192c9beccb0baba7525e184c79a72afe2867f28a7b\",\"blockHeight\":3,\"timestamp\":\"\",\"metadata\":\"{\\\"hashType\\\":\\\"file\\\",\\\"hashAlgorithm\\\":\\\"sha256\\\",\\\"username\\\":\\\"taifu\\\",\\\"timestamp\\\":\\\"1672048892\\\",\\\"proveTimestamp\\\":\\\"\\\"}\"}"
  },
  "message": "SUCCESS",
  "tx_id": "18349a7283a11c33ca8030552781fe1f97095403af5a4090a775f4706510b239"
}

```
```sh
echo "query evidence.FindById(id string) id:02"
./cmc client contract user get \
--contract-name=evidence \
--method=FindById \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"id\":\"02\"}" \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 3717510,
    "result": "{\"id\":\"02\",\"hash\":\"a2\",\"txId\":\"18349a97fc13c02dca00bc0d1fa296f7a8a5968adb33414980f058635cac3f4e\",\"blockHeight\":3,\"timestamp\":\"\",\"metadata\":\"{\\\"hashType\\\":\\\"file\\\",\\\"hashAlgorithm\\\":\\\"sha256\\\",\\\"username\\\":\\\"taifu\\\",\\\"timestamp\\\":\\\"1672048892\\\",\\\"proveTimestamp\\\":\\\"\\\"}\"}"
  },
  "message": "SUCCESS",
  "tx_id": "18349a9cb2c0b72dcacaad908fa8b27e20884238f52d4bf98b5b7b1700e0adc9"
}


```
```sh
echo "查询合约 evidence.SupportStandard，是否支持合约标准CMEVI"
./cmc client contract user get \
--contract-name=evidence \
--method=SupportStandard \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMEVI\"}" \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 2364161,
    "result": "true"
  },
  "message": "SUCCESS",
  "tx_id": "18349b960f75c9b2ca3555cecba733de531622522396462388f9d7c7dd6fa404"
}

```
```sh
echo "查询合约 evidence.SupportStandard，是否支持合约标准CMBC"
./cmc client contract user get \
--contract-name=evidence \
--method=SupportStandard \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMBC\"}" \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 2376805,
    "result": "true"
  },
  "message": "SUCCESS",
  "tx_id": "18349b7a0b053084ca4d970be12151b365cf4c52e02744049b7570793c34c8c5"
}

```
```sh
echo "查询合约 evidence.Standards，支持的合约标准列表"
./cmc client contract user get \
--contract-name=evidence \
--method=Standards \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 2478019,
    "result": "[\"CMEVI\",\"CMBC\"]"
  },
  "message": "SUCCESS",
  "tx_id": "18349b98fb57b923caf18b3d6595155b15e78c5fd42a4481827c5fee7267d7d0"
}
```

