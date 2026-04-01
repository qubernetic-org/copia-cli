package cmdutil

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

// SplitOwnerRepo parses "owner/repo" into separate components.
func SplitOwnerRepo(nwo string) (string, string, error) {
	i := strings.IndexByte(nwo, '/')
	if i > 0 && i < len(nwo)-1 {
		return nwo[:i], nwo[i+1:], nil
	}
	return "", "", fmt.Errorf("expected owner/repo format, got %q", nwo)
}

// ParseRemoteURL extracts owner and repo from a git remote URL.
// Supports HTTPS (https://host/owner/repo.git) and SSH (git@host:owner/repo.git).
func ParseRemoteURL(rawURL string) (owner, repo string, err error) {
	// SSH format: git@host:owner/repo.git
	if strings.HasPrefix(rawURL, "git@") {
		parts := strings.SplitN(rawURL, ":", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("could not parse SSH remote URL: %s", rawURL)
		}
		path := strings.TrimSuffix(parts[1], ".git")
		return SplitOwnerRepo(path)
	}

	// HTTPS format: https://host/owner/repo.git
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("could not parse remote URL: %s", rawURL)
	}

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, ".git")

	return SplitOwnerRepo(path)
}

// DetectBaseRepo returns a function that detects owner/repo from the git remote "origin".
func DetectBaseRepo() func() (string, string, error) {
	return func() (string, string, error) {
		out, err := exec.Command("git", "remote", "get-url", "origin").Output()
		if err != nil {
			return "", "", fmt.Errorf("not in a git repository or no origin remote configured")
		}
		remoteURL := strings.TrimSpace(string(out))
		return ParseRemoteURL(remoteURL)
	}
}
