package main

import (
	"github.com/nsf/termbox-go"
)

const maxLineLen = 5
const maxLines = 5

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
    text [][]rune
    cursor Cursor
}

func NewEditor() *Editor {
    var ed Editor
    ed.text = make([][]rune, 1)
    ed.text[0] = make([]rune, 0)
    ed.cursor.x = 0
    ed.cursor.y = 0
    return &ed
}

func (ed *Editor) insertRune(r rune) {
    cursor := &ed.cursor
    line := ed.text[cursor.y]
    // TODO Better solution
    if cursor.x == maxLineLen - 1 &&
            cursor.y== maxLines - 1 {
        return
    }
    if cursor.x == len(line) {
        line = append(line, r)
    } else {
        line = append(line, ' ')
        copy(line[cursor.x+1:], line[cursor.x:])
        line[cursor.x] = r
    }
    ed.text[cursor.y] = line
    cursor.x += 1
    if cursor.x == maxLineLen {
        if len(ed.text) < maxLines {
            ed.insertLine()
        } else {
            // TODO Better solution
            cursor.x -= 1
        }
    }
}

func (ed *Editor) insertLine() {
    cursor := &ed.cursor
    if len(ed.text) == maxLines {
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
    if cursor.x <= len(currentLine) {
        left := make([]rune, cursor.x + 1)
        copy(left, currentLine[:cursor.x])
        right := make([]rune, len(currentLine) - cursor.x)
        copy(right, currentLine[cursor.x:])
        ed.text[cursor.y] = left
        ed.text[cursor.y+1] = right
    }
    cursor.y += 1
    cursor.x = 0
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

func main() {
    err := termbox.Init()
    check(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
    ed := NewEditor()
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
        ed.Draw()
    }
}
