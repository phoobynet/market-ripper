package ticker

import "strings"

// Clean removes any invalid characters from the ticker symbols, trims whitespace and converts to uppercase.
func Clean(ticker string) string {
	var sb strings.Builder
	for _, c := range ticker {
		if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) || c == 45 || c == 46 || c == 42 || c == 47 {
			sb.WriteRune(c)
		}
	}

	return strings.ToUpper(sb.String())
}
