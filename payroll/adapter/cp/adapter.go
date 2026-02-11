package cp

import (
	"daily_check_in/payroll/dto"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Adapter struct {
	CPPath string
}

func NewAdapter(path string) *Adapter {
	return &Adapter{
		CPPath: path,
	}
}
func (a *Adapter) ExportPayrollEntries(entries []dto.PayrollEntry) error {
	cpCSVPath := filepath.Join(a.CPPath, "PAYROLLTIMECARD.CSV")

	file, err := os.Create(cpCSVPath)
	if err != nil {
		return fmt.Errorf("failed to create csv file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, entry := range entries {
		row := []string{
			entry.EmployeeCode,
			entry.CurrentDate,
			entry.CraftLevel,
			entry.JobNumber,
			entry.Phase,
			entry.CostCode,
			entry.ChangeOrder,
			strconv.FormatFloat(entry.RegularHours, 'f', -1, 64),
			strconv.FormatFloat(entry.OvertimeHours, 'f', -1, 64),
			strconv.FormatFloat(entry.PremiumHours, 'f', -1, 64),
			strconv.Itoa(entry.Day),
			entry.DownFlag,
			entry.SpecialPayType,
			entry.SpecialPayCode,
			strconv.FormatFloat(entry.SpecialUnits, 'f', -1, 64),
			strconv.FormatFloat(entry.SpecialRate, 'f', -1, 64),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	return nil
}
