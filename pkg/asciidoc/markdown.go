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

// This file contains the code to convert from AsciiDoc to Markdown.

package asciidoc

import (
	"regexp"
	"strings"
)

// Markdown converts the given text from AsciiDoc to Markdown.
func Markdown(text string) string {
	// Split the text into lines:
	lines := strings.Split(text, "\n")

	// Transform the lines:
	lines = transformLines(lines, transformHeader)
	lines = transformLines(lines, transformBlockHeader)
	lines = transformLines(lines, transformBlockDelimiter)

	// Join the lines and add line break at the end:
	return strings.Join(lines, "\n") + "\n"
}

func transformLines(lines []string, transformLine func(string) string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = transformLine(line)
	}
	return result
}

func transformHeader(line string) string {
	if headerRE.MatchString(line) {
		data := []byte(line)
		for i := 0; data[i] == '='; i++ {
			data[i] = '#'
		}
		line = string(data)
	}
	return line
}

func transformBlockHeader(line string) string {
	if blockHeaderRE.MatchString(line) {
		line = ""
	}
	return line
}

func transformBlockDelimiter(line string) string {
	if blockDelimiterRE.MatchString(line) {
		line = "```"
	}
	return line
}

var headerRE = regexp.MustCompile(`^=+[^=]*$`)
var blockHeaderRE = regexp.MustCompile(`^\[\s*source\s*(,.*)\s*]\s*$`)
var blockDelimiterRE = regexp.MustCompile(`^(\.\.\.\.|----)\s*$`)
