package domain

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

var (
	errInvalidURL = errors.New("invalid URL")
)

func ExtractDomain(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", errInvalidURL
	}
	if parsedURL.Hostname() == "" {
		return "", errInvalidURL
	}
	// Split the hostname into parts separated by dots
	parts := strings.Split(parsedURL.Hostname(), ".")
	if len(parts) > 2 && !isIP(parts[len(parts)-2]) {
		// If the last two parts of the hostname are not an IP address, join the last two parts
		return strings.Join(parts[len(parts)-2:], "."), nil
	}
	// Otherwise, return the full hostname
	return parsedURL.Hostname(), nil
}

// isIP checks whether the given string is a valid IP address
func isIP(str string) bool {
	return net.ParseIP(str) != nil
}
