package txt

import (
	"strings"
)

// https://github.com/go-dedup/metaphone/blob/master/metaphone.go

func Metaphone(word string) string {
	word     = strings.ToUpper(word)
	word     = RmDup(word)
	wordLen := len(word)

	if wordLen > 1 {
		switch word[0:2] {
		case "PN", "AE", "KN", "GN", "WR":
			word = word[1:]
		case "WH":
			word = "W" + word[2:]
		}
		if word[0:1] == "X" {
			word = "W" + word[1:]
		}
	}

	result := ""
	for i, rune := range word {
		switch rune {
		case 'B':
			{
				if i != wordLen-1 || SubStr(word, i-1, 2) != "MB" {
					result = result + "B"
				}
			}
		case 'C':
			{
				if SubStr(word, i, 3) == "CIA" || SubStr(word, i, 2) == "CH" {
					result = result + "X"
				} else if SubStr(word, i, 2) == "CI" || SubStr(word, i, 2) == "CE" || SubStr(word, i, 2) == "CY" {
					result = result + "S"
				} else if SubStr(word, i-1, 3) != "SCI" || SubStr(word, i-1, 3) != "SCE" || SubStr(word, i-1, 3) != "SCY" {
					result = result + "K"
				}
			}
		case 'D':
			{
				if SubStr(word, i, 3) == "DGE" || SubStr(word, i, 3) == "DGY" || SubStr(word, i, 3) == "DGI" {
					result = result + "J"
				} else {
					result = result + "T"
				}
			}
		case 'F':
			result = result + "F"
		case 'G':
			{
				prev := SubStr(word, i+1, 1)
				if (SubStr(word, i, 2) == "GH" && !isVowel(SubStr(word, i+2, 1))) ||
					SubStr(word, i, 2) == "GN" ||
					SubStr(word, i, 4) == "GNED" ||
					SubStr(word, i, 3) == "GDE" ||
					SubStr(word, i, 3) == "GDY" ||
					SubStr(word, i, 3) == "GDI" {
				} else if prev == "I" || prev == "E" || prev == "Y" {
					result = result + "J"
				} else {
					result = result + "K"
				}
			}
		case 'H':
			{
				if !isVowel(SubStr(word, i+1, 1)) &&
					SubStr(word, i-2, 2) != "CH" &&
					SubStr(word, i-2, 2) != "SH" &&
					SubStr(word, i-2, 2) != "PH" &&
					SubStr(word, i-2, 2) != "TH" &&
					SubStr(word, i-2, 2) != "GH" {
					result = result + "H"
				}
			}
		case 'J':
			result = result + "J"
		case 'K':
			{
				if SubStr(word, i-1, 1) != "C" {
					result = result + "K"
				}
			}
		case 'L':
			result = result + "L"
		case 'M':
			result = result + "M"
		case 'N':
			result = result + "N"
		case 'P':
			{
				if SubStr(word, i+1, 1) == "H" {
					result = result + "F"
				} else {
					result = result + "P"
				}
			}
		case 'Q':
			result = result + "K"
		case 'R':
			result = result + "R"
		case 'S':
			{
				if SubStr(word, i+1, 1) == "H" || SubStr(word, i, 3) == "SIO" || SubStr(word, i, 3) == "SIA" {
					result = result + "X"
				} else {
					result = result + "S"
				}
			}
		case 'T':
			{
				if SubStr(word, i, 3) == "TIO" || SubStr(word, i, 3) == "TIA" {
					result = result + "X"
				} else if SubStr(word, i+1, 1) == "H" {
					result = result + "0"
				} else if SubStr(word, i, 3) != "TCH" {
					result = result + "T"
				}
			}
		case 'V':
			result = result + "F"
		case 'W':
			{
				if isVowel(SubStr(word, i+1, 1)) {
					result = result + "W"
				}
			}
		case 'X':
			result = result + "KS"
		case 'Y':
			{
				if isVowel(SubStr(word, i+1, 1)) {
					result = result + "Y"
				}
			}
		case 'Z':
			result = result + "S"
		}
	}
	return result
}

func isVowel(char string) bool {
	return strings.Index("AEIOU", char) > -1
}