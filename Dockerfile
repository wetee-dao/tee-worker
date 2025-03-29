FROM wetee/ego-ubuntu-24-04:1.7.0
WORKDIR /

ADD bin/manager  /

RUN mkdir -p /opt/wetee-worker

EXPOSE 8880 8883 

CMD ["/bin/sh", "-c" ,"ego run manager"]