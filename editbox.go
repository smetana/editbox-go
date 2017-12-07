package main

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
	x, y int
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
	ed.cursor.x = 0
	ed.cursor.y = 0
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
	return &ed.lines[ed.cursor.y]
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
	line.insertRune(cursor.x, r)
	cursor.x += 1
	if r == '\n' {
		ed.splitLine(cursor.x, cursor.y)
		cursor.y += 1
		cursor.x = 0
	}
	ed.lastx = cursor.x
}

func (ed *Editor) checkYPosition(y int) {
	if y < 0 || y > len(ed.lines) {
		panic("y position out of range")
	}
}

func (ed *Editor) deleteRuneBeforeCursor() {
	cursor := &ed.cursor
	if cursor.x == 0 && cursor.y == 0 {
		return
	}
	ed.moveCursorLeft()
	ed.deleteRuneAtCursor()
}

func (ed *Editor) deleteRuneAtCursor() {
	cursor := &ed.cursor
	line := ed.currentLine()
	r := line.deleteRune(cursor.x)
	if r == '\n' && cursor.y < len(ed.lines)-1 {
		left := &ed.lines[cursor.y]
		right := &ed.lines[cursor.y+1]
		left.text = append(left.text, right.text...)
		if cursor.y == len(ed.lines)-2 {
			ed.lines = ed.lines[:cursor.y+1]
		} else {
			copy(ed.lines[cursor.y+1:], ed.lines[cursor.y+2:])
			ed.lines[len(ed.lines)-1] = *(new(Line))
			ed.lines = ed.lines[:len(ed.lines)-1]
		}
	}
}

