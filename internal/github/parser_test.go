package github

import (
	"testing"
)

func TestParsePRURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *PRInfo
		wantErr bool
	}{
		{
			name: "valid PR URL",
			url:  "https://github.com/lightningnetwork/lnd/pull/1234",
			want: &PRInfo{
				Owner:    "lightningnetwork",
				Repo:     "lnd",
				PRNumber: 1234,
			},
			wantErr: false,
		},
		{
			name: "valid PR URL with trailing slash",
			url:  "https://github.com/bitcoin/bitcoin/pull/5678/",
			want: &PRInfo{
				Owner:    "bitcoin",
				Repo:     "bitcoin",
				PRNumber: 5678,
			},
			wantErr: false,
		},
		{
			name:    "invalid URL - not a PR",
			url:     "https://github.com/lightningnetwork/lnd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - wrong domain",
			url:     "https://gitlab.com/lightningnetwork/lnd/pull/1234",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - missing PR number",
			url:     "https://github.com/lightningnetwork/lnd/pull/",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePRURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePRURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if got.Owner != tt.want.Owner || got.Repo != tt.want.Repo || got.PRNumber != tt.want.PRNumber {
					t.Errorf("ParsePRURL() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *RepoInfo
		wantErr bool
	}{
		{
			name: "valid repo URL",
			url:  "https://github.com/myuser/lnd",
			want: &RepoInfo{
				Owner: "myuser",
				Repo:  "lnd",
			},
			wantErr: false,
		},
		{
			name: "valid repo URL with trailing slash",
			url:  "https://github.com/myuser/lnd/",
			want: &RepoInfo{
				Owner: "myuser",
				Repo:  "lnd",
			},
			wantErr: false,
		},
		{
			name:    "invalid URL - includes path",
			url:     "https://github.com/myuser/lnd/tree/main",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRepoURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRepoURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if got.Owner != tt.want.Owner || got.Repo != tt.want.Repo {
					t.Errorf("ParseRepoURL() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}
