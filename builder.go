// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package stringx

import "strings"

// Builder extends strings.Builder from the standard library.
type Builder struct {
	strings.Builder
}

// joinLen calculates the required len for a join operation.
// add is added to the lenght of each element, to allow for extra space when required by the caller.
func joinLen(elems []string, add int, sep string) int {
	if len(elems) == 0 {
		return 0
	}

	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i]) + add
	}

	return n
}

// writeJoin uses writeElem to write each element.
func (b *Builder) writeJoin(elems []string, sep string, writeElem func(string)) {
	switch len(elems) {
	case 0:
		return
	case 1:
		writeElem(elems[0])
		return
	}

	writeElem(elems[0])
	for _, s := range elems[1:] {
		b.WriteString(sep)
		writeElem(s)
	}
}

// WriteJoin writes the elements of its first argument to the Builder.
// The separator string sep is written between elements.
func (b *Builder) WriteJoin(elems []string, sep string) int {
	n := joinLen(elems, 0, sep)
	b.Grow(n)
	b.writeJoin(elems, sep, func(s string) { b.WriteString(s) })
	return n
}

// Punction marks can be used to enclose strings.
type Punctuation [2]byte

func (p Punctuation) open() byte  { return p[0] }
func (p Punctuation) close() byte { return p[1] }

var (
	SingleQuotes   = Punctuation{'\'', '\''}
	DoubleQuotes   = Punctuation{'"', '"'}
	Backticks      = Punctuation{'`', '`'}
	RoundBrackets  = Punctuation{'(', ')'}
	SquareBrackets = Punctuation{'[', ']'}
	CurlyBrackets  = Punctuation{'{', '}'}
	AngleBrackets  = Punctuation{'<', '>'}
)

func (b *Builder) writeEnclosedString(s string, p Punctuation) {
	b.WriteByte(p.open())
	b.WriteString(s)
	b.WriteByte(p.close())
}

// WriteEnclosedString writes a string to the builder,
// with Punctuation marks before and after.
// Nothing is written if s is an empty string.
func (b *Builder) WriteEnclosedString(s string, p Punctuation) int {
	if s == "" {
		return 0
	}

	n := len(s) + len(p)
	b.Grow(n)

	b.writeEnclosedString(s, p)

	return n
}

// WriteEnclosedElems writes the elements of its first argument to the Builder.
// Each element is enclosed by Punctuation marks before and after.
// The separator string sep is written between elements.
func (b *Builder) WriteEnclosedElements(elems []string, sep string, p Punctuation) int {
	n := joinLen(elems, len(p), sep)
	b.Grow(n)
	b.writeJoin(elems, sep, func(s string) { b.writeEnclosedString(s, p) })
	return n
}

// WriteEnclosedJoin writes the elements of its first argument to the Builder.
// The separator string sep is written between elements.
// The completely joined string is enclosed by Punctuation marks before and after.
func (b *Builder) WriteEnclosedJoin(elems []string, sep string, p Punctuation) int {
	if len(elems) == 0 {
		return 0
	}

	n := joinLen(elems, 0, sep) + len(p)
	b.Grow(n)

	b.WriteByte(p.open())
	b.writeJoin(elems, sep, func(s string) { b.WriteString(s) })
	b.WriteByte(p.close())

	return n
}
