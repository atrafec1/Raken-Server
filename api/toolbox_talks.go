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
	Project Project `json:"project"`
	Attendees []AttendeeRecord `json:"attendees"`
	Date string `json:"scheduleDate"`
}

type AttendeeRecord struct {
	Employee Employee `json:"member"`
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

type CrewAllocationEntry struct {
	Date string
	Project Project 
	Employees []Employee 
}

func (c *Client) GetCrewAllocationData() ([]CrewAllocationEntry, error) {
	resp, err := c.GetToolboxTalks()
	if err != nil {
		fmt.Printf("Error fetching toolbox talks: %v\n", err)
		return nil, err
	}
	var crewAllocations []CrewAllocationEntry

	for _, entry := range resp.Collection {
		var crew CrewAllocationEntry
		crew.Date = entry.Date
		project, err := c.GetProjectByUUID(entry.Project.UUID)
		if err != nil {
			fmt.Printf("Error retrieving project with UUID %s: %v\n", entry.Project.UUID, err)
		}

		crew.Project = project
		fmt.Printf("Project number: %s\n", entry.Project.Number)
		for _, attendee := range entry.Attendees {
			employee, err := c.GetEmployeeByUUID(attendee.Employee.UUID)
			if err != nil {
				fmt.Printf("Error retrieving employee with UUID %s: %v\n", entry.Employee.UUID, err)
			}
			attendee
			crew.Employees = append(crew.Employees, employee)
		}
		crewAllocations = append(crewAllocations, crew)
	}
	return crewAllocations, nil 
}

func GetTodaysCrewAllocations(crews []CrewAllocationEntry) []CrewAllocationEntry {
	var todaysCrews []CrewAllocationEntry
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

func GetCrewMemberHistory(employeeUUID string, history []CrewAllocationEntry) CrewMemberHistory {
	// IMPORTANT: You must initialize the map with make() or it will panic on assignment
	memberHistory := CrewMemberHistory{
		Projects: make(map[string][]string),
	}

	for _, entry := range history {
		for _, employee := range entry.Employees {
			if employee.UUID == employeeUUID {
				// Map the Project Number to the Date string
				memberHistory.Projects[entry.Date] = append(memberHistory.Projects[entry.Date], entry.Project.Number)
			}
		}
	}
	return memberHistory
}


func GetMockData() []CrewAllocationEntry {
	// --- Projects ---
	p1 := Project{UUID: "5", Number: "222001", Name: "San Dieguito Lagoon"}
	p2 := Project{UUID: "6", Number: "225010", Name: "Beador Route 15"}
	p3 := Project{UUID: "7", Number: "224026", Name: "FSSW I-15 Maint"}
	p4 := Project{UUID: "8", Number: "225036", Name: "RCRCD Dos Lagos"}
	p5 := Project{UUID: "9", Number: "222003", Name: "SD LAB APP"}

	// --- Employees ---
		
	e1 := Employee{UUID: "1a2b3c", FirstName: "Santiago", LastName: "Aguirre", EmployeeID: "SA1001"}
	e2 := Employee{UUID: "2b3c4d", FirstName: "Fernado", LastName: "Correa", EmployeeID: "FC1002"}
	e3 := Employee{UUID: "3c4d5e", FirstName: "Jesus", LastName: "Marin", EmployeeID: "JM1003"}
	e4 := Employee{UUID: "4d5e6f", FirstName: "Salvador", LastName: "Ortiz", EmployeeID: "SO1004"}
	e5 := Employee{UUID: "5e6f7g", FirstName: "Ivan", LastName: "Valdivia", EmployeeID: "IV1005"}
	e6 := Employee{UUID: "6f7g8h", FirstName: "EM Aguinaga", LastName: "Maint", EmployeeID: "EA1006"}
	e7 := Employee{UUID: "7g8h9i", FirstName: "Doug", LastName: "Richards", EmployeeID: "DR1007"}
	e8 := Employee{UUID: "8h9i0j", FirstName: "Cristoval", LastName: "Delgado", EmployeeID: "CD1008"}
	e9 := Employee{UUID: "9i0j1k", FirstName: "EdSan", LastName: "Morteson 2A", EmployeeID: "EM1009"}
	e10 := Employee{UUID: "0j1k2l", FirstName: "Jesus", LastName: "Garcia", EmployeeID: "JG1010"}
	e11 := Employee{UUID: "1k2l3m", FirstName: "Luis", LastName: "Jauregui", EmployeeID: "LJ1011"}
	e12 := Employee{UUID: "2l3m4n", FirstName: "Jose Santos", LastName: "Marin", EmployeeID: "JM1012"}
	e13 := Employee{UUID: "3m4n5o", FirstName: "FSSW I-15", LastName: "Maint", EmployeeID: "FI1013"}
	e14 := Employee{UUID: "4n5o6p", FirstName: "Raul", LastName: "Arreguin", EmployeeID: "RA1014"}
	e15 := Employee{UUID: "5o6p7q", FirstName: "Agustin", LastName: "Avila", EmployeeID: "AA1015"}
	e16 := Employee{UUID: "6p7q8r", FirstName: "Oscar", LastName: "Calixto", EmployeeID: "OC1016"}
	e17 := Employee{UUID: "7q8r9s", FirstName: "Bernardino", LastName: "De La Cruz", EmployeeID: "BD1017"}
	e18 := Employee{UUID: "8r9s0t", FirstName: "Luis", LastName: "Espinoza", EmployeeID: "LE1018"}
	e19 := Employee{UUID: "9s0t1u", FirstName: "Martin", LastName: "Garcia", EmployeeID: "MG1019"}
	e20 := Employee{UUID: "0t1u2v", FirstName: "Gerardo", LastName: "Hernandez", EmployeeID: "GH1020"}
	e21 := Employee{UUID: "1u2v3w", FirstName: "Trinidad", LastName: "Lomeli", EmployeeID: "TL1021"}
	e22 := Employee{UUID: "2v3w4x", FirstName: "Salvador", LastName: "Martinez", EmployeeID: "SM1022"}
	e23 := Employee{UUID: "3w4x5y", FirstName: "Gilberto", LastName: "Ortiz", EmployeeID: "GO1023"}
	e24 := Employee{UUID: "4x5y6z", FirstName: "Sergio", LastName: "Palafox", EmployeeID: "SP1024"}
	e25 := Employee{UUID: "5y6z7a", FirstName: "Bernardo", LastName: "Ramirez", EmployeeID: "BR1025"}
	e26 := Employee{UUID: "6z7a8b", FirstName: "Matilde", LastName: "Torres", EmployeeID: "MT1026"}
	e27 := Employee{UUID: "7a8b9c", FirstName: "Paulin", LastName: "Marin", EmployeeID: "PM1027"}
	e28 := Employee{UUID: "8b9c0d", FirstName: "Pablo", LastName: "Marin", EmployeeID: "PM1028"}
	e29 := Employee{UUID: "9c0d1e", FirstName: "DeAnna", LastName: "Jessup", EmployeeID: "DJ1029"}
	e30 := Employee{UUID: "0d1e2f", FirstName: "Edgar", LastName: "Marin", EmployeeID: "EM1030"}
	e31 := Employee{UUID: "1e2f3g", FirstName: "Efrain", LastName: "Oropeza", EmployeeID: "EO1031"}
	e32 := Employee{UUID: "2f3g4h", FirstName: "Alex", LastName: "Regalado", EmployeeID: "AR1032"}
	e33 := Employee{UUID: "3g4h5i", FirstName: "Adam", LastName: "Trafecanty", EmployeeID: "AT1033"}
	e34 := Employee{UUID: "4h5i6j", FirstName: "Georgia", LastName: "Vitous", EmployeeID: "GV1034"}
	e35 := Employee{UUID: "5i6j7k", FirstName: "Guy", LastName: "Gray", EmployeeID: "GG1035"}
	e36 := Employee{UUID: "6j7k8l", FirstName: "Salvador", LastName: "Fernandez", EmployeeID: "SF1036"}
	e37 := Employee{UUID: "7k8l9m", FirstName: "Cesar", LastName: "Ramirez", EmployeeID: "CR1037"}
	e38 := Employee{UUID: "8l9m0n", FirstName: "Guadalupe", LastName: "Marin", EmployeeID: "GM1038"}
	e39 := Employee{UUID: "9m0n1o", FirstName: "Roberto", LastName: "Lastra", EmployeeID: "RL1039"}
	e40 := Employee{UUID: "0n1o2p", FirstName: "Salvador", LastName: "Mora", EmployeeID: "SM1040"}
	e41 := Employee{UUID: "1o2p3q", FirstName: "Ricardo", LastName: "Roblado", EmployeeID: "RR1041"}
	e42 := Employee{UUID: "2p3q4r", FirstName: "Jose A.", LastName: "Flores", EmployeeID: "JF1042"}
	e43 := Employee{UUID: "3q4r5s", FirstName: "Cipriano", LastName: "Flores-Leon", EmployeeID: "CF1043"}
	e44 := Employee{UUID: "4r5s6t", FirstName: "Julian", LastName: "Lopez", EmployeeID: "JL1044"}
	e45 := Employee{UUID: "5s6t7u", FirstName: "Jose", LastName: "Lopez", EmployeeID: "JL1045"}
	e46 := Employee{UUID: "6t7u8v", FirstName: "Fernando", LastName: "Jimenez", EmployeeID: "FJ1046"}
	e47 := Employee{UUID: "7u8v9w", FirstName: "Jesus", LastName: "Velasquez", EmployeeID: "JV1047"}


	// --- Dates (Monday - Friday) ---
	d1 := "2025-12-15"
	d2 := "2025-12-16"
	d3 := "2025-12-17"
	d4 := "2025-12-18"
	d5 := "2025-12-19"

	//Create random crew allocations for each project on each day
	
	// --- Crew Allocation ---
	return []CrewAllocationEntry{
    // --- MONDAY (d1) ---
    {Date: d1, Project: p1, Employees: []Employee{e14, e16, e17, e1, e8, e28, e36, e1, e2, e3, e4, e5, e6}}, 
    {Date: d1, Project: p2, Employees: []Employee{e2, e10, e11, e12, e21, e9, e13, e18, e19, e20}},         
    {Date: d1, Project: p3, Employees: []Employee{e7, e30, e32, e33, e22, e24, e25, e26}},             
    {Date: d1, Project: p4, Employees: []Employee{e15, e23, e46, e47, e27, e29, e31}},             
    {Date: d1, Project: p5, Employees: []Employee{e34, e35, e37, e38, e39, e40, e41, e42, e43, e44, e45}},

    // --- TUESDAY (d2) ---
    {Date: d2, Project: p1, Employees: []Employee{e14, e16, e17, e1, e8, e28, e36, e3, e4, e5, e6}},
    {Date: d2, Project: p2, Employees: []Employee{e2, e10, e11, e12, e21, e9, e13, e18, e19, e20}},
    {Date: d2, Project: p3, Employees: []Employee{e7, e30, e32, e33, e22}},
    {Date: d2, Project: p4, Employees: []Employee{e15, e23, e46, e47, e24, e25, e26, e27, e29, e31}},
    {Date: d2, Project: p5, Employees: []Employee{e34, e35, e37, e38, e39, e40, e41, e42, e43, e44, e45}},

    // --- WEDNESDAY (d3) ---
    {Date: d3, Project: p1, Employees: []Employee{e14, e16, e17, e1, e8, e28, e36, e15, e23, e3, e4, e5, e6}},
    {Date: d3, Project: p2, Employees: []Employee{e2, e10, e11, e12, e21, e9, e13, e18, e19, e20}},
    {Date: d3, Project: p3, Employees: []Employee{e7, e30, e32, e33, e22, e24, e25}},
    {Date: d3, Project: p5, Employees: []Employee{e46, e47, e40, e41, e34, e35, e37, e38, e39, e42, e43, e44, e45}},

    // --- THURSDAY (d4) ---
    {Date: d4, Project: p1, Employees: []Employee{e14, e16, e17, e15, e3, e4, e5, e6, e9}},
    {Date: d4, Project: p2, Employees: []Employee{e2, e10, e15, e13, e18, e19, e20}}, 
    {Date: d4, Project: p3, Employees: []Employee{e7, e30, e32, e33, e22, e24, e25, e26}},
    {Date: d4, Project: p5, Employees: []Employee{e46, e47, e23, e34, e35, e37, e38, e39, e40, e41, e42, e43, e44, e45}},

    // --- FRIDAY (d5 - TODAY) ---
    {Date: d5, Project: p1, Employees: []Employee{e14, e16, e17, e1, e8, e28, e36, e3, e4, e5, e6}},
    {Date: d5, Project: p2, Employees: []Employee{e2, e10, e11, e12, e21, e15, e9, e13, e18, e19, e20}}, 
    {Date: d5, Project: p3, Employees: []Employee{e7, e30, e32, e33, e22}},
    {Date: d5, Project: p4, Employees: []Employee{e23, e46, e47, e42, e43, e24, e25, e26, e27, e29, e31}},
    {Date: d5, Project: p5, Employees: []Employee{e34, e35, e37, e38, e39, e40, e41, e44, e45}},
}
}
