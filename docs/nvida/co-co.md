```bash
sudo modprobe br_netfilter
run as root $ echo 1 > /proc/sys/net/ipv4/ip_forward
sudo kubeadm init --apiserver-advertise-address=192.168.111.109 --pod-network-cidr=10.244.0.0/16

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config


kubectl label node wetee node-role.kubernetes.io/worker=m
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
kubectl apply -f ./co-co/canal.yaml

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

```
  Normal   Scheduled                 10m                   default-scheduler  Successfully assigned default/cuda-vectoradd-kata to wetee
  Warning  FailedCreatePodSandBox    8m19s (x20 over 10m)  kubelet            Failed to create pod sandbox: rpc error: code = Unknown desc = failed to create containerd task: failed to create shim task: QMP command failed: vfio 0000:06:00.0: group 5 is not viable: unknown
  Warning  UnexpectedAdmissionError  23s                   kubelet            Allocate failed due to no healthy devices present; cannot allocate unhealthy devices nvidia.com/GP104_GEFORCE_GTX_1070_TI, which is unexpected
  Warning  FailedMount               22s (x2 over 23s)     kubelet            MountVolume.SetUp failed for volume "kube-api-access-bckdn" : object "default"/"kube-root-ca.crt" not registered
```
