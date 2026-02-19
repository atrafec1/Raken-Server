package excel

import (
	"fmt"
	"path/filepath"
	"prg_tools/helpers"
	"prg_tools/material/domain"
	"strings"
)

type Adapter struct {
	estimateProgDir string
}

func NewAdapter(estimateProgDir string) *Adapter {
	return &Adapter{
		//base directory with all the job information: currently "I:\\Raken"
		estimateProgDir: estimateProgDir,
	}
}

// We need to convert the JobMaterialCollection into not one progressSheet but multiple.
// Each project will be composed of a weekly progress sheet in a single excel file
func (a *Adapter) ExportJobMaterialInfo(allJobMaterialInfo []domain.JobMaterialInfo) error {
	fmt.Println("Exporting job material info to excel...")
	fmt.Println("Number of job material info collections to export: ", len(allJobMaterialInfo))
	for _, jobMatInfo := range allJobMaterialInfo {
		progSheets := CreateProgressSheets(jobMatInfo)
		if len(progSheets) == 0 {
			continue
		}
		jobFolder := a.jobEstimateFolder(jobMatInfo.Job)
		if err := CreateProgressWorkbook(jobFolder, "ProgressEstimate.xlsx", progSheets); err != nil {
			return fmt.Errorf("error exporting progress sheets to excel for job %s: %w", jobMatInfo.Job.Name, err)
		}
	}
	return nil
}

func (a *Adapter) jobEstimateFolder(job domain.Job) string {

	if strings.TrimSpace(job.Number) == "" {
		return filepath.Join(a.estimateProgDir, helpers.SanitizeFileName(job.Name), "Estimate Progress")
	}

	sanitizedJobName := helpers.SanitizeFileName(job.Name)
	return filepath.Join(a.estimateProgDir, fmt.Sprintf("%s-%s", job.Number, sanitizedJobName), "Estimate Progress")
}
