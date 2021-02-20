gen:
	protoc --proto_path=proto --go_out=pb --go-grpc_out=pb proto/*.proto

clean:
	rm pb/*.go

build:
	go build

tidy:
	go mod tidy
