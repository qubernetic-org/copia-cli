package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionInfo(t *testing.T) {
	// Save and restore
	origVersion := Version
	origDate := Date
	defer func() { Version = origVersion; Date = origDate }()

	Version = "1.0.0"
	Date = "2026-04-03"

	info := VersionInfo()
	assert.Contains(t, info, "1.0.0")
	assert.Contains(t, info, "2026-04-03")
	assert.Contains(t, info, "go")
}

func TestVersionInfo_DEV(t *testing.T) {
	origVersion := Version
	origDate := Date
	defer func() { Version = origVersion; Date = origDate }()

	Version = "DEV"
	Date = ""

	info := VersionInfo()
	assert.Contains(t, info, "DEV")
}
