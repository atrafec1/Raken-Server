package raken

import (
	"fmt"
	"prg_tools/external/rakenapi"
	"prg_tools/payroll/dto"
	"prg_tools/payroll/port"
)

type RakenAPIAdapter struct {
	Client *rakenapi.Client
}

func NewRakenAPIAdapter() (*RakenAPIAdapter, error) {
	config, err := rakenapi.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load raken api config: %w", err)
	}
	client, err := rakenapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to make raken api client")
	}

	return &RakenAPIAdapter{
		Client: client,
	}, nil
}

func (r *RakenAPIAdapter) GetPayrollEntries(fromDate, toDate string) (port.PayrollEntryResult, error) {

	timeCardResponse, err := r.Client.GetTimeCards(fromDate, toDate)
	if err != nil {
		return port.PayrollEntryResult{}, err
	}
	equipLogResponse, err := r.Client.GetEquipmentLogs(fromDate, toDate)
	if err != nil {
		return port.PayrollEntryResult{}, err
	}

	projectMap, err := r.makeProjectMap()
	if err != nil {
		return port.PayrollEntryResult{}, err
	}
	employeeMap, err := r.makeEmployeeMap()
	if err != nil {
		return port.PayrollEntryResult{}, err
	}
	adapterTimeCards, err := normalizeTimeCardResponse(*timeCardResponse, projectMap, employeeMap)
	if err != nil {
		return port.PayrollEntryResult{}, err
	}
	adapterEquipLogs, err := normalizeEquipLogResponse(*equipLogResponse, projectMap, employeeMap)
	if err != nil {
		return port.PayrollEntryResult{}, err
	}
	mergedLogs, err := mergeTimeAndEquipLogs(adapterTimeCards, adapterEquipLogs)
	if err != nil {
		return port.PayrollEntryResult{}, fmt.Errorf("failed to merge time cards and equip logs: %w", err)
	}
	applyPayrollRules(mergedLogs)
	warnings := collectWarnings(adapterTimeCards, adapterEquipLogs)
	return port.PayrollEntryResult{
		Entries:  CopySlice(mergedLogs),
		Warnings: warnings}, nil
}

type adapterTimeCard struct {
	EmployeeCode        string
	EmployeeName        string
	Date                string
	JobNumber           string
	Class               string
	CostCode            string
	PayType             string
	Hours               float64
	CostCodeDescription string
}

type adapterEquipLog struct {
	EmployeeName string
	Hours        float64
	Date         string
	JobNumber    string
	EquipNumber  string
	CostCode     string
}

// Need to fetch all projects (to map their uuid to metadata)
func (r *RakenAPIAdapter) makeProjectMap() (map[string]rakenapi.Project, error) {
	projectMap := make(map[string]rakenapi.Project)
	projectsResp, err := r.Client.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to get projects from raken api: %w", err)
	}
	for _, project := range projectsResp.Collection {
		projectMap[project.UUID] = project
	}
	return projectMap, nil
}

func (r *RakenAPIAdapter) makeEmployeeMap() (map[string]rakenapi.Employee, error) {
	employeeMap := make(map[string]rakenapi.Employee)
	employeeResp, err := r.Client.GetEmployees()

	if err != nil {
		return nil, fmt.Errorf("failed to get employees from raken api: %w", err)
	}
	for _, employee := range employeeResp.Collection {
		employeeMap[employee.UUID] = employee
	}
	return employeeMap, nil
}
func applyPayrollRules(entries []*dto.PayrollEntry) {
	for _, entry := range entries {
		applySpecialPhaseRules(entry)
		applyCAOvertimeRules(entry)
	}
}

func applySpecialPhaseRules(entry *dto.PayrollEntry) {
	var specialPhasePayCodes = map[string]string{
		"VACNJB": "VACNJB",
		"SKLVCA": "SKLVCA",
		// add more as needed
	}
	if specialPayCode, exists := specialPhasePayCodes[entry.Phase]; exists {
		entry.SpecialPayType = "PAY"
		entry.SpecialPayCode = specialPayCode
		entry.SpecialUnits = entry.RegularHours
		entry.RegularHours = 0
		entry.OvertimeHours = 0
		entry.PremiumHours = 0
		entry.Phase = "0"
	}
}

func applyCAOvertimeRules(entry *dto.PayrollEntry) {
	switch entry.Day {
	case 6:
		threshold := 12.0
		overTimeHours := entry.RegularHours + entry.OvertimeHours
		if overTimeHours > threshold {
			overAmount := overTimeHours - threshold
			entry.PremiumHours += overAmount
			entry.OvertimeHours = threshold
		}
	case 7:
		totalHours := entry.RegularHours + entry.OvertimeHours + entry.PremiumHours
		entry.PremiumHours = totalHours
		entry.RegularHours = 0
		entry.OvertimeHours = 0
	}
}

func normalizeTimeCardResponse(
	timeCardResponse rakenapi.TimeCardResponse,
	projectMap map[string]rakenapi.Project,
	employeeMap map[string]rakenapi.Employee) ([]adapterTimeCard, error) {

	timeCards := timeCardResponse.Collection
	var adapterTimeCards []adapterTimeCard

	for _, timeCard := range timeCards {
		for _, timeEntry := range timeCard.TimeEntries {
			employee, exists := employeeMap[timeCard.Worker.UUID]
			if !exists {
				fmt.Printf("Employee with uuid %s not found in employee map\n", timeCard.Worker.UUID)
			}
			project, exists := projectMap[timeCard.Project.UUID]
			if !exists {
				fmt.Printf("Project with uuid %s not found in project map\n", timeCard.Project.UUID)
			}

			adapterTimeCards = append(adapterTimeCards,
				adapterTimeCard{
					EmployeeCode:        employee.EmployeeID,
					EmployeeName:        fmt.Sprintf("%s %s", employee.FirstName, employee.LastName),
					Date:                timeCard.Date,
					Class:               timeEntry.Classification.Name,
					JobNumber:           project.Number,
					CostCode:            timeEntry.CostCode.Code,
					PayType:             timeEntry.PayType.Code,
					Hours:               timeEntry.Hours,
					CostCodeDescription: timeEntry.CostCode.Division,
				})
		}
	}
	return adapterTimeCards, nil
}

func normalizeEquipLogResponse(
	equipLogResponse rakenapi.EquipmentLogResponse,
	projectMap map[string]rakenapi.Project,
	employeeMap map[string]rakenapi.Employee) ([]adapterEquipLog, error) {

	equipAssignments := equipLogResponse.Collection
	var adapterEquipLogs []adapterEquipLog
	for _, assignment := range equipAssignments {
		equipment := assignment.Equipment
		projectUuid := assignment.ProjectUUID

		for _, log := range assignment.Logs {
			employeeName := fmt.Sprintf("%s %s", employeeMap[log.EmployeeID].FirstName, employeeMap[log.EmployeeID].LastName)
			adapterEquipLogs = append(adapterEquipLogs,
				adapterEquipLog{
					EmployeeName: employeeName,
					Date:         log.Date,
					JobNumber:    projectMap[projectUuid].Number,
					EquipNumber:  equipment.Code,
					CostCode:     log.CostCode.Code,
				})
		}
	}
	return adapterEquipLogs, nil
}
