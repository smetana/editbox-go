package main

import (
	"testing"
)

func assertEqual(t *testing.T, actual, expected interface{}) {
    if actual != expected {
        t.Errorf("Expected (%T)%+q got (%T)%+q", expected, expected, actual, actual)
    }
}

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
        if r := recover(); r != "position out of range" {
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
        if r := recover(); r != "position out of range" {
            t.Errorf("Wrong panic: %+q", r)
        }
    }()
    l := new(Line)
    l.text = []rune("Sick")
    _,_ = l.split(10)
}

func TestLineDeleteOnWrongPosition(t *testing.T) {
    defer func() {
        if r := recover(); r != "position out of range" {
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

func TestEditorInsertRune(t *testing.T) {
    ed := NewEditor(5, 5)
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
    ed := NewEditor(5, 5)
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

// TODO Add tests for cursor navigation
