package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func lowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func upperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func CamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	if strings.Contains(s, "_") {
		parts := strings.Split(strings.ToLower(s), "_")
		r := ""
		for _, p := range parts {
			r += upperFirst(p)
		}
		return lowerFirst(r)
	}
	if strings.ToUpper(s) == s {
		return strings.ToLower(s)
	}
	return lowerFirst(s)
}

func PascalCase(s string) string {
	return upperFirst(CamelCase(s))
}

func Indent(s string) string {
	if len(s) == 0 {
		return s
	}
	enl := false
	if s[len(s)-1] == '\n' {
		enl = true
		s = s[:len(s)-1]
	}
	var ind []string
	for _, line := range strings.Split(s, "\n") {
		ind = append(ind, fmt.Sprintf("  %s", line))
	}
	r := strings.Join(ind, "\n")
	if enl {
		r += "\n"
	}
	return r
}

func ToJson(o interface{}) string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}
