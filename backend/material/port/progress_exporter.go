package port

type ProgressExporter interface {
	ExportProgressEstimate(projectName string, jobNumber string, date string, materialName string, quantity float64, unit string, costCode string) error
}
