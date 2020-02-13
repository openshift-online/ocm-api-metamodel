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

// ServersGeneratorBuilder is an object used to configure and build the servers generator. Don't create
// instances directly, use the ServersGeneratorBuilder function instead.
type ServersGeneratorBuilder struct {
	reporter *reporter.Reporter
	model    *concepts.Model
	output   string
	packages *PackagesCalculator
	names    *NamesCalculator
	types    *TypesCalculator
	binding  *http.BindingCalculator
}

// ServersGenerator generate resources for the model resources.
// Don't create instances directly, use the builder instead.
type ServersGenerator struct {
	reporter *reporter.Reporter
	errors   int
	model    *concepts.Model
	output   string
	packages *PackagesCalculator
	names    *NamesCalculator
	types    *TypesCalculator
	binding  *http.BindingCalculator
	buffer   *Buffer
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

// Packages sets the object that will be used to calculate package names.
func (b *ServersGeneratorBuilder) Packages(
	value *PackagesCalculator) *ServersGeneratorBuilder {
	b.packages = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *ServersGeneratorBuilder) Names(value *NamesCalculator) *ServersGeneratorBuilder {
	b.names = value
	return b
}

// Types sets the object that will be used to calculate types.
func (b *ServersGeneratorBuilder) Types(value *TypesCalculator) *ServersGeneratorBuilder {
	b.types = value
	return b
}

// Binding sets the object that will by used to do HTTP binding calculations.
func (b *ServersGeneratorBuilder) Binding(value *http.BindingCalculator) *ServersGeneratorBuilder {
	b.binding = value
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
	if b.packages == nil {
		err = fmt.Errorf("packages calculator is mandatory")
		return
	}
	if b.names == nil {
		err = fmt.Errorf("names calculator is mandatory")
		return
	}
	if b.types == nil {
		err = fmt.Errorf("types calculator is mandatory")
		return
	}
	if b.binding == nil {
		err = fmt.Errorf("binding calculator is mandatory")
		return
	}

	// Create the generator:
	generator = &ServersGenerator{
		reporter: b.reporter,
		model:    b.model,
		output:   b.output,
		packages: b.packages,
		names:    b.names,
		types:    b.types,
		binding:  b.binding,
	}

	return
}

// Run executes the code generator.
func (g *ServersGenerator) Run() error {
	var err error

	// Calculate the file name:
	fileName := g.names.File(nomenclator.Server)

	// Create the buffer for the model server:
	g.buffer, err = NewBuffer().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		File(fileName).
		Function("serviceName", g.serviceName).
		Function("serviceSelector", g.packages.ServiceSelector).
		Function("serviceSegment", g.binding.ServiceSegment).
		Build()
	if err != nil {
		return err
	}

	// Generate the source for the model server:
	g.generateMainServerSource()
	g.generateMainDispatcherSource()
	err = g.buffer.Write()
	if err != nil {
		return err
	}

	// Generate the server for each service:
	for _, service := range g.model.Services() {
		err = g.generateServiceServer(service)
		if err != nil {
			return err
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

func (g *ServersGenerator) generateMainServerSource() {
	for _, service := range g.model.Services() {
		g.buffer.Import(g.packages.ServiceImport(service), "")
	}
	g.buffer.Emit(`
		// Server is the interface of the top level server.
		type Server interface {
			{{ range .Model.Services }}
				{{ $serviceName := serviceName . }}
				{{ $serviceSelector := serviceSelector . }}

				// {{ $serviceName }} returns the server for service '{{ .Name }}'.
				{{ $serviceName }}() {{ $serviceSelector }}.Server
			{{ end }}
		}
		`,
		"Model", g.model,
	)
}

func (g *ServersGenerator) generateMainDispatcherSource() {
	g.buffer.Import("net/http", "")
	g.buffer.Import(g.packages.ErrorsImport(), "")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Emit(`
		// Dispatch navigates the servers tree till it finds one that matches the given set
		// of path segments, and then invokes it.
		func Dispatch(w http.ResponseWriter, r *http.Request, server Server, segments []string) {
			if len(segments) > 0 && segments[0] == "api" {
				dispatch(w, r, server, segments[1:])
				return
			}
			errors.SendNotFound(w, r)
		}

		func dispatch(w http.ResponseWriter, r *http.Request, server Server, segments []string) {
			if len(segments) == 0 {
				// TODO: This should send the metadata.
				errors.SendMethodNotAllowed(w, r)
				return
			} else {
				switch segments[0] {
				{{ range .Model.Services }}
					{{ $serviceName := serviceName . }}
					{{ $serviceSelector := serviceSelector . }}

					case "{{ serviceSegment . }}":
						service := server.{{ $serviceName }}()
						if service == nil {
							errors.SendNotFound(w, r)
							return
						}
						{{ $serviceSelector }}.Dispatch(w, r, service, segments[1:])
				{{ end }}
				default:
					errors.SendNotFound(w, r)
					return
				}
			}
		}

		// Adapter is an HTTP handler that knows how to translate HTTP requests into calls
		// to the methods of an object that implements the Server interface.
		type Adapter struct {
			server Server
		}

		// NewAdapter creates a new adapter that will translate HTTP requests into calls to
		// the given server.
		func NewAdapter(server Server) *Adapter {
			return &Adapter{
				server: server,
			}
		}

		// ServeHTTP is the implementation of the http.Handler interface.
		func (a *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
			Dispatch(w, r, a.server, helpers.Segments(r.URL.Path))
		}
		`,
		"Model", g.model,
	)
}

func (g *ServersGenerator) generateServiceServer(service *concepts.Service) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.ServicePackage(service)
	fileName := g.names.File(nomenclator.Server)

	// Create the buffer for the service:
	g.buffer, err = NewBuffer().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("serverName", g.serverName).
		Function("versionName", g.versionName).
		Function("versionSelector", g.packages.VersionSelector).
		Function("versionSegment", g.binding.VersionSegment).
		Build()
	if err != nil {
		return err
	}

	// Generate the source for the service:
	g.generateServiceServerSource(service)
	g.generateServiceDispatcherSource(service)
	err = g.buffer.Write()
	if err != nil {
		return err
	}

	// Generate the clients for the versions:
	for _, version := range service.Versions() {
		err = g.generateVersionServer(version)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *ServersGenerator) generateServiceServerSource(service *concepts.Service) {
	for _, version := range service.Versions() {
		g.buffer.Import(g.packages.VersionImport(version), "")
	}
	g.buffer.Emit(`
		// Server is the interface for the '{{ .Service.Name }}' service.
		type Server interface {
			{{ range .Service.Versions }}
				{{ $versionName := versionName . }}
				{{ $versionSelector := versionSelector . }}
				{{ $rootName := serverName .Root }}

				// {{ $versionName }} returns the server for version '{{ .Name }}'.
				{{ $versionName }}() {{ $versionSelector }}.{{ $rootName }}
			{{ end }}
		}
		`,
		"Service", service,
	)
}

func (g *ServersGenerator) generateServiceDispatcherSource(service *concepts.Service) {
	g.buffer.Import("net/http", "")
	g.buffer.Import(g.packages.ErrorsImport(), "")
	g.buffer.Emit(`
		// Dispatch navigates the servers tree till it finds one that matches the given set
		// of path segments, and then invokes it.
		func Dispatch(w http.ResponseWriter, r *http.Request, server Server, segments []string) {
			if len(segments) == 0 {
				// TODO: This should send the service metadata.
				errors.SendMethodNotAllowed(w, r)
				return
			} else {
				switch segments[0] {
				{{ range .Service.Versions }}
					{{ $versionName := versionName . }}
					{{ $versionSelector := versionSelector . }}

					case "{{ versionSegment . }}":
						version := server.{{ $versionName }}()
						if version == nil {
							errors.SendNotFound(w, r)
							return
						}
						{{ $versionSelector }}.Dispatch(w, r, version, segments[1:])
				{{ end }}
				default:
					errors.SendNotFound(w, r)
					return
				}
			}
		}
		`,
		"Service", service,
	)
}

func (g *ServersGenerator) generateVersionServer(version *concepts.Version) error {
	for _, resource := range version.Resources() {
		err := g.generateResourceServer(resource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *ServersGenerator) generateResourceServer(resource *concepts.Resource) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.VersionPackage(resource.Owner())
	fileName := g.fileName(resource)

	// Create the buffer for the generated code:
	g.buffer, err = NewBuffer().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("adaptRequestName", g.adaptRequestName).
		Function("defaultStatus", g.binding.DefaultStatus).
		Function("dispatchName", g.dispatchName).
		Function("fieldName", g.fieldName).
		Function("fieldType", g.fieldType).
		Function("getterName", g.getterName).
		Function("getterType", g.getterType).
		Function("httpMethod", g.binding.Method).
		Function("jsonFieldName", g.jsonFieldName).
		Function("jsonFieldType", g.jsonFieldType).
		Function("locatorName", g.locatorName).
		Function("locatorSegment", g.binding.LocatorSegment).
		Function("methodName", g.methodName).
		Function("methodSegment", g.binding.MethodSegment).
		Function("parameterName", g.binding.ParameterName).
		Function("readRequestFunc", g.readRequestFunc).
		Function("readerName", g.readerName).
		Function("requestBodyParameters", g.binding.RequestBodyParameters).
		Function("requestName", g.requestName).
		Function("requestParameters", g.binding.RequestParameters).
		Function("responseBodyParameters", g.binding.ResponseBodyParameters).
		Function("responseName", g.responseName).
		Function("responseParameters", g.binding.ResponseParameters).
		Function("serverName", g.serverName).
		Function("setterName", g.setterName).
		Function("setterType", g.setterType).
		Function("structName", g.types.StructName).
		Function("writeFunc", g.writeFunc).
		Function("writeResponseFunc", g.writeResponseFunc).
		Function("zeroValue", g.types.ZeroValue).
		Build()
	if err != nil {
		return err
	}

	// Generate the source:
	g.generateResourceServerSource(resource)
	g.generateResourceDispatcherSource(resource)

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
				{{ $methodName }}(ctx context.Context, request *{{$requestName}}, response *{{$responseName}}) error
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

func (g *ServersGenerator) generateResourceDispatcherSource(resource *concepts.Resource) {
	g.buffer.Import("fmt", "")
	g.buffer.Import("net/http", "")
	g.buffer.Import(g.packages.ErrorsImport(), "")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Emit(`
		{{ $serverName := serverName .Resource }}
		{{ $dispatchName := dispatchName .Resource }}

		// {{ $dispatchName }} navigates the servers tree rooted at the given server
		// till it finds one that matches the given set of path segments, and then invokes
		// the corresponding server.
		func {{ $dispatchName }}(w http.ResponseWriter, r *http.Request, server {{ $serverName }}, segments []string) {
			if len(segments) == 0 {
				switch r.Method {
				{{ range .Resource.Methods }}
					{{ $methodSegment := methodSegment . }}
					{{ if not $methodSegment }}
						case "{{ httpMethod . }}":
							{{ adaptRequestName . }}(w, r, server)
							return
					{{ end }}
				{{ end }}
				default:
					errors.SendMethodNotAllowed(w, r)
					return
				}
			}
			switch segments[0] {
			{{ range .Resource.Methods }}
				{{ $methodSegment := methodSegment . }}
				{{ if $methodSegment }}
					case "{{ methodSegment . }}":
						if r.Method != "POST" {
							errors.SendMethodNotAllowed(w, r)
							return
						}
						{{ adaptRequestName . }}(w, r, server)
						return
				{{ end }}
			{{ end }}
			{{ range .Resource.ConstantLocators }}
				case "{{ locatorSegment . }}":
					target := server.{{ locatorName . }}()
					if target == nil {
						errors.SendNotFound(w, r)
						return
					}
					{{ dispatchName .Target }}(w, r, target, segments[1:])
			{{ end }}
			default:
				{{ if .Resource.VariableLocator }}
					{{ with .Resource.VariableLocator }}
						target := server.{{ locatorName . }}(segments[0])
						if target == nil {
							errors.SendNotFound(w, r)
							return
						}
						{{ dispatchName .Target }}(w, r, target, segments[1:])
					{{ end }}
				{{ else }}
					errors.SendNotFound(w, r)
					return
				{{ end }}
			}
		}

		{{ range .Resource.Methods }}
			{{ $methodName := methodName . }}
			{{ $adaptRequestName := adaptRequestName . }}
			{{ $requestName := requestName . }}
			{{ $responseName := responseName . }}
			{{ $requestBodyParameters := requestBodyParameters . }}
			{{ $requestBodyLen := len $requestBodyParameters }}
			{{ $responseParameters := responseParameters . }}

			// {{ $adaptRequestName }} translates the given HTTP request into a call to
			// the corresponding method of the given server. Then it translates the
			// results returned by that method into an HTTP response.
			func {{ $adaptRequestName }}(w http.ResponseWriter, r *http.Request, server {{ $serverName }}) {
				request := &{{ $requestName }}{}
				err := {{ readRequestFunc . }}(request, r)
				if err != nil {
					glog.Errorf(
						"Can't read request for method '%s' and path '%s': %v",
						r.Method, r.URL.Path, err,
					)
					errors.SendInternalServerError(w, r)
					return
				}
				response := &{{ $responseName }}{}
				response.status = {{ defaultStatus . }}
				err = server.{{ $methodName }}(r.Context(), request, response)
				if err != nil {
					glog.Errorf(
						"Can't process request for method '%s' and path '%s': %v",
						r.Method, r.URL.Path, err,
					)
					errors.SendInternalServerError(w, r)
					return
				}
				err = {{ writeResponseFunc . }}(response, w)
				if err != nil {
					glog.Errorf(
						"Can't write response for method '%s' and path '%s': %v",
						r.Method, r.URL.Path, err,
					)
					return
				}
			}
		{{ end }}
		`,
		"Resource", resource,
	)
}

func (g *ServersGenerator) generateRequestSource(method *concepts.Method) {
	// Classify the parameters:
	all := g.binding.RequestBodyParameters(method)
	var main *concepts.Parameter
	var others []*concepts.Parameter
	for _, parameter := range all {
		if parameter.IsItems() || parameter.IsBody() {
			main = parameter
		} else {
			others = append(others, parameter)
		}
	}

	// Generate the code:
	g.buffer.Import("encoding/json", "")
	g.buffer.Import("io", "")
	g.buffer.Emit(`
		{{ $requestName := requestName .Method }}
		{{ $requestParameters := requestParameters .Method }}
		{{ $requestBodyParameters := requestBodyParameters .Method }}
		{{ $requestBodyLen := len $requestBodyParameters }}

		// {{ $requestName }} is the request for the '{{ .Method.Name }}' method.
		type {{ $requestName }} struct {
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
		"Main", main,
		"Others", others,
	)
}

func (g *ServersGenerator) generateResponseSource(method *concepts.Method) {
	// Classify the parameters:
	all := g.binding.ResponseBodyParameters(method)
	var main *concepts.Parameter
	var others []*concepts.Parameter
	for _, parameter := range all {
		if parameter.IsItems() || parameter.IsBody() {
			main = parameter
		} else {
			others = append(others, parameter)
		}
	}

	// Generate the code:
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Import("github.com/json-iterator/go", "")
	g.buffer.Emit(`
		{{ $responseName := responseName .Method }}
		{{ $responseParameters := responseParameters .Method }}
		{{ $responseLen := len $responseParameters }}
		{{ $isAction := .Method.IsAction }}

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
			{{ $fieldType := fieldType . }}
			{{ $setterName := setterName . }}
			{{ $setterType := setterType . }}

			// {{ $setterName }} sets the value of the '{{ .Name }}' parameter.
			//
			{{ lineComment .Doc }}
			func (r *{{ $responseName }}) {{ $setterName }}(value {{ $setterType }}) *{{ $responseName }} {
				{{ if or .IsItems .Type.IsStruct }}
					r.{{ $fieldName }} = value
				{{ else if .Type.IsScalar }}
					r.{{ $fieldName }} = &value
				{{ else if .Type.IsList }}
					if value == nil {
						r.{{ $fieldName }} = nil
					} else {
						r.{{ $fieldName }} = make({{ $fieldType }}, len(value))
						for i, v := range value {
							r.{{ $fieldName }}[i] = v
						}
					}
				{{ else if .Type.IsMap }}
					if value == nil {
						r.{{ $fieldName }} = nil
					} else {
						r.{{ $fieldName }} = {{ $fieldType }}{}
						for k, v := range value {
							r.{{ $fieldName }}[k] = v
						}
					}
				{{ end }}
				return r
			}
		{{ end }}

		// Status sets the status code.
		func (r *{{ $responseName }}) Status(value int) *{{ $responseName }} {
			r.status = value
			return r
		}
		`,
		"Method", method,
		"Main", main,
		"Others", others,
	)
}

func (g *ServersGenerator) serversFile() string {
	return g.names.File(nomenclator.Servers)
}

func (g *ServersGenerator) fileName(resource *concepts.Resource) string {
	return g.names.File(names.Cat(resource.Name(), nomenclator.Server))
}

func (g *ServersGenerator) serviceName(service *concepts.Service) string {
	return g.names.Public(service.Name())
}

func (g *ServersGenerator) versionName(version *concepts.Version) string {
	return g.names.Public(version.Name())
}

func (g *ServersGenerator) serverName(resource *concepts.Resource) string {
	root := resource.Owner().Root()
	if resource == root {
		return g.names.Public(nomenclator.Server)
	}
	return g.names.Public(names.Cat(resource.Name(), nomenclator.Server))
}

func (g *ServersGenerator) locatorName(locator *concepts.Locator) string {
	return g.names.Public(locator.Name())
}

func (g *ServersGenerator) methodName(method *concepts.Method) string {
	return g.names.Public(method.Name())
}

func (g *ServersGenerator) dispatchName(resource *concepts.Resource) string {
	root := resource.Owner().Root()
	if resource == root {
		return g.names.Public(nomenclator.Dispatch)
	}
	return g.names.Private(names.Cat(nomenclator.Dispatch, resource.Name()))
}

func (g *ServersGenerator) adaptRequestName(method *concepts.Method) string {
	name := names.Cat(
		nomenclator.Adapt,
		method.Owner().Name(),
		method.Name(),
		nomenclator.Request,
	)
	return g.names.Private(name)
}

func (g *ServersGenerator) requestName(method *concepts.Method) string {
	resource := method.Owner()
	var name *names.Name
	if resource.IsRoot() {
		name = names.Cat(method.Name(), nomenclator.Server, nomenclator.Request)
	} else {
		name = names.Cat(
			resource.Name(),
			method.Name(),
			nomenclator.Server,
			nomenclator.Request,
		)
	}
	return g.names.Public(name)
}

func (g *ServersGenerator) responseName(method *concepts.Method) string {
	resource := method.Owner()
	var name *names.Name
	if resource.IsRoot() {
		name = names.Cat(
			method.Name(),
			nomenclator.Server,
			nomenclator.Response,
		)
	} else {
		name = names.Cat(
			resource.Name(),
			method.Name(),
			nomenclator.Server,
			nomenclator.Response,
		)
	}
	return g.names.Public(name)
}

func (g *ServersGenerator) fieldName(parameter *concepts.Parameter) string {
	name := g.names.Private(parameter.Name())
	name = g.avoidBuiltin(name, builtinFields)
	return name
}

func (g *ServersGenerator) fieldType(parameter *concepts.Parameter) *TypeReference {
	if parameter.IsItems() {
		return g.types.ListReference(parameter.Type())
	}
	return g.types.NullableReference(parameter.Type())
}

func (g *ServersGenerator) getterName(parameter *concepts.Parameter) string {
	name := g.names.Public(parameter.Name())
	name = g.avoidBuiltin(name, builtinGetters)
	return name
}

func (g *ServersGenerator) getterType(parameter *concepts.Parameter) *TypeReference {
	return g.accessorType(parameter)
}

func (g *ServersGenerator) setterName(parameter *concepts.Parameter) string {
	name := g.names.Public(parameter.Name())
	name = g.avoidBuiltin(name, builtinSetters)
	return name
}

func (g *ServersGenerator) setterType(parameter *concepts.Parameter) *TypeReference {
	return g.accessorType(parameter)
}

func (g *ServersGenerator) jsonFieldName(parameter *concepts.Parameter) string {
	return g.names.Public(parameter.Name())
}

func (g *ServersGenerator) jsonFieldType(parameter *concepts.Parameter) *TypeReference {
	return g.types.JSONTypeReference(parameter.Type())
}

func (g *ServersGenerator) accessorType(parameter *concepts.Parameter) *TypeReference {
	var ref *TypeReference
	typ := parameter.Type()
	switch {
	case parameter.IsItems():
		ref = g.types.ListReference(typ)
	case typ.IsScalar():
		ref = g.types.ValueReference(typ)
	case typ.IsStruct() || typ.IsList() || typ.IsMap():
		ref = g.types.NullableReference(typ)
	}
	if ref == nil {
		g.reporter.Errorf(
			"Don't know how to calculate accessor type for parameter '%s'",
			parameter,
		)
	}
	return ref
}

func (g *ServersGenerator) avoidBuiltin(name string, builtins map[string]interface{}) string {
	_, ok := builtins[name]
	if ok {
		name = name + "_"
	}
	return name
}

func (g *ServersGenerator) readerName(typ *concepts.Type) string {
	version := typ.Owner()
	switch typ {
	case version.Boolean():
		return "ParseBoolean"
	case version.IntegerType():
		return "ParseInteger"
	case version.FloatType():
		return "ParseFloat"
	case version.StringType():
		return "ParseString"
	case version.DateType():
		return "ParseDate"
	default:
		g.reporter.Errorf("We do not know how to handle type '%s'", typ)
		return ""
	}
}

func (g *ServersGenerator) writeFunc(typ *concepts.Type) string {
	name := names.Cat(nomenclator.Write, typ.Name())
	return g.names.Private(name)
}

func (g *ServersGenerator) readRequestFunc(method *concepts.Method) string {
	resource := method.Owner()
	var name *names.Name
	if resource.IsRoot() {
		name = names.Cat(
			nomenclator.Read,
			method.Name(),
			nomenclator.Request,
		)
	} else {
		name = names.Cat(
			nomenclator.Read,
			resource.Name(),
			method.Name(),
			nomenclator.Request,
		)
	}
	return g.names.Private(name)
}

func (g *ServersGenerator) writeResponseFunc(method *concepts.Method) string {
	resource := method.Owner()
	var name *names.Name
	if resource.IsRoot() {
		name = names.Cat(
			nomenclator.Write,
			method.Name(),
			nomenclator.Response,
		)
	} else {
		name = names.Cat(
			nomenclator.Write,
			resource.Name(),
			method.Name(),
			nomenclator.Response,
		)
	}
	return g.names.Private(name)
}
