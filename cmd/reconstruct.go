package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	_, fileName, err := prompt.Run()
	if err != nil {
		log.Fatalf("Reconstruct prompt failed %v\n", err)
	}

	filePath := filepath.Join("./.data", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Could not open snapshot file: %v\n", err)
	}
	defer file.Close()

	byteData, err := utils.Reconstruct(file)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("Here's the reconstructed file:\n%v", string(byteData))
}

func init() {
	rootCmd.AddCommand(reconstructCmd)
}
