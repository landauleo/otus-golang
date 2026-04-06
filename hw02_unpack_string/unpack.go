package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) > 0 && unicode.IsDigit([]rune(str)[0]) {
		return "", ErrInvalidString
	}

	var (
		prevElem              rune
		builder               strings.Builder
		hasPrevNonDigitSymbol bool
	)

	for _, currElem := range str {
		if unicode.IsDigit(currElem) {
			if !hasPrevNonDigitSymbol {
				return "", ErrInvalidString
			}

			digit := int(currElem - '0')
			unpackSymbols(&builder, prevElem, digit)
			hasPrevNonDigitSymbol = false
		} else {
			if hasPrevNonDigitSymbol {
				unpackSymbols(&builder, prevElem, 1)
			}
			prevElem = currElem
			hasPrevNonDigitSymbol = true
		}
	}

	//последний символ тоже не забываем
	if hasPrevNonDigitSymbol {
		unpackSymbols(&builder, prevElem, 1)
	}

	return builder.String(), nil
}

// не забываем про указатель
func unpackSymbols(builder *strings.Builder, runeElem rune, runeNumber int) {
	for i := 0; i < runeNumber; i++ {
		builder.WriteString(string(runeElem))
	}
}
