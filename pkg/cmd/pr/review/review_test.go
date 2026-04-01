package review

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviewRun_Approve(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/pulls/7/reviews"),
		httpmock.StringResponse(http.StatusOK, `{"id":1}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ReviewOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     7,
		Event:      "APPROVED",
	}

	err := reviewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Approved PR #7")
}

func TestReviewRun_RequestChanges(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/pulls/7/reviews"),
		httpmock.StringResponse(http.StatusOK, `{"id":2}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ReviewOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     7,
		Event:      "REQUEST_CHANGES",
		Body:       "Please fix the tests.",
	}

	err := reviewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Requested changes on PR #7")
}

func TestReviewRun_Comment(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/pulls/7/reviews"),
		httpmock.StringResponse(http.StatusOK, `{"id":3}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ReviewOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     7,
		Event:      "COMMENT",
		Body:       "Looks good overall.",
	}

	err := reviewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Commented on PR #7")
}
