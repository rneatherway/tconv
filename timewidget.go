package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TimeWidget struct {
	*tview.Box

	content string

	cursorPos int

	changed func(t time.Time, err error)
}

func NewTimeWidget(t time.Time) *TimeWidget {
	return &TimeWidget{
		Box:       tview.NewBox(),
		content:   t.UTC().Format(time.RFC3339),
		cursorPos: 0,
	}
}

func (w *TimeWidget) parse() (time.Time, error) {
	return time.Parse(time.RFC3339, w.content)
}

func (w *TimeWidget) Draw(screen tcell.Screen) {
	w.Box.DrawForSubclass(screen, w)

	x, y, _, _ := w.GetInnerRect()

	text := w.content
	if _, err := w.parse(); err != nil {
		text = text + " [yellow]"
		idx := strings.Index(err.Error(), w.content)
		if idx == -1 {
			text += err.Error()
		} else {
			text += err.Error()[idx+len(w.content)+3:]
		}
	}

	tview.PrintSimple(screen, text, x, y)

	if w.HasFocus() {
		screen.ShowCursor(x+w.cursorPos, y)
	}
}

func moveRight(cursorPos int) int {
	switch cursorPos {
	case 18:
		// do nothing
	case 3, 6, 9, 12, 15:
		cursorPos += 2
	default:
		cursorPos += 1
	}
	return cursorPos
}

func (w *TimeWidget) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return w.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		changed := false

		switch event.Key() {
		case tcell.KeyLeft:
			switch w.cursorPos {
			case 0:
				// do nothing
			case 5, 8, 11, 14, 17:
				w.cursorPos -= 2
			default:
				w.cursorPos -= 1
			}
		case tcell.KeyRight:
			w.cursorPos = moveRight(w.cursorPos)
		case tcell.KeyBacktab:
			switch w.cursorPos {
			case 0, 1, 2, 3:
				w.cursorPos = 17
			case 5, 6:
				w.cursorPos = 0
			case 8, 9:
				w.cursorPos = 5
			case 11, 12:
				w.cursorPos = 8
			case 14, 15:
				w.cursorPos = 11
			case 17, 18:
				w.cursorPos = 14
			}
		case tcell.KeyTab:
			switch w.cursorPos {
			case 0, 1, 2, 3:
				w.cursorPos = 5
			case 5, 6:
				w.cursorPos = 8
			case 8, 9:
				w.cursorPos = 11
			case 11, 12:
				w.cursorPos = 14
			case 14, 15:
				w.cursorPos = 17
			case 17, 18:
				w.cursorPos = 0
			}
		case tcell.KeyUp, tcell.KeyDown:
			if t, err := w.parse(); err == nil {
				n := 1
				if event.Key() == tcell.KeyUp {
					n = -1
				}

				switch w.cursorPos {
				case 0, 1, 2, 3:
					t = t.AddDate(n, 0, 0)
				case 5, 6:
					t = t.AddDate(0, n, 0)
				case 8, 9:
					t = t.AddDate(0, 0, n)
				case 11, 12:
					t = t.Add(time.Duration(n) * time.Hour)
				case 14, 15:
					t = t.Add(time.Duration(n) * time.Minute)
				case 17, 18:
					t = t.Add(time.Duration(n) * time.Second)
				}
				w.content = t.Format(time.RFC3339)
				changed = true
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				valid := true
				n, _ := strconv.Atoi(string(event.Rune()))
				switch w.cursorPos {
				case 5:
					valid = n <= 1
				case 8:
					valid = n >= 1 && n <= 3
				case 11:
					valid = n <= 2
				case 14, 17:
					valid = n <= 5
				}

				if valid {
					rs := []rune(w.content)
					rs[w.cursorPos] = event.Rune()
					w.content = string(rs)
					changed = true

					w.cursorPos = moveRight(w.cursorPos)
				}
			}
		}

		if changed {
			w.changed(w.parse())
		}
	})
}

func (w *TimeWidget) SetChangedFunc(handler func(t time.Time, err error)) *TimeWidget {
	w.changed = handler
	return w
}
