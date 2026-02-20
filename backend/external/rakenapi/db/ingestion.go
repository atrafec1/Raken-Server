package db

import "prg_tools/external/rakenapi"

func (p *Pipeline) loadData() {

}

func (p *Pipeline) fetchProjects() (*rakenapi.ProjectResponse, error) {
	return p.Client.GetProjects()
}

func (p *Pipeline) fetchEmployees() (*rakenapi.EmployeeResponse, error) {
	return p.Client.GetEmployees()
}

func (p *Pipeline) fetchClasses() (*rakenapi.ClassResponse, error) {
	return p.Client.GetClasses()
}

func (p *Pipeline) fetchCostCodes() (*rakenapi.CostCodeResponse, error) {
	return p.Client.GetCostCodes()
}
