package editbox

import (
	"bytes"
)

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

// TODO Optimize
func (ed *editor) setText(text string) {
	for _, s := range text {
		ed.insertRune(rune(s))
	}
}
