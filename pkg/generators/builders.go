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

package generators

import (
	"fmt"
	"path/filepath"

	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/golang"
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// BuildersGeneratorBuilder is an object used to configure and build the builders generator. Don't
// create instances directly, use the NewBuildersGenerator function instead.
type BuildersGeneratorBuilder struct {
	reporter *reporter.Reporter
	model    *concepts.Model
	output   string
	base     string
	names    *golang.NamesCalculator
	types    *golang.TypesCalculator
}

// BuildersGenerator generates code for the builders of the model types. Don't create instances
// directly, use the builder instead.
type BuildersGenerator struct {
	reporter *reporter.Reporter
	errors   int
	model    *concepts.Model
	output   string
	base     string
	names    *golang.NamesCalculator
	types    *golang.TypesCalculator
	buffer   *golang.Buffer
}

// NewBuildersGenerator creates a new builder for builders generators.
func NewBuildersGenerator() *BuildersGeneratorBuilder {
	return new(BuildersGeneratorBuilder)
}

// Reporter sets the object that will be used to report information about the generation process,
// including errors.
func (b *BuildersGeneratorBuilder) Reporter(value *reporter.Reporter) *BuildersGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the types generator.
func (b *BuildersGeneratorBuilder) Model(value *concepts.Model) *BuildersGeneratorBuilder {
	b.model = value
	return b
}

// Output sets the directory where the source will be generated.
func (b *BuildersGeneratorBuilder) Output(value string) *BuildersGeneratorBuilder {
	b.output = value
	return b
}

// Base sets the import import path of the base output package.
func (b *BuildersGeneratorBuilder) Base(value string) *BuildersGeneratorBuilder {
	b.base = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *BuildersGeneratorBuilder) Names(value *golang.NamesCalculator) *BuildersGeneratorBuilder {
	b.names = value
	return b
}

// Types sets the object that will be used to calculate types.
func (b *BuildersGeneratorBuilder) Types(value *golang.TypesCalculator) *BuildersGeneratorBuilder {
	b.types = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new
// builders generator using it.
func (b *BuildersGeneratorBuilder) Build() (generator *BuildersGenerator, err error) {
	// Check that the mandatory parameters have been provided:
	if b.reporter == nil {
		err = fmt.Errorf("reporter is mandatory")
		return
	}
	if b.model == nil {
		err = fmt.Errorf("model is mandatory")
		return
	}
	if b.output == "" {
		err = fmt.Errorf("output is mandatory")
		return
	}
	if b.base == "" {
		err = fmt.Errorf("base is mandatory")
		return
	}
	if b.names == nil {
		err = fmt.Errorf("names is mandatory")
		return
	}
	if b.types == nil {
		err = fmt.Errorf("types is mandatory")
		return
	}

	// Create the generator:
	generator = new(BuildersGenerator)
	generator.reporter = b.reporter
	generator.model = b.model
	generator.output = b.output
	generator.base = b.base
	generator.names = b.names
	generator.types = b.types

	return
}

// Run executes the code generator.
func (g *BuildersGenerator) Run() error {
	var err error

	// Generate the code for each type:
	for _, service := range g.model.Services() {
		for _, version := range service.Versions() {
			for _, typ := range version.Types() {
				switch {
				case typ.IsStruct():
					err = g.generateStructBuilderFile(typ)
				case typ.IsList() && typ.Element().IsStruct():
					err = g.generateListBuilderFile(typ)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	// Check if there were errors:
	if g.errors > 0 {
		if g.errors > 1 {
			err = fmt.Errorf("there were %d errors", g.errors)
		} else {
			err = fmt.Errorf("there was 1 error")
		}
		return err
	}

	return nil
}

func (g *BuildersGenerator) generateStructBuilderFile(typ *concepts.Type) error {
	var err error

	// Calculate the package and file names:
	pkgName := g.pkgName(typ.Owner())
	fileName := g.fileName(typ)

	// Create the buffer for the generated code:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Base(g.base).
		Package(pkgName).
		File(fileName).
		Function("builderCtor", g.builderCtor).
		Function("builderName", g.builderName).
		Function("fieldName", g.fieldName).
		Function("fieldType", g.fieldType).
		Function("objectName", g.objectName).
		Function("setterName", g.setterName).
		Function("setterType", g.setterType).
		Function("valueType", g.valueType).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateStructBuilderSource(typ)

	// Write the generated code:
	return g.buffer.Write()

}

func (g *BuildersGenerator) generateStructBuilderSource(typ *concepts.Type) {
	g.buffer.Emit(`
		{{ $builderName := builderName .Type }}
		{{ $builderCtor := builderCtor .Type }}
		{{ $objectName := objectName .Type }}

		// {{ $builderName }} contains the data and logic needed to build '{{ .Type.Name }}' objects.
		//
		{{ lineComment .Type.Doc }}
		type  {{ $builderName }} struct {
			{{ if .Type.IsClass }}
				id   *string
				href *string
				link  bool
			{{ end }}
			{{ range .Type.Attributes }}
				{{ fieldName . }} {{ fieldType . }}
			{{ end }}
		}

		// {{ $builderCtor }} creates a new builder of '{{ .Type.Name }}' objects.
		func {{ $builderCtor }}() *{{ $builderName }} {
			return new({{ $builderName }})
		}

		{{ if .Type.IsClass }}
			// ID sets the identifier of the object.
			func (b *{{ $builderName }}) ID(value string) *{{ $builderName }} {
				b.id = &value
				return b
			}

			// HREF sets the link to the object.
			func (b *{{ $builderName }}) HREF(value string) *{{ $builderName }} {
				b.href = &value
				return b
			}

			// Link sets the flag that indicates if this is a link.
			func (b *{{ $builderName }}) Link(value bool) *{{ $builderName }} {
				b.link = value
				return b
			}
		{{ end }}

		{{ range .Type.Attributes }}
			{{ $fieldName := fieldName . }}
			{{ $setterName := setterName . }}
			{{ $setterType := setterType . }}

			{{ if .Type.IsList }}
				{{ $elementType := valueType .Type.Element }}

				{{ if .Type.Element.IsScalar }}
					// {{ $setterName }} sets the value of the '{{ .Name }}' attribute
					// to the given values.
					//
					{{ lineComment .Type.Doc }}
					func (b *{{ $builderName }}) {{ $setterName }}(values ...{{ $elementType }}) *{{ $builderName }} {
						b.{{ $fieldName }} = make([]{{ $elementType }}, len(values))
						copy(b.{{ $fieldName }}, values)
						return b
					}
				{{ else }}
					{{ $elementBuilderName := builderName .Type.Element }}

					// {{ $setterName }} sets the value of the '{{ .Name }}' attribute
					// to the given values.
					//
					{{ lineComment .Type.Doc }}
					func (b *{{ $builderName }}) {{ $setterName }}(values ...*{{ $elementBuilderName }}) *{{ $builderName }} {
						b.{{ $fieldName }} = make([]*{{ $elementBuilderName }}, len(values))
						copy(b.{{ $fieldName }}, values)
						return b
					}
				{{ end }}
			{{ else }}
				// {{ $setterName }} sets the value of the '{{ .Name }}' attribute
				// to the given value.
				//
				{{ lineComment .Type.Doc }}
				func (b *{{ $builderName }}) {{ $setterName }}(value {{ $setterType }}) *{{ $builderName }} {
					{{ if .Type.IsScalar }}
						b.{{ $fieldName }} = &value
					{{ else }}
						b.{{ $fieldName }} = value
					{{ end }}
					return b
				}
			{{ end }}
		{{ end }}

		// Copy copies the attributes of the given object into this builder, discarding any previous values.
		func (b *{{ $builderName }}) Copy(object *{{ $objectName }}) *{{ $builderName }} {
			if object == nil {
				return b
			}
			{{ if .Type.IsClass }}
				b.id = object.id
				b.href = object.href
				b.link = object.link
			{{ end }}
			{{ range .Type.Attributes }}
				{{ $fieldName := fieldName . }}
				{{ $fieldType := fieldType . }}
				{{ if .Type.IsStruct }}
					if object.{{ $fieldName }} != nil {
						b.{{ $fieldName }} = {{ builderCtor .Type }}().Copy(object.{{ $fieldName }})
					} else {
						b.{{ $fieldName }} = nil
					}
				{{ else if .Type.IsList }}
					{{ if .Type.Element.IsScalar }}
						if len(object.{{ $fieldName }}) > 0 {
							b.{{ $fieldName }} = make({{ $fieldType }}, len(object.{{ $fieldName }}))
							copy(b.{{ $fieldName }}, object.{{ $fieldName }})
						} else {
							b.{{ $fieldName }} = nil
						}
					{{ else if .Type.Element.IsStruct }}
						if object.{{ $fieldName }} != nil && len(object.{{ $fieldName }}.items) > 0 {
							b.{{ $fieldName }} = make([]*{{ builderName .Type.Element }}, len(object.{{ $fieldName }}.items))
							for i, item := range object.{{ $fieldName }}.items {
								b.{{ $fieldName }}[i] = {{ builderCtor .Type.Element }}().Copy(item)
							}
						} else {
							b.{{ $fieldName }} = nil
						}
					{{ end }}
				{{ else if .Type.IsMap }}
					{{ if .Type.Element.IsScalar }}
						if len(object.{{ $fieldName }}) > 0 {
							b.{{ $fieldName }} = make({{ $fieldType }})
							for key, value := range object.{{ $fieldName }} {
								b.{{ $fieldName }}[key] = value
							}
						} else {
							b.{{ $fieldName }} = nil
						}
					{{ else if .Type.Element.IsStruct }}
						if len(object.{{ $fieldName }}) > 0 {
							b.{{ $fieldName }} = make(map[string]*{{ builderName .Type.Element }})
							for key, value := range object.{{ $fieldName }} {
								b.{{ $fieldName }}[key] = {{ builderCtor .Type.Element }}().Copy(value)
							}
						} else {
							b.{{ $fieldName }} = nil
						}
					{{ end }}
				{{ else }}
					b.{{ $fieldName }} = object.{{ $fieldName }}
				{{ end }}
			{{ end }}
			return b
		}

		// Build creates a '{{ .Type.Name }}' object using the configuration stored in the builder.
		func (b *{{ $builderName }}) Build() (object *{{ $objectName }}, err error) {
			object = new({{ $objectName }})
			{{ if .Type.IsClass }}
				object.id = b.id
				object.href = b.href
				object.link = b.link
			{{ end }}
			{{ range .Type.Attributes }}
				{{ $fieldName := fieldName . }}
				{{ $fieldType := fieldType . }}
				if b.{{ $fieldName }} != nil {
					{{ if .Type.IsStruct }}
						object.{{ $fieldName }}, err = b.{{ $fieldName }}.Build()
						if err != nil {
							return
						}
					{{ else if .Type.IsList}}
						{{ if .Type.Element.IsScalar }}
							object.{{ $fieldName }} = make({{ $fieldType }}, len(b.{{ $fieldName }}))
							copy(object.{{ $fieldName }}, b.{{ $fieldName }})
						{{ else if .Type.Element.IsStruct }}
							object.{{ $fieldName }} = new({{ objectName .Type }})
							object.{{ $fieldName }}.items = make([]*{{ objectName .Type.Element }}, len(b.{{ $fieldName }}))
							for i, item := range b.{{ $fieldName }} {
								object.{{ $fieldName }}.items[i], err = item.Build()
								if err != nil {
									return
								}
							}
						{{ end }}
					{{ else if .Type.IsMap}}
						{{ if .Type.Element.IsScalar }}
							object.{{ $fieldName }} = make({{ $fieldType }})
							for key, value := range b.{{ $fieldName }} {
								object.{{ $fieldName }}[key] = value
							}
						{{ else if .Type.Element.IsStruct }}
							object.{{ $fieldName }}  = make(map[string]*{{ objectName .Type.Element }})
							for key, value := range b.{{ $fieldName }} {
								object.{{ $fieldName }}[key], err = value.Build()
								if err != nil {
									return
								}
							}
						{{ end }}
					{{ else }}
						object.{{ $fieldName }} = b.{{ $fieldName }}
					{{ end }}
				}
			{{ end }}
			return
		}
		`,
		"Type", typ,
	)
}

func (g *BuildersGenerator) generateListBuilderFile(typ *concepts.Type) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.pkgName(typ.Owner())
	fileName := g.fileName(typ)

	// Create the buffer for the generated code:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Base(g.base).
		Package(pkgName).
		File(fileName).
		Function("builderCtor", g.builderCtor).
		Function("builderName", g.builderName).
		Function("objectName", g.objectName).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateListBuilderSource(typ)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *BuildersGenerator) generateListBuilderSource(typ *concepts.Type) {
	if typ.Element().IsStruct() {
		g.generateStructListBuilderSource(typ)
	}
}

func (g *BuildersGenerator) generateStructListBuilderSource(typ *concepts.Type) {
	g.buffer.Emit(`
		{{ $objectName := objectName .Type }}
		{{ $builderName := builderName .Type }}
		{{ $builderCtor := builderCtor .Type }}
		{{ $elementObjectName := objectName .Type.Element }}
		{{ $elementBuilderName := builderName .Type.Element }}

		// {{ $builderName }} contains the data and logic needed to build
		// '{{ .Type.Element.Name }}' objects.
		type  {{ $builderName }} struct {
			items []*{{ $elementBuilderName }}
		}

		// {{ $builderCtor }} creates a new builder of '{{ .Type.Element.Name }}' objects.
		func {{ $builderCtor }}() *{{ $builderName }} {
			return new({{ $builderName }})
		}

		// Items sets the items of the list.
		func (b *{{ $builderName }}) Items(values ...*{{ $elementBuilderName }}) *{{ $builderName }} {
                        b.items = make([]*{{ $elementBuilderName }}, len(values))
                        copy(b.items, values)
			return b
		}

		// Build creates a list of '{{ .Type.Element.Name }}' objects using the
		// configuration stored in the builder.
		func (b *{{ $builderName }}) Build() (list *{{ $objectName }}, err error) {
			items := make([]*{{ $elementObjectName }}, len(b.items))
			for i, item := range b.items {
				items[i], err = item.Build()
				if err != nil {
					return
				}
			}
			list = new({{ $objectName }})
			list.items = items
			return
		}
		`,
		"Type", typ,
	)
}

func (g *BuildersGenerator) pkgName(version *concepts.Version) string {
	servicePkg := g.names.Package(version.Owner().Name())
	versionPkg := g.names.Package(version.Name())
	return filepath.Join(servicePkg, versionPkg)
}

func (g *BuildersGenerator) fileName(typ *concepts.Type) string {
	return g.names.File(names.Cat(typ.Name(), nomenclator.Builder))
}

func (g *BuildersGenerator) objectName(typ *concepts.Type) string {
	if typ.IsStruct() || typ.IsList() {
		return g.names.Public(typ.Name())
	}
	g.reporter.Errorf(
		"Don't know how to calculate object type name for type '%s'",
		typ.Name(),
	)
	return ""
}

func (g *BuildersGenerator) builderName(typ *concepts.Type) string {
	if typ.IsStruct() || typ.IsList() {
		name := names.Cat(typ.Name(), nomenclator.Builder)
		return g.names.Public(name)
	}
	g.reporter.Errorf(
		"Don't know how to calculate builder type name for type '%s'",
		typ.Name(),
	)
	return ""
}

func (g *BuildersGenerator) builderCtor(typ *concepts.Type) string {
	name := names.Cat(nomenclator.New, typ.Name())
	return g.names.Public(name)
}

func (g *BuildersGenerator) fieldName(attribute *concepts.Attribute) string {
	return g.names.Private(attribute.Name())
}

func (g *BuildersGenerator) fieldType(attribute *concepts.Attribute) *golang.TypeReference {
	typ := attribute.Type()
	switch {
	case typ.IsScalar():
		return g.types.NullableReference(typ)
	case typ.IsList():
		element := typ.Element()
		if element.IsScalar() {
			return g.types.NullableReference(typ)
		}
		ref := g.types.BuilderReference(element)
		ref = g.types.Reference(
			ref.Import(),
			ref.Selector(),
			ref.Name(),
			fmt.Sprintf("[]%s", ref.Text()),
		)
		return ref
	case typ.IsMap():
		element := typ.Element()
		if element.IsScalar() {
			return g.types.NullableReference(typ)
		}
		ref := g.types.BuilderReference(element)
		ref = g.types.Reference(
			ref.Import(),
			ref.Selector(),
			ref.Name(),
			fmt.Sprintf("map[string]%s", ref.Text()),
		)
		return ref
	default:
		return g.types.BuilderReference(typ)
	}
}

func (g *BuildersGenerator) setterName(attribute *concepts.Attribute) string {
	return g.names.Public(attribute.Name())
}

func (g *BuildersGenerator) setterType(attribute *concepts.Attribute) *golang.TypeReference {
	typ := attribute.Type()
	switch {
	case typ.IsScalar():
		return g.types.ValueReference(typ)
	case typ.IsList():
		element := typ.Element()
		if element.IsScalar() {
			return g.types.ValueReference(typ)
		}
		return g.types.BuilderReference(element)
	case typ.IsMap():
		element := typ.Element()
		if element.IsScalar() {
			return g.types.ValueReference(typ)
		}
		ref := g.types.BuilderReference(element)
		ref = g.types.Reference(
			ref.Import(),
			ref.Selector(),
			ref.Name(),
			fmt.Sprintf("map[string]%s", ref.Text()),
		)
		return ref
	default:
		return g.types.BuilderReference(typ)
	}
}

func (g *BuildersGenerator) valueType(typ *concepts.Type) *golang.TypeReference {
	return g.types.ValueReference(typ)
}
