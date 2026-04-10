package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPRDetailsWithDeletedFork(t *testing.T) {
	tests := []struct {
		name          string
		responseBody  string
		expectedError string
	}{
		{
			name: "fork repository deleted (repo is null)",
			responseBody: `{
				"title": "Test PR",
				"state": "open",
				"head": {
					"ref": "feature-branch",
					"repo": null
				}
			}`,
			expectedError: "fork repository not available — the contributor may have deleted their fork",
		},
		{
			name: "fork repository deleted (empty clone_url)",
			responseBody: `{
				"title": "Test PR",
				"state": "open",
				"head": {
					"ref": "feature-branch",
					"repo": {
						"clone_url": ""
					}
				}
			}`,
			expectedError: "fork repository not available — the contributor may have deleted their fork",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that returns the mock response
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create client with test server URL
			client := &Client{
				httpClient: server.Client(),
				baseURL:    server.URL,
			}

			// Test the function
			_, err := client.GetPRDetails("owner", "repo", 123)

			if err == nil {
				t.Errorf("expected error, got nil")
				return
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestGetPRDetailsSuccess(t *testing.T) {
	responseBody := `{
		"title": "Test PR Title",
		"state": "open",
		"head": {
			"ref": "feature-branch",
			"repo": {
				"clone_url": "https://github.com/user/repo.git"
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	details, err := client.GetPRDetails("owner", "repo", 123)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := &PRDetails{
		ForkURL: "https://github.com/user/repo.git",
		Branch:  "feature-branch",
		Title:   "Test PR Title",
		State:   "open",
	}

	if details.ForkURL != expected.ForkURL {
		t.Errorf("expected ForkURL %q, got %q", expected.ForkURL, details.ForkURL)
	}
	if details.Branch != expected.Branch {
		t.Errorf("expected Branch %q, got %q", expected.Branch, details.Branch)
	}
	if details.Title != expected.Title {
		t.Errorf("expected Title %q, got %q", expected.Title, details.Title)
	}
	if details.State != expected.State {
		t.Errorf("expected State %q, got %q", expected.State, details.State)
	}
}
