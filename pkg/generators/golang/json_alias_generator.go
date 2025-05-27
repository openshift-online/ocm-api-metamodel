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

// JSONSupportAliasGeneratorBuilder is an object used to configure and build the JSON support generator.
// Don't create instances directly, use the NewJSONSupportAliasGenerator function instead.
type JSONSupportAliasGeneratorBuilder struct {
	JSONSupportGeneratorBuilder
}

// JSONSupportAliasGenerator generates JSON support code. Don't create instances directly, use the
// builder instead.
type JSONSupportAliasGenerator struct {
	JSONSupportGenerator
}

// NewJSONSupportAliasGenerator creates a new builder JSON support code generators.
func NewJSONSupportAliasGenerator() *JSONSupportAliasGeneratorBuilder {
	return &JSONSupportAliasGeneratorBuilder{}
}

// Reporter sets the object that will be used to report information about the generation process,
// including errors.
func (b *JSONSupportAliasGeneratorBuilder) Reporter(
	value *reporter.Reporter) *JSONSupportAliasGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the types generator.
func (b *JSONSupportAliasGeneratorBuilder) Model(value *concepts.Model) *JSONSupportAliasGeneratorBuilder {
	b.model = value
	return b
}

// Output sets import path of the output package.
func (b *JSONSupportAliasGeneratorBuilder) Output(value string) *JSONSupportAliasGeneratorBuilder {
	b.output = value
	return b
}

// Package sets the object that will be used to calculate package names.
func (b *JSONSupportAliasGeneratorBuilder) Packages(
	value *PackagesCalculator) *JSONSupportAliasGeneratorBuilder {
	b.packages = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *JSONSupportAliasGeneratorBuilder) Names(value *NamesCalculator) *JSONSupportAliasGeneratorBuilder {
	b.names = value
	return b
}

// Types sets the object that will be used to calculate types.
func (b *JSONSupportAliasGeneratorBuilder) Types(value *TypesCalculator) *JSONSupportAliasGeneratorBuilder {
	b.types = value
	return b
}

