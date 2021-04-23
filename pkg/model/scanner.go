package model

import (
	"strings"
	"unicode"
)

type Scanner struct {
	reader *strings.Reader
}

func NewScanner(reader *strings.Reader) *Scanner {
	return &Scanner{reader: reader}
}

func (scanner *Scanner) Next() bool {
	_, _, e := scanner.reader.ReadRune()
	scanner.reader.UnreadRune()
	return e == nil
}

func (scanner *Scanner) Scan() (text string, nodeType int8) {
	r, _, _ := scanner.reader.ReadRune()
	var builder strings.Builder
	if unicode.IsLetter(r) {
		nodeType = WORD
	} else if unicode.IsNumber(r) {
		nodeType = NUMBER
	} else if r == '\n' {
		nodeType = PARAGRAPH
	} else if unicode.IsSpace(r) {
		nodeType = SPACE
	} else {
		nodeType = OTHER
	}
	builder.WriteRune(r)
	for r, _, e := scanner.reader.ReadRune(); e == nil; r, _, e = scanner.reader.ReadRune() {
		if nodeType == WORD && !(unicode.IsLetter(r) || unicode.IsNumber(r)) {
			if r == '\'' {
				next, _, e := scanner.reader.ReadRune()
				scanner.reader.UnreadRune()
				if e != nil && !unicode.IsLetter(next) {
					scanner.reader.UnreadRune()
					return builder.String(), nodeType
				} else {
					builder.WriteRune(r)
					return builder.String(), nodeType
				}
			} else {
				scanner.reader.UnreadRune()
				return builder.String(), nodeType
			}
		} else if nodeType == NUMBER && !unicode.IsNumber(r) {
			scanner.reader.UnreadRune()
			return builder.String(), nodeType
		} else if nodeType == OTHER && (unicode.IsLetter(r) || unicode.IsNumber(r) || r == '\n') {
			scanner.reader.UnreadRune()
			return builder.String(), nodeType
		} else if nodeType == PARAGRAPH && r != '\n' {
			scanner.reader.UnreadRune()
			return builder.String(), nodeType
		} else if nodeType == SPACE {
			if r == '\n' {
				nodeType = PARAGRAPH
			} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
				scanner.reader.UnreadRune()
				return builder.String(), nodeType
			} else if !unicode.IsSpace(r) {
				nodeType = OTHER
			}
		}
		builder.WriteRune(r)
	}
	return builder.String(), nodeType
}
