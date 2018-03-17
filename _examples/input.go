package main

import (
	".."
	"fmt"
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("Press Esc, Enter, or Tab to Exit")
	fmt.Println("Type here:")
	input := editbox.Input(11, 1, 25, termbox.ColorWhite, termbox.ColorBlue)
	input.WaitExit()
	termbox.Close()
	fmt.Printf("Input was: %s\n", input.Text())
}
