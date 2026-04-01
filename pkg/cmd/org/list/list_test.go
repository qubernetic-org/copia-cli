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
		httpmock.REST("GET", "/api/v1/user/orgs"),
		httpmock.StringResponse(http.StatusOK, `[
			{"username":"my-org","full_name":"My Organization","description":"Industrial automation"},
			{"username":"other-org","full_name":"Other Org","description":"Testing"}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "my-org")
	assert.Contains(t, stdout.String(), "other-org")
}

func TestListRun_JSON(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user/orgs"),
		httpmock.StringResponse(http.StatusOK, `[{"username":"my-org","full_name":"My Org","description":""}]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		JSON:       cmdutil.JSONFlags{Fields: []string{"username"}},
	}

	err := listRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `"username"`)
}
