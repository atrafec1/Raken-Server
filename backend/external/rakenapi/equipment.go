package rakenapi

import (
	"fmt"
	"net/http"
)

type EquipmentLogResponse struct {
	Collection []EquipmentAssignment `json:"collection"`
}

type EquipmentAssignment struct {
	ProjectUUID string         `json:"projectUuid"`
	Equipment   Equipment      `json:"equipment"`
	Logs        []EquipmentLog `json:"logs"`
}

type Equipment struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Code string `json:"equipmentId"`
}

type EquipmentLog struct {
	Date       string   `json:"date"`
	Hours      float64  `json:"hours"`
	EmployeeID string   `json:"operator"`
	Status     string   `json:"status"`
	CostCode   CostCode `json:"costCode"`
}

func (c *Client) GetEquipmentLogs(fromDate, toDate string) (*EquipmentLogResponse, error) {
	limit := "1000"
	requestURL := c.config.BaseURL + "equipmentLogs"

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make new http request for equip log: %w", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("fromDate", fromDate)
	queryParams.Set("toDate", toDate)
	queryParams.Set("limit", limit)
	req.URL.RawQuery = queryParams.Encode()

	var response EquipmentLogResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, fmt.Errorf("failed to complete equip logs request: %w", err)
	}

	return &response, nil
}
