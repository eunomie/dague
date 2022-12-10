package internal

import "runtime/debug"

var (
	// Version is the version of the CLI injected in compilation time
	Version = "dev"
	// Commit if the git commit at compilation time
	Commit = ""
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				Commit = setting.Value
			}
		}
	}
}
