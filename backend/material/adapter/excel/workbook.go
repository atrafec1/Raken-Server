package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/xuri/excelize/v2"
)

/*
===========================================
ORCHESTRATION (MULTI-SHEET WORKBOOK)
===========================================
*/

// CreateProgressWorkbook creates a single Excel file containing
// multiple ProgressSheets (one sheet per ProgressSheet).
func CreateProgressWorkbook(saveDir, fileName string, sheets []ProgressSheet) error {
	if len(sheets) == 0 {
		return fmt.Errorf("no sheets provided")
	}

	if !pathExists(saveDir) {
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed creating directory: %w", err)
		}
	}

	f := excelize.NewFile()
	renameBaseSheet(f, "Estimate - Daily Analysis")
	// Rename default sheet to first sheet name
	firstName := safeSheetName(sheets[0].SheetName, 1)
	renameBaseSheet(f, firstName)

	sheets[0].SheetName = firstName
	if err := sheets[0].writeToExistingSheet(f); err != nil {
		return err
	}

	// Remaining sheets
	for i := 1; i < len(sheets); i++ {
		name := safeSheetName(sheets[i].SheetName, i+1)

		if _, err := f.NewSheet(name); err != nil {
			return err
		}

		sheets[i].SheetName = name
		if err := sheets[i].writeToExistingSheet(f); err != nil {
			return err
		}
	}

	fullPath := filepath.Join(saveDir, fileName)
	fmt.Println("Saving progress workbook to: ", fullPath)
	return f.SaveAs(fullPath)
}

// Backward compatibility: single sheet export
func (p *ProgressSheet) CreateEstimateProgressSheet(saveDir, fileName string) error {
	return CreateProgressWorkbook(saveDir, fileName, []ProgressSheet{*p})
}

// Writes into an already-created sheet
func (p *ProgressSheet) writeToExistingSheet(f *excelize.File) error {
	currentRow := 3

	// Header
	if err := p.createHeader(f); err != nil {
		return err
	}
	currentRow += 3

	// Bid header
	if err := p.createBidHeader(f, currentRow); err != nil {
		return err
	}
	currentRow += 3

	// Rows (no sections anymore)
	for _, r := range p.Rows {
		if err := p.createRow(f, r, currentRow); err != nil {
			return err
		}
		currentRow++
	}

	// Totals row
	if err := p.createTotalsRow(f, currentRow); err != nil {
		return err
	}

	// Build string representation of rows
	for _, r := range p.Rows {
		rowVals := []string{r.Date}
		for _, bi := range p.BidItems {
			if val, ok := r.Quantities[bi.Number]; ok {
				rowVals = append(rowVals, fmt.Sprintf("%.2f", val))
			} else {
				rowVals = append(rowVals, "")
			}
		}
	}
	// Add totals row
	totals := p.CalculateTotals()
	totalVals := []string{"TOTAL"}
	for _, bi := range p.BidItems {
		totalVals = append(totalVals, fmt.Sprintf("%.2f", totals[bi.Number]))
	}
	if err := autoSizeColumnsByBidName(f, p.SheetName, p.BidItems); err != nil {
		return err
	}

	return nil
}

/*
===========================================
LAYOUT FUNCTIONS (UNCHANGED VISUALLY)
===========================================
*/

func (p *ProgressSheet) createHeader(f *excelize.File) error {
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
	if err != nil {
		return err
	}

	for col, bi := range p.BidItems {
		letter := colLetter(col + 3)

		cell := fmt.Sprintf("%s%d", letter, startRow)
		if err := f.SetCellValue(p.SheetName, cell, bi.Number); err != nil {
			return err
		}
		if err := f.SetCellStyle(p.SheetName, cell, cell, bidNumberStyle); err != nil {
			return err
		}
		if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, startRow+1), bi.Name); err != nil {
			return err
		}
		if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, startRow+2), bi.UnitOfMeasure); err != nil {
			return err
		}
	}

	// Bottom border
	style, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return err
	}

	startCol := colLetter(3)
	endCol := colLetter(len(p.BidItems) + 2)
	borderRow := startRow + 2

	startCell := fmt.Sprintf("%s%d", startCol, borderRow)
	endCell := fmt.Sprintf("%s%d", endCol, borderRow)

	return f.SetCellStyle(p.SheetName, startCell, endCell, style)
}

func (p *ProgressSheet) createRow(f *excelize.File, row ProgressRow, rowIndex int) error {
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("B%d", rowIndex), row.Date); err != nil {
		return err
	}

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

func (p *ProgressSheet) createTotalsRow(f *excelize.File, rowIndex int) error {
	totals := p.CalculateTotals()

	// Label in column A
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("A%d", rowIndex), "TOTAL"); err != nil {
		return err
	}

	// Style (bold + top border)
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return err
	}

	// Write totals aligned to bid columns
	for col, bi := range p.BidItems {
		letter := colLetter(col + 3) // C = first bid column

		if val, ok := totals[bi.Number]; ok {
			cell := fmt.Sprintf("%s%d", letter, rowIndex)

			if err := f.SetCellValue(p.SheetName, cell, val); err != nil {
				return err
			}
			if err := f.SetCellStyle(p.SheetName, cell, cell, style); err != nil {
				return err
			}
		}
	}

	// Style label cell
	labelCell := fmt.Sprintf("A%d", rowIndex)
	return f.SetCellStyle(p.SheetName, labelCell, labelCell, style)
}

func (p *ProgressSheet) CalculateTotals() map[string]float64 {
	totals := make(map[string]float64)

	for _, row := range p.Rows {
		for bid, qty := range row.Quantities {
			totals[bid] += qty
		}
	}

	return totals
}

/*
===========================================
HELPERS
===========================================
*/

// Auto-size columns based on string length of content
func autoSizeColumnsByBidName(f *excelize.File, sheet string, bidItems []BidItem) error {
	for col, bi := range bidItems {
		letter := colLetter(col + 3) // first bid column is C

		parts := strings.Split(bi.Name, " ")
		maxLenWord := slices.Max(parts)
		width := float64(len(maxLenWord) + 2)
		// Set the column width
		if err := f.SetColWidth(sheet, letter, letter, width); err != nil {
			return err
		}
	}
	return nil
}

func safeSheetName(name string, idx int) string {
	if name == "" {
		return fmt.Sprintf("Sheet-%d", idx)
	}
	return name
}
