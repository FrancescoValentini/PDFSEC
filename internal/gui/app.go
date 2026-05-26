package gui

import (
	"context"
	"fmt"
	"pdfsec/internal/core"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	inputFile  string
	outputFile string
	ctx        context.Context
}

type EncryptPayload struct {
	InputPath  string   `json:"inputPath"`
	OutputPath string   `json:"outputPath"`
	ReaderPwd  string   `json:"readerPwd"`
	OwnerPwd   string   `json:"ownerPwd"`
	OwnerOnly  bool     `json:"ownerOnly"`
	Perms      []string `json:"perms"`
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
	a.inputFile = path
	a.outputFile = core.ResolveOutputFile(a.inputFile, a.outputFile)
	a.UI_UpdateOutputPath()
	return path, nil
}

func (a *App) UI_UpdateOutputPath() {
	payload := map[string]interface{}{
		"outputPath": a.outputFile,
	}
	runtime.EventsEmit(a.ctx, "PDFSEC:input-selected", payload)
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
	a.outputFile = path
	return path, nil
}

func (a *App) RandomPassword() (string, error) {
	return core.GeneratePassword()
}

func (a *App) UI_EncryptionFinished(success bool, err error) {

	payload := map[string]interface{}{
		"success": success,
	}

	if err != nil {
		payload["error"] = err.Error()
	}

	runtime.EventsEmit(a.ctx, "PDFSEC:encryption-finished", payload)
}

func (a *App) EncryptPDF(payload EncryptPayload) error {
	var err error
	var perms model.PermissionFlags
	permString := strings.Join(payload.Perms, ",")

	if permString == "" {
		perms = model.PermissionsNone
	} else {
		perms, err = core.ParsePermissions(permString)
	}

	if err != nil {
		a.UI_EncryptionFinished(false, err)
		return err
	}

	err = core.EncryptPDF(core.EncryptOptions{
		InputFile:     payload.InputPath,
		OutputFile:    payload.OutputPath,
		UserPassword:  payload.ReaderPwd,
		OwnerPassword: payload.OwnerPwd,
		Permissions:   perms,
	})

	if err != nil {
		a.UI_EncryptionFinished(false, err)
		return err
	}

	a.UI_EncryptionFinished(true, nil)

	return nil
}
