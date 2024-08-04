### 文档

- kubebuilder https://cloudnative.to/kubebuilder/cronjob-tutorial/new-api.html
- 远程证明 https://github.com/confidential-containers/attestation-service
- 使用 DKG 管理项目的去中心化密钥，用于加密数据的传输和存储项目的部署和更新

### ubuntu 22.04 install

```bash
### error
[get_driver_type /home/sgx/jenkins/ubuntuServer2204-release-build-trunk-223/build_target/PROD/label/Builder-UbuntuSrv2204/label_exp/ubuntu64/linux-trunk-opensource/psw/urts/linux/edmm_utility.cpp:116] Failed to open Intel SGX device.
ERROR: enclave_create with ENCLAVE_TYPE_SGX1 type failed (err=0x1) (oe_result_t=OE_PLATFORM_ERROR) [openenclave-src/host/sgx/sgxload.c:oe_sgx_create_enclave:454]
ERROR: oe_create_enclave failed. (Set OE_SIMULATION=1 for simulation mode.) [src/tools/erthost/erthost.cpp:main:265]
ERROR: failed to open Intel SGX device

### 解决办法
sudo apt install linux-base-sgx
sudo udevadm trigger
```
