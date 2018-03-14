package main

import (
	".."
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
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				ebox.Draw()
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
	ebox := editbox.NewEditbox(12, 10, 20, 3, editbox.Options{
		Wrap:       true,
		Autoexpand: true,
		MaxHeight:  6,
		Fg:         12,
		Bg:         63})
	ebox.Draw()
	mainLoop(ebox)
}
