# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

tag=`date "+%Y-%m-%d-%H_%M"`

# 清理旧的部署
make undeploy

make manifests
make install

# 创建暂时无用的目录
sudo mkdir -p /var/run/secrets/kubernetes.io/serviceaccount/

# 编译程序
make build

# 构建镜像
make docker-build docker-push IMG=wetee/worker:$tag

# 部署镜像
make deploy IMG=wetee/worker:$tag

# 创建内部服务
kubectl create -f ./hack/manager_headless.yaml

# 创建外部服务
kubectl create -f ./hack/manager_nodeport.yaml