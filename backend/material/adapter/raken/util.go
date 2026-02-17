package raken

import (
	"prg_tools/external/rakenapi"
	"prg_tools/material/domain"
	"strings"
)

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
