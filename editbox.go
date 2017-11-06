package main

import (
	"github.com/nsf/termbox-go"
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
    text [][]rune
    cursor Cursor
}

func NewEditor() *Editor {
    var ed Editor
    ed.text = make([][]rune, 0)
    ed.cursor.x = 0
    ed.cursor.y = 0
    return &ed
}

func (ed *Editor) addRune(r rune) {
    cursor := &ed.cursor
    if cursor.y > len(ed.text) - 1 {
        ed.text = append(ed.text, make([]rune, 0))
    }
    line := ed.text[cursor.y]
    if cursor.x > len(line) - 1 {
        line = append(line, r)
    } else {
        line[cursor.x] = r
    }
    ed.text[cursor.y] = line
    cursor.x += 1
}

func (ed *Editor) addLine() {
    cursor := &ed.cursor
    cursor.y += 1
    cursor.x = 0
    if cursor.y > len(ed.text) - 1 {
        line := make([]rune, 0)
        ed.text = append(ed.text, line)
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
			//case termbox.KeyArrowLeft:
            //     ed.MoveCursor(-1)
			//case termbox.KeyArrowRight:
            //    ed.MoveCursor(1)
			case termbox.KeyEnter:
                 ed.addLine()
			case termbox.KeySpace:
                 ed.addRune(' ')
			default:
				if ev.Ch != 0 {
                    ed.addRune(ev.Ch)
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
