package domain

type JobMaterialInfo struct {
	Logs      []MaterialLog
	Materials []Material
	Job       Job
	FromDate  string
	ToDate    string
}

type MaterialLog struct {
	Job      Job
	Date     string
	Quantity float64
	Material Material
}

type Material struct {
	BidNumber string
	Name      string
	Unit      string
}
type Job struct {
	Name   string
	Number string
}
