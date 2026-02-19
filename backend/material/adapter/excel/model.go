package excel

type ProgressRow struct {
	Date       string
	Quantities map[string]float64 // BidNumber -> Quantity
}

type ProgressSheet struct {
	BidItems  []BidItem
	SheetName string
	JobDetail string
	Rows      []ProgressRow
	FromDate  string
	ToDate    string
}

type BidItem struct {
	Number        string
	Name          string
	UnitOfMeasure string
}
