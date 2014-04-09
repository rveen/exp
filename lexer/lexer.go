// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer

import (
	"fmt"
	"strings"
	"bytes"
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// lexer holds the state of the scanner.
type Lexer struct {
	name       string    // the name of the input; used only for error reports
	input      string    // the string being scanned
	state      stateFn   // the next lexing function to enter
	pos        int       // current position in the input
	start      int       // start position of this item
	lastPos    int       // position of most recent item returned by nextItem
	items      chan item // channel of scanned items
}

// lex creates a new scanner for the input string.
func NewLexer(name, input string) *Lexer {
    l := &Lexer{
        name:  name,
        input: input,
        state: lexText,
        items: make(chan item, 2),
    }
    return l
}

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType // The type of this item.
	pos int      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	default:
		return fmt.Sprintf("%q", i.val)
	}
}

func (i item) IsEOF() bool {
    return i.typ == itemEOF
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError        itemType = iota // error occurred; value is text of error
	itemEOF
	itemText
	itemSpace
	itemString
	itemLeftDelim
	itemRightDelim
)

// next returns the next byte in the input.
func (l *Lexer) next() int {
	if int(l.pos) >= len(l.input) {
		return eof
	}
	r := l.input[l.pos]
	l.pos++
	
	return int(r)
}

// peek returns but does not consume the next byte in the input.
func (l *Lexer) peek() int {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one byte. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos--
}

// emit passes an item back to the client.
func (l *Lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next byte if it's from the valid set.
func (l *Lexer) accept(valid []byte) bool {
    if bytes.IndexByte(valid,byte(l.next())) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of bytes from the valid set.
func (l *Lexer) acceptRun(valid []byte) {
	for bytes.IndexByte(valid,byte(l.next())) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *Lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

// NextItem returns the next item from the input.
func (l *Lexer) NextItem() item {
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


