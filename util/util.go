package util

import (
	"strings"
)

func ParseArgs(text string) string {
	query := strings.SplitAfterN(text, " ", 2)
	if len(query) < 2 {
		return ""
	}
	return query[1]
}
