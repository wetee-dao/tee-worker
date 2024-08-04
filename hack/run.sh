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
sudo chmod 777 /opt/wetee-worker

echo $KUBE_CONFIG_PATH
make run