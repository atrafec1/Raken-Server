package raken

import (
	"fmt"
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
)

type Adapter struct {
	client     rakenapi.RakenClient
	projectMap map[string]rakenapi.Project
}

func NewAdapter() (*Adapter, error) {
	rakenConfig, err := rakenapi.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load raken config: %w", err)
	}
	rakenClient, err := rakenapi.NewClient(rakenConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to make new raken client: %w", err)
	}

	adapter := Adapter{
		client:     rakenClient,
		projectMap: make(map[string]rakenapi.Project),
	}
	if err := adapter.makeProjectMap(); err != nil {
		return nil, err
	}
	return &adapter, nil
}
func (a *Adapter) GetJobMaterialInfo(fromDate, toDate string) ([]domain.JobMaterialInfo, error) {
	projectsWorkedOn, err := a.fetchProjectsWorkedOn(fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error fetching projects worked on: %w", err)
	}
	var everyJobMatInfo []domain.JobMaterialInfo
	for projUuid := range projectsWorkedOn {
		jobMatInfo, err := a.fetchJobMaterialInfo(projUuid, fromDate, toDate)
		if err != nil {
			return nil, fmt.Errorf("error fetching job material info for project %s: %w", projUuid, err)
		}
		everyJobMatInfo = append(everyJobMatInfo, jobMatInfo)
	}
	return everyJobMatInfo, nil
}

// Get material logs, and materials and aggregate to domain job material info
func (a *Adapter) fetchJobMaterialInfo(projectUuid string, fromDate, toDate string) (domain.JobMaterialInfo, error) {
	rakenMaterialLogResp, err := a.client.GetMaterialLogs(projectUuid, fromDate, toDate)
	if err != nil {
		return domain.JobMaterialInfo{}, fmt.Errorf("error getting material logs for project %s: %w", projectUuid, err)
	}

	MatResp, err := a.fetchMaterialsForProject(projectUuid)
	if err != nil {
		return domain.JobMaterialInfo{}, fmt.Errorf("error getting materials for project %s: %w", projectUuid, err)
	}

	jobMatInfo := a.toJobMaterialInfo(rakenMaterialLogResp, MatResp, fromDate, toDate, projectUuid)
	fmt.Printf("Job material info %+v\n", jobMatInfo)
	return jobMatInfo, nil

}

// Get Projects UUIDs from timecards within the date range
func (a *Adapter) fetchProjectsWorkedOn(fromDate, toDate string) (map[string]struct{}, error) {
	projectsWorkedOn := make(map[string]struct{})
	timecards, err := a.client.GetTimeCards(fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error getting timecards: %w", err)
	}
	for _, timecard := range timecards.Collection {
		projectsWorkedOn[timecard.Project.UUID] = struct{}{}
	}
	return projectsWorkedOn, nil
}

// Get materials set for this specific project
func (a *Adapter) fetchMaterialsForProject(projectUuid string) (*rakenapi.MaterialResponse, error) {
	materials, err := a.client.GetMaterialsForProject(projectUuid)
	if err != nil {
		return nil, fmt.Errorf("error getting materials for project %s: %w", projectUuid, err)
	}
	return materials, nil
}

//Need to find all projects worked on
//For each project get material logs for project
//Convert material logs to domain material log collections
