package editbox

import (
	"github.com/nsf/termbox-go"
)

type SelectBox struct {
	items         []string
	cursor        int
	scroll        int
	x, y          int
	width, height int
	fg, bg        termbox.Attribute
	sfg, sbg      termbox.Attribute
}

func (sbox *SelectBox) scrollToCursor() {
	if sbox.cursor-sbox.scroll >= sbox.height {
		sbox.scroll = sbox.cursor - sbox.height + 1
	}
	if sbox.cursor < sbox.scroll {
		sbox.scroll = sbox.cursor
	}
}

func (sbox *SelectBox) cursorDown() {
	if sbox.cursor < len(sbox.items)-1 {
		sbox.cursor++
		sbox.scrollToCursor()
	}
}

func (sbox *SelectBox) cursorUp() {
	if sbox.cursor > 0 {
		sbox.cursor--
	}
	sbox.scrollToCursor()
}

func (sbox *SelectBox) pageDown() {
	sbox.cursor = sbox.cursor + sbox.height - 1
	if sbox.cursor >= len(sbox.items) {
		sbox.cursor = len(sbox.items) - 1
	}
	sbox.scrollToCursor()
}

func (sbox *SelectBox) pageUp() {
	sbox.cursor = sbox.cursor - sbox.height + 1
	if sbox.cursor < 0 {
		sbox.cursor = 0
	}
	sbox.scrollToCursor()
}

func (sbox *SelectBox) SelectedIndex() int {
	return sbox.cursor
}

func (sbox *SelectBox) Text() string {
	return sbox.items[sbox.cursor]
}

func (sbox *SelectBox) Render() {
	var index int
	var fg, bg termbox.Attribute
	for i := 0; i < sbox.height; i++ {
		index = i + sbox.scroll
		if index > len(sbox.items)-1 {
			break
		}
		if index == sbox.cursor {
			fg, bg = sbox.sfg, sbox.sbg
		} else {
			fg, bg = sbox.fg, sbox.bg
		}
		Label(sbox.x, sbox.y+i, sbox.width, fg, bg, sbox.items[index])
	}
}

// Processes termbox events.
// Useful if you poll them by yourself.
// Returns false on unknown event.
func (sbox *SelectBox) HandleEvent(ev termbox.Event) bool {
	if ev.Type != termbox.EventKey {
		return false
	}
	switch {
	case ev.Key == termbox.KeyArrowDown:
		sbox.cursorDown()
	case ev.Key == termbox.KeyArrowUp:
		sbox.cursorUp()
	case ev.Key == termbox.KeyPgdn:
		sbox.pageDown()
	case ev.Key == termbox.KeyPgup:
		sbox.pageUp()
	default:
		return false
	}
	return true
}

// Make widget to listen for termbox events
// Blocks until exit event.
// Returns event which made SelectBox to exit, selected index,
// and selected value
func (sbox *SelectBox) WaitExit() termbox.Event {
	sbox.Render()
	termbox.Flush()
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventError {
			panic(ev.Err)
		}
		if !sbox.HandleEvent(ev) && ev.Type == termbox.EventKey {
			switch {
			case ev.Key == termbox.KeyEnter:
				return ev
			case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyTab:
				return ev
			default:
				// do nothing
			}
		}
		sbox.Render()
		termbox.Flush()
	}
}

//----------------------------------------------------------------------------
// Widgets
//----------------------------------------------------------------------------

// Create new Select widget. This DOES NOT call termbox.Flush().
func Select(
	x, y, width, height int,
	fg, bg, sfg, sbg termbox.Attribute,
	items []string,
) *SelectBox {
	var sbox SelectBox
	sbox.items = items
	sbox.x = x
	sbox.y = y
	sbox.width = width
	sbox.height = height
	sbox.fg = fg
	sbox.bg = bg
	sbox.sfg = sfg
	sbox.sbg = sbg
	return &sbox
}
