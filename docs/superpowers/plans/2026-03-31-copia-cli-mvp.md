# Copia CLI MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a working Copia CLI that provides auth, repo, issue, PR, and label management against the Gitea-compatible Copia API.

**Architecture:** Go CLI using Cobra for command routing, Gitea Go SDK for API calls, gh-style factory injection for testability. Domain-driven command packages under `pkg/cmd/`, shared utilities in `pkg/cmdutil/` and `pkg/iostreams/`.

**Tech Stack:** Go 1.23+, Cobra, Gitea SDK (`code.gitea.io/sdk/gitea`), testify, GoReleaser

---

## File Structure

```
copia-cli/
├── cmd/copia/main.go                          # Entrypoint
├── internal/
│   ├── build/build.go                         # Version injection
│   ├── config/config.go                       # YAML config management
│   └── copiacmd/root.go                       # Root command factory
├── pkg/
│   ├── api/client.go                          # Gitea SDK wrapper
│   ├── cmd/
│   │   ├── auth/
│   │   │   ├── auth.go                        # Parent command
│   │   │   ├── login/login.go                 # Login subcommand
│   │   │   ├── login/login_test.go
│   │   │   ├── logout/logout.go
│   │   │   ├── logout/logout_test.go
│   │   │   ├── status/status.go
│   │   │   └── status/status_test.go
│   │   ├── repo/
│   │   │   ├── repo.go
│   │   │   ├── list/list.go
│   │   │   ├── list/list_test.go
│   │   │   ├── view/view.go
│   │   │   ├── view/view_test.go
│   │   │   ├── clone/clone.go
│   │   │   └── clone/clone_test.go
│   │   ├── issue/
│   │   │   ├── issue.go
│   │   │   ├── list/list.go
│   │   │   ├── list/list_test.go
│   │   │   ├── list/fixtures/issueList.json
│   │   │   ├── create/create.go
│   │   │   ├── create/create_test.go
│   │   │   ├── view/view.go
│   │   │   ├── view/view_test.go
│   │   │   ├── view/fixtures/issueView.json
│   │   │   ├── close/close.go
│   │   │   ├── close/close_test.go
│   │   │   ├── comment/comment.go
│   │   │   └── comment/comment_test.go
│   │   ├── pr/
│   │   │   ├── pr.go
│   │   │   ├── list/list.go
│   │   │   ├── list/list_test.go
│   │   │   ├── list/fixtures/prList.json
│   │   │   ├── create/create.go
│   │   │   ├── create/create_test.go
│   │   │   ├── view/view.go
│   │   │   ├── view/view_test.go
│   │   │   ├── view/fixtures/prView.json
│   │   │   ├── merge/merge.go
│   │   │   ├── merge/merge_test.go
│   │   │   ├── close/close.go
│   │   │   └── close/close_test.go
│   │   └── label/
│   │       ├── label.go
│   │       ├── list/list.go
│   │       ├── list/list_test.go
│   │       ├── create/create.go
│   │       └── create/create_test.go
│   ├── cmdutil/
│   │   ├── factory.go                         # Dependency injection
│   │   ├── flags.go                           # Common flag helpers
│   │   └── json.go                            # --json flag support
│   ├── iostreams/
│   │   ├── iostreams.go                       # TTY-aware I/O
│   │   └── iostreams_test.go
│   └── httpmock/
│       ├── registry.go                        # HTTP transport mock
│       └── registry_test.go
├── go.mod
├── go.sum
├── Makefile
├── .goreleaser.yml
└── .github/workflows/ci.yml
```

---

### Task 1: Project Skeleton

**Files:**
- Create: `go.mod`
- Create: `cmd/copia/main.go`
- Create: `internal/build/build.go`
- Create: `Makefile`

- [ ] **Step 1: Initialize Go module**

Run:
```bash
cd /home/cbiro/Git/copia-cli
go mod init github.com/qubernetic-org/copia-cli
```

Expected: `go.mod` created with module path.

- [ ] **Step 2: Create internal/build/build.go**

```go
// internal/build/build.go
package build

import "runtime/debug"

// Version is set at build time via ldflags.
var Version = "DEV"

// Date is set at build time via ldflags.
var Date = ""

func init() {
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
```

- [ ] **Step 3: Create minimal entrypoint cmd/copia/main.go**

```go
// cmd/copia/main.go
package main

import (
	"fmt"
	"os"

	"github.com/qubernetic-org/copia-cli/internal/build"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("copia version %s (%s)\n", build.Version, build.Date)
		os.Exit(0)
	}
	fmt.Println("copia: command-line interface for Copia")
	os.Exit(0)
}
```

- [ ] **Step 4: Create Makefile**

```makefile
# Makefile
BIN := copia
VERSION ?= DEV
DATE := $(shell date -u +%Y-%m-%d)
LDFLAGS := -s -w \
	-X github.com/qubernetic-org/copia-cli/internal/build.Version=$(VERSION) \
	-X github.com/qubernetic-org/copia-cli/internal/build.Date=$(DATE)

.PHONY: build test integration acceptance clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BIN) ./cmd/copia

test:
	go test ./...

integration:
	go test -tags=integration ./...

acceptance:
	go test -tags=acceptance ./acceptance/...

clean:
	rm -rf bin/
```

- [ ] **Step 5: Build and verify**

Run:
```bash
make build && ./bin/copia --version
```

Expected: `copia version DEV (2026-03-31)`

- [ ] **Step 6: Commit**

```bash
git add go.mod cmd/ internal/build/ Makefile
git commit -m "feat: initialize project skeleton with Go module, entrypoint, and Makefile"
```

---

### Task 2: IOStreams Abstraction

**Files:**
- Create: `pkg/iostreams/iostreams.go`
- Create: `pkg/iostreams/iostreams_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/iostreams/iostreams_test.go
package iostreams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest_ReturnsWorkingStreams(t *testing.T) {
	ios, stdin, stdout, stderr := Test()

	stdin.WriteString("input data")
	assert.NotNil(t, ios)
	assert.Equal(t, "input data", stdin.String())
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())
	assert.False(t, ios.IsStdoutTTY())
	assert.False(t, ios.IsStdinTTY())
}

func TestIOStreams_TTYDetection(t *testing.T) {
	ios, _, _, _ := Test()

	ios.SetStdoutTTY(true)
	assert.True(t, ios.IsStdoutTTY())

	ios.SetStdoutTTY(false)
	assert.False(t, ios.IsStdoutTTY())
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/iostreams/ -v`
Expected: FAIL — package does not exist yet.

- [ ] **Step 3: Write implementation**

```go
// pkg/iostreams/iostreams.go
package iostreams

import (
	"bytes"
	"io"
	"os"
)

// IOStreams provides TTY-aware I/O for commands.
type IOStreams struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer

	stdinIsTTY  bool
	stdoutIsTTY bool
	stderrIsTTY bool
}

// System returns IOStreams connected to real stdin/stdout/stderr.
func System() *IOStreams {
	return &IOStreams{
		In:          os.Stdin,
		Out:         os.Stdout,
		ErrOut:      os.Stderr,
		stdinIsTTY:  isTerminal(os.Stdin),
		stdoutIsTTY: isTerminal(os.Stdout),
		stderrIsTTY: isTerminal(os.Stderr),
	}
}

// Test returns IOStreams with in-memory buffers for testing.
func Test() (*IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	return &IOStreams{
		In:     io.NopCloser(stdin),
		Out:    stdout,
		ErrOut: stderr,
	}, stdin, stdout, stderr
}

func (s *IOStreams) IsStdinTTY() bool  { return s.stdinIsTTY }
func (s *IOStreams) IsStdoutTTY() bool { return s.stdoutIsTTY }
func (s *IOStreams) IsStderrTTY() bool { return s.stderrIsTTY }

func (s *IOStreams) SetStdinTTY(v bool)  { s.stdinIsTTY = v }
func (s *IOStreams) SetStdoutTTY(v bool) { s.stdoutIsTTY = v }
func (s *IOStreams) SetStderrTTY(v bool) { s.stderrIsTTY = v }

func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./pkg/iostreams/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/iostreams/
git commit -m "feat: add iostreams TTY-aware I/O abstraction"
```

---

### Task 3: Config Management

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing test**

```go
// internal/config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Empty(t *testing.T) {
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

	// Verify file permissions
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	loaded, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "abc123", loaded.Hosts["app.copia.io"].Token)
	assert.Equal(t, "john", loaded.Hosts["app.copia.io"].User)
}

func TestConfig_DefaultHost(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"first.example.com": {Token: "t1", User: "u1"},
			"second.example.com": {Token: "t2", User: "u2"},
		},
	}
	host, hc := cfg.DefaultHost()
	assert.NotEmpty(t, host)
	assert.NotNil(t, hc)
}

func TestConfig_DefaultHost_Empty(t *testing.T) {
	cfg := &Config{Hosts: map[string]*HostConfig{}}
	host, hc := cfg.DefaultHost()
	assert.Empty(t, host)
	assert.Nil(t, hc)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/config/ -v`
