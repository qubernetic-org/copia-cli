package clone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCloneURL_FullURL(t *testing.T) {
	got := buildCloneURL("app.copia.io", "https://app.copia.io/my-org/my-repo.git")
	assert.Equal(t, "https://app.copia.io/my-org/my-repo.git", got)
}

func TestBuildCloneURL_OwnerRepo(t *testing.T) {
	got := buildCloneURL("app.copia.io", "my-org/my-repo")
	assert.Equal(t, "https://app.copia.io/my-org/my-repo.git", got)
}

func TestBuildCloneURL_SSHPassthrough(t *testing.T) {
	got := buildCloneURL("", "git@app.copia.io:my-org/my-repo.git")
	assert.Equal(t, "git@app.copia.io:my-org/my-repo.git", got)
}
