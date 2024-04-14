package version

import "fmt"

// Version is the semantic version of the build.
var Version = "unset"

// BuildDate is the date the executable was built.
var BuildDate = "unset"

// GitCommit is set to git rev-parse --short HEAD.
var GitCommit = "unset"

// UserAgent is a versioned user agent string to use for remote requests.
var UserAgent = fmt.Sprintf("metar-ws2811/%s/go", Version)

// BuildInfo contains the version information of the last compilation.
type BuildInfo struct {
	Version   string
	BuildDate string
	GitCommit string
}

// GetBuildInfo returns the all of the build and version information set during build.
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		BuildDate: BuildDate,
		GitCommit: GitCommit,
	}
}
