package cmd

import (
	"fmt"
	"log"

	"github.com/bsach64/goback/client"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	userClient client.Client
	clientArgs struct {
		user     string
		password string
		host     string
	}
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Connect to server and perform actions like upload, list, etc.",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Connecting to Server...")

		userClient = client.NewClient(clientArgs.user, clientArgs.password)

		c, err := userClient.ConnectToServer(clientArgs.host)
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}

		fmt.Println("Connected to Server ! ")

		userClient.SSHClient = c

		defer userClient.SSHClient.Close()

		for {
			prompt := promptui.Select{
				Label: "Select Command",
				Items: []string{"Upload File", "List Directory", "Exit"},
				Templates: &promptui.SelectTemplates{
					Active:   "* {{ . | bold | green }}", // Green color for the selected item
					Inactive: "{{ . }}",
					Selected: "* {{ . | bold | green }}", // Green color for the selected item
					Details:  "{{ . }}",
				},
			}

			_, result, err := prompt.Run()

			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}

			switch result {
			case "Upload File":
				filepath, err := promptForFilePath()
				if err != nil {
					log.Fatalf("Prompt failed %v\n", err)
				}

				err = client.Upload(userClient.SSHClient, filepath)
				if err != nil {
					log.Fatalf("Failed to upload file: %v", err)
				} else {
					fmt.Println("File uploaded successfully.")
				}

			case "List Directory":
				listRemoteDir()

			case "Exit":
				fmt.Println("Exiting client.")
				return
			}
		}
	},
}

func promptForFilePath() (string, error) {
	filePrompt := promptui.Prompt{
		Default:  "test_files/example.txt",
		Label:    "Enter File Path",
		Validate: validateFilePath,
	}

	filepath, err := filePrompt.Run()
	if err != nil {
		return "test_files/example.txt", nil //Currently hardcoded the value but in production this shall be validated
                                             // Also autocomplete is required
	}                                        // Maybe its better to change this to bubbletea later
	return filepath, nil
}

func validateFilePath(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("file path cannot be empty")
	}
	return nil
}

func listRemoteDir() {
	fmt.Println("ls doesn't work as of now")
}

func init() {
	//Persistent flags for subcommands
	clientCmd.PersistentFlags().StringVarP(&clientArgs.user, "user", "u", "demo", "username")
	clientCmd.PersistentFlags().StringVarP(&clientArgs.password, "password", "p", "password", "password")
	clientCmd.PersistentFlags().StringVarP(&clientArgs.host, "host", "H", "127.0.0.1:2022", "host address")

	rootCmd.AddCommand(clientCmd)
}
