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

func (a *Adapter) ExportMaterialLogs(logs []domain.MaterialLogCollection) error {
	counter := 0
	for _, log := range logs {
		fmt.Printf("%+v\n", log)
		fmt.Println("Converting to progress sheet...")
		progressSheet := convertToProgressSheet(log)
		printProgressSheet(progressSheet)
		fileName := filepath.Join(a.estimateProgDir, fmt.Sprintf("%d_%s.xlsx", counter, log.FromDate))
		if err := createEstimateProgressSheet(progressSheet, fileName); err != nil {
			return fmt.Errorf("failed to create progress sheet for job %s: %w", log.Job.Name, err)
		}
		counter++
	}
	return nil
}

func printProgressSheet(sheet ProgressSheetSection) {
	fmt.Printf("Progress Sheet: From %s To %s\n", sheet.FromDate, sheet.ToDate)
	for _, row := range sheet.Rows {
		fmt.Printf("  Date: %s, Quantity: %.2f, Bid Item: %s - %s (%s)\n",
			row.Date, row.Quantity, row.BidItem.Number, row.BidItem.Name, row.BidItem.UnitOfMeasure)
	}
}
