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
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/labels"),
		httpmock.StringResponse(http.StatusCreated, `{"id":1,"name":"critical","color":"#e11d48"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:          ios,
		HTTPClient:  &http.Client{Transport: reg},
		Host:        "app.copia.io",
		Token:       "test-token",
		Owner:       "my-org",
		Repo:        "my-repo",
		Name:        "critical",
		Color:       "#e11d48",
		Description: "Critical issue",
	}

	err := createRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "critical")
}

func TestCreateRun_MissingName(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:   ios,
		Name: "",
	}

	err := createRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name required")
}
