package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"runtime"
)

func startCont() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	imageName := "bfirsh/reticulate-splines"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}

func getAllConts() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

func stopCont() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// Получение списка запуцщенных контейнеров(docker ps)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		fmt.Print("Stopping container ", c.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, c.ID, container.StopOptions{Signal: "SIGKILL", Timeout: nil}); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}

}

func main() {

	// выполнить в терминале команду docker context ls, выбрать DOCKER ENDPOINT, соответствующий Docker Desktop
	// и подставить значение в dockerHost

	if len(os.Getenv("DOCKER_HOST")) == 0 {
		if runtime.GOOS == "windows" {
			os.Setenv("DOCKER_HOST", "\\\\.\\pipe\\docker_engine:\\\\.\\pipe\\docker_engine")
		} else if runtime.GOOS == "linux" {
			os.Setenv("DOCKER_HOST", "unix:///home/"+os.Getenv("USER")+"/.docker/desktop/docker.sock")
		} else if runtime.GOOS == "darwin" {
			os.Setenv("DOCKER_HOST", "unix:///Users/"+os.Getenv("USER")+"/.docker/desktop/docker.sock")
		}
	}

	//fmt.Println(CTX.EndpointMetaBase{})

	//startCont()

	getAllConts()

	//stopCont()
}
