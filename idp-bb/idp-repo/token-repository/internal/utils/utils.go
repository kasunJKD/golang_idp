package utils

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func PresentAuthScreen(w http.ResponseWriter, r *http.Request, flow int) {
	authScreenStruct := struct {
		Flow      int
	}{
		Flow:      flow,
	}

	tmpl, err := template.ParseFiles(
		"public/templates/authScreen.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.ExecuteTemplate(w, "auth", authScreenStruct)
	if err != nil {
		log.Fatal(err)
	}
}

func PresentSignInScreen(w http.ResponseWriter, r *http.Request, flow int) {
	tmpl, err := template.ParseFiles(
		"public/templates/authLogin.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.ExecuteTemplate(w, "login", flow)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowError(w http.ResponseWriter, r *http.Request, status int, title string, desc string) {
	tmpl, err := template.ParseFiles(
		"public/templates/404.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(status)
	err = tmpl.ExecuteTemplate(w, "404", struct {
		Title string
		Desc  string
	}{Title: title, Desc: desc})
	if err != nil {
		log.Fatal(err)
	}
}

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

func Bool2int(b bool) (int, error) {
	if b {
	   return 1, nil
	}
	return 0, nil
 } 

 // ByteSlices is a helper that converts an array command reply to a [][]byte.
// If err is not equal to nil, then ByteSlices returns nil, err. Nil array
// items are stay nil. ByteSlices returns an error if an array item is not a
// bulk string or nil.

type Error string
func (err Error) Error() string { return string(err) }

func ByteSlices(reply interface{}, err error) ([][]byte, error) {
	var result [][]byte
	err = sliceHelper(reply, err, "ByteSlices", func(n int) { result = make([][]byte, n) }, func(i int, v interface{}) error {
		switch v := v.(type) {
		case []byte:
			result[i] = v
			return nil
		case Error:
			return v
		default:
			return fmt.Errorf("unexpected element type for ByteSlices, got type %T", v)
		}
	})
	return result, err
}

func sliceHelper(reply interface{}, err error, name string, makeSlice func(int), assign func(int, interface{}) error) error {
	if err != nil {
		return err
	}
	switch reply := reply.(type) {
	case []interface{}:
		makeSlice(len(reply))
		for i := range reply {
			if reply[i] == nil {
				continue
			}
			if err := assign(i, reply[i]); err != nil {
				return err
			}
		}
		return nil
	case nil:
		return fmt.Errorf("nil error")
	case Error:
		return reply
	}
	return fmt.Errorf("unexpected type for %s, got type %T", name, reply)
}