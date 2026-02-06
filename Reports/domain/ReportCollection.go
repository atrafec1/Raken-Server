package domain

type ReportCollection struct {
	DailyReports []DailyReport
	FromDate     string
	ToDate       string
	Project      Project
}

type DailyReport struct {
	Creator string
	Date    string
	Link    string
}

type Project struct {
	Name   string
	Number string
}
