package raken

import (
	"fmt"
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
	"time"
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

func (a *Adapter) GetMaterialLogs(fromDate, toDate string) ([]domain.MaterialLogCollection, error) {
	projectUuids, err := a.fetchProjectsWorkedOn(fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects worked on: %w", err)
	}

	var materialLogCollections []domain.MaterialLogCollection

	for uuid := range projectUuids {
		materialLogResponse, err := a.client.GetMaterialLogs(uuid, fromDate, toDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get material logs for project %s: %w", uuid, err)
		}

		if len(materialLogResponse.Collection) == 0 {
			continue
		}
		collection := a.toDomainMaterialLogCollection(*materialLogResponse, fromDate, toDate, uuid)
		materialLogCollections = append(materialLogCollections, collection)
		time.Sleep(500 * time.Millisecond)
	}
	return materialLogCollections, nil

}

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

func (a *Adapter) toDomainMaterialLogCollection(
	materialLogResponse rakenapi.MaterialLogResponse,
	fromDate, toDate,
	projectUuid string) domain.MaterialLogCollection {

	var materialLogs []domain.MaterialLog
	job := domain.Job{
		Name:   a.projectMap[projectUuid].Name,
		Number: a.projectMap[projectUuid].Number,
	}

	for _, log := range materialLogResponse.Collection {
		material := domain.Material{
			Name: log.Material.Name,
			Unit: log.Material.Unit.Name,
		}
		materialLogs = append(materialLogs, domain.MaterialLog{
			Job:      job,
			Date:     log.Date,
			Material: material,
			Quantity: log.Quantity,
		})
	}

	return domain.MaterialLogCollection{
		Logs:     materialLogs,
		FromDate: fromDate,
		ToDate:   toDate,
		Job:      job,
	}
}
func (a *Adapter) makeProjectMap() error {
	projects, err := a.client.GetProjects()
	if err != nil {
		return fmt.Errorf("error getting projects: %w", err)
	}

	for _, proj := range projects.Collection {
		a.projectMap[proj.UUID] = proj
	}
	return nil
}
