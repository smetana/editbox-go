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
// Line
//----------------------------------------------------------------------------

type Line struct {
	text []rune
}

func (l *Line) checkXPosition(x int) {
	if x < 0 || x > len(l.text) {
		panic("x position out of range")
	}
}

func (l *Line) insertRune(pos int, r rune) {
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

func (l *Line) split(pos int) (left, right *Line) {
	l.checkXPosition(pos)
	left, right = l, new(Line)
	right.text = make([]rune, len(l.text)-pos)
	copy(right.text, l.text[pos:len(l.text)])
	left.text = left.text[:pos]
	return
}

func (l *Line) deleteRune(pos int) rune {
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

func (l *Line) lastRune() rune {
	return l.text[len(l.text)-1]
}

func (l *Line) lastRuneX() int {
	if l.lastRune() == '\n' {
		return (len(l.text) - 1)
	} else {
		return (len(l.text))
	}
}

//----------------------------------------------------------------------------
// Editor
//----------------------------------------------------------------------------

type Editor struct {
	lines  []Line
	cursor Cursor
	lastx  int
}

func NewEditor() *Editor {
	var ed Editor
	ed.lines = make([]Line, 1)
	ed.cursor.X = 0
	ed.cursor.Y = 0
	return &ed
}

func (ed *Editor) Text() string {
	var b bytes.Buffer
	for _, l := range ed.lines {
		b.WriteString(string(l.text))
	}
	return b.String()
}

func (ed *Editor) currentLine() *Line {
	return &ed.lines[ed.cursor.Y]
}

func (ed *Editor) splitLine(x, y int) {
	line := &ed.lines[y]
	left, right := line.split(x)
	ed.lines = append(ed.lines, *(new(Line)))
	copy(ed.lines[y+2:], ed.lines[y+1:])
	ed.lines[y] = *left
	ed.lines[y+1] = *right
}

func (ed *Editor) insertRune(r rune) {
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

func (ed *Editor) checkYPosition(y int) {
	if y < 0 || y > len(ed.lines) {
		panic("y position out of range")
	}
}

func (ed *Editor) deleteRuneBeforeCursor() {
	cursor := &ed.cursor
	if cursor.X == 0 && cursor.Y == 0 {
		return
	}
	ed.moveCursorLeft()
	ed.deleteRuneAtCursor()
}

func (ed *Editor) deleteRuneAtCursor() {
	cursor := &ed.cursor
	line := ed.currentLine()
	r := line.deleteRune(cursor.X)
	if r == '\n' && cursor.Y < len(ed.lines)-1 {
		left := &ed.lines[cursor.Y]
		right := &ed.lines[cursor.Y+1]
		left.text = append(left.text, right.text...)
		if cursor.Y == len(ed.lines)-2 {
			ed.lines = ed.lines[:cursor.Y+1]
		} else {
			copy(ed.lines[cursor.Y+1:], ed.lines[cursor.Y+2:])
			ed.lines[len(ed.lines)-1] = *(new(Line))
			ed.lines = ed.lines[:len(ed.lines)-1]
		}
	}
}

func (ed *Editor) moveCursorRight() {
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

func (ed *Editor) moveCursorLeft() {
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

func (ed *Editor) moveCursorToLineStart() {
	ed.cursor.X, ed.lastx = 0, 0
}

func (ed *Editor) moveCursorToLineEnd() {
	line := ed.currentLine()
	if line.lastRune() == '\n' {
		ed.cursor.X = len(line.text) - 1
	} else {
		ed.cursor.X = len(line.text)
	}
	ed.lastx = ed.cursor.X
}

func (ed *Editor) moveCursorVert(dy int) {
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

//----------------------------------------------------------------------------
// Editbox
//----------------------------------------------------------------------------

type Options struct {
	Fg         termbox.Attribute
	Bg         termbox.Attribute
	Wrap       bool
	Autoexpand bool
	MaxHeight  int
	PrintNL    bool
}

type Editbox struct {
	editor        *Editor
	view          [][]rune
	cursor        Cursor
	x, y          int
	width, height int
	wrap          bool
	autoexpand    bool
	fg, bg        termbox.Attribute
	printNL       bool
	// Line y coord in box in wrap mode
	lineBoxY      []int
	virtualHeight int
	minHeight     int
	maxHeight     int
	scroll        Cursor
}

func NewEditbox(x, y, width, height int, options Options) *Editbox {
	var ebox Editbox
	ebox.x = x
	ebox.y = y
	ebox.width = width
	ebox.height = height
	ebox.fg = options.Fg
	ebox.bg = options.Bg
	ebox.editor = NewEditor()
	ebox.wrap = options.Wrap
	ebox.autoexpand = options.Autoexpand
	if ebox.autoexpand {
		ebox.minHeight = height
		if options.MaxHeight <= 0 {
			ebox.maxHeight = ebox.minHeight
		} else {
			ebox.maxHeight = options.MaxHeight
		}
	}
	ebox.printNL = options.PrintNL
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
	ebox.cursor.X, ebox.cursor.Y = ebox.editorToBox(ed.cursor.X, ed.cursor.Y)
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
		if ed.cursor.X+ebox.width < len(line.text) {
			ed.cursor.X += ebox.width
			return
		}
		if ebox.cursor.X+(len(line.text)-ed.cursor.X)-1 >= ebox.width {
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
	if ebox.wrap {
		ed := ebox.editor
		lastx, _ := ebox.editorToBox(ed.lastx, 0)
		x, _ := ebox.editorToBox(ed.cursor.X, 0)
		if x == lastx && ed.cursor.X-ebox.width >= 0 {
			ed.cursor.X -= ebox.width
			return
		}
		d := ebox.width + x - lastx
		if x < lastx && ed.cursor.X-d >= 0 {
			ed.cursor.X -= d
			return
		}
		if ed.cursor.Y-1 < 0 {
			return
		}
		ed.cursor.Y -= 1
		line := ed.currentLine()
		if ed.lastx < ebox.width {
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
		if ebox.cursor.X-ebox.scroll.X > ebox.width-1 {
			ebox.scroll.X = ebox.cursor.X - ebox.width + 1
		} else if ebox.cursor.X-ebox.scroll.X < 0 {
			ebox.scroll.X = ebox.cursor.X
		}
	}
	if ebox.virtualHeight > ebox.height {
		if ebox.cursor.Y-ebox.scroll.Y > ebox.height-1 {
			ebox.scroll.Y = ebox.cursor.Y - ebox.height + 1
		} else if ebox.cursor.Y-ebox.scroll.Y < 0 {
			ebox.scroll.Y = ebox.cursor.Y
		} else if ebox.virtualHeight-ebox.scroll.Y <= ebox.height-1 {
			ebox.scroll.Y = ebox.virtualHeight - ebox.height
		}
	} else {
		ebox.scroll.Y = 0
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
			if boxY < ebox.scroll.Y || boxX < ebox.scroll.X {
				continue
			}
			viewX = boxX - ebox.scroll.X
			viewY = boxY - ebox.scroll.Y
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

func (ebox *Editbox) Draw() {
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
	termbox.SetCursor(ebox.x+ebox.cursor.X-ebox.scroll.X,
		ebox.y+ebox.cursor.Y-ebox.scroll.Y)
}

func (ebox *Editbox) HandleEvent(ev *termbox.Event) bool {
	ed := ebox.editor
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEsc:
			// Quit
			return false
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
	return true
}
