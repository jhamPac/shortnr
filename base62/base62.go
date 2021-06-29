package base62

import "strings"

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
