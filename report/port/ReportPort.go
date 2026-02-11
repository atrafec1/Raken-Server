package port

import "daily_check_in/report/domain"

type ReportFetcher interface {
	GetReports(fromDate, toDate string) ([]domain.ReportCollection, error)
}
