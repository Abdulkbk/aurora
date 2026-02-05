package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is a GitHub API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new GitHub API client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.github.com",
	}
}

// PRDetails contains the information needed to build from a PR.
type PRDetails struct {
	ForkURL string // Clone URL of the fork (e.g., "https://github.com/user/lnd.git")
	Branch  string // Branch name in the fork
	Title   string // PR title for display
	State   string // PR state (open, closed, merged)
}

// prAPIResponse represents the GitHub API response for a PR.
type prAPIResponse struct {
	Title string `json:"title"`
	State string `json:"state"`
	Head  struct {
		Ref  string `json:"ref"` // Branch name
		Repo struct {
			CloneURL string `json:"clone_url"`
		} `json:"repo"`
	} `json:"head"`
}

// GetPRDetails fetches PR details from the GitHub API.
func (c *Client) GetPRDetails(owner, repo string, prNumber int) (*PRDetails, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.baseURL, owner, repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "aurora-cli")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("PR #%d not found in %s/%s", prNumber, owner, repo)
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("GitHub API rate limit exceeded. Try again later or use --repo/--branch flags")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var apiResp prAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	if apiResp.Head.Repo.CloneURL == "" {
		return nil, fmt.Errorf("PR fork repository not available (may have been deleted)")
	}

	return &PRDetails{
		ForkURL: apiResp.Head.Repo.CloneURL,
		Branch:  apiResp.Head.Ref,
		Title:   apiResp.Title,
		State:   apiResp.State,
	}, nil
}
