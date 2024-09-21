package cmd

import (
	"log"
	"github.com/charmbracelet/huh"
	"github.com/bsach64/goback/server"
	"github.com/spf13/cobra"
)
var mainOptions string
var serverCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose an option:").
					Options(
						huh.NewOption("Listen","Listen"),
						huh.NewOption("Log","Log"),
						huh.NewOption("Reconstruct","Reconstruct"),
						huh.NewOption("Exit","Exit"),
					).
					Value(&mainOptions),
			),
		)
		err := form.Run()

		if err != nil {
			log.Fatal(err)
		}
		
		switch mainOptions {
		case "Listen":
			s := server.New("0.0.0.0", "private/id_rsa", 2022)
			err := server.Listen(s)
			if err != nil {
				log.Println("Could not listen on server")
			}
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