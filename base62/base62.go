package base62

import (
	"errors"
	"math"
	"strings"
)

const (
	alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length       = uint64(len(alphanumeric))
)

func Encode(number uint64) string {
	var encodeBuilder strings.Builder
	encodeBuilder.Grow(11)

	for {
		if number <= 0 {
			break
		}

		encodeBuilder.WriteByte(alphanumeric[(number % length)])

		number = number / length
	}

	return encodeBuilder.String()
}

func Decode(encoded string) (uint64, error) {
	var number uint64

	for i, symbol := range encoded {
		alphaPosition := strings.IndexRune(alphanumeric, symbol)

		if alphaPosition == -1 {
			return uint64(alphaPosition), errors.New("invalid character: " + string(symbol))
		}

		number += uint64(alphaPosition) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number, nil
}
