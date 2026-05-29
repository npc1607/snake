package helper

import (
	"net"
	"unicode/utf8"
)

// NetworkForAddress returns the best TCP network for a host:port address.
func NetworkForAddress(address string) string {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return "tcp"
	}
	if host == "" {
		return "tcp"
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return "tcp"
	}
	if ip.To4() != nil {
		return "tcp4"
	}

	return "tcp6"
}

// CopySlice returns a shallow copy of a slice.
func CopySlice[T any](values []T) []T {
	copied := make([]T, len(values))
	copy(copied, values)

	return copied
}

// TruncateText returns text shortened to width using a trailing dot.
func TruncateText(text string, width int) string {
	if width <= 0 || utf8.RuneCountInString(text) <= width {
		return text
	}
	if width == 1 {
		return "."
	}

	runes := []rune(text)
	return string(runes[:width-1]) + "."
}
