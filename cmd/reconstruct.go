package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/bsach64/goback/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var reconstructCmd = &cobra.Command{
	Use:   "reconstruct",
	Short: "Reconstruct a File",
	Long:  "Reconstruct a File",
	Run:   reconstruct,
}

func reconstruct(cmd *cobra.Command, args []string) {
	contents, err := os.ReadDir("./.data")
	if err != nil {
		log.Printf("Could not open Data dir, %v\n", err)
		return
	}

	snapshots := make([]string, 0)
	for _, entry := range contents {
		if strings.HasSuffix(entry.Name(), ".snapshot") {
			snapshots = append(snapshots, entry.Name())
		}
	}

	prompt := promptui.Select{
		Label: "Select Snapshot File",
		Items: snapshots,
		Templates: &promptui.SelectTemplates{
			Active:   "* {{ . | bold | green }}", // Green color for the selected item
			Inactive: "{{ . }}",
			Selected: "* {{ . | bold | green }}", // Green color for the selected item
			Details:  "{{ . }}",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Reconstruct prompt failed %v\n", err)
	}

	err = utils.Reconstruct(result)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

func init() {
	rootCmd.AddCommand(reconstructCmd)
}
