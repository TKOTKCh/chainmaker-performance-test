# exchange合约
交易所合约，用于转移ERC721(NFT) token的时候同时转移ERC20的token
## 主要合约接口

```go
// 将nft tokenId从 from 转移到to
// 并将erc20 amount数量的token从to转移到from
buyNow(tokenId, from, to, amount string) protogo.Response // return "true","false"
```
## cmc使用示例

命令行工具使用示例

```sh
echo
echo "安装合约 exchange"

./cmc client contract user create \
--contract-name=exchange \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/exchange-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params=""

echo
echo "执行合约 exchange，购买nft"
./cmc client contract user invoke \
--contract-name=exchange \
--method=buyNow \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"from\":\"2a230c7ea2110446e6320a44091089f111cb5028\",\"to\":\"335338acd9f5a757a12cf4a35bf327d6051e7068\",\"tokenId\":\"1\", \"metadata\":\"url:http://chainmaker.org.cn/\"}" \
--sync-result=true \
--result-to-string=true


{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc721",
        "contract_version": "1.0",
        "event_data": [
          "2a230c7ea2110446e6320a44091089f111cb5028",
          "b0b9c232e14b79e8ab2208fdccea4e1fc64837db21c1c489caa563d217460e33",
          "1"
        ],
        "topic": "ApprovalForAll2",
        "tx_id": "18363561cd51eb7dca289f68dd935009a233eeeba6884de0a6d54b2fcaadffa6"
      },
      {
        "contract_name": "erc721",
        "contract_version": "1.0",
        "event_data": [
          "2a230c7ea2110446e6320a44091089f111cb5028",
          "335338acd9f5a757a12cf4a35bf327d6051e7068",
          "1"
        ],
        "topic": "transfer",
        "tx_id": "18363561cd51eb7dca289f68dd935009a233eeeba6884de0a6d54b2fcaadffa6"
      },
      {
        "contract_name": "exchange",
        "contract_version": "1.0",
        "event_data": [
          "2a230c7ea2110446e6320a44091089f111cb5028335338acd9f5a757a12cf4a35bf327d6051e7068b0b9c232e14b79e8ab2208fdccea4e1fc64837db21c1c489caa563d217460e33"
        ],
        "topic": "buyNow",
        "tx_id": "18363561cd51eb7dca289f68dd935009a233eeeba6884de0a6d54b2fcaadffa6"
      }
    ],
    "gas_used": 86482373,
    "result": "buyNow success"
  },
  "tx_block_height": 7,
  "tx_id": "18363561cd51eb7dca289f68dd935009a233eeeba6884de0a6d54b2fcaadffa6",
  "tx_timestamp": 1744640599
}
```