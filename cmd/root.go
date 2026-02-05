package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information - will be set at build time
var (
	version = "0.1.0"
	commit  = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "aurora",
	Short: "Build custom Docker images from GitHub PRs for Polar",
	Long: `Aurora is a CLI tool that enables code reviewers to build Docker images
from GitHub PRs and forks for Lightning Polar's supported node implementations.

Supported node types: bitcoind, lnd, eclair, cln, litd, tapd

Example:
  aurora build --pr https://github.com/lightningnetwork/lnd/pull/1234 --tag my-lnd-test`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Aurora",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Aurora v%s (commit: %s)\n", version, commit)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
