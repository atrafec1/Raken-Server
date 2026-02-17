package main

import (
	"embed"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed *.env
var envFile string

func main() {
	println("ENV FILE CONTENT:")
	println(envFile)
	println("END ENV FILE CONTENT")
	// Parse and load embedded .env
	envMap, _ := godotenv.Parse(strings.NewReader(envFile))
	for key, value := range envMap {
		os.Setenv(key, value)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ReportingTool",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
