package excel

import (
	"fmt"
	"prg_tools/payroll/dto"
	"sort"
	"time"
)

type ExcelPayrollEntry struct {
	EmployeeCode   string
	Day            string
	Date           string
	Class          string
	JobNumber      string
	CostCodeNumber string
	CostCode       string
	RegularHours   float64
	OvertimeHours  float64
	PremiumHours   float64
	EquipmentCode  []string
	EquipmentHours []string
}
type payrollKey struct {
	EmployeeCode string
	Date         string
	JobNumber    string
	CostCode     string
}

func transformPayrollEntries(entries []dto.PayrollEntry) []ExcelPayrollEntry {
	grouped := make(map[payrollKey]*ExcelPayrollEntry)

	for _, e := range entries {

		key := payrollKey{
			EmployeeCode: e.EmployeeCode,
			Date:         e.CurrentDate,
			JobNumber:    e.JobNumber,
			CostCode:     getCostCodeNumber(e.Phase, e.CostCode, e.ChangeOrder),
		}

		if existing, ok := grouped[key]; ok {
			// Append equipment if present
			if e.SpecialPayType == "EQP" {
				existing.EquipmentCode = append(existing.EquipmentCode, e.SpecialPayCode)
				existing.EquipmentHours = append(existing.EquipmentHours,
					fmt.Sprintf("%.2f", e.SpecialUnits))
			}
			continue
		}

		// Create new Excel entry
		excelEntry := newExcelEntry(e)

		if e.SpecialPayCode != "" && e.SpecialPayType == "EQP" {
			excelEntry.EquipmentCode = []string{e.SpecialPayCode}
			excelEntry.EquipmentHours = []string{
				fmt.Sprintf("%.2f", e.SpecialUnits),
			}
		}

		grouped[key] = excelEntry
	}

	// Convert map to slice
	result := make([]ExcelPayrollEntry, 0, len(grouped))
	for _, v := range grouped {
		result = append(result, *v)
	}

	sortExcelEntries(result)
	return result
}
func newExcelEntry(payrollEntry dto.PayrollEntry) *ExcelPayrollEntry {
	if payrollEntry.SpecialPayType == "PAY" {
		return &ExcelPayrollEntry{
			EmployeeCode:   payrollEntry.EmployeeCode,
			Day:            getDayName(payrollEntry.Day),
			Date:           payrollEntry.CurrentDate,
			Class:          payrollEntry.CraftLevel,
			JobNumber:      payrollEntry.JobNumber,
			CostCodeNumber: payrollEntry.SpecialPayCode,
			RegularHours:   payrollEntry.SpecialUnits,
			OvertimeHours:  0,
			PremiumHours:   0,
			CostCode:       payrollEntry.CostCodeDivision,
		}
	}
	return &ExcelPayrollEntry{
		EmployeeCode:   payrollEntry.EmployeeCode,
		Day:            getDayName(payrollEntry.Day),
		Date:           payrollEntry.CurrentDate,
		Class:          payrollEntry.CraftLevel,
		JobNumber:      payrollEntry.JobNumber,
		CostCodeNumber: getCostCodeNumber(payrollEntry.Phase, payrollEntry.CostCode, payrollEntry.ChangeOrder),
		RegularHours:   payrollEntry.RegularHours,
		OvertimeHours:  payrollEntry.OvertimeHours,
		PremiumHours:   payrollEntry.PremiumHours,
		CostCode:       payrollEntry.CostCodeDivision,
	}
}

func sortExcelEntries(entries []ExcelPayrollEntry) {
	sort.Slice(entries, func(i, j int) bool {

		// 1️⃣ EmployeeCode
		if entries[i].EmployeeCode != entries[j].EmployeeCode {
			return entries[i].EmployeeCode < entries[j].EmployeeCode
		}

		// 2️⃣ Date (must parse, not string compare)
		dateI, _ := time.Parse("2006-01-02", entries[i].Date)
		dateJ, _ := time.Parse("2006-01-02", entries[j].Date)

		if !dateI.Equal(dateJ) {
			return dateI.Before(dateJ)
		}

		// 3️⃣ JobNumber
		if entries[i].JobNumber != entries[j].JobNumber {
			return entries[i].JobNumber < entries[j].JobNumber
		}

		// 4️⃣ CostCodeNumber
		return entries[i].CostCodeNumber < entries[j].CostCodeNumber
	})
}

func getWeekEndingDate(entries []dto.PayrollEntry) string {
	//getsunday of the week given a date
	var weekEnd string
	for _, entry := range entries {
		if entry.CurrentDate != "" {
			date, err := time.Parse("2006-01-02", entry.CurrentDate)
			if err != nil {
				return ""
			}
			weekEnd = date.AddDate(0, 0, (7-int(date.Weekday()))%7).Format("2006-01-02")
			break
		}
	}
	return weekEnd
}

func getWeekBeginningDate(entries []dto.PayrollEntry) string {
	//getsunday of the week given a date
	var weekStart string
	for _, entry := range entries {
		if entry.CurrentDate != "" {
			date, err := time.Parse("2006-01-02", entry.CurrentDate)
			if err != nil {
				return ""
			}
			weekStart = date.AddDate(0, 0, -int(date.Weekday())).Format("2006-01-02")
			break
		}
	}
	return weekStart
}

func getDayName(day int) string {
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	if day >= 0 && day < 7 {
		return days[day]
	}
	return ""
}

func getCostCodeNumber(phase, costCode, changeOrder string) string {
	costCodeStr := phase
	if costCode != "" {
		costCodeStr += fmt.Sprintf("/%s", costCode)
	}
	if changeOrder != "" {
		costCodeStr += fmt.Sprintf("-%s", changeOrder)
	}
	return costCodeStr
}
