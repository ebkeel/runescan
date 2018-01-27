package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
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

func getUCDPath() string {
	ucdPath := os.Getenv("UCD_PATH")
	if ucdPath == "" { // ➊
		user, err := user.Current() // ➋
		if err != nil {                // ➌
			panic(err) // não sei em que situação user.Current pode dar erro...
		}
		ucdPath = user.HomeDir + "/UnicodeData.txt" // ➍
	}
	return ucdPath
}

func fetchUCD(url, path string) error {
	response, err := http.Get(url) // ➊
	if err != nil {                // ➋
		return err
	}
	defer response.Body.Close()        // ➌
	file, err := os.Create(path) // ➍
	if err != nil {                    // ➋
		return err
	}
	defer file.Close()                    // ➎
	_, err = io.Copy(file, response.Body) // ➏
	if err != nil {                          // ➋
		return err
	}
	return nil
}

// UCD_URL fica em http://www.unicode.org/Public/UNIDATA/UnicodeData.txt
// mas unicode.org não é confiável, então esta URL alternativa pode ser usada:
// http://turing.com.br/etc/UnicodeData.txt
const UCD_URL = "http://turing.com.br/etc/UnicodeData.txt"

func openUCD(path string) (*os.File, error) {
	ucd, err := os.Open(path)
	if os.IsNotExist(err) { // ➊
		fmt.Printf("%s não encontrado\nbaixando %s\n", path, UCD_URL)
		err = fetchUCD(UCD_URL, path) // ➋
		if err != nil {
			return nil, err
		}
		ucd, err = os.Open(path) // ➌
	}
	return ucd, err // ➍
}

func main() {
	ucd, err := openUCD(getUCDPath()) // ➊
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ucd.Close()
	query := strings.Join(os.Args[1:], " ")
	List(ucd, strings.ToUpper(query))
}
