package main

import (
	"Docker/backend"
	"fmt"
	"os"
	"runtime"
)

func main() {

	if len(os.Getenv("DOCKER_HOST")) == 0 {
		if runtime.GOOS == "windows" {
			os.Setenv("DOCKER_HOST", "\\\\.\\pipe\\docker_engine:\\\\.\\pipe\\docker_engine")
		} else if runtime.GOOS == "linux" {
			os.Setenv("DOCKER_HOST", "unix:///home/"+os.Getenv("USER")+"/.docker/desktop/docker.sock")
		} else if runtime.GOOS == "darwin" {
			os.Setenv("DOCKER_HOST", "unix:///Users/"+os.Getenv("USER")+"/.docker/run/docker.sock")
		}
	}

	var option int
	for {
		option = 123
		fmt.Println("Choose an option: ")
		fmt.Println("1 - start new container.")
		fmt.Println("2 - start exist container.")
		fmt.Println("3 - get a list of running containers")
		fmt.Println("4 - get a list of all containers")
		fmt.Println("5 - stop all running containers")
		fmt.Println("6 - get a list of all uploaded Images")
		fmt.Println("7 - upload new Image")
		fmt.Println("8 - get container logs")
		fmt.Println("9 - commit container")
		fmt.Println("0 - exit")

		fmt.Print("Option: ")
		fmt.Scan(&option)
		fmt.Println()

		switch option {
		case 1:
			backend.StartNewCont()

		case 2:
			backend.StartExistCont()

		case 3:
			backend.GetRunningConts()

		case 4:
			backend.GetAllConts()

		case 5:
			backend.StopAllConts()

		case 6:
			backend.GetAllImages()

		case 7:
			backend.PullImage()

		case 8:
			backend.GetContLogs()

		case 9:
			backend.CommitCont()

		case 0:
			return

		default:
			fmt.Println("Not an option!")

		}
	}

}
