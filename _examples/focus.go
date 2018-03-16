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

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	currentInput := -1
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					termbox.Close()
					fmt.Println("Input 1: " + inputs[0].Text())
					fmt.Println("Input 2: " + inputs[1].Text())
					return
				case termbox.KeyTab:
					currentInput++
					if currentInput > 1 {
						currentInput = 0
					}
					ev = inputs[currentInput].WaitExit()
					go func() { eventQueue <- ev }()
				}
			}
		}
	}
}
