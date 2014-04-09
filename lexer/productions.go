// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer

import (
    "fmt"
    "strings"
)

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.NextItem.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

const (
	leftDelim    = "/*"
	rightDelim   = "*/"
)

// lexText scans until an opening action delimiter, "{{".
func lexText(l *Lexer) stateFn {
	for {

		if strings.HasPrefix(l.input[l.pos:], leftDelim) {
			if l.pos > l.start {
			
				l.emit(itemText)
			}
			return lexLeftDelim
		}
	
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

// lexText scans until an opening action delimiter, "{{".
func lexInside(l *Lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], rightDelim) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexRightDelim
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

// lexLeftDelim scans the left delimiter, which is known to be present.
func lexLeftDelim(l *Lexer) stateFn {
	l.pos += len(leftDelim)
	l.emit(itemLeftDelim)
	return lexInside
}

// lexRightDelim scans the right delimiter, which is known to be present.
func lexRightDelim(l *Lexer) stateFn {
	l.pos += len(rightDelim)
	l.emit(itemRightDelim)
	return lexText// lexText scans until an opening action delimiter, "{{".
}

func lexComment(l *Lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], rightDelim) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexRightDelim
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *Lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(itemSpace)
	return lexInside
}


// lexQuote scans a quoted string.
func lexQuote(l *Lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(itemString)
	return lexInside
}



