package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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
// cujo name contem o text da query.
func List(text io.Reader, query string) {
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		char, name := ParseLine(line)
		if strings.Contains(name, query) {
			fmt.Printf("U+%04X\t%[1]c\t%s\n", char, name)
		}
	}
}

func main() { // ➊
	ucd, err := os.Open("UnicodeData.txt") // ➋
	if err != nil {                        // ➌
		log.Fatal(err.Error()) // ➍
	}
	defer ucd.Close()                          // ➎
	query := strings.Join(os.Args[1:], " ") // ➏
	List(ucd, strings.ToUpper(query))     // ➐
}
