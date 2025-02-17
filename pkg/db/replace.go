package db

import (
	"strconv"
	"strings"
)

func ReplaceCharacters(text, searchPattern string) string {
	tmpCount := strings.Count(text, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		text = strings.Replace(text, searchPattern, "$"+strconv.Itoa(m), 1)
	}

	return text
}
