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

	"github.com/standupdev/strset"
)

// ParseLine parses a line in the UnicodeData.txt file returning
// the rune, the name and a set of words build from the name.
func ParseLine(line string) (rune, string, strset.Set) {
	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	name := fields[1]
	wordStr := strings.Replace(fields[1], "-", " ", -1)
	words := strset.MakeFromText(wordStr)
	if fields[10] != "" { // ➊
		name += fmt.Sprintf(" (%s)", fields[10])
		wordStr = strings.Replace(fields[10], "-", " ", -1)
		words.AddAll(strings.Fields(wordStr)...)
	}
	return rune(code), name, words
}

// filter returns a list where each item is a [3]string with the
// U+XXXX codepoint, the character (as a string) and the name of the
// Unicode characters whose name cointains all words in the query.
func filter(text io.Reader, query string) [][3]string {
	result := [][3]string{}
	query = strings.Replace(query, "-", " ", -1)
	terms := strset.MakeFromText(strings.ToUpper(query))
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		char, name, nameWords := ParseLine(line) // ➊
		if terms.SubsetOf(nameWords) {           // ➋
			result = append(result,
				[3]string{fmt.Sprintf("U+%04X", char),
					string(char), name})
		}
	}
	return result
}

// List displays the codepoint, the character and the name of the
// Unicode characters whose name cointain all words in the query.
func List(text io.Reader, query string) {
	for _, fields := range filter(text, query) {
		fmt.Printf("%s\t%s\t%s\n", fields[0], fields[1], fields[2])
	}
}

func getUCDPath() string {
	ucdPath := os.Getenv("UCD_PATH")
	if ucdPath == "" { // ➊
		user, err := user.Current()                 // ➋
		failIf(err)                                 // ➌
		ucdPath = user.HomeDir + "/UnicodeData.txt" // ➍
	}
	return ucdPath
}

func failIf(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func fetchUCD(url, path string, done chan<- bool) { // ➊
	response, err := http.Get(url)
	failIf(err)
	defer response.Body.Close()
	file, err := os.Create(path)
	failIf(err)
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	failIf(err)
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

// UCD_URL is the canonical URL for the Unicode Database.
// If unicode.org is off-line, use this alternative URL
// https://standupdev.com/data/UnicodeData.txt
const UCD_URL = "http://www.unicode.org/Public/UNIDATA/UnicodeData.txt"

func openUCD(path string) (*os.File, error) {
	ucd, err := os.Open(path)
	if os.IsNotExist(err) { // ➊
		fmt.Printf("%s not found\ndownloading %s\n", path, UCD_URL)
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
	List(ucd, query)
}
