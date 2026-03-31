package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_NonExistentFile(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(filepath.Join(dir, "config.yml"))
	require.NoError(t, err)
	assert.Empty(t, cfg.Hosts)
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")

	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"app.copia.io": {
				Token: "abc123",
				User:  "john",
			},
		},
	}

	err := Save(path, cfg)
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	loaded, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "abc123", loaded.Hosts["app.copia.io"].Token)
	assert.Equal(t, "john", loaded.Hosts["app.copia.io"].User)
}

func TestSaveAndLoad_MultipleHosts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")

	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"app.copia.io":     {Token: "t1", User: "u1"},
			"on-prem.corp.com": {Token: "t2", User: "u2"},
		},
	}

	require.NoError(t, Save(path, cfg))

	loaded, err := Load(path)
	require.NoError(t, err)
	assert.Len(t, loaded.Hosts, 2)
	assert.Equal(t, "t1", loaded.Hosts["app.copia.io"].Token)
	assert.Equal(t, "t2", loaded.Hosts["on-prem.corp.com"].Token)
}

func TestConfig_DefaultHost(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"app.copia.io": {Token: "t1", User: "u1"},
		},
	}
	host, hc := cfg.DefaultHost()
	assert.Equal(t, "app.copia.io", host)
	assert.Equal(t, "t1", hc.Token)
}

func TestConfig_DefaultHost_Empty(t *testing.T) {
	cfg := &Config{Hosts: map[string]*HostConfig{}}
	host, hc := cfg.DefaultHost()
	assert.Empty(t, host)
	assert.Nil(t, hc)
}

func TestSave_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "config.yml")

	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"app.copia.io": {Token: "abc", User: "john"},
		},
	}

	err := Save(path, cfg)
	require.NoError(t, err)

	loaded, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "abc", loaded.Hosts["app.copia.io"].Token)
}

func TestDefaultPath_RespectsXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	path := DefaultPath()
	assert.Equal(t, "/custom/config/copia/config.yml", path)
}
