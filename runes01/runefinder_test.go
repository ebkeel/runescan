package main // ➊

import "testing" // ➋

const lineLetterA = `0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;` // ➌

func TestParseLine(t *testing.T) { // ➍
	char, _ := ParseLine(lineLetterA) // ➎
	if char != 'A' {                      // ➏
		t.Errorf("Esperava 'A', veio %c", char) // ➐
	}
}
