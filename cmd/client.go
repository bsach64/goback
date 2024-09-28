package cmd

import (
	"fmt"
	"log"

	"github.com/bsach64/goback/client"
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

		userClient = client.NewClient(clientArgs.user, clientArgs.password)

		c, err := userClient.ConnectToServer(clientArgs.host)
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}

		fmt.Println("Connected to Server ! ")

		userClient.SSHClient = c

		defer userClient.SSHClient.Close()

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
				filepath, err := promptForFilePath()
				if err != nil {
					log.Fatalf("Prompt failed %v\n", err)
				}

				err = client.Upload(userClient.SSHClient, filepath)
				if err != nil {
					log.Fatalf("Failed to upload file %v: %v\n", filepath, err)
				} else {
					fmt.Printf("File %v uploaded successfully.\n", filepath)
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
