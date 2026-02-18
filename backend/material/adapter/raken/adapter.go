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


