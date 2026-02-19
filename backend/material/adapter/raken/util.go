package raken

import (
	"fmt"
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
	"strings"
)

// Converts raken material log to a domain material log
func toDomainMaterialLog(log rakenapi.MaterialLog) domain.MaterialLog {
	material := domain.Material{
		BidNumber: parseBidItemNumber(log.Material.Name),
		Name:      parseMaterialName(log.Material.Name),
		Unit:      log.Material.Unit.Name,
	}

	return domain.MaterialLog{
		Date:     log.Date,
		Quantity: log.Quantity,
		Material: material,
	}
}

// parses raken material name to get the bid item number, assuming format is "## - materialName"
func parseBidItemNumber(materialName string) string {
	parts := strings.SplitN(materialName, "-", 2)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

// parses material name from raken materials in ##-materialName format
func parseMaterialName(rakenMaterialName string) string {
	// SplitN with 2 ensures that if the material name itself contains a dash,
	// it doesn't get cut off.
	parts := strings.SplitN(rakenMaterialName, "-", 2)
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

	if a.projectMap == nil {
		a.projectMap = make(map[string]rakenapi.Project)
	}

	for _, proj := range projects.Collection {
		a.projectMap[proj.UUID] = proj
	}
	return nil
}

// converts raken material log response to domain material log list
func toDomainMaterialLogs(materialLogResponse *rakenapi.MaterialLogResponse) []domain.MaterialLog {
	var materialLogs []domain.MaterialLog
	for _, log := range materialLogResponse.Collection {
		materialLogs = append(materialLogs, toDomainMaterialLog(log))
	}
	return materialLogs
}

func (a *Adapter) toDomainMaterials(materials []rakenapi.Material) []domain.Material {
	var domainMats []domain.Material
	for _, mat := range materials {
		domainMats = append(domainMats, domain.Material{
			BidNumber: parseBidItemNumber(mat.Name),
			Name:      parseMaterialName(mat.Name),
			Unit:      mat.Unit.Name,
		})
	}
	return domainMats
}

func (a *Adapter) toJobMaterialInfo(
	materialLogResponse *rakenapi.MaterialLogResponse,
	jobMats *rakenapi.MaterialResponse,
	fromDate, toDate,
	projectUuid string) domain.JobMaterialInfo {

	domainMatLogs := toDomainMaterialLogs(materialLogResponse)

	proj, ok := a.projectMap[projectUuid]
	jobName := "Unknown Project"
	jobNumber := ""
	if ok {
		jobName = proj.Name
		jobNumber = proj.Number
	}

	job := domain.Job{
		Name:   jobName,
		Number: jobNumber,
	}

	jobMatInfo := domain.JobMaterialInfo{
		Job:       job,
		FromDate:  fromDate,
		ToDate:    toDate,
		Logs:      domainMatLogs,
		Materials: a.toDomainMaterials(jobMats.Collection),
	}

	return jobMatInfo
}
