package excel

import (
	"daily_check_in/api"
	"fmt"
)


type CrewAllocationView struct {
	ProjectNumber string
	ProjectName string
	Date string
	Crew []CrewMember
}

type CrewMember struct {
	FirstName string
	LastName string
	Class string
}

func ToCrewAllocationView(resp *api.ToolboxTalkResponse) [] []CrewAllocationEntry  {
	var crews []CrewAllocationView
	for _, entry := range resp.Collection {
		var crew CrewAllocationView
		crew.Date = entry.Date
		project, err := c.GetProjectByUUID(entry.Project.UUID)
		if err != nil {
			fmt.Printf("Error retrieving project with UUID %s: %v\n", entry.Project.UUID, err)
		}
		for _, attendee := range entry.Attendees {
			employee, err := c.GetEmployeeByUUID(attendee.Employee.UUID)
			if err != nil {
				fmt.Printf("Error retrieving employee with UUID %s: %v\n", entry.Employee.UUID, err)
			}
			class, err := c.GetClassByUUID(employee.ClassUUID)
			if err != nil {
				fmt.Printf("Error retrieving class with UUID %s: %v\n", employee.ClassUUID, err)
			}
			
			var crewMember CrewMember
			crewMember.FirstName = employee.FirstName
			crewMember.LastName = employee.LastName
			crewMember.Class = class.Name
			crew.Crew = append(crew.Crew, crewMember)
			
		}
	}
	return crews
}