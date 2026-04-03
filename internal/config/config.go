// Package config manages the CLI configuration file (~/.config/copia/config.yml),
// including host credentials, default host resolution, and file permissions.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// HostConfig stores credentials for a single Copia/Gitea instance.
type HostConfig struct {
	Token string `yaml:"token"`
	User  string `yaml:"user"`
}

// Config is the top-level configuration.
type Config struct {
	Hosts map[string]*HostConfig `yaml:"hosts"`
}

// DefaultPath returns the default config file path.
func DefaultPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "copia", "config.yml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "copia", "config.yml")
}

// Load reads a config file. Returns empty config if file does not exist.
func Load(path string) (*Config, error) {
	cfg := &Config{Hosts: map[string]*HostConfig{}}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Hosts == nil {
		cfg.Hosts = map[string]*HostConfig{}
	}
	return cfg, nil
}

// Save writes config to path with 0600 permissions.
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// DefaultHost returns the first host entry. Returns empty string and nil if no hosts configured.
func (c *Config) DefaultHost() (string, *HostConfig) {
	for host, hc := range c.Hosts {
		return host, hc
	}
	return "", nil
}
