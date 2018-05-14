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
	editbox.Label(0, 0, 0, 0, 0, "Press TAB to focus input box")
	editbox.Label(0, 2, 0, 0, 0, "Input 1:")
	editbox.Label(0, 4, 0, 0, 0, "Input 2:")
	inputs := [2]*editbox.Editbox{
		editbox.Input(10, 2, 25, termbox.ColorWhite, termbox.ColorBlue),
		editbox.Input(10, 4, 25, termbox.ColorWhite, termbox.ColorRed),
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
