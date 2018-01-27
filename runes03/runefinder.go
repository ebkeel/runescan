package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseLine devolve a char e o name de uma line do UnicodeData.txt
func ParseLine(line string) (rune, string) {
	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	return rune(code), fields[1]
}

// List exibe na saída padrão o code, a char e o name dos caracteres Unicode
// cujo name contem o text da query // ➊
func List(text io.Reader, query string) {
	scanner := bufio.NewScanner(text) // ➋
	for scanner.Scan() {               // ➌
		line := scanner.Text()            // ➍
		if strings.TrimSpace(line) == "" { // ➎
			continue
		}
		char, name := ParseLine(line)    // ➏
		if strings.Contains(name, query) { // ➐
			fmt.Printf("U+%04X\t%[1]c\t%s\n", char, name)
		}
	}
}