Expected: FAIL — package does not exist yet.

- [ ] **Step 3: Write implementation**

```go
// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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
	if runtime.GOOS == "windows" {
		return filepath.Join(home, ".config", "copia", "config.yml")
	}
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
```

- [ ] **Step 4: Run tests to verify they pass**

Run:
```bash
go get gopkg.in/yaml.v3
go test ./internal/config/ -v
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/config/ go.mod go.sum
git commit -m "feat: add YAML config management with multi-host support"
```

---

### Task 4: HTTP Mock for Testing

**Files:**
- Create: `pkg/httpmock/registry.go`
- Create: `pkg/httpmock/registry_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/httpmock/registry_test.go
package httpmock

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_MatchesAndResponds(t *testing.T) {
	reg := &Registry{}

	reg.Register(
		REST("GET", "/api/v1/repos/owner/repo"),
		StringResponse(http.StatusOK, `{"name":"repo"}`),
	)

	req, _ := http.NewRequest("GET", "https://app.copia.io/api/v1/repos/owner/repo", nil)
	resp, err := reg.RoundTrip(req)
	require.NoError(t, err)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"name":"repo"}`, string(body))
}

func TestRegistry_Verify_AllCalled(t *testing.T) {
	mockT := &testing.T{}
	reg := &Registry{}

	reg.Register(
		REST("GET", "/api/v1/user"),
		StringResponse(http.StatusOK, `{"login":"john"}`),
	)

	// Not calling RoundTrip — Verify should report failure
	reg.Verify(mockT)
	// mockT would have failed, but we can't easily assert that in a unit test.
	// This test documents the contract.
}

func TestRegistry_NoMatch_ReturnsError(t *testing.T) {
	reg := &Registry{}

	req, _ := http.NewRequest("GET", "https://app.copia.io/api/v1/unknown", nil)
	_, err := reg.RoundTrip(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no mock matched")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/httpmock/ -v`
Expected: FAIL — package does not exist.

- [ ] **Step 3: Write implementation**

```go
// pkg/httpmock/registry.go
package httpmock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
)

// Matcher checks whether an HTTP request matches a stub.
type Matcher func(req *http.Request) bool

// Responder produces a response for a matched request.
type Responder func(req *http.Request) (*http.Response, error)

type stub struct {
	matcher Matcher
	respond Responder
	called  bool
}

// Registry is an http.RoundTripper that returns stubbed responses.
type Registry struct {
	mu    sync.Mutex
	stubs []*stub
}

// Register adds a matcher/responder pair.
func (r *Registry) Register(m Matcher, resp Responder) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stubs = append(r.stubs, &stub{matcher: m, respond: resp})
}

// RoundTrip implements http.RoundTripper.
func (r *Registry) RoundTrip(req *http.Request) (*http.Response, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.stubs {
		if s.matcher(req) {
			s.called = true
			return s.respond(req)
		}
	}
	return nil, fmt.Errorf("no mock matched for %s %s", req.Method, req.URL.Path)
}

// Verify asserts all registered stubs were called.
func (r *Registry) Verify(t *testing.T) {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.stubs {
		if !s.called {
			t.Errorf("httpmock: registered stub was not called")
		}
	}
}

// REST returns a matcher for a REST API call by method and path suffix.
func REST(method, pathSuffix string) Matcher {
	return func(req *http.Request) bool {
		return req.Method == method && strings.HasSuffix(req.URL.Path, pathSuffix)
	}
}

// StringResponse returns a responder that sends a fixed status and body.
func StringResponse(status int, body string) Responder {
	return func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./pkg/httpmock/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/httpmock/
git commit -m "feat: add HTTP mock registry for testing"
```

---

### Task 5: Factory and JSON Utilities

**Files:**
- Create: `pkg/cmdutil/factory.go`
- Create: `pkg/cmdutil/json.go`
- Create: `pkg/cmdutil/flags.go`
- Create: `pkg/api/client.go`

- [ ] **Step 1: Create pkg/api/client.go**

```go
// pkg/api/client.go
package api

import (
	"code.gitea.io/sdk/gitea"
)

// NewClient creates a Gitea SDK client for the given host and token.
func NewClient(host, token string) (*gitea.Client, error) {
	url := "https://" + host
	return gitea.NewClient(url, gitea.SetToken(token))
}
```

- [ ] **Step 2: Create pkg/cmdutil/factory.go**

```go
// pkg/cmdutil/factory.go
package cmdutil

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/api"
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

// Client creates a Gitea API client using resolved host and token.
func (f *Factory) Client() (*gitea.Client, error) {
	host, token, err := f.resolveAuth()
	if err != nil {
		return nil, err
	}
	return api.NewClient(host, token)
}

func (f *Factory) resolveAuth() (host, token string, err error) {
	// 1. Flag overrides
	host = f.Host
	token = f.Token

	// 2. Config fallback
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
```

- [ ] **Step 3: Create pkg/cmdutil/json.go**

```go
// pkg/cmdutil/json.go
package cmdutil

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// JSONFlags holds --json flag state.
type JSONFlags struct {
	JSON   bool
	Fields []string
}

// AddJSONFlags adds --json flag to a command.
func AddJSONFlags(cmd *cobra.Command, jf *JSONFlags, validFields []string) {
	cmd.Flags().StringSliceVar(&jf.Fields, "json", nil,
		fmt.Sprintf("Output JSON with selected fields: %v", validFields))
}

// IsJSON returns true if --json was specified.
func (jf *JSONFlags) IsJSON() bool {
	return jf.Fields != nil
}

// PrintJSON writes v as indented JSON to w.
func PrintJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
```

- [ ] **Step 4: Create pkg/cmdutil/flags.go**

```go
// pkg/cmdutil/flags.go
package cmdutil

import "github.com/spf13/cobra"

// AddRepoOverride adds the --repo flag to a command.
func AddRepoOverride(cmd *cobra.Command, f *Factory) {
	cmd.PersistentFlags().StringVarP(&f.Host, "host", "", "", "Target Copia host")
}
```

- [ ] **Step 5: Fetch dependencies and verify build**

Run:
```bash
go get code.gitea.io/sdk/gitea
go get github.com/spf13/cobra
go get github.com/stretchr/testify
go build ./...
```

Expected: Build succeeds.

- [ ] **Step 6: Commit**

```bash
git add pkg/api/ pkg/cmdutil/ go.mod go.sum
git commit -m "feat: add factory, API client, JSON flags, and cmdutil helpers"
```

---

### Task 6: Root Command

**Files:**
- Create: `internal/copiacmd/root.go`
- Modify: `cmd/copia/main.go`

- [ ] **Step 1: Create internal/copiacmd/root.go**

```go
// internal/copiacmd/root.go
package copiacmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/internal/build"
	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

// NewRootCmd creates the root `copia` command with all subcommands.
func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "copia <command> <subcommand> [flags]",
		Short:   "Copia CLI — source control for industrial automation",
		Long:    "Work with Copia repositories, issues, pull requests, and more from the command line.",
		Version: build.Version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetVersionTemplate("copia version {{.Version}}\n")

	// Global flags
	cmd.PersistentFlags().StringVar(&f.Host, "host", "", "Target Copia host")
	cmd.PersistentFlags().StringVar(&f.Token, "token", "", "Authentication token")

	// Subcommands will be registered here as they are built:
	// cmd.AddCommand(authCmd.NewCmdAuth(f))
	// cmd.AddCommand(repoCmd.NewCmdRepo(f))
	// etc.

	return cmd
}

// Main is the entrypoint called from cmd/copia/main.go.
func Main() int {
	ios := iostreams.System()

	f := &cmdutil.Factory{
		IOStreams: ios,
		Config: func() (*config.Config, error) {
			return config.Load(config.DefaultPath())
		},
	}

	// Override token from env if not set by flag
	if envToken := os.Getenv("COPIA_TOKEN"); envToken != "" && f.Token == "" {
		f.Token = envToken
	}
	if envHost := os.Getenv("COPIA_HOST"); envHost != "" && f.Host == "" {
		f.Host = envHost
	}

	rootCmd := NewRootCmd(f)

	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}
```

- [ ] **Step 2: Update cmd/copia/main.go**

```go
// cmd/copia/main.go
package main

import (
	"os"

	"github.com/qubernetic-org/copia-cli/internal/copiacmd"
)

func main() {
	code := copiacmd.Main()
	os.Exit(code)
}
```

- [ ] **Step 3: Build and verify**

Run:
```bash
make build && ./bin/copia --version && ./bin/copia --help
```

Expected: Version output and help text with global flags.

- [ ] **Step 4: Commit**

```bash
git add internal/copiacmd/ cmd/copia/main.go
git commit -m "feat: add root command with global flags and env var support"
```

---

### Task 7: Auth Login

**Files:**
- Create: `pkg/cmd/auth/auth.go`
- Create: `pkg/cmd/auth/login/login.go`
- Create: `pkg/cmd/auth/login/login_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/auth/login/login_test.go
package login

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
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
	dir := t.TempDir()
	configPath := dir + "/config.yml"

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

	// Verify config was saved
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/auth/login/ -v`
Expected: FAIL — package does not exist.

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/auth/login/login.go
package login

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type LoginOptions struct {
	IO         *iostreams.IOStreams
	Host       string
	Token      string
	ConfigPath string
	HTTPClient *http.Client
}

func NewCmdLogin(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a Copia instance",
		Example: `  # Interactive login
  copia auth login

  # Non-interactive login (CI/agent)
  copia auth login --host app.copia.io --token YOUR_TOKEN`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.ConfigPath = config.DefaultPath()

			if f.Token != "" {
				opts.Token = f.Token
			}
			if f.Host != "" {
				opts.Host = f.Host
			}

			if opts.Host == "" {
				opts.Host = "app.copia.io"
			}

			if opts.Token == "" && opts.IO.IsStdinTTY() {
				fmt.Fprintf(opts.IO.ErrOut, "Enter token for %s: ", opts.Host)
				tokenBytes, err := io.ReadAll(io.LimitReader(opts.IO.In, 256))
				if err != nil {
					return err
				}
				opts.Token = string(tokenBytes)
			}

			if opts.Token == "" {
				return fmt.Errorf("token required. Use --token flag or run interactively")
			}

			opts.HTTPClient = &http.Client{}
			return loginRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Host, "host", "", "Copia instance hostname (default: app.copia.io)")
	cmd.Flags().StringVar(&opts.Token, "token", "", "Personal access token")

	return cmd
}

type userResponse struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
}

