package merge

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/pulls/7/merge"),
		httpmock.StringResponse(http.StatusOK, ``),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &MergeOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     7,
		Method:     "merge",
	}

	err := mergeRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Merged PR #7")
}

func TestMergeRun_Squash(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/pulls/7/merge"),
		httpmock.StringResponse(http.StatusOK, ``),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &MergeOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     7,
		Method:     "squash",
	}

	err := mergeRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Merged PR #7")
}

func TestMergeRun_WithDeleteBranch(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/pulls/7/merge"),
		httpmock.StringResponse(http.StatusOK, ``),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &MergeOptions{
		IO:           ios,
		HTTPClient:   &http.Client{Transport: reg},
		Host:         "app.copia.io",
		Token:        "test-token",
		Owner:        "my-org",
		Repo:         "my-repo",
		Number:       7,
		Method:       "merge",
		DeleteBranch: true,
	}

	err := mergeRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Merged PR #7")
}
