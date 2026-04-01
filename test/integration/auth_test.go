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

func TestAuth_ValidateToken(t *testing.T) {
	env := loadTestEnv(t)

	url := fmt.Sprintf("https://%s/api/v1/user", env.Host)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "token "+env.Token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var user struct {
		Login string `json:"login"`
		ID    int64  `json:"id"`
	}
	require.NoError(t, json.Unmarshal(body, &user))
	assert.NotEmpty(t, user.Login)
	assert.Greater(t, user.ID, int64(0))

	t.Logf("Authenticated as: %s (ID: %d)", user.Login, user.ID)
}

func TestAuth_InvalidToken(t *testing.T) {
	env := loadTestEnv(t)

	url := fmt.Sprintf("https://%s/api/v1/user", env.Host)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "token invalid-token-that-does-not-exist")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Copia/Gitea returns 403 for invalid tokens (not 401)
	assert.Contains(t, []int{http.StatusUnauthorized, http.StatusForbidden}, resp.StatusCode)
}
