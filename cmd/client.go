package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"golang.org/x/crypto/ssh"

	"github.com/bsach64/goback/client"
	"github.com/bsach64/goback/server"
	"github.com/bsach64/goback/utils"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	userClient client.Client
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Connect to server and perform actions like upload, list, etc.",
	Run:   ClientLoop,
}

func ClientLoop(cmd *cobra.Command, args []string) {
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

	worker, err := CreateWorker()
	if err != nil {
		log.Fatal("Could not create worker", "err", err)
	}

	log.Info("Created Worker!")
	go worker.StartSFTPServer()

	log.Info("Sending Worker details to Server!")

	err = SendWorkerDetails(worker, sshC)
	if err != nil {
		log.Fatal("Could not send worker details", "err", err)
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

			stat, err := os.Stat(path)
			if err != nil {
				log.Error("Could not stat on file for upload", "err", err)
			}

			fileInfo := utils.FileInfo{
				Filename: stat.Name(),
				Size:     stat.Size(),
			}

			log.Infof("Sending file metadata to server: name: %v, size: %v", fileInfo.Filename, fileInfo.Size)
			fileInfoDat, err := json.Marshal(&fileInfo)
			if err != nil {
				log.Fatal("Could not marshal json", "err", err)
			}

			success, reply, err := sshC.SendRequest("start-file-upload", true, fileInfoDat)
			if err != nil {
				log.Fatalf("Failed to send %s request: %v", "start-file-upload", err)
			}

			if !success {
				log.Fatalf("Could not start file upload: %v", string(reply))
			}

			log.Infof("Starting file upload to other clients: name: %v, size: %v", fileInfo.Filename, fileInfo.Size)
			success, reply, err = sshC.SendRequest("create-backup", true, []byte("Get Worker IP"))

			if err != nil {
				log.Fatalf("Failed to send %s request: %v", "create-backup", err)
			}

			if !success {
				log.Warn("ssh request for create-backup failed", "reply", string(reply))
				continue
			}

			var otherWorkers []server.Worker
			if err := json.Unmarshal(reply, &otherWorkers); err != nil {
				log.Fatalf("failed to unmarshal response: %v", err)
			}

			// Worker node ip and port
			for _, w := range otherWorkers {
				if w.Ip == worker.Ip {
					continue
				}
				wip := fmt.Sprintf("%s:%d", w.Ip, w.Port)
				// Worker node username and password for login
				// Will change this to digital signature later
				c := client.NewClient(w.SftpUser, w.SftpPass)
				// Connect to sftp server i.e worker node
				sftpClient, err := c.ConnectToServer(wip)
				if err != nil {
					log.Fatal("Could not connect to worker node", "err", err)
				}
				err = client.Upload(sftpClient, path)

				if err != nil {
					log.Fatalf("Cannot upload file to worker node %s at because %s", wip, err)
				}
				log.Info("Successfully Uploaded", "file", path)

				sftpClient.Close()
			}

			success, reply, err = sshC.SendRequest("finish-file-upload", true, fileInfoDat)
			if err != nil {
				log.Fatalf("Failed to send %s request: %v", "finish-file-upload", err)
			}

			if !success {
				log.Warn("ssh request for finish-file-upload failed", "reply", string(reply))
				continue
			}

		case "Add Directory to Sync":
			dir, err := promptForDirectory()
			if err != nil {
				log.Error("Could not get directory for sync", "err", err)
				continue
			}

			if dir == "" {
				dir = "./files" // Default directory
			}
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
}

func CreateWorker() (server.Worker, error) {
	ip, err := utils.GetLocalIP()
	if err != nil {
		return server.Worker{}, err
	}
	return server.Worker{Ip: ip.String(), Port: 2025}, nil
}

func SendWorkerDetails(worker server.Worker, sshC *ssh.Client) error {
	dat, err := json.Marshal(worker)
	if err != nil {
		return err
	}

	success, reply, err := sshC.SendRequest("worker-details", true, dat)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("failed to send worker-details: %v", string(reply))
	}

	return nil
}

func promptForDirectory() (string, error) {
	var directory string
	directoryPrompt := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Directory to Watch (If it does not exist it will be create)").
				Prompt("? ").
				Placeholder("./files").
				Suggestions([]string{"./files"}).
				Value(&directory),
		),
	)
	err := directoryPrompt.Run()
	if err != nil {
		return "", err
	}
	return directory, nil
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
					huh.NewOption("Add Directory to Sync", "Add Directory to Sync"),
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

func init() {
	// Persistent flags for subcommands
	rootCmd.AddCommand(clientCmd)
}
