package main

import "testing"

const lineLetterA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;"

func TestParseLine(t *testing.T) {
	char, name := ParseLine(lineLetterA)
	if char != 'A' {
		t.Errorf("Esperava 'A', veio %q", char) // ➊
	}
	const nameA = "LATIN CAPITAL LETTER A" // ➋
	if name != nameA {
		t.Errorf("Esperava %q, veio %q", nameA, name) // ➌
	}
}
