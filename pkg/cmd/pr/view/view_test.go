package view

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/pulls/7"),
		httpmock.StringResponse(http.StatusOK, `{
			"number":7,"title":"feat: add cylinder wrapper","state":"open",
			"body":"Adds cylinder control.","mergeable":true,
			"html_url":"https://app.copia.io/my-org/my-repo/pulls/7",
			"user":{"login":"john"},"base":{"label":"main"},"head":{"label":"feature/cylinder"},
			"created_at":"2026-03-30T10:00:00Z"
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
		Number:     7,
	}

	err := ViewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "feat: add cylinder wrapper")
	assert.Contains(t, stdout.String(), "john")
	assert.Contains(t, stdout.String(), "main <- feature/cylinder")
}

func TestViewRun_JSON(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/pulls/7"),
		httpmock.StringResponse(http.StatusOK, `{
			"number":7,"title":"feat: test","state":"open","body":"","mergeable":true,
			"html_url":"","user":{"login":"john"},"base":{"label":"main"},"head":{"label":"feat"},
			"created_at":"2026-03-30T10:00:00Z"
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
		Number:     7,
		JSON:       cmdutil.JSONFlags{Fields: []string{"number", "title"}},
	}

	err := ViewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `"number"`)
}
