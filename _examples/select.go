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
	editbox.Text(0, 0, 0, 0, 0, "Press Esc, Enter, or Tab to Exit")
	editbox.Text(0, 2, 0, 0, 0, "Select:")
	input := editbox.Select(8, 2, 10, 4,
		termbox.ColorWhite,
		termbox.ColorBlue,
		termbox.ColorWhite | termbox.AttrReverse,
		termbox.ColorBlue | termbox.AttrReverse,
		[]string{
			"foo",
			"bar",
			"baz",
			"qux",
			"quux",
			"corge",
			"grault",
			"garply",
			"waldo",
			"fred",
			"plugh",
			"xyzzy",
		})
	input.WaitExit()
	termbox.Close()
	fmt.Printf("Input was: %s\n", input.Text())
}
