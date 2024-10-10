package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/bsach64/goback/client"
	"github.com/bsach64/goback/server"
	"github.com/bsach64/goback/utils"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	userClient client.Client
	worker     server.Worker
)

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

		workerIP, err := utils.GetLocalIP()
		if err != nil {
			log.Fatal("Could not get Local IP", "err", err)
		}

		worker.Ip = workerIP.String()
		worker.Port = 2025

		go worker.StartSFTPServer()

		dat, err := json.Marshal(worker)
		if err != nil {
			log.Fatal("Could not Marshal worker info", "err", err)
		}

		success, _, err := sshC.SendRequest("worker-details", true, dat)
		if err != nil {
			log.Fatal("Could not Send Worker details to Master", "err", err)
		}

		if !success {
			log.Fatal("Server Could not Save Worker Details")
		}

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
					log.Warn("ssh request for create-backup failed", "reply", string(reply))
					continue
				}

				var otherWorker server.Worker
				if err := json.Unmarshal(reply, &otherWorker); err != nil {
					log.Fatalf("failed to unmarshal response: %v", err)
				}

				// Worker node ip and port
				host := fmt.Sprintf("%s:%d", otherWorker.Ip, otherWorker.Port)

				// Worker node username and password for login
				// Will change this to digital signature later
				c := client.NewClient(otherWorker.SftpUser, otherWorker.SftpPass)

				// Connect to sftp server i.e worker node
				sftpClient, err := c.ConnectToServer(host)
				if err != nil {
					log.Fatal("Could not connect to worker node", "err", err)
				}
				defer sftpClient.Close()

				err = client.Upload(sftpClient, path)

				if err != nil {
					log.Fatalf("Cannot upload file to worker node %s at because %s", host, err)
				}

				log.Info("Successfully Uploaded", "file", path)

			case "List Directory":
				listRemoteDir()

			case "Exit":
				fmt.Println("Exiting client.")
				_, _, err := sshC.SendRequest("close-connection", false, []byte(worker.Ip))
				if err != nil {
					log.Error("Error while closing the connection with server")
				}
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
				Placeholder("0.0.0.0:2022").
				Suggestions([]string{"0.0.0.0:2022"}).
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
