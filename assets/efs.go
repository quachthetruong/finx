package assets

import (
	"embed"
)

//go:embed "migrations" config.yaml
var EmbeddedFiles embed.FS
