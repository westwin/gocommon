package util

import (
	"errors"

	"fmt"

	"strings"
)

// KeyMarshal add - between id and code string combination
func KeyMarshal(id, code string) (rtn string) {
	return fmt.Sprintf("%s-%s", id, code)
}

// KeyUnmarshal extract original value to id and code
func KeyUnmarshal(value string) (id string, code string, err error) {
	sa := strings.Split(value, "-")
	if len(sa) != 2 {
		return id, code, errors.New("decode error")
	}
	return sa[0], sa[1], nil
}
