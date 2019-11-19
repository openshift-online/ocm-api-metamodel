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
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Conversion to Markdown", func() {
	It("Converts files", func() {
		// Get the list of AsciiDoc files in the `tests` directory:
		adocFiles, err := filepath.Glob("tests/*.adoc")
		Expect(err).ToNot(HaveOccurred())
		Expect(len(adocFiles)).ToNot(BeZero())

		// Check that all the test files can be converted successfully:
		for _, adocFile := range adocFiles {
			mdFile := adocRE.ReplaceAllString(adocFile, ".md")
			adocData, err := ioutil.ReadFile(adocFile)
			Expect(err).ToNot(HaveOccurred())
			adocText := string(adocData)
			mdData, err := ioutil.ReadFile(mdFile)
			Expect(err).ToNot(HaveOccurred())
			mdText := string(mdData)
			actualText := Markdown(adocText)
			Expect(strings.TrimSpace(actualText)).To(Equal(strings.TrimSpace(mdText)))
		}
	})
})

// Regular expression used to replace the `.adoc` extension:
var adocRE = regexp.MustCompile(`\.adoc$`)
