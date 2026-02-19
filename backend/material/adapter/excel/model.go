package excel

type ProgressSheetSection struct {
	FromDate string
	ToDate   string
	Rows     []ProgressRow
}

type ProgressRow struct {
	Date       string
	Quantities map[string]float64 // BidNumber -> Quantity
}

type ProgressSheet struct {
	BidItems  []BidItem
	Sections  []ProgressSheetSection
	SheetName string
	JobDetail string
}

type BidItem struct {
	Number        string
	Name          string
	UnitOfMeasure string
}
