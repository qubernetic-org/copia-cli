package login

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginRun_NonInteractive_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user"),
		httpmock.StringResponse(http.StatusOK, `{"login":"john","id":1}`),
	)

	ios, _, stdout, _ := iostreams.Test()
	configPath := t.TempDir() + "/config.yml"

	opts := &LoginOptions{
		IO:         ios,
		Host:       "app.copia.io",
		Token:      "test-token-123",
		ConfigPath: configPath,
		HTTPClient: &http.Client{Transport: reg},
	}

	err := loginRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Logged in as john")

	cfg, err := config.Load(configPath)
	require.NoError(t, err)
	assert.Equal(t, "test-token-123", cfg.Hosts["app.copia.io"].Token)
	assert.Equal(t, "john", cfg.Hosts["app.copia.io"].User)
}

func TestLoginRun_InvalidToken(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user"),
		httpmock.StringResponse(http.StatusUnauthorized, `{"message":"Unauthorized"}`),
	)

	ios, _, _, _ := iostreams.Test()

	opts := &LoginOptions{
		IO:         ios,
		Host:       "app.copia.io",
		Token:      "bad-token",
		ConfigPath: t.TempDir() + "/config.yml",
		HTTPClient: &http.Client{Transport: reg},
	}

	err := loginRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
}
