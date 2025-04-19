# erc1155合约
ERC1155 Token Standard:
https://eips.ethereum.org/EIPS/eip-1155

用于批量发行nft和token，包含了erc20和erc721的功能。

接口详情请查看接口：IERC1155
此处主要说明和标准合约的区别在于：多了接口：MintBatchNft，用户批量发行nft，有一些nft的校验。


## cmc使用示例

命令行工具使用示例

```sh
b0831c4aebde4eb5a97b0eb5b1310a746763e752
b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674
echo
echo "安装合约 erc1155，设置管理员地址，若不手动设置，则为创建合约的人，管理员可 mint"

./cmc client contract user create \
--contract-name=erc1155 \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/erc1155-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{\"uri\":\"https://www.fxtoon.com/nft/{tokenId}.json\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc1155",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752"
        ],
        "topic": "AlterAdminAddress",
        "tx_id": "1837b5d82fdfc4f9cace574c454b87d53694547a83754b38adbe649885761498"
      }
    ],
    "gas_used": 647357,
    "result": {
      "address": "7c16571011799474b454d612d94d70b18af70a43",
      "creator": {
        "member_id": "client1.sign.wx-org1.chainmaker.org",
        "member_info": "Y9PdAUogrWJ0Z1LsrlgbL7ltvZIi1+MNRPVgholFo9c=",
        "member_type": 1,
        "org_id": "wx-org1.chainmaker.org",
        "role": "CLIENT",
        "uid": "b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674"
      },
      "name": "erc1155",
      "runtime_type": 2,
      "version": "1.0"
    }
  },
  "tx_block_height": 2,
  "tx_id": "1837b5d82fdfc4f9cace574c454b87d53694547a83754b38adbe649885761498",
  "tx_timestamp": 1745063320
}


echo
echo "set admin"
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=AlterAdminAddress \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"adminAddress\":\"5fa92a33364dd5ce26a9814a6aceb240bd6bf083,08cd36c7be843d70bfc585ccd20e101e8bb8bc60,b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc1155",
        "contract_version": "1.0",
        "event_data": [
          "5fa92a33364dd5ce26a9814a6aceb240bd6bf083",
          "08cd36c7be843d70bfc585ccd20e101e8bb8bc60",
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752"
        ],
        "topic": "AlterAdminAddress",
        "tx_id": "1837b8dfe8b8e648cae0bc8ed8d905ed56dc2d9b166e4dfdb6352306e022d7fa"
      }
    ],
    "gas_used": 764222,
    "result": "ok"
  },
  "tx_block_height": 3,
  "tx_id": "1837b8dfe8b8e648cae0bc8ed8d905ed56dc2d9b166e4dfdb6352306e022d7fa",
  "tx_timestamp": 1745066652
}
echo
echo "Mint erc721 给admin2发送token value为1 "
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=Mint \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"to\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"id\":\"1\",\"amount\":\"1\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc1155",
        "contract_version": "1.0",
        "event_data": [
          "6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962",
          "1",
          "1"
        ],
        "topic": "Mint",
        "tx_id": "1837b690d62fb51dcab81d4b91f15395f537a821faa04c7185bf27df4be23093"
      }
    ],
    "gas_used": 979513,
    "result": "ok"
  },
  "tx_block_height": 4,
  "tx_id": "1837b690d62fb51dcab81d4b91f15395f537a821faa04c7185bf27df4be23093",
  "tx_timestamp": 1745064113
}


echo
echo "Mint erc20 给admin2发送token value为2 "
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=Mint \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"to\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"id\":\"2\",\"amount\":\"2\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc1155",
        "contract_version": "1.0",
        "event_data": [
          "6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962",
          "2",
          "2"
        ],
        "topic": "Mint",
        "tx_id": "1837b69911981c65ca6d3020dd6b1253b0f703c68948464db2fe4b7298cb865d"
      }
    ],
    "gas_used": 979513,
    "result": "ok"
  },
  "tx_block_height": 5,
  "tx_id": "1837b69911981c65ca6d3020dd6b1253b0f703c68948464db2fe4b7298cb865d",
  "tx_timestamp": 1745064149
}

echo
echo "MintBatchNft admin2 2 3 4"
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=MintBatchNft \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"to\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"amount\":\"3\",\"idStart\":\"3\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc1155",
        "contract_version": "1.0",
        "event_data": [
          "6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962",
          "[\"3\",\"4\",\"5\"]",
          "[\"1\",\"1\",\"1\"]"
        ],
        "topic": "MintBatch",
        "tx_id": "1837b6a6b3078103cabccdcdf2e59c974bbadae1dc4a43cbb23127e857ec153b"
      }
    ],
    "gas_used": 1884238,
    "result": "ok"
  },
  "tx_block_height": 6,
  "tx_id": "1837b6a6b3078103cabccdcdf2e59c974bbadae1dc4a43cbb23127e857ec153b",
  "tx_timestamp": 1745064207
}

echo
echo "BalanceOf admin1 token1"
./cmc client contract user get \
--contract-name=erc1155 \
--method=BalanceOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"owner\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"id\":\"1\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 332700,
    "result": "0"
  },
  "message": "SUCCESS",
  "tx_id": "1837b6b5454df1bbcaf81f8adfda5a00f36ba74e40c8487ea3fba2800db797d5"
}

echo
echo "BalanceOf admin2 token1"
./cmc client contract user get \
--contract-name=erc1155 \
--method=BalanceOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"owner\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"id\":\"1\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 356401,
    "result": "1"
  },
  "message": "SUCCESS",
  "tx_id": "1837b6c39caa6b73ca7b2ed250a3ab27feb122d5c5ff451e891d5d5253db56f2"
}

echo
正常admin1没有资产，且sender没权限转不了
echo "SafeTransferFrom token 1 admin1 to admin2 err"
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=SafeTransferFrom \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"from\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"to\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"id\":\"1\",\"amount\":\"1\"}" \
--sync-result=true \
--result-to-string=true

{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 480665,
    "message": "contract message:check isApprovedOrOwner failed, err:Is not Approved For All",
    "result": "0"
  },
  "tx_block_height": 7,
  "tx_id": "1837b6cf43a4a310cace5e3dee41618208841227c2c5407fad77ee5cb3a797a5",
  "tx_timestamp": 1745064381
}


echo
echo "SafeTransferFrom token 1 admin2 to admin1"
sender没权限转不了
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=SafeTransferFrom \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"from\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"to\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"id\":\"1\",\"amount\":\"1\"}" \
--sync-result=true \
--result-to-string=true

{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 480665,
    "message": "contract message:check isApprovedOrOwner failed, err:Is not Approved For All",
    "result": "0"
  },
  "tx_block_height": 7,
  "tx_id": "1837b74ab5e1f7daca52c38fd8f4871831812876284844b8a4f624708edf6af4",
  "tx_timestamp": 1745064912
}

echo
echo "BalanceOfBatch admin2 admin1 token1"
./cmc client contract user get \
--contract-name=erc1155 \
--method=BalanceOfBatch \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"owner\":\"sdk_config,6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"id\":\"1,1\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 485333,
    "result": "[\"0\",\"1\"]"
  },
  "message": "SUCCESS",
  "tx_id": "1837b75a2967355ecad0597f8a35686de9b826a1cacd4560b7ca33e6270a9c9e"
}


echo
echo "SafeTransferFrom token 1 admin1 to admin2, err sender not owner"
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=SafeTransferFrom \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"from\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"to\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"id\":\"1\",\"amount\":\"1\"}" \
--sync-result=true \
--result-to-string=true

{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 416507,
    "message": "contract message:check isApprovedOrOwner failed, err:Is not Approved For All",
    "result": "0"
  },
  "tx_block_height": 8,
  "tx_id": "1837b763bcd16c1bcabeb35c9fb9657607f358b0a57c404ab7f7a709d84be7c4",
  "tx_timestamp": 1745065019
}

echo
echo "SetApprovalForAll admin1 to admin2 no sender"
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=SetApprovalForAll \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"operator\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"approved\":\"true\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc1155",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962",
          "1"
        ],
        "topic": "ApprovalForAll",
        "tx_id": "1837b770f63acd37ca16de688fb8938fb346cc0a69684122a7ac1cd732fd64ae"
      }
    ],
    "gas_used": 402437,
    "result": "ok"
  },
  "tx_block_height": 9,
  "tx_id": "1837b770f63acd37ca16de688fb8938fb346cc0a69684122a7ac1cd732fd64ae",
  "tx_timestamp": 1745065076
}
echo
echo "SetApprovalForAll admin1 to admin2 with sender"
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=SetApprovalForAll \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"operator\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"sender\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"approved\":\"true\",\"sign\":\"3046022100a4afcc04659d5f8bbf831daf0e6e98c08520cf95a482b58ccd5f804b195ff430022100f71d1ff8c9cdc6f90570d88dc3daa39c3dbd158b2bc1ace0b40cae054c3a9444\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "erc1155",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962",
          "1"
        ],
        "topic": "ApprovalForAll",
        "tx_id": "1837b8fefd11ea41ca47d1e785fd8829605c6a03793741648edf3e64edde3d04"
      }
    ],
    "gas_used": 422399,
    "result": "ok"
  },
  "tx_block_height": 9,
  "tx_id": "1837b8fefd11ea41ca47d1e785fd8829605c6a03793741648edf3e64edde3d04",
  "tx_timestamp": 1745066785
}



echo
echo "IsApprovedForAll admin1 to admin2"
./cmc client contract user get \
--contract-name=erc1155 \
--method=IsApprovedForAll \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"operator\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"owner\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 339286,
    "result": "1"
  },
  "message": "SUCCESS",
  "tx_id": "1837b8ac3b44653ecaf8aafe1278e76ef17c863d8555418d80d3d905e3859d83"
}

echo
echo "SafeTransferFrom token 1 admin1 to admin2"
./cmc client contract user invoke \
--contract-name=erc1155 \
--method=SafeTransferFrom \
--sdk-conf-path=./testdata/sdk_config.yml \
--gas-limit=100000000 \
--params="{\"from\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"to\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"id\":\"1\",\"amount\":\"1\"}" \
--sync-result=true \
--result-to-string=true

{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 510371,
    "message": "contract message:insufficient balance from:0",
    "result": ""
  },
  "tx_block_height": 11,
  "tx_id": "1837b913a50162cccabe1e05ff68d05eba486e5307094528b0a04dcc6eb2454e",
  "tx_timestamp": 1745066874
}

echo
echo "Uri 1"
./cmc client contract user get \
--contract-name=erc1155 \
--method=Uri \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"id\":\"1\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 332019,
    "result": "https://www.fxtoon.com/nft/1.json"
  },
  "message": "SUCCESS",
  "tx_id": "1837b91d9f5649c9caf0ce3b54d0864707b1288d901b4875b73d9c2a25821dc8"
}

echo
echo "OwnerOf 1"
./cmc client contract user get \
--contract-name=erc1155 \
--method=OwnerOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"id\":\"1\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 325432,
    "result": "6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962"
  },
  "message": "SUCCESS",
  "tx_id": "1837b923f1b3843dca3b21a0303865a76e5e4a4dacaa4422b13912fb4cb50716"
}
```