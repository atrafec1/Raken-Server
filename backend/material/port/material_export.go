package port

import "prg_tools/material/domain"

type MaterialExporter interface {
	ExportMaterialLogs([]domain.MaterialLogCollection) error
}
