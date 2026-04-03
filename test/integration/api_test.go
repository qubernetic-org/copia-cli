//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apicmd "github.com/qubernetic/copia-cli/pkg/cmd/api"
)

func TestAPI_GetUser(t *testing.T) {
	env := loadTestEnv(t)
	ios, stdout, _ := testIO()

	opts := &apicmd.APIOptions{
		IO:         ios,
		HTTPClient: &http.Client{},
		Host:       env.Host,
		Token:      env.Token,
		Path:       "/user",
	}

	err := apicmd.APIRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "login")
	t.Logf("API /user: %s", stdout.String()[:100])
}

func TestAPI_GetRepo(t *testing.T) {
	env := loadTestEnv(t)
	ios, stdout, _ := testIO()

	opts := &apicmd.APIOptions{
		IO:         ios,
		HTTPClient: &http.Client{},
		Host:       env.Host,
		Token:      env.Token,
		Path:       fmt.Sprintf("/repos/%s/%s", env.Owner, env.Repo),
	}

	err := apicmd.APIRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), env.Repo)
	t.Logf("API repo: %s", stdout.String()[:100])
}

func TestAPI_NotFound(t *testing.T) {
	env := loadTestEnv(t)
	ios, _, _ := testIO()

	opts := &apicmd.APIOptions{
		IO:         ios,
		HTTPClient: &http.Client{},
		Host:       env.Host,
		Token:      env.Token,
		Path:       "/repos/nonexistent/nonexistent",
	}

	err := apicmd.APIRun(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}
