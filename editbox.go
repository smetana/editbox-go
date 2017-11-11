package main

import (
	"github.com/nsf/termbox-go"
//    "fmt"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Cursor struct {
    x int
    y int
}

type Line struct {
    text []rune
    nl bool
}

func (l *Line) insertRune(pos int, r rune) {
    // TODO Raise error on invalid position
    // Append
    if pos == len(l.text) - 1 {
        l.text = append(l.text, r)
    // Insert
    } else {
        l.text = append(l.text, ' ')
        copy(l.text[pos+1:], l.text[pos:])
        l.text[pos] = r
    }
}

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
    ed.lines[0].text = make([]rune, 0)
    ed.lines[0].nl = false
    ed.cursor.x = 0
    ed.cursor.y = 0
    return &ed
}

func (ed *Editor) insertRune(r rune) {
    cursor := &ed.cursor
    line := ed.lines[cursor.y]
    switch {
    // Cursor is at the end of the box
    // and last symbol already exists
    case cursor.x == ed.width - 1 &&
            cursor.y == ed.height - 1 &&
            len(line.text) == ed.width:
        return
    // TODO Move last line character to
    /// next line
    case len(line.text) + 1 > ed.width:
        return
    default:
        line.insertRune(cursor.x, r)
    }
    ed.lines[cursor.y] = line
    cursor.x += 1
    if cursor.x == ed.width {
        if len(ed.lines) < ed.height {
            ed.insertLine(false)
        } else {
            // TODO Better solution
            cursor.x -= 1
        }
    }
    ed.lastx = cursor.x
}

func (ed *Editor) insertLine(nl bool) {
    cursor := &ed.cursor
    if len(ed.lines) == ed.height {
        // TODO Handle this
        return
    }
    line := new(Line)
    ed.lines = append(ed.lines, *line)
    if cursor.y < len(ed.lines) - 1 {
        copy(ed.lines[cursor.y+2:], ed.lines[cursor.y+1:])
        ed.lines[cursor.y+1] = *line
    }
    currentLine := &ed.lines[cursor.y]
    if cursor.x < len(currentLine.text) {
        left, right := new(Line), new(Line)
        left.text = make([]rune, cursor.x)
        copy(left.text, currentLine.text[:cursor.x])
        left.nl = nl
        right.text = make([]rune, len(currentLine.text) - cursor.x)
        copy(right.text, currentLine.text[cursor.x:])
        ed.lines[cursor.y] = *left
        ed.lines[cursor.y+1] = *right
    } else {
        currentLine.nl = nl
    }
    cursor.y += 1
    cursor.x = 0
    ed.lastx = cursor.x
}

func (ed *Editor) moveCursorRight() {
    if len(ed.lines) == 0 {
        return
    }
    cursor := &ed.cursor
    line := ed.lines[cursor.y]
    cursor.x += 1
    if cursor.x > len(line.text) {
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
            line := ed.lines[cursor.y]
            cursor.x = len(line.text)
        } else {
            cursor.x = 0
        }
    }
    ed.lastx = cursor.x
}

func (ed *Editor) moveCursorUp() {
    cursor := &ed.cursor
    if cursor.y == 0 {
        return
    }
    cursor.y -= 1
    line := ed.lines[cursor.y]
    if ed.lastx > len(line.text) {
        cursor.x = len(line.text)
    } else {
        cursor.x = ed.lastx
    }
}

// TODO Code duplucation
func (ed *Editor) moveCursorDown() {
    cursor := &ed.cursor
    if cursor.y == len(ed.lines) - 1 {
        return
    }
    cursor.y += 1
    line := ed.lines[cursor.y]
    if ed.lastx > len(line.text) {
        cursor.x = len(line.text)
    } else {
        cursor.x = ed.lastx
    }
}

func (ed *Editor) Draw() {
    coldef := termbox.ColorDefault
    termbox.Clear(coldef, coldef);
    cursor := ed.cursor

    for y, line := range ed.lines {
        for x, r := range line.text {
            termbox.SetCell(x, y, r, coldef, coldef)
        }
        if line.nl {
            termbox.SetCell(ed.width+2, y, '$', coldef, coldef)
        }
    }
    termbox.SetCursor(cursor.x, cursor.y)
    termbox.Flush()
}

/*
func formatEditor(ed *Editor) {
    for i:=0; i<=len(ed.lines); i++ {
        fmt.Println("")
    }
    for i:=0; i<=10; i++ {
        fmt.Printf("%50s\n", " ")
    }
    fmt.Printf("%v %40s\n", ed.cursor, " ")
    for i, line := range ed.lines {
        fmt.Printf("%v: %v\n", i, line)
    }
}
*/

func main() {
    err := termbox.Init()
    check(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
    ed := NewEditor(5, 5)
    //formatEditor(ed)
    ed.Draw()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			//case termbox.KeyBackspace, termbox.KeyBackspace2:
            //     ed.DeleteRuneBeforeCursor()
			case termbox.KeyArrowLeft:
                 ed.moveCursorLeft()
			case termbox.KeyArrowRight:
                 ed.moveCursorRight()
			case termbox.KeyArrowUp:
                 ed.moveCursorUp()
			case termbox.KeyArrowDown:
                 ed.moveCursorDown()
			case termbox.KeyEnter:
                 ed.insertLine(true)
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
            // TODO Eats CPU. Use time.Sleep ?
        }
        //formatEditor(ed)
        ed.Draw()
    }
}
