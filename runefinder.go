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
	"time"

	"github.com/standupdev/wordset"
)

// ParseLine devolve a rune, o name e uma slice de words que
// ocorrem no campo name de uma line do UnicodeData.txt
func ParseLine(line string) (rune, string, wordset.Set) {
	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	name := fields[1]
	wordStr := strings.Replace(fields[1], "-", " ", -1)
	words := wordset.MakeFromText(wordStr)
	if fields[10] != "" { // ➊
		name += fmt.Sprintf(" (%s)", fields[10])
		wordStr = strings.Replace(fields[10], "-", " ", -1)
		words.Update(wordset.MakeFromText(wordStr))
	}
	return rune(code), name, words
}

// List exibe na saída padrão o code, a rune e o name dos caracteres Unicode
// cujo name contem as words da query.
func List(text io.Reader, query string) {
	terms := wordset.MakeFromText(query)
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		rune, name, wordsName := ParseLine(line) // ➊
		if terms.IsSubSetOf(wordsName) {       // ➋
			fmt.Printf("U+%04X\t%[1]c\t%s\n", rune, name)
		}
	}
}

func getUCDPath() string {
	ucdPath := os.Getenv("UCD_PATH")
	if ucdPath == "" { // ➊
		user, err := user.Current()                 // ➋
		terminarSe(err)                             // ➌
		ucdPath = user.HomeDir + "/UnicodeData.txt" // ➍
	}
	return ucdPath
}

func terminarSe(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func fetchUCD(url, path string, done chan<- bool) { // ➊
	response, err := http.Get(url)
	terminarSe(err)
	defer response.Body.Close()
	file, err := os.Create(path)
	terminarSe(err)
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	terminarSe(err)
	done <- true // ➋
}

func progress(done <-chan bool) { // ➊
	for { // ➋
		select { // ➌
		case <-done: // ➍
			fmt.Println()
			return
		default: // ➎
			fmt.Print(".")
			time.Sleep(150 * time.Millisecond)
		}
	}
}

// UCD_URL fica em http://www.unicode.org/Public/UNIDATA/UnicodeData.txt
// mas unicode.org não é confiável, então esta URL alternativa pode ser usada:
// http://turing.com.br/etc/UnicodeData.txt
const UCD_URL = "http://turing.com.br/etc/UnicodeData.txt"

func openUCD(path string) (*os.File, error) {
	ucd, err := os.Open(path)
	if os.IsNotExist(err) { // ➊
		fmt.Printf("%s não encontrado\nbaixando %s\n", path, UCD_URL)
		done := make(chan bool)          // ➊
		go fetchUCD(UCD_URL, path, done) // ➋
		progress(done)                   // ➌
		ucd, err = os.Open(path)         // ➌
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
