// Copyright 2012-2014, Rolf Veen and contributors.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer

import (
	"unicode"
)

// IsText returns true for all integers > 32 and
// are not OGDL separators (parenthesis and comma)
func isText(c int) bool {
	return c > 32 && c != '(' && c != ')' && c != ','
}

// IsEnd returns true for all integers < 32 that are not newline,
// carriage return or tab.
func isEnd(c int) bool {
	return c < 32 && c != '\t' && c != '\n' && c != '\r' 
}

// IsBreak returns true for 10 and 13 (newline and carriage return)
func isBreak(c int) bool {
	return c == 10 || c == 13
}

// IsSpace returns true for space and tab
func isSpace(c int) bool {
	return c == 32 || c == 9 
}

// ---- The following functions depend on Unicode --------

// IsLetter returns true if the given character is a letter, as per Unicode.
func isLetter(c rune) bool {
	return unicode.IsLetter(c)
}

// IsDigit returns true if the given character a numeric digit, as per Unicode.
func isDigit(c rune) bool {
	return unicode.IsDigit(c)
}

// isAlnum returns true for letters, digits and _ (as per Unicode).
func isAlnum(c rune) bool {
	return c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c) 
}
