package view

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/orgs/my-org"),
		httpmock.StringResponse(http.StatusOK, `{"username":"my-org","full_name":"My Organization","description":"Industrial automation","website":"https://example.com"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ViewOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Name:       "my-org",
	}

	err := viewRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "my-org")
	assert.Contains(t, stdout.String(), "Industrial automation")
}
