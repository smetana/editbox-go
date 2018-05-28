package main

import (
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil { panic(err) }
	editbox.Label(1, 1, 2, termbox.ColorWhite, termbox.ColorBlue, "foobar")
	editbox.Label(1, 2, 0, termbox.ColorWhite, termbox.ColorBlue, "foobar")
	editbox.Label(1, 3, 20, termbox.ColorWhite, termbox.ColorBlue, "foobar")
	editbox.Text(1, 5, 2, 2, termbox.ColorWhite, termbox.ColorRed, `foobar
foobar
foobar`)
	editbox.Text(5, 5, 0, 0, termbox.ColorWhite, termbox.ColorRed, `foo
foobar
foobar foobar`)
	editbox.Text(20, 5, 20, 5, termbox.ColorWhite, termbox.ColorRed, `foo
foobar
foobar foobar`)
	termbox.Flush()
	termbox.PollEvent()
	termbox.Close()
}
