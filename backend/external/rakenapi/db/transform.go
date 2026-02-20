package db

import (
	"prg_tools/database"
	"prg_tools/external/rakenapi"
)

func transformProjects(projectResp *rakenapi.ProjectResponse) []database.Job {
	var jobs []database.Job
	for _, project := range projectResp.Collection {
		job := database.Job{
			ExternalId: project.UUID,
			Name:       project.Name,
			Number:     project.Number,
		}
		jobs = append(jobs, job)
	}
	return jobs
}

func transformEmployees(empResp *rakenapi.EmployeeResponse) []database.Employee {
	var employees []database.Employee
	for _, emp := range empResp.Collection {
		employee := database.Employee{
			ExternalId: emp.UUID,
			EmployeeID: emp.EmployeeID,
			FirstName:  emp.FirstName,
			LastName:   emp.LastName,
		}
		employees = append(employees, employee)
	}
	return employees
}

func transformCostCodes(ccResp *rakenapi.CostCodeResponse) []database.CostCode {
	var costCodes []database.CostCode
	for _, cc := range ccResp.Collection {
		costCode := database.CostCode{
			ExternalId:  cc.UUID,
			Code:        cc.Code,
			Description: cc.Description,
		}
		costCodes = append(costCodes, costCode)
	}
	return costCodes
}

func transformClasses(classResp *rakenapi.ClassResponse) []database.Classification {
	var classes []database.Classification
	for _, class := range classResp.Collection {
		classCode := database.Classification{
			ExternalId: class.UUID,
			Name:       class.Name,
		}
		classes = append(classes, classCode)
	}
	return classes
}
