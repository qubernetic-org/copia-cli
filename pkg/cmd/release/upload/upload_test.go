package upload

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/releases/tags/v1.0.0"),
		httpmock.StringResponse(http.StatusOK, `{"id":42,"tag_name":"v1.0.0"}`),
	)
	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/releases/42/assets"),
		httpmock.StringResponse(http.StatusCreated, `{"id":1,"name":"binary.tar.gz","size":1024}`),
	)

	// Create temp file
	dir := t.TempDir()
	filePath := filepath.Join(dir, "binary.tar.gz")
	require.NoError(t, os.WriteFile(filePath, []byte("fake binary content"), 0644))

	ios, _, stdout, _ := iostreams.Test()

	opts := &UploadOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Tag:        "v1.0.0",
		Files:      []string{filePath},
	}

	err := UploadRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "binary.tar.gz")
}

func TestUploadRun_FileNotFound(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &UploadOptions{
		IO:    ios,
		Tag:   "v1.0.0",
		Files: []string{"/nonexistent/file.tar.gz"},
	}

	err := UploadRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file")
}
