package configs

import "embed"

// FS contains the embedded default configuration files.
//
//go:embed *.textproto
var FS embed.FS
