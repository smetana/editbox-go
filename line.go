package editbox

type line struct {
	text []rune
}

func (l *line) checkXPosition(x int) {
	if x < 0 || x > len(l.text) {
		panic("x position out of range")
	}
}

func (l *line) insertRune(pos int, r rune) {
	l.checkXPosition(pos)
	// Append
	if pos == len(l.text) {
		l.text = append(l.text, r)
		// Insert
	} else {
		l.text = append(l.text, rune(0))
		copy(l.text[pos+1:], l.text[pos:])
		l.text[pos] = r
	}
}

func (l *line) split(pos int) (left, right *line) {
	l.checkXPosition(pos)
	left, right = l, new(line)
	right.text = make([]rune, len(l.text)-pos)
	copy(right.text, l.text[pos:len(l.text)])
	left.text = left.text[:pos]
	return
}

func (l *line) deleteRune(pos int) rune {
	l.checkXPosition(pos)
	if pos < len(l.text) {
		r := l.text[pos]
		copy(l.text[pos:], l.text[pos+1:])
		l.text[len(l.text)-1] = rune(0)
		l.text = l.text[:len(l.text)-1]
		return r
	} else {
		return rune(0)
	}
}

func (l *line) lastRune() rune {
	if len(l.text) == 0 {
		return 0
	} else {
		return l.text[len(l.text)-1]
	}
}

func (l *line) lastRuneX() int {
	if l.lastRune() == '\n' {
		return (len(l.text) - 1)
	} else {
		return (len(l.text))
	}
}
