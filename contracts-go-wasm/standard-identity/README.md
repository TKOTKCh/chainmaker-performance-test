# identity合约
身份认证合约，提供设置权限，查询权限、公钥功能。可在实际合约中，使用跨合约调用identity合约检查该使用者的权限。


## 主要合约接口

参考： 长安链CMID（CM-CS-221221-Identity）存证合约标准实现：
https://git.chainmaker.org.cn/contracts/standard/-/blob/master/living/CM-CS-221221-Identity.md

```go
// Identities 获取该合约支持的所有认证类型
// @return metas, 所有的认证类型编号和认证类型描述
Identities() (metas []IdentityMeta)

// SetIdentity 为地址设置认证类型，管理员可调用
// @param address 必填，公钥/证书的地址。一个地址仅能绑定一个公钥和认证类型编号，重复输入则覆盖。
// @param pkPem 选填,pem格式公钥，可用于验签
// @param level 必填,认证类型编号
// @param metadata 选填,其他信息，json格式字符串，比如：地址类型，上链人身份、组织信息，上链可信时间，上链批次等等
// @return error 返回错误信息
// @event topic: setIdentity(address, level, pkPem)
SetIdentity(address, pkPem string, level int, metadata string) error

// IdentityOf 获取认证信息
// @param address 地址
// @return int 返回当前认证类型编号
// @return identity 认证信息
// @return err 返回错误信息
IdentityOf(address string) (identity IdentityMeta, err error)

// LevelOf 获取认证编号
// @param address 地址
// @return level 返回当前认证类型编号
// @return err 返回错误信息
LevelOf(address string) (level int, err error)

// EmitSetIdentityEvent 发送设置认证类型事件
// @param address 地址
// @param pkPem pem格式公钥
// @param level 认证类型编号
EmitSetIdentityEvent(address, pkPem string, level int)

// PkPemOf 获取公钥
// @param address 地址
// @return string 返回当前地址绑定的公钥
// @return error 返回错误信息
PkPemOf(address string) (string, error)

// SetIdentityBatch 设置多个认证类型，管理员可调用
// @param identities, 入参json格式字符串
// @event topic: setIdentity(address, level, pkPem)
SetIdentityBatch(identities []Identity) error

// AlterAdminAddress 修改管理员，管理员可调用
// @param adminAddresses 管理员地址，可为空，默认为创建人地址。入参为以逗号分隔的地址字符串"addr1,addr2"
// @return error 返回错误信息
// @event topic: alterAdminAddress（adminAddresses）
AlterAdminAddress(adminAddresses string) error
```

## 跨合约使用

```go
args := make(map[string][]byte)
args["address"] = []byte("55e5a6d2f6a66b6c3e6a6f55e5a6d2f6a66b6c3e")
resp := sdk.Instance.CallContract("identity", "levelOf", args)
if string(resp.Payload) != "4" {
    return Error("address not registered")
}
```

## cmc使用示例

命令行工具使用示例
```sh
echo
echo "安装合约 identity，创建者为管理员"
./cmc client contract user create \
--contract-name=identity \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/standard_identity-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--sync-result=true \
--params="{}"
{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "identity",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674"
        ],
        "topic": "alterAdminAddress",
        "tx_id": "1837c5a098f3e17cca5a3816fdb1102e6235de625a2847df949822aa6608f29c"
      }
    ],
    "gas_used": 557803,
    "result": {
      "address": "ce2b56808e19806079a1275b1bae7002aa3ce973",
      "creator": {
        "member_id": "client1.sign.wx-org1.chainmaker.org",
        "member_info": "Y9PdAUogrWJ0Z1LsrlgbL7ltvZIi1+MNRPVgholFo9c=",
        "member_type": 1,
        "org_id": "wx-org1.chainmaker.org",
        "role": "CLIENT",
        "uid": "b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674"
      },
      "name": "identity",
      "runtime_type": 2,
      "version": "1.0"
    }
  },
  "tx_block_height": 2,
  "tx_id": "1837c5a098f3e17cca5a3816fdb1102e6235de625a2847df949822aa6608f29c",
  "tx_timestamp": 1745080674
}
```

