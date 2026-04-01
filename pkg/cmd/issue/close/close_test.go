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
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{"number":12,"state":"closed"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CloseOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
	}

	err := closeRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed issue #12")
}

func TestCloseRun_WithComment(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues/12/comments"),
		httpmock.StringResponse(http.StatusCreated, `{"id":1}`),
	)
	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{"number":12,"state":"closed"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CloseOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		Comment:    "Fixed in PR #7",
	}

	err := closeRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Closed issue #12")
}
