/*
Copyright (c) 2024 Red Hat, Inc.

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

var _ = Describe("Read Model with ref annotation", func() {

	It("Reads referenced class scalar attribute", func() {
		model := MakeModel(
			"my_service/v1_alpha1/root.model",
			`
			resource Root {
			}
			`,
			"my_service/v1_alpha1/my_class.model",
			`
			@ref(path="other_service/v1/my_class")
			class MyClass {
			}
			`,
			"other_service/v1/root.model",
			`
			resource Root{
			}
			`,
			"other_service/v1/my_class.model",
			`
			class MyClass {
				MyAttribute Integer
			}
			`,
		)
		// Check the attribute and its owner
		service := model.FindService(names.ParseUsingSeparator("my_service", "_"))
		Expect(service).ToNot(BeNil())
		version := service.FindVersion(names.ParseUsingSeparator("v1_alpha1", "_"))
		Expect(version).ToNot(BeNil())
		class := version.FindType(names.ParseUsingCase("MyClass"))
		Expect(class).ToNot(BeNil())
		attribute := class.FindAttribute(names.ParseUsingCase("MyAttribute"))
		Expect(attribute).ToNot(BeNil())
		Expect(attribute.Type().Owner().Name().String()).To(Equal("v1"))
	})

	It("References respect link attribute", func() {
		model := MakeModel(
			"my_service/v1_alpha1/root.model",
			`
			resource Root {
			}
			`,
			"my_service/v1_alpha1/my_class.model",
			`
			@ref(path="other_service/v1/my_class")
			class MyClass {
			}
			`,
			"other_service/v1/root.model",
			`
			resource Root{
			}
			`,
			"other_service/v1/my_class.model",
			`
			class MyClass {
				link MyAttribute MyAttribute
			}
			`,
			"other_service/v1/my_attribute.model",
			`
			class MyAttribute{
			}
			`,
		)
		// Check the attribute and its owner
		service := model.FindService(names.ParseUsingSeparator("my_service", "_"))
		Expect(service).ToNot(BeNil())
		version := service.FindVersion(names.ParseUsingSeparator("v1_alpha1", "_"))
		Expect(version).ToNot(BeNil())
		class := version.FindType(names.ParseUsingCase("MyClass"))
		Expect(class).ToNot(BeNil())
		myAttribute := class.FindAttribute(names.ParseUsingCase("MyAttribute"))
		Expect(myAttribute).ToNot(BeNil())
		Expect(myAttribute.Type().Owner().Name().String()).To(Equal("v1"))
		Expect(myAttribute.Link()).To(BeTrue())
	})

	It("Reads referenced class list attribute", func() {
		model := MakeModel(
			"my_service/v1_alpha1/root.model",
			`
			resource Root {
			}
			`,
			"my_service/v1_alpha1/my_class.model",
			`
			@ref(path="other_service/v1/my_class")
			class MyClass {
			}
			`,
			"other_service/v1/root.model",
			`
			resource Root{
			}
			`,
			"other_service/v1/my_class.model",
			`
			class MyClass {
				Foo []MyAttribute
			}`,
			"other_service/v1/my_attribute.model",
			`
			class MyAttribute{
			}
			`,
		)
		// Check the attribute and its owner
		service := model.FindService(names.ParseUsingSeparator("my_service", "_"))
		Expect(service).ToNot(BeNil())
		version := service.FindVersion(names.ParseUsingSeparator("v1_alpha1", "_"))
		Expect(version).ToNot(BeNil())
		class := version.FindType(names.ParseUsingCase("MyClass"))
		Expect(class).ToNot(BeNil())
		attributeType := version.FindType(names.ParseUsingCase("MyAttribute"))
		Expect(attributeType).ToNot(BeNil())
		Expect(attributeType.Owner().Name().String()).To(Equal("v1_alpha1"))
		attributeList := class.FindAttribute(names.ParseUsingCase("Foo"))
		Expect(attributeList).ToNot(BeNil())
		Expect(attributeList.Type().IsList()).To(BeTrue())
		Expect(attributeList.Type().Owner().Name().String()).To(Equal("v1_alpha1"))
	})

	It("Overrides class with other class definition", func() {
		model := MakeModel(
			"my_service/v1_alpha1/root.model",
			`
			resource Root {
			}
			`,
			"my_service/v1_alpha1/my_class.model",
			`
			@ref(path="other_service/v1/my_class")
			class MyClass {
			}
			`,
			"my_service/v1_alpha1/my_attribute.model",
			`
			@ref(path="other_service/v1/my_attribute")
			class MyAttribute {
			}
			`,
			"other_service/v1/root.model",
			`
			resource Root{
			}
			`,
			"other_service/v1/my_class.model",
			`
			class MyClass {
				link Foo []MyAttribute
			}`,
			"other_service/v1/my_attribute.model",
			`
			class MyAttribute{
				Goo Bar
			}
			`,
			"other_service/v1/bar.model",
			`
			class Bar {
			}
			`,
		)
		// Check the attribute and its owner
		service := model.FindService(names.ParseUsingSeparator("my_service", "_"))
		Expect(service).ToNot(BeNil())
		version := service.FindVersion(names.ParseUsingSeparator("v1_alpha1", "_"))
		Expect(version).ToNot(BeNil())
		class := version.FindType(names.ParseUsingCase("MyClass"))
		Expect(class).ToNot(BeNil())
		Expect(class.Owner().Name().String()).To(Equal("v1_alpha1"))
		attributeList := class.FindAttribute(names.ParseUsingCase("Foo"))
		Expect(attributeList).ToNot(BeNil())
		Expect(attributeList.Type().IsList()).To(BeTrue())
		Expect(attributeList.Type().Owner().Name().String()).To(Equal("v1_alpha1"))
		Expect(attributeList.Type().Element().Owner().Name().String()).To(Equal("v1_alpha1"))
		barType := version.FindType(names.ParseUsingCase("Bar"))
		Expect(barType.Owner().Name().String()).To(Equal("v1_alpha1"))
	})
})
