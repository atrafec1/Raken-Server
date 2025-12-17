package excel

import (
    "fmt"
    "github.com/xuri/excelize/v2"
)

func CreateCrewAllocationSheet() {
    f := excelize.NewFile()

    // Create a new sheet
    index := f.NewSheet("Sheet1")

    // Set value in a cell
    f.SetCellValue("Sheet1", "A1", "Hello")
    f.SetCellValue("Sheet1", "B1", 123)

    // Set active sheet
    f.SetActiveSheet(index)

    // Save the file
    if err := f.SaveAs("Book1.xlsx"); err != nil {
        fmt.Println(err)
    }
}
