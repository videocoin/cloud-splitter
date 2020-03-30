FROM golang:1.14 as builder
WORKDIR /go/src/github.com/videocoin/cloud-splitter
COPY . .
RUN make build

FROM jrottenberg/ffmpeg:4.1-ubuntu
RUN apt-get update && \
    apt-get -y --force-yes install \
        mediainfo \
        libmediainfo-dev \
        ffmpeg \
        curl

COPY --from=builder /go/src/github.com/videocoin/cloud-splitter/bin/splitter /opt/videocoin/bin/splitter

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
   curl -L -k https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 --output /bin/grpc_health_probe && chmod +x /bin/grpc_health_probe

ENTRYPOINT ["/opt/videocoin/bin/splitter"]
CMD [""]
