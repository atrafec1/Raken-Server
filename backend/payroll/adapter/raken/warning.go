package raken

import (
	"fmt"
	"prg_tools/payroll/dto"
)

func collectWarnings(timeCards []adapterTimeCard, equipLogs []adapterEquipLog) []dto.Warning {
	var warnings []dto.Warning
	warnings = append(warnings, checkMissingCostCode(timeCards, equipLogs)...)
	warnings = append(warnings, checkForDuplicateTimeCards(timeCards)...)
	warnings = append(warnings, findOrphanEquipLogs(equipLogs, timeCards)...)
	warnings = append(warnings, checkForMissingEquipOperator(equipLogs)...)
	return warnings
}

type timeCardDupeKey struct {
	EmployeeID string
	JobNumber  string
	Date       string
	CostCode   string
	PayType    string
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
			PayType:    timeCard.PayType,
		}
		if _, exists := seenEntries[key]; exists {
			warnings = append(warnings, dto.Warning{
				Message: fmt.Sprintf("%s %s %s %s %s",
					timeCard.EmployeeName, timeCard.JobNumber, timeCard.Date, timeCard.CostCode, timeCard.PayType),
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
				Message:     fmt.Sprintf("%s %s %s", timeCard.EmployeeName, timeCard.JobNumber, timeCard.Date),
				WarningType: "Time Card Missing Cost Code",
			})
		}
	}

	for _, equipLog := range equipLog {
		if equipLog.CostCode == "" {
			warnings = append(warnings, dto.Warning{
				Message: fmt.Sprintf("%s %s %s %s",
					equipLog.EmployeeName, equipLog.JobNumber, equipLog.Date, equipLog.EquipNumber),
				WarningType: "Equip Log Missing Cost Code",
			})
		}
	}
	return warnings
}

func checkForMissingEquipOperator(equipLogs []adapterEquipLog) []dto.Warning {
	var warnings []dto.Warning
	for _, equipLog := range equipLogs {
		if equipLog.EmployeeName == "" {
			warnings = append(warnings, dto.Warning{
				Message: fmt.Sprintf(
					"Equip #: %s Job: %s, Date: %s", equipLog.EquipNumber, equipLog.JobNumber, equipLog.Date),
				WarningType: "Equipment log with no operator",
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

		if _, exists := timecardSet[k]; !exists && eq.CostCode != "" {
			warnings = append(warnings,
				dto.Warning{
					Message:     eq.EmployeeName + eq.JobNumber + eq.Date + eq.CostCode,
					WarningType: "Equipment log entry with no matching time card entry",
				})
		}
	}

	return warnings
}
