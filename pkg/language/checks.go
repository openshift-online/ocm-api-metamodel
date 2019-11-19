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

// checkMethod performs method level semantics checks.
func (r *Reader) checkMethod(method *concepts.Method) {
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
}

func (r *Reader) checkAdd(method *concepts.Method) {
	// Get the reference the resource:
	resource := method.Owner()

	// Only scalar and struct parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() && !parameter.Type().IsStruct() {
			r.reporter.Errorf(
				"Method '%s' of resource '%s' can only have parameters that are "+
					"scalars or structs, but type of parameter '%s' is %s",
				method, resource, parameter, parameter.Type().Kind(),
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
			"Method '%s' of resource '%s' should have exactly one struct parameter "+
				"but it has %d",
			method, resource, count,
		)
	}

	// Non scalar parameters should be input and output:
	for _, parameter := range structs {
		if !parameter.In() || !parameter.Out() {
			r.reporter.Errorf(
				"Parameter '%s' of method '%s' of resource '%s' should be 'in out'",
				parameter, method, resource,
			)
		}
	}

	// Struct parameters should be named `body`:
	for _, parameter := range structs {
		if !nomenclator.Body.Equals(parameter.Name()) {
			r.reporter.Errorf(
				"Name of parameter '%s' of method '%s' of resource '%s' should "+
					"be '%s'",
				parameter, method, resource, nomenclator.Body,
			)
		}
	}
}

func (r *Reader) checkDelete(method *concepts.Method) {
	// Get the reference the resource:
	resource := method.Owner()

	// Only scalar parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() {
			r.reporter.Errorf(
				"Method '%s' of resource '%s' can only have scalar parameters, "+
					"parameter '%s' is %s",
				method, resource, parameter, parameter.Type().Kind(),
			)
		}
	}

	// Only input parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.In() {
			r.reporter.Errorf(
				"Method '%s' of resource '%s' can only have 'in' parameters, but "+
					"parameter '%s' isn't",
				method, resource, parameter,
			)
		}
		if parameter.Out() {
			r.reporter.Errorf(
				"Method '%s' of resource '%s' can only have 'in' parameters, but "+
					"parameter '%s' is 'out'",
				method, resource, parameter,
			)
		}
	}
}

