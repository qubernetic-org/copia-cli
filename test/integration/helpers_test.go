//go:build integration

package integration

import (
	"os"
	"testing"
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
