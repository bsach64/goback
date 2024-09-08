package cmd

import (
	"log"

	"github.com/bsach64/goback/server"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts a server",
	Long:  "starts a server",
	Run: func(cmd *cobra.Command, args []string) {
		prompt := promptui.Select{
			Label: "Select Command",
			Items: []string{"Listen", "Log", "Reconstruct", "Exit"},
			Templates: &promptui.SelectTemplates{
				Active:   "* {{ . | bold | green }}", // Green color for the selected item
				Inactive: "{{ . }}",
				Selected: "* {{ . | bold | green }}", // Green color for the selected item
				Details:  "{{ . }}",
			},
		}
		_, result, err := prompt.Run()

		if err != nil {
			log.Fatalf("Server prompt failed %v\n", err)
		}

		switch result {
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
