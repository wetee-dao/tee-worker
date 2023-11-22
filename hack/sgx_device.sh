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
      - image: registry.cn-hangzhou.aliyuncs.com/acs/sgx-device-plugin:v1.0.0-fb467e2-aliyun
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