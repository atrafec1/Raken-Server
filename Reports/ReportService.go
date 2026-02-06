package reports

import (
	"daily_check_in/Reports/domain"
	"daily_check_in/api"
	"fmt"
)

type ReportService struct {
	Client *api.Client
}

func (r *ReportService) getProjectMap() (map[string]domain.Project, error) {
	projectMap := make(map[string]domain.Project)
	projectsResp, err := r.Client.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("error getting projects: %v", err)
	}
	fmt.Printf("%+v\n", projectsResp)
	for _, project := range projectsResp.Collection {
		domainProject := domain.Project{
			Name:   project.Name,
			Number: project.Number,
			UUUID:  project.UUID,
		}
		projectMap[project.UUID] = domainProject
	}
	return projectMap, nil
}

func (r *ReportService) GetProjectsWorkedOn(startDate, endDate string) ([]domain.Project, error) {
	var Projects []domain.Project
	seenProjects := make(map[string]struct{})
	projectMap, err := r.getProjectMap()
	if err != nil {
		return nil, fmt.Errorf("error listing all projects: %w", err)
	}

	timeCards, err := r.Client.GetTimecards(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting timecards: %v", err)
	}

	for _, timeCard := range timeCards.Collection {
		if project, exists := projectMap[timeCard.Project.UUID]; exists {
			if _, seen := seenProjects[project.UUUID]; !seen {
				Projects = append(Projects, project)
				seenProjects[project.UUUID] = struct{}{}
			}
		} else {
			fmt.Printf("Project UUID %s not found in project map\n", timeCard.Project.UUID)
		}
	}
	return Projects, nil
}

func (r *ReportService) GetReports(fromDate, toDate string, projects []domain.Project) ([]domain.ReportCollection, error) {
	var allReports []domain.ReportCollection

	for _, project := range projects {
		var projectReportCollection domain.ReportCollection
		dailyReports, err := r.Client.GetDailyReports(project.UUUID, fromDate, toDate)
		if err != nil {
			return nil, fmt.Errorf("error getting reports for project %s: %v", project.Name, err)
		}
		for _, report := range dailyReports.Collection {
			domainReport := domain.DailyReport{
				Creator: report.SignedBy.Name,
				Date:    report.ReportDate,
				Link:    report.ReportLinks.Link,
			}
			projectReportCollection.DailyReports = append(projectReportCollection.DailyReports, domainReport)
		}
		allReports = append(allReports, projectReportCollection)
	}
	return allReports, nil
}
