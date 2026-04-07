package cmd

import (
	"context"
	"fmt"

	"github.com/abdulkbk/aurora/internal/docker"
	"github.com/abdulkbk/aurora/internal/github"
	"github.com/spf13/cobra"
)

// Supported node types in Polar
var supportedNodeTypes = []string{"lnd", "bitcoind", "cln", "btcd"}

// repoToNodeType maps known upstream GitHub repos to their node type.
var repoToNodeType = map[string]string{
	"lightningnetwork/lnd":      "lnd",
	"btcsuite/btcd":             "btcd",
	"bitcoin/bitcoin":           "bitcoind",
	"ElementsProject/lightning": "cln",
}

// Build command flags
var (
	prURL    string
	repoURL  string
	branch   string
	nodeType string
	imageTag string
	noCache  bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a Docker image from a GitHub PR or fork",
	Long: `Build a Docker image from a GitHub PR or fork for use in Lightning Polar.

You can specify either a PR URL or a repository URL with a branch.
The node type is auto-detected from the repository URL. Use --node-type to
override this when building from a fork with a non-standard repository name.

Examples:
  # Build from a PR (node type auto-detected)
  aurora build --pr https://github.com/lightningnetwork/lnd/pull/1234 --tag my-lnd-test

  # Build from a fork/branch (node type auto-detected)
  aurora build --repo https://github.com/myuser/lnd --branch feature-x --tag my-lnd-fork

  # Build from a fork with a custom repo name (node type cannot be auto-detected)
  aurora build --repo https://github.com/myuser/my-lnd-fork --branch feature-x --node-type lnd --tag test`,
	RunE: runBuild,
}

func init() {
	buildCmd.Flags().StringVar(&prURL, "pr", "", "GitHub PR URL (e.g., https://github.com/owner/repo/pull/123)")
	buildCmd.Flags().StringVar(&repoURL, "repo", "", "GitHub repository URL (e.g., https://github.com/owner/repo)")
	buildCmd.Flags().StringVar(&branch, "branch", "", "Branch name (required with --repo)")
	buildCmd.Flags().StringVar(&nodeType, "node-type", "", fmt.Sprintf("Node type override: %v (auto-detected if omitted)", supportedNodeTypes))
	buildCmd.Flags().StringVar(&imageTag, "tag", "", "Custom tag for the Docker image (required)")
	buildCmd.Flags().BoolVar(&noCache, "no-cache", false, "Build the image without using Docker cache")

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

	fmt.Println("[aurora] build")
	fmt.Println("===============")

	var gitURL, gitBranch, owner, repo string

	if prURL != "" {
		// Parse PR URL and fetch details from GitHub API
		prInfo, err := github.ParsePRURL(prURL)
		if err != nil {
			return err
		}

		owner, repo = prInfo.Owner, prInfo.Repo

		fmt.Printf(">> PR:     %s/%s#%d\n", prInfo.Owner, prInfo.Repo, prInfo.PRNumber)
		fmt.Println(".. fetching PR details")

		client := github.NewClient()
		prDetails, err := client.GetPRDetails(prInfo.Owner, prInfo.Repo, prInfo.PRNumber)
		if err != nil {
			return fmt.Errorf("failed to fetch PR details: %w", err)
		}

		fmt.Printf(">> Title:  %s\n", prDetails.Title)
		fmt.Printf(">> State:  %s\n", prDetails.State)
		fmt.Printf(">> Fork:   %s\n", prDetails.ForkURL)
		fmt.Printf(">> Branch: %s\n", prDetails.Branch)

		gitURL = prDetails.ForkURL
		gitBranch = prDetails.Branch
	} else {
		// Use direct repo/branch input
		repoInfo, err := github.ParseRepoURL(repoURL)
		if err != nil {
			return err
		}

		owner, repo = repoInfo.Owner, repoInfo.Repo
		gitURL = repoInfo.CloneURL()
		gitBranch = branch

		fmt.Printf(">> Repo:   %s\n", gitURL)
		fmt.Printf(">> Branch: %s\n", gitBranch)
	}

	// Auto-detect node type if not provided
	if nodeType == "" {
		detected, ok := detectNodeType(owner, repo)
		if !ok {
			return fmt.Errorf(
				"could not auto-detect node type for %s/%s, please specify --node-type (one of: %v)",
				owner, repo, supportedNodeTypes,
			)
		}
		nodeType = detected
	}

	fmt.Printf(">> type:   %s\n", nodeType)

	fmt.Printf(">> tag:    %s\n", imageTag)
	fmt.Println()

	// Create Docker builder
	builder, err := docker.NewBuilder()
	if err != nil {
		return fmt.Errorf("failed to initialize Docker: %w", err)
	}

	// Build the image
	fmt.Println(".. building image")
	fmt.Println("----------------------------")

	ctx := context.Background()
	err = builder.Build(ctx, docker.BuildOptions{
		GitURL:   gitURL,
		Checkout: gitBranch,
		Tag:      imageTag + ":aurora",
		NodeType: nodeType,
		NoCache:  noCache,
	})
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Printf("[ok] build complete! Image: %s\n", imageTag)
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

func detectNodeType(owner, repo string) (string, bool) {
	nt, ok := repoToNodeType[owner+"/"+repo]
	return nt, ok
}
