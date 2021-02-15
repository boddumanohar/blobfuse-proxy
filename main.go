package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

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

func parseEndpoint(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") || strings.HasPrefix(strings.ToLower(ep), "tcp://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
	}
	return "", "", fmt.Errorf("Invalid endpoint: %v", ep)
}

func main() {
	endpoint := flag.String("endpoint", "unix://tmp/blobfuseproxy.sock", "CSI endpoint")
	flag.Parse()
	proto, addr, err := parseEndpoint(*endpoint)
	if err != nil {
		log.Fatal(err.Error())
	}

	if proto == "unix" {
		addr = "/" + addr
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			log.Fatalf("Failed to remove %s, error: %s", addr, err.Error())
		}
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	mountServer := server.NewMountServiceServer()

	log.Printf("Listening for connections on address: %#v\n", listener.Addr())
	if err = runGRPCServer(mountServer, false, listener); err != nil {
		log.Fatalf("Listening for connections on address: %#v, error: %v", listener.Addr(), err)
	}
}