func loginRun(opts *LoginOptions) error {
	// Validate token by calling /api/v1/user
	url := fmt.Sprintf("https://%s/api/v1/user", opts.Host)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed for %s (HTTP %d)", opts.Host, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var user userResponse
	if err := json.Unmarshal(body, &user); err != nil {
		return fmt.Errorf("parsing user response: %w", err)
	}

	// Save to config
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}

	cfg.Hosts[opts.Host] = &config.HostConfig{
		Token: opts.Token,
		User:  user.Login,
	}

	if err := config.Save(opts.ConfigPath, cfg); err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, "Logged in as %s on %s\n", user.Login, opts.Host)
	return nil
}
```

- [ ] **Step 4: Create auth parent command**

```go
// pkg/cmd/auth/auth.go
package auth

import (
	"github.com/spf13/cobra"
	loginCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/auth/login"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate with Copia",
		Long:  "Manage authentication state for Copia instances.",
	}

	cmd.AddCommand(loginCmd.NewCmdLogin(f))
	// logout and status will be added in subsequent tasks

	return cmd
}
```

- [ ] **Step 5: Run tests**

Run: `go test ./pkg/cmd/auth/login/ -v`
Expected: PASS

- [ ] **Step 6: Register auth command in root**

In `internal/copiacmd/root.go`, add import and registration:

```go
import (
	// ... existing imports
	authCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/auth"
)

// Inside NewRootCmd, after global flags:
cmd.AddCommand(authCmd.NewCmdAuth(f))
```

- [ ] **Step 7: Build and smoke test**

Run:
```bash
make build && ./bin/copia auth login --help
```

Expected: Help text for `copia auth login` with `--host` and `--token` flags.

- [ ] **Step 8: Commit**

```bash
git add pkg/cmd/auth/ internal/copiacmd/root.go
git commit -m "feat: add copia auth login command with token validation"
```

---

### Task 8: Auth Logout

**Files:**
- Create: `pkg/cmd/auth/logout/logout.go`
- Create: `pkg/cmd/auth/logout/logout_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/auth/logout/logout_test.go
package logout

import (
	"testing"

	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutRun_RemovesHost(t *testing.T) {
	dir := t.TempDir()
	configPath := dir + "/config.yml"

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
	dir := t.TempDir()
	configPath := dir + "/config.yml"

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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/auth/logout/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/auth/logout/logout.go
package logout

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type LogoutOptions struct {
	IO         *iostreams.IOStreams
	Host       string
	ConfigPath string
}

func NewCmdLogout(f *cmdutil.Factory) *cobra.Command {
	opts := &LogoutOptions{}

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of a Copia instance",
		Example: "  copia auth logout --host app.copia.io",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.ConfigPath = config.DefaultPath()
			if f.Host != "" {
				opts.Host = f.Host
			}
			if opts.Host == "" {
				cfg, err := config.Load(opts.ConfigPath)
				if err != nil {
					return err
				}
				h, _ := cfg.DefaultHost()
				opts.Host = h
			}
			if opts.Host == "" {
				return fmt.Errorf("no host specified and no default configured")
			}
			return logoutRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Host, "host", "", "Copia instance hostname")

	return cmd
}

func logoutRun(opts *LogoutOptions) error {
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}

	if _, ok := cfg.Hosts[opts.Host]; !ok {
		return fmt.Errorf("not logged in to %s", opts.Host)
	}

	delete(cfg.Hosts, opts.Host)

	if err := config.Save(opts.ConfigPath, cfg); err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, "Logged out of %s\n", opts.Host)
	return nil
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./pkg/cmd/auth/logout/ -v`
Expected: PASS

- [ ] **Step 5: Register in auth parent**

In `pkg/cmd/auth/auth.go`, add:

```go
import (
	loginCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/auth/login"
	logoutCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/auth/logout"
)

// Inside NewCmdAuth:
cmd.AddCommand(logoutCmd.NewCmdLogout(f))
```

- [ ] **Step 6: Commit**

```bash
git add pkg/cmd/auth/logout/ pkg/cmd/auth/auth.go
git commit -m "feat: add copia auth logout command"
```

---

### Task 9: Auth Status

**Files:**
- Create: `pkg/cmd/auth/status/status.go`
- Create: `pkg/cmd/auth/status/status_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/auth/status/status_test.go
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

	dir := t.TempDir()
	configPath := dir + "/config.yml"
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
	dir := t.TempDir()
	configPath := dir + "/config.yml"
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/auth/status/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/auth/status/status.go
package status

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/internal/config"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type StatusOptions struct {
	IO         *iostreams.IOStreams
	ConfigPath string
	HTTPClient *http.Client
}

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	opts := &StatusOptions{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "View authentication status",
		Example: "  copia auth status",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.ConfigPath = config.DefaultPath()
			opts.HTTPClient = &http.Client{}
			return statusRun(opts)
		},
	}

	return cmd
}

func statusRun(opts *StatusOptions) error {
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}

	if len(cfg.Hosts) == 0 {
		return fmt.Errorf("not logged in to any Copia instance. Run 'copia auth login'")
	}

	for host, hc := range cfg.Hosts {
		tokenStatus := "Token valid"
		if opts.HTTPClient != nil {
			url := fmt.Sprintf("https://%s/api/v1/user", host)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				tokenStatus = fmt.Sprintf("Error: %v", err)
			} else {
				req.Header.Set("Authorization", "token "+hc.Token)
				resp, err := opts.HTTPClient.Do(req)
				if err != nil {
					tokenStatus = fmt.Sprintf("Error: %v", err)
				} else {
					resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						tokenStatus = fmt.Sprintf("Token invalid (HTTP %d)", resp.StatusCode)
					}
				}
			}
		}

		fmt.Fprintf(opts.IO.Out, "%s\n  User: %s\n  %s\n", host, hc.User, tokenStatus)
	}

	return nil
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./pkg/cmd/auth/status/ -v`
Expected: PASS

- [ ] **Step 5: Register in auth parent**

In `pkg/cmd/auth/auth.go`, add:

```go
import (
	// ... existing imports
	statusCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/auth/status"
)

// Inside NewCmdAuth:
cmd.AddCommand(statusCmd.NewCmdStatus(f))
```

- [ ] **Step 6: Run all tests**

Run: `make test`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add pkg/cmd/auth/
git commit -m "feat: add copia auth status command"
```

---

### Task 10: Repo List

**Files:**
- Create: `pkg/cmd/repo/repo.go`
- Create: `pkg/cmd/repo/list/list.go`
- Create: `pkg/cmd/repo/list/list_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/repo/list/list_test.go
package list

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListRun_UserRepos(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user/repos"),
		httpmock.StringResponse(http.StatusOK, `[
			{"full_name":"john/plc-project","description":"PLC code","html_url":"https://app.copia.io/john/plc-project","private":false,"updated_at":"2026-03-30T10:00:00Z"},
			{"full_name":"john/hmi-config","description":"HMI setup","html_url":"https://app.copia.io/john/hmi-config","private":true,"updated_at":"2026-03-29T10:00:00Z"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Limit:      30,
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "john/plc-project")
	assert.Contains(t, stdout.String(), "john/hmi-config")
}

