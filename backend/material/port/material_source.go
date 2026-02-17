package port

import "prg_tools/material/domain"

type MaterialSource interface {
	GetMaterialLogs(fromDate, toDate string) ([]domain.MaterialLogCollection, error)
}
