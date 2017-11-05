package main

import (
	"github.com/nsf/termbox-go"
)

const EditorPageSize = 1024

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Editor struct {
    text []rune
    cursor int
}

func NewEditor() *Editor {
    var ed Editor
    ed.text = make([]rune, 0, EditorPageSize)
    ed.cursor = 0
    return &ed
}

// Increase editor's memory footprint if editor's text
// does not fit into undelying array
func (ed *Editor) addPage() {
    newSlice := make([]rune, len(ed.text), cap(ed.text) + EditorPageSize)
    copy(newSlice, ed.text)
    ed.text = newSlice
}

func (ed *Editor) InsertRune(r rune) {
    if len(ed.text) == cap(ed.text) {
        ed.addPage()
    }
    ed.text = ed.text[:len(ed.text)+1]
    ed.text[ed.cursor] = r
    ed.cursor += 1
}

func (ed *Editor) DeleteRuneBeforeCursor() {
    if ed.cursor > 0 {
        ed.text = ed.text[:len(ed.text)-1]
        ed.cursor -= 1
    }
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
    ed.Draw()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
                 ed.DeleteRuneBeforeCursor()
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
            // TODO Eats CPU. Use time.Sleep ?
        }
        ed.Draw()
    }
}
