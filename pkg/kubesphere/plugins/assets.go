package plugins

import "embed"

//go:embed files/*
var assets embed.FS

func Assets() embed.FS {
	return assets
}
