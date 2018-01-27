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

// ParseLine devolve a char, o name e uma slice de words que
// ocorrem no campo name de uma line do UnicodeData.txt
func ParseLine(line string) (rune, string, []string) {
	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	name := fields[1]
	words := split(fields[1])
	if fields[10] != "" { // ➊
		name += fmt.Sprintf(" (%s)", fields[10])
		for _, word := range split(fields[10]) { // ➋
			if !contains(words, word) { // ➌
				words = append(words, word) // ➍
			}
		}
	}
	return rune(code), name, words
}

func contains(slice []string, needle string) bool {
	for _, item := range slice {
		if item == needle {
			return true // ➋
		}
	}
	return false // ➌
}

func containsAll(slice []string, needles []string) bool {
	for _, needle := range needles {
		if !contains(slice, needle) {
			return false
		}
	}
	return true
}

func split(s string) []string { // ➊
	separator := func(c rune) bool { // ➋
		return c == ' ' || c == '-'
	}
	return strings.FieldsFunc(s, separator) // ➌
}

// List exibe na saída padrão o code, a char e o name dos caracteres Unicode
// cujo name contem as words da query.
func List(text io.Reader, query string) {
	terms := split(query)
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		char, name, wordsName := ParseLine(line) // ➊
		if containsAll(wordsName, terms) {           // ➋
			fmt.Printf("U+%04X\t%[1]c\t%s\n", char, name)
		}
	}
}

func main() {
	ucd, err := os.Open("UnicodeData.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() { ucd.Close() }()
	query := strings.Join(os.Args[1:], " ")
	List(ucd, strings.ToUpper(query))
}
