ContractStandardNameCMDFA ChainMaker - Contract Standard - Digital Fungible Assets

https://git.chainmaker.org.cn/contracts/standard/-/blob/master/living/CM-CS-221221-DFA.md

## CMC示例
### 安装合约
Token名：CM DFA
符号名：DFA
小数位数：8
```sh
UserA:b0831c4aebde4eb5a97b0eb5b1310a746763e752
b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674
#5fa92a33364dd5ce26a9814a6aceb240bd6bf083
UserB:08cd36c7be843d70bfc585ccd20e101e8bb8bc60
UserC:6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962
./cmc client contract user create \
--contract-name=cmdfa \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/standard_dfa-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{\"name\":\"CM DFA\", \"symbol\":\"DFA\", \"decimals\":\"8\"}"

{
  "contract_result": {
    "gas_used": 568071,
    "result": {
      "address": "ef90bd0f031eef572be98bafd219b0b8d7208ba5",
      "creator": {
        "member_id": "client1.sign.wx-org1.chainmaker.org",
        "member_info": "Y9PdAUogrWJ0Z1LsrlgbL7ltvZIi1+MNRPVgholFo9c=",
        "member_type": 1,
        "org_id": "wx-org1.chainmaker.org",
        "role": "CLIENT",
        "uid": "b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674"
      },
      "name": "cmdfa",
      "runtime_type": 2,
      "version": "1.0"
    }
  },
  "tx_block_height": 2,
  "tx_id": "1837be65e0a12d5acaeeed68b9bfebb7f710277624794ddfb1a49a686b968d07",
  "tx_timestamp": 1745072725
}
```
### 铸造Token
铸造10个Token，因为小数位数为8，所以参数的数量是：1000000000
```sh

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=Mint \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"account\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\", \"amount\":\"1000000000\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "cmdfa",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "1000000000"
        ],
        "topic": "mint",
        "tx_id": "1837c10370eb196eca784ac022629de9ca5ff3fb13674e11ba06d1956c02b602"
      }
    ],
    "gas_used": 765744,
    "result": "OK"
  },
  "tx_block_height": 3,
  "tx_id": "1837c10370eb196eca784ac022629de9ca5ff3fb13674e11ba06d1956c02b602",
  "tx_timestamp": 1745075601
}
```
### 查询余额
```sh
./cmc client contract user invoke \
--contract-name=cmdfa \
--method=BalanceOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"account\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}"
{
  "contract_result": {
    "gas_used": 381049,
    "result": "1000000000"
  },
  "tx_block_height": 4,
  "tx_id": "1837c10a482e1103cae741f2a9a49c6d149429ad4a17415faed269120557fd34",
  "tx_timestamp": 1745075630
}
```
### 转账1个Token
转移1个Token，因为小数位数为8，所以参数的数量是：100000000
```sh
./cmc client contract user invoke \
--contract-name=cmdfa \
--method=Transfer \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"to\":\"08cd36c7be843d70bfc585ccd20e101e8bb8bc60\", \"amount\":\"100000000\"}"
{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "cmdfa",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "08cd36c7be843d70bfc585ccd20e101e8bb8bc60",
          "100000000"
        ],
        "topic": "transfer",
        "tx_id": "1837c10ee06fc44aca808be3c793613ac96993552acb4b018d5abf9b6d97907e"
      }
    ],
    "gas_used": 760531,
    "result": "OK"
  },
  "tx_block_height": 5,
  "tx_id": "1837c10ee06fc44aca808be3c793613ac96993552acb4b018d5abf9b6d97907e",
  "tx_timestamp": 1745075650
}
```

### 销毁1个Token
```sh

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=Burn \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"amount\":\"100000000\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "cmdfa",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "100000000"
        ],
        "topic": "burn",
        "tx_id": "1837c1188f075cddca0e17c01b378b918be59a2cb62c4ebb9e603c15e5b85b42"
      }
    ],
    "gas_used": 763561,
    "result": "OK"
  },
  "tx_block_height": 6,
  "tx_id": "1837c1188f075cddca0e17c01b378b918be59a2cb62c4ebb9e603c15e5b85b42",
  "tx_timestamp": 1745075691
}
```

### 查询Token的名字
```sh

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=Name \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{}"

{
  "contract_result": {
    "gas_used": 324170,
    "result": "CM DFA"
  },
  "tx_block_height": 7,
  "tx_id": "1837c122a85bd266ca2a767fa36a6908518d93f5d64e4aefb42e841bcab90560",
  "tx_timestamp": 1745075735
}
```
### 查询Token的符号
```sh

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=Symbol \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{}"

{
  "contract_result": {
    "gas_used": 292360,
    "result": "DFA"
  },
  "tx_block_height": 8,
  "tx_id": "1837c129884b8ee3cabd15bae193b264486818120b19486c8ce99c1a19f845ac",
  "tx_timestamp": 1745075764
}
```

