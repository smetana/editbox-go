package main

import (
	"github.com/nsf/termbox-go"
)

type Editor struct {
    text []rune
    cursor int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewEditor() *Editor {
    var ed Editor
    ed.text = make([]rune, 0, 1024)
    ed.cursor = 0
    return &ed
}

func (ed *Editor) InsertRune(r rune) {
    ed.text = ed.text[:len(ed.text)+1]
    ed.text[ed.cursor] = r
    ed.cursor += 1
}

func (ed *Editor) DeleteRuneBeforeCursor() {
    ed.text = ed.text[:len(ed.text)-1]
    ed.cursor -= 1
}

func (ed *Editor) Draw() {
    coldef := termbox.ColorDefault
    termbox.Clear(coldef, coldef);
    x,y := 0,0
    for _, r := range ed.text {
        if r == '\n' {
            x = 0
            y += 1
        }
        termbox.SetCell(x, y, r, coldef, coldef)
        if r != '\n' {
            x += 1
        }
    }
    termbox.SetCursor(x, y)
    termbox.Flush()
}

func main() {
	err := termbox.Init()
    check(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
    ed := NewEditor()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
                 ed.DeleteRuneBeforeCursor()
                 ed.Draw()
			case termbox.KeyEnter:
                 ed.InsertRune('\n')
                 ed.Draw()
			case termbox.KeySpace:
                 ed.InsertRune(' ')
                 ed.Draw()
			default:
				if ev.Ch != 0 {
                    ed.InsertRune(ev.Ch)
                    ed.Draw()
                }
			}
		case termbox.EventError:
			panic(ev.Err)
        default:
            // TODO Eats CPU. Use time.Sleep ?
        }
    }
}
