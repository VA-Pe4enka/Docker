package main

import (
	"Docker/cmd"
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

	cmd.Execute()
}
