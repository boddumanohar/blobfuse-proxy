package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	mount_azure_blob "github.com/boddumanohar/blobfuse-proxy/pb"
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

	log.Printf("received request: Mounting the container %s to the path %s \n", req.GetContainerName(), req.GetTargetPath())
	resp = &mount_azure_blob.MountAzureBlobResponse{Err: ""}
	args := fmt.Sprintf("%s --tmp-path=%s --container-name=%s", req.GetTargetPath(), req.GetTmpPath(), req.GetContainerName())
	cmd := exec.Command("blobfuse", strings.Split(args, " ")...)
	// todo: take mount args being passed from storage class

	cmd.Env = append(os.Environ(), "AZURE_STORAGE_ACCOUNT="+req.GetAccountName())
	cmd.Env = append(cmd.Env, "AZURE_STORAGE_ACCESS_KEY="+req.GetAccountKey())
	err = cmd.Run()
	if err != nil {
		log.Println("blobfuse mount failed")
		resp.Err = err.Error()
		return resp, err
	}
	log.Println("successfull mounted")
	return resp, nil
}
