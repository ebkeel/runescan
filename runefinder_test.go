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
)

const lineLetterA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;"

const lines3Dto43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

func TestParseLine(t *testing.T) {
	rune, name, words := ParseLine(lineLetterA) // ➊
	if rune != 'A' {
		t.Errorf("Esperado: 'A'; got: %q", rune)
	}
	const nameA = "LATIN CAPITAL LETTER A"
	if name != nameA {
		t.Errorf("Esperado: %q; got: %q", nameA, name)
	}
	wordsA := []string{"LATIN", "CAPITAL", "LETTER", "A"} // ➋
	if !reflect.DeepEqual(words, wordsA) {             // ➌
		t.Errorf("\n\tEsperado: %q\n\tgot: %q", wordsA, words) // ➍
	}
}

func TestParseLineWithHyphenAndField10(t *testing.T) {
	var testCases = []struct { // ➊
		line    string
		rune     rune
		name     string
		words []string
	}{ // ➋
		{"0021;EXCLAMATION MARK;Po;0;ON;;;;;N;;;;;",
			'!', "EXCLAMATION MARK", []string{"EXCLAMATION", "MARK"}},
		{"002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;",
			'-', "HYPHEN-MINUS", []string{"HYPHEN", "MINUS"}},
		{"0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;",
			'\'', "APOSTROPHE (APOSTROPHE-QUOTE)", []string{"APOSTROPHE", "QUOTE"}},
	}
	for _, tc := range testCases { // ➌
		rune, name, words := ParseLine(tc.line) // ➍
		if rune != tc.rune || name != tc.name ||
			!reflect.DeepEqual(words, tc.words) {
			t.Errorf("\nParseLine(%q)\n-> (%q, %q, %q)", // ➎
				tc.line, rune, name, words)
		}
	}
}

func TestContains(t *testing.T) {
	testCases := []struct { // ➊
		slice     []string
		needle string
		want  bool
	}{ // ➋
		{[]string{"A", "B"}, "B", true},
		{[]string{}, "A", false},
		{[]string{"A", "B"}, "Z", false}, // ➌
	} // ➍
	for _, tc := range testCases { // ➎
		got := contains(tc.slice, tc.needle) // ➏
		if got != tc.want {                 // ➐
			t.Errorf("contains(%#v, %#v) want: %v; got: %v",
				tc.slice, tc.needle, tc.want, got) // ➑
		}
	}
}

func TestContainsAll(t *testing.T) {
	testCases := []struct { // ➊
		slice      []string
		needles []string
		want   bool
	}{ // ➋
		{[]string{"A", "B"}, []string{"B"}, true},
		{[]string{}, []string{"A"}, false},
		{[]string{"A"}, []string{}, true}, // ➌
		{[]string{"A", "B"}, []string{"Z"}, false},
		{[]string{"A", "B", "C"}, []string{"A", "C"}, true},
		{[]string{"A", "B", "C"}, []string{"A", "Z"}, false},
		{[]string{"A", "B"}, []string{"A", "B", "C"}, false},
	}
	for _, tc := range testCases {
		got := containsAll(tc.slice, tc.needles) // ➍
		if got != tc.want {
			t.Errorf("containsAll(%#v, %#v)\nwant: %v; got: %v",
				tc.slice, tc.needles, tc.want, got) // ➎
		}
	}
}

func TestSplit(t *testing.T) {
	testCases := []struct {
		text    string
		want []string
	}{
		{"A", []string{"A"}},
		{"A B", []string{"A", "B"}},
		{"A B-C", []string{"A", "B", "C"}},
	}
	for _, tc := range testCases {
		got := split(tc.text)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("split(%q)\nwant: %#v; got: %#v",
				tc.text, tc.want, got)
		}
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
	pathBefore, existed := os.LookupEnv("UCD_PATH")                            // ➊
	defer restore("UCD_PATH", pathBefore, existed)                           // ➋
	ucdPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano()) // ➌
	os.Setenv("UCD_PATH", ucdPath)                                            // ➍
	got := getUCDPath()                                                  // ➎
	if got != ucdPath {
		t.Errorf("getUCDPath() [setado]\nwant: %q; got: %q", ucdPath, got)
	}
}

func TestGetUCDPath_default(t *testing.T) {
	pathBefore, existed := os.LookupEnv("UCD_PATH")
	defer restore("UCD_PATH", pathBefore, existed)
	os.Unsetenv("UCD_PATH")                // ➊
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
		t.Errorf("AbrirUCD(%q):\n%v", ucdPath, err)
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
	done := make(chan bool)                 // ➊
	go fetchUCD(srv.URL, ucdPath, done) // ➋
	_ = <-done                              // ➌
	ucd, err := os.Open(ucdPath)
	if os.IsNotExist(err) {
		t.Errorf("fetchUCD não gerou:%v\n%v", ucdPath, err)
	}
	ucd.Close()
	os.Remove(ucdPath)
}

func TestOpenUCD_remote(t *testing.T) {
	if testing.Short() { // ➊
		t.Skip("teste ignorado [opção -test.short]") // ➋
	}
	ucdPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano()) // ➌
	ucd, err := openUCD(ucdPath)
	if err != nil {
		t.Errorf("AbrirUCD(%q):\n%v", ucdPath, err)
	}
	ucd.Close()
	os.Remove(ucdPath)
}
