package api

import (
	"fmt"
	"net/http"
	"time"
)

type ToolboxTalkResponse struct {
	Collection []ToolboxTalkEntry `json:"collection"`
}

type ToolboxTalkEntry struct {
	Project   Project    `json:"project"`
	Attendees []Attendee `json:"attendees"`
	Date      string     `json:"scheduleDate"`
	Status    string     `json:"status"`
}

type Attendee struct {
	Member Member `json:"member"`
}

type Member struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (c *Client) GetToolboxTalks() (*ToolboxTalkResponse, error) {
	requestURL := c.config.BaseURL + "toolboxTalks/past"
	limit := "1000"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making toolbox talk request: %v", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("limit", limit)
	req.URL.RawQuery = queryParams.Encode()

	var toolboxTalkResp ToolboxTalkResponse

	if err := c.doRequest(req, &toolboxTalkResp); err != nil {
		return nil, fmt.Errorf("error retrieving toolbox talks: %v", err)
	}
	return &toolboxTalkResp, nil
}

type Crew struct {
	Date        string
	Project     Project
	CrewMembers []CrewMember
}

type CrewMember struct {
	FirstName string
	LastName  string
	Class     string
	UUID      string
}

func (c *Client) GetCrewAllocationData() ([]Crew, error) {
	resp, err := c.GetToolboxTalks()
	if err != nil {
		fmt.Printf("Error fetching toolbox talks: %v\n", err)
		return nil, err
	}
	var crews []Crew
	for _, entry := range resp.Collection {
		if entry.Status != "COMPLETED" {
			continue
		}
		var crew Crew
		crew.Date = entry.Date
		project, err := c.GetProjectByUUID(entry.Project.UUID)
		if err != nil {
			fmt.Printf("Error retrieving project with UUID %s: %v\n", entry.Project.UUID, err)
		}
		for _, attendee := range entry.Attendees {
			employee, err := c.GetEmployeeByUUID(attendee.Member.UUID)
			if err != nil {
				fmt.Printf("could not retrieve employee with UUID %s: %v\n", attendee.Member.UUID, err)
			}
			class, err := c.GetClassByUUID(employee.ClassUUID)
			if err != nil {
				fmt.Printf("could not retrieve class with UUID %s: %v\n", employee.ClassUUID, err)
			}
			crewMember := CrewMember{
				FirstName: employee.FirstName,
				LastName:  employee.LastName,
				Class:     class.Name,
				UUID:      employee.UUID,
			}
			crew.CrewMembers = append(crew.CrewMembers, crewMember)
		}
		crew.Project = project
		crews = append(crews, crew)
	}
	return crews, nil
}

func GetTodaysCrewAllocations(crews []Crew) []Crew {
	var todaysCrews []Crew
	today := time.Now().Format("2006-01-02")

	for _, crew := range crews {
		if crew.Date == today {
			todaysCrews = append(todaysCrews, crew)
		}
	}
	return todaysCrews
}

type CrewMemberHistory struct {
	Projects map[string][]string // Date : []ProjectNumbers
}

func GetCrewMemberHistory(employeeUUID string, history []Crew) CrewMemberHistory {
	// IMPORTANT: You must initialize the map with make() or it will panic on assignment
	memberHistory := CrewMemberHistory{
		Projects: make(map[string][]string),
	}

	for _, entry := range history {
		for _, employee := range entry.CrewMembers {
			if employee.UUID == employeeUUID {
				// Map the Project Number to the Date string
				memberHistory.Projects[entry.Date] = append(memberHistory.Projects[entry.Date], entry.Project.Number)
			}
		}
	}
	return memberHistory
}
