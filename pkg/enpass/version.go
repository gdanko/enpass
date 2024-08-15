package enpass

import (
	"fmt"
)

var (
	Major  = "0"
	Minor  = "1"
	Patch  = "4"
	Suffix = "dev"
)

// Version returns a version string based on the SemVer parts defined at compile time. Dev builds will result in
// 0.0.0-dev. Prefix (v) and suffix can be optionally included, while suffix will only be included if one is defined.
func Version(prefix, suffix bool) string {
	version := fmt.Sprintf("%s.%s.%s", Major, Minor, Patch)

	if prefix {
		version = "v" + version
	}

	if suffix && Suffix != "" {
		version = version + "-" + Suffix
	}

	return version
}
