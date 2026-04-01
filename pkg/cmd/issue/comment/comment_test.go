package comment

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues/12/comments"),
		httpmock.StringResponse(http.StatusCreated, `{"id":42}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CommentOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		Body:       "Investigating this now.",
	}

	err := commentRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Comment added to issue #12")
}

func TestCommentRun_MissingBody(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &CommentOptions{IO: ios, Body: ""}

	err := commentRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "body required")
}
