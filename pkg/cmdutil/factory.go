package cmdutil

import (
	"fmt"

	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

// Factory provides shared dependencies to all commands.
type Factory struct {
	IOStreams *iostreams.IOStreams
	Config   func() (*config.Config, error)
	BaseRepo func() (string, string, error) // returns owner, repo

	// Overrides (flags/env)
	Token string
	Host  string
}

// ResolveAuth returns the host and token, resolving from flags then config.
func (f *Factory) ResolveAuth() (host, token string, err error) {
	host = f.Host
	token = f.Token

	if host == "" || token == "" {
		cfg, err := f.Config()
		if err != nil {
			return "", "", err
		}
		if host == "" {
			h, _ := cfg.DefaultHost()
			host = h
		}
		if token == "" && host != "" {
			if hc, ok := cfg.Hosts[host]; ok {
				token = hc.Token
			}
		}
	}

	if host == "" {
		return "", "", fmt.Errorf("no host configured. Run 'copia auth login' first")
	}
	if token == "" {
		return "", "", fmt.Errorf("no token configured for %s. Run 'copia auth login' first", host)
	}
	return host, token, nil
}
