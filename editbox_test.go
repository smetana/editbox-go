package editbox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------------------------------------------------
// Support
// ----------------------------------------------------------------------------

func (ed *editor) toLines() []string {
	lines := make([]string, len(ed.lines))
	for i, line := range ed.lines {
		lines[i] = string(line.text)
	}
	return lines
}

// ----------------------------------------------------------------------------

func TestEditorToBox(t *testing.T) {
	eb := newEditbox(0, 0, 3, 3, options{wrap: true})
	eb.editor.setText("1234567\n12\n1234\n1")
	eb.updateLineOffsets()
	assert.Equal(t, eb.lineBoxY, []int{0, 3, 4, 6})
	x, y := eb.editorToBox(0, 0)
	assert.Equal(t, x, 0)
	assert.Equal(t, y, 0)
	x, y = eb.editorToBox(4, 0)
	assert.Equal(t, x, 1)
	assert.Equal(t, y, 1)
	x, y = eb.editorToBox(6, 0)
	assert.Equal(t, x, 0)
	assert.Equal(t, y, 2)
	x, y = eb.editorToBox(7, 0)
	assert.Equal(t, x, 1)
	assert.Equal(t, y, 2)

	// TODO Wrong. There is no text there
	x, y = eb.editorToBox(8, 0)
	assert.Equal(t, x, 2)
	assert.Equal(t, y, 2)

	x, y = eb.editorToBox(1, 1)
	assert.Equal(t, x, 1)
	assert.Equal(t, y, 3)
	x, y = eb.editorToBox(2, 1)
	assert.Equal(t, x, 2)
	assert.Equal(t, y, 3)
	x, y = eb.editorToBox(1, 2)
	assert.Equal(t, x, 1)
	assert.Equal(t, y, 4)
	x, y = eb.editorToBox(3, 2)
	assert.Equal(t, x, 0)
	assert.Equal(t, y, 5)
	x, y = eb.editorToBox(4, 2)
	assert.Equal(t, x, 1)
	assert.Equal(t, y, 5)

	// TODO index out of range
	// x,y = eb.editorToBox(5, 5)
}

/*
Text:

111222333
444555
6667778

0
11122233
44

In editor with width = 3 should be written as

111
222
33␤
444
5␤
666
777
8␤
␤
0␤
111
222
33␤
44

*/

func TestMoveDown(t *testing.T) {
	eb := newEditbox(0, 0, 3, 3, options{wrap: true})
	eb.editor.setText(`11122233
4445
6667778

0
11122233
44`)
	eb.editor.cursor.x = 0
	eb.editor.cursor.y = 0
	eb.editor.moveCursorRight()
	eb.editor.moveCursorRight()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 0)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 1)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 2)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 3)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 1)
	assert.Equal(t, eb.cursor.y, 4)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 5)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 6)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 1)
	assert.Equal(t, eb.cursor.y, 7)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 0)
	assert.Equal(t, eb.cursor.y, 8)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 1)
	assert.Equal(t, eb.cursor.y, 9)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 10)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 11)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 12)
}

func TestMoveDownOneLine(t *testing.T) {
	eb := newEditbox(0, 0, 3, 3, options{wrap: true})
	eb.editor.setText(`11122233`)
	eb.editor.cursor.x = 0
	eb.editor.cursor.y = 0
	eb.editor.moveCursorRight()
	eb.editor.moveCursorRight()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 0)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 1)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 2)
}

func TestMoveUp(t *testing.T) {
	eb := newEditbox(0, 0, 3, 3, options{wrap: true})
	eb.editor.setText(`11122233
4445
6667778

0
11122233
44`)
	eb.editor.cursor.x = 0
	eb.editor.cursor.y = 6
	eb.editor.moveCursorRight()
	eb.editor.moveCursorRight()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 13)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 12)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 11)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 10)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 1)
	assert.Equal(t, eb.cursor.y, 9)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 0)
	assert.Equal(t, eb.cursor.y, 8)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 1)
	assert.Equal(t, eb.cursor.y, 7)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 6)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 5)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 1)
	assert.Equal(t, eb.cursor.y, 4)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 3)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 2)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 1)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 0)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assert.Equal(t, eb.cursor.x, 2)
	assert.Equal(t, eb.cursor.y, 0)
}
