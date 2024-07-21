FROM registry.cn-hangzhou.aliyuncs.com/wetee_dao/ego-ubuntu-deploy:22.04
WORKDIR /

ADD bin/*  /
ADD bin/keys /keys

RUN mkdir -p /opt/wetee-worker

EXPOSE 8880 8883 

CMD ["/bin/sh", "-c" ,"ego sign manager && ego run manager"]