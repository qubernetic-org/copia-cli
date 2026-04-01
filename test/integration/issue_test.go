//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssue_FullLifecycle(t *testing.T) {
	env := loadTestEnv(t)
	baseURL := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues", env.Host, env.Owner, env.Repo)

	// Create issue
	createPayload, _ := json.Marshal(map[string]string{
		"title": "[integration-test] lifecycle test",
		"body":  "Created by integration test. Will be closed automatically.",
	})

	createReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(createPayload))
	require.NoError(t, err)
	createReq.Header.Set("Authorization", "token "+env.Token)
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := http.DefaultClient.Do(createReq)
	require.NoError(t, err)
	defer func() { _ = createResp.Body.Close() }()

	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	createBody, err := io.ReadAll(createResp.Body)
	require.NoError(t, err)

	var created struct {
		Number int64  `json:"number"`
		Title  string `json:"title"`
		State  string `json:"state"`
	}
	require.NoError(t, json.Unmarshal(createBody, &created))
	assert.Equal(t, "[integration-test] lifecycle test", created.Title)
	assert.Equal(t, "open", created.State)

	t.Logf("Created issue #%d", created.Number)

	// Cleanup: close issue at end
	defer func() {
		closePayload, _ := json.Marshal(map[string]string{"state": "closed"})
		closeURL := fmt.Sprintf("%s/%d", baseURL, created.Number)
		closeReq, err := http.NewRequest("PATCH", closeURL, bytes.NewReader(closePayload))
		if err != nil {
			t.Logf("Failed to create close request: %v", err)
			return
		}
		closeReq.Header.Set("Authorization", "token "+env.Token)
		closeReq.Header.Set("Content-Type", "application/json")

		closeResp, err := http.DefaultClient.Do(closeReq)
		if err != nil {
			t.Logf("Failed to close issue: %v", err)
			return
		}
		_ = closeResp.Body.Close()
		t.Logf("Closed issue #%d (HTTP %d)", created.Number, closeResp.StatusCode)
	}()

	// View issue
	viewURL := fmt.Sprintf("%s/%d", baseURL, created.Number)
	viewReq, err := http.NewRequest("GET", viewURL, nil)
	require.NoError(t, err)
	viewReq.Header.Set("Authorization", "token "+env.Token)

	viewResp, err := http.DefaultClient.Do(viewReq)
	require.NoError(t, err)
	defer func() { _ = viewResp.Body.Close() }()

	assert.Equal(t, http.StatusOK, viewResp.StatusCode)

	viewBody, err := io.ReadAll(viewResp.Body)
	require.NoError(t, err)

	var viewed struct {
		Number int64  `json:"number"`
		Title  string `json:"title"`
	}
	require.NoError(t, json.Unmarshal(viewBody, &viewed))
	assert.Equal(t, created.Number, viewed.Number)

	// Add comment
	commentPayload, _ := json.Marshal(map[string]string{
		"body": "Integration test comment — will be cleaned up.",
	})
	commentURL := fmt.Sprintf("%s/%d/comments", baseURL, created.Number)
	commentReq, err := http.NewRequest("POST", commentURL, bytes.NewReader(commentPayload))
	require.NoError(t, err)
	commentReq.Header.Set("Authorization", "token "+env.Token)
	commentReq.Header.Set("Content-Type", "application/json")

	commentResp, err := http.DefaultClient.Do(commentReq)
	require.NoError(t, err)
	defer func() { _ = commentResp.Body.Close() }()

	assert.Equal(t, http.StatusCreated, commentResp.StatusCode)

	t.Logf("Added comment to issue #%d", created.Number)
}

func TestIssue_List(t *testing.T) {
	env := loadTestEnv(t)

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues?state=all&limit=5&type=issues",
		env.Host, env.Owner, env.Repo)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "token "+env.Token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var issues []struct {
		Number int64 `json:"number"`
	}
	require.NoError(t, json.Unmarshal(body, &issues))

	t.Logf("Found %d issues", len(issues))
}
