package create

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/releases"),
		httpmock.StringResponse(http.StatusCreated, `{"id":1,"tag_name":"v1.0.0","name":"Release 1.0.0","html_url":"https://app.copia.io/my-org/my-repo/releases/tag/v1.0.0"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Tag:        "v1.0.0",
		Title:      "Release 1.0.0",
		Notes:      "First stable release.",
	}

	err := CreateRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "v1.0.0")
}

func TestCreateRun_MissingTag(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	opts := &CreateOptions{IO: ios, Tag: ""}

	err := CreateRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tag required")
}

func TestCreateRun_Draft(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/releases"),
		httpmock.StringResponse(http.StatusCreated, `{"id":2,"tag_name":"v2.0.0-rc.1","name":"RC1","draft":true,"html_url":"https://app.copia.io/my-org/my-repo/releases/tag/v2.0.0-rc.1"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Tag:        "v2.0.0-rc.1",
		Title:      "RC1",
		Draft:      true,
	}

	err := CreateRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "v2.0.0-rc.1")
}
