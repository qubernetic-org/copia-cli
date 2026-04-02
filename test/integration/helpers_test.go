//go:build integration

package integration

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/qubernetic/copia-cli/pkg/iostreams"
)

type testEnv struct {
	Host  string
	Token string
	Owner string
	Repo  string
}

func loadTestEnv(t *testing.T) testEnv {
	t.Helper()

	host := os.Getenv("COPIA_TEST_HOST")
	token := os.Getenv("COPIA_TEST_TOKEN")
	owner := os.Getenv("COPIA_TEST_OWNER")
	repo := os.Getenv("COPIA_TEST_REPO")

	if host == "" || token == "" || owner == "" || repo == "" {
		t.Skip("Integration test env vars not set (COPIA_TEST_HOST, COPIA_TEST_TOKEN, COPIA_TEST_OWNER, COPIA_TEST_REPO)")
	}

	return testEnv{
		Host:  host,
		Token: token,
		Owner: owner,
		Repo:  repo,
	}
}

// testIO returns IOStreams and stdout buffer for asserting CLI output.
func testIO() (*iostreams.IOStreams, *bytes.Buffer, *bytes.Buffer) {
	ios, _, stdout, stderr := iostreams.Test()
	return ios, stdout, stderr
}

// testHTTPClient returns a real HTTP client for integration tests.
func testHTTPClient() *http.Client {
	return &http.Client{}
}
