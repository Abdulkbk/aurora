// Package docker provides utilities for building Docker images.
package docker

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed dockerfiles/lnd/Dockerfile
var lndDockerfile embed.FS

//go:embed dockerfiles/bitcoind/Dockerfile
var bitcoindDockerfile embed.FS

//go:embed dockerfiles/cln/Dockerfile
var clnDockerfile embed.FS

// Builder handles Docker image building.
type Builder struct{}

// NewBuilder creates a new Docker builder.
func NewBuilder() (*Builder, error) {
	// Check if Docker is available
	if err := CheckDockerAvailable(); err != nil {
		return nil, err
	}
	return &Builder{}, nil
}

// BuildOptions contains options for building an image.
type BuildOptions struct {
	GitURL   string // Git repository URL
	Checkout string // Branch, tag, or commit to checkout
	Tag      string // Image tag
	NodeType string // Node type (lnd, bitcoind, etc.)
}

// Build builds a Docker image with the given options.
func (b *Builder) Build(ctx context.Context, opts BuildOptions) error {
	// Get the appropriate Dockerfile
	dockerfile, err := b.getDockerfile(opts.NodeType)
	if err != nil {
		return err
	}

	// Create a temporary directory for the build context
	tmpDir, err := os.MkdirTemp("", "aurora-build-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write Dockerfile to temp directory
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, dockerfile, 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}

	// Build the Docker image using CLI
	args := []string{
		"build",
		"--build-arg", fmt.Sprintf("git_url=%s", opts.GitURL),
		"--build-arg", fmt.Sprintf("checkout=%s", opts.Checkout),
		"-t", opts.Tag,
		"-f", dockerfilePath,
		tmpDir,
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start docker build: %w", err)
	}

	// Stream stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Stream stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	return nil
}

// getDockerfile returns the embedded Dockerfile for the given node type.
func (b *Builder) getDockerfile(nodeType string) ([]byte, error) {
	switch nodeType {
	case "lnd":
		content, err := lndDockerfile.ReadFile("dockerfiles/lnd/Dockerfile")
		if err != nil {
			return nil, fmt.Errorf("failed to read LND Dockerfile: %w", err)
		}
		return content, nil
	case "bitcoind":
		content, err := bitcoindDockerfile.ReadFile("dockerfiles/bitcoind/Dockerfile")
		if err != nil {
			return nil, fmt.Errorf("failed to read bitcoind Dockerfile: %w", err)
		}
		return content, nil
	case "cln":
		content, err := clnDockerfile.ReadFile("dockerfiles/cln/Dockerfile")
		if err != nil {
			return nil, fmt.Errorf("failed to read CLN Dockerfile: %w", err)
		}
		return content, nil
	default:
		return nil, fmt.Errorf("unsupported node type: %s (supported: lnd, bitcoind, cln)", nodeType)
	}
}

// CheckDockerAvailable checks if Docker is available and running.
func CheckDockerAvailable() error {
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Docker is not available or not running. Please ensure Docker is installed and the daemon is running")
	}
	return nil
}
