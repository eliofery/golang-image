package resources

import "embed"

//go:embed all:views
var Views embed.FS

//go:embed assets
var Assets embed.FS
