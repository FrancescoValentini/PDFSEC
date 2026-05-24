package main

import (
	"embed"
	"log"

	"pdfsec/internal/gui"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := gui.NewApp()

	err := wails.Run(&options.App{
		Title:     "PDFSec",
		Width:     900,
		Height:    600,
		Assets:    assets,
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
