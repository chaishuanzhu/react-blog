package service

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	reCodeFence  = regexp.MustCompile("(?s)```.*?```")
	reInlineCode = regexp.MustCompile("`([^`]*)`")
	reImage      = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	reLink       = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	reHTMLTag    = regexp.MustCompile(`<[^>]+>`)
	reLinePrefix = regexp.MustCompile(`(?m)^\s{0,3}(#{1,6}\s+|>\s?|[-*+]\s+|\d+\.\s+)`)
	reEmphasis   = regexp.MustCompile(`[*_~]{1,3}`)
	reSpaces     = regexp.MustCompile(`\s+`)
)

// Summarize turns markdown into a plain-text excerpt of at most maxRunes characters.
func Summarize(markdown string, maxRunes int) string {
	s := reCodeFence.ReplaceAllString(markdown, " ")
	s = reImage.ReplaceAllString(s, " ")
	s = reLink.ReplaceAllString(s, "$1")
	s = reInlineCode.ReplaceAllString(s, "$1")
	s = reHTMLTag.ReplaceAllString(s, " ")
	s = reLinePrefix.ReplaceAllString(s, "")
	s = reEmphasis.ReplaceAllString(s, "")
	s = strings.TrimSpace(reSpaces.ReplaceAllString(s, " "))

	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	return string([]rune(s)[:maxRunes]) + "…"
}

// likeContains builds a LIKE pattern matching s as a literal substring.
func likeContains(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(s) + "%"
}
