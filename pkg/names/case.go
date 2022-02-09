/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains functions used to change the case of strings.

package names

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ToUpperCamel converts the given text to camel case with the first word starting with an upper
// case character. In order to separate the original text into words it will assume that the text
// uses snake case if it contains at least one underscore character. Otherwise it will assume it
// uses camel case.
func ToUpperCamel(text string) string {
	words := splitWords(text)
	var buffer bytes.Buffer
	buffer.Grow(len(text))
	for _, word := range words {
		buffer.WriteString(toCamelWord(word))
	}
	return buffer.String()
}

// ToLowerCamel converts the given text to camel case with the first word starting with a lower case
// character. In order to separate the original text into words it will assume that the text uses
// snake case if it contains at least one underscore character. Otherwise it will assume it uses
// camel case.
func ToLowerCamel(text string) string {
	words := splitWords(text)
	var buffer bytes.Buffer
	buffer.Grow(len(text))
	for i, word := range words {
		if i == 0 {
			buffer.WriteString(strings.ToLower(word))
		} else {
			buffer.WriteString(toCamelWord(word))
		}
	}
	return buffer.String()
}

func splitWords(text string) []string {
	if strings.Contains(text, "_") {
		return splitSnakeWords(text)
	}
	return splitCamelWords(text)
}

func splitCamelWords(text string) []string {
	// Convert the text to an array of runes so that we can easily access the previous, current
	// and next runes:
	var runes []rune
	for _, r := range text {
		runes = append(runes, r)
	}

	// Split the text at the points where there are transitions:
	var chunks []string
	size := len(runes)
	if size > 1 {
		buffer := new(bytes.Buffer)
		for i := 0; i < size-1; i++ {
			current := runes[i]
			next := runes[i+1]
			buffer.WriteRune(current)
			currentLower := unicode.IsLower(current)
			nextLower := unicode.IsLower(next)
			nextDigit := unicode.IsDigit(next)
			if currentLower != nextLower && !nextDigit {
				chunk := buffer.String()
				chunks = append(chunks, chunk)
				buffer.Reset()
			}
		}
		buffer.WriteRune(runes[size-1])
		chunk := buffer.String()
		chunks = append(chunks, chunk)
	} else {
		chunks = []string{text}
	}

	// Process the chunks:
	var result []string
	size = len(chunks)
	if size > 1 {
		i := 0
		for i < size-1 {
			current := chunks[i]
			next := chunks[i+1]

			// A single upper case character followed by a lower case chunk is one
			// single word, like 'Name`:
			if len(current) == 1 && isUpper(current) && isLower(next) {
				word := strings.ToLower(current) + next
				result = append(result, word)
				i = i + 2
				continue
			}

			// A chunk with two or more upper case characters followed by a single `s`
			// is the plural form of an initialism, like `CPUs` or `IDs`:
			if len(current) >= 2 && isUpper(current) && next == "s" {
				word := current + "s"
				result = append(result, word)
				i = i + 2
				continue
			}

			// An upper case chunk followed by a lower case chunk is two words and the
			// last character of the first chunk is the first character of the second
			// word, like in `CPUList` which corresponds to `CPU` and `List`:
			if isUpper(current) && isLower(next) {
				first := current[0 : len(current)-1]
				second := strings.ToLower(current[len(current)-1:]) + next
				result = append(result, first)
				result = append(result, second)
				i = i + 2
				continue
			}

			// Anything else is a word by itself:
			result = append(result, current)
			i++
		}

		// If there is still a chunk to process then it is a word by itself:
		if i < size {
			result = append(result, chunks[i])
		}
	} else {
		result = []string{text}
	}

	return result
}

func toCamelWord(word string) string {
	if isInitialism(word) {
		return strings.ToUpper(word)
	}
	head, size := utf8.DecodeRuneInString(word)
	if unicode.IsUpper(head) {
		return word
	}
	return string(unicode.ToUpper(head)) + word[size:]
}

// ToSnake converts the given text to snake case. In order to separate the original text into words
// it will assume that the text uses snake case if it contains at least one underscore character.
// Otherwise it will assume it uses camel case.
func ToSnake(text string) string {
	words := splitWords(text)
	var buffer bytes.Buffer
	buffer.Grow(len(text) + len(words) - 1)
	for i, word := range words {
		if i > 0 {
			buffer.WriteString("_")
		}
		buffer.WriteString(toSnakeWord(word))
	}
	return buffer.String()
}

func splitSnakeWords(text string) []string {
	return strings.Split(text, "_")
}

func toSnakeWord(word string) string {
	return strings.ToLower(word)
}

// isLower checks if all the runes of the given string are lower case.
func isLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

// isUpper checks if all the runes of the given string are upper case or digits.
func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// isInitialism checks if the given string is a well know initialism that should be written using
// upper case.
func isInitialism(s string) bool {
	return initialisms[strings.ToLower(s)]
}

var initialisms = map[string]bool{
	"api": true,
	"aws": true,
	"cpu": true,
	"gcp": true,
	"id":  true,
	"ip":  true,
	"url": true,
}
