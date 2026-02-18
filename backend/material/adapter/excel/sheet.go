package excel

import (
	"fmt"
	"prg_tools/material/domain"
	"sort"
	"time"

	"github.com/xuri/excelize/v2"
)

// --- DTOs ---

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

// --- Helpers ---

func sortBidItems(bidItems []BidItem) {
	sort.Slice(bidItems, func(i, j int) bool {
		return bidItems[i].Number < bidItems[j].Number
	})
}

func colLetter(idx int) string {
	name, _ := excelize.ColumnNumberToName(idx)
	return name
}

// --- Excel Writing ---
func (p *ProgressSheet) createHeader(f *excelize.File) error {
	// Job Name in A1
	if err := f.SetCellValue(p.SheetName, "A1", "Pacific Restoration Group, Inc."); err != nil {
		return err
	}
	if err := f.SetCellValue(p.SheetName, "A2", p.JobDetail); err != nil {
		return err
	}
	return nil
}

func (p *ProgressSheet) createBidHeader(f *excelize.File, startRow int) error {
	sortBidItems(p.BidItems)

	// Write "Date"
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("B%d", startRow+2), "Date"); err != nil {
		return err
	}
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("B%d", startRow), "Bid Item"); err != nil {
		return err
	}
	bidNumberStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})

	for col, bi := range p.BidItems {
		col = col + 1
		letter := colLetter(col + 2)

		if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, startRow), bi.Number); err != nil {
			return err
		}
		if err := f.SetCellStyle(p.SheetName, fmt.Sprintf("%s%d", letter, startRow), fmt.Sprintf("%s%d", letter, startRow), bidNumberStyle); err != nil {
			return err
		}
		if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, startRow+1), bi.Name); err != nil {
			return err
		}
		if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, startRow+2), bi.UnitOfMeasure); err != nil {
			return err
		}
	}

	// --- Create contiguous bottom border ---
	style, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "bottom",
				Color: "000000",
				Style: 2, // medium line
			},
		},
	})
	if err != nil {
		return err
	}

	// Apply border across entire header width
	startCol := colLetter(3) // first bid column (C)
	endCol := colLetter(len(p.BidItems) + 2)

	borderRow := startRow + 2 // bottom row of header

	startCell := fmt.Sprintf("%s%d", startCol, borderRow)
	endCell := fmt.Sprintf("%s%d", endCol, borderRow)

	if err := f.SetCellStyle(p.SheetName, startCell, endCell, style); err != nil {
		return err
	}

	return nil
}

func (p *ProgressSheet) createSection(f *excelize.File, section ProgressSheetSection, rowIndex int) error {
	// 1. Update style to include top AND bottom borders
	sectionHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Family: "Aptos Narrow"},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// 2. Calculate the end column letter based on BidItems
	// Column A (1) + Date Column (1) + len(BidItems)
	lastColIdx := len(p.BidItems) + 2
	endCol, _ := excelize.ColumnNumberToName(lastColIdx)

	from, _ := time.Parse("2006-01-02", section.FromDate)
	to, _ := time.Parse("2006-01-02", section.ToDate)
	sectionDateRange := fmt.Sprintf("%s - %s", from.Format("Jan 02"), to.Format("Jan 02"))

	// 3. Set the value in A
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("A%d", rowIndex), sectionDateRange); err != nil {
		return err
	}

	// 4. Apply the style to the entire row range (A to endCol)
	startCell := fmt.Sprintf("A%d", rowIndex)
	endCell := fmt.Sprintf("%s%d", endCol, rowIndex)

	return f.SetCellStyle(p.SheetName, startCell, endCell, sectionHeaderStyle)
}

func (p *ProgressSheet) createRow(f *excelize.File, row ProgressRow, rowIndex int) error {
	// Date
	if err := f.SetCellValue(p.SheetName, fmt.Sprintf("B%d", rowIndex), row.Date); err != nil {
		return err
	}

	// Quantities
	for col, bi := range p.BidItems {
		letter := colLetter(col + 3)
		if val, ok := row.Quantities[bi.Number]; ok {
			if err := f.SetCellValue(p.SheetName, fmt.Sprintf("%s%d", letter, rowIndex), val); err != nil {
				return err
			}
		}
	}

	return nil
}

// --- Orchestration ---

