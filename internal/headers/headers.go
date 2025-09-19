package headers

import (
	"fmt"
	"strings"
)

const crlf = "\r\n"

const validCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func checkForInvalidCharacters(r rune) rune {
	if !strings.Contains(validCharacters, string(r)) {
		return -1
	}
	return r
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headerText := string(data)
	crlfIndex := strings.Index(headerText, crlf)
	if crlfIndex == -1 {
		return 0, false, nil
	}
	if crlfIndex == 0 {
		return 2, true, nil
	}
	header := strings.SplitN(headerText, crlf, 2)[0]
	trimmed_header := strings.TrimSpace(header)
	values := strings.SplitN(trimmed_header, ":", 2)
	key := values[0]
	if strings.Contains(key, " ") {
		return 0, false, fmt.Errorf("invalid spacing header: %v", trimmed_header)
	}
	checked_key := strings.Map(checkForInvalidCharacters, key)
	if key != checked_key {
		return 0, false, fmt.Errorf("invalid character in field name: %v", key)
	}
	value := values[1]
	if val, ok := h[strings.ToLower(key)]; ok {
		h[strings.ToLower(key)] = val + ", " + strings.TrimSpace(value)
	} else {
		h[strings.ToLower(key)] = strings.TrimSpace(value)
	}
	return len(header) + 2, false, nil
}
