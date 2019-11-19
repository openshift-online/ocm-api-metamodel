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

	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/golang"
	"github.com/openshift-online/ocm-api-metamodel/pkg/http"
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// ReadersGeneratorBuilder is an object used to configure and build the JSON readers generator.
// Don't create instances directly, use the NewReadersGenerator function instead.
type ReadersGeneratorBuilder struct {
	reporter *reporter.Reporter
	model    *concepts.Model
	output   string
	packages *golang.PackagesCalculator
	names    *golang.NamesCalculator
	types    *golang.TypesCalculator
	binding  *http.BindingCalculator
}

// ReadersGenerator generates code for the JSON readers. Don't create instances directly, use the
// builder instead.
type ReadersGenerator struct {
	reporter *reporter.Reporter
	errors   int
	model    *concepts.Model
	output   string
	packages *golang.PackagesCalculator
	names    *golang.NamesCalculator
	types    *golang.TypesCalculator
	buffer   *golang.Buffer
	binding  *http.BindingCalculator
}

// NewReadersGenerator creates a new builder JSON readers generators.
func NewReadersGenerator() *ReadersGeneratorBuilder {
	return &ReadersGeneratorBuilder{}
}

// Reporter sets the object that will be used to report information about the generation process,
// including errors.
func (b *ReadersGeneratorBuilder) Reporter(value *reporter.Reporter) *ReadersGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the types generator.
func (b *ReadersGeneratorBuilder) Model(value *concepts.Model) *ReadersGeneratorBuilder {
	b.model = value
	return b
}

// Output sets import path of the output package.
func (b *ReadersGeneratorBuilder) Output(value string) *ReadersGeneratorBuilder {
	b.output = value
	return b
}

// Package sets the object that will be used to calculate package names.
func (b *ReadersGeneratorBuilder) Packages(
	value *golang.PackagesCalculator) *ReadersGeneratorBuilder {
	b.packages = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *ReadersGeneratorBuilder) Names(value *golang.NamesCalculator) *ReadersGeneratorBuilder {
	b.names = value
	return b
}

// Types sets the object that will be used to calculate types.
func (b *ReadersGeneratorBuilder) Types(value *golang.TypesCalculator) *ReadersGeneratorBuilder {
	b.types = value
	return b
}

