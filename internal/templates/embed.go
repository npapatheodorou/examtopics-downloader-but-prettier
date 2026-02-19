package templates

import _ "embed"

// EmbeddedTemplate is the built-in HTML shell used when template.html
// is not available on disk (for distributed standalone executables).
//
//go:embed template.html
var EmbeddedTemplate string
