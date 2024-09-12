package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	contents, err := os.ReadDir("./.data/snapshots")
	if err != nil {
		log.Printf("Could not open Data dir, %v\n", err)
		return
	}

	snapshots := make([]string, 0)
	for _, entry := range contents {
		snapshots = append(snapshots, entry.Name())
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

	_, fileName, err := prompt.Run()
	if err != nil {
		log.Fatalf("Reconstruct prompt failed %v\n", err)
	}

	filePath := filepath.Join("./.data/snapshots", fileName)
	dat, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Could not open snapshot file: %v\n", err)
	}

	var snapshot utils.Snapshot
	err = json.Unmarshal(dat, &snapshot)
	if err != nil {
		log.Fatalf("Could not unmarshal snapshot file: %v\n", err)
	}

	byteData, err := utils.Reconstruct(snapshot)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("Here's the reconstructed file:\n%v", string(byteData))
}

func init() {
	rootCmd.AddCommand(reconstructCmd)
}
