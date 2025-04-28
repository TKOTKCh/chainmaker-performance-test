
## Test

### 部署合约
```sh
./cmc client contract user create \
--contract-name=encdata \
--runtime-type=WASMER \
--byte-code-path=./testdata/go-wasm-demo/enc_data_modify-go.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"

{
  "contract_result": {
    "gas_used": 238938,
    "result": {
      "address": "b8ee7eddfd9de32269a0c228aab8920c9cc16bb4",
      "creator": {
        "member_id": "client1.sign.wx-org1.chainmaker.org",
        "member_info": "Y9PdAUogrWJ0Z1LsrlgbL7ltvZIi1+MNRPVgholFo9c=",
        "member_type": 1,
        "org_id": "wx-org1.chainmaker.org",
        "role": "CLIENT",
        "uid": "b0831c4aebde4eb5a97b0eb5b1310a746763e75225ea6021d7fdf1885b9a8674"
      },
      "name": "encdata",
      "runtime_type": 2,
      "version": "1.0"
    }
  },
  "tx_block_height": 2,
  "tx_id": "183a85650741c9e7ca1da63e4909eb9628efa072f3ac47148efb09329b89d58e",
  "tx_timestamp": 1745854474
}


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

{
  "contract_result": {
    "gas_used": 1773567,
    "result": "authorizedPerson -----BEGIN CERTIFICATE-----\nMIICeDCCAh6gAwIBAgIDDmp3MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjb21tb24xLDAqBgNVBAMTI2NvbW1vbjEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEn4ZMa251\nacwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEema3Zl8J33iXv9BNGyKH1/\n7p+yHYj2ougY2KNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCCsMh3Xbs+H\nqbb7iYyi3G2RhZG0+l8GmYPa/i7NSkIxcDArBgNVHSMEJDAigCDStB+0gbNWFT1p\niPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiAG3fYB1HEu\nGi7aUUNBIOizWBCtOuWWvmR5FMVSuuUYdAIhALqbClSD9Kt2gYwYucCE7iPajc3H\nwyi1e7ZVkH5vjHP8\n-----END CERTIFICATE-----,dataKey dataKey,encAuth {\"dataKey\":\"ZGF0YUtleQ==\",\"authorizedPerson\":\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNlRENDQWg2Z0F3SUJBZ0lERG1wM01Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTVM1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3hMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEkxTURReE9ERTFORFF5T1ZvWERUTXcKTURReE56RTFORFF5T1Zvd2daRXhDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jeExtTm9ZV2x1YldGclpYSXViM0puCk1ROHdEUVlEVlFRTEV3WmpiMjF0YjI0eExEQXFCZ05WQkFNVEkyTnZiVzF2YmpFdWMybG5iaTUzZUMxdmNtY3gKTG1Ob1lXbHViV0ZyWlhJdWIzSm5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUVuNFpNYTI1MQphY3daa21aUS9IQldHeXkxaE1yNDBDaEhKMjlhTnZsQ3A5eFVCamwzU0VlbWEzWmw4SjMzaVh2OUJOR3lLSDEvCjdwK3lIWWoyb3VnWTJLTnFNR2d3RGdZRFZSMFBBUUgvQkFRREFnYkFNQ2tHQTFVZERnUWlCQ0NzTWgzWGJzK0gKcWJiN2lZeWkzRzJSaFpHMCtsOEdtWVBhL2k3TlNrSXhjREFyQmdOVkhTTUVKREFpZ0NEU3RCKzBnYk5XRlQxcAppUFc4K1h6Sit2UzBtM0paMWdLWVNVRVN0N24vcHpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBRzNmWUIxSEV1CkdpN2FVVU5CSU9peldCQ3RPdVdXdm1SNUZNVlN1dVVZZEFJaEFMcWJDbFNEOUt0MmdZd1l1Y0NFN2lQYWpjM0gKd3lpMWU3WlZrSDV2akhQOAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t\",\"encAESKey\":\"ZW5jS2V5\",\"authorizer\":\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNlRENDQWg2Z0F3SUJBZ0lERG1wM01Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTVM1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3hMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEkxTURReE9ERTFORFF5T1ZvWERUTXcKTURReE56RTFORFF5T1Zvd2daRXhDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jeExtTm9ZV2x1YldGclpYSXViM0puCk1ROHdEUVlEVlFRTEV3WmpiMjF0YjI0eExEQXFCZ05WQkFNVEkyTnZiVzF2YmpFdWMybG5iaTUzZUMxdmNtY3gKTG1Ob1lXbHViV0ZyWlhJdWIzSm5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUVuNFpNYTI1MQphY3daa21aUS9IQldHeXkxaE1yNDBDaEhKMjlhTnZsQ3A5eFVCamwzU0VlbWEzWmw4SjMzaVh2OUJOR3lLSDEvCjdwK3lIWWoyb3VnWTJLTnFNR2d3RGdZRFZSMFBBUUgvQkFRREFnYkFNQ2tHQTFVZERnUWlCQ0NzTWgzWGJzK0gKcWJiN2lZeWkzRzJSaFpHMCtsOEdtWVBhL2k3TlNrSXhjREFyQmdOVkhTTUVKREFpZ0NEU3RCKzBnYk5XRlQxcAppUFc4K1h6Sit2UzBtM0paMWdLWVNVRVN0N24vcHpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBRzNmWUIxSEV1CkdpN2FVVU5CSU9peldCQ3RPdVdXdm1SNUZNVlN1dVVZZEFJaEFMcWJDbFNEOUt0MmdZd1l1Y0NFN2lQYWpjM0gKd3lpMWU3WlZrSDV2akhQOAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t\",\"authSignature\":null,\"authLevel\":1}"
  },
  "tx_block_height": 3,
  "tx_id": "183a85681fe35a2aca169e22ca8129ef22d7a35e583c400486fefe86ee42836c",
  "tx_timestamp": 1745854487
}


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

{
  "contract_result": {
    "gas_used": 4221410,
    "result": "store auth info successfully"
  },
  "tx_block_height": 4,
  "tx_id": "183a856c963f5846caa25cd4d89d05212faa7a9794ae413f97b76320cc828750",
  "tx_timestamp": 1745854506
}
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
  \"data_key\": \"dataKey\"
}"

{
  "contract_result": {
    "gas_used": 338456,
    "result": "dataValue"
  },
  "tx_block_height": 5,
  "tx_id": "183a8592a9acbcf4cacba58c91df75ec37ea941e4c244a31a49a32b755e63980",
  "tx_timestamp": 1745854670
}

```

