package issues

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchIssues_DefaultStateAll(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/plc/issues"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":5,"title":"Closed issue","state":"closed"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &SearchOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "plc",
		Query:      "issue",
		Limit:      30,
		// State intentionally empty — should default to "all"
	}

	err := SearchRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed issue")
}

func TestSearchIssues_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/plc/issues"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":12,"title":"Fix PLC timeout","state":"open"},
			{"number":5,"title":"Sensor error","state":"closed"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &SearchOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "plc",
		Query:      "timeout",
		State:      "open",
		Limit:      30,
	}

	err := SearchRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Fix PLC timeout")
}
