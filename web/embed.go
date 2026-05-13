package web

import "embed"

//go:embed index.html static
var StaticFS embed.FS
