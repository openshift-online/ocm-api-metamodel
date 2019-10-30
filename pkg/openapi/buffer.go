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

package openapi

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/json-iterator/go"

	"github.com/openshift-online/ocm-api-metamodel/pkg/golang"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// BufferBuilder is used to create a new OpenAPI buffer. Don't create it directly, use the
// NewBufferBuilder function instead.
type BufferBuilder struct {
	reporter *reporter.Reporter
	output   string
	packages *golang.PackagesCalculator
	pkg      string
}

// Buffer is a type that simplifies the generation of OpenAPI specifications.
type Buffer struct {
	reporter *reporter.Reporter
	output   string
	packages *golang.PackagesCalculator
	pkg      string
	stream   *jsoniter.Stream
	stack    []*int
}

// NewBufferBuilder creates a builder for OpenAPI buffers.
func NewBufferBuilder() *BufferBuilder {
	return &BufferBuilder{}
}

// Reporter sets the object that will be used to report errors and other relevant information.
func (b *BufferBuilder) Reporter(reporter *reporter.Reporter) *BufferBuilder {
	b.reporter = reporter
	return b
}

// Output sets the directory where the OpenAPI specifications will be generated.
func (b *BufferBuilder) Output(value string) *BufferBuilder {
	b.output = value
	return b
}

// Packages sets the object that will be use to calculate package names.
func (b *BufferBuilder) Packages(value *golang.PackagesCalculator) *BufferBuilder {
	b.packages = value
	return b
}

// Package sets the name of the package where the Go file containing the OpenAPI specification will
// be generated.
func (b *BufferBuilder) Package(value string) *BufferBuilder {
	b.pkg = value
	return b
}

// Build creates a new buffer using the configuration stored in the builder.
func (b *BufferBuilder) Build() (buffer *Buffer, err error) {
	// Check that the mandatory parameters have been provided:
	if b.reporter == nil {
		err = fmt.Errorf("reporter is mandatory")
		return
	}
	if b.output == "" {
		err = fmt.Errorf("output directory is mandatory")
		return
	}
	if b.packages == nil {
		err = fmt.Errorf("packages calculator is mandatory")
		return
	}
	if b.pkg == "" {
		err = fmt.Errorf("package name is mandatory")
		return
	}

	// Configure the stream:
	config := jsoniter.Config{
		IndentionStep: 2,
	}
	stream := jsoniter.NewStream(config.Froze(), nil, 0)

	// Create the stack:
	stack := []*int{new(int)}

	// Allocate and populate the buffer:
	buffer = &Buffer{
		reporter: b.reporter,
		output:   b.output,
		packages: b.packages,
		pkg:      b.pkg,
		stream:   stream,
		stack:    stack,
	}

	return
}

func (b *Buffer) peekLevel() *int {
	return b.stack[len(b.stack)-1]
}

func (b *Buffer) pushLevel(level *int) {
	b.stack = append(b.stack, level)
}

func (b *Buffer) popLevel() *int {
	level := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	return level
}

func (b *Buffer) StartObject(names ...string) {
	level := b.peekLevel()
	if *level > 0 {
		b.stream.WriteMore()
	}
	if len(names) > 0 {
		*level++
		b.stream.WriteObjectField(names[0])
		b.stream.WriteObjectStart()
	} else {
		b.stream.WriteObjectStart()
	}
	level = new(int)
	b.pushLevel(level)
}

func (b *Buffer) EndObject() {
	b.stream.WriteObjectEnd()
	b.popLevel()
	level := b.peekLevel()
	*level++
}

func (b *Buffer) Field(name string, value interface{}) {
	level := b.peekLevel()
	if *level > 0 {
		b.stream.WriteMore()
	}
	b.stream.WriteObjectField(name)
	b.stream.WriteVal(value)
	*level++
}

func (b *Buffer) StartArray(names ...string) {
	level := b.peekLevel()
	if *level > 0 {
		b.stream.WriteMore()
	}
	if len(names) > 0 {
		*level++
		b.stream.WriteObjectField(names[0])
		b.stream.WriteArrayStart()
	} else {
		b.stream.WriteArrayStart()
	}
	level = new(int)
	b.pushLevel(level)
}

func (b *Buffer) EndArray() {
	b.stream.WriteArrayEnd()
	b.popLevel()
	level := b.peekLevel()
	*level++
}

func (b *Buffer) Item(value interface{}) {
	level := b.peekLevel()
	if *level > 0 {
		b.stream.WriteMore()
	}
	b.stream.WriteVal(value)
	*level++
}

// Write creates the output file and writes the generated content.
func (b *Buffer) Write() error {
	var err error

	// Check if there were errors:
	errors := b.reporter.Errors()
	if errors > 0 {
		if errors > 1 {
			err = fmt.Errorf("there were %d errors", errors)
		} else {
			err = fmt.Errorf("there was 1 error")
		}
		return err
	}

	// Make sure that the output directory exists:
	file := filepath.Join(b.output, b.pkg, "openapi")
	dir := filepath.Dir(file)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	// Write the JSON file:
	jsonFile := file + ".json"
	b.reporter.Infof("Writing file '%s'", jsonFile)
	jsonData := b.stream.Buffer()
	err = ioutil.WriteFile(jsonFile, jsonData, 0666)
	if err != nil {
		return err
	}

	// Write the YAML file:
	yamlFile := file + ".yaml"
	b.reporter.Infof("Writing file '%s'", yamlFile)
	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(yamlFile, yamlData, 0666)
	if err != nil {
		return err
	}

	// Write the Go file:
	goBuffer, err := golang.NewBufferBuilder().
		Reporter(b.reporter).
		Output(b.output).
		Packages(b.packages).
		Package(b.pkg).
		File(filepath.Base(file)).
		Build()
	if err != nil {
		return err
	}
	goBuffer.Emit(`
		// JSON document containing the OpenAPI specification:
		var openAPI = []byte{
			{{ byteArray .Data }}
		}
		`,
		"Data", jsonData,
	)
	err = goBuffer.Write()
	if err != nil {
		return err
	}

	return nil
}
