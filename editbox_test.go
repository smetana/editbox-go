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
// line Tests
// ----------------------------------------------------------------------------

func TestLineSimpleInsertRune(t *testing.T) {
	l := new(line)
	l.insertRune(0, 'H')
	l.insertRune(1, 'e')
	l.insertRune(2, 'l')
	l.insertRune(3, 'l')
	l.insertRune(4, 'o')
	res := string(l.text)
	assert.Equal(t, res, "Hello")
}

func TestLineInsertRune(t *testing.T) {
	l := new(line)
	l.text = []rune("Sick")
	l.insertRune(1, 'l')
	assert.Equal(t, string(l.text), "Slick")
}

func TestLineInsertPostion(t *testing.T) {
	l := new(line)
	l.text = []rune("1")
	l.insertRune(0, '2')
	assert.Equal(t, string(l.text), "21")
}

func TestLineInsertCornerCase1(t *testing.T) {
	l := new(line)
	l.text = []rune("1")
	l.insertRune(1, '2')
	assert.Equal(t, string(l.text), "12")
}

func TestLineInsertOnWrongPosition(t *testing.T) {
	defer func() {
		if r := recover(); r != "x position out of range" {
			t.Errorf("Wrong panic: %+q", r)
		}
	}()
	l := new(line)
	l.text = []rune("1")
	l.insertRune(2, '2')
}

func TestLineInsertNewLine(t *testing.T) {
	l := new(line)
	l.text = []rune("HelloWorld")
	l.insertRune(5, '\n')
	assert.Equal(t, string(l.text), "Hello\nWorld")
}

func TestLineSplit(t *testing.T) {
	l := new(line)
	l.text = []rune("Hello World")
	left, right := l.split(5)
	assert.Equal(t, string(left.text), "Hello")
	assert.Equal(t, string(right.text), " World")
}

func TestLineSplitOnWrongPosition(t *testing.T) {
	defer func() {
		if r := recover(); r != "x position out of range" {
			t.Errorf("Wrong panic: %+q", r)
		}
	}()
	l := new(line)
	l.text = []rune("Sick")
	_, _ = l.split(10)
}

func TestLineDeleteOnWrongPosition(t *testing.T) {
	defer func() {
		if r := recover(); r != "x position out of range" {
			t.Errorf("Wrong panic: %+q", r)
		}
	}()
	l := new(line)
	l.text = []rune("1")
	l.deleteRune(2)
}

func TestLineDelete(t *testing.T) {
	l := new(line)
	l.text = []rune("12")
	l.deleteRune(1)
	assert.Equal(t, string(l.text), "1")
	l.text = []rune("12")
	l.deleteRune(0)
	assert.Equal(t, string(l.text), "2")
	l.text = []rune("")
	l.deleteRune(0)
	assert.Equal(t, string(l.text), "")
}

func TestLineLastRune(t *testing.T) {
	l := new(line)
	l.text = []rune("12")
	assert.Equal(t, l.lastRune(), '2')
	l.text = []rune("12\n")
	assert.Equal(t, l.lastRune(), '\n')
}

func TestLineLastRuneX(t *testing.T) {
	l := new(line)
	l.text = []rune("12")
	assert.Equal(t, l.lastRuneX(), 2)
	l.text = []rune("12\n")
	assert.Equal(t, l.lastRuneX(), 2)
}

// ----------------------------------------------------------------------------
// editor Tests
// ----------------------------------------------------------------------------

func TestEditorInsertRune(t *testing.T) {
	ed := newEditor()
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 0)
	ed.insertRune('H')
	ed.insertRune('e')
	ed.insertRune('l')
	ed.insertRune('l')
	ed.insertRune('o')
	assert.Equal(t, string(ed.text()), "Hello")
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 5)
	ed.insertRune('!')
	assert.Equal(t, string(ed.text()), "Hello!")
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 6)
}

func TestEditorInsertOnCursorPosition(t *testing.T) {
	ed := newEditor()
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 0)
	ed.insertRune('1')
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 1)
	assert.Equal(t, string(ed.text()), "1")
	ed.moveCursorLeft()
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 0)
	ed.insertRune('2')
	assert.Equal(t, string(ed.text()), "21")
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 1)
}

