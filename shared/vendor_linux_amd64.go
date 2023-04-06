//go:build linux && amd64
package shared

import _ "embed"

//go:embed libhog.so
var LibHog []byte