```sh
echo
echo "执行合约 identity.SetIdentity(address, level, pkPem string)，设置身份"
./cmc client contract user invoke \
--contract-name=identity \
--method=SetIdentity \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"level\":\"4\",\"pkPem\":\"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\n-----END PUBLIC KEY-----\",\"metadata\":\"something\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "identity",
        "contract_version": "1.0",
        "event_data": [
          "b0831c4aebde4eb5a97b0eb5b1310a746763e752",
          "4",
          "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\n-----END PUBLIC KEY-----"
        ],
        "topic": "setIdentity",
        "tx_id": "1837c5aafd031282cab90da3d825cfab16b36d92ae554dafb71814409d858a7b"
      }
    ],
    "gas_used": 1515471,
    "result": "ok"
  },
  "tx_block_height": 3,
  "tx_id": "1837c5aafd031282cab90da3d825cfab16b36d92ae554dafb71814409d858a7b",
  "tx_timestamp": 1745080718
}
```

```sh
echo
echo "执行合约 identity.SetIdentityBatch(identities []standard.Identity)，设置身份"
./cmc client contract user invoke \
--contract-name=identity \
--method=SetIdentityBatch \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"identities\":\"[{\\\"address\\\":\\\"5fa92a33364dd5ce26a9814a6aceb240bd6bf082\\\",\\\"level\\\":4,\\\"pkPem\\\":\\\"-----BEGIN PUBLIC KEY-----\\\\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\\\\n-----END PUBLIC KEY-----\\\",\\\"metadata\\\":\\\"something\\\"},{\\\"address\\\":\\\"5fa92a33364dd5ce26a9814a6aceb240bd6bf081\\\",\\\"level\\\":3,\\\"pkPem\\\":\\\"-----BEGIN PUBLIC KEY-----\\\\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\\\\n-----END PUBLIC KEY-----\\\",\\\"metadata\\\":\\\"something\\\"}]\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "identity",
        "contract_version": "1.0",
        "event_data": [
          "5fa92a33364dd5ce26a9814a6aceb240bd6bf082",
          "4",
          "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\n-----END PUBLIC KEY-----"
        ],
        "topic": "setIdentity",
        "tx_id": "1837c5bd7a70f51bca2b069b7bed375cb7b8fecda87e427e81a54a39fc786997"
      },
      {
        "contract_name": "identity",
        "contract_version": "1.0",
        "event_data": [
          "5fa92a33364dd5ce26a9814a6aceb240bd6bf081",
          "3",
          "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\n-----END PUBLIC KEY-----"
        ],
        "topic": "setIdentity",
        "tx_id": "1837c5bd7a70f51bca2b069b7bed375cb7b8fecda87e427e81a54a39fc786997"
      }
    ],
    "gas_used": 2628472,
    "result": "ok"
  },
  "tx_block_height": 4,
  "tx_id": "1837c5bd7a70f51bca2b069b7bed375cb7b8fecda87e427e81a54a39fc786997",
  "tx_timestamp": 1745080798
}
```

```sh
echo
echo "执行合约 identity.AlterAdminAddress(adminAddress string)，设置管理员"
./cmc client contract user invoke \
--contract-name=identity \
--method=AlterAdminAddress \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"adminAddress\":\"e21b6661e9fa7d0056b6c3b0f97995f95ac4d540\"}" \
--sync-result=true \
--result-to-string=true

{
  "contract_result": {
    "contract_event": [
      {
        "contract_name": "identity",
        "contract_version": "1.0",
        "event_data": [
          "e21b6661e9fa7d0056b6c3b0f97995f95ac4d540"
        ],
        "topic": "alterAdminAddress",
        "tx_id": "1837c5c44dbc63d6caa607d50091da23c03785d65057407391e49fa1e852bc08"
      }
    ],
    "gas_used": 645840,
    "result": "ok"
  },
  "tx_block_height": 5,
  "tx_id": "1837c5c44dbc63d6caa607d50091da23c03785d65057407391e49fa1e852bc08",
  "tx_timestamp": 1745080827
}
```

