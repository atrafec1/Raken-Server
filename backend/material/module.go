package material

import (
	"prg_tools/material/adapter/excel"
	"prg_tools/material/adapter/raken"
	"prg_tools/material/adapter/test_material_source"
)

func NewTestProgressEstimateService() *ProgressEstimateService {
	testMaterialSource := test_material_source.NewAdapter()
	materialExporter := excel.NewAdapter("./test_output/material")
	ProgressService := NewProgressEstimateService(testMaterialSource, materialExporter)
	return ProgressService
}

func RakenProgressEstimateService(estimateProgDir string) (*ProgressEstimateService, error) {
	rakenAdapter, err := raken.NewAdapter()
	if err != nil {
		return nil, err
	}
	materialExporter := excel.NewAdapter(estimateProgDir)
	ProgressService := NewProgressEstimateService(rakenAdapter, materialExporter)
	return ProgressService, nil
}
