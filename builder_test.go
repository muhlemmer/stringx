// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package stringx

import (
	"testing"
)

func Test_joinLen(t *testing.T) {
	type args struct {
		elems []string
		add   int
		sep   string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"nil elems, empty sep",
			args{
				nil,
				0,
				"",
			},
			0,
		},
		{
			"2 elems, empty sep",
			args{
				[]string{"foo", "bar"},
				0,
				"",
			},
			6,
		},
		{
			"2 elems, comma sep",
			args{
				[]string{"foo", "bar"},
				0,
				", ",
			},
			8,
		},
		{
			"4 elems, comma sep",
			args{
				[]string{"foo", "bar", "foo", "bar"},
				0,
				", ",
			},
			18,
		},
		{
			"4 elems, 2 add, comma sep",
			args{
				[]string{"foo", "bar", "foo", "bar"},
				2,
				", ",
			},
			26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinLen(tt.args.elems, tt.args.add, tt.args.sep); got != tt.want {
				t.Errorf("joinLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WriteJoin(t *testing.T) {
	type args struct {
		elems []string
		sep   string
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantString string
	}{
		{
			"nil elems, empty sep",
			args{
				nil,
				"",
			},
			0,
			"",
		},
		{
			"2 elems, empty sep",
			args{
				[]string{"foo", "bar"},
				"",
			},
			6,
			"foobar",
		},
		{
			"1 elems, comma sep",
			args{
				[]string{"foo"},
				", ",
			},
			3,
			"foo",
		},
		{
			"4 elems, comma sep",
			args{
				[]string{"foo", "bar", "foo", "bar"},
				", ",
			},
			18,
			"foo, bar, foo, bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Builder
			if got := b.WriteJoin(tt.args.elems, tt.args.sep); got != tt.want {
				t.Errorf("Builder.WriteJoin() = %v, want %v", got, tt.want)
			}

			if gotString := b.String(); gotString != tt.wantString {
				t.Errorf("Builder.WriteJoin() = %v, want %v", gotString, tt.wantString)
			}
		})
	}
}

func BenchmarkBuilder_WriteJoin(b *testing.B) {
	var buf Builder

	for i := 0; i < b.N; i++ {
		buf.WriteJoin([]string{"foo", "bar", "foo", "bar"}, ", ")
		buf.Reset()
	}
}

func TestBuilder_writeEnclosedString(t *testing.T) {
	var b Builder

	b.writeEnclosedString("foobar", Backticks)

	const want = "`foobar`"
	if got := b.String(); got != want {
		t.Errorf("Builder.writeQuote() = %v, want %v", got, want)
	}
}

func TestBuilder_WriteEnclosedString(t *testing.T) {
	type args struct {
		s string
		p Punctuation
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantString string
	}{
		{
			"Empty string",
			args{"", DoubleQuotes},
			0,
			"",
		},
		{
			"With string",
			args{"fabulous", DoubleQuotes},
			10,
			"\"fabulous\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Builder
			if got := b.WriteEnclosedString(tt.args.s, tt.args.p); got != tt.want {
				t.Errorf("Builder.WriteEnclosedString() = %v, want %v", got, tt.want)
			}

			if gotString := b.String(); gotString != tt.wantString {
				t.Errorf("Builder.WriteEnclosedString() = %v, want %v", gotString, tt.wantString)
			}
		})
	}
}

func BenchmarkBuilder_WriteEnclosedString(b *testing.B) {
	var buf Builder

	for i := 0; i < b.N; i++ {
		buf.WriteEnclosedString("fabulous", DoubleQuotes)
		buf.Reset()
	}
}

func TestBuilder_WriteEnclosedElements(t *testing.T) {
	type args struct {
		elems []string
		sep   string
		p     Punctuation
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantString string
	}{
		{
			"nil elems, empty sep",
			args{
				nil,
				"",
				SingleQuotes,
			},
			0,
			"",
		},
		{
			"2 elems, empty sep",
			args{
				[]string{"foo", "bar"},
				"",
				SingleQuotes,
			},
			10,
			"'foo''bar'",
		},
		{
			"1 elems, comma sep",
			args{
				[]string{"foo"},
				", ",
				DoubleQuotes,
			},
			5,
			"\"foo\"",
		},
		{
			"4 elems, comma sep",
			args{
				[]string{"foo", "bar", "foo", "bar"},
				", ",
				DoubleQuotes,
			},
			26,
			"\"foo\", \"bar\", \"foo\", \"bar\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Builder
			if got := b.WriteEnclosedElements(tt.args.elems, tt.args.sep, tt.args.p); got != tt.want {
				t.Errorf("Builder.WriteEnclosedElements() = %v, want %v", got, tt.want)
			}

			if gotString := b.String(); gotString != tt.wantString {
				t.Errorf("Builder.WriteEnclosedElements() = %v, want %v", gotString, tt.wantString)
			}
		})
	}
}

func BenchmarkBuilder_WriteEnclosedElements(b *testing.B) {
	var buf Builder

	for i := 0; i < b.N; i++ {
		buf.WriteEnclosedElements([]string{"foo", "bar", "foo", "bar"}, ", ", DoubleQuotes)
		buf.Reset()
	}
}

func TestBuilder_WriteEnclosedJoin(t *testing.T) {
	type args struct {
		elems []string
		sep   string
		p     Punctuation
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantString string
	}{
		{
			"nil elems, empty sep",
			args{
				nil,
				"",
				RoundBrackets,
			},
			0,
			"",
		},
		{
			"2 elems, empty sep",
			args{
				[]string{"foo", "bar"},
				"",
				RoundBrackets,
			},
			8,
			"(foobar)",
		},
		{
			"1 elems, comma sep",
			args{
				[]string{"foo"},
				", ",
				SquareBrackets,
			},
			5,
			"[foo]",
		},
		{
			"4 elems, comma sep",
			args{
				[]string{"foo", "bar", "foo", "bar"},
				", ",
				CurlyBrackets,
			},
			20,
			"{foo, bar, foo, bar}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Builder
			if got := b.WriteEnclosedJoin(tt.args.elems, tt.args.sep, tt.args.p); got != tt.want {
				t.Errorf("Builder.WriteEnclosedJoin() = %v, want %v", got, tt.want)
			}

			if gotString := b.String(); gotString != tt.wantString {
				t.Errorf("Builder.WriteEnclosedJoin() = %v, want %v", gotString, tt.wantString)
			}
		})
	}
}

func BenchmarkBuilder_WriteEnclosedJoin(b *testing.B) {
	var buf Builder

	for i := 0; i < b.N; i++ {
		buf.WriteEnclosedJoin([]string{"foo", "bar", "foo", "bar"}, ", ", AngleBrackets)
		buf.Reset()
	}
}
