//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	releasecreate "github.com/qubernetic/copia-cli/pkg/cmd/release/create"
	releasedelete "github.com/qubernetic/copia-cli/pkg/cmd/release/delete"
	releaselist "github.com/qubernetic/copia-cli/pkg/cmd/release/list"
)

func TestRelease_Lifecycle(t *testing.T) {
	env := loadTestEnv(t)
	tag := "v0.0.0-integration-test"

	// Create release
	createIO, createOut, _ := testIO()
	err := releasecreate.CreateRun(&releasecreate.CreateOptions{
		IO:         createIO,
		HTTPClient: &http.Client{},
		Host:       env.Host,
		Token:      env.Token,
		Owner:      env.Owner,
		Repo:       env.Repo,
		Tag:        tag,
		Title:      "Integration Test Release",
		Notes:      "Created by integration test. Will be deleted.",
	})
	require.NoError(t, err)
	assert.Contains(t, createOut.String(), tag)
	t.Logf("Create: %s", createOut.String())

	// Cleanup: delete at end
	defer func() {
		deleteIO, _, _ := testIO()
		_ = releasedelete.DeleteRun(&releasedelete.DeleteOptions{
			IO:         deleteIO,
			HTTPClient: &http.Client{},
			Host:       env.Host,
			Token:      env.Token,
			Owner:      env.Owner,
			Repo:       env.Repo,
			Tag:        tag,
		})
		t.Logf("Deleted release %s", tag)
	}()

	// List releases
	listIO, listOut, _ := testIO()
	err = releaselist.ListRun(&releaselist.ListOptions{
		IO:         listIO,
		HTTPClient: &http.Client{},
		Host:       env.Host,
		Token:      env.Token,
		Owner:      env.Owner,
		Repo:       env.Repo,
		Limit:      10,
	})
	require.NoError(t, err)
	assert.Contains(t, listOut.String(), tag)
	t.Logf("List: %s", listOut.String())

	// Delete release
	deleteIO, deleteOut, _ := testIO()
	err = releasedelete.DeleteRun(&releasedelete.DeleteOptions{
		IO:         deleteIO,
		HTTPClient: &http.Client{},
		Host:       env.Host,
		Token:      env.Token,
		Owner:      env.Owner,
		Repo:       env.Repo,
		Tag:        tag,
	})
	require.NoError(t, err)
	assert.Contains(t, deleteOut.String(), "Deleted")
	t.Logf("Delete: %s", deleteOut.String())
}
