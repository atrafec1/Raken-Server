package material

import (
	"prg_tools/material/domain"
	"prg_tools/material/port"
)

type ProgressEstimateService struct {
	materialSource   port.MaterialSource
	estimateExporter port.ProgressExporter
}

func NewProgressEstimateService(materialSource port.MaterialSource, estimateExporter port.ProgressExporter) *ProgressEstimateService {
	return &ProgressEstimateService{
		materialSource:   materialSource,
		estimateExporter: estimateExporter,
	}

}

func (s *ProgressEstimateService) GetMaterialLogs(fromDate, toDate string) ([]domain.MaterialLogCollection, error) {
	return s.materialSource.GetMaterialLogs(fromDate, toDate)
}
