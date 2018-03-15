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
	defer termbox.Close()
	fmt.Println("\n  Type here:")

	ebox := editbox.NewEditbox(13, 1, 25, 3, editbox.Options{
		Wrap:       true,
		Autoexpand: true,
		MaxHeight:  6,
		Fg:         termbox.ColorWhite,
		Bg:         termbox.ColorBlue})

	// ebox.Render() only puts its editor's content into termbox
	// CellBuffer but does not flush it
	termbox.Flush()

	for {
		ev := termbox.PollEvent();
		if ev.Type == termbox.EventKey {
			ok := ebox.HandleEvent(&ev)
			if !ok {
				return // Quit
			}
		}
		ebox.Render()
		termbox.Flush()
	}
}
