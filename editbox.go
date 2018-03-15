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
	Text []rune
}

func (l *Line) checkXPosition(x int) {
	if x < 0 || x > len(l.Text) {
		panic("x position out of range")
	}
}

func (l *Line) InsertRune(pos int, r rune) {
	l.checkXPosition(pos)
	// Append
	if pos == len(l.Text) {
		l.Text = append(l.Text, r)
		// Insert
	} else {
		l.Text = append(l.Text, rune(0))
		copy(l.Text[pos+1:], l.Text[pos:])
		l.Text[pos] = r
	}
}

func (l *Line) Split(pos int) (left, right *Line) {
	l.checkXPosition(pos)
	left, right = l, new(Line)
	right.Text = make([]rune, len(l.Text)-pos)
	copy(right.Text, l.Text[pos:len(l.Text)])
	left.Text = left.Text[:pos]
	return
}

func (l *Line) DeleteRune(pos int) rune {
	l.checkXPosition(pos)
	if pos < len(l.Text) {
		r := l.Text[pos]
		copy(l.Text[pos:], l.Text[pos+1:])
		l.Text[len(l.Text)-1] = rune(0)
		l.Text = l.Text[:len(l.Text)-1]
		return r
	} else {
		return rune(0)
	}
}

func (l *Line) lastRune() rune {
	return l.Text[len(l.Text)-1]
}

func (l *Line) lastRuneX() int {
	if l.lastRune() == '\n' {
		return (len(l.Text) - 1)
	} else {
		return (len(l.Text))
	}
}

//----------------------------------------------------------------------------
// Editor
//----------------------------------------------------------------------------

type Editor struct {
	Lines  []Line
	Cursor Cursor
	lastx  int
}

func NewEditor() *Editor {
	var ed Editor
	ed.Lines = make([]Line, 1)
	ed.Cursor.X = 0
	ed.Cursor.Y = 0
	return &ed
}

func (ed *Editor) Text() string {
	var b bytes.Buffer
	for _, l := range ed.Lines {
		b.WriteString(string(l.Text))
	}
	return b.String()
}

func (ed *Editor) CurrentLine() *Line {
	return &ed.Lines[ed.Cursor.Y]
}

func (ed *Editor) splitLine(x, y int) {
	line := &ed.Lines[y]
	left, right := line.Split(x)
	ed.Lines = append(ed.Lines, *(new(Line)))
	copy(ed.Lines[y+2:], ed.Lines[y+1:])
	ed.Lines[y] = *left
	ed.Lines[y+1] = *right
}

func (ed *Editor) InsertRune(r rune) {
	cursor := &ed.Cursor
	line := ed.CurrentLine()
	line.InsertRune(cursor.X, r)
	cursor.X += 1
	if r == '\n' {
		ed.splitLine(cursor.X, cursor.Y)
		cursor.Y += 1
		cursor.X = 0
	}
	ed.lastx = cursor.X
}

func (ed *Editor) checkYPosition(y int) {
	if y < 0 || y > len(ed.Lines) {
		panic("y position out of range")
	}
}

func (ed *Editor) DeleteRuneBeforeCursor() {
	cursor := &ed.Cursor
	if cursor.X == 0 && cursor.Y == 0 {
		return
	}
	ed.moveCursorLeft()
	ed.DeleteRuneAtCursor()
}

func (ed *Editor) DeleteRuneAtCursor() {
	cursor := &ed.Cursor
	line := ed.CurrentLine()
	r := line.DeleteRune(cursor.X)
	if r == '\n' && cursor.Y < len(ed.Lines)-1 {
		left := &ed.Lines[cursor.Y]
		right := &ed.Lines[cursor.Y+1]
		left.Text = append(left.Text, right.Text...)
		if cursor.Y == len(ed.Lines)-2 {
			ed.Lines = ed.Lines[:cursor.Y+1]
		} else {
			copy(ed.Lines[cursor.Y+1:], ed.Lines[cursor.Y+2:])
			ed.Lines[len(ed.Lines)-1] = *(new(Line))
			ed.Lines = ed.Lines[:len(ed.Lines)-1]
		}
	}
}

func (ed *Editor) moveCursorRight() {
	cursor := &ed.Cursor
	line := ed.CurrentLine()
	cursor.X += 1
	if cursor.X >= len(line.Text) {
		if cursor.Y < len(ed.Lines)-1 {
			cursor.Y += 1
			cursor.X = 0
		} else {
			cursor.X = len(line.Text)
		}
	}
	ed.lastx = cursor.X
}

func (ed *Editor) moveCursorLeft() {
	cursor := &ed.Cursor
	cursor.X -= 1
	if cursor.X < 0 {
		if cursor.Y > 0 {
			cursor.Y -= 1
			line := ed.CurrentLine()
			cursor.X = len(line.Text) - 1
		} else {
			cursor.X = 0
		}
	}
	ed.lastx = cursor.X
}

func (ed *Editor) moveCursorToLineStart() {
	ed.Cursor.X, ed.lastx = 0, 0
}

func (ed *Editor) moveCursorToLineEnd() {
	line := ed.CurrentLine()
	if line.lastRune() == '\n' {
		ed.Cursor.X = len(line.Text) - 1
	} else {
		ed.Cursor.X = len(line.Text)
	}
	ed.lastx = ed.Cursor.X
}

