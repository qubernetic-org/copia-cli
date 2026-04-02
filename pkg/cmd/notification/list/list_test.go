package list

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/notifications"),
		httpmock.StringResponse(http.StatusOK, `[
			{"id":1,"subject":{"title":"Fix PLC timeout","type":"Issue"},"repository":{"full_name":"my-org/plc"},"unread":true},
			{"id":2,"subject":{"title":"Add sensor","type":"Pull"},"repository":{"full_name":"my-org/plc"},"unread":true}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
	}

	err := ListRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Fix PLC timeout")
	assert.Contains(t, stdout.String(), "Add sensor")
}

func TestListRun_Empty(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/notifications"),
		httpmock.StringResponse(http.StatusOK, `[]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
	}

	err := ListRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "No unread notifications")
}

func TestListRun_AllFlag(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("GET", "/api/v1/notifications"),
		httpmock.StringResponse(http.StatusOK, `[
			{"id":1,"subject":{"title":"Old PR","type":"Pull"},"repository":{"full_name":"my-org/plc"},"unread":false}
		]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &ListOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		All:        true,
	}

	err := ListRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Old PR")
}
