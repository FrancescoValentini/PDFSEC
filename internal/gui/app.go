package gui

import (
	"context"
	"fmt"
	"pdfsec/internal/core"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time! %s", name, core.AppName)
}

func (a *App) OpenFileDialog() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select PDF",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF Files",
				Pattern:     "*.pdf",
			},
		},
	})

	if err != nil {
		return "", err
	}

	return path, nil
}

func (a *App) SaveFileDialog() (string, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save PDF",
		DefaultFilename: "document.pdf",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF Files",
				Pattern:     "*.pdf",
			},
		},
	})

	if err != nil {
		return "", err
	}

	return path, nil
}
