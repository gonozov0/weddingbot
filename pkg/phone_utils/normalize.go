package phone_utils

import (
	"strings"
	"unicode"
)

func NormalizePhoneNumber(phoneNumber string) string {
	var normalized strings.Builder
	for _, r := range phoneNumber {
		if unicode.IsDigit(r) {
			normalized.WriteRune(r)
		}
	}
	return normalized.String()
}