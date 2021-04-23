package parser

import (
	"strings"
	"unicode"
)

type Scanner struct {
	reader *strings.Reader
}

const (
	WORD   = iota
	NUMBER = iota
	OTHER  = iota
)

func (scanner *Scanner) Scan() (s string, isWord bool, err error) {
	var builder strings.Builder
	var t int
	isWord = false

	r, _, e := scanner.reader.ReadRune()
	if e != nil {
		return "", isWord, e
	}

	if r == '\n' {
		return "\n", isWord, nil
	}

	if unicode.IsLetter(r) {
		t = WORD
		isWord = true
	} else if unicode.IsNumber(r) {
		t = NUMBER
	} else {
		t = OTHER
	}

	for {
		builder.WriteRune(r)
		r, _, e = scanner.reader.ReadRune()
		if e != nil {
			return builder.String(), isWord, e
		}
		if t == WORD && !unicode.IsLetter(r) {
			scanner.reader.UnreadRune()
			return builder.String(), isWord, nil
		}
		if t == NUMBER && !unicode.IsNumber(r) {
			scanner.reader.UnreadRune()
			return builder.String(), isWord, nil
		}
		if t == OTHER && (unicode.IsLetter(r) || unicode.IsNumber(r) || r == '\n') {
			scanner.reader.UnreadRune()
			return builder.String(), isWord, nil
		}
	}
}
