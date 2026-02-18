package excel

import (
	"fmt"
	"path/filepath"
	"prg_tools/material/domain"
)

type Adapter struct {
	estimateProgDir string
}

func NewAdapter(estimateProgDir string) *Adapter {
	return &Adapter{
		estimateProgDir: estimateProgDir,
	}
}

func (a *Adapter) ExportMaterialLogs(collections []domain.MaterialLogCollection) error {
	counter := 0
	for _, collection := range collections {
		fmt.Println("Converting to progress sheet...")
		progressSheet, err := ConvertToProgressSheet(collection, 21)
		if err != nil {
			return fmt.Errorf("failed to convert material logs to progress sheet for job %s: %w", collection.Job.Name, err)
		}
		fmt.Println("CreatedProgress sheet")

		fileName := filepath.Join(a.estimateProgDir, fmt.Sprintf("%d_%s.xlsx", counter, collection.FromDate))
		fmt.Println("Creating progress sheet for job ", collection.Job.Name, " with filename ", fileName)

		if err := progressSheet.CreateEstimateProgressSheet(fileName); err != nil {
			return fmt.Errorf("failed to create progress sheet for job %s: %w", collection.Job.Name, err)
		}
		fmt.Println("Created progress sheet for job ", collection.Job.Name)
		counter++
	}
	return nil
}
