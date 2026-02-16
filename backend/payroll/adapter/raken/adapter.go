package raken

import (
	"fmt"
	"prg_tools/external/rakenapi"
	"prg_tools/payroll/dto"
	"strings"
)

type RakenAPIAdapter struct {
	Client rakenapi.RakenClient
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

func (r *RakenAPIAdapter) GetPayrollEntries(fromDate, toDate string) (dto.PayrollEntryResult, error) {

	timeCardResponse, err := r.Client.GetTimeCards(fromDate, toDate)
	if err != nil {
		return dto.PayrollEntryResult{}, err
	}
	equipLogResponse, err := r.Client.GetEquipmentLogs(fromDate, toDate)
	if err != nil {
		return dto.PayrollEntryResult{}, err
	}

	projectMap, err := r.makeProjectMap()
	if err != nil {
		return dto.PayrollEntryResult{}, err
	}
	employeeMap, err := r.makeEmployeeMap()
	if err != nil {
		return dto.PayrollEntryResult{}, err
	}
	adapterTimeCards, err := r.normalizeTimeCardResponse(*timeCardResponse, projectMap, employeeMap)
	if err != nil {
		return dto.PayrollEntryResult{}, err
	}
	fmt.Println("ADAPTER TIME CARDS:", adapterTimeCards)
	adapterEquipLogs, err := r.normalizeEquipLogResponse(*equipLogResponse, projectMap, employeeMap)
	if err != nil {
		return dto.PayrollEntryResult{}, err
	}
	fmt.Println("ADAPTER EQUIP LOGS:", adapterEquipLogs)
	mergedLogs, err := r.mergeTimeAndEquipLogs(adapterTimeCards, adapterEquipLogs)
	if err != nil {
		return dto.PayrollEntryResult{}, fmt.Errorf("failed to merge time cards and equip logs: %w", err)
	}

	fmt.Printf("MERGED LOGS: %+v\n", mergedLogs)
	r.applyPayrollRules(mergedLogs)
	warnings := collectWarnings(adapterTimeCards, adapterEquipLogs)
	return dto.PayrollEntryResult{
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

func (r *RakenAPIAdapter) makeClassMap() (map[string]rakenapi.Class, error) {
	classMap := make(map[string]rakenapi.Class)
	classResp, err := r.Client.GetClasses()
	if err != nil {
		return nil, fmt.Errorf("failed to get classifications from raken api: %w", err)
	}
	for _, class := range classResp.Collection {
		classMap[class.UUID] = class
	}
	return classMap, nil
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
func (r *RakenAPIAdapter) applyPayrollRules(entries []*dto.PayrollEntry) {
	for _, entry := range entries {
		r.applySpecialPhaseRules(entry)
		r.applyCAOvertimeRules(entry)
	}
}

func (r *RakenAPIAdapter) applySpecialPhaseRules(entry *dto.PayrollEntry) {
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

func (r *RakenAPIAdapter) applyCAOvertimeRules(entry *dto.PayrollEntry) {
	switch entry.Day {
	case 6:
		threshold := 12.0
		overTimeHours := entry.RegularHours + entry.OvertimeHours
		if overTimeHours > threshold {
			overAmount := overTimeHours - threshold
			entry.PremiumHours += overAmount
			entry.OvertimeHours = threshold
		} else {
			entry.OvertimeHours = overTimeHours
			entry.RegularHours = 0
		}
	case 7:
		totalHours := entry.RegularHours + entry.OvertimeHours + entry.PremiumHours
		entry.PremiumHours = totalHours
		entry.RegularHours = 0
		entry.OvertimeHours = 0
	}
}

func (r *RakenAPIAdapter) normalizeTimeCardResponse(
	timeCardResponse rakenapi.TimeCardResponse,
	projectMap map[string]rakenapi.Project,
	employeeMap map[string]rakenapi.Employee,
) ([]adapterTimeCard, error) {
	fmt.Println("TIMECARD RESPONSE:", timeCardResponse)

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
			if !exists {
				fmt.Printf("Classification with uuid %s not found in class map\n", timeEntry.Classification.Name)
			}

			adapterTimeCards = append(adapterTimeCards,
				adapterTimeCard{
					EmployeeCode:        employee.EmployeeID,
					EmployeeName:        strings.TrimSpace(fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)),
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

func (r *RakenAPIAdapter) normalizeEquipLogResponse(
	equipLogResponse rakenapi.EquipmentLogResponse,
	projectMap map[string]rakenapi.Project,
	employeeMap map[string]rakenapi.Employee) ([]adapterEquipLog, error) {

	equipAssignments := equipLogResponse.Collection
	var adapterEquipLogs []adapterEquipLog
	for _, assignment := range equipAssignments {
		equipment := assignment.Equipment
		projectUuid := assignment.ProjectUUID

		for _, log := range assignment.Logs {
			employeeName :=
				strings.TrimSpace(
					fmt.Sprintf(
						"%s %s", employeeMap[log.EmployeeID].FirstName, employeeMap[log.EmployeeID].LastName))
			adapterEquipLogs = append(adapterEquipLogs,
				adapterEquipLog{
					EmployeeName: employeeName,
					Date:         log.Date,
					JobNumber:    projectMap[projectUuid].Number,
					EquipNumber:  equipment.Code,
					CostCode:     log.CostCode.Code,
					Hours:        log.Hours,
				})
		}
	}
	return adapterEquipLogs, nil
}
