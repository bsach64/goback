package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/bsach64/goback/client"
	"github.com/bsach64/goback/server"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var userClient client.Client

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Connect to server and perform actions like upload, list, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		ip, err := promptForIP()
		if err != nil {
			log.Fatal("Could not get Server IP", "err", err)
		}

		log.Info("Connecting to Server...")

		sshC, err := userClient.ConnectToServer(ip)
		if err != nil {
			log.Fatal("Failed to connect to server:", "err", err)
		}

		log.Info("Connected to Server!")

		defer sshC.Close()

		for {
			selectedOption, err := promptForAction()
			if err != nil {
				log.Fatal("Could not get action", "err", err)
			}

			switch selectedOption {
			case "Upload File":
				path, err := promptForFilePath()

				if err != nil {
					log.Error("Could not get file path for upload", "err", err)
					continue
				}

				success, reply, err := sshC.SendRequest("create-backup", true, []byte("Get Worker IP"))

				if err != nil {
					log.Fatalf("Failed to send %s request: %v", "create-backup", err)
				}

				if !success {
					log.Info("ssh request for create-backup failed")
					continue
				}

				var workerNode server.Worker
				if err := json.Unmarshal(reply, &workerNode); err != nil {
					log.Fatalf("failed to unmarshal response: %v", err)
				}

				// Worker node ip and port
				host := fmt.Sprintf("%s:%d", workerNode.Ip, workerNode.Port)

				// Worker node username and password for login
				// Will change this to digital signature later
				c := client.NewClient(workerNode.SftpUser, workerNode.SftpPass)

				// Connect to sftp server i.e worker node
				sftpClient, err := c.ConnectToServer(host)
				defer sftpClient.Close()

				if err != nil {
					log.Fatal("Could not connect to worker node", "err", err)
				}

				err = client.Upload(sftpClient, path)

				if err != nil {
					log.Fatalf("Cannot upload file to worker node %s at because %s", host, err)
				}

				log.Info("Successfully Uploaded", "file", path)

			case "List Directory":
				listRemoteDir()

			case "Exit":
				fmt.Println("Exiting client.")
				sshC.Close()
				return
			}
		}
	},
}

func promptForIP() (string, error) {
	var ip string
	ipPrompt := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Server IP:").
				Prompt("? ").
				Placeholder("0.0.0.0:8080").
				Suggestions([]string{"0.0.0.0:8080"}).
				Value(&ip),
		),
	)

	err := ipPrompt.Run()
	if err != nil {
		return "", err
	}
	return ip, nil
}

func promptForAction() (string, error) {
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
		return "", err
	}
	return selectedOption, err
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
	rootCmd.AddCommand(clientCmd)
}
