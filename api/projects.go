package api

import (
	"fmt"
	"net/http"
)

type ProjectResponse struct {
	Collection []Project `json:"collection"`
}

type Project struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Number string `json:"number"`
}

func (c *Client) GetProjects() (*ProjectResponse, error) {
	requestURL := c.config.BaseURL + "projects"
	limit := "1000"
	req, err := http.NewRequest("GET",requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making project request: %v", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("limit",limit)
	req.URL.RawQuery = queryParams.Encode()	
	var projects ProjectResponse
	if err := c.doRequest(req, &projects); err != nil {
		return nil, fmt.Errorf("error retrieving projects %v", err)
	}
	return &projects, nil
}

func (c *Client) UpdateProjectMap() error {
	projectsResp, err := c.GetProjects()
	if err != nil {
		return fmt.Errorf("error getting projects: %v", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.projectMap = make(map[string]Project)
	for _, project := range projectsResp.Collection {
		c.projectMap[project.UUID] = project
	}
	return nil
}

func (c *Client) GetProjectByUUID(uuid string) (Project, error) {
	project, exists := c.projectMap[uuid]
	if !exists {
		if err := c.UpdateProjectMap(); err != nil {
			return Project{}, fmt.Errorf("failed to refresh project map: %w", err)
		}
		project, exists = c.projectMap[uuid]
		if !exists {
			return Project{}, fmt.Errorf("project with UUID %s not found after refresh", uuid)
		}
	}
	return project, nil
}