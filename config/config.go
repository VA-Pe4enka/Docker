package config

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"log"
	"os"
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

// GetContConfig getting configuration from config file
func GetContConfig() container.Config {
	config := container.Config{}

	fileData, err := os.ReadFile("config/config-json")
	if err != nil {
		log.Printf("ERROR reading config file: %s", err)
	}

	err = json.Unmarshal(fileData, &config)
	if err != nil {
		log.Printf("ERROR unmarshaling config: %s", err)
	}

	return config
}
