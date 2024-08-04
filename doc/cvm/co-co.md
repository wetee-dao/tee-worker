```bash
sudo modprobe br_netfilter
run as root $ echo 1 > /proc/sys/net/ipv4/ip_forward

containerd config default | sudo tee /etc/containerd/config.toml
如果容器不断重启，按照这个页面调整 https://blog.csdn.net/roadtohacker/article/details/134654399

# init images
sudo kubeadm config images pull --image-repository registry.aliyuncs.com/google_containers

# pause
sudo ctr -n k8s.io images pull -k registry.aliyuncs.com/google_containers/pause:3.8
sudo ctr -n k8s.io images tag registry.aliyuncs.com/google_containers/pause:3.8 registry.k8s.io/pause:3.8

# metrics
sudo ctr -n k8s.io images pull -k registry.aliyuncs.com/google_containers/metrics-server:v0.7.1
sudo ctr -n k8s.io images tag registry.aliyuncs.com/google_containers/metrics-server:v0.7.1 registry.k8s.io/metrics-server:v0.7.1

# kube-rbac-proxy:v0.14.1 错误
sudo ctr -n k8s.io images pull -k docker.io/wetee/kube-rbac-proxy:v0.14.1
sudo ctr -n k8s.io images tag docker.io/wetee/kube-rbac-proxy:v0.14.1 gcr.io/kubebuilder/kube-rbac-proxy:v0.14.1

# init
sudo kubeadm init --apiserver-advertise-address=192.168.111.121 --pod-network-cidr=10.244.0.0/16  --image-repository registry.aliyuncs.com/google_containers

# 命令环境设置
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

kubectl label node sev1 node.kubernetes.io/worker=
kubectl taint nodes --all node-role.kubernetes.io/control-plane-

# 网络插件
kubectl apply -f ./co-co/canal.yaml

# 监控
wget https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
vim ./components.yaml ===> registry.k8s.io/metrics-server:v0.7.1 to registry.aliyuncs.com/google_containers/metrics-server:v0.7.1
kubectl apply -f components.yaml

# 显卡部分
nvidia-device-plugin启动错误
 -  nvidia-container-cli: initialization error: nvml error: driver not loade
    ```
    # 验证具体的错误
    nvidia-container-cli -k -d /dev/tty info
    ```

export VERSION=v0.7.0
kubectl label node wetee nvidia.com/gpu.workload.config=vm-passthrough
kubectl apply -k "github.com/confidential-containers/operator/config/release?ref=${VERSION}"
# kubectl apply --dry-run=client -o yaml \
#     -k "github.com/confidential-containers/operator/config/samples/ccruntime/default?ref=${VERSION}" > ./co-co/ccruntime.yaml
kubectl apply -f ./co-co/ccruntime.yaml

helm repo add nvidia https://helm.ngc.nvidia.com/nvidia \
   && helm repo update
helm install --wait --generate-name \
   -n gpu-operator --create-namespace \
   nvidia/gpu-operator \
   --set sandboxWorkloads.enabled=true \
   --set kataManager.enabled=true

sudo kubeadm join 192.168.111.109:6443 --token lbv3j8.ohs18v5bhkqs7zdm \
    --discovery-token-ca-cert-hash sha256:cb36bc948497bdc9a79a0553f91fc26ae8066d4a1851d0aea88df63577c4757f
```

#### 查看驱动是否已经变成 vfio-pci
```bash
lspci -nnk -d 10de:
```

#### 检查 dmesg, 查询是否有vfio-pci驱动启动
```bash
sudo dmesg | grep vfio
```

#### 检查 dmesg, 查询是否有vfio-pci驱动启动
```bash
sudo dmesg | grep vfio
```

##### 问题1
``` bash
Normal   Scheduled                 10m                   default-scheduler  Successfully assigned default/cuda-vectoradd-kata to wetee
Warning  FailedCreatePodSandBox    8m19s (x20 over 10m)  kubelet            Failed to create pod sandbox: rpc error: code = Unknown desc = failed to create containerd task: failed to create shim task: QMP command failed: vfio 0000:06:00.0: group 5 is not viable: unknown
Warning  UnexpectedAdmissionError  23s                   kubelet            Allocate failed due to no healthy devices present; cannot allocate unhealthy devices nvidia.com/GP104_GEFORCE_GTX_1070_TI, which is unexpected
Warning  FailedMount               22s (x2 over 23s)     kubelet            MountVolume.SetUp failed for volume "kube-api-access-bckdn" : object "default"/"kube-root-ca.crt" not registered
```
  Normal   Scheduled                 10m                   default-scheduler  Successfully assigned default/cuda-vectoradd-kata to wetee
  Warning  FailedCreatePodSandBox    8m19s (x20 over 10m)  kubelet            Failed to create pod sandbox: rpc error: code = Unknown desc = failed to create containerd task: failed to create shim task: QMP command failed: vfio 0000:06:00.0: group 5 is not viable: unknown
  Warning  UnexpectedAdmissionError  23s                   kubelet            Allocate failed due to no healthy devices present; cannot allocate unhealthy devices nvidia.com/GP104_GEFORCE_GTX_1070_TI, which is unexpected
  Warning  FailedMount               22s (x2 over 23s)     kubelet            MountVolume.SetUp failed for volume "kube-api-access-bckdn" : object "default"/"kube-root-ca.crt" not registered
```
