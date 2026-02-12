package main

import (
	"context"
	"fmt"
	"prg_tools/report"
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
	svc, err := report.NewReportExporter("C:\\Users\\adamt\\OneDrive\\Desktop\\Raken")
	if err != nil {
		return fmt.Errorf("error initializing ReportExporterService: %v", err)
	}
	a.ReportExporter = svc
	return nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
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

func (a *App) ChangeExportDir(newDir string) error {
	if err := a.ensureReportExporter(); err != nil {
		return err
	}
	a.ReportExporter.SetBaseDir(newDir)
	return nil
}