func TestListRun_OrgRepos(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/orgs/my-org/repos"),
		httpmock.StringResponse(http.StatusOK, `[
			{"full_name":"my-org/main-plc","description":"Main PLC project","html_url":"https://app.copia.io/my-org/main-plc","private":false,"updated_at":"2026-03-30T10:00:00Z"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Org:        "my-org",
		Limit:      30,
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "my-org/main-plc")
}

func TestListRun_JSON(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user/repos"),
		httpmock.StringResponse(http.StatusOK, `[
			{"full_name":"john/plc-project","description":"PLC code","html_url":"https://app.copia.io/john/plc-project","private":false,"updated_at":"2026-03-30T10:00:00Z"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Limit:      30,
		JSON:       cmdutil.JSONFlags{Fields: []string{"fullName", "description"}},
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "john/plc-project")
	assert.Contains(t, stdout.String(), "PLC code")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/repo/list/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/repo/list/list.go
package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

var validJSONFields = []string{"fullName", "description", "private", "updatedAt", "htmlUrl"}

type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Org        string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type repoEntry struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Private     bool   `json:"private"`
	UpdatedAt   string `json:"updated_at"`
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List repositories",
		Aliases: []string{"ls"},
		Example: `  copia repo list
  copia repo list --org my-org
  copia repo list --json fullName,description`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.Host = f.Host
			opts.Token = f.Token

			if opts.Host == "" || opts.Token == "" {
				cfg, err := f.Config()
				if err != nil {
					return err
				}
				if opts.Host == "" {
					h, _ := cfg.DefaultHost()
					opts.Host = h
				}
				if opts.Token == "" && opts.Host != "" {
					if hc, ok := cfg.Hosts[opts.Host]; ok {
						opts.Token = hc.Token
					}
				}
			}

			opts.HTTPClient = &http.Client{}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Org, "org", "o", "", "List repositories for an organization")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of repositories")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func listRun(opts *ListOptions) error {
	endpoint := "/api/v1/user/repos"
	if opts.Org != "" {
		endpoint = fmt.Sprintf("/api/v1/orgs/%s/repos", opts.Org)
	}

	url := fmt.Sprintf("https://%s%s?limit=%d", opts.Host, endpoint, opts.Limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var repos []repoEntry
	if err := json.Unmarshal(body, &repos); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, repos)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, r := range repos {
		visibility := "public"
		if r.Private {
			visibility = "private"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.FullName, r.Description, visibility)
	}
	return w.Flush()
}
```

- [ ] **Step 4: Create repo parent command**

```go
// pkg/cmd/repo/repo.go
package repo

import (
	"github.com/spf13/cobra"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo/list"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

func NewCmdRepo(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo <command>",
		Short: "Manage repositories",
		Long:  "Work with Copia repositories.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	// view and clone will be added in subsequent tasks

	return cmd
}
```

- [ ] **Step 5: Register in root command**

In `internal/copiacmd/root.go`, add:

```go
import (
	// ... existing
	repoCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo"
)

// Inside NewRootCmd:
cmd.AddCommand(repoCmd.NewCmdRepo(f))
```

- [ ] **Step 6: Run tests**

Run: `go test ./pkg/cmd/repo/list/ -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add pkg/cmd/repo/ internal/copiacmd/root.go
git commit -m "feat: add copia repo list command with JSON output"
```

---

### Task 11: Repo View

**Files:**
- Create: `pkg/cmd/repo/view/view.go`
- Create: `pkg/cmd/repo/view/view_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/repo/view/view_test.go
package view

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo"),
		httpmock.StringResponse(http.StatusOK, `{
			"full_name":"my-org/my-repo",
			"description":"Main PLC project",
			"html_url":"https://app.copia.io/my-org/my-repo",
			"private":false,
			"default_branch":"main",
			"stars_count":5,
			"forks_count":2,
			"open_issues_count":3,
			"updated_at":"2026-03-30T10:00:00Z"
		}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ViewOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
	}

	err := viewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "my-org/my-repo")
	assert.Contains(t, stdout.String(), "Main PLC project")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/repo/view/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/repo/view/view.go
package view

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

var validJSONFields = []string{"fullName", "description", "private", "defaultBranch", "stars", "forks", "openIssues"}

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	JSON       cmdutil.JSONFlags
}

type repoDetail struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	HTMLURL         string `json:"html_url"`
	Private         bool   `json:"private"`
	DefaultBranch   string `json:"default_branch"`
	StarsCount      int    `json:"stars_count"`
	ForksCount      int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	UpdatedAt       string `json:"updated_at"`
}

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	opts := &ViewOptions{}

	cmd := &cobra.Command{
		Use:   "view [<owner/repo>]",
		Short: "View a repository",
		Example: `  copia repo view
  copia repo view my-org/my-repo
  copia repo view --json fullName,description`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.Host = f.Host
			opts.Token = f.Token

			if opts.Host == "" || opts.Token == "" {
				cfg, err := f.Config()
				if err != nil {
					return err
				}
				if opts.Host == "" {
					h, _ := cfg.DefaultHost()
					opts.Host = h
				}
				if opts.Token == "" && opts.Host != "" {
					if hc, ok := cfg.Hosts[opts.Host]; ok {
						opts.Token = hc.Token
					}
				}
			}

			if len(args) > 0 {
				parts := splitOwnerRepo(args[0])
				if parts == nil {
					return fmt.Errorf("expected owner/repo format")
				}
				opts.Owner, opts.Repo = parts[0], parts[1]
			} else if f.BaseRepo != nil {
				owner, repo, err := f.BaseRepo()
				if err != nil {
					return fmt.Errorf("could not determine repository. Use argument or --repo flag: %w", err)
				}
				opts.Owner, opts.Repo = owner, repo
			} else {
				return fmt.Errorf("could not determine repository. Provide owner/repo as argument")
			}

			opts.HTTPClient = &http.Client{}
			return viewRun(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func viewRun(opts *ViewOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s", opts.Host, opts.Owner, opts.Repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository %s/%s not found", opts.Owner, opts.Repo)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var repo repoDetail
	if err := json.Unmarshal(body, &repo); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, repo)
	}

	visibility := "public"
	if repo.Private {
		visibility = "private"
	}

	fmt.Fprintf(opts.IO.Out, "%s (%s)\n", repo.FullName, visibility)
	if repo.Description != "" {
		fmt.Fprintf(opts.IO.Out, "%s\n", repo.Description)
	}
	fmt.Fprintf(opts.IO.Out, "\nDefault branch: %s\n", repo.DefaultBranch)
	fmt.Fprintf(opts.IO.Out, "Stars: %d  Forks: %d  Open issues: %d\n", repo.StarsCount, repo.ForksCount, repo.OpenIssuesCount)
	fmt.Fprintf(opts.IO.Out, "URL: %s\n", repo.HTMLURL)

	return nil
}

func splitOwnerRepo(nwo string) []string {
	for i, c := range nwo {
		if c == '/' {
			if i > 0 && i < len(nwo)-1 {
				return []string{nwo[:i], nwo[i+1:]}
			}
		}
	}
	return nil
}
```

- [ ] **Step 4: Register in repo parent and run tests**

In `pkg/cmd/repo/repo.go`, add:

```go
import (
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo/list"
	viewCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo/view"
)

// Inside NewCmdRepo:
cmd.AddCommand(viewCmd.NewCmdView(f))
```

Run: `go test ./pkg/cmd/repo/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/repo/
git commit -m "feat: add copia repo view command"
```

---

### Task 12: Repo Clone

**Files:**
- Create: `pkg/cmd/repo/clone/clone.go`
- Create: `pkg/cmd/repo/clone/clone_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/repo/clone/clone_test.go
package clone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCloneURL(t *testing.T) {
	tests := []struct {
		name  string
		host  string
		nwo   string
		want  string
	}{
		{"full URL passthrough", "", "https://app.copia.io/my-org/my-repo.git", "https://app.copia.io/my-org/my-repo.git"},
		{"owner/repo format", "app.copia.io", "my-org/my-repo", "https://app.copia.io/my-org/my-repo.git"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCloneURL(tt.host, tt.nwo)
			assert.Equal(t, tt.want, got)
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/repo/clone/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/repo/clone/clone.go
package clone

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type CloneOptions struct {
	IO    *iostreams.IOStreams
	Host  string
	Token string
	Repo  string
	Dir   string
}

func NewCmdClone(f *cmdutil.Factory) *cobra.Command {
	opts := &CloneOptions{}

	cmd := &cobra.Command{
		Use:   "clone <owner/repo | URL> [<directory>]",
		Short: "Clone a repository",
		Example: `  copia repo clone my-org/my-repo
  copia repo clone https://app.copia.io/my-org/my-repo.git`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.Host = f.Host
			opts.Token = f.Token
			opts.Repo = args[0]

			if opts.Host == "" {
				cfg, err := f.Config()
				if err != nil {
					return err
				}
				h, _ := cfg.DefaultHost()
				opts.Host = h
			}

			if len(args) > 1 {
				opts.Dir = args[1]
			}

			return cloneRun(opts)
		},
	}

	return cmd
}

