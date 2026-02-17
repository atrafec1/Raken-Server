package domain

type MaterialLogCollection struct {
	Logs     []MaterialLog
	Job      Job
	FromDate string
	ToDate   string
}

type MaterialLog struct {
	Job      Job
	Date     string
	Quantity float64
	Material Material
}

type Material struct {
	Name string
	Unit string
}
type Job struct {
	Name   string
	Number string
}
