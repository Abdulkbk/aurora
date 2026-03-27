package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/abdulkbk/aurora/internal/docker"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Docker images built by Aurora",
	Long: `List all Docker images that were built by Aurora.
	Aurora tags all built images with the ":aurora" suffix, so this command
	filters local Docker images by that tag and displays them in a table.
	Example:
		aurora list`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	images, err := docker.ListAuroraImages()
	if err != nil {
		return err
	}

	if len(images) == 0 {
		fmt.Println("No Aurora images found. Build one with: " +
			"aurora build")

		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "IMAGE\tTAG\tCREATED")

	for _, img := range images {
		fmt.Fprintf(
			w, "%s\t%s\t%s\n", img.Repository, img.Tag,
			img.CreatedSince,
		)
	}

	return w.Flush()
}
