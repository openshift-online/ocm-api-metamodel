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

package names

import (
	"bytes"
	"strings"
	"unicode"
)

// ParseUsingSeparator separates the given text into words, using the given separator character, and
// creates a new name containing those words. For example, to convert the text 'my_favorite_fruit'
// into a name the method can be used as follows:
//
//	name := names.ParseUsingSeparator("my_favorite_fruit", "_")
func ParseUsingSeparator(text string, separator string) *Name {
	chunks := strings.Split(text, separator)
	words := make([]*Word, len(chunks))
	for i, chunk := range chunks {
		words[i] = NewWord(chunk)
	}
	return NewName(words...)
}

// ParseUsingCase separates the given text into words, using the case transitions as separators, and
// creates a new name containing those words.
func ParseUsingCase(text string) *Name {
	// Convert the text to an array of runes so that we can easily access the previous, current
	// and next runes:
	var runes []rune
	for _, r := range text {
		runes = append(runes, r)
	}

	// Iterate the runes looking for case transitions and storing the words:
	buffer := new(bytes.Buffer)
	var words []*Word
	size := len(runes)
	for i := 0; i < size; i++ {
		var previous rune
		if i > 0 {
			previous = runes[i-1]
		}
		current := runes[i]
		var next rune
		if i < size-1 {
			next = runes[i+1]
		}
		currentUpper := unicode.IsUpper(current)
		previousLower := unicode.IsLower(previous)
		nextLower := unicode.IsLower(next)
		if currentUpper && (previousLower || nextLower) {
			if buffer.Len() > 0 {
				chunk := buffer.String()
				word := NewWord(chunk)
				words = append(words, word)
				buffer.Reset()
			}
		}
		buffer.WriteRune(current)
	}
	if buffer.Len() > 0 {
		chunk := buffer.String()
		word := NewWord(chunk)
		words = append(words, word)
	}

	// Create the name from the stored words:
	return NewName(words...)
}
