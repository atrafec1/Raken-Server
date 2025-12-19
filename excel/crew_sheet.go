package excel

import (
	"fmt"
	"strings"
	"time"
	"sort"
	"github.com/xuri/excelize/v2"
	"daily_check_in/api"
)

type MyExcel struct {
	File   *excelize.File
	Styles map[string]int
}

func (m *MyExcel) InitializeStyles() error {
	m.Styles = make(map[string]int)

	dateHeaderStyle, _ := m.File.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "bottom", Color: "000000", Style: 2}},
		Font:   &excelize.Font{Bold: true},
	})

	mainDateFormat := "dddd, mmmm dd, yyyy"
	mainDateStyle, _ := m.File.NewStyle(&excelize.Style{
		CustomNumFmt: &mainDateFormat,
		Font:         &excelize.Font{Bold: true, Size: 22, Underline: "single"},
	})

	titleStyle, _ := m.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 14},
	})

	jobHeaderStyle, _ := m.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Underline: "single"},
	})

	m.Styles["DateHeader"] = dateHeaderStyle
	m.Styles["MainDate"] = mainDateStyle
	m.Styles["Title"] = titleStyle
	m.Styles["JobHeader"] = jobHeaderStyle

	return nil
}

func (m *MyExcel) SetHeaderValues(week []time.Time) {
	// 1. Static Headers
	m.File.SetCellValue("Sheet1", "A1", "TIMESHEET REVIEW / RECAP")
	m.File.SetCellStyle("Sheet1", "A1", "A1", m.Styles["Title"])

	// 2. Main Date Display (Using the Friday/Sunday as the "Week Ending")
	weDate := week[len(week)-1].Format("Monday, January 02, 2006")
	m.File.SetCellValue("Sheet1", "F1", weDate)
	m.File.SetCellStyle("Sheet1", "F1", "F1", m.Styles["MainDate"])

	// 3. Table Column Headers
	headers := map[string]string{
		"A5": "Last Name", "B5": "First Name", "C5": "Class", "D5": "Equip",
	}
	for cell, val := range headers {
		m.File.SetCellValue("Sheet1", cell, val)
	}

	// 4. M-S Days and Dates (F-K)
	dayNames := []string{"M", "T", "W", "Th", "F", "S"}
	for i, t := range week {
		if i >= 6 { break } // Only map M-S
		col, _ := excelize.ColumnNumberToName(6 + i) // Starts at F (6)
		
		// Day Name (Row 4)
		m.File.SetCellValue("Sheet1", col+"4", dayNames[i])
		
		// Date String (Row 5)
		m.File.SetCellValue("Sheet1", col+"5", t.Format("01/02"))
		m.File.SetCellStyle("Sheet1", col+"5", col+"5", m.Styles["DateHeader"])
	}
}

func (m *MyExcel) BuildTable(todaysCrews []api.CrewAllocationEntry, allEntries []api.CrewAllocationEntry, week []time.Time) {
	// Map dates to columns for fast lookup
	dateToCol := make(map[string]string)
	for i, t := range week {
		col, _ := excelize.ColumnNumberToName(6 + i)
		dateToCol[t.Format("2006-01-02")] = col
	}

	// 1. Sort today's crews by Job Number
	sort.Slice(todaysCrews, func(i, j int) bool {
		return todaysCrews[i].Project.Number < todaysCrews[j].Project.Number
	})

	currentRow := 7
	for _, crew := range todaysCrews {
		// 2. Job Header Row
		jobLabel := fmt.Sprintf("JOB %s %s", crew.Project.Number, crew.Project.Name)
		m.File.SetCellValue("Sheet1", "A"+fmt.Sprint(currentRow), jobLabel)
		m.File.SetCellStyle("Sheet1", "A"+fmt.Sprint(currentRow), "A"+fmt.Sprint(currentRow), m.Styles["JobHeader"])
		currentRow++

		for i, emp := range crew.Employees {
			// 3. Employee Info (Left Side)
			m.File.SetCellValue("Sheet1", "A"+fmt.Sprint(currentRow), i+1) // Index
			m.File.SetCellValue("Sheet1", "A"+fmt.Sprint(currentRow), emp.LastName)
			m.File.SetCellValue("Sheet1", "A"+fmt.Sprint(currentRow), emp.FirstName)
			m.Files.SetCellValue("Sheet1", "A"+fmt.Sprint(currentRow), emp.Class)

			// 4. History (Right Side) - Aligned to current row
			history := api.GetCrewMemberHistory(emp.UUID, allEntries)
			for dateStr, colLetter := range dateToCol {
				if projects, ok := history.Projects[dateStr]; ok {
					displayVal := strings.Join(projects, " / ")
					m.File.SetCellValue("Sheet1", colLetter+fmt.Sprint(currentRow), displayVal)
				}
			}
			currentRow++
		}
		currentRow++ // Spacer row between jobs
	}
}

// --- Entry Point ---

func CreateCrewAllocationSheet(filename string, allCrews []api.CrewAllocationEntry) error {
	f := excelize.NewFile()
	defer f.Close()

	myExcel := &MyExcel{File: f}
	myExcel.InitializeStyles()

	// Calculate the week (Monday to Saturday)
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 { offset = -6 } // Handle Sundays
	
	monday := now.AddDate(0, 0, offset)
	week := make([]time.Time, 6)
	for i := 0; i < 6; i++ {
		week[i] = monday.AddDate(0, 0, i)
	}

	myExcel.SetHeaderValues(week)

	todaysCrews := api.GetTodaysCrewAllocations(allCrews)
	myExcel.BuildTable(todaysCrews, allCrews, week)

	// Set column widths for readability
	f.SetColWidth("Sheet1", "A", "A", 30)
	f.SetColWidth("Sheet1", "B", "D", 15)

	return f.SaveAs(filename)
}