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
	"gitlab.cee.redhat.com/service/ocm-api-metamodel/pkg/names"
)

// Attribute is the representation of an attribute of an structured type.
type Attribute struct {
	owner *Type
	doc   string
	name  *names.Name
	link  bool
	typ   *Type
}

// NewAttribute creates a new attribute.
func NewAttribute() *Attribute {
	return new(Attribute)
}

// Owner returns the type that owns this attribute.
func (a *Attribute) Owner() *Type {
	return a.owner
}

// SetOwner sets the type that owns this attribute.
func (a *Attribute) SetOwner(value *Type) {
	a.owner = value
}

// Doc returns the documentation of this attribute.
func (a *Attribute) Doc() string {
	return a.doc
}

// SetDoc sets the documentation of this attribute.
func (a *Attribute) SetDoc(value string) {
	a.doc = value
}

// Name returns the name of the attribute.
func (a *Attribute) Name() *names.Name {
	return a.name
}

// SetName sets the name of the attribute.
func (a *Attribute) SetName(value *names.Name) {
	a.name = value
}

// Link returns true if the attribute is a link, false otherwise.
func (a *Attribute) Link() bool {
	return a.link
}

// SetLink sets the flag that indicates if this attribute is a link.
func (a *Attribute) SetLink(value bool) {
	a.link = value
}

// Type returns the type of the attribute.
func (a *Attribute) Type() *Type {
	return a.typ
}

// SetType sets the type of the attribute.
func (a *Attribute) SetType(value *Type) {
	a.typ = value
}
