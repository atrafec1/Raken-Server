package report

import (
	"prg_tools/report/adapter/raken"
	"prg_tools/report/domain"
	"prg_tools/report/port"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func (r *ReportExporterService) SetBaseDir(dir string) {
	r.baseDir = filepath.Clean(dir)
}

func (r *ReportExporterService) ExportToBaseDir(fromDate, toDate string) error {
	reportCollections, err := r.reportFetcher.GetReports(fromDate, toDate)
	if err != nil {
		return fmt.Errorf("failed to fetch reports: %w", err)
	}

	for _, reportCollection := range reportCollections {
		projectDir := r.getProjectDirectory(reportCollection.Reports[0].Project)

		for _, report := range reportCollection.Reports {
			savePath := filepath.Join(projectDir, report.ToFileName())
			if err := downloadPDF(report.PDFLink, savePath); err != nil {
				return fmt.Errorf("failed to export report")
			}
		}
	}
	return nil
}

func (r *ReportExporterService) getProjectDirectory(project domain.Project) string {
	projectFolder := sanitize(fmt.Sprintf("%s-%s", project.Number, project.Name))
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

func sanitize(name string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	return replacer.Replace(name)
}