### 查询Token的小数位数
```sh

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=Decimals \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{}"
{
  "contract_result": {
    "gas_used": 291789,
    "result": "8"
  },
  "tx_block_height": 9,
  "tx_id": "1837c1308b052b46ca036638a642fd12a0f5e867b2ab43679942a8f852e2a5f3",
  "tx_timestamp": 1745075794
}
```

### 查询Token的发行量
```sh

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=TotalSupply \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{}"

{
  "contract_result": {
    "gas_used": 323104,
    "result": "900000000"
  },
  "tx_block_height": 10,
  "tx_id": "1837c13843c3efe1ca8f59bbc696b5c4d539e304befe44058a739ce12c32d05b",
  "tx_timestamp": 1745075828
}
```

## 授权转移操作
### 授权3个Token给用户C
```sh
./cmc client contract user invoke \
--contract-name=cmdfa \
--method=Approve \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"spender\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\", \"amount\":\"300000000\"}"

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "cmdfa",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962",
          "300000000"
        ],
        "topic": "approve",
        "tx_id": "1837c1951bf3ee16ca3d17a74ba6da039f0bbb85ef65433a92df00d9ad87c942"
      }
    ],
    "gas_used": 459466,
    "result": "OK"
  },
  "tx_block_height": 12,
  "tx_id": "1837c1951bf3ee16ca3d17a74ba6da039f0bbb85ef65433a92df00d9ad87c942",
  "tx_timestamp": 1745076226
}
```
### C用户转移Token给用户B
因为是用的UserC的身份，所以需要切换SDK_Config到UserC的对应配置
```sh

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=TransferFrom \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"from\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"to\":\"08cd36c7be843d70bfc585ccd20e101e8bb8bc60\", \"amount\":\"100000000\"}"
{
  "code": 4,
  "contract_result": {
    "code": 1,
    "contract_event": [
      {
        "contract_name": "cmdfa",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "08cd36c7be843d70bfc585ccd20e101e8bb8bc60",
          "100000000"
        ],
        "topic": "transfer",
        "tx_id": "1837c22cf27b811aca223a94fcd93be2e71a84a78afb4957adca864f207a1bb2"
      }
    ],
    "gas_used": 1017563,
    "message": "contract message:spend allowance failed, err:CMDFA: insufficient allowance",
    "result": ""
  },
  "tx_block_height": 7,
  "tx_id": "1837c22cf27b811aca223a94fcd93be2e71a84a78afb4957adca864f207a1bb2",
  "tx_timestamp": 1745076878
}
```
### 查询用户A授权给用户C的转移额度
本来授权了3Token，现在转移了1个，所以结果应该2Token
```sh
./cmc client contract user get --contract-name=cmdfa --method=Allowance --sdk-conf-path=./testdata/sdk_config_UserC.yml --params="{\"spender\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"owner\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}"

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=Allowance \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"spender\":\"6fa2bac1c8f5c2c7225c0ec5a3f9262dd315d962\",\"owner\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}"
{
  "contract_result": {
    "gas_used": 384919,
    "result": "300000000"
  },
  "tx_block_height": 8,
  "tx_id": "1837c235efb5c247ca281bb2f65e363591e8f7603aaa43148a46f7438b462773",
  "tx_timestamp": 1745076917
}
```
### UserC从授权账户UserA销毁1Token
因为是用的UserC的身份，所以需要切换SDK_Config到UserC的对应配置
```sh

./cmc client contract user invoke \
--contract-name=cmdfa \
--method=BurnFrom \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--result-to-string=true \
--params="{\"account\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\", \"amount\":\"100000000\"}"

{
  "code": 4,
  "contract_result": {
    "code": 1,
    "contract_event": [
      {
        "contract_name": "cmdfa",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "100000000"
        ],
        "topic": "burn",
        "tx_id": "1837c23d1a6c9cb7ca957bca6386dd3833d5a751c5064d0f8d2b35d799931da0"
      }
    ],
    "gas_used": 855869,
    "message": "contract message:CMDFA: insufficient allowance",
    "result": ""
  },
  "tx_block_height": 9,
  "tx_id": "1837c23d1a6c9cb7ca957bca6386dd3833d5a751c5064d0f8d2b35d799931da0",
  "tx_timestamp": 1745076948
}
```
