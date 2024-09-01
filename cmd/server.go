package cmd

import (
	"github.com/bsach64/goback/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts a server",
	Long:  "starts a server",
	Run: func(cmd *cobra.Command, args []string) {
		go server.Listen("")
    select {}
	},
}


func init(){
  rootCmd.AddCommand(serverCmd)
}
