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

type Editor struct {
    width, height int
    text [][]rune
    cursor Cursor
    lastx int
}

func NewEditor(width, height int) *Editor {
    var ed Editor
    ed.width = width
    ed.height = height
    ed.text = make([][]rune, 1)
    ed.text[0] = make([]rune, 0)
    ed.cursor.x = 0
    ed.cursor.y = 0
    return &ed
}

func (ed *Editor) insertRune(r rune) {
    cursor := &ed.cursor
    line := ed.text[cursor.y]
    switch {
    // Cursor is at the end of the box
    // and last symbol already exists
    case cursor.x == ed.width - 1 &&
            cursor.y == ed.height - 1 &&
            len(line) == ed.width:
        line[cursor.x] = r
    case cursor.x == len(line) - 1:
        line = append(line, r)
    default:
        line = append(line, ' ')
        copy(line[cursor.x+1:], line[cursor.x:])
        line[cursor.x] = r
    }
    ed.text[cursor.y] = line
    cursor.x += 1
    if cursor.x == ed.width {
        if len(ed.text) < ed.height {
            ed.insertLine()
        } else {
            // TODO Better solution
            cursor.x -= 1
        }
    }
    ed.lastx = cursor.x
}

func (ed *Editor) insertLine() {
    cursor := &ed.cursor
    if len(ed.text) == ed.height {
        // TODO Handle this
        return
    }
    if cursor.y == len(ed.text) - 1 {
        line := make([]rune, 0)
        ed.text = append(ed.text, line)
    } else {
        newLine := make([]rune, 0)
        ed.text = append(ed.text, newLine)
        copy(ed.text[cursor.y+2:], ed.text[cursor.y+1:])
        ed.text[cursor.y+1] = newLine
    }
    currentLine := ed.text[cursor.y]
    if cursor.x < len(currentLine) {
        left := make([]rune, cursor.x)
        copy(left, currentLine[:cursor.x])
        right := make([]rune, len(currentLine) - cursor.x)
        copy(right, currentLine[cursor.x:])
        ed.text[cursor.y] = left
        ed.text[cursor.y+1] = right
    }
    cursor.y += 1
    cursor.x = 0
    ed.lastx = cursor.x
}

func (ed *Editor) moveCursorRight() {
    if len(ed.text) == 0 {
        return
    }
    cursor := &ed.cursor
    line := ed.text[cursor.y]
    cursor.x += 1
    if cursor.x >= len(line) {
        if cursor.y < len(ed.text) - 1 {
            cursor.y += 1
            cursor.x = 0
        } else {
            cursor.x = len(line)
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
            line := ed.text[cursor.y]
            cursor.x = len(line) - 1
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
    line := ed.text[cursor.y]
    if ed.lastx > len(line) {
        cursor.x = len(line)
    } else {
        cursor.x = ed.lastx
    }
}

// TODO Code duplucation
func (ed *Editor) moveCursorDown() {
    cursor := &ed.cursor
    if cursor.y == len(ed.text) - 1 {
        return
    }
    cursor.y += 1
    line := ed.text[cursor.y]
    if ed.lastx > len(line) {
        cursor.x = len(line)
    } else {
        cursor.x = ed.lastx
    }
}

func (ed *Editor) Draw() {
    coldef := termbox.ColorDefault
    termbox.Clear(coldef, coldef);
    cursor := ed.cursor

    for y, line := range ed.text {
        for x, r := range line {
            termbox.SetCell(x, y, r, coldef, coldef)
        }
    }
    termbox.SetCursor(cursor.x, cursor.y)
    termbox.Flush()
}

/*
func formatEditor(ed *Editor) {
    for i:=0; i<=len(ed.text); i++ {
        fmt.Println("")
    }
    for i:=0; i<=10; i++ {
        fmt.Printf("%50s\n", " ")
    }
    fmt.Printf("%v %40s\n", ed.cursor, " ")
    for i, line := range ed.text {
        fmt.Printf("%v: %v\n", i, line)
    }
}
*/

func main() {
    err := termbox.Init()
    check(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
    ed := NewEditor(40, 20)
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
                 ed.insertLine()
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
