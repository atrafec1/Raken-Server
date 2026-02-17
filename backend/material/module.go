package material

import (
	"fmt"
	"prg_tools/material/adapter/excel"
	"prg_tools/material/adapter/raken"
)

func NewTestProgressEstimateService() *ProgressEstimateService {
	materialSource, err := raken.NewAdapter()
	if err != nil {
		panic(fmt.Sprintf("failed to create raken adapter: %v", err))
	}

	materialExporter := excel.NewAdapter("./test_output/material")
	ProgressService := NewProgressEstimateService(materialSource, materialExporter)
	return ProgressService
}
