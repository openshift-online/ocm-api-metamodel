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
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// HelpersGeneratorBuilder is an object used to configure and build a helpers generator. Don't
// create instances directly, use the NewHelpersGenerator function instead.
type HelpersGeneratorBuilder struct {
	reporter *reporter.Reporter
	model    *concepts.Model
	output   string
	base     string
	packages *golang.PackagesCalculator
	names    *golang.NamesCalculator
}

// HelpersGenerator generates helper code. Don't create instances directly, use the builder instead.
type HelpersGenerator struct {
	reporter *reporter.Reporter
	errors   int
	model    *concepts.Model
	output   string
	base     string
	packages *golang.PackagesCalculator
	names    *golang.NamesCalculator
	buffer   *golang.Buffer
}

// NewHelpersGenerator creates a new builder for helpers generators.
func NewHelpersGenerator() *HelpersGeneratorBuilder {
	return new(HelpersGeneratorBuilder)
}

// Reporter sets the object that will be used to report information about the generation process,
// including errors.
func (b *HelpersGeneratorBuilder) Reporter(value *reporter.Reporter) *HelpersGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the client generator.
func (b *HelpersGeneratorBuilder) Model(value *concepts.Model) *HelpersGeneratorBuilder {
	b.model = value
	return b
}

// Output sets the output directory.
func (b *HelpersGeneratorBuilder) Output(value string) *HelpersGeneratorBuilder {
	b.output = value
	return b
}

// Base sets the output base package.
func (b *HelpersGeneratorBuilder) Base(value string) *HelpersGeneratorBuilder {
	b.base = value
	return b
}

// Packages sets the object that will be used to calculate package names.
func (b *HelpersGeneratorBuilder) Packages(
	value *golang.PackagesCalculator) *HelpersGeneratorBuilder {
	b.packages = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *HelpersGeneratorBuilder) Names(value *golang.NamesCalculator) *HelpersGeneratorBuilder {
	b.names = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new client
// generator using it.
func (b *HelpersGeneratorBuilder) Build() (generator *HelpersGenerator, err error) {
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
	if b.packages == nil {
		err = fmt.Errorf("packages calculator is mandatory")
		return
	}
	if b.names == nil {
		err = fmt.Errorf("names calculator is mandatory")
		return
	}

	// Create the generator:
	generator = &HelpersGenerator{
		reporter: b.reporter,
		model:    b.model,
		output:   b.output,
		base:     b.base,
		packages: b.packages,
		names:    b.names,
	}

	return
}

// Run executes the code generator.
func (g *HelpersGenerator) Run() error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.HelpersPackage()
	fileName := g.helpersFile()

	// Create the buffer for the generated code:
	g.buffer, err = golang.NewBufferBuilder().
		Reporter(g.reporter).
		Output(g.output).
		Base(g.base).
		Package(pkgName).
		File(fileName).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	g.buffer.Import("fmt", "")
	g.buffer.Import("net/http", "")
	g.buffer.Import("net/url", "")
	g.buffer.Import("strings", "")
	g.buffer.Import("time", "")
	g.buffer.Emit(`
		// AddValue creates the given set of query parameters if needed, an then adds
		// the given parameter.
		func AddValue(query *url.Values, name string, value interface{}) {
			if *query == nil {
				*query = make(url.Values)
			}
			query.Add(name, fmt.Sprintf("%v", value))
		}

		// CopyQuery creates a copy of the given set of query parameters.
		func CopyQuery(query url.Values) url.Values {
			if query == nil {
				return nil
			}
			result := make(url.Values)
			for name, values := range query {
				result[name] = CopyValues(values)
			}
			return result
		}

		// AddHeader creates the given set of headers if needed, and then adds the given
		// header:
		func AddHeader(header *http.Header, name string, value interface{}) {
			if *header == nil {
				*header = make(http.Header)
			}
			header.Add(name, fmt.Sprintf("%v", value))
		}

		// SetHeader creates a copy of the given set of headers, and adds the header
		// containing the given metrics path.
		func SetHeader(header http.Header, metric string) http.Header {
			result := make(http.Header)
			for name, values := range header {
				result[name] = CopyValues(values)
			}
			result.Set(metricHeader, metric)
			return result
		}

		// CopyValues copies a slice of strings.
		func CopyValues(values []string) []string {
			if values == nil {
				return nil
			}
			result := make([]string, len(values))
			copy(result, values)
			return result
		}

		// Segments calculates the path segments for the given path.
		func Segments(path string) []string {
			for strings.HasPrefix(path, "/") {
				path = path[1:]
			}
			for strings.HasSuffix(path, "/") {
				path = path[0:len(path)-1]
			}
			return strings.Split(path, "/")
		}

		// Name of the header used to contain the metrics path:
		const metricHeader = "X-Metric"
        `)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *HelpersGenerator) helpersFile() string {
	return g.names.File(nomenclator.Helpers)
}
