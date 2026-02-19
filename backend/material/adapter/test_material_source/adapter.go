package test_material_source

import (
	"fmt"
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
	"regexp"
	"strings"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) GetJobMaterialInfo(fromDate, toDate string) ([]domain.JobMaterialInfo, error) {
	var materialLogCollections []domain.JobMaterialInfo
	rakenMaterialCollection := generateMockCollections(2, fromDate, toDate)
	for _, materialResp := range rakenMaterialCollection {
		collection := toDomainMaterialCollection(materialResp, fromDate, toDate)
		fmt.Printf("%+v\n", collection)
		materialLogCollections = append(materialLogCollections, collection)
	}
	return materialLogCollections, nil
}

func toDomainMaterialCollection(materialResp rakenapi.MaterialLogResponse, from, to string) domain.JobMaterialInfo {
	var materialLogs []domain.MaterialLog
	materialMap := make(map[string]domain.Material) // key: bidNumber+name+unit
	job := generateMockJob()
	for _, log := range materialResp.Collection {
		dml := newDomainMaterialLog(log)
		materialLogs = append(materialLogs, dml)
		// Use a composite key to ensure uniqueness
		key := dml.Material.BidNumber + "|" + dml.Material.Name + "|" + dml.Material.Unit
		materialMap[key] = dml.Material
	}
	var materials []domain.Material
	for _, m := range materialMap {
		materials = append(materials, m)
	}
	return domain.JobMaterialInfo{
		Job: domain.Job{
			Name:   job.Name,
			Number: job.Number,
		},
		FromDate:  from,
		ToDate:    to,
		Logs:      materialLogs,
		Materials: materials,
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
	// Handles: "Item ### - material name", "### - material name", "Item ###-material name", "###-material name"
	re := regexp.MustCompile(`(?i)^\s*(?:item\s*)?(\d+)\s*-?\s*.*$`)
	matches := re.FindStringSubmatch(materialName)
	if len(matches) > 1 {
		return matches[1]
	}
	fmt.Println("No bid item number found in material name: ", materialName)
	return ""
}

func getMaterialName(rakenMaterialName string) string {
	// Handles: "Item ### - material name", "### - material name", "Item ###-material name", "###-material name"
	re := regexp.MustCompile(`(?i)^\s*(?:item\s*)?\d+\s*-?\s*(.*)$`)
	matches := re.FindStringSubmatch(rakenMaterialName)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return strings.TrimSpace(rakenMaterialName)
}
