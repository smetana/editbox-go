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

type Cursor struct {
	X, Y int
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
	cursor Cursor
	lastx  int
}

func newEditor() *editor {
	var ed editor
	ed.lines = make([]line, 1)
	ed.cursor.X = 0
	ed.cursor.Y = 0
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
	return &ed.lines[ed.cursor.Y]
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
	line.insertRune(cursor.X, r)
	cursor.X += 1
	if r == '\n' {
		ed.splitLine(cursor.X, cursor.Y)
		cursor.Y += 1
		cursor.X = 0
	}
	ed.lastx = cursor.X
}

func (ed *editor) checkYPosition(y int) {
	if y < 0 || y > len(ed.lines) {
		panic("y position out of range")
	}
}

func (ed *editor) deleteRuneBeforeCursor() {
	cursor := &ed.cursor
	if cursor.X == 0 && cursor.Y == 0 {
		return
	}
	ed.moveCursorLeft()
	ed.deleteRuneAtCursor()
}

func (ed *editor) deleteRuneAtCursor() {
	cursor := &ed.cursor
	l := ed.currentLine()
	r := l.deleteRune(cursor.X)
	if r == '\n' && cursor.Y < len(ed.lines)-1 {
		left := &ed.lines[cursor.Y]
		right := &ed.lines[cursor.Y+1]
		left.text = append(left.text, right.text...)
		if cursor.Y == len(ed.lines)-2 {
			ed.lines = ed.lines[:cursor.Y+1]
		} else {
			copy(ed.lines[cursor.Y+1:], ed.lines[cursor.Y+2:])
			ed.lines[len(ed.lines)-1] = *(new(line))
			ed.lines = ed.lines[:len(ed.lines)-1]
		}
	}
}

func (ed *editor) moveCursorRight() {
	cursor := &ed.cursor
	line := ed.currentLine()
	cursor.X += 1
	if cursor.X >= len(line.text) {
		if cursor.Y < len(ed.lines)-1 {
			cursor.Y += 1
			cursor.X = 0
		} else {
			cursor.X = len(line.text)
		}
	}
	ed.lastx = cursor.X
}

func (ed *editor) moveCursorLeft() {
	cursor := &ed.cursor
	cursor.X -= 1
	if cursor.X < 0 {
		if cursor.Y > 0 {
			cursor.Y -= 1
			line := ed.currentLine()
			cursor.X = len(line.text) - 1
		} else {
			cursor.X = 0
		}
	}
	ed.lastx = cursor.X
}

func (ed *editor) moveCursorToLineStart() {
	ed.cursor.X, ed.lastx = 0, 0
}

func (ed *editor) moveCursorToLineEnd() {
	line := ed.currentLine()
	if line.lastRune() == '\n' {
		ed.cursor.X = len(line.text) - 1
	} else {
		ed.cursor.X = len(line.text)
	}
	ed.lastx = ed.cursor.X
}

