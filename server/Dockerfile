FROM golang:1.15.8-alpine as builder
WORKDIR /go/src/sigs.k8s.io/blob-csi-driver
ADD . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -mod vendor -o _output/blobfuse-proxy ./server

FROM k8s.gcr.io/build-image/debian-base-amd64:v2.1.3
COPY --from=builder /go/src/sigs.k8s.io/blob-csi-driver/_output/blobfuse-proxy /blobfuse-proxy
RUN apt-get update \
    && clean-install ca-certificates pkg-config libfuse-dev cmake libcurl4-gnutls-dev libgnutls28-dev uuid-dev libgcrypt20-dev wget

RUN wget -O /tmp/packages-microsoft-prod.deb https://packages.microsoft.com/config/ubuntu/16.04/packages-microsoft-prod.deb \
    && dpkg -i /tmp/packages-microsoft-prod.deb \
    && apt-get update \
    && apt install blobfuse fuse -y \
    && rm -f /tmp/packages-microsoft-prod.deb \
    && apt remove wget -y

ENTRYPOINT ["/blobfuse-proxy"]
