package cmd

import (
	"log"

	"github.com/bsach64/goback/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts a server",
	Long:  "starts a server",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.New("0.0.0.0", "private/id_rsa", 2022)
		err := server.Listen(s)
		if err != nil {
			log.Println("Could not listen on server")
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
