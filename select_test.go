package editbox

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// ----------------------------------------------------------------------------
// Support
// ----------------------------------------------------------------------------

func (sbox *SelectBox) toString() string {
	var s strings.Builder
	var cursor string

	fmt.Fprintf(&s, "\n")

	var index int
	for i := 0; i < sbox.height; i++ {
		index = i + sbox.scroll
		if index == sbox.cursor {
			cursor = ">"
		} else {
			cursor = " "
		}
		fmt.Fprintf(&s, "%s %2d %s\n", cursor, index, sbox.Items[index])
	}

	return s.String()
}

// ----------------------------------------------------------------------------

func TestSelectNew(t *testing.T) {
	s := Select(
		0, 0, 3, 3,
		0, 0, 0, 0,
		[]string{"foo", "bar", "baz"},
	)
	assert.Equal(t, s.cursor, 0)
	assert.Equal(t, s.SelectedIndex(), 0)
	assert.Equal(t, s.Text(), "foo")
	assert.Equal(t, s.toString(), `
>  0 foo
   1 bar
   2 baz
`)
}

func TestCursorAndScroll(t *testing.T) {
	s := Select(
		0, 0, 3, 3,
		0, 0, 0, 0,
		[]string{"foo", "bar", "baz", "qux", "xyz", "qwe", "asd", "zxc"},
	)
	assert.Equal(t, s.cursor, 0)
	assert.Equal(t, s.SelectedIndex(), 0)
	assert.Equal(t, s.Text(), "foo")
	assert.Equal(t, s.toString(), `
>  0 foo
   1 bar
   2 baz
`)
	s.cursorDown()
	s.cursorDown()
	s.cursorDown()
	assert.Equal(t, s.cursor, 3)
	assert.Equal(t, s.SelectedIndex(), 3)
	assert.Equal(t, s.Text(), "qux")
	assert.Equal(t, s.toString(), `
   1 bar
   2 baz
>  3 qux
`)

	s.cursorUp()
	s.cursorUp()
	s.cursorUp()
	s.cursorUp()
	s.cursorUp()
	s.cursorDown()
	assert.Equal(t, s.cursor, 1)
	assert.Equal(t, s.SelectedIndex(), 1)
	assert.Equal(t, s.Text(), "bar")
	assert.Equal(t, s.toString(), `
   0 foo
>  1 bar
   2 baz
`)

	s.pageDown()
	assert.Equal(t, s.cursor, 3)
	assert.Equal(t, s.SelectedIndex(), 3)
	assert.Equal(t, s.Text(), "qux")
	assert.Equal(t, s.toString(), `
   1 bar
   2 baz
>  3 qux
`)

	s.pageDown()
	assert.Equal(t, s.cursor, 5)
	assert.Equal(t, s.SelectedIndex(), 5)
	assert.Equal(t, s.Text(), "qwe")
	assert.Equal(t, s.toString(), `
   3 qux
   4 xyz
>  5 qwe
`)

	s.pageDown()
	s.pageDown()
	s.pageDown()
	assert.Equal(t, s.cursor, 7)
	assert.Equal(t, s.SelectedIndex(), 7)
	assert.Equal(t, s.Text(), "zxc")
	assert.Equal(t, s.toString(), `
   5 qwe
   6 asd
>  7 zxc
`)

	s.pageUp()
	assert.Equal(t, s.cursor, 5)
	assert.Equal(t, s.SelectedIndex(), 5)
	assert.Equal(t, s.Text(), "qwe")
	assert.Equal(t, s.toString(), `
>  5 qwe
   6 asd
   7 zxc
`)
	s.pageUp()
	assert.Equal(t, s.cursor, 3)
	assert.Equal(t, s.SelectedIndex(), 3)
	assert.Equal(t, s.Text(), "qux")
	assert.Equal(t, s.toString(), `
>  3 qux
   4 xyz
   5 qwe
`)
}
