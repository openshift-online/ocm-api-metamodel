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
	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/http"
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// TypesAliasGeneratorBuilder is an object used to configure and build the types generator. Don't create
// instances directly, use the NewTypesAliasGenerator function instead.
type TypesAliasGeneratorBuilder struct {
	TypesGeneratorBuilder
}

// TypesAliasGenerator Go types for the model types. Don't create instances directly, use the builder
// instead.
type TypesAliasGenerator struct {
	TypesGenerator
}

// NewTypesAliasGenerator creates a new builder for types generators.
func NewTypesAliasGenerator() *TypesAliasGeneratorBuilder {
	return &TypesAliasGeneratorBuilder{}
}

// Reporter sets the object that will be used to report information about the generation process,
// including errors.
func (b *TypesAliasGeneratorBuilder) Reporter(value *reporter.Reporter) *TypesAliasGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the types generator.
func (b *TypesAliasGeneratorBuilder) Model(value *concepts.Model) *TypesAliasGeneratorBuilder {
	b.model = value
	return b
}

// Output sets import path of the output package.
func (b *TypesAliasGeneratorBuilder) Output(value string) *TypesAliasGeneratorBuilder {
	b.output = value
	return b
}

// Packages sets the object that will be used to calculate package names.
func (b *TypesAliasGeneratorBuilder) Packages(value *PackagesCalculator) *TypesAliasGeneratorBuilder {
	b.packages = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *TypesAliasGeneratorBuilder) Names(value *NamesCalculator) *TypesAliasGeneratorBuilder {
	b.names = value
	return b
}

// Types sets the object that will be used to calculate types.
func (b *TypesAliasGeneratorBuilder) Types(value *TypesCalculator) *TypesAliasGeneratorBuilder {
	b.types = value
	return b
}

// Binding sets the object that will be used to calculate HTTP names.
func (b *TypesAliasGeneratorBuilder) Binding(value *http.BindingCalculator) *TypesAliasGeneratorBuilder {
	b.binding = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new
// types generator using it.
func (b *TypesAliasGeneratorBuilder) Build() (generator *TypesAliasGenerator, err error) {
	delegate, err := b.TypesGeneratorBuilder.Build()
	if err != nil {
		return nil, err
	}

	// Create the generator:
	return &TypesAliasGenerator{
		TypesGenerator: *delegate,
	}, nil
}

// Run executes the code generator.
func (g *TypesAliasGenerator) Run() error {
	var err error

	// Generate the go types:
	for _, service := range g.model.Services() {
		for _, version := range service.Versions() {
			// Generate the version metadata type:
			err := g.generateVersionMetadataTypeAliasFile(version)
			if err != nil {
				return err
			}

			// Generate the Go types that correspond to model types:
			for _, typ := range version.Types() {
				switch {
				case typ.IsEnum() || typ.IsStruct():
					err = g.generateTypeAliasFile(typ)
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

func (g *TypesAliasGenerator) generateVersionMetadataTypeAliasFile(version *concepts.Version) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.VersionPackage(version)
	fileName := g.metadataFile()

	// Create the buffer for the generated code:
	g.buffer, err = NewBuffer().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("apiVersionPackage", func() string { return g.packages.APIVersionSelector(version) }).
		Build()
	if err != nil {
		return err
	}

	// Generate the source:
	g.generateVersionMetadataTypeAliasSource(version)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *TypesAliasGenerator) generateVersionMetadataTypeAliasSource(version *concepts.Version) {
	g.buffer.Import(g.packages.APIVersionImport(version), g.packages.APIVersionSelector(version))
	g.buffer.Emit(`
		type Metadata = {{ apiVersionPackage }}.Metadata
		`,
	)
}

func (g *TypesAliasGenerator) generateTypeAliasFile(typ *concepts.Type) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.VersionPackage(typ.Owner())
	fileName := g.typeFile(typ)

	// Create the buffer for the generated code:
	g.buffer, err = NewBuffer().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("enumName", g.types.EnumName).
		Function("listName", g.listName).
		Function("objectName", g.objectName).
		Function("valueName", g.valueName).
		Function("apiVersionPackage", func() string { return g.packages.APIVersionSelector(typ.Owner()) }).
		Build()
	if err != nil {
		return err
	}

	// Generate the source:
	g.generateTypeAliasSource(typ)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *TypesAliasGenerator) generateTypeAliasSource(typ *concepts.Type) {
	switch {
	case typ.IsEnum():
		g.generateEnumTypeAliasSource(typ)
	case typ.IsStruct():
		g.generateStructTypeAliasSource(typ)
	}
}

func (g *TypesAliasGenerator) generateEnumTypeAliasSource(typ *concepts.Type) {
	g.buffer.Import(g.packages.APIVersionImport(typ.Owner()), g.packages.APIVersionSelector(typ.Owner()))
	g.buffer.Emit(`
		{{ $enumName := enumName .Type }}

		// {{ $enumName }} represents the values of the '{{ .Type.Name }}' enumerated type.
		type {{ $enumName }} = {{ apiVersionPackage }}.{{ $enumName }}

		const (
			{{ range .Type.Values }}
				{{ lineComment .Doc }}
				{{ valueName . }} {{ $enumName }} = {{ apiVersionPackage }}.{{ valueName . }}
			{{ end }}
		)
		`,
		"Type", typ,
	)
}

func (g *TypesAliasGenerator) generateStructTypeAliasSource(typ *concepts.Type) {
	g.buffer.Import(g.packages.APIVersionImport(typ.Owner()), g.packages.APIVersionSelector(typ.Owner()))
	g.buffer.Import("time", "")
	g.buffer.Emit(`
		{{ $objectName := objectName .Type }}
		{{ $listName := listName .Type }}

		{{ if .Type.IsClass }}
			// {{ $objectName }}Kind is the name of the type used to represent objects
			// of type '{{ .Type.Name }}'.
			const {{ $objectName }}Kind = {{ apiVersionPackage }}.{{ $objectName }}Kind

			// {{ $objectName }}LinkKind is the name of the type used to represent links
			// to objects of type '{{ .Type.Name }}'.
			const {{ $objectName }}LinkKind = {{ apiVersionPackage }}.{{ $objectName }}LinkKind

			// {{ $objectName }}NilKind is the name of the type used to nil references
			// to objects of type '{{ .Type.Name }}'.
			const {{ $objectName }}NilKind = {{ apiVersionPackage }}.{{ $objectName }}NilKind
		{{ end }}

		// {{ $objectName }} represents the values of the '{{ .Type.Name }}' type.
		//
		{{ lineComment .Type.Doc }}
		type  {{ $objectName }} = {{ apiVersionPackage }}.{{ $objectName }}

		type  {{ $listName }} = {{ apiVersionPackage }}.{{ $listName }}
		`,
		"Type", typ,
	)
}

func (g *TypesAliasGenerator) metadataFile() string {
	return g.names.File(names.Cat(nomenclator.Metadata, nomenclator.TypeAlias))
}

func (g *TypesAliasGenerator) typeFile(typ *concepts.Type) string {
	return g.names.File(names.Cat(typ.Name(), nomenclator.TypeAlias))
}
