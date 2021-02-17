FROM alpine:3.13.1

WORKDIR /app
COPY blobfuse-proxy /app/blobfuse-proxy
CMD ["/app/blobfuse-proxy"]
