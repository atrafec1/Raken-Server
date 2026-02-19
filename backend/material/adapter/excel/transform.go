package excel

import (
	"fmt"
	"prg_tools/material/domain"
	"sort"
	"strconv"
	"time"
)

// Entry point — converts JobMaterialInfo into multiple sheets (weekly, monthly, etc.)

func CreateProgressSheets(mLogs domain.JobMaterialInfo) []ProgressSheet {
	if len(mLogs.Logs) == 0 {
		return nil
	}

	bidItems := extractBidItems(mLogs)
	sortBidItems(bidItems)

	ranges := buildWeeklyRanges(mLogs.Logs)

	var sheets []ProgressSheet
	for _, r := range ranges {
		rows := buildRowsForRange(mLogs.Logs, r.From, r.To)
		if len(rows) == 0 {
			continue
		}

		sheets = append(sheets, ProgressSheet{
			SheetName: buildSheetName(r.From, r.To),
			JobDetail: fmt.Sprintf("%s %s", mLogs.Job.Number, mLogs.Job.Name),
			BidItems:  bidItems,
			Rows:      rows,
			FromDate:  r.From,
			ToDate:    r.To,
		})
	}
	sortSheetsMostRecentFirst(sheets)
	return sheets
}

type DateRange struct {
	From string
	To   string
}

func buildWeeklyRanges(logs []domain.MaterialLog) []DateRange {
	layout := "2006-01-02"

	earliest := logs[0].Date
	latest := logs[0].Date

	for _, l := range logs {
		if l.Date < earliest {
			earliest = l.Date
		}
		if l.Date > latest {
			latest = l.Date
		}
	}

	start, _ := time.Parse(layout, getMonday(earliest))
	end, _ := time.Parse(layout, getSunday(latest))

	var ranges []DateRange

	for !start.After(end) {
		weekEnd := start.AddDate(0, 0, 6)
		if weekEnd.After(end) {
			weekEnd = end
		}

		ranges = append(ranges, DateRange{
			From: start.Format(layout),
			To:   weekEnd.Format(layout),
		})

		start = weekEnd.AddDate(0, 0, 1)
	}

	return ranges
}

// ------------------ helpers ------------------

func extractBidItems(mLogs domain.JobMaterialInfo) []BidItem {
	bidMap := make(map[string]BidItem)

	for _, mat := range mLogs.Materials {
		if _, ok := bidMap[mat.BidNumber]; !ok {
			bidMap[mat.BidNumber] = BidItem{
				Number:        mat.BidNumber,
				Name:          mat.Name,
				UnitOfMeasure: mat.Unit,
			}
		}
	}

	var bidItems []BidItem
	for _, b := range bidMap {
		bidItems = append(bidItems, b)
	}
	return bidItems
}

func buildRowsForRange(
	logs []domain.MaterialLog,
	from, to string,
) []ProgressRow {

	layout := "2006-01-02"
	start, _ := time.Parse(layout, from)
	end, _ := time.Parse(layout, to)

	rowMap := make(map[string]map[string]float64)

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

	// Sort dates
	var dates []string
	for d := range rowMap {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	var rows []ProgressRow
	for _, d := range dates {
		rows = append(rows, ProgressRow{
			Date:       d,
			Quantities: rowMap[d],
		})
	}

	return rows
}

func sortBidItems(bidItems []BidItem) {
	sort.Slice(bidItems, func(i, j int) bool {
		numberI, _ := strconv.Atoi(bidItems[i].Number)
		numberJ, _ := strconv.Atoi(bidItems[j].Number)
		return numberI < numberJ
	})
}

func buildSheetName(from, to string) string {
	toDate, err := time.Parse("2006-01-02", to)
	if err != nil {
		return fmt.Sprintf("%s", to)
	}
	formatted := toDate.Format("01-02-06")
	return formatted
}

func sortSheetsMostRecentFirst(sheets []ProgressSheet) {
	layout := "2006-01-02"
	sort.Slice(sheets, func(i, j int) bool {
		dateI, _ := time.Parse(layout, sheets[i].ToDate)
		dateJ, _ := time.Parse(layout, sheets[j].ToDate)
		return dateI.After(dateJ) // descending: most recent first
	})
}
