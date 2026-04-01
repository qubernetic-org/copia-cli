package logout

import (
	"testing"

	"github.com/qubernetic/copia-cli/internal/config"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutRun_RemovesHost(t *testing.T) {
	configPath := t.TempDir() + "/config.yml"
	cfg := &config.Config{
		Hosts: map[string]*config.HostConfig{
			"app.copia.io": {Token: "abc", User: "john"},
		},
	}
	require.NoError(t, config.Save(configPath, cfg))

	ios, _, stdout, _ := iostreams.Test()

	opts := &LogoutOptions{
		IO:         ios,
		Host:       "app.copia.io",
		ConfigPath: configPath,
	}

	err := logoutRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Logged out of app.copia.io")

	loaded, _ := config.Load(configPath)
	assert.Empty(t, loaded.Hosts)
}

func TestLogoutRun_HostNotFound(t *testing.T) {
	configPath := t.TempDir() + "/config.yml"
	cfg := &config.Config{Hosts: map[string]*config.HostConfig{}}
	require.NoError(t, config.Save(configPath, cfg))

	ios, _, _, _ := iostreams.Test()

	opts := &LogoutOptions{
		IO:         ios,
		Host:       "unknown.host.com",
		ConfigPath: configPath,
	}

	err := logoutRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not logged in to unknown.host.com")
}
