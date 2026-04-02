package cmdutil

import (
	"testing"

	"github.com/qubernetic/copia-cli/internal/config"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
)

func TestFactory_ResolveAuth_FlagOverride(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	f := &Factory{
		IOStreams: ios,
		Token:    "flag-token",
		Host:     "flag-host.io",
		Config: func() (*config.Config, error) {
			return &config.Config{Hosts: map[string]*config.HostConfig{}}, nil
		},
	}

	host, token, err := f.ResolveAuth()
	assert.NoError(t, err)
	assert.Equal(t, "flag-host.io", host)
	assert.Equal(t, "flag-token", token)
}

func TestFactory_ResolveAuth_EnvOverride(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	t.Setenv("COPIA_TOKEN", "env-token")
	t.Setenv("COPIA_HOST", "env-host.io")

	f := &Factory{
		IOStreams: ios,
		Config: func() (*config.Config, error) {
			return &config.Config{Hosts: map[string]*config.HostConfig{
				"config-host.io": {Token: "config-token"},
			}}, nil
		},
	}

	host, token, err := f.ResolveAuth()
	assert.NoError(t, err)
	assert.Equal(t, "env-host.io", host)
	assert.Equal(t, "env-token", token)
}

func TestFactory_ResolveAuth_FlagBeatsEnv(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	t.Setenv("COPIA_TOKEN", "env-token")
	t.Setenv("COPIA_HOST", "env-host.io")

	f := &Factory{
		IOStreams: ios,
		Token:    "flag-token",
		Host:     "flag-host.io",
		Config: func() (*config.Config, error) {
			return &config.Config{Hosts: map[string]*config.HostConfig{}}, nil
		},
	}

	host, token, err := f.ResolveAuth()
	assert.NoError(t, err)
	assert.Equal(t, "flag-host.io", host)
	assert.Equal(t, "flag-token", token)
}

func TestFactory_ResolveAuth_ConfigFallback(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	f := &Factory{
		IOStreams: ios,
		Config: func() (*config.Config, error) {
			return &config.Config{
				Hosts: map[string]*config.HostConfig{
					"app.copia.io": {Token: "config-token", User: "john"},
				},
			}, nil
		},
	}

	host, token, err := f.ResolveAuth()
	assert.NoError(t, err)
	assert.Equal(t, "app.copia.io", host)
	assert.Equal(t, "config-token", token)
}

func TestFactory_ResolveAuth_NoConfig(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	f := &Factory{
		IOStreams: ios,
		Config: func() (*config.Config, error) {
			return &config.Config{Hosts: map[string]*config.HostConfig{}}, nil
		},
	}

	_, _, err := f.ResolveAuth()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no host configured")
}

func TestFactory_ResolveRepo_FlagOverride(t *testing.T) {
	f := &Factory{
		Repo: "flag-org/flag-repo",
		BaseRepo: func() (string, string, error) {
			return "git-org", "git-repo", nil
		},
	}

	owner, repo, err := f.ResolveRepo()
	assert.NoError(t, err)
	assert.Equal(t, "flag-org", owner)
	assert.Equal(t, "flag-repo", repo)
}

func TestFactory_ResolveRepo_BaseRepoFallback(t *testing.T) {
	f := &Factory{
		BaseRepo: func() (string, string, error) {
			return "git-org", "git-repo", nil
		},
	}

	owner, repo, err := f.ResolveRepo()
	assert.NoError(t, err)
	assert.Equal(t, "git-org", owner)
	assert.Equal(t, "git-repo", repo)
}

func TestFactory_ResolveRepo_NoRepoNoBaseRepo(t *testing.T) {
	f := &Factory{}

	_, _, err := f.ResolveRepo()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not determine repository")
}

func TestFactory_ResolveRepo_InvalidFormat(t *testing.T) {
	f := &Factory{
		Repo: "invalid-no-slash",
	}

	_, _, err := f.ResolveRepo()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected owner/repo format")
}

func TestParseRemoteURL_HTTPS(t *testing.T) {
	owner, repo, err := ParseRemoteURL("https://app.copia.io/my-org/my-repo.git")
	assert.NoError(t, err)
	assert.Equal(t, "my-org", owner)
	assert.Equal(t, "my-repo", repo)
}

func TestParseRemoteURL_HTTPS_NoGitSuffix(t *testing.T) {
	owner, repo, err := ParseRemoteURL("https://app.copia.io/my-org/my-repo")
	assert.NoError(t, err)
	assert.Equal(t, "my-org", owner)
	assert.Equal(t, "my-repo", repo)
}

func TestParseRemoteURL_SSH(t *testing.T) {
	owner, repo, err := ParseRemoteURL("git@app.copia.io:my-org/my-repo.git")
	assert.NoError(t, err)
	assert.Equal(t, "my-org", owner)
	assert.Equal(t, "my-repo", repo)
}

func TestParseRemoteURL_Invalid(t *testing.T) {
	_, _, err := ParseRemoteURL("not-a-url")
	assert.Error(t, err)
}
