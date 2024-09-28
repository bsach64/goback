package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bsach64/goback/utils"
	"github.com/charmbracelet/huh"
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

	var fileName string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Snapshot File").
				Options(huh.NewOptions(snapshots...)...).
				Value(&fileName),
		),
	)

	err = form.Run()
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
