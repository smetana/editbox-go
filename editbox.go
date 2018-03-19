package editbox

import (
	"bytes"
	"github.com/nsf/termbox-go"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type cursor struct {
	x, y int
}

//----------------------------------------------------------------------------
// line
//----------------------------------------------------------------------------

type line struct {
	text []rune
}

func (l *line) checkXPosition(x int) {
	if x < 0 || x > len(l.text) {
		panic("x position out of range")
	}
}

func (l *line) insertRune(pos int, r rune) {
	l.checkXPosition(pos)
	// Append
	if pos == len(l.text) {
		l.text = append(l.text, r)
		// Insert
	} else {
		l.text = append(l.text, rune(0))
		copy(l.text[pos+1:], l.text[pos:])
		l.text[pos] = r
	}
}

func (l *line) split(pos int) (left, right *line) {
	l.checkXPosition(pos)
	left, right = l, new(line)
	right.text = make([]rune, len(l.text)-pos)
	copy(right.text, l.text[pos:len(l.text)])
	left.text = left.text[:pos]
	return
}

func (l *line) deleteRune(pos int) rune {
	l.checkXPosition(pos)
	if pos < len(l.text) {
		r := l.text[pos]
		copy(l.text[pos:], l.text[pos+1:])
		l.text[len(l.text)-1] = rune(0)
		l.text = l.text[:len(l.text)-1]
		return r
	} else {
		return rune(0)
	}
}

func (l *line) lastRune() rune {
	return l.text[len(l.text)-1]
}

func (l *line) lastRuneX() int {
	if l.lastRune() == '\n' {
		return (len(l.text) - 1)
	} else {
		return (len(l.text))
	}
}

//----------------------------------------------------------------------------
// editor
//----------------------------------------------------------------------------

type editor struct {
	lines  []line
	cursor cursor
	lastx  int
}

func newEditor() *editor {
	var ed editor
	ed.lines = make([]line, 1)
	ed.cursor.x = 0
	ed.cursor.y = 0
	return &ed
}

func (ed *editor) text() string {
	var b bytes.Buffer
	for _, l := range ed.lines {
		b.WriteString(string(l.text))
	}
	return b.String()
}

func (ed *editor) currentLine() *line {
	return &ed.lines[ed.cursor.y]
}

func (ed *editor) splitLine(x, y int) {
	l := &ed.lines[y]
	left, right := l.split(x)
	ed.lines = append(ed.lines, *(new(line)))
	copy(ed.lines[y+2:], ed.lines[y+1:])
	ed.lines[y] = *left
	ed.lines[y+1] = *right
}

func (ed *editor) insertRune(r rune) {
	cursor := &ed.cursor
	line := ed.currentLine()
	line.insertRune(cursor.x, r)
	cursor.x += 1
	if r == '\n' {
		ed.splitLine(cursor.x, cursor.y)
		cursor.y += 1
		cursor.x = 0
	}
	ed.lastx = cursor.x
}

func (ed *editor) checkYPosition(y int) {
	if y < 0 || y > len(ed.lines) {
		panic("y position out of range")
	}
}

func (ed *editor) deleteRuneBeforeCursor() {
	cursor := &ed.cursor
	if cursor.x == 0 && cursor.y == 0 {
		return
	}
	ed.moveCursorLeft()
	ed.deleteRuneAtCursor()
}

func (ed *editor) deleteRuneAtCursor() {
	cursor := &ed.cursor
	l := ed.currentLine()
	r := l.deleteRune(cursor.x)
	if r == '\n' && cursor.y < len(ed.lines)-1 {
		left := &ed.lines[cursor.y]
		right := &ed.lines[cursor.y+1]
		left.text = append(left.text, right.text...)
		if cursor.y == len(ed.lines)-2 {
			ed.lines = ed.lines[:cursor.y+1]
		} else {
			copy(ed.lines[cursor.y+1:], ed.lines[cursor.y+2:])
			ed.lines[len(ed.lines)-1] = *(new(line))
			ed.lines = ed.lines[:len(ed.lines)-1]
		}
	}
}

func (ed *editor) moveCursorRight() {
	cursor := &ed.cursor
	line := ed.currentLine()
	cursor.x += 1
	if cursor.x >= len(line.text) {
		if cursor.y < len(ed.lines)-1 {
			cursor.y += 1
			cursor.x = 0
		} else {
			cursor.x = len(line.text)
		}
	}
	ed.lastx = cursor.x
}

func (ed *editor) moveCursorLeft() {
	cursor := &ed.cursor
	cursor.x -= 1
	if cursor.x < 0 {
		if cursor.y > 0 {
			cursor.y -= 1
			line := ed.currentLine()
			cursor.x = len(line.text) - 1
		} else {
			cursor.x = 0
		}
	}
	ed.lastx = cursor.x
}

func (ed *editor) moveCursorToLineStart() {
	ed.cursor.x, ed.lastx = 0, 0
}

func (ed *editor) moveCursorToLineEnd() {
	line := ed.currentLine()
	if line.lastRune() == '\n' {
		ed.cursor.x = len(line.text) - 1
	} else {
		ed.cursor.x = len(line.text)
	}
	ed.lastx = ed.cursor.x
}

func (ed *editor) moveCursorVert(dy int) {
	cursor := &ed.cursor
	if cursor.y+dy < 0 {
		return
	}
	if cursor.y+dy > len(ed.lines)-1 {
		return
	}
	cursor.y += dy
	line := ed.currentLine()
	switch {
	case len(line.text) == 0:
		cursor.x = 0
	case ed.lastx >= len(line.text):
		cursor.x = len(line.text) - 1
	default:
		cursor.x = ed.lastx
	}
}

// TODO Refactor
func (ed *editor) setText(text string) {
	for _, s := range text {
		ed.insertRune(rune(s))
	}
}

//----------------------------------------------------------------------------
// Editbox
//----------------------------------------------------------------------------

type options struct {
	fg         termbox.Attribute
	bg         termbox.Attribute
	wrap       bool
	autoexpand bool
	maxHeight  int
	printNL    bool
	exitKeys   []termbox.Key
}

// Base type for all editbox widgets.
type Editbox struct {
	editor        *editor
	cursor        cursor
	scroll        cursor
	x, y          int
	width, height int
	wrap          bool
	fg, bg        termbox.Attribute
	autoexpand    bool
	printNL       bool
	exitKeys      []termbox.Key
	view          [][]rune
	// Line y coord in box in wrap mode
	lineBoxY      []int
	virtualHeight int
	minHeight     int
	maxHeight     int
}

func newEditbox(x, y, width, height int, options options) *Editbox {
	var ebox Editbox
	ebox.x = x
	ebox.y = y
	ebox.width = width
	ebox.height = height
	ebox.fg = options.fg
	ebox.bg = options.bg
	ebox.wrap = options.wrap
	ebox.autoexpand = options.autoexpand
	if ebox.autoexpand {
		ebox.minHeight = height
		if options.maxHeight <= 0 {
			ebox.maxHeight = ebox.minHeight
		} else {
			ebox.maxHeight = options.maxHeight
		}
	}
	ebox.printNL = options.printNL
	ebox.exitKeys = options.exitKeys
	ebox.editor = newEditor()
	return &ebox
}

func (ebox *Editbox) updateLineOffsets() {
	ed := ebox.editor
	linesCnt := len(ed.lines)
	ebox.lineBoxY = make([]int, linesCnt)
	dy := 0 // delta between editor y and box Y
	cumulativeOffset := 0
	for y := 0; y < linesCnt; y++ {
		ebox.lineBoxY[y] = y + cumulativeOffset
		if ebox.wrap {
			dy = (len(ed.lines[y].text) - 1) / ebox.width
			cumulativeOffset += dy
		}
	}
	ebox.virtualHeight = ebox.lineBoxY[linesCnt-1] + dy + 1
	if ebox.autoexpand {
		if ebox.virtualHeight > ebox.height {
			if ebox.virtualHeight > ebox.maxHeight {
				ebox.height = ebox.maxHeight
			} else {
				ebox.height = ebox.virtualHeight
			}
		} else if ebox.virtualHeight < ebox.height {
			if ebox.virtualHeight < ebox.minHeight {
				ebox.height = ebox.minHeight
			} else {
				ebox.height = ebox.virtualHeight
			}
		}
		// else Ok. Don't change height
	}
	// else Ok. don't change height
	ebox.cursor.x, ebox.cursor.y = ebox.editorToBox(ed.cursor.x, ed.cursor.y)
}

func (ebox *Editbox) editorToBox(x, y int) (int, int) {
	if ebox.wrap {
		ldy := x / ebox.width
		x = x - (ldy * ebox.width)
		y = ebox.lineBoxY[y] + ldy
	}
	return x, y
}

// Cursor movement in wrap mode is a bit tricky
// TODO Code smell. Refactor
func (ebox *Editbox) moveCursorDown() {
	if ebox.wrap {
		ed := ebox.editor
		line := ed.currentLine()
		// Try to move within current line
		if ed.cursor.x+ebox.width < len(line.text) {
			ed.cursor.x += ebox.width
			return
		}
		if ebox.cursor.x+(len(line.text)-ed.cursor.x)-1 >= ebox.width {
			ed.cursor.x = line.lastRuneX()
			return
		}
		// Jump to next line
		if ed.cursor.y+1 > len(ed.lines)-1 {
			return
		}
		ed.cursor.y += 1
		line = ed.currentLine()
		if len(line.text) == 0 {
			ed.cursor.x = 0
			return
		}
		x, _ := ebox.editorToBox(ed.lastx, 0)
		if x >= len(line.text) {
			ed.cursor.x = line.lastRuneX()
		} else {
			ed.cursor.x = x
		}
	} else {
		ebox.editor.moveCursorVert(+1)
	}
}

func (ebox *Editbox) moveCursorUp() {
	if ebox.wrap {
		ed := ebox.editor
		lastx, _ := ebox.editorToBox(ed.lastx, 0)
		x, _ := ebox.editorToBox(ed.cursor.x, 0)
		if x == lastx && ed.cursor.x-ebox.width >= 0 {
			ed.cursor.x -= ebox.width
			return
		}
		d := ebox.width + x - lastx
		if x < lastx && ed.cursor.x-d >= 0 {
			ed.cursor.x -= d
			return
		}
		if ed.cursor.y-1 < 0 {
			return
		}
		ed.cursor.y -= 1
		line := ed.currentLine()
		if ed.lastx < ebox.width {
			ed.cursor.x = ed.lastx
		}
		if lastx >= line.lastRuneX() {
			ed.cursor.x = line.lastRuneX()
			return
		}
		x, _ = ebox.editorToBox(line.lastRuneX(), 0)
		if x <= lastx {
			ed.cursor.x = line.lastRuneX()
		} else {
			ed.cursor.x = line.lastRuneX() - x + lastx
		}
	} else {
		ebox.editor.moveCursorVert(-1)
	}
}

func (ebox *Editbox) moveCursorPageUp() {
	for i := 1; i <= ebox.height; i++ {
		ebox.moveCursorUp()
	}
}

func (ebox *Editbox) moveCursorPageDown() {
	for i := 1; i <= ebox.height; i++ {
		ebox.moveCursorDown()
	}
}

func (ebox *Editbox) scrollToCursor() {
	if !ebox.wrap {
		if ebox.cursor.x-ebox.scroll.x > ebox.width-1 {
			ebox.scroll.x = ebox.cursor.x - ebox.width + 1
		} else if ebox.cursor.x-ebox.scroll.x < 0 {
			ebox.scroll.x = ebox.cursor.x
		}
	}
	if ebox.virtualHeight > ebox.height {
		if ebox.cursor.y-ebox.scroll.y > ebox.height-1 {
			ebox.scroll.y = ebox.cursor.y - ebox.height + 1
		} else if ebox.cursor.y-ebox.scroll.y < 0 {
			ebox.scroll.y = ebox.cursor.y
		} else if ebox.virtualHeight-ebox.scroll.y <= ebox.height-1 {
			ebox.scroll.y = ebox.virtualHeight - ebox.height
		}
	} else {
		ebox.scroll.y = 0
	}
}

func (ebox *Editbox) renderView() {
	ebox.updateLineOffsets()
	ebox.scrollToCursor()
	ed := ebox.editor
	var (
		boxX, boxY   int
		viewX, viewY int
	)
	ebox.view = make([][]rune, ebox.height)
	for i := range ebox.view {
		ebox.view[i] = make([]rune, ebox.width)
	}
	for y, line := range ed.lines {
		for x, r := range line.text {
			boxX, boxY = ebox.editorToBox(x, y)
			//TODO Optimize
			if boxY < ebox.scroll.y || boxX < ebox.scroll.x {
				continue
			}
			viewX = boxX - ebox.scroll.x
			viewY = boxY - ebox.scroll.y
			if viewX > ebox.width-1 {
				break
			}
			if viewY > ebox.height-1 {
				break
			}
			if r == '\n' {
				if ebox.printNL {
					r = '␤'
				} else {
					r = ' '
				}
			}
			ebox.view[viewY][viewX] = r
		}
		if viewY > ebox.height-1 {
			break
		}
	}
}

//----------------------------------------------------------------------------
// API
//----------------------------------------------------------------------------

// Returns widget content.
func (ebox *Editbox) Text() string {
	return ebox.editor.text()
}

// Puts widget contents into termbox' cell buffer.
// This function DOES NOT call termbox.Flush().
func (ebox *Editbox) Render() {
	ebox.renderView()
	var r rune
	for y := 0; y < ebox.height; y++ {
		for x := 0; x < ebox.width; x++ {
			if ebox.view[y][x] != 0 {
				r = ebox.view[y][x]
			} else {
				r = ' ' // Fill empty cells with background color
			}
			termbox.SetCell(ebox.x+x, ebox.y+y, r, ebox.fg, ebox.bg)
		}
	}
	termbox.SetCursor(ebox.x+ebox.cursor.x-ebox.scroll.x,
		ebox.y+ebox.cursor.y-ebox.scroll.y)
}

// Processes termbox events.
// Useful if you poll them by yourself.
func (ebox *Editbox) HandleEvent(ev *termbox.Event) {
	ed := ebox.editor
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowLeft:
			ed.moveCursorLeft()
		case termbox.KeyArrowRight:
			ed.moveCursorRight()
		case termbox.KeyArrowUp:
			ebox.moveCursorUp()
		case termbox.KeyArrowDown:
			ebox.moveCursorDown()
		case termbox.KeyHome:
			ed.moveCursorToLineStart()
		case termbox.KeyEnd:
			ed.moveCursorToLineEnd()
		case termbox.KeyPgup:
			ebox.moveCursorPageUp()
		case termbox.KeyPgdn:
			ebox.moveCursorPageDown()
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			ed.deleteRuneBeforeCursor()
		case termbox.KeyDelete:
			ed.deleteRuneAtCursor()
		case termbox.KeyEnter:
			ed.insertRune('\n')
		case termbox.KeySpace:
			ed.insertRune(' ')
		default:
			if ev.Ch != 0 {
				ed.insertRune(ev.Ch)
			}
		}
	case termbox.EventError:
		panic(ev.Err)
	default:
		// TODO
	}
}

