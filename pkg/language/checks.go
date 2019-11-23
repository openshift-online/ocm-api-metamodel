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

// This file contains the semantics checks.

package language

import (
	"github.com/openshift-online/ocm-api-metamodel/pkg/concepts"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
)

func (r *Reader) checkModel() {
	for _, service := range r.model.Services() {
		r.checkService(service)
	}
}

func (r *Reader) checkService(service *concepts.Service) {
	for _, version := range service.Versions() {
		r.checkVersion(version)
	}
}

func (r *Reader) checkVersion(version *concepts.Version) {
	// Check that there is a root resource:
	if version.Root() == nil {
		r.reporter.Errorf("Version '%s' doesn't have a root resource", version)
	}

	// Check the resources:
	for _, resource := range version.Resources() {
		r.checkResource(resource)
	}
}

func (r *Reader) checkResource(resource *concepts.Resource) {
	for _, method := range resource.Methods() {
		r.checkMethod(method)
	}
}

func (r *Reader) checkMethod(method *concepts.Method) {
	// Run specific checks according tot he type of method:
	name := method.Name()
	switch {
	case name.Equals(nomenclator.Add):
		r.checkAdd(method)
	case name.Equals(nomenclator.Delete):
		r.checkDelete(method)
	case name.Equals(nomenclator.Get):
		r.checkGet(method)
	case name.Equals(nomenclator.List):
		r.checkList(method)
	case name.Equals(nomenclator.Post):
		r.checkPost(method)
	case name.Equals(nomenclator.Update):
		r.checkUpdate(method)
	default:
		r.checkAction(method)
	}

	// Check the parameters:
	for _, parameter := range method.Parameters() {
		r.checkParameter(parameter)
	}
}

func (r *Reader) checkAdd(method *concepts.Method) {
	// Only scalar and struct parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() && !parameter.Type().IsStruct() {
			r.reporter.Errorf(
				"Type of parameter '%s' should be scalar or struct but it is %s",
				parameter, parameter.Type().Kind(),
			)
		}
	}

	// Exactly one struct parameter:
	var structs []*concepts.Parameter
	for _, parameter := range method.Parameters() {
		if parameter.Type().IsStruct() {
			structs = append(structs, parameter)
		}
	}
	count := len(structs)
	if count != 1 {
		r.reporter.Errorf(
			"Method '%s' should have exactly one struct parameter but it has %d",
			method, count,
		)
	}

	// Non scalar parameters should be input and output:
	for _, parameter := range structs {
		if !parameter.In() || !parameter.Out() {
			r.reporter.Errorf(
				"Direction of parameter '%s' should be 'in out'",
				parameter,
			)
		}
	}

	// Struct parameters should be named `body`:
	for _, parameter := range structs {
		if !nomenclator.Body.Equals(parameter.Name()) {
			r.reporter.Errorf(
				"Name of parameter '%s' should be '%s'",
				parameter, nomenclator.Body,
			)
		}
	}
}

func (r *Reader) checkDelete(method *concepts.Method) {
	// Only scalar parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() {
			r.reporter.Errorf(
				"Type of parameter '%s' should be scalar but it is %s",
				method, parameter.Type().Kind(),
			)
		}
	}

	// Only input parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.In() || parameter.Out() {
			r.reporter.Errorf("Direction of parameter '%s' must be 'in'", parameter)
		}
	}
}

func (r *Reader) checkGet(method *concepts.Method) {
	// Only scalar and struct parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() && !parameter.Type().IsStruct() {
			r.reporter.Errorf(
				"Type of parameter '%s' must be scalar or struct but it is %s",
				parameter, parameter.Type().Kind(),
			)
		}
	}

	// Exactly one struct parameter:
	var structs []*concepts.Parameter
	for _, parameter := range method.Parameters() {
		if parameter.Type().IsStruct() {
			structs = append(structs, parameter)
		}
	}
	count := len(structs)
	if count != 1 {
		r.reporter.Errorf(
			"Method '%s' should have exactly one struct parameter but it has %d",
			method, count,
		)
	}

	// Scalar parameters should be input only:
	for _, parameter := range method.Parameters() {
		if parameter.Type().IsScalar() {
			if !parameter.In() || parameter.Out() {
				r.reporter.Errorf(
					"Direction of parameter '%s' should be 'in'",
					parameter,
				)
			}
		}
	}

	// Struct parameters should be output only:
	for _, parameter := range structs {
		if !parameter.Out() || parameter.In() {
			r.reporter.Errorf("Direction of parameter '%s' should be 'out'", parameter)
		}
	}

	// Struct parameters should be named `body`:
	for _, parameter := range structs {
		if !nomenclator.Body.Equals(parameter.Name()) {
			r.reporter.Errorf(
				"Name of parameter '%s' should be '%s'",
				parameter, nomenclator.Body,
			)
		}
	}
}

