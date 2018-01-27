package main

import (
	"os"
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
	char, name := ParseLine(lineLetterA)
	if char != 'A' {
		t.Errorf("Esperava 'A', veio %q", char)
	}
	const nameA = "LATIN CAPITAL LETTER A"
	if name != nameA {
		t.Errorf("Esperava %q, veio %q", nameA, name)
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

func Example() { // ➊
	oldArgs := os.Args                   // ➋
	defer func() { os.Args = oldArgs }() // ➌
	os.Args = []string{"", "cruzeiro"}   // ➍
	main()                               // ➎
	// Output:
	// U+20A2	₢	CRUZEIRO SIGN
}
