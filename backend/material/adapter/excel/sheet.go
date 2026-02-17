package excel

import (
	"fmt"
	"prg_tools/material/domain"

	"github.com/xuri/excelize/v2"
)

// Excel DTOs

type ProgressSheetSection struct {
	FromDate string
	ToDate   string
	Rows     []ProgressRow
}

type ProgressRow struct {
	Date     string
	Quantity float64
	BidItem  BidItem
}

type BidItem struct {
	Number        string
	Name          string
	UnitOfMeasure string
}

func createEstimateProgressSheet(sheet ProgressSheetSection, fileName string) error {
	f := excelize.NewFile()
	sheetName := "Progress"
	f.SetSheetName("Sheet1", sheetName)
	// Set header row
	headers := []string{"Date", "Quantity", "Bid Item Number", "Bid Item Name", "Unit of Measure"}
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	// Fill rows
	for i, row := range sheet.Rows {
		rowIndex := i + 2 // Excel rows start at 1, header is row 1

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), row.Date)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), row.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), row.BidItem.Number)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), row.BidItem.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), row.BidItem.UnitOfMeasure)
	}

	// Save file
	return f.SaveAs(fileName)
}

// Convert domain MaterialLogCollection -> Excel DTO
func convertToProgressSheet(mLogs domain.MaterialLogCollection) ProgressSheetSection {
	section := ProgressSheetSection{
		FromDate: mLogs.FromDate,
		ToDate:   mLogs.ToDate,
		Rows:     make([]ProgressRow, 0, len(mLogs.Logs)),
	}

	for _, log := range mLogs.Logs {
		row := ProgressRow{
			Date:     log.Date,
			Quantity: log.Quantity,
			BidItem: BidItem{
				Number:        log.Material.Number,
				Name:          log.Material.Name,
				UnitOfMeasure: log.Material.Unit,
			},
		}
		section.Rows = append(section.Rows, row)
	}

	return section
}
