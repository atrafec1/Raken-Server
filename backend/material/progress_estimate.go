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

func (s *ProgressEstimateService) GetMaterialLogs(fromDate, toDate string) ([]domain.JobMaterialInfo, error) {
	return s.materialSource.GetJobMaterialInfo(fromDate, toDate)
}

func (s *ProgressEstimateService) ExportMaterialLogs(logs []domain.JobMaterialInfo) error {
	return s.materialExporter.ExportJobMaterialInfo(logs)
}
