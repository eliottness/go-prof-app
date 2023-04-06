//go:build darwin && arm64
package shared

import _ "embed"

//go:embed libhog.dylib
var LibHog []byte