// Binding sets the object that will by used to do HTTP binding calculations.
func (b *JSONSupportAliasGeneratorBuilder) Binding(
	value *http.BindingCalculator) *JSONSupportAliasGeneratorBuilder {
	b.binding = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new types
// generator using it.
func (b *JSONSupportAliasGeneratorBuilder) Build() (generator *JSONSupportAliasGenerator, err error) {
	delegate, err := b.JSONSupportGeneratorBuilder.Build()
	if err != nil {
		return nil, err
	}

	// Create the generator:
	return &JSONSupportAliasGenerator{
		JSONSupportGenerator: *delegate,
	}, nil
}

// Run executes the code generator.
func (g *JSONSupportAliasGenerator) Run() error {
	var err error

	// Generate the code for each type:
	for _, service := range g.model.Services() {
		for _, version := range service.Versions() {
			// Generate the code for the version metadata type:
			err := g.generateVersionMetadataAliasSupport(version)
			if err != nil {
				return err
			}

			// Generate the code for the model types:
			var importRefs []struct {
				path     string
				selector string
			}
			for _, typ := range version.Types() {
				for _, att := range typ.Attributes() {
					if att.LinkOwner() != nil {
						importRefs = append(importRefs,
							struct {
								path     string
								selector string
							}{
								path:     g.packages.VersionImport(att.LinkOwner()),
								selector: g.packages.VersionSelector(att.LinkOwner()),
							})
					}
				}
				switch {
				case typ.IsStruct():
					err = g.generateStructTypeAliasSupport(version, typ, importRefs)
				case typ.IsList():
					element := typ.Element()
					if element.IsScalar() || element.IsStruct() {
						err = g.generateListTypeAliasSupport(version, typ, importRefs)
					}
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

func (g *JSONSupportAliasGenerator) generateVersionMetadataAliasSupport(version *concepts.Version) error {
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
		Function("readTypeFunc", g.readTypeFunc).
		Function("writeTypeFunc", g.writeTypeFunc).
		Function("apiVersionPackage", func() string { return g.packages.APIVersionSelector(version) }).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateVersionMetadataAliasSource(version)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *JSONSupportAliasGenerator) generateVersionMetadataAliasSource(version *concepts.Version) {
	g.buffer.Import("fmt", "")
	g.buffer.Import("io", "")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Import("github.com/json-iterator/go", "jsoniter")
	g.buffer.Import(g.packages.APIVersionImport(version), g.packages.APIVersionSelector(version))
	g.buffer.Emit(`
		var MarshalMetadata = {{ apiVersionPackage }}.MarshalMetadata
		var UnmarshalMetadata = {{ apiVersionPackage }}.UnmarshalMetadata
		`,
		"Version", version,
	)
}

func (g *JSONSupportAliasGenerator) generateStructTypeAliasSupport(version *concepts.Version, typ *concepts.Type,
	importRefs []struct {
		path     string
		selector string
	}) error {
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
		Function("marshalTypeFunc", g.marshalTypeFunc).
		Function("readTypeFunc", g.readTypeFunc).
		Function("unmarshalTypeFunc", g.unmarshalTypeFunc).
		Function("writeTypeFunc", g.writeTypeFunc).
		Function("apiVersionPackage", func() string { return g.packages.APIVersionSelector(version) }).
		Build()
	if err != nil {
		return err
	}

	for _, ref := range importRefs {
		g.buffer.Import(ref.path, ref.selector)
	}

	// Generate the code:
	g.generateStructTypeAliasSource(version, typ)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *JSONSupportAliasGenerator) generateStructTypeAliasSource(version *concepts.Version, typ *concepts.Type) {
	g.buffer.Import("fmt", "")
	g.buffer.Import("io", "")
	g.buffer.Import("time", "")
	g.buffer.Import("github.com/json-iterator/go", "jsoniter")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Import(g.packages.APIVersionImport(version), g.packages.APIVersionSelector(version))
	g.buffer.Emit(`
		{{ $marshalTypeFunc := marshalTypeFunc .Type }}
		{{ $writeTypeFunc := writeTypeFunc .Type }}
		{{ $unmarshalTypeFunc := unmarshalTypeFunc .Type }}
		{{ $readTypeFunc := readTypeFunc .Type }}

		var {{ $marshalTypeFunc }} = {{ apiVersionPackage }}.{{ $marshalTypeFunc }}
		var {{ $writeTypeFunc }} = {{ apiVersionPackage }}.{{ $writeTypeFunc }}
		var {{ $unmarshalTypeFunc }} = {{ apiVersionPackage }}.{{ $unmarshalTypeFunc }}
		var {{ $readTypeFunc }} = {{ apiVersionPackage }}.{{ $readTypeFunc }}
		`,
		"Type", typ,
	)
}

func (g *JSONSupportAliasGenerator) generateListTypeAliasSupport(version *concepts.Version, typ *concepts.Type,
	importRefs []struct {
		path     string
		selector string
	}) error {
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
		Function("marshalTypeFunc", g.marshalTypeFunc).
		Function("readTypeFunc", g.readTypeFunc).
		Function("unmarshalTypeFunc", g.unmarshalTypeFunc).
		Function("writeTypeFunc", g.writeTypeFunc).
		Function("apiVersionPackage", func() string { return g.packages.APIVersionSelector(version) }).
		Build()
	if err != nil {
		return err
	}

	for _, ref := range importRefs {
		g.buffer.Import(ref.path, ref.selector)
	}
	g.buffer.Import(g.packages.APIVersionImport(version), g.packages.APIVersionSelector(version))

	// Generate the code:
	g.generateListTypeAliasSource(version, typ)
	// Write the generated code:
	return g.buffer.Write()
}

func (g *JSONSupportAliasGenerator) generateListTypeAliasSource(
	version *concepts.Version,
	typ *concepts.Type,
) {
	var linkOwner *concepts.Version
	g.buffer.Import("fmt", "")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Import("github.com/json-iterator/go", "jsoniter")
	g.buffer.Emit(`
		{{ $marshalTypeFunc := marshalTypeFunc .Type }}
		{{ $writeTypeFunc := writeTypeFunc .Type }}
		{{ $unmarshalTypeFunc := unmarshalTypeFunc .Type }}
		{{ $readTypeFunc := readTypeFunc .Type }}

		var {{ $marshalTypeFunc }} = {{ apiVersionPackage }}.{{ $marshalTypeFunc }}
		var {{ $writeTypeFunc }} = {{ apiVersionPackage }}.{{ $writeTypeFunc }}
		var {{ $unmarshalTypeFunc }} = {{ apiVersionPackage }}.{{ $unmarshalTypeFunc }}
		var {{ $readTypeFunc }} = {{ apiVersionPackage }}.{{ $readTypeFunc }}
		`,
		"Type", typ,
		"linkOwner", linkOwner,
	)
}

func (g *JSONSupportAliasGenerator) metadataFile() string {
	return g.names.File(names.Cat(nomenclator.MetadataAlias, nomenclator.Reader))
}

func (g *JSONSupportAliasGenerator) typeFile(typ *concepts.Type) string {
	return g.names.File(names.Cat(typ.Name(), nomenclator.TypeAlias, nomenclator.JSONAlias))
}
