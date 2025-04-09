# 一、接口设计
&nbsp;&nbsp;主要定义了溯源最基本的流程，从生产商品、运输商品、上架商品、销售商品，整体过程都记录了相应的溯源信息。
同时提供了获取商品当前状态以及获取商品所有溯源记录的接口。
  ```     
 type Trace interface {

	//创建商品
	newGoods() protogo.Response

	//运输商品
	transferGoods() protogo.Response

	//上架商品
	uploadGoods() protogo.Response

	//销售商品
	sellGoods() protogo.Response

	//获取商品当前状态
	getGoodsStatus() protogo.Response

	//获取溯源信息
	getTraceInfo() protogo.Response
    }
```
# 二、结构体定义
## 1.商品信息
商品信息做了简化，只定义了商品id、商品名称、商品状态、商品溯源信息：
```
// Goods 商品信息
type Goods struct {
	//商品id
	GoodsId string
	//商品名称
	Name string
	//商品状态
	Status uint8
	//商品溯源信息
	TraceDatas []*TraceData
}
```
## 2.溯源信息
溯源信息包含操作人、状态、操作时间、备注
```
//TraceData 溯源信
type TraceData struct {
	//操作人
	Operator string
	//状态{0 :生产  1： 运输 2：上架  3：售卖}
	Status uint8
	//操作时间
	OperatorTime string
	//备注
	Remark string
}
```
# 三、接口实现及测试
## 0.部署合约
部署合约用于实现合约部署，具体看源码
```sh
./cmc client contract user create \
--contract-name=TraceConstract \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/trace-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"
```
## 1.生产商品
&nbsp;&nbsp;创建一件商品，并存储到商品库里，同时修改商品状态，以及记录溯源信息，具体代码看源码
```sh
./cmc client contract user invoke \
--contract-name=TraceConstract \
--method=newGoods \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"goodsId\":\"2000\",\"name\": \"apple\",\"factory\":\"chainMakerFC\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 3534221,
    "result": "newGoods success"
  },
  "tx_block_height": 3,
  "tx_id": "18349c6db5be5464cafd541bd029b44c54518236ed874af7bcf8bc76a025fd92",
  "tx_timestamp": 1744190950
}

```
## 2.运输商品
&nbsp;&nbsp;修改商品运输状态，记录溯源信息， 具体代码看源码
```sh
./cmc client contract user invoke \
--contract-name=TraceConstract \
--method=transferGoods \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"goodsId\":\"2000\",\"from\":\"hunan\",\"to\":\"hubei\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 3910712,
    "result": "transferGoods success"
  },
  "tx_block_height": 4,
  "tx_id": "18349c7b72473e02caa219ec347fb8e57b4a065b9da44187bf21ba9ee5fb42a8",
  "tx_timestamp": 1744191009
}
```

## 3.上架商品
&nbsp;&nbsp;修改商品上架状态，记录溯源信息， 具体代码看源码
```sh
./cmc client contract user invoke \
--contract-name=TraceConstract \
--method=uploadGoods \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"goodsId\":\"2000\",\"uploader\": \"lihongyi\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 4104642,
    "result": "uploadGoods success"
  },
  "tx_block_height": 5,
  "tx_id": "18349c8c9286c69fca8fa6f7369764bd1d76be333dde4d8fa9cd69bfd6468942",
  "tx_timestamp": 1744191083
}
```
## 4.销售商品
&nbsp;&nbsp;修改商品销售状态，记录溯源信息， 具体代码看源码
```sh
./cmc client contract user invoke \
--contract-name=TraceConstract \
--method=sellGoods \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"goodsId\":\"2000\",\"seller\": \"lihongyi\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 4310536,
    "result": "sellGoods success"
  },
  "tx_block_height": 6,
  "tx_id": "18349c9b78b73120cad9d4eda077dffbdf072097712c4ab7b0bbbc46310ca863",
  "tx_timestamp": 1744191147
}
```
## 5.获取商品当前状态
&nbsp;&nbsp;获取商品当前状态， 具体代码看源码
```sh
./cmc client contract user invoke \
--contract-name=TraceConstract \
--method=goodsStatus \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"goodsId\":\"2000\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 4137369,
    "result": "11"
  },
  "tx_block_height": 7,
  "tx_id": "18349cfef7941ba3cab46ffaf039bb019718092a84404cb293d40f3819b43e32",
  "tx_timestamp": 1744191574
}
```

## 6.获取商品溯源信息
&nbsp;&nbsp;获取商品溯源信息， 具体代码看源码
```sh
./cmc client contract user invoke \
--contract-name=TraceConstract \
--method=traceGoods \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"goodsId\":\"2000\"}" \
--sync-result=true \
--result-to-string=true

期望输出
{
  "contract_result": {
    "gas_used": 4336515,
    "result": "[{\"Operator\":\"7a83769df9cdfe9c96bf8e01c623e9686a7dc1e796ce12c25ef327d7fd1871ee\",\"Status\":8,\"OperatorTime\":\"\",\"Remark\":\"2000:chainMakerFC created\"},{\"Operator\":\"7a83769df9cdfe9c96bf8e01c623e9686a7dc1e796ce12c25ef327d7fd1871ee\",\"Status\":9,\"OperatorTime\":\"\",\"Remark\":\"2000:hunan-\\u003ehubei\"},{\"Operator\":\"7a83769df9cdfe9c96bf8e01c623e9686a7dc1e796ce12c25ef327d7fd1871ee\",\"Status\":10,\"OperatorTime\":\"\",\"Remark\":\"2000:lihongyi upload\"},{\"Operator\":\"7a83769df9cdfe9c96bf8e01c623e9686a7dc1e796ce12c25ef327d7fd1871ee\",\"Status\":11,\"OperatorTime\":\"\",\"Remark\":\"2000:selled by lihongyi\"}]"
  },
  "tx_block_height": 7,
  "tx_id": "18349d812dc73faacaf5b594030e50c5226b0d564c7a42a390a4395dbd5fb106",
  "tx_timestamp": 1744192133
}
```