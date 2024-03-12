# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

# 部署主服务
make deploy IMG=wetee/worker:2024-03-12-09_09

# 创建内部服务
kubectl create -f ./hack/manager_headless.yaml

# 创建外部服务
kubectl create -f ./hack/manager_nodeport.yaml

# 创建外部服务
# kubectl create -f ./hack/manager_for_localdev.yaml