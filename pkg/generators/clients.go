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

	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/golang"
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// ClientsGeneratorBuilder is an object used to configure and build a client generator. Don't create
// instances directly, use the NewClientsGenerator function instead.
type ClientsGeneratorBuilder struct {
	reporter *reporter.Reporter
	model    *concepts.Model
	output   string
	base     string
	names    *golang.NamesCalculator
	types    *golang.TypesCalculator
}

// ClientsGenerator generates client code. Don't create instances directly, use the builder instead.
type ClientsGenerator struct {
	reporter *reporter.Reporter
	errors   int
	model    *concepts.Model
	output   string
	base     string
	names    *golang.NamesCalculator
	types    *golang.TypesCalculator
	buffer   *golang.Buffer
}

// NewClientsGenerator creates a new builder for client generators.
func NewClientsGenerator() *ClientsGeneratorBuilder {
	return new(ClientsGeneratorBuilder)
}

// Reporter sets the object that will be used to report information about the generation process,
// including errors.
func (b *ClientsGeneratorBuilder) Reporter(value *reporter.Reporter) *ClientsGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the client generator.
func (b *ClientsGeneratorBuilder) Model(value *concepts.Model) *ClientsGeneratorBuilder {
	b.model = value
	return b
}

// Output sets the output directory.
func (b *ClientsGeneratorBuilder) Output(value string) *ClientsGeneratorBuilder {
	b.output = value
	return b
}

// Base sets the output base package.
func (b *ClientsGeneratorBuilder) Base(value string) *ClientsGeneratorBuilder {
	b.base = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *ClientsGeneratorBuilder) Names(value *golang.NamesCalculator) *ClientsGeneratorBuilder {
	b.names = value
	return b
}

