package config

import "embed"

//go:embed "config.yaml"
var EmbeddedFiles embed.FS
