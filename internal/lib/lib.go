package lib

import (
	"math"
	"slices"
	"strings"
)

const ALPHANUMERIC = "abcdefghijklmnopqrstuvwxyzABCDEFGHTJKLMNOPQRSTUVWXYZ0123456789"

func Encode(number int) string {
	digits := toBase62Digits(number)
	return digitsToString(digits)
}

func Decode(s string) int {
	totalChars := float64(len(ALPHANUMERIC))
	number := 0
	for idx, char := range s {
		digit := strings.IndexRune(ALPHANUMERIC, char)
		power := float64(len(s) - idx - 1)
		digitFactor := int(math.Pow(totalChars, power))
		number += digit * digitFactor
	}
	return number
}

func toBase62Digits(number int) []int {
	var digits []int
	totalChars := len(ALPHANUMERIC)
	for number > 0 {
		remainder := number % totalChars
		digits = append(digits, remainder)
		number /= totalChars
	}
	slices.Reverse(digits)
	return digits
}

func digitsToString(digits []int) string {
	var sb strings.Builder
	for _, digit := range digits {
		sByte := ALPHANUMERIC[digit]
		sb.WriteByte(sByte)
	}
	return sb.String()
}
