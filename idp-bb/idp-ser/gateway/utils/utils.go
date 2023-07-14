package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
)

// ParseParams parses a URL string containing application/x-www-urlencoded
// parameters and returns a map of string key-value pairs of the same
func ParseParams(str string) (map[string]string, error) {
	str, err := url.QueryUnescape(str)
	if err != nil {
		return nil, err
	}

	if strings.Contains(str, "?") {
		str = strings.Split(str, "?")[1]
	}

	if !strings.Contains(str, "=") {
		return nil, fmt.Errorf("\"%s\" contains no key-value pairs", str)
	}

	pairs := make(map[string]string)
	for _, pair := range strings.Split(string(str), "&") {
		items := strings.Split(pair, "=")
		pairs[items[0]] = items[1]
	}

	return pairs, nil
}

// ParseBasicAuthHeader decodes the Basic Auth header.
// First checks if the string contains the substring "Basic"
// and strips it off if present.
// Returns the username:password pair
func ParseBasicAuthHeader(header string) (string, string) {
	// Trimming leading and trailing whitespace
	header = strings.TrimSpace(header)

	// Check if the entire header value was used as the argument
	// eg: Basic Y2xpZW50SUQ6Y2xpZW50U2VjcmV0
	// If yes, strip off "Basic "
	if strings.HasPrefix(header, "Basic ") {
		header = strings.Split(header, " ")[1]
	}

	bytes, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		log.Println(err)
		return "", ""
	}

	str := string(bytes)
	pair := strings.Split(str, ":")
	if len(pair) != 2 {
		return "", ""
	}

	return pair[0], pair[1]
}