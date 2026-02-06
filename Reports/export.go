package reports

import (
	"daily_check_in/api"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ExportReports(fromDate, toDate, outputDir string, client *api.Client) (string, error) {
	if client == nil {
		return "", fmt.Errorf("client cannot be nil")
	}
	if outputDir == "" {
		return "", fmt.Errorf("outputDir cannot be empty")
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("error creating output directory: %w", err)
	}

	reportService := ReportService{Client: client}
	projects, err := reportService.GetProjectsWorkedOn(fromDate, toDate)
	if err != nil {
		return "", fmt.Errorf("error getting projects worked on: %w", err)
	}

	reportCollections, err := reportService.GetReports(fromDate, toDate, projects)
	if err != nil {
		return "", fmt.Errorf("error getting reports: %w", err)
	}

	data, err := json.MarshalIndent(reportCollections, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error serializing reports: %w", err)
	}

	fileName := fmt.Sprintf("daily_reports_%s_to_%s.json", fromDate, toDate)
	filePath := filepath.Join(outputDir, fileName)
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return "", fmt.Errorf("error writing report file: %w", err)
	}

	return filePath, nil
}