func cloneRun(opts *CloneOptions) error {
	cloneURL := buildCloneURL(opts.Host, opts.Repo)

	args := []string{"clone", cloneURL}
	if opts.Dir != "" {
		args = append(args, opts.Dir)
	}

	gitCmd := exec.Command("git", args...)
	gitCmd.Stdout = opts.IO.Out
	gitCmd.Stderr = opts.IO.ErrOut
	gitCmd.Stdin = os.Stdin

	return gitCmd.Run()
}

func buildCloneURL(host, repo string) string {
	if strings.HasPrefix(repo, "https://") || strings.HasPrefix(repo, "git@") {
		return repo
	}
	return fmt.Sprintf("https://%s/%s.git", host, repo)
}
```

- [ ] **Step 4: Register in repo parent and run tests**

In `pkg/cmd/repo/repo.go`, add:

```go
import (
	// ... existing
	cloneCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/repo/clone"
)

// Inside NewCmdRepo:
cmd.AddCommand(cloneCmd.NewCmdClone(f))
```

Run: `go test ./pkg/cmd/repo/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/repo/
git commit -m "feat: add copia repo clone command"
```

---

### Task 13: Label List

**Files:**
- Create: `pkg/cmd/label/label.go`
- Create: `pkg/cmd/label/list/list.go`
- Create: `pkg/cmd/label/list/list_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/label/list/list_test.go
package list

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/labels"),
		httpmock.StringResponse(http.StatusOK, `[
			{"id":1,"name":"bug","color":"#e11d48","description":"Something isn't working"},
			{"id":2,"name":"feature","color":"#0969da","description":"New feature request"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "bug")
	assert.Contains(t, stdout.String(), "feature")
}

func TestListRun_JSON(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/labels"),
		httpmock.StringResponse(http.StatusOK, `[
			{"id":1,"name":"bug","color":"#e11d48","description":"Something isn't working"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		JSON:       cmdutil.JSONFlags{Fields: []string{"name", "color"}},
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `"name"`)
	assert.Contains(t, stdout.String(), "bug")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/label/list/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/label/list/list.go
package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

var validJSONFields = []string{"name", "color", "description"}

type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	JSON       cmdutil.JSONFlags
}

type labelEntry struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List labels in a repository",
		Aliases: []string{"ls"},
		Example: `  copia label list
  copia label list --json name,color`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams

			client, host, token, owner, repo, err := resolveContext(f)
			if err != nil {
				return err
			}
			_ = client
			opts.Host = host
			opts.Token = token
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return listRun(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func resolveContext(f *cmdutil.Factory) (interface{}, string, string, string, string, error) {
	host := f.Host
	token := f.Token

	if host == "" || token == "" {
		cfg, err := f.Config()
		if err != nil {
			return nil, "", "", "", "", err
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
		return nil, "", "", "", "", fmt.Errorf("no host configured. Run 'copia auth login'")
	}

	var owner, repo string
	if f.BaseRepo != nil {
		var err error
		owner, repo, err = f.BaseRepo()
		if err != nil {
			return nil, "", "", "", "", fmt.Errorf("could not determine repository: %w", err)
		}
	}
	if owner == "" || repo == "" {
		return nil, "", "", "", "", fmt.Errorf("could not determine repository. Use --repo flag")
	}

	return nil, host, token, owner, repo, nil
}

func listRun(opts *ListOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/labels", opts.Host, opts.Owner, opts.Repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var labels []labelEntry
	if err := json.Unmarshal(body, &labels); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, labels)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, l := range labels {
		fmt.Fprintf(w, "%s\t%s\t%s\n", l.Name, l.Color, l.Description)
	}
	return w.Flush()
}
```

- [ ] **Step 4: Create label parent command**

```go
// pkg/cmd/label/label.go
package label

import (
	"github.com/spf13/cobra"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/label/list"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

func NewCmdLabel(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label <command>",
		Short: "Manage labels",
		Long:  "Work with repository labels.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))

	return cmd
}
```

- [ ] **Step 5: Register in root and run tests**

In `internal/copiacmd/root.go`, add:

```go
import (
	// ... existing
	labelCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/label"
)

// Inside NewRootCmd:
cmd.AddCommand(labelCmd.NewCmdLabel(f))
```

Run: `go test ./pkg/cmd/label/... -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add pkg/cmd/label/ internal/copiacmd/root.go
git commit -m "feat: add copia label list command"
```

---

### Task 14: Label Create

**Files:**
- Create: `pkg/cmd/label/create/create.go`
- Create: `pkg/cmd/label/create/create_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/label/create/create_test.go
package create

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/labels"),
		httpmock.StringResponse(http.StatusCreated, `{"id":1,"name":"critical","color":"#e11d48","description":"Critical issue"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:          ios,
		HTTPClient:  &http.Client{Transport: reg},
		Host:        "app.copia.io",
		Token:       "test-token",
		Owner:       "my-org",
		Repo:        "my-repo",
		Name:        "critical",
		Color:       "#e11d48",
		Description: "Critical issue",
	}

	err := createRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "critical")
}

func TestCreateRun_MissingName(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:    ios,
		Owner: "my-org",
		Repo:  "my-repo",
		Name:  "",
	}

	err := createRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name required")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/label/create/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/label/create/create.go
package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type CreateOptions struct {
	IO          *iostreams.IOStreams
	HTTPClient  *http.Client
	Host        string
	Token       string
	Owner       string
	Repo        string
	Name        string
	Color       string
	Description string
}

type createRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a label",
		Example: `  copia label create --name bug --color "#e11d48"
  copia label create --name feature --color "#0969da" --description "New feature"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			opts.Host = f.Host
			opts.Token = f.Token

			if opts.Host == "" || opts.Token == "" {
				cfg, err := f.Config()
				if err != nil {
					return err
				}
				if opts.Host == "" {
					h, _ := cfg.DefaultHost()
					opts.Host = h
				}
				if opts.Token == "" && opts.Host != "" {
					if hc, ok := cfg.Hosts[opts.Host]; ok {
						opts.Token = hc.Token
					}
				}
			}

			if f.BaseRepo != nil {
				owner, repo, err := f.BaseRepo()
				if err != nil {
					return fmt.Errorf("could not determine repository: %w", err)
				}
				opts.Owner, opts.Repo = owner, repo
			}

			opts.HTTPClient = &http.Client{}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Label name (required)")
	cmd.Flags().StringVarP(&opts.Color, "color", "c", "", "Label color in hex (e.g. #e11d48)")
	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Label description")

	return cmd
}

func createRun(opts *CreateOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("name required")
	}

	payload := createRequest{
		Name:        opts.Name,
		Color:       opts.Color,
		Description: opts.Description,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/labels", opts.Host, opts.Owner, opts.Repo)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create label (HTTP %d)", resp.StatusCode)
	}

	fmt.Fprintf(opts.IO.Out, "Label %q created\n", opts.Name)
	return nil
}
```

- [ ] **Step 4: Register in label parent and run tests**

In `pkg/cmd/label/label.go`, add:

```go
import (
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/label/list"
	createCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/label/create"
)

// Inside NewCmdLabel:
cmd.AddCommand(createCmd.NewCmdCreate(f))
```

Run: `go test ./pkg/cmd/label/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/label/
git commit -m "feat: add copia label create command"
```

---

### Task 15: Issue List

**Files:**
- Create: `pkg/cmd/issue/issue.go`
- Create: `pkg/cmd/issue/list/list.go`
- Create: `pkg/cmd/issue/list/list_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/issue/list/list_test.go
package list

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/issues"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":12,"title":"Fix PLC connection timeout","state":"open","updated_at":"2026-03-30T10:00:00Z","labels":[{"name":"bug"}]},
			{"number":11,"title":"Add safety interlock","state":"open","updated_at":"2026-03-29T10:00:00Z","labels":[]}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		State:      "open",
		Limit:      30,
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Fix PLC connection timeout")
	assert.Contains(t, stdout.String(), "Add safety interlock")
	assert.Contains(t, stdout.String(), "12")
}

func TestListRun_JSON(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/issues"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":12,"title":"Fix PLC connection timeout","state":"open","updated_at":"2026-03-30T10:00:00Z","labels":[]}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		State:      "open",
		Limit:      30,
		JSON:       cmdutil.JSONFlags{Fields: []string{"number", "title"}},
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `"number"`)
}

func TestListRun_FilterByState(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/issues"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":10,"title":"Closed issue","state":"closed","updated_at":"2026-03-28T10:00:00Z","labels":[]}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		State:      "closed",
		Limit:      30,
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed issue")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/issue/list/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/issue/list/list.go
package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

var validJSONFields = []string{"number", "title", "state", "labels", "updatedAt"}

