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
	"github.com/openshift-online/ocm-api-metamodel/pkg/names"
)

// Parameter represents a parameter of a method.
type Parameter struct {
	owner *Method
	doc   string
	name  *names.Name
	typ   *Type
	in    bool
	out   bool
}

// NewParameter creates a new parameter.
func NewParameter() *Parameter {
	return new(Parameter)
}

// Owner returns the method that owns this parameter.
func (p *Parameter) Owner() *Method {
	return p.owner
}

// SetOwner sets the method that owns this parameter.
func (p *Parameter) SetOwner(value *Method) {
	p.owner = value
}

// Doc returns the documentation of this parameter.
func (p *Parameter) Doc() string {
	return p.doc
}

// SetDoc sets the documentation of this parameter.
func (p *Parameter) SetDoc(value string) {
	p.doc = value
}

// Name returns the name of the parameter.
func (p *Parameter) Name() *names.Name {
	return p.name
}

// SetName sets the name of the parameter.
func (p *Parameter) SetName(value *names.Name) {
	p.name = value
}

// Type returns the type of the parameter.
func (p *Parameter) Type() *Type {
	return p.typ
}

// SetType sets the type of the parameter.
func (p *Parameter) SetType(value *Type) {
	p.typ = value
}

// In returns true if this is an input parameter.
func (p *Parameter) In() bool {
	return p.in
}

// SetIn sets the input direction flag.
func (p *Parameter) SetIn(value bool) {
	p.in = value
}

// Out returns true if this is an output parameter.
func (p *Parameter) Out() bool {
	return p.out
}

// SetOut sets the output direction flag.
func (p *Parameter) SetOut(value bool) {
	p.out = value
}

// ParameterSlice is used to simplify sorting of slices of attributes by name.
type ParameterSlice []*Parameter

func (s ParameterSlice) Len() int {
	return len(s)
}

func (s ParameterSlice) Less(i, j int) bool {
	return names.Compare(s[i].name, s[j].name) == -1
}

func (s ParameterSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
