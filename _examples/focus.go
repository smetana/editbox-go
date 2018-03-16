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
	fmt.Println("Press TAB to focus input box\n\n" +
		"Input 1:\n\n" +
		"Input 2:")
	inputs := [2]*editbox.Editbox{
		editbox.NewInputbox(10, 2, 25, termbox.ColorWhite, termbox.ColorBlue),
		editbox.NewInputbox(10, 4, 25, termbox.ColorWhite, termbox.ColorRed),
	}
	termbox.SetCursor(-1, -1)
	termbox.Flush()

	currentInput := inputs[1]
	ev := termbox.PollEvent()
	for {
		switch ev.Key {
		case termbox.KeyEsc:
			termbox.Close()
			fmt.Println("Input 1: " + inputs[0].Text())
			fmt.Println("Input 2: " + inputs[1].Text())
			return
		case termbox.KeyTab, termbox.KeyEnter:
			if currentInput == inputs[1] {
				currentInput = inputs[0]
			} else {
				currentInput = inputs[1]
			}
			ev = currentInput.WaitExit()
		default:
			ev = termbox.PollEvent()
		}
	}
}
