package delete

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/releases/tags/v1.0.0"),
		httpmock.StringResponse(http.StatusOK, `{"id":42,"tag_name":"v1.0.0"}`),
	)
	reg.Register(
		httpmock.REST("DELETE", "/api/v1/repos/my-org/my-repo/releases/42"),
		httpmock.StringResponse(http.StatusNoContent, ``),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &DeleteOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Tag:        "v1.0.0",
	}

	err := DeleteRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Deleted release v1.0.0")
}

func TestDeleteRun_NotFound(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/releases/tags/v9.9.9"),
		httpmock.StringResponse(http.StatusNotFound, `{}`),
	)

	ios, _, _, _ := iostreams.Test()

	opts := &DeleteOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Tag:        "v9.9.9",
	}

	err := DeleteRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
