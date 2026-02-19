package excel

import (
	"fmt"
	"prg_tools/material/domain"
	"sort"
	"strconv"
	"time"
)

func ConvertToProgressSheet(mLogs domain.JobMaterialInfo, sectionEndDay int) (ProgressSheet, error) {
	if len(mLogs.Logs) == 0 {
		return ProgressSheet{}, nil
	}

	// 1. Extract unique bid items
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

func sortBidItems(bidItems []BidItem) {
	sort.Slice(bidItems, func(i, j int) bool {
		numberI, _ := strconv.Atoi(bidItems[i].Number)
		numberJ, _ := strconv.Atoi(bidItems[j].Number)
		return numberI < numberJ
	})
}
