package api

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientWithHTTP_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	reg.Register(
		httpmock.REST("GET", "/api/v1/version"),
		httpmock.StringResponse(http.StatusOK, `{"version":"1.21.0"}`),
	)

	client, err := NewClientWithHTTP("app.copia.io", "test-token", &http.Client{Transport: reg})
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithHTTP_EmptyHost(t *testing.T) {
	_, err := NewClientWithHTTP("", "test-token", &http.Client{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "host is required")
}
