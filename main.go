package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "OpenClaw-Sifu",
		Width:            1440,
		Height:           920,
		MinWidth:         360,
		MinHeight:        220,
		DisableResize:    false,
		Frameless:        true,
		AssetServer:      &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 245, G: 240, B: 231, A: 1},
		Windows: &windows.Options{
			DisableFramelessWindowDecorations: false,
			Theme:                             windows.Light,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
