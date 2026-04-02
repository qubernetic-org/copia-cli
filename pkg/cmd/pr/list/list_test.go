package list

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/cmdutil"
	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
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

	err := ListRun(opts)
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

	err := ListRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `"number"`)
}
