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

package concepts

import (
	"sort"

	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
	"github.com/openshift-online/ocm-api-metamodel/pkg/nomenclator"
)

// Method represents a method of a resource.
type Method struct {
	owner      *Resource
	doc        string
	name       *names.Name
	parameters ParameterSlice
}

// NewMethod creates a new method.
func NewMethod() *Method {
	return new(Method)
}

// Owner returns the resource that owns this method.
func (m *Method) Owner() *Resource {
	return m.owner
}

// SetOwner sets the resource that owns this method.
func (m *Method) SetOwner(value *Resource) {
	m.owner = value
}

// Doc returns the documentation of this method.
func (m *Method) Doc() string {
	return m.doc
}

// SetDoc sets the documentation of this type.
func (m *Method) SetDoc(value string) {
	m.doc = value
}

// Name returns the name of the method.
func (m *Method) Name() *names.Name {
	return m.name
}

// SetName sets the name of the method.
func (m *Method) SetName(value *names.Name) {
	m.name = value
}

// Parameters returns the parameters of the method.
func (m *Method) Parameters() ParameterSlice {
	return m.parameters
}

// AddParameter adds a parameter to the method.
func (m *Method) AddParameter(parameter *Parameter) {
	if parameter != nil {
		m.parameters = append(m.parameters, parameter)
		sort.Sort(m.parameters)
		parameter.SetOwner(m)
	}
}

// Get parameter returns the parameter with the given name, or nil if there is no parameter with
// that name.
func (m *Method) GetParameter(name *names.Name) *Parameter {
	if name == nil {
		return nil
	}
	for _, parameter := range m.parameters {
		if parameter.Name().Equals(name) {
			return parameter
		}
	}
	return nil
}

// IsAction determined if this method is an action instead of a regular REST method.
func (m *Method) IsAction() bool {
	switch {
	case m.name.Equals(nomenclator.Add):
		return false
	case m.name.Equals(nomenclator.Delete):
		return false
	case m.name.Equals(nomenclator.Get):
		return false
	case m.name.Equals(nomenclator.List):
		return false
	case m.name.Equals(nomenclator.Post):
		return false
	case m.name.Equals(nomenclator.Update):
		return false
	default:
		return true
	}
}

// MethodSlice is used to simplify sorting of slices of methods by name.
type MethodSlice []*Method

func (s MethodSlice) Len() int {
	return len(s)
}

func (s MethodSlice) Less(i, j int) bool {
	return names.Compare(s[i].name, s[j].name) == -1
}

func (s MethodSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
