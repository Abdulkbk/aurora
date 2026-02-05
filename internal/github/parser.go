// Package github provides utilities for interacting with GitHub repositories and PRs.
package github

import (
	"fmt"
	"regexp"
	"strconv"
)

// PRInfo contains information extracted from a GitHub PR URL.
type PRInfo struct {
	Owner    string // Repository owner (e.g., "lightningnetwork")
	Repo     string // Repository name (e.g., "lnd")
	PRNumber int    // Pull request number (e.g., 1234)
}

// RepoInfo contains information extracted from a GitHub repository URL.
type RepoInfo struct {
	Owner  string // Repository owner
	Repo   string // Repository name
	Branch string // Branch name
}

// prURLPattern matches GitHub PR URLs like:
// https://github.com/owner/repo/pull/123
var prURLPattern = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/pull/(\d+)/?$`)

// repoURLPattern matches GitHub repository URLs like:
// https://github.com/owner/repo
var repoURLPattern = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/?$`)

// ParsePRURL parses a GitHub PR URL and extracts owner, repo, and PR number.
func ParsePRURL(url string) (*PRInfo, error) {
	matches := prURLPattern.FindStringSubmatch(url)
	if matches == nil {
		return nil, fmt.Errorf("invalid GitHub PR URL: %s\nExpected format: https://github.com/owner/repo/pull/123", url)
	}

	prNum, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid PR number: %s", matches[3])
	}

	return &PRInfo{
		Owner:    matches[1],
		Repo:     matches[2],
		PRNumber: prNum,
	}, nil
}

// ParseRepoURL parses a GitHub repository URL and extracts owner and repo.
func ParseRepoURL(url string) (*RepoInfo, error) {
	matches := repoURLPattern.FindStringSubmatch(url)
	if matches == nil {
		return nil, fmt.Errorf("invalid GitHub repository URL: %s\nExpected format: https://github.com/owner/repo", url)
	}

	return &RepoInfo{
		Owner: matches[1],
		Repo:  matches[2],
	}, nil
}

// CloneURL returns the HTTPS clone URL for the repository.
func (r *RepoInfo) CloneURL() string {
	return fmt.Sprintf("https://github.com/%s/%s.git", r.Owner, r.Repo)
}
