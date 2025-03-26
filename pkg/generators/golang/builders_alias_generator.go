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
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"

	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
)

// BuildersAliasGeneratorBuilder is an object used to configure and build the builders generator. Don't
// create instances directly, use the NewBuildersAliasGenerator function instead.
type BuildersAliasGeneratorBuilder struct {
	BuildersGeneratorBuilder
}

// BuildersAliasGenerator generates code for the builders of the model types. Don't create instances
// directly, use the builder instead.
type BuildersAliasGenerator struct {
	BuildersGenerator
}

// NewBuildersAliasGenerator creates a new builder for builders generators.
func NewBuildersAliasGenerator() *BuildersAliasGeneratorBuilder {
	return &BuildersAliasGeneratorBuilder{}
}

// Reporter sets the object that will be used to report information about the generation process,
// including errors.
func (b *BuildersAliasGeneratorBuilder) Reporter(value *reporter.Reporter) *BuildersAliasGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the types generator.
func (b *BuildersAliasGeneratorBuilder) Model(value *concepts.Model) *BuildersAliasGeneratorBuilder {
	b.model = value
	return b
}

// Output sets the directory where the source will be generated.
func (b *BuildersAliasGeneratorBuilder) Output(value string) *BuildersAliasGeneratorBuilder {
	b.output = value
	return b
}

// Packages sets the object that will be used to calculate package names.
func (b *BuildersAliasGeneratorBuilder) Packages(
	value *PackagesCalculator) *BuildersAliasGeneratorBuilder {
	b.packages = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *BuildersAliasGeneratorBuilder) Names(value *NamesCalculator) *BuildersAliasGeneratorBuilder {
	b.names = value
	return b
}

// Types sets the object that will be used to calculate types.
func (b *BuildersAliasGeneratorBuilder) Types(value *TypesCalculator) *BuildersAliasGeneratorBuilder {
	b.types = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new
// builders generator using it.
func (b *BuildersAliasGeneratorBuilder) Build() (generator *BuildersAliasGenerator, err error) {
	delegate, err := b.BuildersGeneratorBuilder.Build()
	if err != nil {
		return nil, err
	}

	// Create the generator:
	return &BuildersAliasGenerator{
		BuildersGenerator: *delegate,
	}, nil
}

// Run executes the code generator.
func (g *BuildersAliasGenerator) Run() error {
	var err error

	// Generate the code for each type:
	for _, service := range g.model.Services() {
		for _, version := range service.Versions() {
			for _, typ := range version.Types() {
				switch {
				case typ.IsStruct():
					err = g.generateStructBuilderAliasFile(typ)
				case typ.IsList() && typ.Element().IsStruct():
					err = g.generateListBuilderAliasFile(typ)
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

func (g *BuildersAliasGenerator) generateStructBuilderAliasFile(typ *concepts.Type) error {
	var err error

	// Calculate the package and file names:
	pkgName := g.packages.VersionPackage(typ.Owner())
	fileName := g.fileName(typ)

	// Create the buffer for the generated code:
	g.buffer, err = NewBuffer().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("builderCtor", g.builderCtor).
		Function("builderName", g.builderName).
		Function("apiVersionPackage", func() string { return g.packages.APIVersionSelector(typ.Owner()) }).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateStructBuilderAliasSource(typ)

	// Write the generated code:
	return g.buffer.Write()

}

func (g *BuildersAliasGenerator) generateStructBuilderAliasSource(typ *concepts.Type) {
	g.buffer.Import(g.packages.APIVersionImport(typ.Owner()), g.packages.APIVersionSelector(typ.Owner()))
	g.buffer.Emit(`
		{{ $builderName := builderName .Type }}
		{{ $builderCtor := builderCtor .Type }}

		type  {{ $builderName }} = {{ apiVersionPackage }}.{{ $builderName }}
		var {{ $builderCtor }} = {{ apiVersionPackage }}.{{ $builderCtor }}
		`,
		"Type", typ,
	)
}

func (g *BuildersAliasGenerator) generateListBuilderAliasFile(typ *concepts.Type) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.VersionPackage(typ.Owner())
	fileName := g.fileName(typ)

	// Create the buffer for the generated code:
	g.buffer, err = NewBuffer().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("builderCtor", g.builderCtor).
		Function("builderName", g.builderName).
		Function("apiVersionPackage", func() string { return g.packages.APIVersionSelector(typ.Owner()) }).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateListBuilderAliasSource(typ)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *BuildersAliasGenerator) generateListBuilderAliasSource(typ *concepts.Type) {
	if typ.Element().IsStruct() {
		g.generateStructListBuilderAliasSource(typ)
	}
}

func (g *BuildersAliasGenerator) generateStructListBuilderAliasSource(typ *concepts.Type) {
	g.buffer.Import(g.packages.APIVersionImport(typ.Owner()), g.packages.APIVersionSelector(typ.Owner()))
	g.buffer.Emit(`
		{{ $builderName := builderName .Type }}
		{{ $builderCtor := builderCtor .Type }}

		type  {{ $builderName }} = {{ apiVersionPackage }}.{{ $builderName }}
		var {{ $builderCtor }} = {{ apiVersionPackage }}.{{ $builderCtor }}
		`,
		"Type", typ,
	)
}

func (g *BuildersAliasGenerator) fileName(typ *concepts.Type) string {
	return g.names.File(names.Cat(typ.Name(), nomenclator.BuilderAlias))
}
