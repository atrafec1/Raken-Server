package adapter

import (
	"daily_check_in/external/rakenapi"
	"daily_check_in/payroll/dto"
	"fmt"
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

func (r *RakenAPIAdapter) GetPayrollEntries(fromDate, toDate string) ([]dto.PayrollEntry, error) {

	timeCardResponse, err := r.Client.GetTimeCards(fromDate, toDate)
	if err != nil {
		return nil, err
	}
	equipLogResponse, err := r.Client.GetEquipmentLogs(fromDate, toDate)
	if err != nil {
		return nil, err
	}

	projectMap, err := r.makeProjectMap()
	if err != nil {
		return nil, err
	}
	employeeMap, err := r.makeEmployeeMap()
	if err != nil {
		return nil, err
	}
	adapterTimeCards, err := normalizeTimeCardResponse(*timeCardResponse, projectMap, employeeMap)
	if err != nil {
		return nil, err
	}
	adapterEquipLogs, err := normalizeEquipLogResponse(*equipLogResponse, projectMap, employeeMap)
	if err != nil {
		return nil, err
	}
	return mergeTimeAndEquipLogs(adapterTimeCards, adapterEquipLogs)
}

func mergeTimeAndEquipLogs(
	timeCards []adapterTimeCard,
	equipLogs []adapterEquipLog,
) ([]dto.PayrollEntry, error) {

	entries := make(map[mergeKey]*dto.PayrollEntry)

	// Timecards → labor hours
	for _, tc := range timeCards {
		key := mergeKey{
			EmployeeName: tc.EmployeeName,
			Date:         tc.Date,
			JobNumber:    tc.JobNumber,
			CostCode:     tc.CostCode,
		}

		entry, exists := entries[key]
		if !exists {
			day, err := convertDateToInt(tc.Date)
			if err != nil {
				return nil, fmt.Errorf("failed to convert date to int: %w", err)
			}
			entry = &dto.PayrollEntry{
				EmployeeCode: tc.EmployeeCode,
				CurrentDate:  tc.Date,
				CraftLevel:   tc.Class,
				JobNumber:    tc.JobNumber,
				CostCode:     tc.CostCode,
				Day:          day,
			}
			entries[key] = entry
		}

		payRoute := routePay(tc)
		entry.RegularHours += payRoute.RegularHours
		entry.OvertimeHours += payRoute.OvertimeHours
		entry.PremiumHours += payRoute.PremiumHours

	}

	// Equipment logs
	for _, el := range equipLogs {
		key := mergeKey{
			EmployeeName: el.EmployeeName,
			Date:         el.Date,
			JobNumber:    el.JobNumber,
			CostCode:     el.CostCode,
		}

		entry, exists := entries[key]
		if !exists {
			day, err := convertDateToInt(el.Date)
			if err != nil {
				return nil, fmt.Errorf("failed to convert date to int: %w", err)
			}
			entry = &dto.PayrollEntry{
				CurrentDate:    el.Date,
				JobNumber:      el.JobNumber,
				CostCode:       el.CostCode,
				SpecialPayType: "EQP",
				SpecialPayCode: el.EquipNumber,
				SpecialUnits:   el.Hours,
				Day:            day,
			}
			entries[key] = entry
		}

		entry.EquipmentCode = el.EquipNumber
	}

	// Convert map → slice
	result := make([]dto.PayrollEntry, 0, len(entries))
	for _, v := range entries {
		result = append(result, *v)
	}

	return result, nil
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

type adapterTimeCard struct {
	EmployeeCode string
	EmployeeName string
	Date         string
	JobNumber    string
	Class        string
	CostCode     string
	PayType      string
	Hours        float64
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

func normalizeTimeCardResponse(
	timeCardResponse rakenapi.TimeCardResponse,
	projectMap map[string]rakenapi.Project,
	employeeMap map[string]rakenapi.Employee) ([]adapterTimeCard, error) {

	timeCards := timeCardResponse.Collection
	var adapterTimeCards []adapterTimeCard

	for _, timeCard := range timeCards {
		for _, timeEntry := range timeCard.TimeEntries {
			employee := employeeMap[timeCard.Worker.UUID]
			project := projectMap[timeCard.Project.UUID]

			adapterTimeCards = append(adapterTimeCards,
				adapterTimeCard{
					EmployeeCode: employee.EmployeeID,
					EmployeeName: fmt.Sprintf("%s %s", employee.FirstName, employee.LastName),
					Date:         timeCard.Date,
					Class:        timeEntry.Classification.Name,
					JobNumber:    project.Number,
					CostCode:     timeEntry.CostCode.Code,
					PayType:      timeEntry.PayType.Code,
					Hours:        timeEntry.Hours,
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
