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

	input.WaitExit()

	termbox.Close()
	fmt.Printf("Text entered: %s\n", input.Text())
}
