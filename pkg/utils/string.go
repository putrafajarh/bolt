package utils

import (
	"strings"
	"unicode"
)

// Capitalize func capitalizes the first letter of a string
func Capitalize(str string) string {
	if str == "" {
		return str
	}
	firstLetter := rune(str[0])
	return string(unicode.ToUpper(firstLetter)) + str[1:]
}

// Slug func generates a URL friendly "slug" from the given string
func Slug(str string, separator string) string {
	sepChar := "-"
	if separator != "" {
		sepChar = separator
	}

	var b strings.Builder
	prevSep := false
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(unicode.ToLower(r))
			prevSep = false
		} else {
			if !prevSep {
				b.WriteString(sepChar)
				prevSep = true
			}
		}
	}
	s := b.String()
	s = strings.TrimPrefix(s, sepChar)
	s = strings.TrimSuffix(s, sepChar)
	return s
}

// Words func limits the number of words in a string. An additional string may be passed to this method via its third argument
// to specify which string should be appended to the end of the truncated string
func Words(text string, limit uint, end string) string {
	words := strings.Fields(text)
	if uint(len(words)) <= limit {
		return text
	}

	return strings.Join(words[:limit], " ") + end
}
