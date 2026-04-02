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
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues"),
		httpmock.StringResponse(http.StatusCreated, `{"number":13,"title":"Fix sensor mapping","html_url":"https://app.copia.io/my-org/my-repo/issues/13"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Title:      "Fix sensor mapping",
		Body:       "The sensor I/O mapping is incorrect.",
	}

	err := CreateRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "#13")
}

func TestCreateRun_MissingTitle(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &CreateOptions{IO: ios, Title: ""}

	err := CreateRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title required")
}
