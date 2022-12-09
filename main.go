package main

import (
	"embed"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed resources/iconLightMode.png
var lightModeIcon []byte

//go:embed resources/iconDarkMode.png
var darkModeIcon []byte

func main() {
	setup()
}
