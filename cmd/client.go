package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bsach64/goback/client"
	"github.com/bsach64/goback/server"
	"github.com/charmbracelet/huh"
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

		c, err := userClient.ConnectToServer(clientArgs.host)
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}

		fmt.Println("Connected to Server ! ")

		defer c.Close()

		for {
			var selectedOption string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Choose an option:").
						Options(
							huh.NewOption("Upload File", "Upload File"),
							huh.NewOption("List Directory", "List Directory"),
							huh.NewOption("Exit", "Exit"),
						).
						Value(&selectedOption),
				),
			)
			err := form.Run()

			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}

			switch selectedOption {
			case "Upload File":
				path, err := promptForFilePath()

				if err != nil {
					log.Fatalf("Error while reading file path")
				}

				createBackupPayload := []byte("Get Server IP")
				success, reply, err := c.SendRequest("create-backup", true, createBackupPayload)

				if err != nil {
					log.Fatalf("failed to send %s request: %v", "create-backup", err)
				}

				if success {
					var workerNode server.Worker
					if err := json.Unmarshal(reply, &workerNode); err != nil {
						log.Fatalf("failed to unmarshal response: %v", err)
					}

					//Worker node ip and port
					host := fmt.Sprintf("%s:%d", workerNode.Ip, workerNode.Port)

					//Worker node username and password for login
					// Will change this to digital signature later
					c := client.NewClient(workerNode.SftpUser, workerNode.SftpPass)

					//Connect to sftp server i.e worker node
					sftpClient, err := c.ConnectToServer(host)

					if err != nil {
						log.Fatalf("Cannot connect to worker node")
					}

					err = client.Upload(sftpClient, path)

					if err != nil {
						log.Printf("Cannot upload file to worker node %s at because %s", host, err)
					}

					sftpClient.Close() //using defer for this doesn't seem to work for some reason

				} else {
					fmt.Println("create-backup request failed")
				}

			case "List Directory":
				listRemoteDir()

			case "Exit":
				fmt.Println("Exiting client.")
				c.Close()
				return
			}
		}
	},
}

func promptForFilePath() (string, error) {
	var filepath string
	filePrompt := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter File Path:").
				Prompt("? ").
				Placeholder("test_files/example.txt").
				Suggestions([]string{"test_files/example.txt"}).
				// Default("test_files/example.txt").
				Validate(validateFilePath).
				Value(&filepath),
		),
	)

	err := filePrompt.Run()
	if err != nil {
		return "", err
		// Also autocomplete is required
	}
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
	// Persistent flags for subcommands
	clientCmd.PersistentFlags().StringVarP(&clientArgs.user, "user", "u", "demo", "username")
	clientCmd.PersistentFlags().StringVarP(&clientArgs.password, "password", "p", "password", "password")
	clientCmd.PersistentFlags().StringVarP(&clientArgs.host, "host", "H", "127.0.0.1:2022", "host address")
	rootCmd.AddCommand(clientCmd)
}
