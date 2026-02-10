package adapter

import (
	"daily_check_in/payroll/dto"
	"fmt"
	"strings"
	"time"
)

func convertDateToInt(date string) (int, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %w", err)
	}
	day := t.Weekday()
	if day == time.Sunday {
		return 7, nil
	}
	return int(day), nil
}

//Merge rules:

//Merge on same day, cost code, job number, employee
//For equipment logs with same merge key only one row gets employee hours rest are rows without regular time etc
//Raken cost code is phase/code/change order

type rakenCostCode struct {
	Phase       string
	Code        string
	ChangeOrder string
}

func parseRakenCostCode(costCode string) rakenCostCode {
	parts := strings.Split(costCode, "/")
	var result rakenCostCode
	result.Phase = parts[0]
	if len(pats) >= 2 {
		result.Code = parts[1]
	}
	if len(parts) >= 3 {
		result.ChangeOrder = parts[2]
	}
}

type mergeKey struct {
	EmployeeName string
	JobNumber    string
	Date         string
	CostCode     string
}

func mergeTimeAndEquipLogs(timeCards []adapterTimeCard, equipLogs []adapterEquipLog) {
	var payrollEntries []dto.PayrollEntry
	groupedEntries := make(map[mergeKey]*dto.PayrollEntry)

	for _, timeCard := range timeCards {
		key := mergeKey{
			EmployeeName: timeCard.EmployeeName,
			JobNumber:    timeCard.JobNumber,
			Date:         timeCard.Date,
			CostCode:     timeCard.CostCode,
		}
		entry, exists := groupedEntries[key]

		if !exists {
			rakenCode := parseRakenCostCode(timeCard.CostCode)

			groupedEntries[key] = &dto.PayrollEntry{
				CurrentDate: timeCard.Date,
				Phase:       rakenCode.Phase,
				CostCode:    rakenCode.Code,
				ChangeOrder: rakenCode.ChangeOrder,
			}
		} else {
			updatedEntry = dto.PayrollEntry{
,
		}

	}

}
