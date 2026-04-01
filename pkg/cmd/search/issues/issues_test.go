package issues

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchIssues_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/search"),
		httpmock.StringResponse(http.StatusOK, `[
			{"number":12,"title":"Fix PLC timeout","state":"open","repository":{"full_name":"my-org/plc"}},
			{"number":5,"title":"Sensor error","state":"closed","repository":{"full_name":"my-org/plc"}}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &SearchOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Query:      "timeout",
		State:      "open",
		Limit:      30,
	}

	err := searchRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Fix PLC timeout")
}
