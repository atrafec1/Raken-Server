package cp

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"prg_tools/payroll/dto"
	"sort"
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
	entries = sortPayrollEntries(entries)
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

func sortPayrollEntries(entries []dto.PayrollEntry) []dto.PayrollEntry {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].EmployeeCode == entries[j].EmployeeCode {
			return entries[i].CurrentDate < entries[j].CurrentDate
		}
		return entries[i].EmployeeCode < entries[j].EmployeeCode
	})
	return entries
}
