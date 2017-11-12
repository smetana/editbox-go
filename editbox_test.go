package main

import (
	"testing"
)

func assertEqual(t *testing.T, actual, expected interface{}) {
    if actual != expected {
        t.Errorf("Expected (%T)%v got (%T)%v", expected, expected, actual, actual)
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

func TestLineInsertOnWrongPosition(t *testing.T) {
    defer func() {
        if r := recover(); r != "position out of range" {
            t.Errorf("Wrong panic: %+q", r)
        }
    }()
    l := new(Line)
    l.text = []rune("Sick")
    l.insertRune(10, 'l')
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

func TestEditorInsertRune(t *testing.T) {
    ed := NewEditor(5, 5)
    ed.insertRune('H')
    ed.insertRune('e')
    ed.insertRune('l')
    ed.insertRune('l')
    ed.insertRune('o')
    assertEqual(t, string(ed.Text()), "Hello")
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 5)
    ed.insertRune('W')
    assertEqual(t, string(ed.Text()), "HelloW")
    assertEqual(t, ed.cursor.y, 0)
    assertEqual(t, ed.cursor.x, 5)
}

