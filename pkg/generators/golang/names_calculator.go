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
	"strings"

	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// NamesCalculatorBuilder is an object used to configure and build the Go names calculators. Don't
// create instances directly, use the NewNamesCalculator function instead.
type NamesCalculatorBuilder struct {
	reporter *reporter.Reporter
}

// NamesCalculator is an object used to calculate Go names. Don't create instances directly, use the
// builder instead.
type NamesCalculator struct {
	reporter *reporter.Reporter
}

// NewNamesCalculator creates a Go names calculator builder.
func NewNamesCalculator() *NamesCalculatorBuilder {
	return &NamesCalculatorBuilder{}
}

// Reporter sets the object that will be used to report information about the calculation processes,
// including errors.
func (b *NamesCalculatorBuilder) Reporter(value *reporter.Reporter) *NamesCalculatorBuilder {
	b.reporter = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new
// calculator using it.
func (b *NamesCalculatorBuilder) Build() (calculator *NamesCalculator, err error) {
	// Check that the mandatory parameters have been provided:
	if b.reporter == nil {
		err = fmt.Errorf("reporter is mandatory")
		return
	}

	// Create the calculator:
	calculator = &NamesCalculator{
		reporter: b.reporter,
	}

	return
}

// Public concatenates the given list of values and converts them into a string that follows the
// conventions for Go public names. The elements of the list can be strings or objects that
// implement the names.Named interface. Strings will be used directly. Objects that implement the
// names.Named interface will be replaced by the result of calling the Name() method. Objects that
// have the @go annotation will be replaced by the value of the name paratemer of that annotation.
func (c *NamesCalculator) Public(values ...interface{}) string {
	result := c.concatenateValues(values)
	result = names.ToUpperCamel(result)
	result = AvoidReservedWord(result)
	return result
}

// Private concatenates the given list of values and converts them into a string that follows the
// conventions for Go private names. The elements of the list can be strings or objects that
// implement the names.Named interface. Strings will be used directly. Objects that implement the
// names.Named interface will be replaced by the result of calling the Name() method. Objects that
// have the @go annotation will be replaced by the value of the name paratemer of that annotation.
func (c *NamesCalculator) Private(values ...interface{}) string {
	result := c.concatenateValues(values)
	result = names.ToLowerCamel(result)
	result = AvoidReservedWord(result)
	return result
}

func (c *NamesCalculator) concatenateValues(values []interface{}) string {
	texts := make([]string, len(values))
	for i, value := range values {
		texts[i] = c.replaceValue(value)
	}
	return strings.Join(texts, "")
}

func (c *NamesCalculator) replaceValue(value interface{}) string {
	var result string
	switch typed := value.(type) {
	case string:
		result = typed
	case names.Named:
		result = typed.Name()
		annotated, ok := typed.(concepts.Annotated)
		if ok {
			annotation := annotated.GetAnnotation("go")
			if annotation != nil {
				name := annotation.GetString("name")
				if name != "" {
					result = name
				}
			}
		}
	default:
		c.reporter.Errorf(
			"Don't know how to concatenate object of type '%T'",
			value,
		)
		result = "!"
	}
	return result
}

// File converts the given name into an string, following the rules for Go source files.
func (c *NamesCalculator) File(name string) string {
	return names.ToSnake(name)
}
