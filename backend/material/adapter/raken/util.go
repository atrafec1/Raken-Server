package raken

import (
	"fmt"
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
	"regexp"
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

func parseBidItemNumber(materialName string) string {
	re := regexp.MustCompile(`(?i)^\s*(?:item\s*)?(\d+)\s*-?\s*.*$`)
	matches := re.FindStringSubmatch(materialName)
	if len(matches) > 1 {
		return matches[1]
	}
	fmt.Println("No bid item number found in material name: ", materialName)
	return ""
}

func parseMaterialName(rakenMaterialName string) string {
	re := regexp.MustCompile(`(?i)^\s*(?:item\s*)?\d+\s*-?\s*(.*)$`)
	matches := re.FindStringSubmatch(rakenMaterialName)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
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