type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	State      string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type labelRef struct {
	Name string `json:"name"`
}

type issueEntry struct {
	Number    int64      `json:"number"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	UpdatedAt string     `json:"updated_at"`
	Labels    []labelRef `json:"labels"`
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List issues in a repository",
		Aliases: []string{"ls"},
		Example: `  copia issue list
  copia issue list --state closed
  copia issue list --json number,title,state`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			host, token, owner, repo, err := resolveRepoContext(f)
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state: {open|closed|all}")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of issues")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func resolveRepoContext(f *cmdutil.Factory) (host, token, owner, repo string, err error) {
	host = f.Host
	token = f.Token

	if host == "" || token == "" {
		cfg, cfgErr := f.Config()
		if cfgErr != nil {
			return "", "", "", "", cfgErr
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
		return "", "", "", "", fmt.Errorf("no host configured. Run 'copia auth login'")
	}

	if f.BaseRepo != nil {
		owner, repo, err = f.BaseRepo()
		if err != nil {
			return "", "", "", "", fmt.Errorf("could not determine repository: %w", err)
		}
	}
	if owner == "" || repo == "" {
		return "", "", "", "", fmt.Errorf("could not determine repository. Use --repo flag")
	}

	return host, token, owner, repo, nil
}

func listRun(opts *ListOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues?state=%s&limit=%d&type=issues",
		opts.Host, opts.Owner, opts.Repo, opts.State, opts.Limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var issues []issueEntry
	if err := json.Unmarshal(body, &issues); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, issues)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, i := range issues {
		labels := ""
		for j, l := range i.Labels {
			if j > 0 {
				labels += ", "
			}
			labels += l.Name
		}
		fmt.Fprintf(w, "#%d\t%s\t%s\t%s\n", i.Number, i.Title, i.State, labels)
	}
	return w.Flush()
}
```

- [ ] **Step 4: Create issue parent command**

```go
// pkg/cmd/issue/issue.go
package issue

import (
	"github.com/spf13/cobra"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/list"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage issues",
		Long:  "Work with Copia repository issues.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))

	return cmd
}
```

- [ ] **Step 5: Register in root and run tests**

In `internal/copiacmd/root.go`, add:

```go
import (
	// ... existing
	issueCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue"
)

// Inside NewRootCmd:
cmd.AddCommand(issueCmd.NewCmdIssue(f))
```

Run: `go test ./pkg/cmd/issue/... -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add pkg/cmd/issue/ internal/copiacmd/root.go
git commit -m "feat: add copia issue list command with state filtering and JSON output"
```

---

### Task 16: Issue Create

**Files:**
- Create: `pkg/cmd/issue/create/create.go`
- Create: `pkg/cmd/issue/create/create_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/issue/create/create_test.go
package create

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues"),
		httpmock.StringResponse(http.StatusCreated, `{"number":13,"title":"Fix sensor mapping","html_url":"https://app.copia.io/my-org/my-repo/issues/13"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Title:      "Fix sensor mapping",
		Body:       "The sensor I/O mapping is incorrect.",
		Labels:     []string{"bug"},
	}

	err := createRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "#13")
	assert.Contains(t, stdout.String(), "Fix sensor mapping")
}

func TestCreateRun_MissingTitle(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:    ios,
		Title: "",
	}

	err := createRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title required")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/issue/create/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/issue/create/create.go
package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Title      string
	Body       string
	Labels     []string
	Assignees  []string
}

type createRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

type createResponse struct {
	Number  int64  `json:"number"`
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an issue",
		Example: `  copia issue create --title "Fix sensor mapping" --label bug
  copia issue create --title "Add feature" --body "Description here"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			host, token, owner, repo, err := resolveRepoContext(f)
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Issue title (required)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Issue body")
	cmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", nil, "Add labels")
	cmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", nil, "Assign users")

	return cmd
}

func resolveRepoContext(f *cmdutil.Factory) (host, token, owner, repo string, err error) {
	host = f.Host
	token = f.Token

	if host == "" || token == "" {
		cfg, cfgErr := f.Config()
		if cfgErr != nil {
			return "", "", "", "", cfgErr
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

	if f.BaseRepo != nil {
		owner, repo, err = f.BaseRepo()
		if err != nil {
			return "", "", "", "", err
		}
	}
	return host, token, owner, repo, nil
}

func createRun(opts *CreateOptions) error {
	if opts.Title == "" {
		return fmt.Errorf("title required")
	}

	payload := createRequest{
		Title:  opts.Title,
		Body:   opts.Body,
		Labels: opts.Labels,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues", opts.Host, opts.Owner, opts.Repo)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create issue (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result createResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, "Created issue #%d: %s\n%s\n", result.Number, result.Title, result.HTMLURL)
	return nil
}
```

- [ ] **Step 4: Register in issue parent and run tests**

In `pkg/cmd/issue/issue.go`, add:

```go
import (
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/list"
	createCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/issue/create"
)

// Inside NewCmdIssue:
cmd.AddCommand(createCmd.NewCmdCreate(f))
```

Run: `go test ./pkg/cmd/issue/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/issue/
git commit -m "feat: add copia issue create command"
```

---

### Task 17: Issue View

**Files:**
- Create: `pkg/cmd/issue/view/view.go`
- Create: `pkg/cmd/issue/view/view_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/issue/view/view_test.go
package view

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{
			"number":12,
			"title":"Fix PLC connection timeout",
			"body":"The PLC connection times out after 30 seconds.",
			"state":"open",
			"html_url":"https://app.copia.io/my-org/my-repo/issues/12",
			"user":{"login":"john"},
			"labels":[{"name":"bug"}],
			"created_at":"2026-03-30T10:00:00Z",
			"comments":2
		}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ViewOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
	}

	err := viewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Fix PLC connection timeout")
	assert.Contains(t, stdout.String(), "john")
	assert.Contains(t, stdout.String(), "open")
}

func TestViewRun_JSON(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{
			"number":12,"title":"Fix PLC","body":"","state":"open",
			"html_url":"https://app.copia.io/my-org/my-repo/issues/12",
			"user":{"login":"john"},"labels":[],"created_at":"2026-03-30T10:00:00Z","comments":0
		}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ViewOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		JSON:       cmdutil.JSONFlags{Fields: []string{"number", "title"}},
	}

	err := viewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `"number"`)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/issue/view/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/issue/view/view.go
package view

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

var validJSONFields = []string{"number", "title", "body", "state", "author", "labels", "createdAt", "comments"}

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Number     int64
	JSON       cmdutil.JSONFlags
}

type userRef struct {
	Login string `json:"login"`
}

type labelRef struct {
	Name string `json:"name"`
}

type issueDetail struct {
	Number    int64      `json:"number"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	State     string     `json:"state"`
	HTMLURL   string     `json:"html_url"`
	User      userRef    `json:"user"`
	Labels    []labelRef `json:"labels"`
	CreatedAt string     `json:"created_at"`
	Comments  int        `json:"comments"`
}

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	opts := &ViewOptions{}

	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View an issue",
		Example: `  copia issue view 12
  copia issue view 12 --json number,title,state`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid issue number: %s", args[0])
			}
			opts.Number = num
			opts.IO = f.IOStreams

			host, token, owner, repo, err := resolveRepoContext(f)
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return viewRun(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func resolveRepoContext(f *cmdutil.Factory) (host, token, owner, repo string, err error) {
	host = f.Host
	token = f.Token

	if host == "" || token == "" {
		cfg, cfgErr := f.Config()
		if cfgErr != nil {
			return "", "", "", "", cfgErr
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

	if f.BaseRepo != nil {
		owner, repo, err = f.BaseRepo()
		if err != nil {
			return "", "", "", "", err
		}
	}
	return host, token, owner, repo, nil
}

func viewRun(opts *ViewOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues/%d",
		opts.Host, opts.Owner, opts.Repo, opts.Number)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("issue #%d not found", opts.Number)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var issue issueDetail
	if err := json.Unmarshal(body, &issue); err != nil {
		return err
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, issue)
	}

	fmt.Fprintf(opts.IO.Out, "#%d %s\n", issue.Number, issue.Title)
	fmt.Fprintf(opts.IO.Out, "State: %s  Author: %s  Comments: %d\n", issue.State, issue.User.Login, issue.Comments)

	labels := ""
	for i, l := range issue.Labels {
		if i > 0 {
			labels += ", "
		}
		labels += l.Name
	}
	if labels != "" {
		fmt.Fprintf(opts.IO.Out, "Labels: %s\n", labels)
	}

	if issue.Body != "" {
		fmt.Fprintf(opts.IO.Out, "\n%s\n", issue.Body)
	}

	fmt.Fprintf(opts.IO.Out, "\n%s\n", issue.HTMLURL)
	return nil
}
```

- [ ] **Step 4: Register and run tests**

In `pkg/cmd/issue/issue.go`, add `viewCmd` import and registration.

Run: `go test ./pkg/cmd/issue/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/issue/
git commit -m "feat: add copia issue view command"
```

---

### Task 18: Issue Close

**Files:**
- Create: `pkg/cmd/issue/close/close.go`
- Create: `pkg/cmd/issue/close/close_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/issue/close/close_test.go
package close

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloseRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{"number":12,"state":"closed"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CloseOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
	}

	err := closeRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed issue #12")
}

