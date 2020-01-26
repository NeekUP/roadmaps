package core

import (
	"fmt"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

const Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const AlphabetLen = 62

func EncodeNumToString(i int) string {
	if i == 0 {
		return string(Alphabet[0])
	}

	s := make([]byte, 0, 0)

	for i > 0 {
		s = append(s, Alphabet[i%AlphabetLen])
		i = i / AlphabetLen
	}

	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}

	return string(s)
}

func DecodeStringToNum(s string) (int, error) {
	num := 0

	for _, c := range s {
		position := -1
		for i := 0; i < AlphabetLen; i++ {
			if Alphabet[i] == byte(c) {
				position = i
				break
			}
		}

		if position == -1 {
			return 0, fmt.Errorf("Invalid input string.")
		}

		num = (num * AlphabetLen) + position
	}

	return num, nil
}

// func RenderMarkdown(text string) []byte {
// 	extentions := blackfriday.NoIntraEmphasis
// 	extentions |= blackfriday.Tables
// 	extentions |= blackfriday.FencedCode
// 	extentions |= blackfriday.DefinitionLists

// 	options := blackfriday.WithExtensions(extentions)
// 	return blackfriday.Run([]byte(text), options)
// }

func SanitizeString(text string) string {
	return string(SanitizeInput([]byte(text)))
}

func SanitizeInput(data []byte) []byte {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return p.SanitizeBytes(data)
}
