package cmd

import (
	"github.com/bsach64/goback/server"
	"github.com/bsach64/goback/utils"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts a server",
	Long:  "starts a server",
	Run: func(cmd *cobra.Command, args []string) {
		localip, err := utils.GetLocalIP()
		if err != nil {
			log.Fatal("Could Not Local IP for server", "err", err)
		}

		log.Info("Starting Master Server")
		server.NewMaster(localip.String())
		if err != nil {
			log.Fatal("Could Not Listen on Server", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
