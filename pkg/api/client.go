// Package api provides a Gitea SDK client wrapper for the Copia REST API.
package api

import (
	"fmt"
	"net/http"

	"code.gitea.io/sdk/gitea"
)

// NewClient creates a Gitea SDK client for the given host and token.
func NewClient(host, token string) (*gitea.Client, error) {
	if host == "" {
		return nil, fmt.Errorf("host is required")
	}
	url := "https://" + host
	return gitea.NewClient(url, gitea.SetToken(token))
}

// NewClientWithHTTP creates a Gitea SDK client with a custom HTTP client.
// Useful for testing with mocked HTTP transport.
func NewClientWithHTTP(host, token string, httpClient *http.Client) (*gitea.Client, error) {
	if host == "" {
		return nil, fmt.Errorf("host is required")
	}
	url := "https://" + host
	return gitea.NewClient(url, gitea.SetToken(token), gitea.SetHTTPClient(httpClient))
}
