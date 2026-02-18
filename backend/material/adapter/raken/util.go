package raken

import (
	"fmt"
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
	"strings"
)

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
		materialLogs = append(materialLogs, newDomainMaterialLog(log))
	}

	return domain.MaterialLogCollection{
		Logs:     materialLogs,
		FromDate: fromDate,
		ToDate:   toDate,
		Job:      job,
	}
}

func newDomainMaterialLog(log rakenapi.MaterialLog) domain.MaterialLog {
	material := domain.Material{
		BidNumber: getBidItemNumber(log.Material.Name),
		Name:      getMaterialName(log.Material.Name),
		Unit:      log.Material.Unit.Name,
	}

	return domain.MaterialLog{
		Date:     log.Date,
		Quantity: log.Quantity,
		Material: material,
	}

}

func getBidItemNumber(materialName string) string {
	parts := strings.Split(materialName, "-")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

// For bid items ##-materialName
func getMaterialName(rakenMaterialName string) string {
	parts := strings.Split(rakenMaterialName, "-")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(rakenMaterialName)
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
