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
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// ErrorsGeneratorBuilder is an object used to configure and build an errors generator. Don't create
// instances directly, use the NewErrorsGenerator function instead.
type ErrorsGeneratorBuilder struct {
	reporter *reporter.Reporter
	model    *concepts.Model
	output   string
	packages *PackagesCalculator
	names    *NamesCalculator
}

// ErrorsGenerator generates errors code. Don't create instances directly, use the builder instead.
type ErrorsGenerator struct {
	reporter *reporter.Reporter
	errors   int
	model    *concepts.Model
	output   string
	packages *PackagesCalculator
	names    *NamesCalculator
	buffer   *Buffer
}

// NewErrorsGenerator creates a new builder for errors generators.
func NewErrorsGenerator() *ErrorsGeneratorBuilder {
	return new(ErrorsGeneratorBuilder)
}

// Reporter sets the object that will be used to report information about the generation process.
func (b *ErrorsGeneratorBuilder) Reporter(value *reporter.Reporter) *ErrorsGeneratorBuilder {
	b.reporter = value
	return b
}

// Model sets the model that will be used by the client generator.
func (b *ErrorsGeneratorBuilder) Model(value *concepts.Model) *ErrorsGeneratorBuilder {
	b.model = value
	return b
}

// Output sets the output directory.
func (b *ErrorsGeneratorBuilder) Output(value string) *ErrorsGeneratorBuilder {
	b.output = value
	return b
}

// Packages sets the object that will be used to calculate package names.
func (b *ErrorsGeneratorBuilder) Packages(
	value *PackagesCalculator) *ErrorsGeneratorBuilder {
	b.packages = value
	return b
}

// Names sets the object that will be used to calculate names.
func (b *ErrorsGeneratorBuilder) Names(value *NamesCalculator) *ErrorsGeneratorBuilder {
	b.names = value
	return b
}

// Build checks the configuration stored in the builder and, if it is correct, creates a new errors
// generator using it.
func (b *ErrorsGeneratorBuilder) Build() (generator *ErrorsGenerator, err error) {
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
	if b.packages == nil {
		err = fmt.Errorf("packages calculator is mandatory")
		return
	}
	if b.names == nil {
		err = fmt.Errorf("names calculator is mandatory")
		return
	}

	// Create the generator:
	generator = &ErrorsGenerator{
		reporter: b.reporter,
		model:    b.model,
		output:   b.output,
		packages: b.packages,
		names:    b.names,
	}

	return
}

