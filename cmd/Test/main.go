package main

import (
	"daily_check_in/payroll/adapter"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

func main() {
	adapter, err := adapter.NewRakenAPIAdapter()
	if err != nil {
		panic(err)
	}

	entries, err := adapter.GetPayrollEntries("2026-02-02", "2026-02-08")
	if err != nil {
		panic(err)
	}
	if err := WriteCSV("payroll_entries.csv", entries); err != nil {
		panic(err)
	}

}

func WriteCSV[T any](filename string, data []T) error {
	if len(data) == 0 {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	val := reflect.ValueOf(data[0])
	typ := val.Type()

	// headers
	headers := make([]string, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		headers[i] = typ.Field(i).Name
	}
	if err := w.Write(headers); err != nil {
		return err
	}

	// rows
	for _, item := range data {
		v := reflect.ValueOf(item)
		row := make([]string, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			row[i] = toString(v.Field(i).Interface())
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func toString(v any) string {
	return fmt.Sprint(v)
}