func (ed *Editor) moveCursorVert(dy int) {
	cursor := &ed.Cursor
	if cursor.Y+dy < 0 {
		return
	}
	if cursor.Y+dy > len(ed.Lines)-1 {
		return
	}
	cursor.Y += dy
	line := ed.CurrentLine()
	switch {
	case len(line.Text) == 0:
		cursor.X = 0
	case ed.lastx >= len(line.Text):
		cursor.X = len(line.Text) - 1
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
	Editor        *Editor
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
	ebox.Editor = NewEditor()
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
	ed := ebox.Editor
	linesCnt := len(ed.Lines)
	ebox.lineBoxY = make([]int, linesCnt)
	dy := 0 // delta between editor y and box Y
	cumulativeOffset := 0
	for y := 0; y < linesCnt; y++ {
		ebox.lineBoxY[y] = y + cumulativeOffset
		if ebox.wrap {
			dy = (len(ed.Lines[y].Text) - 1) / ebox.width
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
	ebox.cursor.X, ebox.cursor.Y = ebox.editorToBox(ed.Cursor.X, ed.Cursor.Y)
}

func (ebox *Editbox) editorToBox(x, y int) (int, int) {
	if ebox.wrap {
		ldy := x / ebox.width
		x = x - (ldy * ebox.width)
		y = ebox.lineBoxY[y] + ldy
	}
	return x, y
}


func (ebox *Editbox) moveCursorLeft() {
	ebox.Editor.moveCursorLeft()
}

func (ebox *Editbox) moveCursorRight() {
	ebox.Editor.moveCursorRight()
}

func (ebox *Editbox) moveCursorToLineStart() {
	ebox.Editor.moveCursorToLineStart()
}

func (ebox *Editbox) moveCursorToLineEnd() {
	ebox.Editor.moveCursorToLineEnd()
}

// Cursor movement in wrap mode is a bit tricky
// TODO Code smell. Refactor
func (ebox *Editbox) moveCursorDown() {
	if ebox.wrap {
		ed := ebox.Editor
		line := ed.CurrentLine()
		// Try to move within current line
		if ed.Cursor.X+ebox.width < len(line.Text) {
			ed.Cursor.X += ebox.width
			return
		}
		if ebox.cursor.X+(len(line.Text)-ed.Cursor.X)-1 >= ebox.width {
			ed.Cursor.X = line.lastRuneX()
			return
		}
		// Jump to next line
		if ed.Cursor.Y+1 > len(ed.Lines)-1 {
			return
		}
		ed.Cursor.Y += 1
		line = ed.CurrentLine()
		if len(line.Text) == 0 {
			ed.Cursor.X = 0
			return
		}
		x, _ := ebox.editorToBox(ed.lastx, 0)
		if x >= len(line.Text) {
			ed.Cursor.X = line.lastRuneX()
		} else {
			ed.Cursor.X = x
		}
	} else {
		ebox.Editor.moveCursorVert(+1)
	}
}

func (ebox *Editbox) moveCursorUp() {
	if ebox.wrap {
		ed := ebox.Editor
		lastx, _ := ebox.editorToBox(ed.lastx, 0)
		x, _ := ebox.editorToBox(ed.Cursor.X, 0)
		if x == lastx && ed.Cursor.X-ebox.width >= 0 {
			ed.Cursor.X -= ebox.width
			return
		}
		d := ebox.width + x - lastx
		if x < lastx && ed.Cursor.X-d >= 0 {
			ed.Cursor.X -= d
			return
		}
		if ed.Cursor.Y-1 < 0 {
			return
		}
		ed.Cursor.Y -= 1
		line := ed.CurrentLine()
		if ed.lastx < ebox.width {
			ed.Cursor.X = ed.lastx
		}
		if lastx >= line.lastRuneX() {
			ed.Cursor.X = line.lastRuneX()
			return
		}
		x, _ = ebox.editorToBox(line.lastRuneX(), 0)
		if x <= lastx {
			ed.Cursor.X = line.lastRuneX()
		} else {
			ed.Cursor.X = line.lastRuneX() - x + lastx
		}
	} else {
		ebox.Editor.moveCursorVert(-1)
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
	ed := ebox.Editor
	var (
		boxX, boxY   int
		viewX, viewY int
	)
	ebox.view = make([][]rune, ebox.height)
	for i := range ebox.view {
		ebox.view[i] = make([]rune, ebox.width)
	}
	for y, line := range ed.Lines {
		for x, r := range line.Text {
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
	ed := ebox.Editor
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEsc:
			// Quit
			return false
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
			ed.DeleteRuneBeforeCursor()
		case termbox.KeyDelete:
			ed.DeleteRuneAtCursor()
		case termbox.KeyEnter:
			ed.InsertRune('\n')
		case termbox.KeySpace:
			ed.InsertRune(' ')
		default:
			if ev.Ch != 0 {
				ed.InsertRune(ev.Ch)
			}
		}
	case termbox.EventError:
		panic(ev.Err)
	default:
		// TODO
	}
	return true
}
