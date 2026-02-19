package excel

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func colLetter(idx int) string {
	name, _ := excelize.ColumnNumberToName(idx)
	return name
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	//default to true. Don't want to overwrite files if there's an error checking existence
	return true
}

func renameFileSheet(f *excelize.File, index int, newName string) error {
	oldSheetName := f.GetSheetName(index)
	if err := f.SetSheetName(oldSheetName, newName); err != nil {
		return fmt.Errorf("initial excel setup: failed to rename sheet from %s to %s: %w", oldSheetName, newName, err)
	}
	return nil
}
func renameBaseSheet(f *excelize.File, newName string) error {
	return renameFileSheet(f, 0, newName)
}
