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

func TestLabel_CreateListDelete(t *testing.T) {
	env := loadTestEnv(t)
	baseURL := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/labels", env.Host, env.Owner, env.Repo)

	// Create label
	createPayload, _ := json.Marshal(map[string]string{
		"name":        "integration-test-label",
		"color":       "#ff0000",
		"description": "Created by integration test",
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
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	require.NoError(t, json.Unmarshal(createBody, &created))
	assert.Equal(t, "integration-test-label", created.Name)
	assert.Greater(t, created.ID, int64(0))

	t.Logf("Created label: %s (ID: %d)", created.Name, created.ID)

	// Cleanup: delete label
	defer func() {
		deleteURL := fmt.Sprintf("%s/%d", baseURL, created.ID)
		deleteReq, err := http.NewRequest("DELETE", deleteURL, nil)
		if err != nil {
			t.Logf("Failed to create delete request: %v", err)
			return
		}
		deleteReq.Header.Set("Authorization", "token "+env.Token)

		deleteResp, err := http.DefaultClient.Do(deleteReq)
		if err != nil {
			t.Logf("Failed to delete label: %v", err)
			return
		}
		_ = deleteResp.Body.Close()
		t.Logf("Deleted label ID: %d (HTTP %d)", created.ID, deleteResp.StatusCode)
	}()

	// List labels and verify our label exists
	listReq, err := http.NewRequest("GET", baseURL, nil)
	require.NoError(t, err)
	listReq.Header.Set("Authorization", "token "+env.Token)

	listResp, err := http.DefaultClient.Do(listReq)
	require.NoError(t, err)
	defer func() { _ = listResp.Body.Close() }()

	assert.Equal(t, http.StatusOK, listResp.StatusCode)

	listBody, err := io.ReadAll(listResp.Body)
	require.NoError(t, err)

	var labels []struct {
		Name string `json:"name"`
	}
	require.NoError(t, json.Unmarshal(listBody, &labels))

	found := false
	for _, l := range labels {
		if l.Name == "integration-test-label" {
			found = true
			break
		}
	}
	assert.True(t, found, "created label not found in list")
}
