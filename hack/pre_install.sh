# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

cat <<EOF | kubectl create -f -
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: sgx-device-plugin-ds
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: sgx-device-plugin
  template:
    metadata:
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ""
      labels:
        k8s-app: sgx-device-plugin
    spec:
      containers:
      - image: registry.cn-hangzhou.aliyuncs.com/acs/sgx-device-plugin:v1.1.0-bb1f5f9-aliyun
        imagePullPolicy: IfNotPresent
        name: sgx-device-plugin
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
        volumeMounts:
        - mountPath: /var/lib/kubelet/device-plugins
          name: device-plugin
        - mountPath: /dev
          name: dev
      tolerations:
      - effect: NoSchedule
        key: alibabacloud.com/sgx_epc_MiB
        operator: Exists
      volumes:
      - hostPath:
          path: /var/lib/kubelet/device-plugins
          type: DirectoryOrCreate
        name: device-plugin
      - hostPath:
          path: /dev
          type: Directory
        name: dev
EOF

cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Namespace
metadata:
  name: worker-addon
  annotations:
    field.cattle.io/containerDefaultResourceLimit: '{}'
  labels:
    {}
EOF


# 为wetee-worker赋予集群管理权限
kubectl create clusterrolebinding wetee-admin --clusterrole=cluster-admin --user=system:serviceaccount:worker-system:worker-controller-manager

# docker run -e APIKEY=8d375b775de84c2593e7694dedbbbe90 -p 8081:8081 --name pccs -d --restart=unless-stopped  wetee/pccs