// Types sets the object that will be used to calculate types.
func (b *ClientsGeneratorBuilder) Types(value *golang.TypesCalculator) *ClientsGeneratorBuilder {
	b.types = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new client
// generator using it.
func (b *ClientsGeneratorBuilder) Build() (generator *ClientsGenerator, err error) {
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
		err = fmt.Errorf("path is mandatory")
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
	generator = new(ClientsGenerator)
	generator.reporter = b.reporter
	generator.model = b.model
	generator.output = b.output
	generator.base = b.base
	generator.names = b.names
	generator.types = b.types

	return
}

// Run executes the code generator.
func (g *ClientsGenerator) Run() error {
	var err error

	// Generate the client for each service:
	for _, service := range g.model.Services() {
		g.reporter.Infof("Generating client for service '%s'", service.Name())
		err = g.generateServiceClient(service)
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

func (g *ClientsGenerator) generateServiceClient(service *concepts.Service) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.names.Package(service.Name())
	fileName := g.names.File(nomenclator.Client)

	// Create the buffer for the service:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Base(g.base).
		Package(pkgName).
		File(fileName).
		Function("clientName", g.clientName).
		Function("urlSegment", g.urlSegment).
		Function("versionName", g.versionName).
		Function("versionSelector", g.versionSelector).
		Build()
	if err != nil {
		return err
	}

	// Generate the source for the service:
	err = g.generateServiceClientSource(service)
	if err != nil {
		return err
	}
	err = g.buffer.Write()
	if err != nil {
		return err
	}

	// Generate the clients for the versions:
	for _, version := range service.Versions() {
		err = g.generateVersionClient(version)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *ClientsGenerator) generateServiceClientSource(service *concepts.Service) error {
	g.buffer.Import("net/http", "")
	g.buffer.Import("path", "")
	for _, version := range service.Versions() {
		g.buffer.Import(g.versionImport(version), "")
	}
	g.buffer.Emit(`
		// Client is the client for service '{{ .Service.Name }}'.
		type Client struct {
			transport http.RoundTripper
			path string
			metric string
		}

		// NewClient creates a new client for the service '{{ .Service.Name }}' using the
		// given transport to send the requests and receive the responses.
		func NewClient(transport http.RoundTripper, path string, metric string) *Client {
			client := new(Client)
			client.transport = transport
			client.path = path
			client.metric = metric
			return client
		}

		{{ range .Service.Versions }}
			{{ $versionName := versionName . }}
			{{ $versionSelector := versionSelector . }}
			{{ $versionSegment := urlSegment .Name }}
			{{ $rootName := clientName .Root }}

			// {{ $versionName }} returns a reference to a client for version '{{ .Name }}'.
			func (c *Client) {{ $versionName }}() *{{ $versionSelector }}.{{ $rootName }} {
				return {{ $versionSelector }}.New{{ $rootName }}(
					c.transport,
					path.Join(c.path, "{{ $versionSegment }}"),
					path.Join(c.metric, "{{ $versionSegment }}"),
				)
			}
		{{ end }}
		`,
		"Service", service,
	)

	return nil
}

func (g *ClientsGenerator) generateVersionClient(version *concepts.Version) error {
	for _, resource := range version.Resources() {
		err := g.generateResourceClient(resource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *ClientsGenerator) generateResourceClient(resource *concepts.Resource) error {
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
		Function("clientName", g.clientName).
		Function("dataFieldName", g.dataFieldName).
		Function("dataFieldType", g.dataFieldType).
		Function("dataStruct", g.dataStruct).
		Function("enumName", g.enumName).
		Function("fieldName", g.fieldName).
		Function("fieldTag", g.fieldTag).
		Function("fieldType", g.fieldType).
		Function("getterName", g.getterName).
		Function("getterType", g.getterType).
		Function("httpMethod", g.httpMethod).
		Function("locatorName", g.locatorName).
		Function("methodName", g.methodName).
		Function("requestBodyParameters", g.requestBodyParameters).
		Function("requestData", g.requestData).
		Function("requestName", g.requestName).
		Function("requestParameters", g.requestParameters).
		Function("requestQueryParameters", g.requestQueryParameters).
		Function("responseBodyParameters", g.responseBodyParameters).
		Function("responseData", g.responseData).
		Function("responseName", g.responseName).
		Function("responseParameters", g.responseParameters).
		Function("setterName", g.setterName).
		Function("setterType", g.setterType).
		Function("urlSegment", g.urlSegment).
		Function("zeroValue", g.types.ZeroValue).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateResourceClientSource(resource)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ClientsGenerator) generateResourceClientSource(resource *concepts.Resource) {
	g.buffer.Import("net/http", "")
	g.buffer.Import("path", "")
	g.buffer.Emit(`
		{{ $clientName := clientName .Resource }}

		// {{ $clientName }} is the client of the '{{ .Resource.Name }}' resource.
		//
		{{ lineComment .Resource.Doc }}
		type {{ $clientName }} struct {
			transport http.RoundTripper
			path string
			metric string
		}

		// New{{ $clientName }} creates a new client for the '{{ .Resource.Name }}'
		// resource using the given transport to sned the requests and receive the
		// responses.
		func New{{ $clientName }}(transport http.RoundTripper, path string, metric string) *{{ $clientName }} {
			client := new({{ $clientName }})
			client.transport = transport
			client.path = path
			client.metric = metric
			return client
		}

		{{ range .Resource.Methods }}
			{{ $methodName := methodName . }}
			{{ $requestName := requestName . }}

			// {{ $methodName }} creates a request for the '{{ .Name }}' method.
			//
			{{ lineComment .Doc }}
			func (c *{{ $clientName }}) {{ $methodName }}() *{{ $requestName }} {
				request := new({{ $requestName }})
				request.transport = c.transport
				request.path = c.path
				request.metric = c.metric
				return request
			}
		{{ end }}

		{{ range .Resource.Locators }}
			{{ $locatorName := locatorName . }}
			{{ $locatorSegment := urlSegment .Name }}
			{{ $targetName := clientName .Target }}

			{{ if .Variable }}
				// {{ $locatorName }} returns the target '{{ .Target.Name }}' resource for the given identifier.
				//
				{{ lineComment .Doc }}
				func (c *{{ $clientName }}) {{ $locatorName }}(id string) *{{ $targetName }} {
					return New{{ $targetName }}(
						c.transport,
						path.Join(c.path, id),
						path.Join(c.metric, "-"),
					)
				}
			{{ else }}
				// {{ $locatorName }} returns the target '{{ .Target.Name }}' resource.
				//
				{{ lineComment .Doc }}
				func (c *{{ $clientName }}) {{ $locatorName }}() *{{ $targetName }} {
					return New{{ $targetName }}(
						c.transport,
						path.Join(c.path, "{{ $locatorSegment }}"),
						path.Join(c.metric, "{{ $locatorSegment }}"),
					)
				}
			{{ end }}
		{{ end }}
		`,
		"Resource", resource,
	)

	// Generate the request and response types:
	for _, method := range resource.Methods() {
		g.generateRequestSource(method)
		g.generateResponseSource(method)
	}
}

func (g *ClientsGenerator) generateRequestSource(method *concepts.Method) {
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
		{{ $requestData := requestData .Method }}
		{{ $requestName := requestName .Method }}
		{{ $requestParameters := requestParameters .Method }}
		{{ $requestQueryParameters := requestQueryParameters .Method }}
		{{ $requestBodyParameters := requestBodyParameters .Method }}
		{{ $requestBodyLen := len $requestBodyParameters }}
		{{ $responseData := responseData .Method }}
		{{ $responseName := responseName .Method }}
		{{ $responseParameters := responseParameters .Method }}
		{{ $responseBodyParameters := responseBodyParameters .Method }}

		// {{ $requestName }} is the request for the '{{ .Method.Name }}' method.
		type {{ $requestName }} struct {
			transport http.RoundTripper
			path      string
			metric    string
			query     url.Values
			header    http.Header
			{{ range $requestParameters }}
				{{ fieldName . }} {{ fieldType . }}
			{{ end }}
		}

		// Parameter adds a query parameter.
		func (r *{{ $requestName }}) Parameter(name string, value interface{}) *{{ $requestName }} {
			helpers.AddValue(&r.query, name, value)
			return r
		}

		// Header adds a request header.
		func (r *{{ $requestName }}) Header(name string, value interface{}) *{{ $requestName }} {
			helpers.AddHeader(&r.header, name, value)
			return r
		}

		{{ range $requestParameters }}
			{{ $fieldName := fieldName . }}
			{{ $setterName := setterName . }}
			{{ $setterType := setterType . }}

			// {{ $setterName }} sets the value of the '{{ .Name }}' parameter.
			//
			{{ lineComment .Doc }}
			func (r *{{ $requestName }}) {{ $setterName }}(value {{ $setterType }}) *{{ $requestName }} {
				{{ if or .Type.IsStruct .Type.IsList }}
					r.{{ $fieldName }} = value
				{{ else }}
					r.{{ $fieldName }} = &value
				{{ end }}
				return r
			}
		{{ end }}

		// Send sends this request, waits for the response, and returns it.
		//
		// This is a potentially lengthy operation, as it requires network communication.
		// Consider using a context and the SendContext method.
		func (r *{{ $requestName }}) Send() (result *{{ $responseName }}, err error) {
			return r.SendContext(context.Background())
		}

		// SendContext sends this request, waits for the response, and returns it.
		func (r *{{ $requestName }}) SendContext(ctx context.Context) (result *{{ $responseName }}, err error) {
			query := helpers.CopyQuery(r.query)
			{{ range $requestQueryParameters }}
				{{ $fieldName := fieldName . }}
				{{ $fieldTag := fieldTag . }}
				if r.{{ $fieldName }} != nil {
					helpers.AddValue(&query, "{{ $fieldTag }}", *r.{{ $fieldName }})
				}
			{{ end }}
			header := helpers.SetHeader(r.header, r.metric)
			{{ if $requestBodyParameters }}
				buffer := new(bytes.Buffer)
				err = r.marshal(buffer)
				if err != nil {
					return
				}
			{{ end }}
			uri := &url.URL{
				Path: r.path,
				RawQuery: query.Encode(),
			}
			request := &http.Request{
				Method: {{ httpMethod .Method }},
				URL:    uri,
				Header: header,
				{{ if $requestBodyParameters }}
					Body: ioutil.NopCloser(buffer),
				{{ end }}
			}
			if ctx != nil {
				request = request.WithContext(ctx)
			}
			response, err := r.transport.RoundTrip(request)
			if err != nil {
				return
			}
			defer response.Body.Close()
			result = new({{ $responseName }})
			result.status = response.StatusCode
			result.header = response.Header
			if result.status >= 400 {
				result.err, err = errors.UnmarshalError(response.Body)
				if err != nil {
					return
				}
				err = result.err
				return
			}
			{{ if $responseBodyParameters }}
				err = result.unmarshal(response.Body)
				if err != nil {
					return
				}
			{{ end }}
			return
		}

		{{ if $requestBodyParameters }}
			// marshall is the method used internally to marshal requests for the
			// '{{ .Method.Name }}' method.
			func (r *{{ $requestName }}) marshal(writer io.Writer) error {
				var err error
				encoder := json.NewEncoder(writer)
				{{ if eq $requestBodyLen 1 }}
					{{ with index $requestBodyParameters 0 }}
						data, err := r.{{ fieldName . }}.wrap()
						if err != nil {
							return err
						}
					{{ end }}
				{{ else }}
					data := new({{ $requestData }})
					{{ range $requestBodyParameters }}
						{{ $dataFieldName := dataFieldName . }}
						{{ $fieldName := fieldName . }}
						{{ if or .Type.IsScalar }}
							data.{{ $dataFieldName }} = r.{{ $fieldName }}
						{{ else }}
							data.{{ $dataFieldName }}, err = r.{{ $fieldName }}.wrap()
							if err != nil {
								return err
							}
						{{ end }}
					{{ end }}
				{{ end }}
				err = encoder.Encode(data)
				return err
			}

			{{ if gt $requestBodyLen 1 }}
				// {{ $requestData }} is the structure used internally to write the request of the
				// '{{ .Method.Name }}' method.
				type {{ $requestData }} struct {
					{{ range $requestBodyParameters }}
						{{ dataFieldName . }} {{ dataFieldType . }} "json:\"{{ fieldTag . }},omitempty\""
					{{ end }}
				}
			{{ end }}
		{{ end }}
		`,
		"Method", method,
	)
}

func (g *ClientsGenerator) generateResponseSource(method *concepts.Method) {
	g.buffer.Import("io", "")
	g.buffer.Import("net/http", "")
	g.buffer.Import(path.Join(g.base, g.errorsPkg()), "")
	g.buffer.Emit(`
		{{ $responseName := responseName .Method }}
		{{ $responseData := responseData .Method }}
		{{ $responseParameters := responseParameters .Method }}
		{{ $responseBodyParameters := responseBodyParameters .Method }}
		{{ $responseBodyLen := len $responseBodyParameters }}

		// {{ $responseName }} is the response for the '{{ .Method.Name }}' method.
		type  {{ $responseName }} struct {
			status int
			header http.Header
			err    *errors.Error
			{{ range $responseParameters }}
				{{ fieldName . }} {{ fieldType . }}
			{{ end }}
		}

		// Status returns the response status code.
		func (r *{{ $responseName }}) Status() int {
			return r.status
		}

		// Header returns header of the response.
		func (r *{{ $responseName }}) Header() http.Header {
			return r.header
		}

		// Error returns the response error.
		func (r *{{ $responseName }}) Error() *errors.Error {
			return r.err
		}

		{{ range $responseParameters }}
			{{ $parameterType := .Type.Name.String }}
			{{ $fieldName := fieldName . }}
			{{ $getterName := getterName . }}
			{{ $getterType := getterType . }}

			// {{ $getterName }} returns the value of the '{{ .Name }}' parameter.
			//
			{{ lineComment .Doc }}
			func (r *{{ $responseName }}) {{ $getterName }}() {{ $getterType }} {
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
			func (r *{{ $responseName }}) Get{{ $getterName }}() (value {{ $getterType }}, ok bool) {
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

		{{ if $responseBodyParameters }}
			// unmarshal is the method used internally to unmarshal responses to the
			// '{{ .Method.Name }}' method.
			func (r *{{ $responseName }}) unmarshal(reader io.Reader) error {
				var err error
				decoder := json.NewDecoder(reader)
				{{ if eq $responseBodyLen 1 }}
					{{ with index $responseBodyParameters 0 }}
						data := new({{ dataStruct . }})
					{{ end }}
				{{ else }}
					data := new({{ $responseData }})
				{{ end }}
				err = decoder.Decode(data)
				if err != nil {
					return err
				}
				{{ if eq $responseBodyLen 1 }}
					{{ with index $responseBodyParameters 0 }}
						r.{{ fieldName . }}, err = data.unwrap()
						if err != nil {
							return err
						}
					{{ end }}
				{{ else }}
					{{ range $responseBodyParameters }}
						{{ $dataFieldName := dataFieldName . }}
						{{ $fieldName := fieldName . }}
						{{ if or .Type.IsScalar }}
							r.{{ $fieldName }} = data.{{ $dataFieldName }}
						{{ else }}
							r.{{ $fieldName }}, err = data.{{ $dataFieldName }}.unwrap()
							if err != nil {
								return err
							}
						{{ end }}
					{{ end }}
				{{ end }}
				return err
			}

			{{ if gt $responseBodyLen 1 }}
				// {{ $responseData }} is the structure used internally to unmarshal
				// the response of the '{{ .Method.Name }}' method.
				type {{ $responseData }} struct {
					{{ range $responseBodyParameters }}
						{{ dataFieldName . }} {{ dataFieldType . }} "json:\"{{ fieldTag . }},omitempty\""
					{{ end }}
				}
			{{ end }}
		{{ end }}
		`,
		"Method", method,
	)
}

func (g *ClientsGenerator) errorsPkg() string {
	return g.names.Package(nomenclator.Errors)
}

func (g *ClientsGenerator) helpersPkg() string {
	return g.names.Package(nomenclator.Helpers)
}

func (g *ClientsGenerator) clientsFile() string {
	return g.names.File(nomenclator.Clients)
}

func (g *ClientsGenerator) versionName(version *concepts.Version) string {
	return g.names.Public(version.Name())
}

func (g *ClientsGenerator) versionSelector(version *concepts.Version) string {
	return g.names.Package(version.Name())
}

func (g *ClientsGenerator) serviceImport(service *concepts.Service) string {
	serviceSegment := g.names.Package(service.Name())
	return path.Join(g.base, serviceSegment)
}

func (g *ClientsGenerator) versionImport(version *concepts.Version) string {
	serviceSegment := g.names.Package(version.Owner().Name())
	versionSegment := g.names.Package(version.Name())
	return path.Join(g.base, serviceSegment, versionSegment)
}

func (g *ClientsGenerator) pkgName(version *concepts.Version) string {
	servicePkg := g.names.Package(version.Owner().Name())
	versionPkg := g.names.Package(version.Name())
	return path.Join(servicePkg, versionPkg)
}

func (g *ClientsGenerator) fileName(resource *concepts.Resource) string {
	return g.names.File(names.Cat(resource.Name(), nomenclator.Client))
}

func (g *ClientsGenerator) enumName(typ *concepts.Type) string {
	return g.names.Public(typ.Name())
}

func (g *ClientsGenerator) fieldName(parameter *concepts.Parameter) string {
	name := g.names.Private(parameter.Name())
	name = g.avoidBuiltin(name, builtinFields)
	return name
}

func (g *ClientsGenerator) fieldType(parameter *concepts.Parameter) *golang.TypeReference {
	return g.types.NullableReference(parameter.Type())
}

func (g *ClientsGenerator) dataStruct(parameter *concepts.Parameter) string {
	return g.types.DataReference(parameter.Type()).Name()
}

func (g *ClientsGenerator) dataFieldName(parameter *concepts.Parameter) string {
	return g.names.Public(parameter.Name())
}

func (g *ClientsGenerator) dataFieldType(parameter *concepts.Parameter) *golang.TypeReference {
	return g.types.DataReference(parameter.Type())
}

func (g *ClientsGenerator) fieldTag(parameter *concepts.Parameter) string {
	return g.names.Tag(parameter.Name())
}

func (g *ClientsGenerator) urlSegment(name *names.Name) string {
	return g.names.Tag(name)
}

func (g *ClientsGenerator) getterName(parameter *concepts.Parameter) string {
	name := g.names.Public(parameter.Name())
	name = g.avoidBuiltin(name, builtinGetters)
	return name
}

func (g *ClientsGenerator) getterType(parameter *concepts.Parameter) *golang.TypeReference {
	return g.accessorType(parameter.Type())
}

func (g *ClientsGenerator) setterName(parameter *concepts.Parameter) string {
	name := g.names.Public(parameter.Name())
	name = g.avoidBuiltin(name, builtinSetters)
	return name
}

func (g *ClientsGenerator) setterType(parameter *concepts.Parameter) *golang.TypeReference {
	return g.accessorType(parameter.Type())
}

func (g *ClientsGenerator) accessorType(typ *concepts.Type) *golang.TypeReference {
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

func (g *ClientsGenerator) locatorName(locator *concepts.Locator) string {
	return g.names.Public(locator.Name())
}

func (g *ClientsGenerator) methodName(method *concepts.Method) string {
	return g.names.Public(method.Name())
}

func (g *ClientsGenerator) clientName(resource *concepts.Resource) string {
	name := names.Cat(resource.Name(), nomenclator.Client)
	return g.names.Public(name)
}

func (g *ClientsGenerator) requestName(method *concepts.Method) string {
	name := names.Cat(method.Owner().Name(), method.Name(), nomenclator.Request)
	return g.names.Public(name)
}

func (g *ClientsGenerator) requestData(method *concepts.Method) string {
	name := names.Cat(method.Owner().Name(), method.Name(), nomenclator.Request, nomenclator.Data)
	return g.names.Private(name)
}

func (g *ClientsGenerator) responseName(method *concepts.Method) string {
	name := names.Cat(method.Owner().Name(), method.Name(), nomenclator.Response)
	return g.names.Public(name)
}

func (g *ClientsGenerator) responseData(method *concepts.Method) string {
	name := names.Cat(method.Owner().Name(), method.Name(), nomenclator.Response, nomenclator.Data)
	return g.names.Private(name)
}

func (g *ClientsGenerator) httpMethod(method *concepts.Method) string {
	name := method.Name()
	switch {
	case nomenclator.Get.Equals(name) || nomenclator.List.Equals(name):
		return "http.MethodGet"
	case nomenclator.Update.Equals(name):
		return "http.MethodPatch"
	case nomenclator.Delete.Equals(name):
		return "http.MethodDelete"
	default:
		return "http.MethodPost"
	}
}

func (g *ClientsGenerator) requestParameters(method *concepts.Method) []*concepts.Parameter {
	result := make([]*concepts.Parameter, 0)
	for _, parameter := range method.Parameters() {
		if parameter.In() {
			result = append(result, parameter)
		}
	}
	return result
}

func (g *ClientsGenerator) responseParameters(method *concepts.Method) []*concepts.Parameter {
	result := make([]*concepts.Parameter, 0)
	for _, parameter := range method.Parameters() {
		if parameter.Out() {
			result = append(result, parameter)
		}
	}
	return result
}

func (g *ClientsGenerator) requestQueryParameters(method *concepts.Method) []*concepts.Parameter {
	result := make([]*concepts.Parameter, 0)
	for _, parameter := range method.Parameters() {
		if parameter.In() && parameter.Type().IsScalar() {
			result = append(result, parameter)
		}
	}
	return result
}

func (g *ClientsGenerator) requestBodyParameters(method *concepts.Method) []*concepts.Parameter {
	result := make([]*concepts.Parameter, 0)
	for _, parameter := range method.Parameters() {
		if parameter.In() && (parameter.Type().IsStruct() || parameter.Type().IsList()) {
			result = append(result, parameter)
		}
	}
	return result
}

func (g *ClientsGenerator) responseBodyParameters(method *concepts.Method) []*concepts.Parameter {
	result := make([]*concepts.Parameter, 0)
	for _, parameter := range method.Parameters() {
		if parameter.Out() {
			result = append(result, parameter)
		}
	}
	return result
}

func (g *ClientsGenerator) avoidBuiltin(name string, builtins map[string]interface{}) string {
	_, ok := builtins[name]
	if ok {
		name = name + "_"
	}
	return name
}

var builtinFields = map[string]interface{}{
	"err":    nil,
	"status": nil,
}

var builtinGetters = map[string]interface{}{
	"Error":  nil,
	"Status": nil,
}

var builtinSetters = map[string]interface{}{
	"Error":  nil,
	"Status": nil,
}
