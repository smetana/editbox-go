package main

import (
	"github.com/nsf/termbox-go"
    "bytes"
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
    right.text = make([]rune, len(l.text) - pos)
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

//----------------------------------------------------------------------------
// Editor
//----------------------------------------------------------------------------

type Editor struct {
    width, height int
    lines []Line
    cursor Cursor
    lastx int
}

func NewEditor(width, height int) *Editor {
    var ed Editor
    ed.width = width
    ed.height = height
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

func (ed *Editor) insertRune(r rune) {
    cursor := &ed.cursor
    line := ed.currentLine()
    line.insertRune(cursor.x, r)
    cursor.x += 1
    if r == '\n' {
        left, right := line.split(cursor.x)
        ed.lines = append(ed.lines, *(new(Line)))
        copy(ed.lines[cursor.y+2:], ed.lines[cursor.y+1:])
        ed.lines[cursor.y] = *left
        ed.lines[cursor.y+1] = *right
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

// TODO Better name
func (ed *Editor) concatNextLine() {
	cursor := &ed.cursor
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
	if r == '\n' && cursor.y < len(ed.lines) - 1 {
		ed.concatNextLine()
	}
}

func (ed *Editor) moveCursorRight() {
    cursor := &ed.cursor
    line := ed.currentLine()
    cursor.x += 1
    if cursor.x >= len(line.text) {
        if cursor.y < len(ed.lines) - 1 {
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

func (ed *Editor) moveCursorVert(dy int) {
    cursor := &ed.cursor
    if cursor.y + dy < 0 { return }
    if cursor.y + dy > len(ed.lines) - 1 { return }
    cursor.y += dy
    line := ed.currentLine()
    if ed.lastx >= len(line.text) {
        cursor.x = len(line.text) - 1
    } else {
        cursor.x = ed.lastx
    }
}


func (ed *Editor) Draw() {
    coldef := termbox.ColorDefault
    termbox.Clear(coldef, coldef);
    for y, line := range ed.lines {
        for x, r := range line.text {
            if r == '\n' { r = '$' } // TODO remove debug
            termbox.SetCell(x, y, r, coldef, coldef)
        }
    }
    termbox.SetCursor(ed.cursor.x, ed.cursor.y)
    termbox.Flush()
}


func main() {
    err := termbox.Init()
    check(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
    ed := NewEditor(10, 10)
    ed.Draw()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowLeft:
                 ed.moveCursorLeft()
			case termbox.KeyArrowRight:
                 ed.moveCursorRight()
			case termbox.KeyArrowUp:
                 ed.moveCursorVert(-1)
			case termbox.KeyArrowDown:
                 ed.moveCursorVert(+1)
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
        ed.Draw()
    }
}
