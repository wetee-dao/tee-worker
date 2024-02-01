## Run in Docker

### 1. Hardware and Software Requirements
- [CPU List - click to see cpu list](https://ark.intel.com/content/www/us/en/ark/search/featurefilter.html?productType=873&2_SoftwareGuardExtensions=Yes%20with%20Intel%C2%AE%20ME)
    - Intel 7th generation (Kaby Lake) Core i3, i5, i7, and i9 processors
    - Intel 8th generation (Cannon Lake) Core i3, i5, i7, and i9 processors
    - Intel 9th generation (Cascade Lake) Core i3, i5, i7, and i9 processors
    - Intel 10th generation (Comet Lake) Core i3, i5, i7, and i9 processors
    - 2nd Generation Xeon Scalable processors (Cascade Lake) and later generations generally provide SGX capabilities.
- OS ubuntu 20.04 or ubuntu 22.04 (not in docker)

### 2. install [Docker](https://docs.docker.com/get-docker/)
### 3. install Intel Sgx on Ubuntu 20.04/Ubuntu 22.04 and Ego Setup
> For more information about Ego, please refer to https://docs.edgeless.systems/ego/getting-started/install
```bash
sudo apt install build-essential libssl-dev

sudo mkdir -p /etc/apt/keyrings
wget -qO- https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | sudo tee /etc/apt/keyrings/intel-sgx-keyring.asc > /dev/null
echo "deb [signed-by=/etc/apt/keyrings/intel-sgx-keyring.asc arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/intel-sgx.list
sudo apt update

EGO_DEB=ego_1.4.1_amd64_ubuntu-$(lsb_release -rs).deb
wget https://github.com/edgelesssys/ego/releases/download/v1.4.1/$EGO_DEB
sudo apt install ./$EGO_DEB build-essential libssl-dev

sudo mkdir /opt/wetee-worker
sudo chmod 777 /opt/wetee-worker
```

Then run the following command to start a single node development chain.  

A few useful ones are as follow:  

```bash
# Use Docker to build (ego build must be run in sgx)
docker run --device /dev/sgx/enclave --device /dev/sgx/provision \
    -v ${PWD}:/srv wetee/ego-ubuntu:20.04 \
    bash -c "cd /srv && ego-go build -o ./bin/manager ./cmd/main.go"

# Build wetee-worker image
docker build -t wetee/worker:dev .

# Run worker in docker
docker run --device /dev/sgx/enclave --device /dev/sgx/provision \
     --network host \
     -e KUBECONFIG=/etc/kube/config \
     -v /etc/rancher/k3s:/etc/rancher/k3s \
     wetee/worker:dev
```