package diff

import (
	"net/http"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/httpmock"
	"github.com/qubernetic/copia-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffRun_Success(t *testing.T) {
	reg := &httpmock.Registry{}
	defer reg.Verify(t)

	diffContent := `diff --git a/main.go b/main.go
index abc1234..def5678 100644
--- a/main.go
+++ b/main.go
@@ -1,3 +1,4 @@
 package main
+
+import "fmt"
`

	reg.Register(
		httpmock.REST("GET", "/api/v1/repos/my-org/my-repo/pulls/7.diff"),
		httpmock.StringResponse(http.StatusOK, diffContent),
	)

	ios, _, stdout, _ := iostreams.Test()

	opts := &DiffOptions{
		IO:         ios,
		HTTPClient: &http.Client{Transport: reg},
		Host:       "app.copia.io",
		Token:      "test-token",
		Owner:      "my-org",
		Repo:       "my-repo",
		Number:     7,
	}

	err := DiffRun(opts)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "diff --git")
	assert.Contains(t, stdout.String(), "+import \"fmt\"")
}
