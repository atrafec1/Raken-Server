package api

import (
	"fmt"
	"net/http"
)

type EmployeeResponse struct {
	Collection []Employee `json:"collection"`
}

type Employee struct {
	UUID string `json:"uuid"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	EmployeeID string `json:"employeeId"`
}

func (c *Client) GetEmployees() (*EmployeeResponse, error) {
	limit := "1000"
	requestURL := c.config.BaseURL + "members"

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making employees request: %v", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("limit",limit)
	req.URL.RawQuery = queryParams.Encode()

	var employees EmployeeResponse
	if err := c.doRequest(req, &employees); err != nil {
		return nil, fmt.Errorf("error retrieving all employees: %v", err)
	}
	return &employees, nil
}

func (c *Client) UpdateEmployeeMap() error {
	employeesResp, err := c.GetEmployees()
	if err != nil {
		return fmt.Errorf("error getting employees: %v", err)
	}	
	c.mu.Lock()
	defer c.mu.Unlock()
	c.employeeMap = make(map[string]Employee)
	for _, employee := range employeesResp.Collection {
		c.employeeMap[employee.UUID] = employee
	}
	return nil
}