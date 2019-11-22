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

package reporter

import (
	"errors"
	"fmt"
	"os"

	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
)

// Reporter is the reported used by the metamodel tools. It prints the messages to the standard
// output stream. Don't create instances directly, use the NewReporter function instead.
type Reporter struct {
	errors int
}

// NewReporter createsa a new reporter.
func NewReporter() *Reporter {
	reporter := new(Reporter)
	return reporter
}

// Infof prints an informative message with the given format and arguments.
func (r *Reporter) Infof(format string, args ...interface{}) {
	message := r.printf(format, args)
	fmt.Fprintf(os.Stdout, "%s%s\n", infoPrefix, message)
}

// Warnf prints an warning message with the given format and arguments.
func (r *Reporter) Warnf(format string, args ...interface{}) {
	message := r.printf(format, args)
	fmt.Fprintf(os.Stdout, "%s%s\n", warnPrefix, message)
}

// Errorf prints an error message with the given format and arguments. It also return an error
// containing the same information, which will be usually discarded, except when the caller needs to
// report the error and also return it.
func (r *Reporter) Errorf(format string, args ...interface{}) error {
	message := r.printf(format, args)
	fmt.Fprintf(os.Stdout, "%s%s\n", errorPrefix, message)
	r.errors++
	return errors.New(message)
}

func (r *Reporter) printf(format string, args []interface{}) string {
	// Replace arguments that are named objects or names with their camel case equivalent,
	// as that is what users type in the source files:
	for i, arg := range args {
		switch typed := arg.(type) {
		case names.Named:
			name := typed.Name()
			args[i] = name.Camel()
		case *names.Name:
			args[i] = typed.Camel()
		}
	}

	// Format the message:
	return fmt.Sprintf(format, args...)
}

// Errors returns the number of errors that have been reported via this reporter.
func (r *Reporter) Errors() int {
	return r.errors
}

// Message prefix using ANSI scape seequences to set colors:
const (
	infoPrefix  = "\033[0;32mI:\033[m "
	warnPrefix  = "\033[0;33mW:\033[m "
	errorPrefix = "\033[0;31mE:\033[m "
)
