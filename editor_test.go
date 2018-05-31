package editbox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

	ed.moveCursorToLineEnd()
	assert.Equal(t, ed.cursor.y, 0)
	assert.Equal(t, ed.cursor.x, 0)

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
