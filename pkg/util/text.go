package util

import "strings"

func SanitizeText(t string) string {
	t = strings.ReplaceAll(t, "\r\n", "\n")
	return t
}