// Start listen for termbox events and edit text.
// Blocks until exit event. Returns event which made Editbox to exit.
func (ebox *Editbox) WaitExit() termbox.Event {
	// Buffered channel processes paste from buffer faster
	// because render is called less often
	events := make(chan termbox.Event, 256)
	exitEvent := make(chan termbox.Event)
	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey {
				for _, key := range ebox.exitKeys {
					if ev.Key == key {
						exitEvent <- ev
						return
					}
				}
			}
			events <- ev
		}
	}()
	ebox.Render()
	termbox.Flush()
	for {
		select {
		case ev := <-events:
			ebox.HandleEvent(&ev)
			// re-render on empty events buffer
			if len(events) == 0 {
				ebox.Render()
				termbox.Flush()
			}
		case ev := <-exitEvent:
			return ev
		}
	}
}

//----------------------------------------------------------------------------
// Widgets
//----------------------------------------------------------------------------

// Create new Input widget. This DOES NOT call termbox.Flush().
func Input(x, y, width int, fg, bg termbox.Attribute) *Editbox {
	ebox := newEditbox(x, y, width, 1, options{
		fg:   fg,
		bg:   bg,
		wrap: false,
		exitKeys: []termbox.Key{
			termbox.KeyEsc,
			termbox.KeyTab,
			termbox.KeyEnter,
		},
		autoexpand: false,
	})
	ebox.Render()
	return ebox
}

// Create new Textarea widget. This DOES NOT call termbox.Flush().
func Textarea(
	x, y, width, height int,
	fg, bg termbox.Attribute,
	wrap bool,
) *Editbox {
	ebox := newEditbox(x, y, width, height, options{
		fg:   fg,
		bg:   bg,
		wrap: wrap,
		exitKeys: []termbox.Key{
			termbox.KeyEsc,
			termbox.KeyTab,
		},
		autoexpand: false,
	})
	ebox.Render()
	return ebox
}