func TestCloseRun_WithComment(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues/12/comments"),
		httpmock.StringResponse(http.StatusCreated, `{"id":1}`),
	)
	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{"number":12,"state":"closed"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CloseOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		Comment:    "Fixed in PR #7",
	}

	err := closeRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed issue #12")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/issue/close/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/issue/close/close.go
package close

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type CloseOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Number     int64
	Comment    string
}

func NewCmdClose(f *cmdutil.Factory) *cobra.Command {
	opts := &CloseOptions{}

	cmd := &cobra.Command{
		Use:   "close <number>",
		Short: "Close an issue",
		Example: `  copia issue close 12
  copia issue close 12 --comment "Fixed in PR #7"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid issue number: %s", args[0])
			}
			opts.Number = num
			opts.IO = f.IOStreams

			host, token, owner, repo, err := resolveRepoContext(f)
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return closeRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Comment, "comment", "c", "", "Add a comment before closing")

	return cmd
}

func resolveRepoContext(f *cmdutil.Factory) (host, token, owner, repo string, err error) {
	host = f.Host
	token = f.Token

	if host == "" || token == "" {
		cfg, cfgErr := f.Config()
		if cfgErr != nil {
			return "", "", "", "", cfgErr
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

	if f.BaseRepo != nil {
		owner, repo, err = f.BaseRepo()
		if err != nil {
			return "", "", "", "", err
		}
	}
	return host, token, owner, repo, nil
}

func closeRun(opts *CloseOptions) error {
	if opts.Comment != "" {
		commentPayload, _ := json.Marshal(map[string]string{"body": opts.Comment})
		commentURL := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues/%d/comments",
			opts.Host, opts.Owner, opts.Repo, opts.Number)

		req, err := http.NewRequest("POST", commentURL, bytes.NewReader(commentPayload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "token "+opts.Token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := opts.HTTPClient.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to add comment (HTTP %d)", resp.StatusCode)
		}
	}

	closePayload, _ := json.Marshal(map[string]string{"state": "closed"})
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues/%d",
		opts.Host, opts.Owner, opts.Repo, opts.Number)

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(closePayload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to close issue (HTTP %d)", resp.StatusCode)
	}

	fmt.Fprintf(opts.IO.Out, "Closed issue #%d\n", opts.Number)
	return nil
}
```

- [ ] **Step 4: Register and run tests**

Run: `go test ./pkg/cmd/issue/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/issue/
git commit -m "feat: add copia issue close command with optional comment"
```

---

### Task 19: Issue Comment

**Files:**
- Create: `pkg/cmd/issue/comment/comment.go`
- Create: `pkg/cmd/issue/comment/comment_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/issue/comment/comment_test.go
package comment

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues/12/comments"),
		httpmock.StringResponse(http.StatusCreated, `{"id":42,"body":"Investigating this now.","html_url":"https://app.copia.io/my-org/my-repo/issues/12#issuecomment-42"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CommentOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		Body:       "Investigating this now.",
	}

	err := commentRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Comment added to issue #12")
}

func TestCommentRun_MissingBody(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &CommentOptions{
		IO:   ios,
		Body: "",
	}

	err := commentRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "body required")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/issue/comment/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/issue/comment/comment.go
package comment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type CommentOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Number     int64
	Body       string
}

func NewCmdComment(f *cmdutil.Factory) *cobra.Command {
	opts := &CommentOptions{}

	cmd := &cobra.Command{
		Use:   "comment <number>",
		Short: "Add a comment to an issue",
		Example: `  copia issue comment 12 --body "Investigating this now."`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid issue number: %s", args[0])
			}
			opts.Number = num
			opts.IO = f.IOStreams

			host, token, owner, repo, err := resolveRepoContext(f)
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return commentRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Comment body (required)")

	return cmd
}

func resolveRepoContext(f *cmdutil.Factory) (host, token, owner, repo string, err error) {
	host = f.Host
	token = f.Token

	if host == "" || token == "" {
		cfg, cfgErr := f.Config()
		if cfgErr != nil {
			return "", "", "", "", cfgErr
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

	if f.BaseRepo != nil {
		owner, repo, err = f.BaseRepo()
		if err != nil {
			return "", "", "", "", err
		}
	}
	return host, token, owner, repo, nil
}

func commentRun(opts *CommentOptions) error {
	if opts.Body == "" {
		return fmt.Errorf("body required")
	}

	payload, _ := json.Marshal(map[string]string{"body": opts.Body})
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/issues/%d/comments",
		opts.Host, opts.Owner, opts.Repo, opts.Number)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add comment (HTTP %d)", resp.StatusCode)
	}

	fmt.Fprintf(opts.IO.Out, "Comment added to issue #%d\n", opts.Number)
	return nil
}
```

- [ ] **Step 4: Register all issue subcommands and run tests**

In `pkg/cmd/issue/issue.go`, register all subcommands (list, create, view, close, comment).

Run: `go test ./pkg/cmd/issue/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/issue/
git commit -m "feat: add copia issue comment command"
```

---

### Task 20: PR List

**Files:**
- Create: `pkg/cmd/pr/pr.go`
- Create: `pkg/cmd/pr/list/list.go`
- Create: `pkg/cmd/pr/list/list_test.go`

The PR commands follow the same pattern as issue commands but use the Gitea pull request API endpoints (`/repos/{owner}/{repo}/pulls`).

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/pr/list/list_test.go
package list

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/pulls"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":7,"title":"feat: add cylinder wrapper","state":"open","user":{"login":"john"},"base":{"label":"main"},"head":{"label":"feature/cylinder"},"updated_at":"2026-03-30T10:00:00Z"},
			{"number":6,"title":"fix: sensor timeout","state":"open","user":{"login":"jane"},"base":{"label":"main"},"head":{"label":"fix/sensor"},"updated_at":"2026-03-29T10:00:00Z"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		State:      "open",
		Limit:      30,
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "add cylinder wrapper")
	assert.Contains(t, stdout.String(), "sensor timeout")
}

func TestListRun_JSON(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/pulls"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":7,"title":"feat: add cylinder wrapper","state":"open","user":{"login":"john"},"base":{"label":"main"},"head":{"label":"feature/cylinder"},"updated_at":"2026-03-30T10:00:00Z"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		State:      "open",
		Limit:      30,
		JSON:       cmdutil.JSONFlags{Fields: []string{"number", "title"}},
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `"number"`)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/pr/list/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/pr/list/list.go
package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

var validJSONFields = []string{"number", "title", "state", "author", "base", "head", "updatedAt"}

type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	State      string
	Limit      int
	JSON       cmdutil.JSONFlags
}

type branchRef struct {
	Label string `json:"label"`
}

type userRef struct {
	Login string `json:"login"`
}

type prEntry struct {
	Number    int64     `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	User      userRef   `json:"user"`
	Base      branchRef `json:"base"`
	Head      branchRef `json:"head"`
	UpdatedAt string    `json:"updated_at"`
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List pull requests",
		Aliases: []string{"ls"},
		Example: `  copia pr list
  copia pr list --state closed
  copia pr list --json number,title,state`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			host, token, owner, repo, err := resolveRepoContext(f)
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state: {open|closed|all}")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of pull requests")
	cmdutil.AddJSONFlags(cmd, &opts.JSON, validJSONFields)

	return cmd
}

