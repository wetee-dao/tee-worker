# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

# export KUBERNETES_SERVICE_HOST=0.0.0.0
# export KUBERNETES_SERVICE_PORT=6443
# export KUBE_CONFIG_PATH=/etc/rancher/k3s/k3s.yaml


echo $KUBE_CONFIG_PATH
make run