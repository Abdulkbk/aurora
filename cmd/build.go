package cmd

import (
	"context"
	"fmt"

	"github.com/abdulkbk/aurora/internal/docker"
	"github.com/abdulkbk/aurora/internal/github"
	"github.com/spf13/cobra"
)

// Supported node types in Polar
var supportedNodeTypes = []string{"lnd", "bitcoind"}

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
	buildCmd.Flags().StringVar(&nodeType, "node-type", "", fmt.Sprintf("Node type: %v (required)", supportedNodeTypes))
	buildCmd.Flags().StringVar(&imageTag, "tag", "", "Custom tag for the Docker image (required)")

	buildCmd.MarkFlagRequired("tag")
	buildCmd.MarkFlagRequired("node-type")

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

	fmt.Println("🚀 Aurora Build")
	fmt.Println("===============")

	var gitURL, gitBranch string

	if prURL != "" {
		// Parse PR URL and fetch details from GitHub API
		prInfo, err := github.ParsePRURL(prURL)
		if err != nil {
			return err
		}

		fmt.Printf("📋 PR:     %s/%s#%d\n", prInfo.Owner, prInfo.Repo, prInfo.PRNumber)
		fmt.Println("🔍 Fetching PR details from GitHub...")

		client := github.NewClient()
		prDetails, err := client.GetPRDetails(prInfo.Owner, prInfo.Repo, prInfo.PRNumber)
		if err != nil {
			return fmt.Errorf("failed to fetch PR details: %w", err)
		}

		fmt.Printf("📝 Title:  %s\n", prDetails.Title)
		fmt.Printf("📊 State:  %s\n", prDetails.State)
		fmt.Printf("🔗 Fork:   %s\n", prDetails.ForkURL)
		fmt.Printf("🌿 Branch: %s\n", prDetails.Branch)

		gitURL = prDetails.ForkURL
		gitBranch = prDetails.Branch
	} else {
		// Use direct repo/branch input
		repoInfo, err := github.ParseRepoURL(repoURL)
		if err != nil {
			return err
		}

		gitURL = repoInfo.CloneURL()
		gitBranch = branch

		fmt.Printf("🔗 Repo:   %s\n", gitURL)
		fmt.Printf("🌿 Branch: %s\n", gitBranch)
	}

	// Node type is now required
	fmt.Printf("📦 Type:   %s\n", nodeType)

	fmt.Printf("🏷️  Tag:    %s\n", imageTag)
	fmt.Println()

	// Create Docker builder
	builder, err := docker.NewBuilder()
	if err != nil {
		return fmt.Errorf("failed to initialize Docker: %w", err)
	}

	// Build the image
	fmt.Println("🔨 Building Docker image...")
	fmt.Println("----------------------------")

	ctx := context.Background()
	err = builder.Build(ctx, docker.BuildOptions{
		GitURL:   gitURL,
		Checkout: gitBranch,
		Tag:      imageTag + ":aurora",
		NodeType: nodeType,
	})
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Printf("✅ Build complete! Image: %s\n", imageTag)
	fmt.Println()
	fmt.Println("To use in Polar, add this as a custom node image.")

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
