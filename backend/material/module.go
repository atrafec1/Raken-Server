package material

import (
	"fmt"
	"prg_tools/material/adapter/raken"
)

func NewTestProgressEstimateService() *ProgressEstimateService {
	materialSource, err := raken.NewAdapter()
	if err != nil {
		panic(fmt.Sprintf("failed to create raken adapter: %v", err))
	}
	ProgressService := NewProgressEstimateService(materialSource, nil)
	return ProgressService
}
