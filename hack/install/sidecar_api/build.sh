# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/

tag=`date "+%Y-%m-%d-%H_%M"`

docker build -t registry.cn-hangzhou.aliyuncs.com/wetee_dao/substrate_sidecar:$tag .
docker push registry.cn-hangzhou.aliyuncs.com/wetee_dao/substrate_sidecar:$tag