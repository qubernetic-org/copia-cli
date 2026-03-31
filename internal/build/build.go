package build

import "runtime/debug"

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
