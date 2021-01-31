/*
Copyright (c) 2022 Red Hat, Inc.

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

package language

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

func TestLanguage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Language")
}

// MakeModel creates a temporary directory, and inside it all files whose names and contents are
// given in the list of pairs. It then reads those model files and returns the resulting model.
func MakeModel(pairs ...string) *concepts.Model {
	// Create a temporary directory for the model files:
	root, err := ioutil.TempDir("", "model-*")
	Expect(err).ToNot(HaveOccurred())
	defer func() {
		err = os.RemoveAll(root)
		Expect(err).ToNot(HaveOccurred())
	}()

	// Write the model files into the temporary directory:
	count := len(pairs) / 2
	for i := 0; i < count; i++ {
		name := pairs[2*i]
		data := pairs[2*i+1]
		path := filepath.Join(root, name)
		dir := filepath.Dir(path)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
			Expect(err).ToNot(HaveOccurred())
		}
		err = ioutil.WriteFile(path, []byte(data), 0600)
		Expect(err).ToNot(HaveOccurred())
	}

	// Create a reporter that writes to the Ginkgo output stream:
	reporter, err := reporter.New().
		Streams(GinkgoWriter, GinkgoWriter).
		Build()
	Expect(err).ToNot(HaveOccurred())

	// Read the model:
	model, err := NewReader().
		Reporter(reporter).
		Input(root).
		Read()
	Expect(err).ToNot(HaveOccurred())
	Expect(model).ToNot(BeNil())

	// Check that there are no errors reported:
	Expect(reporter.Errors()).To(BeZero())

	// read the model:
	return model
}
