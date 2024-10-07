package cmd

import (
	"github.com/bsach64/goback/server"
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
						huh.NewOption("Reconstruct", "Reconstruct"),
						huh.NewOption("Exit", "Exit"),
					).
					Value(&mainOptions),
			),
		)
        err := form.Run()

		if err != nil {
			log.Fatal("failed", "err", err)
			// log.Fatal(err)
		}

		switch mainOptions {
		case "Listen":
			if err != nil {
				// log.Fatal(err)
				log.Fatal("failed", "err", err)
			}
			server.NewMaster()
			if err != nil {
				log.Info("Could not listen on server")
			}

        case "Worker":
            
		case "Log":
			//TODO
		case "Exit":
			//TODO
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
