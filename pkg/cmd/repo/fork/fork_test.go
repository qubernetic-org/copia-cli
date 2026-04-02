package fork

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForkRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/upstream-org/project/forks"),
		httpmock.StringResponse(http.StatusAccepted, `{"full_name":"john/project","html_url":"https://app.copia.io/john/project","clone_url":"https://app.copia.io/john/project.git"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ForkOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "upstream-org",
		Repo:       "project",
	}

	err := ForkRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "john/project")
}

func TestForkRun_ToOrg(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/upstream-org/project/forks"),
		httpmock.StringResponse(http.StatusAccepted, `{"full_name":"my-org/project","html_url":"https://app.copia.io/my-org/project"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ForkOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "upstream-org",
		Repo:       "project",
		Org:        "my-org",
	}

	err := ForkRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "my-org/project")
}