```sh
echo
echo "执行合约 identity.LevelOf(address string)，获取身份类型编号"
./cmc client contract user get \
--contract-name=identity \
--method=LevelOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 350137,
    "result": "4"
  },
  "message": "SUCCESS",
  "tx_id": "1837c5cc39e8dd86ca50afe04c76551fe9444e5a1123444ca593d03c25fe68f6"
}
```

```sh
echo
echo "执行合约 identity.PkPemOf(address string)，获取pem公钥串"
./cmc client contract user get \
--contract-name=identity \
--method=PkPemOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 361533,
    "result": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\n-----END PUBLIC KEY-----"
  },
  "message": "SUCCESS",
  "tx_id": "1837c5d4a6219566ca70fed7a6103afec873c7e01ee64a44a85627b4c685d322"
}
```

```sh
echo
echo "执行合约 identity.IdentityOf(address string)，获取身份信息"
./cmc client contract user get \
--contract-name=identity \
--method=IdentityOf \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"address\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 700298,
    "result": "{\"address\":\"b0831c4aebde4eb5a97b0eb5b1310a746763e752\",\"pkPem\":\"-----BEGIN PUBLIC KEY-----\\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECr/ijK264TXbHfIvhJJz43z9hZLroyWUZY371pfQqaToo2by5ljMj3Ot8/XM2n5Xr/6xIwVLJ7t+C5cwcGRjqA==\\n-----END PUBLIC KEY-----\",\"level\":4,\"metadata\":\"something\"}"
  },
  "message": "SUCCESS",
  "tx_id": "1837c5e309283c94caf1d3d2d6ff51945449fdd40cd342309efdc095619eb5cb"
}
```

```sh
echo
echo "查询合约 identity.SupportStandard(standardName string)，是否支持合约标准CMID"
./cmc client contract user get \
--contract-name=identity \
--method=SupportStandard \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMID\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 187891,
    "result": "true"
  },
  "message": "SUCCESS",
  "tx_id": "1837c5e9970004f5cafaa3ee882a8991d5e71b78fd35459ea8b09bf9b2fc4a8d"
}
```

```sh
echo
echo "查询合约 identity.SupportStandard(standardName string)，是否支持合约标准CMBC"
./cmc client contract user get \
--contract-name=identity \
--method=SupportStandard \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMBC\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 186865,
    "result": "true"
  },
  "message": "SUCCESS",
  "tx_id": "1837c5f167c65f3dca2e5b02d603e043bf4ea51290e645208b5bc651f02c5608"
}
```

```sh
echo
echo "查询合约 identity.Standards()，支持的合约标准列表"
./cmc client contract user get \
--contract-name=identity \
--method=Standards \
--sdk-conf-path=./testdata/sdk_config.yml \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 254564,
    "result": "[\"CMID\",\"CMBC\"]"
  },
  "message": "SUCCESS",
  "tx_id": "1837c5f85f385a7dca0168626127f14b4f52f6a87009491aaae077512e683279"
}
```

```sh
echo
echo "查询合约 identity.Identities()，当前合约支持的身份类型列表"
./cmc client contract user get \
--contract-name=identity \
--method=Identities \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"standardName\":\"CMBC\"}" \
--result-to-string=true

{
  "contract_result": {
    "gas_used": 862382,
    "result": "[{\"level\":0,\"description\":\"未认证\"},{\"level\":2,\"description\":\"个人手机号注册用户\"},{\"level\":3,\"description\":\"个人实名用户\"},{\"level\":1,\"description\":\"企业的用户\"},{\"level\":4,\"description\":\"企业实名用户\"}]"
  },
  "message": "SUCCESS",
  "tx_id": "1837c5fe01c47152ca2611498af010f90a95b003af6349e5a40411a4c3c9ea98"
}
```
