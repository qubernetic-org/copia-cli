package build

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Version is set at build time via ldflags.
var Version = "DEV"

// Date is set at build time via ldflags.
var Date = ""

func init() {
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}

// VersionInfo returns a formatted version string with build details.
func VersionInfo() string {
	version := fmt.Sprintf("copia-cli version %s", Version)
	if Date != "" {
		version += fmt.Sprintf(" (%s)", Date)
	}
	version += fmt.Sprintf("\ngo: %s\nos/arch: %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	return version
}
