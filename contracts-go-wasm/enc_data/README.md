
## Test

### 部署合约
```sh
./cmc client contract user create \
--contract-name=encdata \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/enc_data-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"




```

### enc_data
```sh
./cmc client contract user invoke \
--contract-name=encdata \
--method=enc_data \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{
  \"data_key\": \"dataKey\",
  \"data_value\": \"dataValue\",
  \"enc_key\": \"encKey\",
  \"authorized_person\": \"-----BEGIN CERTIFICATE-----\nMIICeDCCAh6gAwIBAgIDDmp3MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjb21tb24xLDAqBgNVBAMTI2NvbW1vbjEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEn4ZMa251\nacwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEema3Zl8J33iXv9BNGyKH1/\n7p+yHYj2ougY2KNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCCsMh3Xbs+H\nqbb7iYyi3G2RhZG0+l8GmYPa/i7NSkIxcDArBgNVHSMEJDAigCDStB+0gbNWFT1p\niPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiAG3fYB1HEu\nGi7aUUNBIOizWBCtOuWWvmR5FMVSuuUYdAIhALqbClSD9Kt2gYwYucCE7iPajc3H\nwyi1e7ZVkH5vjHP8\n-----END CERTIFICATE-----\"
}"




```

### enc_auth
```sh

./cmc client contract user invoke \
--contract-name=encdata \
--method=enc_auth \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{
  \"data_key\": \"dataKey\",
  \"enc_key\": \"encKey\",
  \"authorized_person\": \"-----BEGIN CERTIFICATE-----\nMIICfjCCAiSgAwIBAgIDCgn6MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZcxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMRIwEAYDVQQLEwljb25zZW5zdXMxLzAtBgNVBAMTJmNvbnNlbnN1czEuc2lnbi53\neC1vcmcxLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nXJBsVjVS5zcdQk2RhdA7eRs1DXdVq8xXRCD8G9CQ+YoDp/3bWLTBj7nw2ZYQHdxq\nBp1iPP0tIbv4S/LAw1WbCqNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCB0\noajU1EwCPAWpcyBwnuaUUo98H4W75/0IyqmbvrXuEDArBgNVHSMEJDAigCDStB+0\ngbNWFT1piPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiEA\nzQIb4bTapNnTqbEyr0B2VahFunoFThRZrZG1PXSicTUCIBk3x7Z/PRR9Q/agNuJI\nNaH1gyFpD5XW1nlTQa4xdrML\n-----END CERTIFICATE-----\",
  \"authorizer\": \"-----BEGIN CERTIFICATE-----\nMIICeDCCAh6gAwIBAgIDDmp3MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjb21tb24xLDAqBgNVBAMTI2NvbW1vbjEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEn4ZMa251\nacwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEema3Zl8J33iXv9BNGyKH1/\n7p+yHYj2ougY2KNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCCsMh3Xbs+H\nqbb7iYyi3G2RhZG0+l8GmYPa/i7NSkIxcDArBgNVHSMEJDAigCDStB+0gbNWFT1p\niPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiAG3fYB1HEu\nGi7aUUNBIOizWBCtOuWWvmR5FMVSuuUYdAIhALqbClSD9Kt2gYwYucCE7iPajc3H\nwyi1e7ZVkH5vjHP8\n-----END CERTIFICATE-----\",
  \"auth_sign\": \"-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIK0M179niQ0F5+iZAjIWSa+frPiYGyrktwUKln/gGOCWoAoGCCqGSM49\nAwEHoUQDQgAEn4ZMa251acwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEem\na3Zl8J33iXv9BNGyKH1/7p+yHYj2ougY2A==\n-----END EC PRIVATE KEY-----\",
  \"auth_level\": 2
}"


```

### 查询get_enc_data

```sh
./cmc client contract user invoke \
--contract-name=encdata \
--method=get_enc_data \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{
  \"data_key\": \"dataKey\",
}"


```

