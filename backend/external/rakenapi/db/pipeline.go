package db

import (
	"fmt"
	"prg_tools/database"
	"prg_tools/external/rakenapi"
)

type Pipeline struct {
	Client *rakenapi.Client
}

func NewPipeline() (*Pipeline, error) {
	cfg, err := rakenapi.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}
	rakenClient, err := rakenapi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Pipeline{
		Client: rakenClient,
	}, nil
}

func (p *Pipeline) LoadDatabase() error {
	data, err := p.gatherData()
	if err != nil {
		return fmt.Errorf("error gathering data: %v", err)
	}
	fmt.Println(data)

	tx := database.DB.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	if err := tx.CreateInBatches(data["jobs"], 100).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting jobs: %v", err)
	}
	if err := tx.CreateInBatches(data["employees"], 100).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting employees: %v", err)
	}
	if err := tx.CreateInBatches(data["costCodes"], 100).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting cost codes: %v", err)
	}
	return tx.Commit().Error
}

func (p *Pipeline) gatherData() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	jobs, err := p.Jobs()
	if err != nil {
		return nil, fmt.Errorf("error fetching jobs: %v", err)
	}
	data["jobs"] = jobs
	employees, err := p.Employees()
	if err != nil {
		return nil, fmt.Errorf("error fetching employees: %v", err)
	}
	data["employees"] = employees
	costCodes, err := p.CostCodes()
	if err != nil {
		return nil, fmt.Errorf("error fetching cost codes: %v", err)
	}
	data["costCodes"] = costCodes
	return data, nil
}

func (p *Pipeline) Jobs() ([]database.Job, error) {
	projectResp, err := p.fetchProjects()
	if err != nil {
		return nil, err
	}

	return transformProjects(projectResp), nil
}

func (p *Pipeline) Employees() ([]database.Employee, error) {
	empResp, err := p.fetchEmployees()
	if err != nil {
		return nil, err
	}
	return transformEmployees(empResp), nil
}

func (p *Pipeline) CostCodes() ([]database.CostCode, error) {
	ccResp, err := p.fetchCostCodes()
	if err != nil {
		return nil, err
	}
	return transformCostCodes(ccResp), nil
}
