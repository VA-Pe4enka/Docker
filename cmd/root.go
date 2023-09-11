package cmd

import (
	"Docker/backend"
	"Docker/config"
	"github.com/spf13/cobra"
	"os"
)

var container = &cobra.Command{
	Use:   "container",
	Short: "Use this command to rule your containers",
	Run: func(cmd *cobra.Command, args []string) {
		flag, _ := cmd.Flags().GetString("start")
		if flag == "-new" {
			cfg := config.GetContConfig()
			backend.StartNewCont(cfg)
		} else if flag == "-exist" {
			backend.GetStoppedConts()
			contName := config.GetContainerName()
			backend.StartExistCont(contName)
		}

		flag, _ = cmd.Flags().GetString("stop")
		if flag == "-all" {
			backend.StopAllConts()
		} else if flag == "-s" {
			backend.GetRunningConts()
			contName := config.GetContainerName()
			backend.StopOneCont(contName)
		}

		flag, _ = cmd.Flags().GetString("ls")
		if flag == "-all" {
			backend.GetAllContainers()
		} else if flag == "-r" {
			backend.GetRunningConts()
		}

		flag, _ = cmd.Flags().GetString("log")
		if flag == "." {
			backend.GetRunningConts()
			contName := config.GetContainerName()
			backend.GetContLogs(contName)
		}

		flag, _ = cmd.Flags().GetString("commit")
		if flag == "." {
			backend.GetRunningConts()
			contName := config.GetContainerName()
			backend.CommitCont(contName)
		}

	},
}

var image = &cobra.Command{
	Use:   "image",
	Short: "Use this command to rule your images",
	Run: func(cmd *cobra.Command, args []string) {
		flag, _ := cmd.Flags().GetString("ls")
		if flag == "-a" {
			backend.GetAllImages()
		}

		flag, _ = cmd.Flags().GetString("pull")
		if flag == "." {
			imageName := config.GetImageName()
			backend.PullImage(imageName)
		}
	},
}

var rootCmd = &cobra.Command{
	Use: " ",
	//	Short: "A brief description of your application",
	//	Long: `A longer description that spans multiple lines and likely contains
	//examples and usage of using your application. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(container)
	container.PersistentFlags().String("start", "", "Use this flag to star new or exist containers")
	container.PersistentFlags().String("stop", "", "Use this flag to stop one or all containers")
	container.PersistentFlags().String("ls", "", "Use this flag to get containers list")
	container.PersistentFlags().String("log", ".", "Use this flag to get container logs")
	container.PersistentFlags().String("commit", ".", "Use this flag to commit container")

	rootCmd.AddCommand(image)
	image.PersistentFlags().String("ls", "", "Use this flag to get images list")
	image.PersistentFlags().String("pull", ".", "Use this flag to pull new image")
}
