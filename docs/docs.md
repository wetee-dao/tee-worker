### 文档
 - kubebuilder  https://cloudnative.to/kubebuilder/cronjob-tutorial/new-api.html
 - 远程证明 https://github.com/confidential-containers/attestation-service


| Project | Technical Solutions | Programming Languages | Execution Hardware | Supported Programming Languages |
| ----- | ----------- | ------------- | ------------- | ------------- |
| WeTEE | • Kubernetes and Docker as computing cluster solutions <br>• Gramine or Ego as confidential container solutions<br>• Worker with TEE（run as K8s Operator） provides sgx remote attestation, key management, program Confidentiality Injection and uploads sgx remote attestation as part of proof of work to blockchain network consensus.| golang | x86 Server with Intel SGX | C, Python, Go, Rust, Javascript and all program codes that support Gramine | • Web2 service deployment in Web2.5 solutions<br>• Decentralized team's Web2 aplication deployment platform<br>• Deploy website applications<br>• Deploy server-side APIs |
| Integritee | Woker with TEE provides Sidechain、Off-Chain Worker、Oracle as confidential solutions | rust | x86 Server with Intel SGX | Rust（Fork integritee worker and add own Rust code）| • use as Sidechain<br>• use as Off-Chain Worker<br>• use as Oracle |
| Acurast | Acurast  Zero-Trust Execution Layer as End-to-End Zero Trust Job Execution as confidential solutions | rust | mobile device | Javascript | • job to get access to the legacy Web2 world through this module and bring off-chain data of public or permissioned APIs and off-chain computation to their Web3 application<br>• job to building block allowing for general bidirectional message passing between two blockchain networks.<br>• job to built in mind to perform learning tasks of language models on mobile devices