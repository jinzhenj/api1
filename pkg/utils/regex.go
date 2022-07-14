package utils

import (
	"fmt"
	"regexp"
)

func Compile(s string, v ...interface{}) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(s, v...))
}
