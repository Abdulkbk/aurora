package docker

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ImageInfo holds information about a Docker image built by Aurora.
type ImageInfo struct {
	Repository   string
	Tag          string
	CreatedSince string
}

// ListAuroraImages returns all Docker images tagged with the ":aurora" suffix.
func ListAuroraImages() ([]ImageInfo, error) {
	cmd := exec.Command(
		"docker", "images",
		"--filter", "reference=*:aurora",
		"--format", "{{.Repository}}\t{{.Tag}}\t{{.CreatedSince}}",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list Docker images: %w", err)
	}

	var images []ImageInfo
	for line := range strings.SplitSeq(
		strings.TrimSpace(out.String()), "\n") {

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			continue
		}

		images = append(images, ImageInfo{
			Repository:   parts[0],
			Tag:          parts[1],
			CreatedSince: parts[2],
		})
	}

	return images, nil
}
