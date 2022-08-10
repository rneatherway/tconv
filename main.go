package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/github/go-freetime"
	"github.com/rivo/tview"
)

func label(s string) *tview.TextView {
	return tview.NewTextView().
		SetText(s + ":").
		SetTextColor(tcell.ColorYellow)
}

func realMain() error {
	t := time.Now()
	if len(os.Args) == 2 {
		parser, err := freetime.NewParser()
		if err != nil {
			return err
		}

		t, err = parser.Parse(os.Args[1])
		if err != nil {
			return err
		}
	}

	app := tview.NewApplication()
	pages := tview.NewPages()
	grid := tview.NewGrid().
		SetRows(1, 1).
		SetColumns(15, 0)

	timeWidget := NewTimeWidget(t)
	unixFld := tview.NewInputField()
	unixView := tview.NewTextView().SetText(strconv.FormatInt(t.Unix(), 10))

	timeWidget.SetChangedFunc(func(t time.Time, err error) {
		if err == nil {
			text := strconv.FormatInt(t.Unix(), 10)
			unixFld.SetText(text)
			unixView.SetText(text)
		}
	})

	grid.Box = tview.NewBox()
	grid.SetBorders(true).
		AddItem(label("ISO time"), 0, 0, 1, 1, 0, 0, false).
		AddItem(timeWidget, 0, 1, 1, 1, 0, 0, true).
		AddItem(label("UNIX timestamp"), 1, 0, 1, 1, 0, 0, false).
		AddItem(unixView, 1, 1, 1, 1, 0, 0, false)
	pages.AddPage("main", grid, false, true)

	return app.SetRoot(grid, true).Run()
}

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: %s [TIME]", os.Args[0])
		os.Exit(1)
	}

	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
