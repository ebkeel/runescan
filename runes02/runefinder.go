package main

import ( // ➊
	"strconv"
	"strings"
)

func ParseLine(ucdLine string) (rune, string) {
	fields := strings.Split(ucdLine, ";")            // ➋
	code, _ := strconv.ParseInt(fields[0], 16, 32) // ➌
	return rune(code), fields[1]                   // ➍
}
