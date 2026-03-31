package status

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusRun_LoggedIn(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user"),
		httpmock.StringResponse(http.StatusOK, `{"login":"john","id":1}`),
	)

	configPath := t.TempDir() + "/config.yml"
	cfg := &config.Config{
		Hosts: map[string]*config.HostConfig{
			"app.copia.io": {Token: "abc123", User: "john"},
		},
	}
	require.NoError(t, config.Save(configPath, cfg))

	ios, _, stdout, _ := iostreams.Test()

	opts := &StatusOptions{
		IO:         ios,
		ConfigPath: configPath,
		HTTPClient: &http.Client{Transport: reg},
	}

	err := statusRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "app.copia.io")
	assert.Contains(t, stdout.String(), "john")
	assert.Contains(t, stdout.String(), "Token valid")
}

func TestStatusRun_NoHosts(t *testing.T) {
	configPath := t.TempDir() + "/config.yml"
	cfg := &config.Config{Hosts: map[string]*config.HostConfig{}}
	require.NoError(t, config.Save(configPath, cfg))

	ios, _, _, _ := iostreams.Test()

	opts := &StatusOptions{
		IO:         ios,
		ConfigPath: configPath,
	}

	err := statusRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not logged in")
}

func TestStatusRun_InvalidToken(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user"),
		httpmock.StringResponse(http.StatusUnauthorized, `{}`),
	)

	configPath := t.TempDir() + "/config.yml"
	cfg := &config.Config{
		Hosts: map[string]*config.HostConfig{
			"app.copia.io": {Token: "expired", User: "john"},
		},
	}
	require.NoError(t, config.Save(configPath, cfg))

	ios, _, stdout, _ := iostreams.Test()

	opts := &StatusOptions{
		IO:         ios,
		ConfigPath: configPath,
		HTTPClient: &http.Client{Transport: reg},
	}

	err := statusRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Token invalid")
}
