package parser

import (
	"strings"
	"unicode"

	"github.com/hslapr/leglog/pkg/data"
)

type Parser interface {
	Parse(text string) *data.Text
}

// Parse parses a sentence
func Parse(text string, lang string) *data.Text {
	switch lang {
	case "en":
		return parseDefault(text, data.ENGLISH)
	case "it":
		return parseItalian(text)
	}
	return nil
}

func parseItalian(text string) *data.Text {
	var t data.Text
	t.Language = data.ITALIAN
	var reader = strings.NewReader(strings.ReplaceAll(text, "\r\n", "\n"))

	var builder strings.Builder
	var c int
	var isWord bool

	for {
		builder.Reset()
		isWord = false
		r, _, e := reader.ReadRune()
		if e != nil {
			return &t
		}
		if r == '\n' {
			t.Append(data.NewSegment("\n", isWord))
			continue
		}
		if unicode.IsLetter(r) {
			c = WORD
			isWord = true
		} else if unicode.IsNumber(r) {
			c = NUMBER
		} else {
			c = OTHER
		}
		for {
			builder.WriteRune(r)
			r, _, e = reader.ReadRune()
			if e != nil {
				t.Append(data.NewSegment(builder.String(), isWord))
				return &t
			}
			if c == WORD && !unicode.IsLetter(r) {
				if unicode.IsNumber(r) {
					continue
				} else if r == '\'' {
					next, _, e := reader.ReadRune()
					reader.UnreadRune()
					if e != nil && !unicode.IsLetter(next) {
						reader.UnreadRune()
						t.Append(data.NewSegment(builder.String(), isWord))
						break
					} else {
						builder.WriteRune(r)
						t.Append(data.NewSegment(builder.String(), isWord))
						break
					}
				} else {
					reader.UnreadRune()
					t.Append(data.NewSegment(builder.String(), isWord))
					break
				}
			}
			if c == NUMBER && !unicode.IsNumber(r) {
				reader.UnreadRune()
				t.Append(data.NewSegment(builder.String(), isWord))
				break
			}
			if c == OTHER && (unicode.IsLetter(r) || unicode.IsNumber(r) || r == '\n') {
				reader.UnreadRune()
				t.Append(data.NewSegment(builder.String(), isWord))
				break
			}
		}
	}
}

func parseDefault(text string, language data.Language) *data.Text {
	var t data.Text
	t.Language = language
	var reader = strings.NewReader(strings.ReplaceAll(text, "\r\n", "\n"))
	scanner := Scanner{reader: reader}

	for {
		s, isWord, e := scanner.Scan()
		if e != nil {
			break
		}
		t.Append(data.NewSegment(s, isWord))
	}
	return &t
}
