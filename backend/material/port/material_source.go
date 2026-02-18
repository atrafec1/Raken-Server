package port

import "prg_tools/material/domain"

type MaterialSource interface {
	GetJobMaterialInfo(fromDate, toDate string) ([]domain.JobMaterialInfo, error)
}
