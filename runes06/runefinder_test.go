package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
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
	char, name, words := ParseLine(lineLetterA) // ➊
	if char != 'A' {
		t.Errorf("Esperado: 'A'; got: %q", char)
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
