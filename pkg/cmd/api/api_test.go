package api

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiRun_GET(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/user"),
		httpmock.StringResponse(http.StatusOK, `{"login":"john","id":1}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &APIOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Method:     "GET",
		Path:       "/user",
	}

	err := apiRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "john")
}

func TestApiRun_POST_WithFields(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues"),
		httpmock.StringResponse(http.StatusCreated, `{"number":1,"title":"test issue"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &APIOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Method:     "POST",
		Path:       "/repos/my-org/my-repo/issues",
		Fields:     []string{"title=test issue", "body=description"},
	}

	err := apiRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "test issue")
}

func TestApiRun_DELETE(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("DELETE", "/api/v1/repos/my-org/my-repo"),
		httpmock.StringResponse(http.StatusNoContent, ``),
	)

	ios, _, _, _ := iostreams.Test()

	opts := &APIOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Method:     "DELETE",
		Path:       "/repos/my-org/my-repo",
	}

	err := apiRun(opts)
	require.NoError(t, err)
}

func TestApiRun_DefaultMethodGET(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/version"),
		httpmock.StringResponse(http.StatusOK, `{"version":"1.21.0"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &APIOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Path:       "/version",
	}

	err := apiRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "1.21.0")
}

func TestApiRun_MissingPath(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &APIOptions{IO: ios, Path: ""}

	err := apiRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path required")
}
