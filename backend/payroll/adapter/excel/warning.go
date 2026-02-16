package excel

import (
	"fmt"
	"path/filepath"
	"prg_tools/payroll/dto"

	"github.com/xuri/excelize/v2"
)

func (e *ExcelPayrollExporter) ExportWarnings(warnings []dto.Warning) error {
	if len(warnings) == 0 {
		return nil
	}

	filePath := filepath.Join(e.warningsDir, "payrollwarnings.xlsx")

	f := excelize.NewFile()

	sheet := "PayrollWarnings"
	f.SetSheetName("Sheet1", sheet)

	// Create styles
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 16,
		},
	})
	if err != nil {
		return err
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return err
	}

	grouped := sortWarningsByType(warnings)

	row := 1

	// Title
	f.SetCellValue(sheet, "A1", "Payroll Warnings")
	f.SetCellStyle(sheet, "A1", "A1", titleStyle)
	row += 2 // spacing

	// Iterate groups
	for warningType, list := range grouped {
		// Header
		cell := fmt.Sprintf("A%d", row)
		f.SetCellValue(sheet, cell, warningType)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		row++

		// Values
		for _, w := range list {
			cell := fmt.Sprintf("A%d", row)
			f.SetCellValue(sheet, cell, "- "+w.Message)
			row++
		}

		row++ // blank line between groups
	}
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save warnings excel file: %w", err)
	}
	return f.Close()
}

func sortWarningsByType(warnings []dto.Warning) map[string][]dto.Warning {
	grouped := make(map[string][]dto.Warning)

	for _, w := range warnings {
		grouped[w.WarningType] = append(grouped[w.WarningType], w)
	}

	return grouped
}
