// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer

import (
	"strings"
	"unicode/utf8"
)

// stateFnU represents the state of the scanner as a function that returns the next state.
type stateFnU func(*ULexer) stateFnU

// ULexer holds the state of the scanner.
type ULexer struct {
	name       string    // the name of the input; used only for error reports
	input      string    // the string being scanned
	state      stateFnU   // the next lexing function to enter
	pos        int       // current position in the input
	start      int       // start position of this item
	lastPos    int       // position of most recent item returned by nextItem
	items      chan item // channel of scanned items
	width      int       // width of last rune read from input
}

// NewUlexer creates a new scanner for the input string.
func NewULexer(name, input string) *ULexer {
    l := &ULexer{
        name:  name,
        input: input,
        state: nil, // XXX
        items: make(chan item, 2),
    }
    return l
}

// next returns the next rune in the input.
func (l *ULexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *ULexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *ULexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *ULexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *ULexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *ULexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *ULexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *ULexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

// nextItem returns the next item from the input.
func (l *ULexer) NextItem() item {
    for {
        select {
        case item := <-l.items:
            return item
        default:
            l.state = l.state(l)
        }
    }
    panic("not reached")
}


