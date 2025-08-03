package language

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// MakeModelWithErrors creates a model and expects errors to be reported
func MakeModelWithErrors(pairs ...string) int {
	// Create a temporary directory for the model files:
	root, err := os.MkdirTemp("", "model-*")
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
		err = os.WriteFile(path, []byte(data), 0600)
		Expect(err).ToNot(HaveOccurred())
	}

	// Create a reporter that writes to the Ginkgo output stream:
	reporter, err := reporter.New().
		Streams(GinkgoWriter, GinkgoWriter).
		Build()
	Expect(err).ToNot(HaveOccurred())

	// Read the model (may return an error when there are validation errors):
	model, err := NewReader().
		Reporter(reporter).
		Input(root).
		Read()

	// If there are validation errors, the Read() method may return an error
	// but we still want to return the count of validation errors from the reporter
	if err != nil {
		if reporter.Errors() > 0 {
			// This is expected when there are validation errors
			return reporter.Errors()
		}
		// This is an unexpected error
		Expect(err).ToNot(HaveOccurred())
	}
	Expect(model).ToNot(BeNil())

	// Return the number of errors reported:
	return reporter.Errors()
}

var _ = Describe("Validation checks", func() {

	Describe("Class Id field validation", func() {
		It("Should detect explicit Id field in class", func() {
			errors := MakeModelWithErrors(
				"my_service/v1/root.model",
				`
				resource Root {
				}
				`,
				"my_service/v1/my_class.model",
				`
				class MyClass {
					Id String
					Name String
				}
				`,
			)
			Expect(errors).To(Equal(1))
		})

		It("Should detect explicit ID field in class", func() {
			errors := MakeModelWithErrors(
				"my_service/v1/root.model",
				`
				resource Root {
				}
				`,
				"my_service/v1/my_class.model",
				`
				class MyClass {
					ID String
					Name String
				}
				`,
			)
			Expect(errors).To(Equal(1))
		})

		It("Should allow Id field in struct", func() {
			model := MakeModel(
				"my_service/v1/root.model",
				`
				resource Root {
				}
				`,
				"my_service/v1/my_struct.model",
				`
				struct MyStruct {
					Id String
					Name String
				}
				`,
			)
			Expect(model).ToNot(BeNil())
		})

		It("Should allow ID field in struct", func() {
			model := MakeModel(
				"my_service/v1/root.model",
				`
				resource Root {
				}
				`,
				"my_service/v1/my_struct.model",
				`
				struct MyStruct {
					ID String
					Name String
				}
				`,
			)
			Expect(model).ToNot(BeNil())
		})

		It("Should allow other fields in class", func() {
			model := MakeModel(
				"my_service/v1/root.model",
				`
				resource Root {
				}
				`,
				"my_service/v1/my_class.model",
				`
				class MyClass {
					Name String
					Description String
				}
				`,
			)
			Expect(model).ToNot(BeNil())
		})

		It("Should detect multiple classes with Id fields", func() {
			errors := MakeModelWithErrors(
				"my_service/v1/root.model",
				`
				resource Root {
				}
				`,
				"my_service/v1/class_one.model",
				`
				class ClassOne {
					Id String
					Name String
				}
				`,
				"my_service/v1/class_two.model",
				`
				class ClassTwo {
					ID String
					Description String
				}
				`,
			)
			Expect(errors).To(Equal(2))
		})
	})
})