### 查询get_enc_auth

```sh
./cmc client contract user invoke \
--contract-name=encdata \
--method=get_enc_data \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{
  \"data_key\": \"dataKey\",
  \"authorizer\": \"-----BEGIN CERTIFICATE-----\nMIICeDCCAh6gAwIBAgIDDmp3MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjb21tb24xLDAqBgNVBAMTI2NvbW1vbjEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEn4ZMa251\nacwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEema3Zl8J33iXv9BNGyKH1/\n7p+yHYj2ougY2KNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCCsMh3Xbs+H\nqbb7iYyi3G2RhZG0+l8GmYPa/i7NSkIxcDArBgNVHSMEJDAigCDStB+0gbNWFT1p\niPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiAG3fYB1HEu\nGi7aUUNBIOizWBCtOuWWvmR5FMVSuuUYdAIhALqbClSD9Kt2gYwYucCE7iPajc3H\nwyi1e7ZVkH5vjHP8\n-----END CERTIFICATE-----\"
}"



```

### 更新update_enc_auth
```sh

./cmc client contract user invoke \
--contract-name=erc721 \
--method=ownerOf \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params=--params="{
  \"data_key\": \"dataKey\",
  \"authorized_person\": \"-----BEGIN CERTIFICATE-----\nMIICfjCCAiSgAwIBAgIDCgn6MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZcxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMRIwEAYDVQQLEwljb25zZW5zdXMxLzAtBgNVBAMTJmNvbnNlbnN1czEuc2lnbi53\neC1vcmcxLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nXJBsVjVS5zcdQk2RhdA7eRs1DXdVq8xXRCD8G9CQ+YoDp/3bWLTBj7nw2ZYQHdxq\nBp1iPP0tIbv4S/LAw1WbCqNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCB0\noajU1EwCPAWpcyBwnuaUUo98H4W75/0IyqmbvrXuEDArBgNVHSMEJDAigCDStB+0\ngbNWFT1piPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiEA\nzQIb4bTapNnTqbEyr0B2VahFunoFThRZrZG1PXSicTUCIBk3x7Z/PRR9Q/agNuJI\nNaH1gyFpD5XW1nlTQa4xdrML\n-----END CERTIFICATE-----\",
  \"authorizer\": \"-----BEGIN CERTIFICATE-----\nMIICeDCCAh6gAwIBAgIDDmp3MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjb21tb24xLDAqBgNVBAMTI2NvbW1vbjEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEn4ZMa251\nacwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEema3Zl8J33iXv9BNGyKH1/\n7p+yHYj2ougY2KNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCCsMh3Xbs+H\nqbb7iYyi3G2RhZG0+l8GmYPa/i7NSkIxcDArBgNVHSMEJDAigCDStB+0gbNWFT1p\niPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiAG3fYB1HEu\nGi7aUUNBIOizWBCtOuWWvmR5FMVSuuUYdAIhALqbClSD9Kt2gYwYucCE7iPajc3H\nwyi1e7ZVkH5vjHP8\n-----END CERTIFICATE-----\",
  \"auth_sign\": \"-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIK0M179niQ0F5+iZAjIWSa+frPiYGyrktwUKln/gGOCWoAoGCCqGSM49\nAwEHoUQDQgAEn4ZMa251acwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEem\na3Zl8J33iXv9BNGyKH1/7p+yHYj2ougY2A==\n-----END EC PRIVATE KEY-----\",
  \"auth_level\": \"2\"
}"

{
  "code": 4,
  "contract_result": {
    "code": 1,
    "gas_used": 2493988,
    "message": "contract message:get owner failed",
    "result": ""
  },
  "tx_block_height": 4,
  "tx_id": "1834a288b1539ac2caaba6485172b2ddb53d2d0eda834acf88aaa373ffbd977e",
  "tx_timestamp": 1744197663
}
```