func (ed *Editor) moveCursorRight() {
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

func (ed *Editor) moveCursorLeft() {
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

func (ed *Editor) moveCursorToLineStart() {
	ed.cursor.x, ed.lastx = 0, 0
}

func (ed *Editor) moveCursorToLineEnd() {
	line := ed.currentLine()
	if line.lastRune() == '\n' {
		ed.cursor.x = len(line.text) - 1
	} else {
		ed.cursor.x = len(line.text)
	}
	ed.lastx = ed.cursor.x
}

func (ed *Editor) moveCursorVert(dy int) {
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

//----------------------------------------------------------------------------
// Editbox
//----------------------------------------------------------------------------

type Options struct {
	fg         termbox.Attribute
	bg         termbox.Attribute
	wrap       bool
	autoexpand bool
	maxHeight  int
}

type Editbox struct {
	editor        *Editor
	cursor        Cursor
	width, height int
	wrap          bool
	autoexpand    bool
	fg, bg        termbox.Attribute
	// Line y coord in box in wrap mode
	lineBoxY      []int
	visibleHeight int
	virtualHeight int
	minHeight, maxHeight int
	scroll        Cursor
}

func NewEditbox(width, height int, options Options) *Editbox {
	var ebox Editbox
	ebox.width = width
	ebox.height = height
	ebox.fg = options.fg
	ebox.bg = options.bg
	ebox.editor = NewEditor()
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
		switch {
		case ebox.virtualHeight > ebox.height:
			if ebox.virtualHeight > ebox.maxHeight {
				ebox.visibleHeight = ebox.maxHeight
			} else {
				ebox.visibleHeight = ebox.virtualHeight
			}
			ebox.height = ebox.visibleHeight
		case ebox.virtualHeight < ebox.height:
			if ebox.virtualHeight < ebox.minHeight {
				ebox.visibleHeight = ebox.minHeight
			} else {
				ebox.visibleHeight = ebox.virtualHeight
			}
			ebox.height = ebox.visibleHeight
		default:
			ebox.visibleHeight = ebox.height
		}
	} else {
		ebox.visibleHeight = ebox.height
	}
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
	for i:=1; i <= ebox.height; i++ {
		ebox.moveCursorUp()
	}
}

func (ebox *Editbox) moveCursorPageDown() {
	for i:=1; i <= ebox.height; i++ {
		ebox.moveCursorDown()
	}
}

func (ebox *Editbox) scrollToCursor() {
	if !ebox.wrap {
		if ebox.cursor.x - ebox.scroll.x > ebox.width - 1 {
			ebox.scroll.x = ebox.cursor.x - ebox.width + 1
		} else if ebox.cursor.x - ebox.scroll.x < 0 {
			ebox.scroll.x = ebox.cursor.x
		}
	}
	if ebox.virtualHeight > ebox.height {
		if ebox.cursor.y-ebox.scroll.y > ebox.height-1 {
			ebox.scroll.y = ebox.cursor.y - ebox.height + 1
		} else if ebox.cursor.y-ebox.scroll.y < 0 {
			ebox.scroll.y = ebox.cursor.y
		} else if ebox.virtualHeight - ebox.scroll.y <= ebox.height - 1 {
			ebox.scroll.y = ebox.virtualHeight - ebox.height
		}
	} else {
		ebox.scroll.y = 0
	}
}

func (ebox *Editbox) Draw() {
	ebox.updateLineOffsets()
	ebox.scrollToCursor()
	ed := ebox.editor
	coldef := termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	var (
		x, y          int
		boxX, boxY    int
		viewX, viewY  int
	)
	// Fill background. TODO Optimize with next for
	for y = 0; y < ebox.visibleHeight; y++ {
		for x = 0; x < ebox.width; x++ {
			termbox.SetCell(x, y, ' ', ebox.fg, ebox.bg)
		}
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
			if viewX > ebox.width - 1 {
				break
			}
			if viewY > ebox.height-1 {
				break
			}
			// TODO Remove debug ???
			if r == '\n' {
				r = '␤'
			}
			termbox.SetCell(viewX, viewY, r, ebox.fg, ebox.bg)
		}
		if viewY > ebox.height-1 {
			break
		}
	}
	ebox.indicateScrolling()
	termbox.SetCursor(ebox.cursor.x - ebox.scroll.x, ebox.cursor.y - ebox.scroll.y)
	termbox.Flush()
}

// TODO Better solution
func (ebox *Editbox) indicateScrolling() {
	if ebox.scroll.y > 0 {
		if ebox.cursor.x != 0 || ebox.cursor.y > ebox.scroll.y {
			termbox.SetCell(0, 0, '↑', ebox.fg, ebox.bg)
		}
		if (ebox.cursor.x != ebox.width - 1) || ebox.cursor.y > ebox.scroll.y {
			termbox.SetCell(ebox.width - 1, 0, '↑', ebox.fg, ebox.bg)
		}
	}
	if ebox.virtualHeight > ebox.visibleHeight + ebox.scroll.y {
		if ebox.cursor.x != 0 || ebox.cursor.y < ebox.scroll.y + ebox.visibleHeight {
			termbox.SetCell(0, ebox.height - 1, '↓', ebox.fg, ebox.bg)
		}
		if (ebox.cursor.x != ebox.width - 1) || ebox.cursor.y < ebox.scroll.y + ebox.visibleHeight {
			termbox.SetCell(ebox.width - 1, ebox.height - 1, '↓', ebox.fg, ebox.bg)
		}
	}
}

func (ebox *Editbox) handleEvent(ev *termbox.Event) bool {
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

//----------------------------------------------------------------------------
// main() and support
//----------------------------------------------------------------------------

func mainLoop(ebox *Editbox) {
	eventQueue := make(chan termbox.Event, 256)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	for {
		select {
		case ev := <-eventQueue:
			ok := ebox.handleEvent(&ev)
			if !ok {
				return
			}
			if len(eventQueue) == 0 {
				ebox.Draw()
			}
		}
	}
}

func main() {
	err := termbox.Init()
	check(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)
	ebox := NewEditbox(20, 3, Options{
		wrap:       true,
		autoexpand: true,
		maxHeight:  6,
		fg:         12,
		bg:         63})
	ebox.Draw()
	mainLoop(ebox)
}
