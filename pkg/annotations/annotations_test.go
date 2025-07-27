/*
Copyright (c) 2025 Red Hat, Inc.

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

package annotations

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
)

var _ = Describe("Deprecation annotation support", func() {
	Describe("IsDeprecated function", func() {
		var annotatedConcept concepts.Annotated

		BeforeEach(func() {
			// Create a test type that implements concepts.Annotated
			annotatedConcept = concepts.NewType()
		})

		Context("when concept has deprecated annotation", func() {
			It("should return true for deprecated annotation without parameters", func() {
				// Create and add deprecated annotation
				deprecatedAnnotation := concepts.NewAnnotation()
				deprecatedAnnotation.SetName("deprecated")
				annotatedConcept.AddAnnotation(deprecatedAnnotation)

				// Test IsDeprecated function
				result := IsDeprecated(annotatedConcept)
				Expect(result).To(BeTrue())
			})

			It("should return true for deprecated annotation with parameters", func() {
				// Create deprecated annotation with parameters
				deprecatedAnnotation := concepts.NewAnnotation()
				deprecatedAnnotation.SetName("deprecated")

				// Add a parameter (like a deprecation message)
				deprecatedAnnotation.AddParameter("message", "This feature will be removed in v2")

				annotatedConcept.AddAnnotation(deprecatedAnnotation)

				// Test IsDeprecated function
				result := IsDeprecated(annotatedConcept)
				Expect(result).To(BeTrue())
			})

			It("should return true when deprecated annotation exists among multiple annotations", func() {
				// Add multiple annotations including deprecated
				jsonAnnotation := concepts.NewAnnotation()
				jsonAnnotation.SetName("json")
				annotatedConcept.AddAnnotation(jsonAnnotation)

				deprecatedAnnotation := concepts.NewAnnotation()
				deprecatedAnnotation.SetName("deprecated")
				annotatedConcept.AddAnnotation(deprecatedAnnotation)

				httpAnnotation := concepts.NewAnnotation()
				httpAnnotation.SetName("http")
				annotatedConcept.AddAnnotation(httpAnnotation)

				// Test IsDeprecated function
				result := IsDeprecated(annotatedConcept)
				Expect(result).To(BeTrue())
			})
		})

		Context("when concept does not have deprecated annotation", func() {
			It("should return false for concept with no annotations", func() {
				// Test IsDeprecated function on clean concept
				result := IsDeprecated(annotatedConcept)
				Expect(result).To(BeFalse())
			})

			It("should return false for concept with other annotations but no deprecated", func() {
				// Add non-deprecated annotations
				jsonAnnotation := concepts.NewAnnotation()
				jsonAnnotation.SetName("json")
				annotatedConcept.AddAnnotation(jsonAnnotation)

				httpAnnotation := concepts.NewAnnotation()
				httpAnnotation.SetName("http")
				annotatedConcept.AddAnnotation(httpAnnotation)

				goAnnotation := concepts.NewAnnotation()
				goAnnotation.SetName("go")
				annotatedConcept.AddAnnotation(goAnnotation)

				// Test IsDeprecated function
				result := IsDeprecated(annotatedConcept)
				Expect(result).To(BeFalse())
			})

			It("should return false for annotation with similar but different name", func() {
				// Add annotation with similar name
				similarAnnotation := concepts.NewAnnotation()
				similarAnnotation.SetName("deprecation") // Note: not "deprecated"
				annotatedConcept.AddAnnotation(similarAnnotation)

				// Test IsDeprecated function
				result := IsDeprecated(annotatedConcept)
				Expect(result).To(BeFalse())
			})
		})

		Context("edge cases", func() {
			It("should handle case-sensitive annotation name correctly", func() {
				// Add annotation with different case
				upperCaseAnnotation := concepts.NewAnnotation()
				upperCaseAnnotation.SetName("DEPRECATED") // Wrong case
				annotatedConcept.AddAnnotation(upperCaseAnnotation)

				// Test IsDeprecated function
				result := IsDeprecated(annotatedConcept)
				Expect(result).To(BeFalse())
			})

			It("should handle empty annotation name", func() {
				// Add annotation with empty name
				emptyAnnotation := concepts.NewAnnotation()
				emptyAnnotation.SetName("")
				annotatedConcept.AddAnnotation(emptyAnnotation)

				// Test IsDeprecated function
				result := IsDeprecated(annotatedConcept)
				Expect(result).To(BeFalse())
			})
		})
	})

	Describe("Integration with different concept types", func() {
		Context("when used with different annotated concepts", func() {
			It("should work correctly with Attribute", func() {
				attribute := concepts.NewAttribute()
				attribute.SetName(names.ParseUsingCase("testAttribute"))

				deprecatedAnnotation := concepts.NewAnnotation()
				deprecatedAnnotation.SetName("deprecated")
				attribute.AddAnnotation(deprecatedAnnotation)

				result := IsDeprecated(attribute)
				Expect(result).To(BeTrue())
			})

			It("should work correctly with Method", func() {
				method := concepts.NewMethod()
				method.SetName(names.ParseUsingCase("testMethod"))

				deprecatedAnnotation := concepts.NewAnnotation()
				deprecatedAnnotation.SetName("deprecated")
				method.AddAnnotation(deprecatedAnnotation)

				result := IsDeprecated(method)
				Expect(result).To(BeTrue())
			})

			It("should work correctly with Resource", func() {
				resource := concepts.NewResource()
				resource.SetName(names.ParseUsingCase("testResource"))

				deprecatedAnnotation := concepts.NewAnnotation()
				deprecatedAnnotation.SetName("deprecated")
				resource.AddAnnotation(deprecatedAnnotation)

				result := IsDeprecated(resource)
				Expect(result).To(BeTrue())
			})

			It("should work correctly with Parameter", func() {
				parameter := concepts.NewParameter()
				parameter.SetName(names.ParseUsingCase("testParameter"))

				deprecatedAnnotation := concepts.NewAnnotation()
				deprecatedAnnotation.SetName("deprecated")
				parameter.AddAnnotation(deprecatedAnnotation)

				result := IsDeprecated(parameter)
				Expect(result).To(BeTrue())
			})
		})
	})
})
