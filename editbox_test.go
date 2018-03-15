package editbox

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// ----------------------------------------------------------------------------
// Support
// ----------------------------------------------------------------------------

// CallerInfo is borrowed from https://github.com/stretchr/testify

/* CallerInfo is necessary because the assert functions use the testing object
internally, causing it to print the file:line of the assert method, rather than where
the problem actually occurred in calling code.*/

// CallerInfo returns an array of strings containing the file and line number
// of each stack frame leading from the current test to the assert call that
// failed.
func CallerInfo() string {

	pc := uintptr(0)
	file := ""
	line := 0
	ok := false
	name := ""

	callers := []string{}
	for i := 2; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			// The breaks below failed to terminate the loop, and we ran off the
			// end of the call stack.
			break
		}

		// This is a huge edge case, but it will panic if this is the case, see #180
		if file == "<autogenerated>" {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		name = f.Name()

		// testing.tRunner is the standard library function that calls
		// tests. Subtests are called directly by tRunner, without going through
		// the Test/Benchmark/Example function that contains the t.Run calls, so
		// with subtests we should break when we hit tRunner, without adding it
		// to the list of callers.
		if name == "testing.tRunner" {
			break
		}

		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
		if len(parts) > 1 {
			dir := parts[len(parts)-2]
			if (dir != "assert" && dir != "mock" && dir != "require") || file == "mock_test.go" {
				callers = append(callers, fmt.Sprintf("%s:%d %s", file, line, name))
			}
		}

		// Drop the package
		segments := strings.Split(name, ".")
		name = segments[len(segments)-1]
		if isTest(name, "Test") ||
			isTest(name, "Benchmark") ||
			isTest(name, "Example") {
			break
		}
	}

	return strings.Join(callers, "\n")
}

// Stolen from the `go test` tool.
// isTest tells whether name looks like a test (or benchmark, according to prefix).
// It is a Test (say) if there is a character after Test that is not a lower-case letter.
// We don't want TesticularCancer.
func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		trace := CallerInfo()
		t.Errorf("Expected (%T)%v got (%T)%v\n%s", expected, expected, actual, actual, trace)
	}
}

func (ed *Editor) toLines() []string {
	lines := make([]string, len(ed.Lines))
	for i, line := range ed.Lines {
		lines[i] = string(line.Text)
	}
	return lines
}

// ----------------------------------------------------------------------------
// Line Tests
// ----------------------------------------------------------------------------

func TestLineSimpleInsertRune(t *testing.T) {
	l := new(Line)
	l.InsertRune(0, 'H')
	l.InsertRune(1, 'e')
	l.InsertRune(2, 'l')
	l.InsertRune(3, 'l')
	l.InsertRune(4, 'o')
	res := string(l.Text)
	assertEqual(t, res, "Hello")
}

func TestLineInsertRune(t *testing.T) {
	l := new(Line)
	l.Text = []rune("Sick")
	l.InsertRune(1, 'l')
	assertEqual(t, string(l.Text), "Slick")
}

func TestLineInsertPostion(t *testing.T) {
	l := new(Line)
	l.Text = []rune("1")
	l.InsertRune(0, '2')
	assertEqual(t, string(l.Text), "21")
}

func TestLineInsertCornerCase1(t *testing.T) {
	l := new(Line)
	l.Text = []rune("1")
	l.InsertRune(1, '2')
	assertEqual(t, string(l.Text), "12")
}

func TestLineInsertOnWrongPosition(t *testing.T) {
	defer func() {
		if r := recover(); r != "x position out of range" {
			t.Errorf("Wrong panic: %+q", r)
		}
	}()
	l := new(Line)
	l.Text = []rune("1")
	l.InsertRune(2, '2')
}

func TestLineInsertNewLine(t *testing.T) {
	l := new(Line)
	l.Text = []rune("HelloWorld")
	l.InsertRune(5, '\n')
	assertEqual(t, string(l.Text), "Hello\nWorld")
}

func TestLineSplit(t *testing.T) {
	l := new(Line)
	l.Text = []rune("Hello World")
	left, right := l.Split(5)
	assertEqual(t, string(left.Text), "Hello")
	assertEqual(t, string(right.Text), " World")
}

func TestLineSplitOnWrongPosition(t *testing.T) {
	defer func() {
		if r := recover(); r != "x position out of range" {
			t.Errorf("Wrong panic: %+q", r)
		}
	}()
	l := new(Line)
	l.Text = []rune("Sick")
	_, _ = l.Split(10)
}

func TestLineDeleteOnWrongPosition(t *testing.T) {
	defer func() {
		if r := recover(); r != "x position out of range" {
			t.Errorf("Wrong panic: %+q", r)
		}
	}()
	l := new(Line)
	l.Text = []rune("1")
	l.DeleteRune(2)
}

func TestLineDelete(t *testing.T) {
	l := new(Line)
	l.Text = []rune("12")
	l.DeleteRune(1)
	assertEqual(t, string(l.Text), "1")
	l.Text = []rune("12")
	l.DeleteRune(0)
	assertEqual(t, string(l.Text), "2")
	l.Text = []rune("")
	l.DeleteRune(0)
	assertEqual(t, string(l.Text), "")
}

