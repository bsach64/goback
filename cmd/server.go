package cmd

import (
	"github.com/bsach64/goback/server"
	"github.com/bsach64/goback/utils"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var mainOptions string
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts a server",
	Long:  "starts a server",
	Run: func(cmd *cobra.Command, args []string) {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose an option:").
					Options(
						huh.NewOption("Listen", "Listen"),
						huh.NewOption("Worker", "Worker"),
						huh.NewOption("Log", "Log"),
						huh.NewOption("Exit", "Exit"),
					).
					Value(&mainOptions),
			),
		)
		err := form.Run()

		if err != nil {
			log.Fatal("Failed to Run Server Action Form", "err", err)
		}

		switch mainOptions {
		case "Listen":
			localip, err := utils.GetLocalIP()
			if err != nil {
				log.Fatal("Could Not Local IP for server", "err", err)
			}

			server.NewMaster(localip.String())
			if err != nil {
				log.Fatal("Could Not Listen on Server", "err", err)
			}

		case "Worker":
			// TODO
		case "Log":
			// TODO
		case "Exit":
			// TODO
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
