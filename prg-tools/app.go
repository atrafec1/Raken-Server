package main

import (
	"context"
	"fmt"
	"prg_tools/material"
	"prg_tools/material/domain"
	"prg_tools/payroll"
	"prg_tools/payroll/dto"
	"prg_tools/report"
	"prg_tools/datbase"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx                     context.Context
	ReportExporter          *report.ReportExporterService
	PayrollService          *payroll.PayrollService
	ProgressEstimateService *material.ProgressEstimateService
	DB                      *gorm.DB
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


func (a *App) InitDB() error {
	
}
func (a *App) ensurePayrollService() error {
	if a.PayrollService != nil {
		return nil
	}

	svc, err := payroll.NewCPService()
	if err != nil {
		return fmt.Errorf("error initializing PayrollService: %v", err)
	}
	a.PayrollService = svc
	return nil
}

func (a *App) FetchPayrollEntries(fromDate, toDate string) (dto.PayrollEntryResult, error) {
	if err := a.ensurePayrollService(); err != nil {
		return dto.PayrollEntryResult{}, err
	}
	result, err := a.PayrollService.GetEntries(fromDate, toDate)
	if err != nil {
		return dto.PayrollEntryResult{}, fmt.Errorf("error fetching payroll entries: %v", err)
	}
	return result, nil
}

// Handles both CP export and Excel export
func (a *App) ProcessPayroll(result dto.PayrollEntryResult) error {
	if err := a.ensurePayrollService(); err != nil {
		return fmt.Errorf("error initializing PayrollService: %v", err)
	}

	if err := a.PayrollService.ExportToPayroll(result.Entries); err != nil {
		return fmt.Errorf("error exporting to payroll: %v", err)
	}
	if err := a.PayrollService.ExportExcel(result); err != nil {
		return fmt.Errorf("error exporting to excel: %v", err)
	}
	return nil
}

func (a *App) ExportPayrollWarnings(warnings []dto.Warning) error {
	if err := a.ensurePayrollService(); err != nil {
		return fmt.Errorf("error initializing PayrollService: %v", err)
	}
	if err := a.PayrollService.ExportWarnings(warnings); err != nil {
		return fmt.Errorf("error exporting warnings: %v", err)
	}
	return nil
}

func (a *App) GroupPayrollWarnings(payrollResult dto.PayrollEntryResult) map[string][]dto.Warning {
	grouped := make(map[string][]dto.Warning)
	for _, w := range payrollResult.Warnings {
		grouped[w.WarningType] = append(grouped[w.WarningType], w)
	}
	return grouped
}

func (a *App) ensureReportExporter() error {
	if a.ReportExporter != nil {
		return nil
	}
	svc, err := report.NewReportExporter("I:\\Raken")
	if err != nil {
		return fmt.Errorf("error initializing ReportExporterService: %v", err)
	}
	a.ReportExporter = svc
	return nil
}

func (a *App) progressEventEmitter(message string) {
	runtime.EventsEmit(a.ctx, "exportProgress", message)
}
func (a *App) ExportDailyReports(fromDate, toDate string) error {
	runtime.EventsEmit(a.ctx, "exportProgress", "Starting export...")
	if err := a.ensureReportExporter(); err != nil {
		return err
	}
	go func() {
		err := a.ReportExporter.ExportToBaseDir(
			fromDate,
			toDate,
			a.progressEventEmitter,
		)

		if err != nil {
			runtime.EventsEmit(a.ctx, "exportError", fmt.Sprintf("Export failed: %v", err))
			return
		}
		runtime.EventsEmit(a.ctx, "exportComplete", fmt.Sprintf("Export completed. Check %s for the results.", a.ReportExporter.GetBaseDir()))
	}()
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

// Progress Estimate Service
func (a *App) ensureProgressEstimateService() error {
	if a.ProgressEstimateService != nil {
		return nil
	}
	svc, err := material.RakenProgressEstimateService("I:\\Raken")
	if err != nil {
		return fmt.Errorf("failed to intitialize Raken Progress Estimate Service: %w", err)
	}
	a.ProgressEstimateService = svc
	return nil
}

func (a *App) ExportProgressEstimate(JobMatInfo []domain.JobMaterialInfo) error {
	a.ensureProgressEstimateService()

	if err := a.ProgressEstimateService.ExportJobMaterialInfo(JobMatInfo); err != nil {
		return fmt.Errorf("error exporting job material info: %v", err)
	}
	return nil
}

func (a *App) GetJobMaterialInfo(fromDate, toDate string) ([]domain.JobMaterialInfo, error) {
	a.ensureProgressEstimateService()
	return a.ProgressEstimateService.GetJobMaterialInfo(fromDate, toDate)
}
