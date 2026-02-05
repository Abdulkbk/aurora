package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Supported node types in Polar
var supportedNodeTypes = []string{"bitcoind", "lnd", "eclair", "cln", "litd", "tapd"}

// Build command flags
var (
	prURL    string
	repoURL  string
	branch   string
	nodeType string
	imageTag string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a Docker image from a GitHub PR or fork",
	Long: `Build a Docker image from a GitHub PR or fork for use in Lightning Polar.

You can specify either a PR URL or a repository URL with a branch.

Examples:
  # Build from a PR
  aurora build --pr https://github.com/lightningnetwork/lnd/pull/1234 --tag my-lnd-test

  # Build from a fork/branch
  aurora build --repo https://github.com/myuser/lnd --branch feature-x --tag my-lnd-fork

  # Build with explicit node type
  aurora build --pr https://github.com/lightningnetwork/lnd/pull/1234 --node-type lnd --tag test`,
	RunE: runBuild,
}

func init() {
	buildCmd.Flags().StringVar(&prURL, "pr", "", "GitHub PR URL (e.g., https://github.com/owner/repo/pull/123)")
	buildCmd.Flags().StringVar(&repoURL, "repo", "", "GitHub repository URL (e.g., https://github.com/owner/repo)")
	buildCmd.Flags().StringVar(&branch, "branch", "", "Branch name (required with --repo)")
	buildCmd.Flags().StringVar(&nodeType, "node-type", "", fmt.Sprintf("Node type: %v (auto-detected if not specified)", supportedNodeTypes))
	buildCmd.Flags().StringVar(&imageTag, "tag", "", "Custom tag for the Docker image (required)")

	buildCmd.MarkFlagRequired("tag")

	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Validate that either --pr or --repo is provided
	if prURL == "" && repoURL == "" {
		return fmt.Errorf("either --pr or --repo must be specified")
	}

	if prURL != "" && repoURL != "" {
		return fmt.Errorf("cannot specify both --pr and --repo")
	}

	if repoURL != "" && branch == "" {
		return fmt.Errorf("--branch is required when using --repo")
	}

	// Validate node type if provided
	if nodeType != "" && !isValidNodeType(nodeType) {
		return fmt.Errorf("invalid node type %q, must be one of: %v", nodeType, supportedNodeTypes)
	}

	// Placeholder - actual implementation will come in future steps
	fmt.Println("🚀 Aurora Build")
	fmt.Println("===============")

	if prURL != "" {
		fmt.Printf("PR URL:    %s\n", prURL)
	} else {
		fmt.Printf("Repo URL:  %s\n", repoURL)
		fmt.Printf("Branch:    %s\n", branch)
	}

	if nodeType != "" {
		fmt.Printf("Node Type: %s\n", nodeType)
	} else {
		fmt.Println("Node Type: (will auto-detect)")
	}

	fmt.Printf("Tag:       %s\n", imageTag)
	fmt.Println()
	fmt.Println("⚠️  Build functionality coming in Step 2-4...")

	return nil
}

func isValidNodeType(nt string) bool {
	for _, valid := range supportedNodeTypes {
		if nt == valid {
			return true
		}
	}
	return false
}