func resolveRepoContext(f *cmdutil.Factory) (host, token, owner, repo string, err error) {
	host = f.Host
	token = f.Token

	if host == "" || token == "" {
		cfg, cfgErr := f.Config()
		if cfgErr != nil {
			return "", "", "", "", cfgErr
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

	if f.BaseRepo != nil {
		owner, repo, err = f.BaseRepo()
		if err != nil {
			return "", "", "", "", err
		}
	}
	if owner == "" || repo == "" {
		return "", "", "", "", fmt.Errorf("could not determine repository. Use --repo flag")
	}
	return host, token, owner, repo, nil
}

func listRun(opts *ListOptions) error {
	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/pulls?state=%s&limit=%d",
		opts.Host, opts.Owner, opts.Repo, opts.State, opts.Limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (HTTP %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var prs []prEntry
	if err := json.Unmarshal(body, &prs); err != nil {
		return err
	}

	if opts.JSON.IsJSON() {
		return cmdutil.PrintJSON(opts.IO.Out, prs)
	}

	w := tabwriter.NewWriter(opts.IO.Out, 0, 0, 2, ' ', 0)
	for _, pr := range prs {
		fmt.Fprintf(w, "#%d\t%s\t%s\t%s\t%s <- %s\n",
			pr.Number, pr.Title, pr.State, pr.User.Login, pr.Base.Label, pr.Head.Label)
	}
	return w.Flush()
}
```

- [ ] **Step 4: Create PR parent command and register in root**

```go
// pkg/cmd/pr/pr.go
package pr

import (
	"github.com/spf13/cobra"
	listCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr/list"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
)

func NewCmdPR(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr <command>",
		Short: "Manage pull requests",
		Long:  "Work with Copia pull requests.",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))

	return cmd
}
```

In `internal/copiacmd/root.go`, add:

```go
import (
	// ... existing
	prCmd "github.com/qubernetic-org/copia-cli/pkg/cmd/pr"
)

// Inside NewRootCmd:
cmd.AddCommand(prCmd.NewCmdPR(f))
```

- [ ] **Step 5: Run tests**

Run: `go test ./pkg/cmd/pr/... -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add pkg/cmd/pr/ internal/copiacmd/root.go
git commit -m "feat: add copia pr list command"
```

---

### Task 21: PR Create

**Files:**
- Create: `pkg/cmd/pr/create/create.go`
- Create: `pkg/cmd/pr/create/create_test.go`

- [ ] **Step 1: Write the failing test**

```go
// pkg/cmd/pr/create/create_test.go
package create

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/pulls"),
		httpmock.StringResponse(http.StatusCreated, `{
			"number":8,
			"title":"feat: add cylinder wrapper",
			"html_url":"https://app.copia.io/my-org/my-repo/pulls/8"
		}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Title:      "feat: add cylinder wrapper",
		Body:       "Adds cylinder control wrapper.",
		Base:       "main",
		Head:       "feature/cylinder",
	}

	err := createRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "#8")
	assert.Contains(t, stdout.String(), "feat: add cylinder wrapper")
}

func TestCreateRun_MissingTitle(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	opts := &CreateOptions{IO: ios, Title: ""}

	err := createRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title required")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/pr/create/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// pkg/cmd/pr/create/create.go
package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/qubernetic-org/copia-cli/pkg/cmdutil"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient *http.Client
	Host       string
	Token      string
	Owner      string
	Repo       string
	Title      string
	Body       string
	Base       string
	Head       string
	Labels     []string
}

type createRequest struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Base  string `json:"base"`
	Head  string `json:"head"`
}

type createResponse struct {
	Number  int64  `json:"number"`
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull request",
		Example: `  copia pr create --title "feat: add wrapper" --base main --head feature/wrapper
  copia pr create --title "fix: timeout" --base develop --head fix/timeout --body "Fixes #12"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IO = f.IOStreams
			host, token, owner, repo, err := resolveRepoContext(f)
			if err != nil {
				return err
			}
			opts.Host = host
			opts.Token = token
			opts.Owner = owner
			opts.Repo = repo
			opts.HTTPClient = &http.Client{}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "PR title (required)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "PR body")
	cmd.Flags().StringVar(&opts.Base, "base", "main", "Base branch")
	cmd.Flags().StringVarP(&opts.Head, "head", "H", "", "Head branch (default: current branch)")

	return cmd
}

func resolveRepoContext(f *cmdutil.Factory) (host, token, owner, repo string, err error) {
	host = f.Host
	token = f.Token
	if host == "" || token == "" {
		cfg, cfgErr := f.Config()
		if cfgErr != nil {
			return "", "", "", "", cfgErr
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
	if f.BaseRepo != nil {
		owner, repo, err = f.BaseRepo()
		if err != nil {
			return "", "", "", "", err
		}
	}
	return host, token, owner, repo, nil
}

func createRun(opts *CreateOptions) error {
	if opts.Title == "" {
		return fmt.Errorf("title required")
	}

	payload := createRequest{
		Title: opts.Title,
		Body:  opts.Body,
		Base:  opts.Base,
		Head:  opts.Head,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/api/v1/repos/%s/%s/pulls", opts.Host, opts.Owner, opts.Repo)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", opts.Host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create PR (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result createResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, "Created PR #%d: %s\n%s\n", result.Number, result.Title, result.HTMLURL)
	return nil
}
```

- [ ] **Step 4: Register and run tests**

Run: `go test ./pkg/cmd/pr/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/pr/
git commit -m "feat: add copia pr create command"
```

---

### Task 22: PR View, Merge, Close

These three commands follow the exact same patterns as Issue View and Close. I'll outline them concisely — the implementation follows the same Options struct + test + run function pattern.

**Files:**
- Create: `pkg/cmd/pr/view/view.go` + `view_test.go`
- Create: `pkg/cmd/pr/merge/merge.go` + `merge_test.go`
- Create: `pkg/cmd/pr/close/close.go` + `close_test.go`

- [ ] **Step 1: PR View — write test and implementation**

`view.go` — `GET /api/v1/repos/{owner}/{repo}/pulls/{number}`, displays PR details (title, state, author, base/head branches, mergeable status). `--json` support.

Test: mock GET response, verify output contains title, state, branches.

- [ ] **Step 2: PR Merge — write test and implementation**

`merge.go` — `POST /api/v1/repos/{owner}/{repo}/pulls/{number}/merge` with body `{"Do":"merge"}`. Flags: `--merge` (default), `--rebase`, `--squash`, `--delete-branch`.

```go
type MergeOptions struct {
	IO           *iostreams.IOStreams
	HTTPClient   *http.Client
	Host, Token, Owner, Repo string
	Number       int64
	Method       string // merge, rebase, squash
	DeleteBranch bool
}
```

Test: mock POST merge response (HTTP 200), verify "Merged PR #N" output. Test `--delete-branch` adds DELETE `/repos/{owner}/{repo}/branches/{branch}` call.

- [ ] **Step 3: PR Close — write test and implementation**

`close.go` — `PATCH /api/v1/repos/{owner}/{repo}/pulls/{number}` with body `{"state":"closed"}`. Same pattern as issue close.

Test: mock PATCH response, verify "Closed PR #N" output.

- [ ] **Step 4: Register all in PR parent**

```go
// pkg/cmd/pr/pr.go
func NewCmdPR(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr <command>",
		Short: "Manage pull requests",
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	cmd.AddCommand(createCmd.NewCmdCreate(f))
	cmd.AddCommand(viewCmd.NewCmdView(f))
	cmd.AddCommand(mergeCmd.NewCmdMerge(f))
	cmd.AddCommand(closeCmd.NewCmdClose(f))

	return cmd
}
```

- [ ] **Step 5: Run all tests**

Run: `go test ./pkg/cmd/pr/... -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add pkg/cmd/pr/
git commit -m "feat: add copia pr view, merge, and close commands"
```

---

### Task 23: GoReleaser and CI

**Files:**
- Create: `.goreleaser.yml`
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Create .goreleaser.yml**

```yaml
version: 2
project_name: copia

builds:
  - id: copia
    main: ./cmd/copia
    binary: copia
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X github.com/qubernetic-org/copia-cli/internal/build.Version={{.Version}}
      - -X github.com/qubernetic-org/copia-cli/internal/build.Date={{time "2006-01-02"}}
    env:
      - CGO_ENABLED=0

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip
    files:
      - LICENSE
      - README.md

release:
  github:
    owner: qubernetic-org
    name: copia-cli
  draft: true
  prerelease: auto

checksum:
  name_template: "checksums.txt"
```

- [ ] **Step 2: Create .github/workflows/ci.yml**

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
      - run: make test

  integration:
    runs-on: ubuntu-latest
    if: github.event_name == 'push'
    needs: test
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
      - run: make integration
        env:
          COPIA_TEST_TOKEN: ${{ secrets.COPIA_TEST_TOKEN }}
          COPIA_TEST_HOST: ${{ secrets.COPIA_TEST_HOST }}

  build:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
      - run: make build
```

- [ ] **Step 3: Verify goreleaser config**

Run:
```bash
go install github.com/goreleaser/goreleaser/v2@latest
goreleaser check
```

Expected: config is valid.

- [ ] **Step 4: Commit**

```bash
git add .goreleaser.yml .github/
git commit -m "feat: add GoReleaser config and GitHub Actions CI workflow"
```

---

### Task 24: Final Validation

- [ ] **Step 1: Run full test suite**

Run: `make test`
Expected: All unit tests PASS.

- [ ] **Step 2: Build binary**

Run: `make build`
Expected: `bin/copia` binary created.

- [ ] **Step 3: Verify all commands registered**

Run:
```bash
./bin/copia --help
./bin/copia auth --help
./bin/copia repo --help
./bin/copia issue --help
./bin/copia pr --help
./bin/copia label --help
```

Expected: All commands and subcommands visible in help output.

- [ ] **Step 4: Verify version**

Run: `./bin/copia --version`
Expected: `copia version DEV`

- [ ] **Step 5: Commit any final fixes and tag**

```bash
git add -A
git commit -m "chore: final MVP validation and cleanup"
```
