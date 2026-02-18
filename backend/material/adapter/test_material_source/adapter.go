package test_material_source

import (
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
	"strings"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) GetMaterialLogs(fromDate, toDate string) ([]domain.MaterialLogCollection, error) {
	var materialLogCollections []domain.MaterialLogCollection
	rakenMaterialCollection := generateMockCollections(2, fromDate, toDate)
	for _, materialResp := range rakenMaterialCollection {
		collection := toDomainMaterialCollection(materialResp, fromDate, toDate)
		materialLogCollections = append(materialLogCollections, collection)
	}
	return materialLogCollections, nil
}

func toDomainMaterialCollection(materialResp rakenapi.MaterialLogResponse, from, to string) domain.MaterialLogCollection {
	var materialLogs []domain.MaterialLog
	job := generateMockJob()
	for _, log := range materialResp.Collection {
		materialLogs = append(materialLogs, newDomainMaterialLog(log))
	}
	return domain.MaterialLogCollection{
		Job: domain.Job{
			Name:   job.Name,
			Number: job.Number,
		},
		FromDate: from,
		ToDate:   to,
		Logs:     materialLogs,
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
	return parts[1]
}

func getMaterialName(rakenMaterialName string) string {
	parts := strings.Split(rakenMaterialName, "-")
	return parts[2]
}
