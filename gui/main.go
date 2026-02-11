package main

import (
	"daily_check_in/payroll"
	"daily_check_in/payroll/adapter/cp"
	"daily_check_in/payroll/adapter/raken"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Initialize payroll service (adapter injected as needed)
	adapter, err := raken.NewRakenAPIAdapter() // adapt as necessary
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Raken adapter: %v", err))
	}
	exporter := cp.NewAdapter("C:\\Users\\EMarin\\Desktop")        // adapt as necessary
	payrollService := payroll.NewPayrollService(adapter, exporter) // adapt as necessary

	// Create Fyne app and window
	a := app.New()
	w := a.NewWindow("Payroll Generator")
	w.Resize(fyne.NewSize(500, 300))

	// Week ending date entry
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("YYYY-MM-DD")

	// Status/output area
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Status / Warnings")
	output.Wrapping = fyne.TextWrapWord
	output.Disable() // read-only

	// Last Week button
	lastWeekBtn := widget.NewButton("Last Week", func() {
		lastSunday := lastSundayDate()
		dateEntry.SetText(lastSunday.Format("2006-01-02"))
	})

	// Generate Payroll Files button
	generateBtn := widget.NewButton("Generate Payroll Files", func() {
		output.SetText("") // clear previous output
		dateText := dateEntry.Text
		if dateText == "" {
			output.SetText("Error: Please enter a week ending date.")
			return
		}

		// Determine payroll week range (last Monday → last Friday)
		weekStart, weekEnd, err := weekRangeFromEndingDate(dateText)
		if err != nil {
			output.SetText(fmt.Sprintf("Error parsing date: %v", err))
			return
		}

		// Call payroll service
		result, err := payrollService.GetEntries(weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))
		if err != nil {
			output.SetText(fmt.Sprintf("Error fetching payroll entries: %v", err))
			return
		}

		// Show warnings
		if len(result.Warnings) > 0 {
			output.SetText("Warnings:\n" + formatWarnings(result.Warning))
		}

		// Export CSV
		err = payrollService.Export(result.Entries)
		if err != nil {
			output.SetText(output.Text + fmt.Sprintf("\nError exporting CSV: %v", err))
			return
		}

		output.SetText(output.Text + fmt.Sprintf("\nSuccess! Exported %d entries.", len(result.Entries)))
	})

	// Layout
	dateContainer := container.NewHBox(dateEntry, lastWeekBtn)
	content := container.NewVBox(
		dateContainer,
		generateBtn,
		output,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

// --- Helpers ---

func lastSundayDate() time.Time {
	today := time.Now()
	offset := (int(today.Weekday()) + 7 - 0) % 7 // 0 = Sunday
	lastSunday := today.AddDate(0, 0, -offset)
	return lastSunday
}

// Compute week start (Monday) and end (Friday) from ending date
func weekRangeFromEndingDate(weekEnding string) (time.Time, time.Time, error) {
	t, err := time.Parse("2006-01-02", weekEnding)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	// Assuming week ends Friday
	// Adjust if ending date is Sunday
	weekday := t.Weekday()
	var weekEnd time.Time
	if weekday == time.Sunday {
		weekEnd = t.AddDate(0, 0, -2) // Friday
	} else {
		weekEnd = t
	}
	weekStart := weekEnd.AddDate(0, 0, -4) // Monday
	return weekStart, weekEnd, nil
}

func formatWarnings(warnings []string) string {
	result := ""
	for _, w := range warnings {
		result += "- " + w + "\n"
	}
	return result
}