func TestLineLastRune(t *testing.T) {
	l := new(Line)
	l.Text = []rune("12")
	assertEqual(t, l.lastRune(), '2')
	l.Text = []rune("12\n")
	assertEqual(t, l.lastRune(), '\n')
}

func TestLineLastRuneX(t *testing.T) {
	l := new(Line)
	l.Text = []rune("12")
	assertEqual(t, l.lastRuneX(), 2)
	l.Text = []rune("12\n")
	assertEqual(t, l.lastRuneX(), 2)
}

// ----------------------------------------------------------------------------
// Editor Tests
// ----------------------------------------------------------------------------

func TestEditorInsertRune(t *testing.T) {
	ed := NewEditor()
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 0)
	ed.InsertRune('H')
	ed.InsertRune('e')
	ed.InsertRune('l')
	ed.InsertRune('l')
	ed.InsertRune('o')
	assertEqual(t, string(ed.Text()), "Hello")
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 5)
	ed.InsertRune('!')
	assertEqual(t, string(ed.Text()), "Hello!")
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 6)
}

func TestEditorInsertOnCursorPosition(t *testing.T) {
	ed := NewEditor()
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 0)
	ed.InsertRune('1')
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 1)
	assertEqual(t, string(ed.Text()), "1")
	ed.moveCursorLeft()
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 0)
	ed.InsertRune('2')
	assertEqual(t, string(ed.Text()), "21")
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 1)
}

func TestEditorCurrentLine(t *testing.T) {
	ed := NewEditor()
	ed.setText("Hello World!\nSecond Line\nThird Line")
	ed.Cursor.X = 2
	ed.Cursor.Y = 1
	assertEqual(t, string(ed.CurrentLine().Text), "Second Line\n")
}

func TestEditorSplitLine(t *testing.T) {
	ed := NewEditor()
	ed.setText("123\n123\n123")
	ed.splitLine(1, 1)
	assertEqual(t, ed.toLines(), []string{"123\n", "1", "23\n", "123"})

	ed.splitLine(3, 3)
	assertEqual(t, ed.toLines(), []string{"123\n", "1", "23\n", "123", ""})
}

func TestEditorInsertNewLine(t *testing.T) {
	ed := NewEditor()
	ed.setText("12345")
	assertEqual(t, len(ed.Lines), 1)
	ed.InsertRune('\n')
	assertEqual(t, len(ed.Lines), 2)
}

func TestMoveCursorLeft(t *testing.T) {
	ed := NewEditor()
	ed.InsertRune('1')
	ed.InsertRune('2')
	ed.InsertRune('\n')
	assertEqual(t, len(ed.Lines), 2)
	assertEqual(t, ed.Cursor.Y, 1)
	assertEqual(t, ed.Cursor.X, 0)
	ed.moveCursorLeft()
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 2)
}

func TestBackspace(t *testing.T) {
	ed := NewEditor()
	ed.InsertRune('1')
	ed.InsertRune('2')
	ed.InsertRune('\n')
	ed.InsertRune('1')

	assertEqual(t, len(ed.Lines), 2)
	assertEqual(t, ed.Cursor.Y, 1)
	assertEqual(t, ed.Cursor.X, 1)
	ed.DeleteRuneBeforeCursor()

	assertEqual(t, ed.Text(), "12\n")
	assertEqual(t, len(ed.Lines), 2)
	assertEqual(t, ed.Cursor.Y, 1)
	assertEqual(t, ed.Cursor.X, 0)

	ed.DeleteRuneBeforeCursor()
	assertEqual(t, ed.Text(), "12")
	assertEqual(t, len(ed.Lines), 1)
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 2)

	ed.DeleteRuneBeforeCursor()
	ed.DeleteRuneBeforeCursor()

	assertEqual(t, ed.Text(), "")
	assertEqual(t, len(ed.Lines), 1)
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 0)
}

func TestDeleteAtCursor(t *testing.T) {
	ed := NewEditor()
	ed.InsertRune('1')
	ed.InsertRune('2')
	ed.InsertRune('\n')
	ed.InsertRune('3')
	ed.InsertRune('\n')
	ed.InsertRune('4')
	ed.InsertRune('5')

	ed.Cursor.X = 0
	ed.Cursor.Y = 1

	assertEqual(t, ed.Text(), "12\n3\n45")
	assertEqual(t, len(ed.Lines), 3)
	assertEqual(t, ed.Cursor.Y, 1)
	assertEqual(t, ed.Cursor.X, 0)

	ed.DeleteRuneAtCursor()
	assertEqual(t, len(ed.Lines), 3)
	assertEqual(t, ed.Text(), "12\n\n45")

	ed.DeleteRuneAtCursor()
	assertEqual(t, len(ed.Lines), 2)
	assertEqual(t, ed.Text(), "12\n45")

	ed.DeleteRuneAtCursor()
	ed.DeleteRuneAtCursor()
	ed.DeleteRuneAtCursor() // No effect

	assertEqual(t, len(ed.Lines), 2)
	assertEqual(t, ed.Text(), "12\n")
}

