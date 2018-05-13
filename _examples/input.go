package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/smetana/editbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	editbox.Text(0, 0, 0, 0, 0, "Press Esc, Enter, or Tab to Exit")
	editbox.Text(0, 1, 0, 0, 0, "Type here:")
	input := editbox.Input(11, 1, 25, termbox.ColorWhite, termbox.ColorBlue)
	input.WaitExit()
	termbox.Close()
	fmt.Printf("Input was: %s\n", input.Text())
}
