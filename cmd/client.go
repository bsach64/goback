package cmd

import (
	"log"

	"github.com/bsach64/goback/client"
	"github.com/spf13/cobra"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client -u [user] -p [password] -H [host_addr] -f [filepath]",
		Short: "starts a client that connects to [host_addr] and sends [filepath]",
		Long:  "starts a client that connects to [host_addr] and sends [filepath]",
		Run:   StartClient,
	}

	clientArgs struct {
		user     string
		password string
		host     string
		f        string
	}
)

func StartClient(cmd *cobra.Command, args []string) {
	user_client := client.NewClient(clientArgs.user, clientArgs.password)

	c, err := user_client.ConnectToServer(clientArgs.host)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer c.Close()
	err = client.Upload(c, clientArgs.f)
	if err != nil {
		log.Fatalf("Failed to upload the file to the server %v", err)
	}
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&clientArgs.user, "user", "u", "demo", "username")
	clientCmd.Flags().StringVarP(&clientArgs.password, "password", "p", "password", "password")
	clientCmd.Flags().StringVarP(&clientArgs.host, "host", "H", "127.0.0.1:2022", "host address")
	clientCmd.Flags().StringVarP(&clientArgs.f, "filepath", "f", "", "file path")
	if clientCmd.MarkFlagRequired("filepath") != nil {
		log.Fatalf("Required flag -f is missing")
	}
}
