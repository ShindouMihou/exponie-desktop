package main

import (
	"embed"
	"exponie/pkg/exponie"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := exponie.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title: "Exponie",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		Frameless: true,
		MinHeight: 667,
		MinWidth:  375,
		BackgroundColour: &options.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 1,
		},
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
