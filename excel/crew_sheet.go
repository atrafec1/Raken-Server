package excel

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type MyExcel struct {
	File   *excelize.File
	Styles map[string]int
}

func (m *MyExcel) InitializeStyles() error {
	DateHeaderStyleID, err := m.File.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 2},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create date header style: %v", err)
	}

	MainDateFormat := "dddd, mmmm dd, yyyy"
	MainDateStyleID, err := m.File.NewStyle(&excelize.Style{
		CustomNumFmt: &MainDateFormat,
		Font: &excelize.Font{
			Bold: true,
			Size: 22,
			Underline: "single",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create main date style: %v", err)
	}
	
	TitleStyleID, err := m.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 14,
		},
	})
	m.Styles["MainDate"] = MainDateStyleID
	m.Styles["DateHeader"] = DateHeaderStyleID
	m.Styles["Title"] =TitleStyleID

	// JobHeaderStyleID, err := m.File.NewStyle(&excelize.Style{
	// 	Font: &excelize.Font{ 
	// 		Underline: "single",
	// 	},
	// })
	return nil
}

// func (m *MyExcel) CreateCrewTable(jobNumber int, startCell, endCell string) error {
// 	m.File.AddTable("Sheet1", &excelize.Table{
// 		Range: fmt.Sprintf("%v:%v", startCell, endCell)
		

// 	})
// }
func (m *MyExcel) SetHeaderValues() error {
    currentTime := time.Now()
    currentDay := currentTime.Weekday()
    endOfWeekSunday := currentTime.AddDate(0, 0, 6-int(currentDay))
    week := make([]time.Time, 7)
    for i := 0; i < 7; i++ {
        week[i] = endOfWeekSunday.AddDate(0, 0, i-6)
    }

    // SECTION 1: The "01/02" Dates (Row 5)
    dates := map[string]string {
        "F5": week[0].Format("01/02"),
        "G5": week[1].Format("01/02"),
        "H5": week[2].Format("01/02"),
        "I5": week[3].Format("01/02"),
        "J5": week[4].Format("01/02"),
        "K5": week[5].Format("01/02"),
    }

    // SECTION 2: The "M, T, W" Days (Row 4)
    days := map[string]string {
        "F4": "M",
        "G4": "T",
        "H4": "W",
        "I4": "Th",
        "J4": "F",
        "K4": "S", // Fixed: replaced "." with ","
    }
	mainDateCell := "F1"
	titleCell := "A1"
    // SECTION 3: General Headers
    headerCells := map[string]string{
        "B5": "Last Name",
        "C5": "First Name",
        "D5": "Class",
        "E5": "Equip",
    }
	fmt.Println("End of week Sunday Formatted and not Formatted:", endOfWeekSunday.Format("2006-01-02"), endOfWeekSunday)
    DateHeaderStyleID := m.Styles["DateHeader"]
	MainDateStyleID := m.Styles["MainDate"]
	TitleStyleID := m.Styles["Title"]
	
	m.File.SetCellValue("Sheet1", titleCell, "TIMESHEET REVIEW / RECAP")
	m.File.SetCellStyle("Sheet1", titleCell, titleCell, TitleStyleID)

	m.File.SetCellValue("Sheet1", mainDateCell, endOfWeekSunday.Format("Monday, January 02, 2006"))
	m.File.SetCellStyle("Sheet1", mainDateCell, mainDateCell, MainDateStyleID)
    // Apply General Headers (No special style)
    for cell, value := range headerCells {
        m.File.SetCellValue("Sheet1", cell, value)
    }

    // Apply Days (Row 4) with the Border/Bold Style
    for cell, value := range days {
        m.File.SetCellValue("Sheet1", cell, value)
    }

    // Apply Dates (Row 5) with the Border/Bold Style
    for cell, value := range dates {
        m.File.SetCellValue("Sheet1", cell, value)
        m.File.SetCellStyle("Sheet1", cell, cell, DateHeaderStyleID)
    }

    return nil
}

func CreateCrewAllocationSheet(filename string) error {
	f := excelize.NewFile()
	styles := make(map[string]int)
	myExcel := &MyExcel{File: f, Styles: styles}
	
	if err := myExcel.InitializeStyles(); err != nil {
		return err
	}

	if err := myExcel.SetHeaderValues(); err != nil {
		return err
	}

	sheet2 := "Sheet2"
	f.NewSheet(sheet2)
	f.SetCellValue(sheet2, "A2", "This is on Sheet 2!")

	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %v", err)
	}

	return nil
}