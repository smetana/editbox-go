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

	fmt.Println("\n  Type here:")
	input := editbox.NewInputbox(
		13, 1, 25, termbox.ColorWhite, termbox.ColorBlue)

	ev := input.WaitExit()

	termbox.Close()
	fmt.Printf("Exit on: %t\n", ev)
	fmt.Printf("Input was: %s\n", input.Text())
}
