package create

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRun_UserRepo(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/user/repos"),
		httpmock.StringResponse(http.StatusCreated, `{"full_name":"john/new-repo","html_url":"https://app.copia.io/john/new-repo","clone_url":"https://app.copia.io/john/new-repo.git"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Name:       "new-repo",
		Private:    true,
	}

	err := CreateRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "john/new-repo")
}

func TestCreateRun_OrgRepo(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/orgs/my-org/repos"),
		httpmock.StringResponse(http.StatusCreated, `{"full_name":"my-org/new-repo","html_url":"https://app.copia.io/my-org/new-repo","clone_url":"https://app.copia.io/my-org/new-repo.git"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Name:       "new-repo",
		Org:        "my-org",
	}

	err := CreateRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "my-org/new-repo")
}

func TestCreateRun_MissingName(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	opts := &CreateOptions{IO: ios, Name: ""}

	err := CreateRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name required")
}
