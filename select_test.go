package editbox

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------------------------------------------------
// Support
// ----------------------------------------------------------------------------

func (sbox *SelectBox) toString() string {
	buf := bytes.NewBufferString("")
	var cursor string

	fmt.Fprintf(buf, "\n")

	var index int
	for i := 0; i < sbox.height; i++ {
		index = i + sbox.scroll
		if index == sbox.SelectedIndex() {
			cursor = ">"
		} else {
			cursor = " "
		}
		fmt.Fprintf(buf, "%s %2d %s\n", cursor, index, sbox.items[index])
	}

	return buf.String()
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


func TestCursorWithNotSelectable(t *testing.T) {
	s := Select(
		0, 0, 3, 4,
		0, 0, 0, 0,
		[]string{
			"foo",
			"",
			"bar",
			"baz",
			"qux",
			"",
			"",
			"xyz",
			"qwe",
			"asd",
			"",
			"",
			"",
			"zxc"},
	)
	assert.Equal(t, s.cursor, 0)
	assert.Equal(t, s.SelectedIndex(), 0)
	assert.Equal(t, s.toString(), `
>  0 foo
   1 
   2 bar
   3 baz
`)
	s.cursorDown()
	s.cursorDown()
	s.cursorDown()
	assert.Equal(t, s.cursor, 3)
	assert.Equal(t, s.SelectedIndex(), 4)
	assert.Equal(t, s.Text(), "qux")
	assert.Equal(t, s.toString(), `
   1 
   2 bar
   3 baz
>  4 qux
`)

	s.cursorUp()
	s.cursorUp()
	s.cursorUp()
	s.cursorUp()
	s.cursorUp()
	s.cursorDown()
	assert.Equal(t, s.cursor, 1)
	assert.Equal(t, s.SelectedIndex(), 2)
	assert.Equal(t, s.Text(), "bar")
	assert.Equal(t, s.toString(), `
   0 foo
   1 
>  2 bar
   3 baz
`)

	s.pageDown()
	assert.Equal(t, s.cursor, 4)
	assert.Equal(t, s.SelectedIndex(), 7)
	assert.Equal(t, s.Text(), "xyz")
	assert.Equal(t, s.toString(), `
   4 qux
   5 
   6 
>  7 xyz
`)

	s.pageDown()
	assert.Equal(t, s.cursor, 7)
	assert.Equal(t, s.SelectedIndex(), 13)
	assert.Equal(t, s.Text(), "zxc")
	assert.Equal(t, s.toString(), `
  10 
  11 
  12 
> 13 zxc
`)

	s.pageDown()
	s.pageDown()
	s.pageDown()
	assert.Equal(t, s.cursor, 7)
	assert.Equal(t, s.SelectedIndex(), 13)
	assert.Equal(t, s.Text(), "zxc")
	assert.Equal(t, s.toString(), `
  10 
  11 
  12 
> 13 zxc
`)

	s.pageUp()
	assert.Equal(t, s.cursor, 4)
	assert.Equal(t, s.SelectedIndex(), 7)
	assert.Equal(t, s.Text(), "xyz")
	assert.Equal(t, s.toString(), `
>  7 xyz
   8 qwe
   9 asd
  10 
`)
	s.pageUp()
	assert.Equal(t, s.cursor, 1)
	assert.Equal(t, s.SelectedIndex(), 2)
	assert.Equal(t, s.Text(), "bar")
	assert.Equal(t, s.toString(), `
>  2 bar
   3 baz
   4 qux
   5 
`)
}
