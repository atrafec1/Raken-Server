package port

import "prg_tools/report/domain"

type ReportFetcher interface {
	GetReports(fromDate, toDate string) ([]domain.ReportCollection, error)
}

