package main

import (
	"context"
	"fmt"
	"prg_tools/report"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	ReportExporter *report.ReportExporterService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ensureReportExporter() error {
	if a.ReportExporter != nil {
		return nil
	}
	svc, err := report.NewReportExporter("C:\\Users\\jdtra\\OneDrive\\Desktop\\raken")
	if err != nil {
		return fmt.Errorf("error initializing ReportExporterService: %v", err)
	}
	a.ReportExporter = svc
	return nil
}

func (a *App) ExportReports(fromDate, toDate string) error {
	if err := a.ensureReportExporter(); err != nil {
		return err
	}

	if err := a.ReportExporter.ExportToBaseDir(fromDate, toDate); err != nil {
		return fmt.Errorf("error exporting reports: %v", err)
	}
	return nil
}

func (a *App) ChangeExportDir(newDir string) (string, error) {
	if err := a.ensureReportExporter(); err != nil {
		return "", err
	}
	a.ReportExporter.SetBaseDir(newDir)
	return newDir, nil
}

func (a *App) SelectFolder() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Export Folder",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (a *App) GetExportDir() string {
	if err := a.ensureReportExporter(); err != nil {
		return ""
	}
	return a.ReportExporter.GetBaseDir()
}
