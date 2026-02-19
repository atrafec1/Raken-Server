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
//We need to convert the JobMaterialCollection into not one progressSheet but multiple.
//Each project will be composed of a weekly progress sheet in a single excel file
func (a *Adapter) ExportJobMaterialInfo(collections []domain.JobMaterialInfo) error {
	counter := 0
	fmt.Println("Exporting job material info to excel...")
	fmt.Println("Number of job material info collections to export: ", len(collections))
	for _, collection := range collections {
		fmt.Println("Converting to progress sheet...")
		progressSheet, err := ConvertToProgressSheet(collection, 21)
		if err != nil {
			return fmt.Errorf("failed to convert material logs to progress sheet for job %s: %w", collection.Job.Name, err)
		}
		fmt.Println("CreatedProgress sheet")
		jobEstimateFolder := a.jobEstimateFolder(collection.Job)
		fileName := "ProgressEstimate.xlsx"
		fmt.Println("Creating progress sheet for job ", collection.Job.Name, " with filename ", fileName)

		if err := progressSheet.CreateEstimateProgressSheet(jobEstimateFolder, fileName); err != nil {
			return fmt.Errorf("failed to create progress sheet for job %s: %w", collection.Job.Name, err)
		}
		fmt.Println("Created progress sheet for job ", collection.Job.Name)
		counter++
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