func TestMoveCursorToLineEnd(t *testing.T) {
	ed := NewEditor()
	ed.InsertRune('1')
	ed.InsertRune('2')
	ed.InsertRune('\n')
	ed.InsertRune('3')
	ed.InsertRune('\n')
	ed.InsertRune('4')
	ed.InsertRune('5')

	ed.Cursor.X = 0
	ed.Cursor.Y = 0
	ed.moveCursorToLineEnd()
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 2)

	ed.Cursor.X = 0
	ed.Cursor.Y = 2
	ed.moveCursorToLineEnd()
	assertEqual(t, ed.Cursor.Y, 2)
	assertEqual(t, ed.Cursor.X, 2)
}

func TestMoveCursorToEmptyLine(t *testing.T) {
	ed := NewEditor()
	ed.InsertRune('1')
	ed.InsertRune('2')
	ed.InsertRune('\n')
	assertEqual(t, ed.Cursor.Y, 1)
	assertEqual(t, ed.Cursor.X, 0)

	ed.moveCursorVert(-1)
	assertEqual(t, ed.Cursor.Y, 0)
	assertEqual(t, ed.Cursor.X, 0)

	ed.moveCursorVert(+1)
	assertEqual(t, ed.Cursor.Y, 1)
	assertEqual(t, ed.Cursor.X, 0)
}

// TODO Add tests for cursor navigation

// ----------------------------------------------------------------------------
// EditBox Tests
// ----------------------------------------------------------------------------

func TestEditorToBox(t *testing.T) {
	eb := NewEditbox(0, 0, 3, 3, Options{Wrap: true})
	eb.Editor.setText("1234567\n12\n1234\n1")
	eb.updateLineOffsets()
	assertEqual(t, eb.lineBoxY, []int{0, 3, 4, 6})
	x, y := eb.editorToBox(0, 0)
	assertEqual(t, x, 0)
	assertEqual(t, y, 0)
	x, y = eb.editorToBox(4, 0)
	assertEqual(t, x, 1)
	assertEqual(t, y, 1)
	x, y = eb.editorToBox(6, 0)
	assertEqual(t, x, 0)
	assertEqual(t, y, 2)
	x, y = eb.editorToBox(7, 0)
	assertEqual(t, x, 1)
	assertEqual(t, y, 2)

	// TODO Wrong. There is no text there
	x, y = eb.editorToBox(8, 0)
	assertEqual(t, x, 2)
	assertEqual(t, y, 2)

	x, y = eb.editorToBox(1, 1)
	assertEqual(t, x, 1)
	assertEqual(t, y, 3)
	x, y = eb.editorToBox(2, 1)
	assertEqual(t, x, 2)
	assertEqual(t, y, 3)
	x, y = eb.editorToBox(1, 2)
	assertEqual(t, x, 1)
	assertEqual(t, y, 4)
	x, y = eb.editorToBox(3, 2)
	assertEqual(t, x, 0)
	assertEqual(t, y, 5)
	x, y = eb.editorToBox(4, 2)
	assertEqual(t, x, 1)
	assertEqual(t, y, 5)

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
	eb := NewEditbox(0, 0, 3, 3, Options{Wrap: true})
	eb.Editor.setText(`11122233
4445
6667778

0
11122233
44`)
	eb.Editor.Cursor.X = 0
	eb.Editor.Cursor.Y = 0
	eb.moveCursorRight()
	eb.moveCursorRight()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 0)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 1)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 2)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 3)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 1)
	assertEqual(t, eb.Cursor.Y, 4)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 5)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 6)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 1)
	assertEqual(t, eb.Cursor.Y, 7)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 0)
	assertEqual(t, eb.Cursor.Y, 8)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 1)
	assertEqual(t, eb.Cursor.Y, 9)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 10)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 11)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 12)
}

func TestMoveDownOneLine(t *testing.T) {
	eb := NewEditbox(0, 0, 3, 3, Options{Wrap: true})
	eb.Editor.setText(`11122233`)
	eb.Editor.Cursor.X = 0
	eb.Editor.Cursor.Y = 0
	eb.moveCursorRight()
	eb.moveCursorRight()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 0)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 1)
	eb.moveCursorDown()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 2)
}

func TestMoveUp(t *testing.T) {
	eb := NewEditbox(0, 0, 3, 3, Options{Wrap: true})
	eb.Editor.setText(`11122233
4445
6667778

0
11122233
44`)
	eb.Editor.Cursor.X = 0
	eb.Editor.Cursor.Y = 6
	eb.moveCursorRight()
	eb.moveCursorRight()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 13)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 12)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 11)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 10)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 1)
	assertEqual(t, eb.Cursor.Y, 9)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 0)
	assertEqual(t, eb.Cursor.Y, 8)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 1)
	assertEqual(t, eb.Cursor.Y, 7)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 6)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 5)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 1)
	assertEqual(t, eb.Cursor.Y, 4)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 3)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 2)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 1)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 0)
	eb.moveCursorUp()
	eb.updateLineOffsets()
	assertEqual(t, eb.Cursor.X, 2)
	assertEqual(t, eb.Cursor.Y, 0)
}