// Run executes the code generator.
func (g *ErrorsGenerator) Run() error {
	var err error

	// Generate the errors:
	err = g.generateCommonErrors()
	if err != nil {
		return err
	}

	// Generate the clients for the versions:
	for _, service := range g.model.Services() {
		for _, version := range service.Versions() {
			err = g.generateVersionErrors(version)
			if err != nil {
				return err
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

func (g *ErrorsGenerator) generateCommonErrors() error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.ErrorsPackage()
	fileName := g.errorsFile()

	// Create the buffer for the generated code:
	g.buffer, err = NewBuffer().
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
	g.buffer.Import("fmt", "")
	g.buffer.Import("io", "")
	g.buffer.Import("strings", "")
	g.buffer.Import("github.com/golang/glog", "")
	g.buffer.Import("github.com/openshift-online/ocm-api-metamodel/pkg/runtime", "")
	g.buffer.Import(g.packages.HelpersImport(), "")
	g.buffer.Emit(`
		// Error kind is the name of the type used to represent errors.
		const ErrorKind = "Error"

		// ErrorNilKind is the name of the type used to nil errors.
		const ErrorNilKind = "ErrorNil"

		// ErrorBuilder is a builder for the error type.
		type ErrorBuilder struct{
			bitmap_     uint32
			id          string
			href        string
			code        string
			reason      string
			operationID string
		}

		// Error represents errors.
		type Error struct {
			bitmap_     uint32
			id          string
			href        string
			code        string
			reason      string
			operationID string
		}

		// NewError creates a new builder that can then be used to create error objects.
		func NewError() *ErrorBuilder {
			return &ErrorBuilder{}
		}

		// ID sets the identifier of the error.
		func (b *ErrorBuilder) ID(value string) *ErrorBuilder {
			b.id = value
			b.bitmap_ |= 1
			return b
		}

		// HREF sets the link of the error.
		func (b *ErrorBuilder) HREF(value string) *ErrorBuilder {
			b.href = value
			b.bitmap_ |= 2
			return b
		}

		// Code sets the code of the error.
		func (b *ErrorBuilder) Code(value string) *ErrorBuilder {
			b.code = value
			b.bitmap_ |= 4
			return b
		}

		// Reason sets the reason of the error.
		func (b *ErrorBuilder) Reason(value string) *ErrorBuilder {
			b.reason = value
			b.bitmap_ |= 8
			return b
		}

		// OperationID sets the identifier of the operation that caused the error.
		func (b *ErrorBuilder) OperationID(value string) *ErrorBuilder {
			b.operationID = value
			b.bitmap_ |= 16
			return b
		}

		// Build uses the information stored in the builder to create a new error object.
		func (b *ErrorBuilder) Build() (result *Error,  err error) {
			result = &Error{
				id:          b.id,
				href:        b.href,
				code:        b.code,
				reason:      b.reason,
				operationID: b.operationID,
				bitmap_:   b.bitmap_,
			}
			return
		}

		// Kind returns the name of the type of the error.
		func (e *Error) Kind() string {
			if e == nil {
				return ErrorNilKind
			}
			return ErrorKind
		}

		// ID returns the identifier of the error.
		func (e *Error) ID() string {
			if e != nil && e.bitmap_&1 != 0 {
				return e.id
			}
			return ""
		}

		// GetID returns the identifier of the error and a flag indicating if the
		// identifier has a value.
		func (e *Error) GetID() (value string, ok bool) {
			ok = e != nil && e.bitmap_&1 != 0
			if ok {
				value = e.id
			}
			return
		}

		// HREF returns the link to the error.
		func (e *Error) HREF() string {
			if e != nil && e.bitmap_&2 != 0 {
				return e.href
			}
			return ""
		}

		// GetHREF returns the link of the error and a flag indicating if the
		// link has a value.
		func (e *Error) GetHREF() (value string, ok bool) {
			ok = e != nil && e.bitmap_&2 != 0
			if ok {
				value = e.href
			}
			return
		}

		// Code returns the code of the error.
		func (e *Error) Code() string {
			if e != nil && e.bitmap_&4 != 0 {
				return e.code
			}
			return ""
		}

		// GetCode returns the link of the error and a flag indicating if the
		// code has a value.
		func (e *Error) GetCode() (value string, ok bool) {
			ok = e != nil && e.bitmap_&4 != 0
			if ok {
				value = e.code
			}
			return
		}

		// Reason returns the reason of the error.
		func (e *Error) Reason() string {
			if e != nil && e.bitmap_&8 != 0 {
				return e.reason
			}
			return ""
		}

		// GetReason returns the link of the error and a flag indicating if the
		// reason has a value.
		func (e *Error) GetReason() (value string, ok bool) {
			ok = e != nil && e.bitmap_&8 != 0
			if ok {
				value = e.reason
			}
			return
		}

		// OperationID returns the identifier of the operation that caused the error.
		func (e *Error) OperationID() string {
			if e != nil && e.bitmap_&16 != 0 {
				return e.operationID
			}
			return ""
		}

		// GetOperationID returns the identifier of the operation that caused the error and
		// a flag indicating if that identifier does have a value.
		func (e *Error) GetOperationID() (value string, ok bool) {
			ok = e != nil && e.bitmap_&16 != 0
			if ok {
				value = e.operationID
			}
			return
		}

		// Error is the implementation of the error interface.
		func (e *Error) Error() string {
			chunks := make([]string, 0, 3)
			if e.id != "" {
				chunks = append(chunks, fmt.Sprintf("identifier is '%s'", e.id))
			}
			if e.code != "" {
				chunks = append(chunks, fmt.Sprintf("code is '%s'", e.code))
			}
			if e.operationID != "" {
				chunks = append(chunks, fmt.Sprintf("operation identifier is '%s'", e.operationID))
			}
			var result string
			size := len(chunks)
			if size == 1 {
				result = chunks[0]
			} else if size > 1 {
				result = strings.Join(chunks[0:size-1], ", ") + " and " + chunks[size-1]
			}
			if e.reason != "" {
				if result != "" {
					result = result + ": "
				}
				result = result + e.reason
			}
			if result == "" {
				result = "unknown error"
			}
			return result
		}

		// String returns a string representing the error.
		func (e *Error) String() string {
			return e.Error()
		}

		// UnmarshalError reads an error from the given source which can be an slice of
		// bytes, a string, a reader or a JSON decoder.
		func UnmarshalError(source interface{}) (object *Error, err error) {
			iterator, err := helpers.NewIterator(source)
			if err != nil {
				return
			}
			object = readError(iterator)
			err = iterator.Error
			return
		}

		func readError(iterator *jsoniter.Iterator) *Error {
			object := &Error{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "id":
					object.id = iterator.ReadString()
					object.bitmap_ |= 1
				case "href":
					object.href = iterator.ReadString()
					object.bitmap_ |= 2
				case "code":
					object.code = iterator.ReadString()
					object.bitmap_ |= 4
				case "reason":
					object.reason = iterator.ReadString()
					object.bitmap_ |= 8
				case "operation_id":
					object.operationID = iterator.ReadString()
					object.bitmap_ |= 16
				default:
					iterator.ReadAny()
				}
			}
			return object
		}

		// MarshalError writes an error to the given writer.
		func MarshalError(e *Error, writer io.Writer) error {
			stream := helpers.NewStream(writer)
			writeError(e, stream)
			stream.Flush()
			return stream.Error
		}

		func writeError(e *Error, stream *jsoniter.Stream) {
			stream.WriteObjectStart()
			stream.WriteObjectField("kind")
			stream.WriteString(ErrorKind)
			if e.bitmap_&1 != 0 {
				stream.WriteMore()
				stream.WriteObjectField("id")
				stream.WriteString(e.id)
			}
			if e.bitmap_&2 != 0 {
				stream.WriteMore()
				stream.WriteObjectField("href")
				stream.WriteString(e.href)
			}
			if e.bitmap_&4 != 0 {
				stream.WriteMore()
				stream.WriteObjectField("code")
				stream.WriteString(e.code)
			}
			if e.bitmap_&8 != 0 {
				stream.WriteMore()
				stream.WriteObjectField("reason")
				stream.WriteString(e.reason)
			}
			if e.bitmap_&16 != 0 {
				stream.WriteMore()
				stream.WriteObjectField("operation_id")
				stream.WriteString(e.operationID)
			}
			stream.WriteObjectEnd()
		}

		var panicID = "1000"
		var panicError, _ = NewError().
			ID(panicID).
			Reason("An unexpected error happened, please check the log of the service " +
			"for details").
			Build()

		// SendError writes a given error and status code to a response writer.
		// if an error occurred it will log the error and exit.
		// This methods is used internaly and no backwards compatibily is guaranteed.
		func SendError(w http.ResponseWriter, r *http.Request, object *Error) {
			status, err := strconv.Atoi(object.ID())
			if err != nil {
				SendPanic(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			err = MarshalError(object, w)
			if err != nil {
				glog.Errorf("Can't send response body for request '%s'", r.URL.Path)
				return
			}
		}

		// SendPanic sends a panic error response to the client, but it doesn't end the process.
		// This methods is used internaly and no backwards compatibily is guaranteed.
		func SendPanic(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			err := MarshalError(panicError, w)
			if err != nil {
				glog.Errorf(
					"Can't send panic response for request '%s': %s",
					r.URL.Path,
					err.Error(),
				)
			}
		}

		// SendNotFound sends a generic 404 error.
		func SendNotFound(w http.ResponseWriter, r *http.Request) {
			reason := fmt.Sprintf(
				"Can't find resource for path '%s'",
				r.URL.Path,
			)
			body, err := NewError().
				ID("404").
				Reason(reason).
				Build()
			if err != nil {
				SendPanic(w, r)
				return
			}
			SendError(w, r, body)
		}

		// SendMethodNotAllowed sends a generic 405 error.
		func SendMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
			reason := fmt.Sprintf(
				"Method '%s' isn't supported for path '%s'",
				r.Method, r.URL.Path,
			)
			body, err := NewError().
				ID("405").
				Reason(reason).
				Build()
			if err != nil {
				SendPanic(w, r)
				return
			}
			SendError(w, r, body)
		}

		// SendInternalServerError sends a generic 500 error.
		func SendInternalServerError(w http.ResponseWriter, r *http.Request) {
			reason := fmt.Sprintf(
				"Can't process '%s' request for path '%s' due to an internal"+
					"server error",
				r.Method, r.URL.Path,
			)
			body, err := NewError().
				ID("500").
				Reason(reason).
				Build()
			if err != nil {
				SendPanic(w, r)
				return
			}
			SendError(w, r, body)
		}
        `)

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ErrorsGenerator) generateVersionErrors(version *concepts.Version) error {
	var err error

	// Calculate the package and file name:
	pkgName := g.packages.VersionPackage(version)
	fileName := g.errorsFile()

	// Create the buffer for the generated code:
	g.buffer, err = NewBuffer().
		Reporter(g.reporter).
		Output(g.output).
		Packages(g.packages).
		Package(pkgName).
		File(fileName).
		Function("errorName", g.errorName).
		Build()
	if err != nil {
		return err
	}

	// Generate the code:
	err = g.generateVersionErrorsSource(version)
	if err != nil {
		return err
	}

	// Write the generated code:
	return g.buffer.Write()
}

func (g *ErrorsGenerator) generateVersionErrorsSource(version *concepts.Version) error {
	g.buffer.Emit(`
		{{ if .Version.Errors }}
			const (
				{{ range .Version.Errors }}
					{{ lineComment .Doc }}
					{{ errorName . }} = {{ .Code }}
				{{ end }}
			)
		{{ end }}
		`,
		"Version", version,
	)
	return nil
}

func (g *ErrorsGenerator) errorsFile() string {
	return g.names.File(nomenclator.Errors)
}

func (g *ErrorsGenerator) errorName(err *concepts.Error) string {
	return g.names.Public(names.Cat(err.Name(), nomenclator.Error))
}
