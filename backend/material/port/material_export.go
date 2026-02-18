package port

import "prg_tools/material/domain"

type MaterialExporter interface {
	ExportJobMaterialInfo([]domain.JobMaterialInfo) error
}
