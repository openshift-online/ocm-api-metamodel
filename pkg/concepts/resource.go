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

// Resource represents an API resource.
type Resource struct {
	owner    *Version
	doc      string
	name     *names.Name
	methods  []*Method
	locators []*Locator
}

// NewResource creates a new resource.
func NewResource() *Resource {
	return new(Resource)
}

// Owner returns the version that owns the resource.
func (r *Resource) Owner() *Version {
	return r.owner
}

// SetOwner sets the version that owns the resource.
func (r *Resource) SetOwner(value *Version) {
	r.owner = value
}

// Doc returns the documentation of this resource.
func (r *Resource) Doc() string {
	return r.doc
}

// SetDoc sets the documentation of this resource.
func (r *Resource) SetDoc(value string) {
	r.doc = value
}

// Name returns the name of the resource.
func (r *Resource) Name() *names.Name {
	return r.name
}

// SetName sets the name of the resource.
func (r *Resource) SetName(value *names.Name) {
	r.name = value
}

// Methods returns the methods of the resource.
func (r *Resource) Methods() []*Method {
	return r.methods
}

// AddMethod adds a method to the resource.
func (r *Resource) AddMethod(method *Method) {
	if method != nil {
		r.methods = append(r.methods, method)
		method.SetOwner(r)
	}
}

// Locators returns the locators of the resource.
func (r *Resource) Locators() []*Locator {
	return r.locators
}

// AddLocator adds a locator to the resource.
func (r *Resource) AddLocator(locator *Locator) {
	if locator != nil {
		r.locators = append(r.locators, locator)
		locator.SetOwner(r)
	}
}
