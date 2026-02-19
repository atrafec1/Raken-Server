package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

// --- Orchestration ---

func (p *ProgressSheet) CreateEstimateProgressSheet(saveDir, fileName string) error {
	if !pathExists(saveDir) {
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			return fmt.Errorf("could not create or update folder at path %s: %w", saveDir, err)
		}
	}

	f := excelize.NewFile()
	if f.GetSheetName(0) == "Sheet1" {
		renameBaseSheet(f, "Estimate - Daily Analysis")
	}
	//Sheet for human modified estimate sheet
	if p.SheetName == "" {
		p.SheetName = "Bid Item Data From Raken"
	}

	workingIndex, err := f.NewSheet(p.SheetName)
	if err != nil {
		return fmt.Errorf("failed to create new sheet for progress data: %w", err)
	}
	f.SetActiveSheet(workingIndex)

	currentRow := 3

	// Header
	if err := p.createHeader(f); err != nil {
		return err
	}
	currentRow += 3

	//Bid Item Header

	if err := p.createBidHeader(f, currentRow); err != nil {
		return err
	}
	currentRow += 3
	// Sections + rows
	for _, sec := range p.Sections {
		if err := p.createSection(f, sec, currentRow); err != nil {
			return err
		}
		currentRow++
		for _, r := range sec.Rows {
			if err := p.createRow(f, r, currentRow); err != nil {
				return err
			}
			currentRow++
		}
		currentRow++
	}
	entireFilePath := filepath.Join(saveDir, fileName)
	return f.SaveAs(entireFilePath)
}

// --- Excel Writing ---
func (p *ProgressSheet) createHeader(f *excelize.File) error {
	// Job Name in A1
	if err := f.SetCellValue(p.SheetName, "A1", "Pacific Restoration Group, Inc."); err != nil {
		return err
	}
	if err := f.SetCellValue(p.SheetName, "A2", p.JobDetail); err != nil {
		return err
	}
	return nil
}

func (p *ProgressSheet) createBidHeader(f *excelize.File, startRow int) error {
	sortBidItems(p.BidItems)

	// Write "Date"
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("B%d", startRow+2), "Date"); err != nil {
		return err
	}
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("B%d", startRow), "Bid Item"); err != nil {
		return err
	}
	bidNumberStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})

	for col, bi := range p.BidItems {
		col = col + 1
		letter := colLetter(col + 2)

		if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, startRow), bi.Number); err != nil {
			return err
		}
		if err := f.SetCellStyle(p.SheetName, fmt.Sprintf("%s%d", letter, startRow), fmt.Sprintf("%s%d", letter, startRow), bidNumberStyle); err != nil {
			return err
		}
		if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, startRow+1), bi.Name); err != nil {
			return err
		}
		if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, startRow+2), bi.UnitOfMeasure); err != nil {
			return err
		}
	}

	// --- Create contiguous bottom border ---
	style, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "bottom",
				Color: "000000",
				Style: 2, // medium line
			},
		},
	})
	if err != nil {
		return err
	}

	// Apply border across entire header width
	startCol := colLetter(3) // first bid column (C)
	endCol := colLetter(len(p.BidItems) + 2)

	borderRow := startRow + 2 // bottom row of header

	startCell := fmt.Sprintf("%s%d", startCol, borderRow)
	endCell := fmt.Sprintf("%s%d", endCol, borderRow)

	if err := f.SetCellStyle(p.SheetName, startCell, endCell, style); err != nil {
		return err
	}

	return nil
}

func (p *ProgressSheet) createSection(f *excelize.File, section ProgressSheetSection, rowIndex int) error {
	// 1. Update style to include top AND bottom borders
	sectionHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Family: "Aptos Narrow"},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// 2. Calculate the end column letter based on BidItems
	// Column A (1) + Date Column (1) + len(BidItems)
	lastColIdx := len(p.BidItems) + 2
	endCol, _ := excelize.ColumnNumberToName(lastColIdx)

	from, _ := time.Parse("2006-01-02", section.FromDate)
	to, _ := time.Parse("2006-01-02", section.ToDate)
	sectionDateRange := fmt.Sprintf("%s - %s", from.Format("Jan 02"), to.Format("Jan 02"))

	// 3. Set the value in A
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("A%d", rowIndex), sectionDateRange); err != nil {
		return err
	}

	// 4. Apply the style to the entire row range (A to endCol)
	startCell := fmt.Sprintf("A%d", rowIndex)
	endCell := fmt.Sprintf("%s%d", endCol, rowIndex)

	return f.SetCellStyle(p.SheetName, startCell, endCell, sectionHeaderStyle)
}

func (p *ProgressSheet) createRow(f *excelize.File, row ProgressRow, rowIndex int) error {
	// Date
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("B%d", rowIndex), row.Date); err != nil {
		return err
	}

	// Quantities
	for col, bi := range p.BidItems {
		letter := colLetter(col + 3)
		if val, ok := row.Quantities[bi.Number]; ok {
			if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, rowIndex), val); err != nil {
				return err
			}
		}
	}

	return nil
}
