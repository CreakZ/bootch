package utils

import "strings"

func CleanIsbn(isbn string) string {
	return strings.ReplaceAll(isbn, "-", "")
}
