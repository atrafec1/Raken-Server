package main

import (
	reports "daily_check_in/Reports"
	"daily_check_in/api"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	fyneApp := app.New()
	window := fyneApp.NewWindow("Daily Report Exporter")
	window.Resize(fyne.NewSize(560, 420))

	fromEntry := widget.NewEntry()
	fromEntry.SetPlaceHolder("YYYY-MM-DD")
	toEntry := widget.NewEntry()
	toEntry.SetPlaceHolder("YYYY-MM-DD")

	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("Select output folder")

	status := widget.NewMultiLineEntry()
	status.SetReadOnly(true)
	status.SetMinRowsVisible(8)

	appendStatus := func(message string) {
		if strings.TrimSpace(status.Text) == "" {
			status.SetText(message)
			return
		}
		status.SetText(status.Text + "\n" + message)
	}

	browseButton := widget.NewButton("Browse...", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if uri == nil {
				return
			}
			outputEntry.SetText(uri.Path())
		}, window).Show()
	})

	runButton := widget.NewButton("Download Reports", func() {
		fromDate := strings.TrimSpace(fromEntry.Text)
		toDate := strings.TrimSpace(toEntry.Text)
		outputDir := strings.TrimSpace(outputEntry.Text)

		if _, err := time.Parse("2006-01-02", fromDate); err != nil {
			dialog.ShowError(fmt.Errorf("from date must be YYYY-MM-DD"), window)
			return
		}
		if _, err := time.Parse("2006-01-02", toDate); err != nil {
			dialog.ShowError(fmt.Errorf("to date must be YYYY-MM-DD"), window)
			return
		}
		if outputDir == "" {
			dialog.ShowError(fmt.Errorf("output folder is required"), window)
			return
		}

		runButton.Disable()
		appendStatus("Starting export...")

		go func() {
			cfg, err := api.LoadConfig()
			if err != nil {
				fyneApp.Driver().RunOnMain(func() {
					runButton.Enable()
					dialog.ShowError(err, window)
				})
				return
			}

			client, err := api.NewClient(cfg)
			if err != nil {
				fyneApp.Driver().RunOnMain(func() {
					runButton.Enable()
					dialog.ShowError(err, window)
				})
				return
			}

			filePath, err := reports.ExportReports(fromDate, toDate, outputDir, client)
			fyneApp.Driver().RunOnMain(func() {
				runButton.Enable()
				if err != nil {
					dialog.ShowError(err, window)
					appendStatus("Export failed: " + err.Error())
					return
				}
				appendStatus("Saved: " + filePath)
			})
		}()
	})

	form := container.NewGridWithColumns(2,
		widget.NewLabel("From date"), fromEntry,
		widget.NewLabel("To date"), toEntry,
		widget.NewLabel("Output folder"), container.NewBorder(nil, nil, nil, browseButton, outputEntry),
	)

	content := container.NewBorder(form, nil, nil, nil,
		container.NewVBox(
			runButton,
			widget.NewSeparator(),
			widget.NewLabel("Status"),
			status,
		),
	)

	window.SetContent(content)
	window.ShowAndRun()
}
