package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/standupdev/wordset"
)

const lineLetterA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;"

func TestParseLine(t *testing.T) {
	rune, name, words := ParseLine(lineLetterA) // ➊
	if rune != 'A' {
		t.Errorf("Want: 'A'; got: %q", rune)
	}
	const nameA = "LATIN CAPITAL LETTER A"
	if name != nameA {
		t.Errorf("Want: %q; got: %q", nameA, name)
	}
	wordsA := wordset.MakeFromText(nameA) // ➋
	if !wordsA.Equal(words) {             // ➌
		t.Errorf("\n\tWant: %q\n\tgot: %q", wordsA, words) // ➍
	}
}

func TestParseLineWithHyphenAndField10(t *testing.T) {
	var testCases = []struct { // ➊
		line  string
		char  rune
		name  string
		words wordset.Set
	}{ // ➋
		{"0021;EXCLAMATION MARK;Po;0;ON;;;;;N;;;;;",
			'!', "EXCLAMATION MARK",
			wordset.MakeFromText("EXCLAMATION MARK")},
		{"002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;",
			'-', "HYPHEN-MINUS",
			wordset.MakeFromText("HYPHEN MINUS")},
		{"0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;",
			'\'', "APOSTROPHE (APOSTROPHE-QUOTE)",
			wordset.MakeFromText("APOSTROPHE QUOTE")},
	}

	for _, tc := range testCases { // ➌
		t.Run("case "+string(tc.char), func(t *testing.T) {
			char, name, words := ParseLine(tc.line) // ➍
			if char != tc.char || name != tc.name ||
				!words.Equal(tc.words) {
				t.Errorf("\nParseLine(%q)\n-> (%q, %q, %q)", // ➎
					tc.line, char, name, words)
			}
		})
	}
}

const lines3Dto43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

func TestFilter(t *testing.T) {
	var testCases = []struct { // ➊
		query string
		want  [][3]string
	}{ // ➋
		{"ZZZZZZ", [][3]string{}},
		{"MARK", [][3]string{
			{"U+003F", "?", "QUESTION MARK"},
		}},
		{"SIGN", [][3]string{
			{"U+003D", "=", "EQUALS SIGN"},
			{"U+003E", ">", "GREATER-THAN SIGN"},
		}},
		{"GREATER-THAN", [][3]string{
			{"U+003E", ">", "GREATER-THAN SIGN"},
		}},
	}
	for _, tc := range testCases { // ➌
		t.Run(tc.query, func(t *testing.T) {
			text := strings.NewReader(lines3Dto43)
			got := filter(text, tc.query) // ➍
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("query: %q\twant: %q\tgot: %q", // ➎
					tc.query, tc.want, got)
			}
		})
	}
}

func ExampleList() {
	text := strings.NewReader(lines3Dto43)
	List(text, "MARK")
	// Output: U+003F	?	QUESTION MARK
}

func ExampleList_2Results() {
	text := strings.NewReader(lines3Dto43)
	List(text, "SIGN")
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}

func ExampleList_2Words() {
	text := strings.NewReader(lines3Dto43)
	List(text, "CAPITAL LATIN")
	// Output:
	// U+0041	A	LATIN CAPITAL LETTER A
	// U+0042	B	LATIN CAPITAL LETTER B
	// U+0043	C	LATIN CAPITAL LETTER C
}

func Example() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "cruzeiro"}
	main()
	// Output:
	// U+20A2	₢	CRUZEIRO SIGN
}

func Example_2WordQuery() { // ➊
	oldArgs := os.Args // ➋
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "cat", "smiling"}
	main() // ➌
	// Output:
	// U+1F638	😸	GRINNING CAT FACE WITH SMILING EYES
	// U+1F63A	😺	SMILING CAT FACE WITH OPEN MOUTH
	// U+1F63B	😻	SMILING CAT FACE WITH HEART-SHAPED EYES
}

func Example_queryWithHiphenAndField10() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "quote"}
	main()
	// Output:
	// U+0027	'	APOSTROPHE (APOSTROPHE-QUOTE)
	// U+2358	⍘	APL FUNCTIONAL SYMBOL QUOTE UNDERBAR
	// U+235E	⍞	APL FUNCTIONAL SYMBOL QUOTE QUAD
}

func restore(nameVar, value string, existed bool) {
	if existed {
		os.Setenv(nameVar, value)
	} else {
		os.Unsetenv(nameVar)
	}
}

func TestGetUCDPath_isSet(t *testing.T) {
	pathBefore, existed := os.LookupEnv("UCD_PATH")                           // ➊
	defer restore("UCD_PATH", pathBefore, existed)                            // ➋
	ucdPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano()) // ➌
	os.Setenv("UCD_PATH", ucdPath)                                            // ➍
	got := getUCDPath()                                                       // ➎
	if got != ucdPath {
		t.Errorf("getUCDPath() [set]\nwant: %q; got: %q", ucdPath, got)
	}
}

func TestGetUCDPath_default(t *testing.T) {
	pathBefore, existed := os.LookupEnv("UCD_PATH")
	defer restore("UCD_PATH", pathBefore, existed)
	os.Unsetenv("UCD_PATH")             // ➊
	ucdPathSuffix := "/UnicodeData.txt" // ➋
	got := getUCDPath()
	if !strings.HasSuffix(got, ucdPathSuffix) { // ➌
		t.Errorf("getUCDPath() [default]\nwant (sufixo): %q; got: %q", ucdPathSuffix, got)
	}
}

func TestOpenUCD_local(t *testing.T) {
	ucdPath := getUCDPath()
	ucd, err := openUCD(ucdPath)
	if err != nil {
		t.Errorf("openUCD(%q):\n%v", ucdPath, err)
	}
	ucd.Close()
}

func TestFetchUCD(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(lines3Dto43))
		}))
	defer srv.Close()

	ucdPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano())
	done := make(chan bool)             // ➊
	go fetchUCD(srv.URL, ucdPath, done) // ➋
	_ = <-done                          // ➌
	ucd, err := os.Open(ucdPath)
	if os.IsNotExist(err) {
		t.Errorf("fetchUCD did not save:%v\n%v", ucdPath, err)
	}
	ucd.Close()
	os.Remove(ucdPath)
}

func TestOpenUCD_remote(t *testing.T) {
	if testing.Short() { // ➊
		t.Skip("skipped test [-test.short option]") // ➋
	}
	ucdPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano()) // ➌
	ucd, err := openUCD(ucdPath)
	if err != nil {
		t.Errorf("openUCD(%q):\n%v", ucdPath, err)
	}
	ucd.Close()
	os.Remove(ucdPath)
}
