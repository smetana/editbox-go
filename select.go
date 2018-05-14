package editbox

import (
	"github.com/nsf/termbox-go"
)

type SelectBox struct {
	Items         []string
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
	if sbox.cursor < len(sbox.Items)-1 {
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
	if sbox.cursor >= len(sbox.Items) {
		sbox.cursor = len(sbox.Items) - 1
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
	return sbox.Items[sbox.cursor]
}

func (sbox *SelectBox) Render() {
	var index int
	var fg, bg termbox.Attribute
	for i := 0; i < sbox.height; i++ {
		index = i + sbox.scroll
		if index == sbox.cursor {
			fg, bg = sbox.sfg, sbox.sbg
		} else {
			fg, bg = sbox.fg, sbox.bg
		}
		Label(sbox.x, sbox.y+i, sbox.width, fg, bg, sbox.Items[index])
	}
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
		switch ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyArrowDown:
				sbox.cursorDown()
			case ev.Key == termbox.KeyArrowUp:
				sbox.cursorUp()
			case ev.Key == termbox.KeyPgdn:
				sbox.pageDown()
			case ev.Key == termbox.KeyPgup:
				sbox.pageUp()
			case ev.Key == termbox.KeyEnter:
				return ev
			case ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyTab:
				return ev
			default:
				// do nothing
			}
		case termbox.EventError:
			panic(ev.Err)
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
	sbox.Items = items
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
