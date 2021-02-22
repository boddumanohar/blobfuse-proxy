gen:
	protoc --proto_path=proto --go_out=pb --go-grpc_out=pb proto/*.proto

clean:
	rm pb/*.go

server:
	go build -mod vendor -o _output/blobfuse-proxy ./server

tidy:
	go mod tidy
