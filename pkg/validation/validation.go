package validation

import (
	"strings"
	"unicode"
)

func IsIsbn10Valid(isbn string) bool {
	if strings.ContainsFunc(isbn, func(r rune) bool {
		return !unicode.IsDigit(r)
	}) {
		return false
	}
	return isIsbn10SumValid(isbn)
}

func IsIsbn13Valid(isbn string) bool {
	if strings.ContainsFunc(isbn, func(r rune) bool {
		return !unicode.IsDigit(r)
	}) {
		return false
	}
	return isIsbn13SumValid(isbn)
}

func isIsbn10SumValid(isbn string) bool {
	if len(isbn) != 10 {
		return false
	}

	sum := 0
	for i, c := range isbn {
		charValue := int(c - '0')
		sum += charValue * (len(isbn) - i)
	}
	return sum%11 == 0
}

func isIsbn13SumValid(isbn string) bool {
	if len(isbn) != 13 {
		return false
	}

	sum := 0
	last := int(isbn[len(isbn)-1] - '0')
	for i, c := range isbn[:len(isbn)-1] {
		charValue := int(c - '0')
		if i%2 != 0 {
			charValue *= 3
		}
		sum += charValue
	}
	return (10-(sum%10))%10 == last
}
