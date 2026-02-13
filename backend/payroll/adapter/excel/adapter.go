package excel

import (
	"fmt"
	"os"
	"prg_tools/payroll/dto"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelPayrollExporter struct {
	OutputPath     string
	WeekEndingDate string
}

func NewExcelPayrollExporter(outputPath string) *ExcelPayrollExporter {
	return &ExcelPayrollExporter{
		OutputPath:     outputPath,
		WeekEndingDate: "",
	}
}

func (e *ExcelPayrollExporter) SetWeekEndingDate(weekEndingDate string) {
	e.WeekEndingDate = weekEndingDate
}

func (e *ExcelPayrollExporter) ExportPayrollEntries(rawEntries []dto.PayrollEntry) error {

	// ---- BEGIN EXCEL CREATION ----
	f := excelize.NewFile()
	defer f.Close()

	weekEnd := getWeekEndingDate(rawEntries)
	e.SetWeekEndingDate(weekEnd)
	entries := transformPayrollEntries(rawEntries)
	sheetName := "Sheet1"

	// Styles
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Family: "Aptos Narrow"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "000000", Style: 1}},
	})

	boldStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Family: "Aptos Narrow"},
	})

	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Family: "Aptos Narrow"},
	})

	regularStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: "Aptos Narrow"},
	})

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 24, Family: "Aptos Narrow"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	thickTopStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Family: "Aptos Narrow"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 2}},
	})

	// Header
	f.SetCellValue(sheetName, "A1", "PACIFIC RESTORATION GROUP, INC.")
	f.SetCellStyle(sheetName, "A1", "A1", boldStyle)

	f.SetCellValue(sheetName, "C2", "PAYROLL")
	f.SetCellStyle(sheetName, "C2", "C2", titleStyle)

	f.SetCellValue(sheetName, "A3", fmt.Sprintf("Week Ending: %s", e.WeekEndingDate))
	f.SetCellStyle(sheetName, "A3", "A3", boldStyle)

	f.SetCellValue(sheetName, "A4", "Notes:")
	f.SetCellStyle(sheetName, "A4", "A4", boldStyle)

	// Column headers
	headers := []string{"EID", "Day", "Date", "Class", "Job #", "Cost Code #", "Cost Code", "RT", "OT", "Premium T", "Equip #", "Equip Hours"}
	headerRow := 6

	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", string(rune('A'+i)), headerRow)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Column widths
	widths := map[string]float64{
		"A": 8, "B": 4, "C": 10, "D": 8, "E": 8, "F": 12,
		"G": 14, "H": 8, "I": 6, "J": 9, "K": 18, "L": 18,
	}
	for col, width := range widths {
		f.SetColWidth(sheetName, col, col, width)
	}

	row := headerRow + 1
	var prevEID string
	var employeeRT, employeeOT, employeePT float64
	var totalRT, totalOT, totalPT float64

	for i, entry := range entries {

		if prevEID != "" && entry.EmployeeCode != prevEID {
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "Totals:")
			f.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), boldStyle)

			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), employeeRT)
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), employeeOT)
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), employeePT)
			f.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("J%d", row), centerStyle)

			row += 2
			employeeRT, employeeOT, employeePT = 0, 0, 0
		}

		equipCodes := strings.Join(entry.EquipmentCode, ", ")
		equipHours := strings.Join(entry.EquipmentHours, ", ")

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), entry.EmployeeCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), entry.Day)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), entry.Date)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), entry.Class)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), entry.JobNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), entry.CostCodeNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), entry.CostCode)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), entry.RegularHours)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), entry.OvertimeHours)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), entry.PremiumHours)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), equipCodes)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), equipHours)

		centerCols := []string{"A", "C", "E", "F", "H", "I", "J", "K", "L"}
		for _, col := range centerCols {
			f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("%s%d", col, row), centerStyle)
		}

		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), regularStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), regularStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), regularStyle)

		employeeRT += entry.RegularHours
		employeeOT += entry.OvertimeHours
		employeePT += entry.PremiumHours

		totalRT += entry.RegularHours
		totalOT += entry.OvertimeHours
		totalPT += entry.PremiumHours

		prevEID = entry.EmployeeCode
		row++

		if i == len(entries)-1 {
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "Totals:")
			f.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), boldStyle)

			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), employeeRT)
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), employeeOT)
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), employeePT)
			f.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("J%d", row), centerStyle)
			row++
		}
	}

	row += 1

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "TOTAL HOURS:")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)

	f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), totalRT)
	f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), totalOT)
	f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), totalPT)
	f.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("J%d", row), thickTopStyle)

	if e.OutputPath == "" {
		e.OutputPath = "."
	}

	if err := os.MkdirAll(e.OutputPath, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/payroll_%s.xlsx",
		e.OutputPath,
		time.Now().Format("20060102_150405"),
	)

	return f.SaveAs(filename)
}
