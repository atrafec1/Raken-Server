package excel

import (
	"fmt"
	"daily_check_in/api"
)

type FormattedToolbox struct {
	Date        string
	ProjectNumber string
	EmployeeNames []string
}
func FormatToolBoxTalkData(client *api.Client) FormattedToolbox {
	resp, err := client.GetToolboxTalks()
	if err != nil {
		fmt.Printf("Error fetching toolbox talks: %v\n", err)
	}
	var formattedData FormattedToolbox

	for _, talk := range resp.Collection {
		formattedData.Date = talk.Date
		formattedData.ProjectNumber = talk.Project.Number
		for _, attendee := range talk.Attendees {
			employee, exists, err := client.GetEmployeeByUUID(attendee.Employee.UUID)
			if !exists {
				fmt.Printf("Employee with UUID %s not found\n", attendee.Employee.UUID)
				continue
			}else if err != nil {
				fmt.Printf("Error retrieving employee with UUID %s: %v\n", attendee.Employee.UUID, err)
			}
			fullName := fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)
			formattedData.EmployeeNames = append(formattedData.EmployeeNames, fullName)
		}
	}
	return formattedData
}