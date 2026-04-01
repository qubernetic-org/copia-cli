package cmdutil

import (
	"fmt"
	"os"

	"github.com/qubernetic/copia-cli/internal/config"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

// Factory provides shared dependencies to all commands.
type Factory struct {
	IOStreams *iostreams.IOStreams
	Config   func() (*config.Config, error)
	BaseRepo func() (string, string, error) // returns owner, repo

	// Overrides (flags only — env vars resolved in ResolveAuth)
	Token string
	Host  string
}

// ResolveAuth returns the host and token, resolving in order:
// 1. --token/--host flags (highest priority)
// 2. COPIA_TOKEN/COPIA_HOST env vars
// 3. Config file (lowest priority)
func (f *Factory) ResolveAuth() (host, token string, err error) {
	// 1. Flags (already set by Cobra if provided)
	host = f.Host
	token = f.Token

	// 2. Env vars (only if flag not set)
	if host == "" {
		host = os.Getenv("COPIA_HOST")
	}
	if token == "" {
		token = os.Getenv("COPIA_TOKEN")
	}

	// 3. Config fallback
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
