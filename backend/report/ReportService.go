package report

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"prg_tools/helpers"
	"prg_tools/report/adapter/raken"
	"prg_tools/report/domain"
	"prg_tools/report/port"
)

type ReportExporterService struct {
	reportFetcher port.ReportFetcher
	baseDir       string
}

func NewReportExporter(baseDir string) (*ReportExporterService, error) {

	adapter, err := raken.NewAdapter()
	if err != nil {
		return nil, err
	}
	return &ReportExporterService{
		reportFetcher: adapter,
		baseDir:       baseDir,
	}, nil

}

func (r *ReportExporterService) GetBaseDir() string {
	return r.baseDir
}

func (r *ReportExporterService) SetBaseDir(dir string) {
	r.baseDir = filepath.Clean(dir)
}

func (r *ReportExporterService) ExportToBaseDir(fromDate, toDate string, onProgress func(message string)) error {
	reportCollections, err := r.reportFetcher.GetReports(fromDate, toDate)
	if err != nil {
		return fmt.Errorf("failed to fetch reports: %w", err)
	}

	if onProgress != nil {
		onProgress(fmt.Sprintf("Fetched %d report collections", len(reportCollections)))
	}

	for _, reportCollection := range reportCollections {
		projectDir := r.getProjectDirectory(reportCollection.Reports[0].Project)

		for _, report := range reportCollection.Reports {
			if onProgress != nil {
				progressMsg := fmt.Sprintf("Exporting report for project %s: %s", report.Project.Name, report.Date)
				onProgress(progressMsg)
			}
			savePath := filepath.Join(projectDir, "Daily Reports", report.ToFileName())
			if err := downloadPDF(report.PDFLink, savePath); err != nil {
				return fmt.Errorf("failed to export report: %w", err)
			}
		}
	}
	return nil
}
func projectFolder(project domain.Project) string {
	if project.Number == "" {
		return helpers.SanitizeFileName(project.Name)
	}
	sanitizedProjectName := helpers.SanitizeFileName(project.Name)
	return fmt.Sprintf("%s-%s", project.Number, sanitizedProjectName)
}
func (r *ReportExporterService) getProjectDirectory(project domain.Project) string {
	projectFolder := helpers.SanitizeFileName(fmt.Sprintf("%s-%s", project.Number, project.Name))
	projectDir := filepath.Join(r.baseDir, projectFolder)
	return projectDir
}

func downloadPDF(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download pdf: %s", resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
