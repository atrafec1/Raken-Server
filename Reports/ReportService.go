package reports

import (
	"daily_check_in/Reports/domain"
	"daily_check_in/api"
)

type ReportService struct {
	Client *api.Client
}

func (r *ReportService) GetReports(projectUuid string, startDate string, endDate string) (domain.ReportCollection, error) {
	var api.DailyReportResponse
	dailyReports, err := r.Client.GetDailyReports(projectUuid, startDate, endDate)
	if err != nil {
		return domain.ReportCollection{}, err
	}

}
