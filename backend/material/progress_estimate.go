package material

import (
	"prg_tools/material/domain"
	"prg_tools/material/port"
)

type ProgressEstimateService struct {
	materialSource   port.MaterialSource
	materialExporter port.MaterialExporter
}

func NewProgressEstimateService(materialSource port.MaterialSource, estimateExporter port.MaterialExporter) *ProgressEstimateService {
	return &ProgressEstimateService{
		materialSource:   materialSource,
		materialExporter: estimateExporter,
	}

}

func (s *ProgressEstimateService) GetMaterialLogs(fromDate, toDate string) ([]domain.MaterialLogCollection, error) {
	return s.materialSource.GetMaterialLogs(fromDate, toDate)
}

func (s *ProgressEstimateService) ExportMaterialLogs(logs []domain.MaterialLogCollection) error {
	return s.materialExporter.ExportMaterialLogs(logs)
}