// Binding sets the object that will by used to do HTTP binding calculations.
func (b *ReadersGeneratorBuilder) Binding(value *http.BindingCalculator) *ReadersGeneratorBuilder {
	b.binding = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new types
// generator using it.
func (b *ReadersGeneratorBuilder) Build() (generator *ReadersGenerator, err error) {
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

	// Create the generator:
	generator = &ReadersGenerator{
		reporter: b.reporter,
		model:    b.model,
		output:   b.output,
		packages: b.packages,
		names:    b.names,
		types:    b.types,
	}

	return
}

// Run executes the code generator.
func (g *ReadersGenerator) Run() error {
	var err error

	// Generate the helpers:
	err = g.generateHelpers()
	if err != nil {
		return err
	}

	// Generate the code for each type:
	for _, service := range g.model.Services() {
		for _, version := range service.Versions() {
			// Generate the code for the version metadata type:
			err := g.generateVersionMetadataReader(version)
			if err != nil {
				return err
			}

			// Generate the code for the model types:
			for _, typ := range version.Types() {
				switch {
				case typ.IsStruct():
					err = g.generateStructReader(typ)
				case typ.IsList() && typ.Element().IsStruct():
					err = g.generateListReader(typ)
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

func (g *ReadersGenerator) generateHelpers() error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.HelpersPackage()
	fileName := g.readersFile()

	// Create the buffer for the generated code:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.buffer.Import("bytes", "")
	g.buffer.Import("encoding/json", "")
	g.buffer.Import("fmt", "")
	g.buffer.Import("io", "")
	g.buffer.Import("net/url", "")
	g.buffer.Import("strconv", "")
	g.buffer.Emit(`
		// NewEncoder creates a new JSON encoder from the given target. The target can be a
		// a writer or a JSON encoder.
		func NewEncoder(target interface{}) (encoder *json.Encoder, err error) {
			switch output := target.(type) {
			case io.Writer:
				encoder = json.NewEncoder(output)
			case *json.Encoder:
				encoder = output
			default:
				err = fmt.Errorf(
					"expected writer or JSON decoder, but got %T",
					output,
				)
			}
			return
		}

		// NewDecoder creates a new JSON decoder from the given source. The source can be a
		// slice of bytes, a string, a reader or a JSON decoder.
		func NewDecoder(source interface{}) (decoder *json.Decoder, err error) {
			switch input := source.(type) {
			case []byte:
				decoder = json.NewDecoder(bytes.NewBuffer(input))
			case string:
				decoder = json.NewDecoder(bytes.NewBufferString(input))
			case io.Reader:
				decoder = json.NewDecoder(input)
			case *json.Decoder:
				decoder = input
			default:
				err = fmt.Errorf(
					"expected bytes, string, reader or JSON decoder, but got %T",
					input,
				)
			}
			return
		}

		// ParseInteger reads a string and parses it to integer,
		// if an error occurred it returns a non-nil error.
		func ParseInteger(query url.Values, parameterName string) (*int, error) {
			values := query[parameterName]
			count := len(values)
			if count == 0 {
				return nil, nil
			}
		   	if count > 1 {
				err := fmt.Errorf(
					"expected at most one value for parameter '%s' but got %d",
					parameterName, count,
				)
				return nil, err
			}
			value := values[0]
			parsedInt64, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"value '%s' isn't valid for the '%s' parameter because it isn't an integer: %v",
					value, parameterName, err,
				)
			}
			parsedInt := int(parsedInt64)
			return &parsedInt, nil
		}

		// ParseFloat reads a string and parses it to float,
		// if an error occurred it returns a non-nil error.
		func ParseFloat(query url.Values, parameterName string) (*float64, error) {
			values := query[parameterName]
			count := len(values)
			if count == 0 {
				return nil, nil
			}
		   	if count > 1 {
				err := fmt.Errorf(
					"expected at most one value for parameter '%s' but got %d",
					parameterName, count,
				)
				return nil, err
			}
			value := values[0]
			parsedFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"value '%s' isn't valid for the '%s' parameter because it isn't a float: %v",
					value, parameterName, err,
				)
			}
			return &parsedFloat, nil
		}

		// ParseString returns a pointer to the string and nil error.
		func ParseString(query url.Values, parameterName string) (*string, error) {
			values := query[parameterName]
			count := len(values)
			if count == 0 {
				return nil, nil
			}
		   	if count > 1 {
				err := fmt.Errorf(
					"expected at most one value for parameter '%s' but got %d",
					parameterName, count,
				)
				return nil, err
			}
			return &values[0], nil
		}

		// ParseBoolean reads a string and parses it to boolean,
		// if an error occurred it returns a non-nil error.
		func ParseBoolean(query url.Values, parameterName string) (*bool, error) {
			values := query[parameterName]
			count := len(values)
			if count == 0 {
				return nil, nil
			}
		   	if count > 1 {
				err := fmt.Errorf(
					"expected at most one value for parameter '%s' but got %d",
					parameterName, count,
				)
				return nil, err
			}
			value := values[0]
			parsedBool, err := strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf(
					"value '%s' isn't valid for the '%s' parameter because it isn't a boolean: %v",
					value, parameterName, err,
				)
			}
			return &parsedBool, nil
		}

		// ParseDate reads a string and parses it to a time.Time,
		// if an error occurred it returns a non-nil error.
		func ParseDate(query url.Values, parameterName string) (*time.Time, error) {
			values := query[parameterName]
			count := len(values)
			if count == 0 {
				return nil, nil
			}
		   	if count > 1 {
				err := fmt.Errorf(
					"expected at most one value for parameter '%s' but got %d",
					parameterName, count,
				)
				return nil, err
			}
			value := values[0]
			parsedTime, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return nil, fmt.Errorf(
					"value '%s' isn't valid for the '%s' parameter because it isn't a date: %v",
					value, parameterName, err,
				)
			}
			return &parsedTime, nil
		}
	`)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ReadersGenerator) generateVersionMetadataReader(version *concepts.Version) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.VersionPackage(version)
	fileName := g.metadataFile()

	// Create the buffer for the generated code:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateVersionMetadataReaderSource(version)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ReadersGenerator) generateVersionMetadataReaderSource(version *concepts.Version) {
	g.buffer.Import("fmt", "")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Emit(`
		// metadataData is the data structure used internally to marshal and unmarshal
		// metadata.
		type metadataData struct {
			ServerVersion *string "json:\"server_version,omitempty\""
		}

		// MarshalMetadata writes a value of the metadata type to the given target, which
		// can be a writer or a JSON encoder.
		func MarshalMetadata(object *Metadata, target interface{}) error {
			encoder, err := helpers.NewEncoder(target)
			if err != nil {
				return err
			}
			data, err := object.wrap()
			if err != nil {
				return err
			}
			return encoder.Encode(data)
		}

		// wrap is the method used internally to convert a metadata object to a JSON
		// document.
		func (m *Metadata) wrap() (data *metadataData, err error) {
			if m == nil {
				return
			}
			data = &metadataData{
				ServerVersion: m.serverVersion,
			}
			return
		}

		// UnmarshalMetadata reads a value of the metadata type from the given source, which
		// which can be an slice of bytes, a string, a reader or a JSON decoder.
		func UnmarshalMetadata(source interface{}) (object *Metadata, err error) {
			decoder, err := helpers.NewDecoder(source)
			if err != nil {
				return
			}
			data := &metadataData{}
			err = decoder.Decode(data)
			if err != nil {
				return
			}
			object, err = data.unwrap()
			return
		}

		// unwrap is the function used internally to convert the JSON unmarshalled data to a
		// value of the metadata type.
		func (d *metadataData) unwrap() (object *Metadata, err error) {
			if d == nil {
				return
			}
			object = &Metadata{
				serverVersion: d.ServerVersion,
			}
			return
		}
		`,
	)
}

func (g *ReadersGenerator) generateStructReader(typ *concepts.Type) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.VersionPackage(typ.Owner())
	fileName := g.readerFile(typ)

	// Create the buffer for the generated code:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("attributeName", g.binding.AttributeName).
		Function("dataFieldName", g.dataFieldName).
		Function("dataFieldType", g.dataFieldType).
		Function("dataStruct", g.dataStruct).
		Function("marshalFunc", g.marshalFunc).
		Function("objectFieldName", g.objectFieldName).
		Function("objectFieldType", g.objectFieldType).
		Function("objectName", g.objectName).
		Function("unmarshalFunc", g.unmarshalFunc).
		Function("valueType", g.valueType).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateStructReaderSource(typ)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ReadersGenerator) generateStructReaderSource(typ *concepts.Type) {
	g.buffer.Import("fmt", "")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Emit(`
		{{ $objectName := objectName .Type }}
		{{ $dataStruct := dataStruct .Type }}
		{{ $marshalFunc := marshalFunc .Type }}
		{{ $unmarshalFunc := unmarshalFunc .Type }}

		// {{ $dataStruct }} is the data structure used internally to marshal and unmarshal
		// objects of type '{{ .Type.Name }}'.
		type {{ $dataStruct }} struct {
			{{ if .Type.IsClass }}
				Kind *string "json:\"kind,omitempty\""
				ID   *string "json:\"id,omitempty\""
				HREF *string "json:\"href,omitempty\""
			{{ end }}
			{{ range .Type.Attributes }}
				{{ dataFieldName . }} {{ dataFieldType . }} "json:\"{{ attributeName . }},omitempty\""
			{{ end }}
		}

		// {{ $marshalFunc }} writes a value of the '{{ .Type.Name }}' to the given target,
		// which can be a writer or a JSON encoder.
		func {{ $marshalFunc }}(object *{{ $objectName }}, target interface{}) error {
			encoder, err := helpers.NewEncoder(target)
			if err != nil {
				return err
			}
			data, err := object.wrap()
			if err != nil {
				return err
			}
			return encoder.Encode(data)
		}

		// wrap is the method used internally to convert a value of the '{{ .Type.Name }}'
		// value to a JSON document.
		func (o *{{ $objectName }}) wrap() (data *{{ $dataStruct }}, err error) {
			if o == nil {
				return
			}
			data = new({{ $dataStruct }})
			{{ if .Type.IsClass }}
				data.ID = o.id
				data.HREF = o.href
				data.Kind = new(string)
				if o.link {
					*data.Kind = {{ $objectName }}LinkKind
				} else {
					*data.Kind = {{ $objectName }}Kind
				}
			{{ end }}
			{{ range .Type.Attributes }}
				{{ $dataFieldName := dataFieldName . }}
				{{ $dataFieldType := dataFieldType . }}
				{{ $objectFieldName := objectFieldName . }}
				{{ if .Type.IsList }}
					{{ if .Type.Element.IsScalar }}
						data.{{ $dataFieldName }} = o.{{ $objectFieldName }}
					{{ else if .Type.Element.IsStruct }}
						{{ if and .Link .Type.Element.IsClass }}
							data.{{ $dataFieldName }}, err = o.{{ $objectFieldName }}.wrapLink()
						{{ else }}
							data.{{ $dataFieldName }}, err = o.{{ $objectFieldName }}.wrap()
						{{ end }}
						if err != nil {
							return
						}
					{{ end }}
				{{ else if .Type.IsMap }}
					{{ if .Type.Element.IsScalar }}
						data.{{ $dataFieldName }} = o.{{ $objectFieldName }}
					{{ else if .Type.Element.IsStruct }}
						data.{{ $dataFieldName }} = make({{ $dataFieldType }})
						for key, value := range o.{{ $objectFieldName }} {
							data.{{ $dataFieldName }}[key], err = value.wrap()
							if err != nil {
								return
							}
						}
					{{ end }}
				{{ else if .Type.IsStruct }}
					data.{{ $dataFieldName }}, err = o.{{ $objectFieldName }}.wrap()
					if err != nil {
						return
					}
				{{ else }}
					data.{{ $dataFieldName }} = o.{{ $objectFieldName }}
				{{ end }}
			{{ end }}
			return
		}

		// {{ $unmarshalFunc }} reads a value of the '{{ .Type.Name }}' type from the given
		// source, which can be an slice of bytes, a string, a reader or a JSON decoder.
		func {{ $unmarshalFunc }}(source interface{}) (object *{{ $objectName }}, err error) {
			decoder, err := helpers.NewDecoder(source)
			if err != nil {
				return
			}
			data := new({{ $dataStruct }})
			err = decoder.Decode(data)
			if err != nil {
				return
			}
			object, err = data.unwrap()
			return
		}

		// unwrap is the function used internally to convert the JSON unmarshalled data to a
		// value of the '{{ .Type.Name }}' type.
		func (d *{{ $dataStruct }}) unwrap() (object *{{ $objectName }}, err error) {
			if d == nil {
				return
			}
			object = new({{ $objectName }})
			{{ if .Type.IsClass }}
				object.id = d.ID
				object.href = d.HREF
				if d.Kind != nil {
					switch *d.Kind {
					case {{ $objectName }}Kind:
						object.link = false
					case {{ $objectName }}LinkKind:
						object.link = true
					default:
						err = fmt.Errorf(
							"expected kind '%s' or '%s' but got '%s'",
							{{ $objectName }}Kind,
							{{ $objectName }}LinkKind,
							*d.Kind,
						)
						return
					}
				}
			{{ end }}
			{{ range .Type.Attributes }}
				{{ $dataFieldName := dataFieldName . }}
				{{ $objectFieldName := objectFieldName . }}
				{{ $objectFieldType := objectFieldType . }}
				{{ if .Type.IsList }}
					{{ if .Type.Element.IsScalar }}
						object.{{ $objectFieldName }} = d.{{ $dataFieldName }}
					{{ else if .Type.Element.IsStruct }}
						{{ if and .Link .Type.Element.IsClass }}
							object.{{ $objectFieldName }}, err = d.{{ $dataFieldName }}.unwrapLink()
						{{ else }}
							object.{{ $objectFieldName }}, err = d.{{ $dataFieldName }}.unwrap()
						{{ end }}
						if err != nil {
							return
						}
					{{ end }}
				{{ else if .Type.IsMap }}
					{{ if .Type.Element.IsScalar }}
						object.{{ $objectFieldName }} = d.{{ $dataFieldName }}
					{{ else if .Type.Element.IsStruct }}
						object.{{ $objectFieldName }} = make({{ $objectFieldType }})
						for key, value := range d.{{ $dataFieldName }} {
							object.{{ $objectFieldName }}[key], err = value.unwrap()
							if err != nil {
								return
							}
						}
					{{ end }}
				{{ else if .Type.IsStruct }}
					object.{{ $objectFieldName }}, err = d.{{ $dataFieldName }}.unwrap()
					if err != nil {
						return
					}
				{{ else }}
					object.{{ $objectFieldName }} = d.{{ $dataFieldName }}
				{{ end }}
			{{ end }}
			return
		}
		`,
		"Type", typ,
	)
}

func (g *ReadersGenerator) generateListReader(typ *concepts.Type) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.VersionPackage(typ.Owner())
	fileName := g.readerFile(typ)

	// Create the buffer for the generated code:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("dataList", g.dataList).
		Function("dataStruct", g.dataStruct).
		Function("linkStruct", g.linkStruct).
		Function("marshalFunc", g.marshalFunc).
		Function("objectList", g.objectList).
		Function("objectName", g.objectName).
		Function("unmarshalFunc", g.unmarshalFunc).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.generateListReaderSource(typ)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ReadersGenerator) generateListReaderSource(typ *concepts.Type) {
	g.buffer.Import("fmt", "")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Emit(`
		{{ $objectName := objectName .Type.Element }}
		{{ $objectList := objectList .Type }}
		{{ $dataStruct := dataStruct .Type.Element }}
		{{ $dataList := dataList .Type }}
		{{ $marshalFunc := marshalFunc .Type }}
		{{ $unmarshalFunc := unmarshalFunc .Type }}

		// {{ $dataList }} is type used internally to marshal and unmarshal lists of objects
		// of type '{{ .Type.Element.Name }}'.
		type {{ $dataList }} []*{{ $dataStruct }}

		// {{ $unmarshalFunc }} reads a list of values of the '{{ .Type.Element.Name }}'
		// from the given source, which can be a slice of bytes, a string, an io.Reader or a
		// json.Decoder.
		func {{ $unmarshalFunc }}(source interface{}) (list *{{ $objectList }}, err error) {
			decoder, err := helpers.NewDecoder(source)
			if err != nil {
				return
			}
			var data {{ $dataList }}
			err = decoder.Decode(&data)
			if err != nil {
				return
			}
			list, err = data.unwrap()
			return
		}

		// wrap is the method used internally to convert a list of values of the
		// '{{ .Type.Element.Name }}' value to a JSON document.
		func (l *{{ $objectList }}) wrap() (data {{ $dataList }}, err error) {
			if l == nil {
				return
			}
			data = make({{ $dataList }}, len(l.items))
			for i, item := range l.items {
				data[i], err = item.wrap()
				if err != nil {
					return
				}
			}
			return
		}

		// unwrap is the function used internally to convert the JSON unmarshalled data to a
		// list of values of the '{{ .Type.Element.Name }}' type.
		func (d {{ $dataList }}) unwrap() (list *{{ $objectList }}, err error) {
			if d == nil {
				return
			}
			items := make([]*{{ $objectName }}, len(d))
			for i, item := range d {
				items[i], err = item.unwrap()
				if err != nil {
					return
				}
			}
			list = new({{ $objectList }})
			list.items = items
			return
		}

		{{ if .Type.Element.IsClass }}
			{{ $linkStruct := linkStruct .Type }}

			// {{ $linkStruct }} is type used internally to marshal and unmarshal links
			// to lists of objects of type '{{ .Type.Element.Name }}'.
			type {{ $linkStruct }} struct {
				Kind *string "json:\"kind,omitempty\""
				HREF *string "json:\"href,omitempty\""
				Items []*{{ $dataStruct }} "json:\"items,omitempty\""
			}

			// wrapLink is the method used internally to convert a list of values of the
			// '{{ .Type.Element.Name }}' value to a link.
			func (l *{{ $objectList }}) wrapLink() (data *{{ $linkStruct }}, err error) {
				if l == nil {
					return
				}
				items := make([]*{{ $dataStruct }}, len(l.items))
				for i, item := range l.items {
					items[i], err = item.wrap()
					if err != nil {
						return
					}
				}
				data = new({{ $linkStruct }})
				data.Items = items
				data.HREF = l.href
				data.Kind = new(string)
				if l.link {
					*data.Kind = {{ $objectName }}ListLinkKind
				} else {
					*data.Kind = {{ $objectName }}ListKind
				}
				return
			}

			// unwrapLink is the function used internally to convert a JSON link to a list
			// of values of the '{{ .Type.Element.Name }}' type to a list.
			func (d *{{ $linkStruct }}) unwrapLink() (list *{{ $objectList }}, err error) {
				if d == nil {
					return
				}
				items := make([]*{{ $objectName }}, len(d.Items))
				for i, item := range d.Items {
					items[i], err = item.unwrap()
					if err != nil {
						return
					}
				}
				list = new({{ $objectList }})
				list.items = items
				list.href = d.HREF
				if d.Kind != nil {
					switch *d.Kind {
					case {{ $objectName }}ListKind:
						list.link = false
					case {{ $objectName }}ListLinkKind:
						list.link = true
					default:
						err = fmt.Errorf(
							"expected kind '%s' or '%s' but got '%s'",
							{{ $objectName }}ListKind,
							{{ $objectName }}ListLinkKind,
							*d.Kind,
						)
						return
					}
				}
				return
			}
		{{ end }}
		`,
		"Type", typ,
	)
}

func (g *ReadersGenerator) readersFile() string {
	return g.names.File(nomenclator.Readers)
}

func (g *ReadersGenerator) metadataFile() string {
	return g.names.File(names.Cat(nomenclator.Metadata, nomenclator.Reader))
}

func (g *ReadersGenerator) readerFile(typ *concepts.Type) string {
	return g.names.File(names.Cat(typ.Name(), nomenclator.Reader))
}

func (g *ReadersGenerator) marshalFunc(typ *concepts.Type) string {
	name := names.Cat(nomenclator.Marshal, typ.Name())
	return g.names.Public(name)
}

func (g *ReadersGenerator) unmarshalFunc(typ *concepts.Type) string {
	name := names.Cat(nomenclator.Unmarshal, typ.Name())
	return g.names.Public(name)
}

func (g *ReadersGenerator) objectName(typ *concepts.Type) string {
	return g.types.ValueReference(typ).Name()
}

func (g *ReadersGenerator) objectList(typ *concepts.Type) string {
	return g.types.ValueReference(typ).Name()
}

func (g *ReadersGenerator) dataStruct(typ *concepts.Type) string {
	return g.types.DataReference(typ).Name()
}

func (g *ReadersGenerator) linkStruct(typ *concepts.Type) string {
	return g.types.LinkDataReference(typ).Name()
}

func (g *ReadersGenerator) dataList(typ *concepts.Type) string {
	return g.types.DataReference(typ).Name()
}

func (g *ReadersGenerator) dataFieldName(attribute *concepts.Attribute) string {
	return g.names.Public(attribute.Name())
}

func (g *ReadersGenerator) dataFieldType(attribute *concepts.Attribute) *golang.TypeReference {
	if attribute.Link() {
		return g.types.LinkDataReference(attribute.Type())
	}
	return g.types.DataReference(attribute.Type())
}

func (g *ReadersGenerator) objectFieldName(attribute *concepts.Attribute) string {
	return g.names.Private(attribute.Name())
}

func (g *ReadersGenerator) objectFieldType(attribute *concepts.Attribute) *golang.TypeReference {
	return g.types.NullableReference(attribute.Type())
}

func (g *ReadersGenerator) valueType(typ *concepts.Type) *golang.TypeReference {
	return g.types.ValueReference(typ)
}
