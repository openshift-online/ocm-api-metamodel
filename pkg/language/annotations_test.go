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
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Annotation", func() {
	It("Reads attribute annotation", func() {
		// Create the model:
		model := MakeModel(
			"my_service/v1/root.model",
			`
			resource Root {
			}
			`,
			"my_service/v1/my_class.model",
			`
			class MyClass {
				@go(name = "my_name")
				attribute MyAttribute Integer
			}
			`,
		)

		// Check the annotation:
		service := model.FindService(names.ParseUsingSeparator("my_service", "_"))
		Expect(service).ToNot(BeNil())
		version := service.FindVersion(names.ParseUsingSeparator("v1", "_"))
		Expect(version).ToNot(BeNil())
		class := version.FindType(names.ParseUsingCase("MyClass"))
		Expect(class).ToNot(BeNil())
		attribute := class.FindAttribute(names.ParseUsingCase("MyAttribute"))
		Expect(attribute).ToNot(BeNil())
		annotation := attribute.GetAnnotation("go")
		Expect(annotation).ToNot(BeNil())
		name := annotation.GetString("name")
		Expect(name).To(Equal("my_name"))
	})

	It("Reads attribute annotation without parameters", func() {
		// Create the model:
		model := MakeModel(
			"my_service/v1/root.model",
			`
			resource Root {
			}
			`,
			"my_service/v1/my_class.model",
			`
			class MyClass {
				@deprecated
				attribute MyAttribute Integer
			}
			`,
		)

		// Check the annotation:
		service := model.FindService(names.ParseUsingSeparator("my_service", "_"))
		Expect(service).ToNot(BeNil())
		version := service.FindVersion(names.ParseUsingSeparator("v1", "_"))
		Expect(version).ToNot(BeNil())
		class := version.FindType(names.ParseUsingCase("MyClass"))
		Expect(class).ToNot(BeNil())
		attribute := class.FindAttribute(names.ParseUsingCase("MyAttribute"))
		Expect(attribute).ToNot(BeNil())
		annotation := attribute.GetAnnotation("deprecated")
		Expect(annotation).ToNot(BeNil())
	})
})
