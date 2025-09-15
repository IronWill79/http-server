package headers

import (
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headerText := string(data)
	if !strings.Contains(headerText, crlf) {
		return 0, false, nil
	}
	if headerText[:2] == crlf {
		return 0, true, nil
	}
	header := strings.SplitN(headerText, crlf, 1)[0]
	trimmed_header := strings.TrimSpace(header)
	values := strings.SplitN(trimmed_header, ":", 2)
	key := values[0]
	value := values[1]
	if strings.Contains(key, " ") {
		return 0, false, fmt.Errorf("invalid spacing header: %v", trimmed_header)
	}
	h[strings.TrimSpace(key)] = strings.TrimSpace(value)
	return len(header) - 2, false, nil
}
