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

	"gitlab.cee.redhat.com/service/ocm-api-metamodel/pkg/concepts"
	"gitlab.cee.redhat.com/service/ocm-api-metamodel/pkg/golang"
	"gitlab.cee.redhat.com/service/ocm-api-metamodel/pkg/names"
	"gitlab.cee.redhat.com/service/ocm-api-metamodel/pkg/nomenclator"
	"gitlab.cee.redhat.com/service/ocm-api-metamodel/pkg/reporter"
)

// ServersGeneratorBuilder is an object used to configure and build the servers generator. Don't create
// instances directly, use the ServersGeneratorBuilder function instead.
type ServersGeneratorBuilder struct {
	reporter *reporter.Reporter
	model    *concepts.Model
	output   string
	base     string
	names    *golang.NamesCalculator
	types    *golang.TypesCalculator
}

// ServersGenerator generate resources for the model resources.
// Don't create instances directly, use the builder instead.
type ServersGenerator struct {
	reporter *reporter.Reporter
	errors   int
	model    *concepts.Model
	output   string
	base     string
	names    *golang.NamesCalculator
	types    *golang.TypesCalculator
	buffer   *golang.Buffer
}

// NewServersGenerator creates a new builder for resource generators.
func NewServersGenerator() *ServersGeneratorBuilder {
	return new(ServersGeneratorBuilder)
}

// Reporter sets the object that will be used to report information about the generation process,
// including errors.
func (b *ServersGeneratorBuilder) Reporter(value *reporter.Reporter) *ServersGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the resource generator.
func (b *ServersGeneratorBuilder) Model(value *concepts.Model) *ServersGeneratorBuilder {
	b.model = value
	return b
}

// Output sets import path of the output package.
func (b *ServersGeneratorBuilder) Output(value string) *ServersGeneratorBuilder {
	b.output = value
	return b
}

// Base sets the import import path of the output package.
func (b *ServersGeneratorBuilder) Base(value string) *ServersGeneratorBuilder {
	b.base = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *ServersGeneratorBuilder) Names(value *golang.NamesCalculator) *ServersGeneratorBuilder {
	b.names = value
	return b
}

// Types sets the object that will be used to calculate types.
func (b *ServersGeneratorBuilder) Types(value *golang.TypesCalculator) *ServersGeneratorBuilder {
	b.types = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new
// types generator using it.
func (b *ServersGeneratorBuilder) Build() (generator *ServersGenerator, err error) {
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
		err = fmt.Errorf("package is mandatory")
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
	generator = new(ServersGenerator)
	generator.reporter = b.reporter
	generator.model = b.model
	generator.output = b.output
	generator.base = b.base
	generator.names = b.names
	generator.types = b.types

	return
}

// Run executes the code generator.
func (g *ServersGenerator) Run() error {
	var err error

	// Generate the Go resource for each model resource:
	for _, service := range g.model.Services() {
		for _, version := range service.Versions() {
			for _, resource := range version.Resources() {
				err = g.generateServerFile(resource)
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

func (g *ServersGenerator) generateServerFile(resource *concepts.Resource) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.pkgName(resource.Owner())
	fileName := g.fileName(resource)

	// Create the buffer for the generated code:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Base(g.base).
		Package(pkgName).
		File(fileName).
		Function("resourceName", g.resourceName).
		Build()
	if err != nil {
		return err
	}

	// Generate the source:
	g.generateServerSource(resource)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ServersGenerator) generateServerSource(resource *concepts.Resource) {
	g.buffer.Emit(`
		{{ $resourceName := resourceName .Resource }}

		// {{ $resourceName }}Server represents the interface the manages the '{{ .Resource.Name }}' resource.
		type {{ $resourceName }}Server interface {}
		`,
		"Resource", resource,
	)
}

func (g *ServersGenerator) pkgName(version *concepts.Version) string {
	servicePkg := g.names.Package(version.Owner().Name())
	versionPkg := g.names.Package(version.Name())
	return filepath.Join(servicePkg, versionPkg)
}

func (g *ServersGenerator) fileName(resource *concepts.Resource) string {
	return g.names.File(names.Cat(resource.Name(), nomenclator.Server))
}

func (g *ServersGenerator) resourceName(resource *concepts.Resource) string {
	return g.names.Public(resource.Name())
}
