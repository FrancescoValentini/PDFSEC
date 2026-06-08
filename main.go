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
		Width:     520,
		Height:    720,
		Assets:    assets,
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
		DisableResize:            true,
		EnableDefaultContextMenu: false,
	})
	if err != nil {
		log.Fatal(err)
	}
}
