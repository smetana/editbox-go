package main

import (
	".."
	"fmt"
	"github.com/nsf/termbox-go"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func mainLoop(ebox *editbox.Editbox) {
	eventQueue := make(chan termbox.Event, 256)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	for {
		select {
		case ev := <-eventQueue:
			ok := ebox.HandleEvent(&ev)
			if !ok {
				return
			}
			if len(eventQueue) == 0 {
				ebox.Draw()
				termbox.Flush()
			}
		}
	}
}

func main() {
	err := termbox.Init()
	check(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)
	fmt.Println("Editor Test\n\n\n\nEnter text here:")
	ebox := editbox.NewEditbox(17, 4, 20, 3, editbox.Options{
		Wrap:       true,
		Autoexpand: true,
		MaxHeight:  6,
		Fg:         12,
		Bg:         63})
	ebox.Draw()
	termbox.Flush()
	mainLoop(ebox)
}
