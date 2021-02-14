package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	mount_azure_blob "github.com/boddumanohar/blobfuse-proxy/pb"
	"github.com/boddumanohar/blobfuse-proxy/server"
	"google.golang.org/grpc"
)

func runGRPCServer(
	mountServer mount_azure_blob.MountServiceServer,
	enableTLS bool,
	listener net.Listener,
) error {
	serverOptions := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(serverOptions...)

	mount_azure_blob.RegisterMountServiceServer(grpcServer, mountServer)

	log.Printf("Start GRPC server at %s, TLS = %t", listener.Addr().String(), enableTLS)
	return grpcServer.Serve(listener)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	mountServer := server.NewMountServiceServer()
	err = runGRPCServer(mountServer, false, listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
