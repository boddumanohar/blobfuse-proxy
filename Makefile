gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb

clean:
	rm pb/*.go

build:
	go build

tidy:
	go mod tidy