func (r *Reader) checkList(method *concepts.Method) {
	// Get the reference to the and to the version:
	resource := method.Owner()
	version := resource.Owner()

	// Only scalar and list parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() && !parameter.Type().IsList() {
			r.reporter.Errorf(
				"Type of parameter '%s' should be scalar or struct but it is '%s'",
				parameter, parameter.Type().Kind(),
			)
		}
	}

	// Exactly one list parameter:
	var lists []*concepts.Parameter
	for _, parameter := range method.Parameters() {
		if parameter.Type().IsList() {
			lists = append(lists, parameter)
		}
	}
	count := len(lists)
	if count != 1 {
		r.reporter.Errorf(
			"Method '%s' should have exactly one list parameter but it has %d",
			method, count,
		)
	}

	// Check the `page` parameter:
	page := method.GetParameter(nomenclator.Page)
	if page == nil {
		r.reporter.Warnf(
			"Method '%s' doesn't have a '%s' parameter",
			method, nomenclator.Page,
		)
	} else {
		if page.Type() != version.Integer() {
			r.reporter.Errorf(
				"Type of parameter '%s' should be integer but it is '%s'",
				page, page.Type(),
			)
		}
		if !page.In() || !page.Out() {
			r.reporter.Errorf(
				"Direction of parameter '%s' should be 'in out'",
				page,
			)
		}
		if page.Default() != 1 {
			r.reporter.Errorf(
				"Default value of parameter `%s` should be 1",
				page,
			)
		}
	}

	// Check the `size` parameter:
	size := method.GetParameter(nomenclator.Size)
	if size == nil {
		r.reporter.Warnf(
			"Method '%s' doesn't have a '%s' parameter",
			method, nomenclator.Size,
		)
	} else {
		if size.Type() != version.Integer() {
			r.reporter.Warnf(
				"Type of parameter '%s' should be integer but it is '%s'",
				size, size.Type(),
			)
		}
		if !size.In() || !size.Out() {
			r.reporter.Errorf(
				"Direction of parameter '%s' should be 'in out'",
				size,
			)
		}
		if size.Default() == nil {
			r.reporter.Errorf(
				"Parameter `%s` should have a default value",
				size,
			)
		}
	}

	// Check the `total` parameter:
	total := method.GetParameter(nomenclator.Total)
	if total == nil {
		r.reporter.Warnf(
			"Method '%s' doesn't have a '%s' parameter",
			method, nomenclator.Total,
		)
	} else {
		if total.Type() != version.Integer() {
			r.reporter.Warnf(
				"Type of parameter '%s' should be integer but it is '%s'",
				total, total.Type(),
			)
		}
		if total.In() || !total.Out() {
			r.reporter.Errorf(
				"Direction of parameter '%s' should be 'out'",
				total,
			)
		}
	}

	// Check the `items` parameter:
	items := method.GetParameter(nomenclator.Items)
	if items == nil {
		r.reporter.Errorf(
			"Method '%s' doesn't have a '%s' parameter",
			method, nomenclator.Items,
		)
	}

	// List parameters should be output only:
	for _, parameter := range lists {
		if parameter.In() || !parameter.Out() {
			r.reporter.Errorf(
				"Direction of parameter '%s' should be 'out'",
				parameter,
			)
		}
	}

	// List parameters should be named `items`:
	for _, parameter := range lists {
		if !nomenclator.Items.Equals(parameter.Name()) {
			r.reporter.Errorf(
				"Name of parameter '%s' should be '%s'",
				parameter, nomenclator.Items,
			)
		}
	}
}

func (r *Reader) checkPost(method *concepts.Method) {
	// Empty on purpose.
}

func (r *Reader) checkUpdate(method *concepts.Method) {
	// Only scalar and struct parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() && !parameter.Type().IsStruct() {
			r.reporter.Errorf(
				"Type of parameter '%s' should be scalar or struct but it is '%s'",
				parameter, parameter.Type().Kind,
			)
		}
	}

	// Exactly one struct parameter:
	var structs []*concepts.Parameter
	for _, parameter := range method.Parameters() {
		if parameter.Type().IsStruct() {
			structs = append(structs, parameter)
		}
	}
	count := len(structs)
	if count != 1 {
		r.reporter.Errorf(
			"Method '%s' should have exactly one struct parameter but it has %d",
			method, count,
		)
	}

	// Scalar parameters should be input only:
	for _, parameter := range method.Parameters() {
		if parameter.Type().IsScalar() {
			if !parameter.In() || parameter.Out() {
				r.reporter.Errorf(
					"Direction of parameter '%s' should be 'in'",
					parameter,
				)
			}
		}
	}

	// Struct parameters should be input and output:
	for _, parameter := range structs {
		if !parameter.In() || !parameter.Out() {
			r.reporter.Errorf(
				"Direction of parameter '%s' should be 'in out'",
				parameter,
			)
		}
	}

	// Struct parameters should be named `body`:
	for _, parameter := range structs {
		if !nomenclator.Body.Equals(parameter.Name()) {
			r.reporter.Errorf(
				"Name of parameter '%s' should be be '%s'",
				parameter, nomenclator.Body,
			)
		}
	}
}

func (r *Reader) checkAction(method *concepts.Method) {
	// Empty on purpose.
}

func (r *Reader) checkParameter(parameter *concepts.Parameter) {
	// Get the version:
	version := parameter.Owner().Owner().Owner()

	// Check that the default value is of a type compatible with the type of the parameter:
	value := parameter.Default()
	if value != nil {
		var typ *concepts.Type
		switch value.(type) {
		case bool:
			typ = version.Boolean()
		case int:
			typ = version.Integer()
		case string:
			typ = version.String()
		default:
			r.reporter.Errorf(
				"Don't know how to check if default value '%v' is compatible "+
					"with type of parameter '%s'",
				value, parameter,
			)
		}
		if typ != nil && typ != parameter.Type() {
			r.reporter.Errorf(
				"Type of default value of parameter '%s' should be '%s'",
				parameter, parameter.Type(),
			)
		}
	}
}
