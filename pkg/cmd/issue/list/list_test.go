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

func TestListRun_InvalidState(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &ListOptions{
		IO:    ios,
		State: "invalid",
		Limit: 30,
	}

	err := ListRun(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid state")
}

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

	err := ListRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Fix PLC connection timeout")
	assert.Contains(t, stdout.String(), "12")
}

func TestListRun_JSON(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/issues"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":12,"title":"Fix PLC","state":"open","updated_at":"2026-03-30T10:00:00Z","labels":[]}
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

func TestListRun_LabelFilter(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/issues"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":12,"title":"Fix PLC timeout","state":"open","updated_at":"2026-03-30T10:00:00Z","labels":[{"name":"bug"}]}
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
		Labels:     []string{"bug"},
	}

	err := ListRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Fix PLC timeout")
}
