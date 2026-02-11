package raken

import (
	"daily_check_in/payroll/dto"
	"fmt"
)

func collectWarnings(timeCards []adapterTimeCard, equipLogs []adapterEquipLog) []dto.Warning {
	var warnings []dto.Warning
	warnings = append(warnings, checkMissingCostCode(timeCards, equipLogs)...)
	warnings = append(warnings, checkForDuplicateTimeCards(timeCards)...)
	warnings = append(warnings, findOrphanEquipLogs(equipLogs, timeCards)...)
	return warnings
}

type timeCardDupeKey struct {
	EmployeeID string
	JobNumber  string
	Date       string
	CostCode   string
}

func checkForDuplicateTimeCards(timeCards []adapterTimeCard) []dto.Warning {
	var warnings []dto.Warning
	seenEntries := make(map[timeCardDupeKey]struct{})

	for _, timeCard := range timeCards {
		key := timeCardDupeKey{
			EmployeeID: timeCard.EmployeeCode,
			JobNumber:  timeCard.JobNumber,
			Date:       timeCard.Date,
			CostCode:   timeCard.CostCode,
		}
		if _, exists := seenEntries[key]; exists {
			warnings = append(warnings, dto.Warning{
				Message:     fmt.Sprintf("Duplicate time card entry for employee %s on job %s for date %s and cost code %s", timeCard.EmployeeCode, timeCard.JobNumber, timeCard.Date, timeCard.CostCode),
				WarningType: "Duplicate Time Card Entry",
			})
		} else {
			seenEntries[key] = struct{}{}
		}
	}
	return warnings
}

func checkMissingCostCode(timeCards []adapterTimeCard, equipLog []adapterEquipLog) []dto.Warning {
	var warnings []dto.Warning
	for _, timeCard := range timeCards {
		if timeCard.CostCode == "" {
			warnings = append(warnings, dto.Warning{
				Message:     timeCard.EmployeeName + timeCard.JobNumber + timeCard.Date,
				WarningType: "Time Card Missing Cost Code",
			})
		}
	}

	for _, equipLog := range equipLog {
		if equipLog.CostCode == "" {
			warnings = append(warnings, dto.Warning{
				Message:     equipLog.EmployeeName + equipLog.JobNumber + equipLog.Date + equipLog.EquipNumber,
				WarningType: "Equip Log Missing Cost Code",
			})
		}
	}
	return warnings
}

func findOrphanEquipLogs(equipLogs []adapterEquipLog, timeCards []adapterTimeCard) []dto.Warning {
	var warnings []dto.Warning

	timecardSet := make(map[mergeKey]struct{})
	for _, tc := range timeCards {
		k := mergeKey{
			EmployeeName: tc.EmployeeName,
			Date:         tc.Date,
			CostCode:     tc.CostCode,
			JobNumber:    tc.JobNumber,
		}
		timecardSet[k] = struct{}{}
	}

	for _, eq := range equipLogs {
		k := mergeKey{
			EmployeeName: eq.EmployeeName,
			Date:         eq.Date,
			CostCode:     eq.CostCode,
			JobNumber:    eq.JobNumber,
		}
		if _, exists := timecardSet[k]; !exists {
			warnings = append(warnings,
				dto.Warning{
					Message:     eq.EmployeeName + eq.JobNumber + eq.Date + eq.CostCode,
					WarningType: "Equipment log entry with no matching time card entry",
				})
		}
	}

	return warnings
}
