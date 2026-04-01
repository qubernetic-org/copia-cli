package read

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadRun_MarkAll(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("PUT", "/api/v1/notifications"),
		httpmock.StringResponse(http.StatusResetContent, `[]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ReadOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		All:        true,
	}

	err := readRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "All notifications marked as read")
}

func TestReadRun_MarkSingle(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("PATCH", "/api/v1/notifications/threads/42"),
		httpmock.StringResponse(http.StatusResetContent, ``),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ReadOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		ThreadID:   42,
	}

	err := readRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Notification #42 marked as read")
}
