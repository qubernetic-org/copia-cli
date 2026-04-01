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
		httpmock.REST("DELETE", "/api/v1/repos/my-org/my-repo"),
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
		Confirmed:  true,
	}

	err := deleteRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Deleted repository my-org/my-repo")
}

func TestDeleteRun_NotConfirmed(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &DeleteOptions{
		IO:        ios,
		Owner:     "my-org",
		Repo:      "my-repo",
		Confirmed: false,
	}

	err := deleteRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "use --yes to confirm")
}

func TestDeleteRun_NotFound(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("DELETE", "/api/v1/repos/my-org/missing"),
		httpmock.StringResponse(http.StatusNotFound, `{}`),
	)

	ios, _, _, _ := iostreams.Test()

	opts := &DeleteOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "missing",
		Confirmed:  true,
	}

	err := deleteRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete")
}
