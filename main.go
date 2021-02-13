package main

// idea:
// a GRPC service that interacts with the csi Node Service and does all the mounting job.
// as a prototype lets create a http server that does the mounting

// POST /mount
// temp
// targetPath
// accountName
// containerName
// accountKey
// validates by actually checking if the mount point exists
// returns {"success": true}

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type DataReq struct {
	TargetPath    string `json:"targetPath"`
	AccountName   string `json:"accountName"`
	ContainerName string `json:"containerName"`
	AccountKey    string `json:"acccountKey"`
	TmpPath       string `json:"tmpPath"`
}

func handleGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "use POST /mount to mount an endpoint"})
}

func handleMount(c *gin.Context) {
	var req DataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	args := fmt.Sprintf("%s --tmp-path=%s --container-name=%s", req.TargetPath, req.TmpPath, req.ContainerName)
	cmd := exec.Command("blobfuse", strings.Split(args, " ")...)
	// todo: take mount args being passed from storage class

	cmd.Env = append(os.Environ(), "AZURE_STORAGE_ACCOUNT="+req.AccountName)
	cmd.Env = append(cmd.Env, "AZURE_STORAGE_ACCESS_KEY="+req.AccountKey)
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "true"})
}

func main() {
	r := gin.Default()
	r.GET("/", handleGet)
	r.POST("/mount", handleMount)
	r.Run()
}
