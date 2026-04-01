//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepo_ListUserRepos(t *testing.T) {
	env := loadTestEnv(t)

	url := fmt.Sprintf("https://%s/api/v1/user/repos?limit=5", env.Host)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "token "+env.Token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var repos []struct {
		FullName string `json:"full_name"`
	}
	require.NoError(t, json.Unmarshal(body, &repos))
	assert.NotEmpty(t, repos)

	t.Logf("Found %d repos", len(repos))
}

func TestRepo_ViewTestRepo(t *testing.T) {
	env := loadTestEnv(t)

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s", env.Host, env.Owner, env.Repo)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "token "+env.Token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var repo struct {
		FullName      string `json:"full_name"`
		DefaultBranch string `json:"default_branch"`
	}
	require.NoError(t, json.Unmarshal(body, &repo))
	assert.Contains(t, repo.FullName, env.Repo)
	assert.NotEmpty(t, repo.DefaultBranch)

	t.Logf("Repo: %s (default branch: %s)", repo.FullName, repo.DefaultBranch)
}
