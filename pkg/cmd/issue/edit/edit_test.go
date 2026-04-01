package edit

import (
	"net/http"
	"testing"

	"github.com/qubernetic-org/copia-cli/pkg/httpmock"
	"github.com/qubernetic-org/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditRun_SetTitle(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{"number":12,"title":"Updated title"}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &EditOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		Title:      "Updated title",
	}

	err := editRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Updated issue #12")
}

func TestEditRun_AddLabels(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("POST", "/api/v1/repos/my-org/my-repo/issues/12/labels"),
		httpmock.StringResponse(http.StatusOK, `[{"id":1,"name":"bug"}]`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &EditOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		AddLabels:  []string{"bug"},
	}

	err := editRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Updated issue #12")
}

func TestEditRun_SetAssignees(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{"number":12}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &EditOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		Assignees:  []string{"john", "jane"},
	}

	err := editRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Updated issue #12")
}

func TestEditRun_SetMilestone(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	reg.Register(
		httpmock.REST("PATCH", "/api/v1/repos/my-org/my-repo/issues/12"),
		httpmock.StringResponse(http.StatusOK, `{"number":12}`),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &EditOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     12,
		Milestone:  1,
	}

	err := editRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Updated issue #12")
}

func TestEditRun_NothingToEdit(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	opts := &EditOptions{
		IO:     ios,
		Number: 12,
	}

	err := editRun(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nothing to edit")
}
