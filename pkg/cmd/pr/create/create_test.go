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
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/pulls"),
		httpmock.StringResponse(http.StatusCreated, `{
			"number":8,"title":"feat: add cylinder wrapper",
			"html_url":"https://app.copia.io/my-org/my-repo/pulls/8"
		}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Title:      "feat: add cylinder wrapper",
		Body:       "Adds cylinder control wrapper.",
		Base:       "main",
		Head:       "feature/cylinder",
	}

	err := createRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "#8")
}

func TestCreateRun_MissingTitle(t *testing.T) {
	ios, _, _, _ := iostreams.Test()
	opts := &CreateOptions{IO: ios, Title: ""}

	err := createRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title required")
}
