# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

sudo chmod 777 /etc/rancher/k3s/k3s.yaml
sudo mkdir /opt/wetee-worker
sudo chmod -R 777 /opt/wetee-worker

export SIDE_CHAIN_PORT=10000
export GQL_PORT=10005
export CHAIN_ADDR=ws://127.0.0.1:9944
export KUBECONFIG=/etc/kube/config

echo $KUBE_CONFIG_PATH
make run