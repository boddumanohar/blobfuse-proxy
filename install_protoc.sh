PROTOCOL_BUF_VERSION=v3.15.1
PROTOC_GEN_GO_GRPC_VERSION=v1.1.0
PROTOC_GEN_GO_VERSION=v1.25.0
wget https://github.com/protocolbuffers/protobuf/releases/download/$PROTOCOL_BUF_VERSION/protoc-$PROTOCOL_BUF_VERSION-linux-x86_64.zip
wget https://github.com/grpc/grpc-go/releases/download/cmd/protoc-gen-go-grpc/$PROTOGEN_GO_GRPC_VERSION/protoc-gen-go-grpc.$PROTOGEN_GO_GRPC_VERSION.linux.amd64.tar.gz
wget wget https://github.com/protocolbuffers/protobuf-go/releases/download/$PROTOC_GEN_GO_VERSION/protoc-gen-go.$PROTOC_GEN_GO_VERSION.linux.amd64.tar.gz
sudo unzip -o protoc-3.15.1-linux-x86_64.zip -d /usr/local bin/protoc
sudo unzip -o protoc-3.15.1-linux-x86_64.zip -d /usr/local 'include/*'
sudo tar -C /usr/local/bin -xvf protoc-gen-go-grpc.$PROTOGEN_GO_GRPC_VERSION.linux.amd64.tar.gz
sudo tar -C /usr/local/bin -xvf protoc-gen-go.$PROTOC_GEN_GO_VERSION.linux.amd64.tar.gz
sudo chmod +x /usr/local/bin/protoc-gen-go
sudo chmod +x /usr/local/bin/protoc-gen-go-grpc
