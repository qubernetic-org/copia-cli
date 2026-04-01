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

func TestListRun_UserRepos(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user/repos"),
		httpmock.StringResponse(http.StatusOK, `[
			{"full_name":"john/plc-project","description":"PLC code","private":false,"updated_at":"2026-03-30T10:00:00Z"},
			{"full_name":"john/hmi-config","description":"HMI setup","private":true,"updated_at":"2026-03-29T10:00:00Z"}
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
			{"full_name":"my-org/main-plc","description":"Main PLC","private":false,"updated_at":"2026-03-30T10:00:00Z"}
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
			{"full_name":"john/plc-project","description":"PLC code","private":false,"updated_at":"2026-03-30T10:00:00Z"}
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
