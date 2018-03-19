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
	fmt.Println("TAB - focus, Esc - exit\n\n" +
		"Textarea (wrap == false):" +
		"\n\n\n\n\n\n\n\n\n" +
		"Textarea (wrap == true):")
	inputs := [2]*editbox.Editbox{
		editbox.Textarea(
			0, 3, 25, 7, termbox.ColorWhite, termbox.ColorBlue, false),
		editbox.Textarea(
			0, 12, 25, 7, termbox.ColorWhite, termbox.ColorRed, true),
	}
	termbox.SetCursor(-1, -1)
	termbox.Flush()

	currentInput := inputs[1]
	ev := termbox.PollEvent()
	for {
		switch ev.Key {
		case termbox.KeyEsc:
			termbox.Close()
			fmt.Println("Textarea 1: " + inputs[0].Text())
			fmt.Println("Textarea 2: " + inputs[1].Text())
			return
		case termbox.KeyTab:
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
