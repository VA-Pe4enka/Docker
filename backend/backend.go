package backend

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"strings"
)

// StartNewCont starts new container
// also upload new Image if one is missing
func StartNewCont(imageName string) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

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

// StartExistCont starts one of existing containers
func StartExistCont(contName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	conts := make(map[string]string)
	for _, container := range containers {
		conts[strings.Join(container.Names, "")] = container.ID
	}

	ID, ok := conts["/"+contName]
	if ok {
		if err := cli.ContainerStart(ctx, ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
	}

}

// GetRunningConts get a list of running containers
func GetRunningConts() {
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
		fmt.Println(container.Names, container.ID)
	}
}

// GetAllConts get a list of all containers
func GetAllConts() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.Names, container.ID)
	}
}

// StopAllConts stops all running containers
func StopAllConts() {
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

// GetAllImages get a list of all uploaded Images
func GetAllImages() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID)
	}
}

// PullImage upload new Image
func PullImage(imageName string) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	io.Copy(os.Stdout, out)
}

// GetContLogs get container logs
func GetContLogs(contName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true}

	conts := make(map[string]string)
	for _, container := range containers {
		conts[strings.Join(container.Names, "")] = container.ID
	}

	ID, ok := conts["/"+contName]
	if ok {
		out, err := cli.ContainerLogs(ctx, ID, options)
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, out)
	}
}

// CommitCont commits container
func CommitCont(contName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	conts := make(map[string]string)
	for _, container := range containers {
		conts[strings.Join(container.Names, "")] = container.ID
	}

	ID, ok := conts["/"+contName]
	if ok {
		commitResp, err := cli.ContainerCommit(ctx, ID, types.ContainerCommitOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Println("Commit success: ", commitResp.ID)
	}
}

// GetImageName returns Image name string
func GetImageName() string {
	var imageName string
	fmt.Print("Enter Image name: ")
	fmt.Scan(&imageName)
	return imageName
}

// GetStoppedConts displays list of running and list of all containers
func GetStoppedConts() {
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

	stoppedContainers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	fmt.Println("Running containers: ")
	for _, container := range containers {
		fmt.Println(container.Names, container.ID)
	}
	fmt.Println()

	fmt.Println("All containers: ")
	for _, container := range stoppedContainers {
		fmt.Println(container.Names, container.ID)
	}
	fmt.Println()
}

// GetContainerName returns container name string
func GetContainerName() string {
	var contName string
	fmt.Print("Enter container /name: ")
	fmt.Scan(&contName)
	return contName
}
