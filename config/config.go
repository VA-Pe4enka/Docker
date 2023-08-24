package config

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
)

// GetImageName returns Image name string
func GetImageName() string {
	var imageName string
	fmt.Print("Enter Image name: ")
	fmt.Scan(&imageName)
	return imageName
}

// GetContainerName returns container name string
func GetContainerName() string {
	var contName string
	fmt.Print("Enter container /name: ")
	fmt.Scan(&contName)
	return contName
}

func SetContConfig() container.Config {
	config := container.Config{}

	return config
}