func (r *Reader) checkGet(method *concepts.Method) {
	// Get the reference to the resource:
	resource := method.Owner()

	// Only scalar and struct parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() && !parameter.Type().IsStruct() {
			r.reporter.Errorf(
				"Method '%s' of resource '%s' can only have parameters that are "+
					"scalars or structs, but type of parameter '%s' is %s",
				method, resource, parameter,
				parameter.Type().Kind(),
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
			"Method '%s' of resource '%s' should have exactly one struct parameter "+
				"but it has %d",
			method, resource, count,
		)
	}

	// Scalar parameters should be input only:
	for _, parameter := range method.Parameters() {
		if parameter.Type().IsScalar() {
			if !parameter.In() {
				r.reporter.Errorf(
					"Scalar parameters of method '%s' of resource '%s' " +
						"should be 'in' but parameter '%s' isn't",
					method, resource, parameter,
				)
			}
			if parameter.Out() {
				r.reporter.Errorf(
					"Scalar parameters of method '%s' of resource '%s" +
						"should be 'in' but parameter '%s' is 'out'",
					method, resource, parameter,
				)
			}
		}
	}

	// Struct parameters should be output only:
	for _, parameter := range structs {
		if !parameter.Out() {
			r.reporter.Errorf(
				"Non scalar parameters of method '%s' of resource '%s' should "+
					"be 'out' but parameter '%s' isn't",
				method, resource, parameter,
			)
		}
		if parameter.In() {
			r.reporter.Errorf(
				"Non scalar parameters of method '%s' of resource '%s' should "+
					"be 'out' but parameter '%s' is 'in'",
				method, resource, parameter,
			)
		}
	}

	// Struct parameters should be named `body`:
	for _, parameter := range structs {
		if !nomenclator.Body.Equals(parameter.Name()) {
			r.reporter.Errorf(
				"Name of parameter '%s' of method '%s' of resource '%s' should "+
					"be '%s'",
				parameter, method, resource,
				nomenclator.Body,
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
				"Method '%s' of resource '%s' can only have parameters that are "+
					"scalars or lists, but type of parameter '%s' is %s",
				method, resource, parameter,
				parameter.Type().Kind(),
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
			"Method '%s' of resource '%s' should have exactly one list parameter "+
				"but it has %d",
			method, resource, count,
		)
	}

	// Check the `page` parameter:
	page := method.GetParameter(nomenclator.Page)
	if page == nil {
		r.reporter.Warnf(
			"Method '%s' of resource '%s' doesn't have a '%s' parameter",
			method, resource, nomenclator.Page,
		)
	} else {
		if page.Type() != version.Integer() {
			r.reporter.Errorf(
				"Parameter '%s' of method '%s' of resource '%s' doesn't have "+
					"should be integer but it is %s",
				page, method, resource, page.Type(),
			)
		}
		if !page.In() || !page.Out() {
			r.reporter.Errorf(
				"Parameter '%s' of method '%s' of resource '%s' should be 'in out'",
				page, method, resource,
			)
		}
	}

	// Check the `size` parameter:
	size := method.GetParameter(nomenclator.Size)
	if size == nil {
		r.reporter.Warnf(
			"Method '%s' of resource '%s' doesn't have a '%s' parameter",
			method, resource, nomenclator.Size,
		)
	} else {
		if size.Type() != version.Integer() {
			r.reporter.Warnf(
				"Parameter '%s' of method '%s' of resource '%s' doesn't have "+
					"should be integer but it is %s",
				size, method, resource, size.Type(),
			)
		}
		if !size.In() || !size.Out() {
			r.reporter.Errorf(
				"Parameter '%s' of method '%s' of resource '%s' should be 'in out'",
				size, method, resource,
			)
		}
	}

	// Check the `total` parameter:
	total := method.GetParameter(nomenclator.Total)
	if total == nil {
		r.reporter.Warnf(
			"Method '%s' of resource '%s' doesn't have a '%s' parameter",
			method, resource, nomenclator.Total,
		)
	} else {
		if total.Type() != version.Integer() {
			r.reporter.Warnf(
				"Parameter '%s' of method '%s' of resource '%s' doesn't have "+
					"should be integer but it is %s",
				total, method, resource, total.Type(),
			)
		}
		if total.In() || !total.Out() {
			r.reporter.Errorf(
				"Parameter '%s' of method '%s' of resource '%s' should be 'out'",
				total, method, resource,
			)
		}
	}

	// Check the `items` parameter:
	items := method.GetParameter(nomenclator.Items)
	if items == nil {
		r.reporter.Errorf(
			"Method '%s' of resource '%s' doesn't have a '%s' parameter",
			method, resource, nomenclator.Items,
		)
	}

	// List parameters should be output only:
	for _, parameter := range lists {
		if parameter.In() || !parameter.Out() {
			r.reporter.Errorf(
				"Parameters '%s' of method '%s' of resource '%s' should be 'out'",
				parameter, method, resource.Name(),
			)
		}
	}

	// List parameters should be named `items`:
	for _, parameter := range lists {
		if !nomenclator.Items.Equals(parameter.Name()) {
			r.reporter.Errorf(
				"Name of parameter '%s' of method '%s' of resource '%s' should "+
					"be '%s'",
				parameter, method, resource, nomenclator.Items,
			)
		}
	}
}

func (r *Reader) checkPost(method *concepts.Method) {
	// Empty on purpose.
}

func (r *Reader) checkUpdate(method *concepts.Method) {
	// Get the reference to the resource:
	resource := method.Owner()

	// Only scalar and struct parameters:
	for _, parameter := range method.Parameters() {
		if !parameter.Type().IsScalar() && !parameter.Type().IsStruct() {
			r.reporter.Errorf(
				"Method '%s' of resource '%s' can only have parameters that are "+
					"scalars or structs, but type of parameter '%s' is %s",
				method.Name(), resource.Name(), parameter.Name(),
				parameter.Type().Kind,
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
			"Method '%s' of resource '%s' should have exactly one struct parameter "+
				"but it has %d",
			method.Name(), resource.Name(), count,
		)
	}

	// Scalar parameters should be input only:
	for _, parameter := range method.Parameters() {
		if parameter.Type().IsScalar() {
			if !parameter.In() {
				r.reporter.Errorf(
					"Scalar parameters of method '%s' of resource '%s' " +
						"should be 'in' but parameter '%s' isn't",
					method.Name(), resource.Name(), parameter.Name(),
				)
			}
			if parameter.Out() {
				r.reporter.Errorf(
					"Scalar parameters of method '%s' of resource '%s" +
						"should be 'in' but parameter '%s' is 'out'",
					method.Name(), resource.Name(), parameter.Name(),
				)
			}
		}
	}

	// Struct parameters should be input and output:
	for _, parameter := range structs {
		if !parameter.In() || !parameter.Out() {
			r.reporter.Errorf(
				"Parameter '%s' of method '%s' of resource '%s' should be 'in out'",
				parameter.Name(), method.Name(), resource.Name(),
			)
		}
	}

	// Struct parameters should be named `body`:
	for _, parameter := range structs {
		if !nomenclator.Body.Equals(parameter.Name()) {
			r.reporter.Errorf(
				"Name of parameter '%s' of method '%s' of resource '%s' should "+
					"be '%s'",
				parameter.Name(), method.Name(), resource.Name(),
				nomenclator.Body,
			)
		}
	}
}

func (r *Reader) checkAction(method *concepts.Method) {
	// Empty on purpose.
}
