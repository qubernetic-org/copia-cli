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

func TestPR_List(t *testing.T) {
	env := loadTestEnv(t)

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/pulls?state=all&limit=5",
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

	var prs []struct {
		Number int64 `json:"number"`
	}
	require.NoError(t, json.Unmarshal(body, &prs))

	t.Logf("Found %d pull requests", len(prs))
}
