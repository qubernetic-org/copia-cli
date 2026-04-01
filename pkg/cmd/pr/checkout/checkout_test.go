package checkout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRefSpec(t *testing.T) {
	ref := buildRefSpec(7)
	assert.Equal(t, "pull/7/head:pr-7", ref)
}

func TestBuildRefSpec_Large(t *testing.T) {
	ref := buildRefSpec(1234)
	assert.Equal(t, "pull/1234/head:pr-1234", ref)
}
