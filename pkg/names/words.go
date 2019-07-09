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
	"strings"
)

// Words represents a single word.
type Word struct {
	text    string
	acronym bool
}

// NewWord creates a new word from the given text. If the text is all upper case then the word will
// be considered an acronym.
func NewWord(text string) *Word {
	word := new(Word)
	word.acronym = len(text) >= 2 && strings.ToUpper(text) == text
	if word.acronym {
		word.text = text
	} else {
		word.text = strings.ToLower(text)
	}
	return word
}

// String returns a string representation of this word.
func (w *Word) String() string {
	return w.text
}

// Capitalize converts this word to a capitalized string: first character using upper case and the
// rest using lower case for normal words, and all upper case for acronyms.
func (w *Word) Capitalize() string {
	if w.acronym {
		return w.text
	}
	return strings.Title(w.text)
}

// Equals check if this word is equal to the given word.
func (w *Word) Equals(word *Word) bool {
	if w == nil && word == nil {
		return true
	}
	if w == nil || word == nil {
		return false
	}
	return w.text == word.text && w.acronym == word.acronym
}

// Words is an slice of names, intended to simplify sorting.
type Words []*Word

// Len implements the Len method of sort.Interface, so that slices of words can be easily sorted.
func (w Words) Len() int {
	return len(w)
}

// Swap implements the Swap method of sort.Interface, so that slices of words can be easily sorted.
func (w Words) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// Less implements the Less method of sort.Interface, so that alices of words can be easily sorted.
func (w Words) Less(i, j int) bool {
	return w[i].text < w[j].text
}
