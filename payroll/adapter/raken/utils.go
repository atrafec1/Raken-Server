package raken

import (
	"prg_tools/payroll/dto"
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
	if len(parts) >= 2 {
		result.Code = parts[1]
	}
	if len(parts) >= 3 {
		result.ChangeOrder = parts[2]
	}
	fmt.Println("Parsed cost code: ", result)
	fmt.Println("Original cost code: ", costCode)
	return result
}

type mergeKey struct {
	EmployeeName string
	JobNumber    string
	Date         string
	CostCode     string
}

func mergeTimeAndEquipLogs(timeCards []adapterTimeCard, equipLogs []adapterEquipLog) ([]*dto.PayrollEntry, error) {
	var payrollEntries []*dto.PayrollEntry
	entries := make(map[mergeKey]*dto.PayrollEntry)

	for _, timeCard := range timeCards {

		key := mergeKey{
			EmployeeName: timeCard.EmployeeName,
			JobNumber:    timeCard.JobNumber,
			Date:         timeCard.Date,
			CostCode:     timeCard.CostCode,
		}

		entry, exists := entries[key]

		if !exists {
			newEntry, err := newBasePayrollEntry(timeCard)
			if err != nil {
				return nil, fmt.Errorf("Failed to make base payroll entry: %w")
			}
			entries[key] = newEntry
			payrollEntries = append(payrollEntries, newEntry)

		} else {
			applyTimeCard(entry, timeCard)
		}
	}
	for _, equipLog := range equipLogs {
		key := mergeKey{
			Date:         equipLog.Date,
			CostCode:     equipLog.CostCode,
			JobNumber:    equipLog.JobNumber,
			EmployeeName: equipLog.EmployeeName,
		}
		entry, exists := entries[key]
		if !exists {
			fmt.Println("Should exist but doesnt: %v", key)
			continue
		}
		if entry.SpecialPayCode == "" {
			applyEquipLog(entry, equipLog)
		} else {
			clone := *entry
			clone.RegularHours = 0
			clone.OvertimeHours = 0
			clone.PremiumHours = 0
			applyEquipLog(&clone, equipLog)
			payrollEntries = append(payrollEntries, &clone)
		}

	}
	return payrollEntries, nil
}

func CopySlice[T any](ptrSlice []*T) []T {
	varSlice := make([]T, len(ptrSlice))
	for i, p := range ptrSlice {
		varSlice[i] = *p
	}
	return varSlice
}
func applyTimeCard(entry *dto.PayrollEntry, timeCard adapterTimeCard) {
	payRoute := routePay(timeCard)
	entry.PremiumHours += payRoute.PremiumHours
	entry.OvertimeHours += payRoute.OvertimeHours
	entry.RegularHours += payRoute.RegularHours
}
func applyEquipLog(entry *dto.PayrollEntry, equipLog adapterEquipLog) {
	entry.SpecialPayType = "EQP"
	entry.SpecialPayCode = equipLog.EquipNumber
	entry.SpecialUnits = equipLog.Hours
}

func newBasePayrollEntry(timeCard adapterTimeCard) (*dto.PayrollEntry, error) {
	payRoute := routePay(timeCard)
	costCodeParts := parseRakenCostCode(timeCard.CostCode)
	day, err := convertDateToInt(timeCard.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to convert date to int: %w")
	}
	return &dto.PayrollEntry{
		CurrentDate:   timeCard.Date,
		EmployeeCode:  timeCard.EmployeeCode,
		JobNumber:     timeCard.JobNumber,
		Day:           day,
		Phase:         costCodeParts.Phase,
		CostCode:      costCodeParts.Code,
		ChangeOrder:   costCodeParts.ChangeOrder,
		PremiumHours:  payRoute.PremiumHours,
		OvertimeHours: payRoute.OvertimeHours,
		RegularHours:  payRoute.RegularHours,
	}, nil
}

type payRouting struct {
	RegularHours  float64
	PremiumHours  float64
	OvertimeHours float64
}

// Handles different pay types
func routePay(timeCard adapterTimeCard) payRouting {
	switch timeCard.PayType {
	case "RT":
		return payRouting{
			RegularHours: timeCard.Hours,
		}
	case "OT":
		return payRouting{
			OvertimeHours: timeCard.Hours,
		}
	case "DT":
		return payRouting{
			PremiumHours: timeCard.Hours,
		}
	default:
		return payRouting{
			RegularHours: timeCard.Hours,
		}
	}
}

