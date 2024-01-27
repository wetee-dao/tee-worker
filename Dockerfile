FROM wetee/ego-ubuntu:20.04
WORKDIR /
# COPY --from=builder /workspace/manager .
ADD bin/*  /

RUN mkdir -p /opt/wetee-worker

EXPOSE 8880 8883 
CMD ["ego","run","/manager"]
