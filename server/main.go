package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	mount_azure_blob "github.com/boddumanohar/blobfuse-proxy/pb"
	"google.golang.org/grpc"
)

type MountServer struct {
	mount_azure_blob.UnimplementedMountServiceServer
}

// NewMountServer returns a new Mountserver
func NewMountServiceServer() *MountServer {
	return &MountServer{}
}

func (server *MountServer) mustEmbedUnimplementedMountServiceServer() {}

// MountAzureBlob mounts an azure blob container to given location
func (server *MountServer) MountAzureBlob(ctx context.Context,
	req *mount_azure_blob.MountAzureBlobRequest,
) (resp *mount_azure_blob.MountAzureBlobResponse, err error) {

	args := req.GetMountArgs()
	log.Printf("received mount request: Mounting with args %v \n", args)
	cmd := exec.Command("blobfuse", strings.Split(args, " ")...)
	var finalRes mount_azure_blob.MountAzureBlobResponse
	cmd.Env = append(os.Environ(), "AZURE_STORAGE_ACCOUNT="+req.GetAccountName())
	cmd.Env = append(cmd.Env, "AZURE_STORAGE_ACCESS_KEY="+req.GetAccountKey())
	cmd.Env = append(cmd.Env, req.GetAuthEnv()...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("blobfuse mount failed with error", err)
	} else {
		log.Println("successfully mounted")
	}
	finalRes.Output = string(output)
	log.Println(finalRes.Output)
	return &finalRes, err
}

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
	endpoint := flag.String("endpoint", "unix://tmp/blobfuse-proxy.sock", "CSI endpoint")
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

	mountServer := NewMountServiceServer()

	log.Printf("Listening for connections on address: %#v\n", listener.Addr())
	if err = runGRPCServer(mountServer, false, listener); err != nil {
		log.Fatalf("Listening for connections on address: %#v, error: %v", listener.Addr(), err)
	}
}