func TestEditorCurrentLine(t *testing.T) {
	ed := newEditor()
	ed.setText("Hello World!\nSecond Line\nThird Line")
	ed.cursor.x = 2
	ed.cursor.y = 1
	assert.Equal(t, string(ed.currentLine().text), "Second Line\n")
}

func TestEditorSplitLine(t *testing.T) {
	ed := newEditor()
	ed.setText("123\n123\n123")
	ed.splitLine(1, 1)
	assert.Equal(t, ed.toLines(), []string{"123\n", "1", "23\n", "123"})

	ed.splitLine(3, 3)
	assert.Equal(t, ed.toLines(), []string{"123\n", "1", "23\n", "123", ""})
}

func TestEditorInsertNewLine(t *testing.T) {
	ed := newEditor()
	ed.setText("12345")
	assert.Equal(t, len(ed.lines), 1)
	ed.insertRune('\n')
	assert.Equal(t, len(ed.lines), 2)
}

func TestMoveCursorLeft(t *testing.T) {
	ed := newEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
	assert.Equal(t, len(ed.lines), 2)
	assert.Equal(t, ed.cursor.y, 1)
	assert.Equal(t, ed.cursor.x, 0)
	ed.moveCursorLeft()
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 2)
}

func TestBackspace(t *testing.T) {
	ed := newEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
	ed.insertRune('1')

	assert.Equal(t, len(ed.lines), 2)
	assert.Equal(t, ed.cursor.y, 1)
	assert.Equal(t, ed.cursor.x, 1)
	ed.deleteRuneBeforeCursor()

	assert.Equal(t, ed.text(), "12\n")
	assert.Equal(t, len(ed.lines), 2)
	assert.Equal(t, ed.cursor.y, 1)
	assert.Equal(t, ed.cursor.x, 0)

	ed.deleteRuneBeforeCursor()
	assert.Equal(t, ed.text(), "12")
	assert.Equal(t, len(ed.lines), 1)
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 2)

	ed.deleteRuneBeforeCursor()
	ed.deleteRuneBeforeCursor()

	assert.Equal(t, ed.text(), "")
	assert.Equal(t, len(ed.lines), 1)
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 0)
}

func TestDeleteAtCursor(t *testing.T) {
	ed := newEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
	ed.insertRune('3')
	ed.insertRune('\n')
	ed.insertRune('4')
	ed.insertRune('5')

	ed.cursor.x = 0
	ed.cursor.y = 1

	assert.Equal(t, ed.text(), "12\n3\n45")
	assert.Equal(t, len(ed.lines), 3)
	assert.Equal(t, ed.cursor.y, 1)
	assert.Equal(t, ed.cursor.x, 0)

	ed.deleteRuneAtCursor()
	assert.Equal(t, len(ed.lines), 3)
	assert.Equal(t, ed.text(), "12\n\n45")

	ed.deleteRuneAtCursor()
	assert.Equal(t, len(ed.lines), 2)
	assert.Equal(t, ed.text(), "12\n45")

	ed.deleteRuneAtCursor()
	ed.deleteRuneAtCursor()
	ed.deleteRuneAtCursor() // No effect

	assert.Equal(t, len(ed.lines), 2)
	assert.Equal(t, ed.text(), "12\n")
}

func TestMoveCursorToLineEnd(t *testing.T) {
	ed := newEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
	ed.insertRune('3')
	ed.insertRune('\n')
	ed.insertRune('4')
	ed.insertRune('5')

	ed.cursor.x = 0
	ed.cursor.y = 0
	ed.moveCursorToLineEnd()
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 2)

	ed.cursor.x = 0
	ed.cursor.y = 2
	ed.moveCursorToLineEnd()
	assert.Equal(t, ed.cursor.y, 2)
	assert.Equal(t, ed.cursor.x, 2)
}

func TestMoveCursorToEmptyLine(t *testing.T) {
	ed := newEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
	assert.Equal(t, ed.cursor.y, 1)
	assert.Equal(t, ed.cursor.x, 0)

	ed.moveCursorVert(-1)
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 0)

	ed.moveCursorVert(+1)
	assert.Equal(t, ed.cursor.y, 1)
	assert.Equal(t, ed.cursor.x, 0)
}

// TODO Add tests for cursor navigation

// ----------------------------------------------------------------------------
// EditBox Tests
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
