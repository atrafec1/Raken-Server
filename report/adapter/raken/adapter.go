package raken

import (
	"daily_check_in/external/rakenapi"
	"daily_check_in/report/domain"
	"fmt"
	"time"
)

type Adapter struct {
	client     *rakenapi.Client
	projectMap map[string]rakenapi.Project
}

func NewAdapter() (*Adapter, error) {
	rakenConfig, err := rakenapi.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load raken config: %w", err)
	}
	rakenClient, err := rakenapi.NewClient(rakenConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to make new raken client: %w", err)
	}

	adapter := Adapter{
		client:     rakenClient,
		projectMap: make(map[string]rakenapi.Project),
	}
	if err := adapter.makeProjectMap(); err != nil {
		return nil, err
	}
	return &adapter, nil
}

func (a *Adapter) GetReports(fromDate, toDate string) ([]domain.ReportCollection, error) {
	projectsWorkedOn, err := a.fetchProjectsWorkedOn(fromDate, toDate)
	if err != nil {
		return nil, err
	}
	var reportCollections []domain.ReportCollection

	for projectUuid := range projectsWorkedOn {
		projReportsResp, err := a.client.GetDailyReports(projectUuid, fromDate, toDate)
		if err != nil {
			return []domain.ReportCollection{}, err
		}
		if len(projReportsResp.Collection) == 0 {
			continue
		}
		time.Sleep(700 * time.Millisecond)
		reportCollection, err := a.convertToDomainReportCollection(projReportsResp.Collection, fromDate, toDate)
		if err != nil {
			return []domain.ReportCollection{}, err
		}
		reportCollections = append(reportCollections, reportCollection)
	}
	return reportCollections, nil
}

func (a *Adapter) makeProjectMap() error {
	projects, err := a.client.GetProjects()
	if err != nil {
		return fmt.Errorf("error getting projects: %w", err)
	}

	for _, proj := range projects.Collection {
		a.projectMap[proj.UUID] = proj
	}
	return nil
}

func (a *Adapter) fetchProjectsWorkedOn(fromDate, toDate string) (map[string]struct{}, error) {
	projectsWorkedOn := make(map[string]struct{})
	timecards, err := a.client.GetTimeCards(fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error getting timecards: %w", err)
	}
	for _, timecard := range timecards.Collection {
		projectsWorkedOn[timecard.Project.UUID] = struct{}{}
	}
	return projectsWorkedOn, nil
}

func (a *Adapter) convertToDomainReportCollection(apiReports []rakenapi.DailyReport, fromDate, toDate string) (domain.ReportCollection, error) {
	var reportCollection domain.ReportCollection

	reportCollection.FromDate = fromDate
	reportCollection.ToDate = toDate

	for _, apiReport := range apiReports {
		domainReport, err := a.convertToDomainReport(apiReport)
		if err != nil {
			return domain.ReportCollection{}, err
		}
		reportCollection.Reports = append(reportCollection.Reports, domainReport)
	}
	return reportCollection, nil
}

func (a *Adapter) convertToDomainReport(apiReport rakenapi.DailyReport) (domain.Report, error) {
	apiProject, ok := a.projectMap[apiReport.ProjectUuid]
	if !ok {
		return domain.Report{}, fmt.Errorf("project map lacked correct project uuid: %s", apiReport.ProjectUuid)
	}
	return domain.Report{
		Date:    apiReport.ReportDate,
		PDFLink: apiReport.ReportLinks.Link,
		Creator: domain.ReportCreator{
			Name: apiReport.SignedBy.Name,
		},
		Project: domain.Project{
			Name:   apiProject.Name,
			Number: apiProject.Number,
		},
	}, nil
}

//Get all project info
//retrieve all the reports given that date
// convert repoorts into domain reports
//return domain report colelction
