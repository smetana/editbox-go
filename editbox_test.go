package main

import (
	"testing"
    "runtime"
	"strings"
	"fmt"
	"unicode"
	"unicode/utf8"
    "reflect"
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

func (ed *Editor) setText(text string) {
    for _, s := range text {
        ed.insertRune(rune(s))
    }
}

func (ed *Editor) toLines() []string {
    lines := make([]string, len(ed.lines))
    for i, line := range ed.lines {
        lines[i] = string(line.text)
    }
    return lines
}

// ----------------------------------------------------------------------------
// Line Tests
// ----------------------------------------------------------------------------

func TestLineSimpleInsertRune(t *testing.T) {
    l := new(Line)
    l.insertRune(0, 'H')
    l.insertRune(1, 'e')
    l.insertRune(2, 'l')
    l.insertRune(3, 'l')
    l.insertRune(4, 'o')
    res := string(l.text)
    assertEqual(t, res, "Hello")
}

func TestLineInsertRune(t *testing.T) {
    l := new(Line)
    l.text = []rune("Sick")
    l.insertRune(1, 'l')
    assertEqual(t, string(l.text), "Slick")
}

func TestLineInsertPostion(t *testing.T) {
    l := new(Line)
    l.text = []rune("1")
    l.insertRune(0, '2')
    assertEqual(t, string(l.text), "21")
}

func TestLineInsertCornerCase1(t *testing.T) {
    l := new(Line)
    l.text = []rune("1")
    l.insertRune(1, '2')
    assertEqual(t, string(l.text), "12")
}

func TestLineInsertOnWrongPosition(t *testing.T) {
    defer func() {
        if r := recover(); r != "x position out of range" {
            t.Errorf("Wrong panic: %+q", r)
        }
    }()
    l := new(Line)
    l.text = []rune("1")
    l.insertRune(2, '2')
}

func TestLineInsertNewLine(t *testing.T) {
    l := new(Line)
    l.text = []rune("HelloWorld")
    l.insertRune(5, '\n')
    assertEqual(t, string(l.text), "Hello\nWorld")
}

func TestLineSplit(t *testing.T) {
    l := new(Line)
    l.text = []rune("Hello World")
    left, right := l.split(5)
    assertEqual(t, string(left.text), "Hello")
    assertEqual(t, string(right.text), " World")
}

func TestLineSplitOnWrongPosition(t *testing.T) {
    defer func() {
        if r := recover(); r != "x position out of range" {
            t.Errorf("Wrong panic: %+q", r)
        }
    }()
    l := new(Line)
    l.text = []rune("Sick")
    _,_ = l.split(10)
}

func TestLineDeleteOnWrongPosition(t *testing.T) {
    defer func() {
        if r := recover(); r != "x position out of range" {
            t.Errorf("Wrong panic: %+q", r)
        }
    }()
    l := new(Line)
    l.text = []rune("1")
    l.deleteRune(2)
}

func TestLineDelete(t *testing.T) {
    l := new(Line)
    l.text = []rune("12")
    l.deleteRune(1)
    assertEqual(t, string(l.text), "1")
    l.text = []rune("12")
    l.deleteRune(0)
    assertEqual(t, string(l.text), "2")
    l.text = []rune("")
    l.deleteRune(0)
    assertEqual(t, string(l.text), "")
}

func TestLineLastRune(t *testing.T) {
    l := new(Line)
    l.text = []rune("12")
    assertEqual(t, l.lastRune(), '2')
    l.text = []rune("12\n")
    assertEqual(t, l.lastRune(), '\n')
}

// ----------------------------------------------------------------------------
// Editor Tests
// ----------------------------------------------------------------------------

func TestEditorInsertRune(t *testing.T) {
    ed := NewEditor()
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 0)
    ed.insertRune('H')
    ed.insertRune('e')
    ed.insertRune('l')
    ed.insertRune('l')
    ed.insertRune('o')
    assertEqual(t, string(ed.Text()), "Hello")
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 5)
    ed.insertRune('!')
    assertEqual(t, string(ed.Text()), "Hello!")
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 6)
}

func TestEditorInsertOnCursorPosition(t *testing.T) {
    ed := NewEditor()
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 0)
    ed.insertRune('1')
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 1)
    assertEqual(t, string(ed.Text()), "1")
    ed.moveCursorLeft()
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 0)
    ed.insertRune('2')
    assertEqual(t, string(ed.Text()), "21")
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 1)
}

func TestEditorCurrentLine(t *testing.T) {
    ed := NewEditor()
    ed.setText("Hello World!\nSecond Line\nThird Line")
    ed.cursor.x = 2
    ed.cursor.y = 1
    assertEqual(t, string(ed.currentLine().text), "Second Line\n")
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
    assertEqual(t, len(ed.lines), 1)
    ed.insertRune('\n')
    assertEqual(t, len(ed.lines), 2)
}

func TestMoveCursorLeft(t *testing.T) {
    ed := NewEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
	assertEqual(t, len(ed.lines), 2)
    assertEqual(t, ed.cursor.y, 1)
    assertEqual(t, ed.cursor.x, 0)
    ed.moveCursorLeft()
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 2)
}

