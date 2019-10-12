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

package golang

import (
	"fmt"
	"path"

	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// TypesCalculatorBuilder is an object used to configure and build the Go type calculators. Don't
// create instances directly, use the NewTypesCalculator function instead.
type TypesCalculatorBuilder struct {
	reporter *reporter.Reporter
	base     string
	names    *NamesCalculator
}

// TypesCalculator is an object used to calculate Go types. Don't create instances directly, use the
// builder instead.
type TypesCalculator struct {
	reporter *reporter.Reporter
	base     string
	names    *NamesCalculator
}

// NewTypesCalculator creates a Go names calculator builder.
func NewTypesCalculator() *TypesCalculatorBuilder {
	builder := new(TypesCalculatorBuilder)
	return builder
}

// Reporter sets the object that will be used to report information about the calculation processes,
// including errors.
func (b *TypesCalculatorBuilder) Reporter(value *reporter.Reporter) *TypesCalculatorBuilder {
	b.reporter = value
	return b
}

// Base sets the import path of the base package were the code will be generated.
func (b *TypesCalculatorBuilder) Base(value string) *TypesCalculatorBuilder {
	b.base = value
	return b
}

// Names sets the names calculator object that will be used to calculate Go names.
func (b *TypesCalculatorBuilder) Names(value *NamesCalculator) *TypesCalculatorBuilder {
	b.names = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new
// calculator using it.
func (b *TypesCalculatorBuilder) Build() (calculator *TypesCalculator, err error) {
	// Check that the mandatory parameters have been provided:
	if b.reporter == nil {
		err = fmt.Errorf("reporter is mandatory")
		return
	}
	if b.base == "" {
		err = fmt.Errorf("package is mandatory")
		return
	}
	if b.names == nil {
		err = fmt.Errorf("names is mandatory")
		return
	}

	// Create the calculator:
	calculator = new(TypesCalculator)
	calculator.reporter = b.reporter
	calculator.base = b.base
	calculator.names = b.names

	return
}

// Reference creates a new type reference with the given package and text.
func (c *TypesCalculator) Reference(imprt, selector, name, text string) *TypeReference {
	reference := new(TypeReference)
	reference.imprt = imprt
	reference.selector = selector
	reference.text = text
	return reference
}

// ValueReference calculates a type reference to a value of the given type.
func (c *TypesCalculator) ValueReference(typ *concepts.Type) *TypeReference {
	version := typ.Owner()
	switch {
	case typ == version.Boolean():
		ref := new(TypeReference)
		ref.name = "bool"
		ref.text = "bool"
		return ref
	case typ == version.Integer():
		ref := new(TypeReference)
		ref.name = "int"
		ref.text = "int"
		return ref
	case typ == version.Long():
		ref := new(TypeReference)
		ref.name = "int64"
		ref.text = "int64"
		return ref
	case typ == version.Float():
		ref := new(TypeReference)
		ref.name = "float64"
		ref.text = "float64"
		return ref
	case typ == version.String():
		ref := new(TypeReference)
		ref.name = "string"
		ref.text = "string"
		return ref
	case typ == version.Date():
		ref := new(TypeReference)
		ref.imprt = "time"
		ref.selector = "time"
		ref.name = "Time"
		ref.text = "time.Time"
		return ref
	case typ.IsEnum():
		ref := new(TypeReference)
		ref.imprt, ref.selector = c.Package(typ)
		ref.name = c.names.Public(typ.Name())
		ref.text = ref.name
		return ref
	case typ.IsList():
		element := typ.Element()
		switch {
		case element.IsScalar():
			ref := c.ValueReference(element)
			ref.text = fmt.Sprintf("[]%s", ref.text)
			return ref
		case element.IsStruct():
			ref := new(TypeReference)
			ref.imprt, ref.selector = c.Package(element)
			ref.name = c.names.Public(names.Cat(element.Name(), nomenclator.List))
			ref.text = fmt.Sprintf("%s.%s", ref.selector, ref.name)
			return ref
		default:
			c.reporter.Errorf(
				"Don't know how to make value reference for list of '%s'",
				element.Name(),
			)
			return new(TypeReference)
		}
	case typ.IsMap():
		element := typ.Element()
		switch {
		case element.IsScalar():
			ref := c.ValueReference(element)
			ref.text = fmt.Sprintf("map[string]%s", ref.text)
			return ref
		case element.IsStruct():
			ref := new(TypeReference)
			ref.imprt, ref.selector = c.Package(element)
			ref.name = c.names.Public(names.Cat(element.Name(), nomenclator.Map))
			ref.text = fmt.Sprintf("%s.%s", ref.selector, ref.name)
			return ref
		default:
			c.reporter.Errorf(
				"Don't know how to make value reference for map of '%s'",
				element.Name(),
			)
			return new(TypeReference)
		}
	case typ.IsStruct():
		ref := new(TypeReference)
		ref.imprt, ref.selector = c.Package(typ)
		ref.name = c.names.Public(typ.Name())
		ref.text = ref.name
		return ref
	default:
		c.reporter.Errorf(
			"Don't know how to make value reference for type '%s' of kind '%s'",
			typ.Name(), typ.Kind(),
		)
		return new(TypeReference)
	}
}

// NullableReference calculates a type reference for a value of the given type that can be assigned
// the nil value.
func (c *TypesCalculator) NullableReference(typ *concepts.Type) *TypeReference {
	switch {
	case typ.IsScalar() || typ.IsStruct():
		ref := c.ValueReference(typ)
		ref.text = fmt.Sprintf("*%s", ref.text)
		return ref
	case typ.IsList():
		element := typ.Element()
		switch {
		case element.IsScalar():
			ref := c.ValueReference(element)
			ref.text = fmt.Sprintf("[]%s", ref.text)
			return ref
		case element.IsStruct():
			ref := new(TypeReference)
			ref.imprt, ref.selector = c.Package(element)
			ref.name = c.names.Public(names.Cat(element.Name(), nomenclator.List))
			ref.text = fmt.Sprintf("*%s.%s", ref.selector, ref.name)
			return ref
		default:
			c.reporter.Errorf(
				"Don't know how to make type reference for list of '%s'",
				element.Name(),
			)
			return new(TypeReference)
		}
	case typ.IsMap():
		element := typ.Element()
		switch {
		case element.IsScalar():
			ref := c.ValueReference(element)
			ref.text = fmt.Sprintf("map[string]%s", ref.text)
			return ref
		case element.IsStruct():
			ref := c.ValueReference(element)
			ref.text = fmt.Sprintf("map[string]*%s", ref.text)
			return ref
		default:
			c.reporter.Errorf(
				"Don't know how to make type reference for map of '%s'",
				element.Name(),
			)
			return new(TypeReference)
		}
	default:
		c.reporter.Errorf(
			"Don't know how to make type reference for type '%s' of kind '%s'",
			typ.Name(), typ.Kind(),
		)
		return new(TypeReference)
	}
}

// DataReference calculates a reference for the type used to marshal and unmarshal JSON documents
// containing the given type.
func (c *TypesCalculator) DataReference(typ *concepts.Type) *TypeReference {
	switch {
	case typ.IsScalar():
		return c.NullableReference(typ)
	case typ.IsList():
		element := typ.Element()
		switch {
		case element.IsScalar():
			return c.NullableReference(typ)
		case element.IsStruct():
			ref := new(TypeReference)
			ref.imprt, ref.selector = c.Package(element)
			ref.name = c.names.Private(names.Cat(
				element.Name(), nomenclator.List, nomenclator.Data),
			)
			ref.text = fmt.Sprintf("%s.%s", ref.selector, ref.name)
			return ref
		default:
			c.reporter.Errorf(
				"Don't know how to make data reference for list of '%s'",
				element.Name(),
			)
			return new(TypeReference)
		}
	case typ.IsMap():
		element := typ.Element()
		switch {
		case element.IsScalar():
			return c.NullableReference(typ)
		case element.IsStruct():
			ref := c.DataReference(element)
			ref.text = fmt.Sprintf("map[string]%s", ref.text)
			return ref
		default:
			c.reporter.Errorf(
				"Don't know how to make data reference for map of '%s'",
				element.Name(),
			)
			return new(TypeReference)
		}
	case typ.IsStruct():
		ref := new(TypeReference)
		ref.imprt, ref.selector = c.Package(typ)
		ref.name = c.names.Private(names.Cat(typ.Name(), nomenclator.Data))
		ref.text = fmt.Sprintf("*%s.%s", ref.selector, ref.name)
		return ref
	default:
		c.reporter.Errorf(
			"Don't know how to make data reference for type '%s'",
			typ.Name(),
		)
		return new(TypeReference)
	}
}

// LinkDataReferece reference calculates a reference for the type used to marshal and unmarshal
// links to objects of the given type.
func (c *TypesCalculator) LinkDataReference(typ *concepts.Type) *TypeReference {
	switch {
	case typ.IsList() && typ.Element().IsClass():
		element := typ.Element()
		ref := new(TypeReference)
		ref.imprt, ref.selector = c.Package(element)
		ref.name = c.names.Private(names.Cat(
			element.Name(),
			nomenclator.List,
			nomenclator.Link,
			nomenclator.Data,
		))
		ref.text = fmt.Sprintf("*%s.%s", ref.selector, ref.name)
		return ref
	case typ.IsStruct():
		return c.DataReference(typ)
	default:
		c.reporter.Errorf(
			"Don't know how to make link data reference for type '%s'",
			typ.Name(),
		)
		return new(TypeReference)
	}
}

// Builder reference calculates a reference for the type used build objects of the given type.
func (c *TypesCalculator) BuilderReference(typ *concepts.Type) *TypeReference {
	switch {
	case typ.IsList():
		element := typ.Element()
		switch {
		case element.IsStruct():
			ref := new(TypeReference)
			ref.imprt, ref.selector = c.Package(element)
			ref.name = c.names.Public(names.Cat(
				element.Name(), nomenclator.List, nomenclator.Builder,
			))
			ref.text = fmt.Sprintf("*%s.%s", ref.selector, ref.name)
			return ref
		default:
			c.reporter.Errorf(
				"Don't know how to make builder reference for list of '%s'",
				element.Name(),
			)
			return new(TypeReference)
		}
	case typ.IsMap():
		element := typ.Element()
		switch {
		case element.IsStruct():
			ref := new(TypeReference)
			ref.imprt, ref.selector = c.Package(element)
			ref.name = c.names.Public(names.Cat(
				element.Name(), nomenclator.Map, nomenclator.Builder,
			))
			ref.text = fmt.Sprintf("*%s.%s", ref.selector, ref.name)
			return ref
		default:
			c.reporter.Errorf(
				"Don't know how to make builder reference for map of '%s'",
				element.Name(),
			)
			return new(TypeReference)
		}
	case typ.IsStruct():
		ref := new(TypeReference)
		ref.imprt, ref.selector = c.Package(typ)
		ref.name = c.names.Public(names.Cat(typ.Name(), nomenclator.Builder))
		ref.text = fmt.Sprintf("*%s.%s", ref.selector, ref.name)
		return ref
	default:
		c.reporter.Errorf(
			"Don't know how to make builder reference for type '%s' of kind '%s'",
			typ.Name(), typ.Kind(),
		)
		return new(TypeReference)
	}
}

// Zero value calculates the zero value for the given type.
func (c *TypesCalculator) ZeroValue(typ *concepts.Type) string {
	version := typ.Owner()
	switch {
	case typ.IsStruct() || typ.IsList() || typ.IsMap():
		return `nil`
	case typ.IsEnum():
		ref := c.ValueReference(typ)
		return fmt.Sprintf(`%s("")`, ref.Name())
	case typ == version.Boolean():
		return `false`
	case typ == version.Integer():
		return `0`
	case typ == version.Long():
		return `0`
	case typ == version.Float():
		return `0.0`
	case typ == version.Date():
		return `time.Time{}`
	case typ == version.String():
		return `""`
	default:
		c.reporter.Errorf(
			"Don't know how to calculate zero value for type '%s'",
			typ.Name(),
		)
		return ""
	}
}

// Package calculates the name of the package for the given type. It returns the full import path of
// the package and the package selector.
func (c *TypesCalculator) Package(typ *concepts.Type) (imprt, selector string) {
	version := typ.Owner()
	service := version.Owner()
	switch {
	case typ == version.Boolean():
		return
	case typ == version.Integer():
		return
	case typ == version.Float():
		return
	case typ == version.String():
		return
	case typ == version.Date():
		imprt = "time"
		selector = "time"
		return
	case typ.IsEnum() || typ.IsStruct() || typ.IsList():
		servicePkg := c.names.Package(service.Name())
		versionPkg := c.names.Package(version.Name())
		imprt = path.Join(c.base, servicePkg, versionPkg)
		selector = path.Base(imprt)
		return
	default:
		c.reporter.Errorf(
			"Don't know how to make package name for type '%s' of kind '%s'",
			typ.Name(), typ.Kind(),
		)
		return
	}
}

// TypeReference represents a reference to a Go type.
type TypeReference struct {
	imprt    string
	selector string
	name     string
	text     string
}

// Import returns the import path of the package that contains the type.
func (r *TypeReference) Import() string {
	return r.imprt
}

// Selector returns the selector that should be used to import the package.
func (r *TypeReference) Selector() string {
	return r.selector
}

// Name returns the local name of the base type, without the selector.
func (r *TypeReference) Name() string {
	return r.name
}

// Text returns the text of the reference.
func (r *TypeReference) Text() string {
	return r.text
}
