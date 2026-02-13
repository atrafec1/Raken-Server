package raken

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"prg_tools/external/rakenapi"
	"prg_tools/payroll/dto"
)

const (
	timeCardJsonPath = "./mocks/timecards.json"
	equipLogJsonPath = "./mocks/equiplogs.json"
	projectJsonPath  = "./mocks/projects.json"
	employeeJsonPath = "./mocks/employees.json"
)

// --- Helpers for loading JSON mocks ---

func loadJSON[T any](path string) (T, error) {
	var result T
	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

// --- Mock client ---

type RakenClientMock struct {
	TimeCards rakenapi.TimeCardResponse
	EquipLogs rakenapi.EquipmentLogResponse
	Projects  rakenapi.ProjectResponse
	Employees rakenapi.EmployeeResponse
}

func (m *RakenClientMock) GetTimeCards(from, to string) (*rakenapi.TimeCardResponse, error) {
	return &m.TimeCards, nil
}

func (m *RakenClientMock) GetEquipmentLogs(from, to string) (*rakenapi.EquipmentLogResponse, error) {
	return &m.EquipLogs, nil
}

func (m *RakenClientMock) GetProjects() (*rakenapi.ProjectResponse, error) {
	return &m.Projects, nil
}

func (m *RakenClientMock) GetEmployees() (*rakenapi.EmployeeResponse, error) {
	return &m.Employees, nil
}

// --- Global adapter ---

var (
	testAdapter     *RakenAPIAdapter
	testProjectMap  map[string]rakenapi.Project
	testEmployeeMap map[string]rakenapi.Employee
)

func TestMain(m *testing.M) {
	timeCards, _ := loadJSON[rakenapi.TimeCardResponse](timeCardJsonPath)
	equipLogs, _ := loadJSON[rakenapi.EquipmentLogResponse](equipLogJsonPath)
	projects, _ := loadJSON[rakenapi.ProjectResponse](projectJsonPath)
	employees, _ := loadJSON[rakenapi.EmployeeResponse](employeeJsonPath)

	mockClient := &RakenClientMock{
		TimeCards: timeCards,
		EquipLogs: equipLogs,
		Projects:  projects,
		Employees: employees,
	}
	testAdapter = &RakenAPIAdapter{Client: mockClient}
	testProjectMap = make(map[string]rakenapi.Project)
	for _, p := range projects.Collection {
		testProjectMap[p.UUID] = p
	}

	testEmployeeMap = make(map[string]rakenapi.Employee)
	for _, e := range employees.Collection {
		testEmployeeMap[e.UUID] = e
	}
	os.Exit(m.Run())
}

// --- Tests ---

func TestGetPayrollEntries(t *testing.T) {
	result, err := testAdapter.GetPayrollEntries("2026-01-01", "2026-01-10")
	if err != nil {
		t.Fatalf("GetPayrollEntries() error: %v", err)
	}

	if len(result.Entries) == 0 {
		t.Errorf("expected at least 1 payroll entry, got 0")
	}

	first := result.Entries[0]
	if first.SpecialPayCode != "" && first.SpecialPayType != "PAY" {
		t.Errorf("expected SpecialPayType 'PAY' for special pay code, got %s", first.SpecialPayType)
	}

	specialUnits := sumSpecialUnits(result.Entries)
	if specialUnits == 0 {
		t.Errorf("not capturing special units")
	}

}

func TestMergeTimeAndEquipLogs(t *testing.T) {
	timeCards, _ := loadJSON[rakenapi.TimeCardResponse](timeCardJsonPath)
	equipLogs, _ := loadJSON[rakenapi.EquipmentLogResponse](equipLogJsonPath)

	// Normalize responses
	normalizedTimeCards, err := testAdapter.normalizeTimeCardResponse(timeCards, testProjectMap, testEmployeeMap)
	if err != nil {
		t.Fatalf("normalizeTimeCardResponse() error: %v", err)
	}

	normalizedEquipLogs, err := testAdapter.normalizeEquipLogResponse(equipLogs, testProjectMap, testEmployeeMap)
	if err != nil {
		t.Fatalf("normalizeEquipLogResponse() error: %v", err)
	}

	// Validate we have data to test
	if len(normalizedTimeCards) == 0 {
		t.Fatal("expected at least 1 normalized time card, got 0")
	}
	if len(normalizedEquipLogs) == 0 {
		t.Fatal("expected at least 1 normalized equipment log, got 0")
	}

	// Build merge key index to identify which equipment logs should be included
	// Equipment logs without matching time cards are intentionally skipped
	timeCardKeys := buildMergeKeySet(normalizedTimeCards)

	// Calculate input hours
	timeCardHours := sumTimeCardHours(normalizedTimeCards)
	matchedEquipLogHours := sumMatchedEquipLogHours(normalizedEquipLogs, timeCardKeys)
	inputTotal := timeCardHours + matchedEquipLogHours

	// Perform merge
	mergedEntries, err := testAdapter.mergeTimeAndEquipLogs(normalizedTimeCards, normalizedEquipLogs)
	if err != nil {
		t.Fatalf("mergeTimeAndEquipLogs() error: %v", err)
	}

	if len(mergedEntries) == 0 {
		t.Fatal("expected at least 1 merged payroll entry, got 0")
	}

	// Calculate output hours
	outputTotal := sumPayrollEntryHours(mergedEntries)

	// Verify hours are preserved exactly (not lost, not added)
	if inputTotal != outputTotal {
		t.Errorf("Hours mismatch: input total = %.2f, output total = %.2f (time cards: %.2f, matched equip logs: %.2f)",
			inputTotal, outputTotal, timeCardHours, matchedEquipLogHours)
	}
}

// Helper functions for hour calculations

func sumTimeCardHours(timeCards []adapterTimeCard) float64 {
	var total float64
	for _, tc := range timeCards {
		total += tc.Hours
	}
	return total
}

func buildMergeKeySet(timeCards []adapterTimeCard) map[mergeKey]bool {
	keySet := make(map[mergeKey]bool)
	for _, tc := range timeCards {
		key := mergeKey{
			EmployeeName: tc.EmployeeName,
			JobNumber:    tc.JobNumber,
			Date:         tc.Date,
			CostCode:     tc.CostCode,
		}
		keySet[key] = true
	}
	return keySet
}

func sumMatchedEquipLogHours(equipLogs []adapterEquipLog, timeCardKeys map[mergeKey]bool) float64 {
	var total float64
	for _, el := range equipLogs {
		key := mergeKey{
			EmployeeName: el.EmployeeName,
			JobNumber:    el.JobNumber,
			Date:         el.Date,
			CostCode:     el.CostCode,
		}
		// Only count equipment log hours that have a matching time card
		if timeCardKeys[key] {
			total += el.Hours
		}
	}
	return total
}

func sumPayrollEntryHours(entries []*dto.PayrollEntry) float64 {
	var total float64
	for _, entry := range entries {
		total += entry.RegularHours + entry.OvertimeHours + entry.PremiumHours + entry.SpecialUnits
	}
	return total
}

func sumSpecialUnits(entries []dto.PayrollEntry) float64 {
	result := 0.0
	for _, entry := range entries {
		result += entry.SpecialUnits
	}
	return result
}

func TestNormalizeTimeCardResponse(t *testing.T) {
	timeCards, _ := loadJSON[rakenapi.TimeCardResponse](timeCardJsonPath)

	normalized, err := testAdapter.normalizeTimeCardResponse(timeCards, testProjectMap, testEmployeeMap)
	if err != nil {
		t.Fatalf("NormalizeTimeCardResponse() error: %v", err)
	}

	if len(normalized) == 0 {
		t.Errorf("expected normalized time cards, got 0")
	}
}
func TestCAOvertimeRule(t *testing.T) {
	saturdayEntry := dto.PayrollEntry{Day: 6, RegularHours: 8}
	sundayEntry := dto.PayrollEntry{Day: 7, RegularHours: 8, OvertimeHours: 4}
	specialSaturdayEntry := dto.PayrollEntry{Day: 6, RegularHours: 8, PremiumHours: 4, OvertimeHours: 5}
	testAdapter.applyCAOvertimeRules(&saturdayEntry)
	if saturdayEntry.RegularHours != 0 {
		fmt.Printf("%+v", saturdayEntry)
		t.Errorf("CA overtime rule not applied: regular hours should be 0 on sat")
	}

	testAdapter.applyCAOvertimeRules(&sundayEntry)
	if sundayEntry.RegularHours != 0 || sundayEntry.OvertimeHours != 0 {
		t.Errorf("CA overtime rule not applied: all hours on sunday should be DT")
	}
	if sundayEntry.PremiumHours != 12 {
		t.Errorf("CA overtime rules added reg and overtime hours to DT incorrectly")
	}

	testAdapter.applyCAOvertimeRules(&specialSaturdayEntry)
	if specialSaturdayEntry.PremiumHours != 5 {
		t.Errorf("CA Overtime rules didn't apply premium hours for saturday")
	}
}
func TestApplyPayrollRules(t *testing.T) {
	entries := []*dto.PayrollEntry{
		{RegularHours: 10, Day: 6, Phase: "VACNJB"},
		{RegularHours: 8, Day: 7},
	}

	testAdapter.applyPayrollRules(entries)

	if entries[0].RegularHours != 0 || entries[0].SpecialPayType != "PAY" {
		t.Errorf("special phase rule not applied correctly")
	}

	if entries[1].RegularHours != 0 || entries[1].PremiumHours != 8 {
		t.Errorf("CA overtime rule not applied correctly for day 7")
	}
}
