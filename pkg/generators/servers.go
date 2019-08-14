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
	"path"
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

	// Generate the Go server for each model resource:
	for _, service := range g.model.Services() {
		for _, version := range service.Versions() {
			for _, resource := range version.Resources() {
				err = g.generateResourceServerFile(resource)
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

func (g *ServersGenerator) generateResourceServerFile(resource *concepts.Resource) error {
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
		Function("methodName", g.methodName).
		Function("serverName", g.serverName).
		Function("locatorName", g.locatorName).
		Function("fieldName", g.fieldName).
		Function("fieldType", g.fieldType).
		Function("getterName", g.getterName).
		Function("getterType", g.getterType).
		Function("requestName", g.requestName).
		Function("responseName", g.responseName).
		Function("methodName", g.methodName).
		Function("requestParameters", g.requestParameters).
		Function("responseParameters", g.responseParameters).
		Function("setterName", g.setterName).
		Function("setterType", g.setterType).
		Function("zeroValue", g.types.ZeroValue).
		Build()
	if err != nil {
		return err
	}

	// Generate the source:
	g.generateResourceServerSource(resource)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ServersGenerator) generateResourceServerSource(resource *concepts.Resource) {
	g.buffer.Emit(`
		{{ $serverName := serverName .Resource }}

		// {{ $serverName }} represents the interface the manages the '{{ .Resource.Name }}' resource.
		type {{ $serverName }} interface {
			{{ range .Resource.Methods }}
				{{ $methodName := methodName . }}
				{{ $responseName := responseName . }}
				{{ $requestName := requestName . }}
				// {{ $methodName }} handles a request for the '{{ .Name }}' method.
				//
				{{ lineComment .Doc }}
				{{ $methodName }}(request *{{$requestName}}, response *{{$responseName}}) error
			{{ end }}

			{{ range .Resource.Locators }}
				{{ $locatorName := locatorName . }}
				{{ $targetName := serverName .Target }}

				{{ if .Variable }}
					// {{ $locatorName }} returns the target '{{ .Target.Name }}' server for the given identifier.
					//
					{{ lineComment .Doc }}
					{{ $locatorName }}(id string) {{ $targetName }}
				{{ else }}
					// {{ $locatorName }} returns the target '{{ .Target.Name }}' resource.
					//
					{{ lineComment .Doc }}
					{{ $locatorName }}() {{ $targetName }}
				{{ end }}
			{{ end }}
		}
		`,
		"Resource", resource,
	)

	// Generate the request and response types:
	for _, method := range resource.Methods() {
		g.generateRequestSource(method)
		g.generateResponseSource(method)
	}
}

func (g *ServersGenerator) generateRequestSource(method *concepts.Method) {
	g.buffer.Import("bytes", "")
	g.buffer.Import("context", "")
	g.buffer.Import("encoding/json", "")
	g.buffer.Import("fmt", "")
	g.buffer.Import("io/ioutil", "")
	g.buffer.Import("net/http", "")
	g.buffer.Import("net/url", "")
	g.buffer.Import(path.Join(g.base, g.errorsPkg()), "")
	g.buffer.Import(path.Join(g.base, g.helpersPkg()), "")
	g.buffer.Emit(`
		{{ $requestName := requestName .Method }}
		{{ $requestParameters := requestParameters .Method }}

		// {{ $requestName }} is the request for the '{{ .Method.Name }}' method.
		type {{ $requestName }} struct {
			path      string
			query     url.Values
			{{ range $requestParameters }}
				{{ fieldName . }} {{ fieldType . }}
			{{ end }}
		}

		{{ range $requestParameters }}
			{{ $parameterType := .Type.Name.String }}
			{{ $fieldName := fieldName . }}
			{{ $getterName := getterName . }}
			{{ $getterType := getterType . }}

			// {{ $getterName }} returns the value of the '{{ .Name }}' parameter.
			//
			{{ lineComment .Doc }}
			func (r *{{ $requestName }}) {{ $getterName }}() {{ $getterType }} {
				{{ if or .Type.IsStruct .Type.IsList .Type.IsMap }}
					if r == nil {
						return nil
					}
					return r.{{ $fieldName }}
				{{ else }}
					if r != nil && r.{{ $fieldName }} != nil {
						return *r.{{ $fieldName }}
					}
					return {{ zeroValue .Type }}
				{{ end }}
			}

			// Get{{ $getterName }} returns the value of the '{{ .Name }}' parameter and
			// a flag indicating if the parameter has a value.
			//
			{{ lineComment .Doc }}
			func (r *{{ $requestName }}) Get{{ $getterName }}() (value {{ $getterType }}, ok bool) {
				ok = r != nil && r.{{ $fieldName }} != nil
				if ok {
					{{ if or .Type.IsStruct .Type.IsList .Type.IsMap }}
						value = r.{{ $fieldName }}
					{{ else }}
						value = *r.{{ $fieldName }}
					{{ end }}
				}
				return
			}
		{{ end }}
		`,
		"Method", method,
	)
}

func (g *ServersGenerator) generateResponseSource(method *concepts.Method) {
	g.buffer.Import("io", "")
	g.buffer.Import("net/http", "")
	g.buffer.Import(path.Join(g.base, g.errorsPkg()), "")
	g.buffer.Emit(`
		{{ $responseName := responseName .Method }}
		{{ $responseParameters := responseParameters .Method }}

		// {{ $responseName }} is the response for the '{{ .Method.Name }}' method.
		type  {{ $responseName }} struct {
			status int
			err    *errors.Error
			{{ range $responseParameters }}
				{{ fieldName . }} {{ fieldType . }}
			{{ end }}
		}
		
		
		{{ range $responseParameters }}
			{{ $fieldName := fieldName . }}
			{{ $setterName := setterName . }}
			{{ $setterType := setterType . }}

			// {{ $setterName }} sets the value of the '{{ .Name }}' parameter.
			//
			{{ lineComment .Doc }}
			func (r *{{ $responseName }}) {{ $setterName }}(value {{ $setterType }}) *{{ $responseName }} {
				{{ if or .Type.IsStruct .Type.IsList }}
					r.{{ $fieldName }} = value
				{{ else }}
					r.{{ $fieldName }} = &value
				{{ end }}
				return r
			}
		{{ end }}
		`,
		"Method", method,
	)
}

func (g *ServersGenerator) errorsPkg() string {
	return g.names.Package(nomenclator.Errors)
}

func (g *ServersGenerator) helpersPkg() string {
	return g.names.Package(nomenclator.Helpers)
}

func (g *ServersGenerator) serversFile() string {
	return g.names.File(nomenclator.Servers)
}

func (g *ServersGenerator) pkgName(version *concepts.Version) string {
	servicePkg := g.names.Package(version.Owner().Name())
	versionPkg := g.names.Package(version.Name())
	return filepath.Join(servicePkg, versionPkg)
}

func (g *ServersGenerator) fileName(resource *concepts.Resource) string {
	return g.names.File(names.Cat(resource.Name(), nomenclator.Server))
}

func (g *ServersGenerator) serverName(resource *concepts.Resource) string {
	return g.names.Public(names.Cat(resource.Name(), nomenclator.Server))
}

func (g *ServersGenerator) locatorName(locator *concepts.Locator) string {
	return g.names.Public(locator.Name())
}

func (g *ServersGenerator) methodName(method *concepts.Method) string {
	return g.names.Public(method.Name())
}

func (g *ServersGenerator) requestName(method *concepts.Method) string {
	name := names.Cat(method.Owner().Name(), method.Name(), nomenclator.Server, nomenclator.Request)
	return g.names.Public(name)
}

func (g *ServersGenerator) requestParameters(method *concepts.Method) []*concepts.Parameter {
	result := make([]*concepts.Parameter, 0)
	for _, parameter := range method.Parameters() {
		if parameter.In() {
			result = append(result, parameter)
		}
	}
	return result
}

func (g *ServersGenerator) responseName(method *concepts.Method) string {
	name := names.Cat(method.Owner().Name(), method.Name(), nomenclator.Server, nomenclator.Response)
	return g.names.Public(name)
}

func (g *ServersGenerator) responseParameters(method *concepts.Method) []*concepts.Parameter {
	result := make([]*concepts.Parameter, 0)
	for _, parameter := range method.Parameters() {
		if parameter.Out() {
			result = append(result, parameter)
		}
	}
	return result
}

func (g *ServersGenerator) fieldName(parameter *concepts.Parameter) string {
	name := g.names.Private(parameter.Name())
	name = g.avoidBuiltin(name, builtinFields)
	return name
}

func (g *ServersGenerator) fieldType(parameter *concepts.Parameter) *golang.TypeReference {
	return g.types.NullableReference(parameter.Type())
}

func (g *ServersGenerator) getterName(parameter *concepts.Parameter) string {
	name := g.names.Public(parameter.Name())
	name = g.avoidBuiltin(name, builtinGetters)
	return name
}

func (g *ServersGenerator) getterType(parameter *concepts.Parameter) *golang.TypeReference {
	return g.accessorType(parameter.Type())
}

func (g *ServersGenerator) setterName(parameter *concepts.Parameter) string {
	name := g.names.Public(parameter.Name())
	name = g.avoidBuiltin(name, builtinSetters)
	return name
}

func (g *ServersGenerator) setterType(parameter *concepts.Parameter) *golang.TypeReference {
	return g.accessorType(parameter.Type())
}

func (g *ServersGenerator) accessorType(typ *concepts.Type) *golang.TypeReference {
	switch {
	case typ.IsList():
		element := typ.Element()
		switch {
		case element.IsStruct():
			name := g.names.Public(names.Cat(element.Name(), nomenclator.List))
			return g.types.Reference("", "", "", "*"+name)
		default:
			return g.types.NullableReference(typ)
		}
	case typ.IsStruct():
		return g.types.NullableReference(typ)
	default:
		return g.types.ValueReference(typ)
	}
}

func (g *ServersGenerator) avoidBuiltin(name string, builtins map[string]interface{}) string {
	_, ok := builtins[name]
	if ok {
		name = name + "_"
	}
	return name
}
