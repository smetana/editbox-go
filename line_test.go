package editbox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
