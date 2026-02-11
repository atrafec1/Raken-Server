package domain

import (
	"fmt"
	"time"
)

type ReportCollection struct {
	FromDate string
	ToDate   string
	Reports  []Report
}

type Report struct {
	Date    string
	Creator ReportCreator
	PDFLink string
	Project Project
}

type ReportCreator struct {
	Name string
}

func (r Report) YearWeek() (int, int, error) {
	t, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing report date: %w", err)
	}
	year, week := t.ISOWeek()
	return year, week, nil
}

// Returns FileName as
func (r Report) ToFileName() string {
	return fmt.Sprintf("%s.pdf", r.Date)
}

