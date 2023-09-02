package backend

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
	"strings"
)

// StartNewCont starts new container
// also upload new Image if one is missing
func StartNewCont(config container.Config) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	out, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("ERROR in pulling new Image: %s", err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &config, nil, nil, nil, "")
	if err != nil {
		log.Printf("ERROR in creating container: %s", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Container start ERROR: %s", err)
	}

	fmt.Println(resp.ID)
}

// StartExistCont starts one of existing containers
func StartExistCont(contName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("ERROR getting container list: %s", err)
	}

	conts := make(map[string]string)
	for _, container := range containers {
		conts[strings.Join(container.Names, "")] = container.ID
	}

	ID, ok := conts["/"+contName]
	if ok {
		if err := cli.ContainerStart(ctx, ID, types.ContainerStartOptions{}); err != nil {
			log.Printf("Container start ERROR: %s", err)
		}
	}

}

// GetRunningConts get a list of running containers
func GetRunningConts() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Printf("ERROR getting running container list: %s", err)
	}

	fmt.Println("List of running containers: ")

	for _, container := range containers {
		if strings.Contains(container.Image, "sha256") {
			container.Image = "<none>"
		}
		fmt.Println(container.Names, container.Image)
	}

	fmt.Println()
}

// GetAllConts get a list of all containers
func GetAllConts() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("ERROR getting container list: %s", err)
	}

	fmt.Println("List of all containers: ")

	for _, container := range containers {
		if strings.Contains(container.Image, "sha256") {
			container.Image = "<none>"
		}
		fmt.Println(container.Names, container.Image)
	}

	fmt.Println()
}

// StopAllConts stops all running containers
func StopAllConts() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}

	// Получение списка запуцщенных контейнеров(docker ps)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Printf("ERROR getting running container list: %s", err)
	}

	for _, c := range containers {
		fmt.Print("Stopping container ", c.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, c.ID, container.StopOptions{Signal: "SIGKILL", Timeout: nil}); err != nil {
			log.Printf("Container was stopped with ERROR: %s", err)
		}
		fmt.Println("All containers stopped successfully")
	}

}

// GetAllImages get a list of all uploaded Images
func GetAllImages() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Printf("ERROR getting Images list: %s", err)
	}

	fmt.Println("List of all Images: ")

	for _, image := range images {
		if len(image.RepoTags) == 0 {
			image.RepoTags = append(image.RepoTags, "<none>")
		}
		fmt.Println(image.RepoTags, image.ID)
	}

	fmt.Println()
}

// PullImage upload new Image
func PullImage(imageName string) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		log.Printf("ERROR in pulling new Image: %s", err)
	}
	defer out.Close()

	io.Copy(os.Stdout, out)
}

// GetContLogs get container logs
func GetContLogs(contName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("ERROR getting containers list: %s", err)
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
			log.Printf("ERROR getting container logs: %s", err)
		}
		io.Copy(os.Stdout, out)
	}
}

// CommitCont commits container
func CommitCont(contName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("ERROR getting containers list: %s", err)
	}

	conts := make(map[string]string)
	for _, container := range containers {
		conts[strings.Join(container.Names, "")] = container.ID
	}

	ID, ok := conts["/"+contName]
	if ok {
		commitResp, err := cli.ContainerCommit(ctx, ID, types.ContainerCommitOptions{})
		if err != nil {
			log.Printf("ERROR commiting container: %s", err)
		}
		fmt.Println("Commit success: ", commitResp.ID)
	}
}

// GetStoppedConts displays list of running and list of all containers
func GetStoppedConts() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Printf("ERROR getting  containers list: %s", err)
	}

	stoppedContainers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("ERROR getting running containers list: %s", err)
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

// StopContViaImage stopping container with the specified images
func StopContViaImage(imageName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Printf("ERROR getting running container list: %s", err)
	}

	for _, c := range containers {
		if strings.Contains(c.Image, "sha256") {
			c.Image = "<none>"
		}
		if strings.Contains(c.Image, imageName) {
			fmt.Print("Stopping container ", c.ID[:10], "... ")
			if err := cli.ContainerStop(ctx, c.ID, container.StopOptions{Signal: "SIGKILL", Timeout: nil}); err != nil {
				log.Printf("Container was stopped with ERROR: %s", err)
			}
			fmt.Printf("%s containers stopped successfully", imageName)
			fmt.Println()
		}
	}

}

// DeleteContViaImage removing container with the specified images
func DeleteContViaImage(imageName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("ERROR in opening client: %s", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("ERROR getting running container list: %s", err)
	}

	for _, c := range containers {
		if strings.Contains(c.Image, "sha256") {
			c.Image = "<none>"
		}
		if strings.Contains(c.Image, imageName) {
			fmt.Print("Removing container ", c.ID[:10], "... ")
			if err := cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
				log.Printf("Container was removed with ERROR: %s", err)
			}
			fmt.Printf("%s containers removed successfully", imageName)
			fmt.Println()

		}
	}
}