### 查询get_enc_auth

```sh
./cmc client contract user invoke \
--contract-name=encdata \
--method=get_enc_auth \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{
  \"data_key\": \"dataKey\",
  \"authorizer\": \"-----BEGIN CERTIFICATE-----\nMIICeDCCAh6gAwIBAgIDDmp3MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjb21tb24xLDAqBgNVBAMTI2NvbW1vbjEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEn4ZMa251\nacwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEema3Zl8J33iXv9BNGyKH1/\n7p+yHYj2ougY2KNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCCsMh3Xbs+H\nqbb7iYyi3G2RhZG0+l8GmYPa/i7NSkIxcDArBgNVHSMEJDAigCDStB+0gbNWFT1p\niPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiAG3fYB1HEu\nGi7aUUNBIOizWBCtOuWWvmR5FMVSuuUYdAIhALqbClSD9Kt2gYwYucCE7iPajc3H\nwyi1e7ZVkH5vjHP8\n-----END CERTIFICATE-----\"
}"

{
  "contract_result": {
    "gas_used": 2868959,
    "result": "encKey"
  },
  "tx_block_height": 7,
  "tx_id": "183a85a8cf4a68b9ca97236f23dbfcad516c5f0a0ce545d58f58421e22ed6e42",
  "tx_timestamp": 1745854765
}

```

### 更新update_enc_auth
```sh

./cmc client contract user invoke \
--contract-name=encdata \
--method=update_enc_auth \
--sync-result=true \
--result-to-string=true \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{
  \"data_key\": \"dataKey\",
  \"authorized_person\": \"-----BEGIN CERTIFICATE-----\nMIICfjCCAiSgAwIBAgIDCgn6MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZcxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMRIwEAYDVQQLEwljb25zZW5zdXMxLzAtBgNVBAMTJmNvbnNlbnN1czEuc2lnbi53\neC1vcmcxLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\nXJBsVjVS5zcdQk2RhdA7eRs1DXdVq8xXRCD8G9CQ+YoDp/3bWLTBj7nw2ZYQHdxq\nBp1iPP0tIbv4S/LAw1WbCqNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCB0\noajU1EwCPAWpcyBwnuaUUo98H4W75/0IyqmbvrXuEDArBgNVHSMEJDAigCDStB+0\ngbNWFT1piPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiEA\nzQIb4bTapNnTqbEyr0B2VahFunoFThRZrZG1PXSicTUCIBk3x7Z/PRR9Q/agNuJI\nNaH1gyFpD5XW1nlTQa4xdrML\n-----END CERTIFICATE-----\",
  \"authorizer\": \"-----BEGIN CERTIFICATE-----\nMIICeDCCAh6gAwIBAgIDDmp3MAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI1MDQxODE1NDQyOVoXDTMw\nMDQxNzE1NDQyOVowgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjb21tb24xLDAqBgNVBAMTI2NvbW1vbjEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEn4ZMa251\nacwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEema3Zl8J33iXv9BNGyKH1/\n7p+yHYj2ougY2KNqMGgwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCCsMh3Xbs+H\nqbb7iYyi3G2RhZG0+l8GmYPa/i7NSkIxcDArBgNVHSMEJDAigCDStB+0gbNWFT1p\niPW8+XzJ+vS0m3JZ1gKYSUESt7n/pzAKBggqhkjOPQQDAgNIADBFAiAG3fYB1HEu\nGi7aUUNBIOizWBCtOuWWvmR5FMVSuuUYdAIhALqbClSD9Kt2gYwYucCE7iPajc3H\nwyi1e7ZVkH5vjHP8\n-----END CERTIFICATE-----\",
  \"auth_sign\": \"-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIK0M179niQ0F5+iZAjIWSa+frPiYGyrktwUKln/gGOCWoAoGCCqGSM49\nAwEHoUQDQgAEn4ZMa251acwZkmZQ/HBWGyy1hMr40ChHJ29aNvlCp9xUBjl3SEem\na3Zl8J33iXv9BNGyKH1/7p+yHYj2ougY2A==\n-----END EC PRIVATE KEY-----\",
  \"auth_level\": 2
}"

{
  "contract_result": {
    "gas_used": 5258284,
    "result": "store auth info successfully"
  },
  "tx_block_height": 8,
  "tx_id": "183a85c1047eca3dca5ba23c8029cbb7f8631a89b05c4e919efe224fedfc016a",
  "tx_timestamp": 1745854869
}
```

