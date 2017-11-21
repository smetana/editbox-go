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

func (l *Line) lastRune() rune {
    return l.text[len(l.text) - 1]
}

//----------------------------------------------------------------------------
// Editor
//----------------------------------------------------------------------------

type Editor struct {
    lines []Line
    cursor Cursor
    lastx int
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
	if r == '\n' && cursor.y < len(ed.lines) - 1 {
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

func (ed *Editor) moveCursorToLineStart() {
    ed.cursor.x, ed.lastx = 0,0
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
    if cursor.y + dy < 0 { return }
    if cursor.y + dy > len(ed.lines) - 1 { return }
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
    fg termbox.Attribute
    bg termbox.Attribute
    wrap bool
    autoexpand bool
}

type Editbox struct {
    editor *Editor
    cursor Cursor
    width, height int
    wrap bool
    autoexpand bool
    fg, bg termbox.Attribute
    // Line y coord in box in wrap mode
    lineBoxY []int
    virtualHeight int
    scroll int
    // Needed to calculate cursor movement direction for scrolling
    prevCursor Cursor
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
    ebox.scroll = 0
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
    ebox.prevCursor = ebox.cursor
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

func (ebox *Editbox) scrollToCursor() {
    if ebox.cursor.y - ebox.scroll > ebox.height - 1 {
        ebox.scroll = ebox.cursor.y - ebox.height + 1
    } else if ebox.cursor.y - ebox.scroll < 0 {
        ebox.scroll = ebox.cursor.y
    }
}

func (ebox *Editbox) Draw() {
    ebox.updateLineOffsets()
    if !ebox.autoexpand {
        ebox.scrollToCursor()
    }
    ed := ebox.editor
    coldef := termbox.ColorDefault
    termbox.Clear(coldef, coldef);
    var (
        x, y int
        boxX, boxY int
        viewX, viewY int
        visibleHeight int
    )
    if ebox.autoexpand {
        if ebox.virtualHeight > ebox.height {
            visibleHeight = ebox.virtualHeight
        } else {
            visibleHeight = ebox.height
        }
    } else {
        visibleHeight = ebox.height
    }
    // Fill background. TODO Optimize with next for
    for y = 0; y < visibleHeight; y++ {
        for x = 0; x < ebox.width; x++ {
	        termbox.SetCell(x, y, ' ', ebox.fg, ebox.bg)
        }
    }
    for y, line := range ed.lines {
        for x, r := range line.text {
            boxX, boxY = ebox.editorToBox(x, y)
            //TODO Optimize
            if boxY < ebox.scroll {
                continue
            }
            viewX = boxX
            viewY = boxY - ebox.scroll
            if viewY > ebox.height - 1 && !ebox.autoexpand {
                break
            }
		    // TODO Remove debug ???
            if r == '\n' { r = '␤'}
	        termbox.SetCell(viewX, viewY, r, ebox.fg, ebox.bg)
        }
    }
    termbox.SetCursor(ebox.cursor.x, ebox.cursor.y - ebox.scroll)
    termbox.Flush()
}

//----------------------------------------------------------------------------
// main() and support
//----------------------------------------------------------------------------

func mainLoop(ebox *Editbox) {
    ed := ebox.editor
    eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	for {
        select {
        case ev := <-eventQueue:
            switch ev.Type {
            case termbox.EventKey:
                switch ev.Key {
                case termbox.KeyEsc:
                    return
                case termbox.KeyArrowLeft:
                     ed.moveCursorLeft()
                case termbox.KeyArrowRight:
                     ed.moveCursorRight()
                case termbox.KeyArrowUp:
                     ed.moveCursorVert(-1)
                case termbox.KeyArrowDown:
                     ed.moveCursorVert(+1)
                case termbox.KeyHome:
                     ed.moveCursorToLineStart()
                case termbox.KeyEnd:
                     ed.moveCursorToLineEnd()
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
            ebox.Draw()
        }
	}
}


func main() {
    err := termbox.Init()
    check(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
    termbox.SetOutputMode(termbox.Output256)
    ebox := NewEditbox(20, 10, Options{
        wrap: true,
        autoexpand: false,
        fg: 12,
        bg: 63})
    ebox.Draw()
    mainLoop(ebox)
}