func (ed *editor) moveCursorVert(dy int) {
	cursor := &ed.cursor
	if cursor.Y+dy < 0 {
		return
	}
	if cursor.Y+dy > len(ed.lines)-1 {
		return
	}
	cursor.Y += dy
	line := ed.currentLine()
	switch {
	case len(line.text) == 0:
		cursor.X = 0
	case ed.lastx >= len(line.text):
		cursor.X = len(line.text) - 1
	default:
		cursor.X = ed.lastx
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

type Options struct {
	Fg         termbox.Attribute
	Bg         termbox.Attribute
	Wrap       bool
	autoexpand bool
	maxHeight  int
	printNL    bool
	exitKeys   []termbox.Key
}

type Editbox struct {
	editor        *editor
	Cursor        Cursor
	Scroll        Cursor
	X, Y          int
	Width, Height int
	Wrap          bool
	Fg, Bg        termbox.Attribute
	autoexpand    bool
	printNL       bool
	exitKeys      []termbox.Key
	view          [][]rune
	// Line y coord in box in Wrap mode
	lineBoxY      []int
	virtualHeight int
	minHeight     int
	maxHeight     int
}

func NewEditbox(x, y, width, height int, options Options) *Editbox {
	var ebox Editbox
	ebox.X = x
	ebox.Y = y
	ebox.Width = width
	ebox.Height = height
	ebox.Fg = options.Fg
	ebox.Bg = options.Bg
	ebox.Wrap = options.Wrap
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

func (ebox *Editbox) Text() string {
	return ebox.editor.text()
}

func (ebox *Editbox) updateLineOffsets() {
	ed := ebox.editor
	linesCnt := len(ed.lines)
	ebox.lineBoxY = make([]int, linesCnt)
	dy := 0 // delta between editor y and box Y
	cumulativeOffset := 0
	for y := 0; y < linesCnt; y++ {
		ebox.lineBoxY[y] = y + cumulativeOffset
		if ebox.Wrap {
			dy = (len(ed.lines[y].text) - 1) / ebox.Width
			cumulativeOffset += dy
		}
	}
	ebox.virtualHeight = ebox.lineBoxY[linesCnt-1] + dy + 1
	if ebox.autoexpand {
		if ebox.virtualHeight > ebox.Height {
			if ebox.virtualHeight > ebox.maxHeight {
				ebox.Height = ebox.maxHeight
			} else {
				ebox.Height = ebox.virtualHeight
			}
		} else if ebox.virtualHeight < ebox.Height {
			if ebox.virtualHeight < ebox.minHeight {
				ebox.Height = ebox.minHeight
			} else {
				ebox.Height = ebox.virtualHeight
			}
		}
		// else Ok. Don't change height
	}
	// else Ok. don't change height
	ebox.Cursor.X, ebox.Cursor.Y = ebox.editorToBox(ed.cursor.X, ed.cursor.Y)
}

func (ebox *Editbox) editorToBox(x, y int) (int, int) {
	if ebox.Wrap {
		ldy := x / ebox.Width
		x = x - (ldy * ebox.Width)
		y = ebox.lineBoxY[y] + ldy
	}
	return x, y
}

func (ebox *Editbox) moveCursorLeft() {
	ebox.editor.moveCursorLeft()
}

func (ebox *Editbox) moveCursorRight() {
	ebox.editor.moveCursorRight()
}

func (ebox *Editbox) moveCursorToLineStart() {
	ebox.editor.moveCursorToLineStart()
}

func (ebox *Editbox) moveCursorToLineEnd() {
	ebox.editor.moveCursorToLineEnd()
}

// Cursor movement in Wrap mode is a bit tricky
// TODO Code smell. Refactor
func (ebox *Editbox) moveCursorDown() {
	if ebox.Wrap {
		ed := ebox.editor
		line := ed.currentLine()
		// Try to move within current line
		if ed.cursor.X+ebox.Width < len(line.text) {
			ed.cursor.X += ebox.Width
			return
		}
		if ebox.Cursor.X+(len(line.text)-ed.cursor.X)-1 >= ebox.Width {
			ed.cursor.X = line.lastRuneX()
			return
		}
		// Jump to next line
		if ed.cursor.Y+1 > len(ed.lines)-1 {
			return
		}
		ed.cursor.Y += 1
		line = ed.currentLine()
		if len(line.text) == 0 {
			ed.cursor.X = 0
			return
		}
		x, _ := ebox.editorToBox(ed.lastx, 0)
		if x >= len(line.text) {
			ed.cursor.X = line.lastRuneX()
		} else {
			ed.cursor.X = x
		}
	} else {
		ebox.editor.moveCursorVert(+1)
	}
}

func (ebox *Editbox) moveCursorUp() {
	if ebox.Wrap {
		ed := ebox.editor
		lastx, _ := ebox.editorToBox(ed.lastx, 0)
		x, _ := ebox.editorToBox(ed.cursor.X, 0)
		if x == lastx && ed.cursor.X-ebox.Width >= 0 {
			ed.cursor.X -= ebox.Width
			return
		}
		d := ebox.Width + x - lastx
		if x < lastx && ed.cursor.X-d >= 0 {
			ed.cursor.X -= d
			return
		}
		if ed.cursor.Y-1 < 0 {
			return
		}
		ed.cursor.Y -= 1
		line := ed.currentLine()
		if ed.lastx < ebox.Width {
			ed.cursor.X = ed.lastx
		}
		if lastx >= line.lastRuneX() {
			ed.cursor.X = line.lastRuneX()
			return
		}
		x, _ = ebox.editorToBox(line.lastRuneX(), 0)
		if x <= lastx {
			ed.cursor.X = line.lastRuneX()
		} else {
			ed.cursor.X = line.lastRuneX() - x + lastx
		}
	} else {
		ebox.editor.moveCursorVert(-1)
	}
}

func (ebox *Editbox) moveCursorPageUp() {
	for i := 1; i <= ebox.Height; i++ {
		ebox.moveCursorUp()
	}
}

func (ebox *Editbox) moveCursorPageDown() {
	for i := 1; i <= ebox.Height; i++ {
		ebox.moveCursorDown()
	}
}

func (ebox *Editbox) scrollToCursor() {
	if !ebox.Wrap {
		if ebox.Cursor.X-ebox.Scroll.X > ebox.Width-1 {
			ebox.Scroll.X = ebox.Cursor.X - ebox.Width + 1
		} else if ebox.Cursor.X-ebox.Scroll.X < 0 {
			ebox.Scroll.X = ebox.Cursor.X
		}
	}
	if ebox.virtualHeight > ebox.Height {
		if ebox.Cursor.Y-ebox.Scroll.Y > ebox.Height-1 {
			ebox.Scroll.Y = ebox.Cursor.Y - ebox.Height + 1
		} else if ebox.Cursor.Y-ebox.Scroll.Y < 0 {
			ebox.Scroll.Y = ebox.Cursor.Y
		} else if ebox.virtualHeight-ebox.Scroll.Y <= ebox.Height-1 {
			ebox.Scroll.Y = ebox.virtualHeight - ebox.Height
		}
	} else {
		ebox.Scroll.Y = 0
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
	ebox.view = make([][]rune, ebox.Height)
	for i := range ebox.view {
		ebox.view[i] = make([]rune, ebox.Width)
	}
	for y, line := range ed.lines {
		for x, r := range line.text {
			boxX, boxY = ebox.editorToBox(x, y)
			//TODO Optimize
			if boxY < ebox.Scroll.Y || boxX < ebox.Scroll.X {
				continue
			}
			viewX = boxX - ebox.Scroll.X
			viewY = boxY - ebox.Scroll.Y
			if viewX > ebox.Width-1 {
				break
			}
			if viewY > ebox.Height-1 {
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
		if viewY > ebox.Height-1 {
			break
		}
	}
}

//----------------------------------------------------------------------------
// API
//----------------------------------------------------------------------------

func (ebox *Editbox) Render() {
	ebox.renderView()
	var r rune
	for y := 0; y < ebox.Height; y++ {
		for x := 0; x < ebox.Width; x++ {
			if ebox.view[y][x] != 0 {
				r = ebox.view[y][x]
			} else {
				r = ' ' // Fill empty cells with background color
			}
			termbox.SetCell(ebox.X+x, ebox.Y+y, r, ebox.Fg, ebox.Bg)
		}
	}
	termbox.SetCursor(ebox.X+ebox.Cursor.X-ebox.Scroll.X,
		ebox.Y+ebox.Cursor.Y-ebox.Scroll.Y)
}

func (ebox *Editbox) HandleEvent(ev *termbox.Event) {
	ed := ebox.editor
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowLeft:
			ebox.moveCursorLeft()
		case termbox.KeyArrowRight:
			ebox.moveCursorRight()
		case termbox.KeyArrowUp:
			ebox.moveCursorUp()
		case termbox.KeyArrowDown:
			ebox.moveCursorDown()
		case termbox.KeyHome:
			ebox.moveCursorToLineStart()
		case termbox.KeyEnd:
			ebox.moveCursorToLineEnd()
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

func (ebox *Editbox) WaitExit() termbox.Event {
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

func NewInputbox(x, y, width int, fg, bg termbox.Attribute) *Editbox {
	ebox := NewEditbox(x, y, width, 1, Options{
		Fg:         fg,
		Bg:         bg,
		Wrap:       false,
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
