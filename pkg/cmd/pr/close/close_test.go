package close

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloseRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/pulls/7"),
		httpmock.StringResponse(http.StatusOK, `{"number":7,"state":"closed"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CloseOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     7,
	}

	err := CloseRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed PR #7")
}

func TestCloseRun_WithComment(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues/7/comments"),
		httpmock.StringResponse(http.StatusCreated, `{"id":1}`),
	)
	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/pulls/7"),
		httpmock.StringResponse(http.StatusOK, `{"number":7,"state":"closed"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CloseOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     7,
		Comment:    "Closing in favor of #8",
	}

	err := CloseRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed PR #7")
}

func TestCloseRun_WithDeleteBranch(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/pulls/7"),
		httpmock.StringResponse(http.StatusOK, `{"number":7,"state":"closed","head":{"label":"feature/old"}}`),
	)
	reg.Register(
		httpmock.REST("DELETE", "/api/v1/repos/my-org/my-repo/branches/feature/old"),
		httpmock.StringResponse(http.StatusNoContent, ``),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CloseOptions{
		IO:           ios,
		HTTPClient:   &http.Client{Transport: reg},
		Host:         "app.copia.io",
		Token:        "test-token",
		Owner:        "my-org",
		Repo:         "my-repo",
		Number:       7,
		DeleteBranch: true,
	}

	err := CloseRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed PR #7")
}
