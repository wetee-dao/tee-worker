FROM wetee/ego-ubuntu:20.04
WORKDIR /
# COPY --from=builder /workspace/manager .
ADD bin/*  /

RUN mkdir -p /opt/wetee-worker
USER 65532:65532

CMD ["ego","run","/manager"]