func (p *ProgressSheet) CreateEstimateProgressSheet(fileName string) error {
	f := excelize.NewFile()
	if p.SheetName == "" {
		p.SheetName = "ProgressSheet"
	}

	index, _ := f.NewSheet(p.SheetName)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	currentRow := 3

	// Header
	if err := p.createHeader(f); err != nil {
		return err
	}
	currentRow += 3

	//Bid Item Header

	if err := p.createBidHeader(f, currentRow); err != nil {
		return err
	}
	currentRow += 3
	// Sections + rows
	for _, sec := range p.Sections {
		if err := p.createSection(f, sec, currentRow); err != nil {
			return err
		}
		currentRow++
		for _, r := range sec.Rows {
			if err := p.createRow(f, r, currentRow); err != nil {
				return err
			}
			currentRow++
		}
		currentRow++
	}

	return f.SaveAs(fileName)
}

// --- Conversion Helpers ---

func ConvertToProgressSheet(mLogs domain.MaterialLogCollection, sectionEndDay int) (ProgressSheet, error) {
	if len(mLogs.Logs) == 0 {
		return ProgressSheet{}, nil
	}

	// 1. Extract unique bid items
	bidMap := make(map[string]BidItem)
	for _, log := range mLogs.Logs {
		if _, ok := bidMap[log.Material.BidNumber]; !ok {
			bidMap[log.Material.BidNumber] = BidItem{
				Number:        log.Material.BidNumber,
				Name:          log.Material.Name,
				UnitOfMeasure: log.Material.Unit,
			}
		}
	}
	var bidItems []BidItem
	for _, b := range bidMap {
		bidItems = append(bidItems, b)
	}

	// 2. Determine sheet start/end based on logs
	layout := "2006-01-02"
	earliest := mLogs.Logs[0].Date
	latest := mLogs.Logs[0].Date
	for _, log := range mLogs.Logs {
		if log.Date < earliest {
			earliest = log.Date
		}
		if log.Date > latest {
			latest = log.Date
		}
	}

	currentStart := earliest
	endDate := latest

	var sections []ProgressSheetSection

	for {
		startTime, _ := time.Parse(layout, currentStart)
		endTime, _ := time.Parse(layout, endDate)

		// Calculate sectionEnd safely
		year, month, _ := startTime.Date()
		day := sectionEndDay
		lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, startTime.Location()).Day()
		if day > lastDay {
			day = lastDay
		}
		sectionEnd := time.Date(year, month, day, 0, 0, 0, 0, startTime.Location())
		if sectionEnd.Before(startTime) {
			sectionEnd = sectionEnd.AddDate(0, 1, 0)
		}
		if sectionEnd.After(endTime) {
			sectionEnd = endTime
		}

		// Create section
		secRows := ConvertLogsToProgressSection(mLogs.Logs, bidItems, startTime.Format(layout), sectionEnd.Format(layout))
		if len(secRows.Rows) > 0 {
			sections = append(sections, secRows)
		}

		if !sectionEnd.Before(endTime) {
			break
		}

		currentStart = sectionEnd.AddDate(0, 0, 1).Format(layout)
	}

	return ProgressSheet{
		JobDetail: fmt.Sprintf("%s %s", mLogs.Job.Number, mLogs.Job.Name),
		BidItems:  bidItems,
		Sections:  sections,
	}, nil
}

// Convert logs into a section with rows aggregated by date
func ConvertLogsToProgressSection(logs []domain.MaterialLog, bidItems []BidItem, sectionStart, sectionEnd string) ProgressSheetSection {
	layout := "2006-01-02"
	start, _ := time.Parse(layout, sectionStart)
	end, _ := time.Parse(layout, sectionEnd)

	rowMap := make(map[string]map[string]float64) // Date -> BidNumber -> Quantity

	for _, log := range logs {
		logDate, _ := time.Parse(layout, log.Date)
		if logDate.Before(start) || logDate.After(end) {
			continue
		}
		if _, ok := rowMap[log.Date]; !ok {
			rowMap[log.Date] = make(map[string]float64)
		}
		rowMap[log.Date][log.Material.BidNumber] += log.Quantity
	}

	// Create ordered rows
	var dates []string
	for date := range rowMap {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	var rows []ProgressRow
	for _, date := range dates {
		rows = append(rows, ProgressRow{
			Date:       date,
			Quantities: rowMap[date],
		})
	}

	return ProgressSheetSection{
		FromDate: sectionStart,
		ToDate:   sectionEnd,
		Rows:     rows,
	}
}