func TestBackspace(t *testing.T) {
    ed := NewEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
	ed.insertRune('1')

	assertEqual(t, len(ed.lines), 2)
    assertEqual(t, ed.cursor.y, 1)
    assertEqual(t, ed.cursor.x, 1)
    ed.deleteRuneBeforeCursor()

	assertEqual(t, ed.Text(), "12\n")
	assertEqual(t, len(ed.lines), 2)
    assertEqual(t, ed.cursor.y, 1)
    assertEqual(t, ed.cursor.x, 0)

	ed.deleteRuneBeforeCursor()
    assertEqual(t, ed.Text(), "12")
	assertEqual(t, len(ed.lines), 1)
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 2)

	ed.deleteRuneBeforeCursor()
	ed.deleteRuneBeforeCursor()

    assertEqual(t, ed.Text(), "")
	assertEqual(t, len(ed.lines), 1)
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 0)
}

func TestDeleteAtCursor(t *testing.T) {
    ed := NewEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
	ed.insertRune('3')
	ed.insertRune('\n')
	ed.insertRune('4')
	ed.insertRune('5')

	ed.cursor.x = 0
	ed.cursor.y = 1

	assertEqual(t, ed.Text(), "12\n3\n45")
	assertEqual(t, len(ed.lines), 3)
    assertEqual(t, ed.cursor.y, 1)
    assertEqual(t, ed.cursor.x, 0)

    ed.deleteRuneAtCursor()
	assertEqual(t, len(ed.lines), 3)
	assertEqual(t, ed.Text(), "12\n\n45")

    ed.deleteRuneAtCursor()
	assertEqual(t, len(ed.lines), 2)
	assertEqual(t, ed.Text(), "12\n45")

    ed.deleteRuneAtCursor()
    ed.deleteRuneAtCursor()
    ed.deleteRuneAtCursor() // No effect

	assertEqual(t, len(ed.lines), 2)
	assertEqual(t, ed.Text(), "12\n")
}

func TestMoveCursorToLineEnd(t *testing.T) {
    ed := NewEditor()
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
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 2)

	ed.cursor.x = 0
	ed.cursor.y = 2
    ed.moveCursorToLineEnd()
    assertEqual(t, ed.cursor.y, 2)
    assertEqual(t, ed.cursor.x, 2)
}


func TestMoveCursorToEmptyLine(t *testing.T) {
    ed := NewEditor()
	ed.insertRune('1')
	ed.insertRune('2')
	ed.insertRune('\n')
    assertEqual(t, ed.cursor.y, 1)
    assertEqual(t, ed.cursor.x, 0)

    ed.moveCursorVert(-1)
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 0)

    ed.moveCursorVert(+1)
    assertEqual(t, ed.cursor.y, 1)
    assertEqual(t, ed.cursor.x, 0)
}


// TODO Add tests for cursor navigation



// ----------------------------------------------------------------------------
// EditBox Tests
// ----------------------------------------------------------------------------

func TestEditorToBox(t *testing.T) {
    eb := NewEditbox(3, 3, Options{wrap:true})
    eb.editor.setText("1234567\n12\n1234\n1")
    eb.updateLineOffsets()
    assertEqual(t, eb.lineBoxY, []int{0,3,4,6})
    x,y := eb.editorToBox(0, 0)
    assertEqual(t, x, 0)
    assertEqual(t, y, 0)
    x,y = eb.editorToBox(4, 0)
    assertEqual(t, x, 1)
    assertEqual(t, y, 1)
    x,y = eb.editorToBox(6, 0)
    assertEqual(t, x, 0)
    assertEqual(t, y, 2)
    x,y = eb.editorToBox(7, 0)
    assertEqual(t, x, 1)
    assertEqual(t, y, 2)

    // TODO Wrong. There is no text there
    x,y = eb.editorToBox(8, 0)
    assertEqual(t, x, 2)
    assertEqual(t, y, 2)

    x,y = eb.editorToBox(1, 1)
    assertEqual(t, x, 1)
    assertEqual(t, y, 3)
    x,y = eb.editorToBox(2, 1)
    assertEqual(t, x, 2)
    assertEqual(t, y, 3)
    x,y = eb.editorToBox(1, 2)
    assertEqual(t, x, 1)
    assertEqual(t, y, 4)
    x,y = eb.editorToBox(3, 2)
    assertEqual(t, x, 0)
    assertEqual(t, y, 5)
    x,y = eb.editorToBox(4, 2)
    assertEqual(t, x, 1)
    assertEqual(t, y, 5)

    // TODO index out of range
    // x,y = eb.editorToBox(5, 5)
}

func TestPrevCursor(t *testing.T) {
    eb := NewEditbox(3, 3, Options{wrap:true})
    eb.updateLineOffsets()
    assertEqual(t, eb.cursor.x, 0)
    assertEqual(t, eb.cursor.y, 0)
    assertEqual(t, eb.prevCursor.x, 0)
    assertEqual(t, eb.prevCursor.y, 0)
    eb.editor.insertRune('a')
    eb.updateLineOffsets()
    assertEqual(t, eb.cursor.x, 1)
    assertEqual(t, eb.cursor.y, 0)
    assertEqual(t, eb.prevCursor.x, 0)
    assertEqual(t, eb.prevCursor.y, 0)
    eb.editor.insertRune('\n')
    eb.updateLineOffsets()
    assertEqual(t, eb.cursor.x, 0)
    assertEqual(t, eb.cursor.y, 1)
    assertEqual(t, eb.prevCursor.x, 1)
    assertEqual(t, eb.prevCursor.y, 0)
}
