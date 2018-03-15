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

	// ebox.Render() only puts its editor's content into
	// termbox CellBuffer but does not flush it
	termbox.Flush()

	ok := true
	for ok {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyEnter:
				ok = false
			default:
				ok = input.HandleEvent(&ev)
			}
		}
		termbox.Flush()
	}

	termbox.Close()
	fmt.Printf("Text entered: %s\n", input.Text())
}
