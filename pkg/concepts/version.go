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
	"gitlab.cee.redhat.com/service/ocm-api-metamodel/pkg/nomenclator"
)

// Version is the representation of a version of a service.
type Version struct {
	// Service that owns this version:
	owner *Service

	// Name of the version:
	name *names.Name

	// All the types of the version, indexed by name:
	types map[string]*Type

	// All the resources of the version, indexed by name:
	resources map[string]*Resource

	// All the error catagories of the version, indexed by name:
	errors map[string]*Error
}

// NewVersion creates a new version containing only the built-in types.
func NewVersion() *Version {
	// Create an empty version:
	version := new(Version)
	version.types = make(map[string]*Type)
	version.resources = make(map[string]*Resource)
	version.errors = make(map[string]*Error)

	// Add the built-in scalar types:
	version.addScalarType(nomenclator.Boolean)
	version.addScalarType(nomenclator.Integer)
	version.addScalarType(nomenclator.Float)
	version.addScalarType(nomenclator.String)
	version.addScalarType(nomenclator.Date)

	return version
}

// Owner returns the service that owns this version.
func (v *Version) Owner() *Service {
	return v.owner
}

// SetOwner sets the service that owns this version.
func (v *Version) SetOwner(value *Service) {
	v.owner = value
}

// Name returns the name of this version.
func (v *Version) Name() *names.Name {
	return v.name
}

// SetName sets the name of this version.
func (v *Version) SetName(value *names.Name) {
	v.name = value
}

// Types returns the list of types that are part of this version.
func (v *Version) Types() []*Type {
	count := len(v.types)
	types := make([]*Type, count)
	index := 0
	for _, typ := range v.types {
		types[index] = typ
		index++
	}
	return types
}

// FindType returns the type with the given name, or nil of there is no such type.
func (v *Version) FindType(name *names.Name) *Type {
	if name == nil {
		return nil
	}
	return v.types[name.String()]
}

// AddType adds the given type to the version.
func (v *Version) AddType(typ *Type) {
	if typ != nil {
		v.types[typ.Name().String()] = typ
		typ.SetOwner(v)
	}
}

// AddTypes adds the given types to the version.
func (v *Version) AddTypes(types []*Type) {
	for _, typ := range types {
		v.AddType(typ)
	}
}

// Boolean returns the boolean type.
func (v *Version) Boolean() *Type {
	return v.FindType(nomenclator.Boolean)
}

// Integer returns the integer type.
func (v *Version) Integer() *Type {
	return v.FindType(nomenclator.Integer)
}

// String returns the string type.
func (v *Version) String() *Type {
	return v.FindType(nomenclator.String)
}

// Float returns the floating point type.
func (v *Version) Float() *Type {
	return v.FindType(nomenclator.Float)
}

// Date returns the date type.
func (v *Version) Date() *Type {
	return v.FindType(nomenclator.Date)
}

// Resources returns the list of resources that are part of this version.
func (v *Version) Resources() []*Resource {
	count := len(v.resources)
	resources := make([]*Resource, count)
	index := 0
	for _, resource := range v.resources {
		resources[index] = resource
		index++
	}
	return resources
}

// FindResource returns the resource with the given name, or nil if there is no such resource.
func (v *Version) FindResource(name *names.Name) *Resource {
	if name == nil {
		return nil
	}
	return v.resources[name.String()]
}

// AddResource adds the given resource to the version.
func (v *Version) AddResource(resource *Resource) {
	if resource != nil {
		v.resources[resource.Name().String()] = resource
		resource.SetOwner(v)
	}
}

// AddResources adds the given resources to the version.
func (v *Version) AddResources(resources []*Resource) {
	for _, resource := range resources {
		v.AddResource(resource)
	}
}

// Root returns the root resource, or nil if there is no such resource.
func (v *Version) Root() *Resource {
	return v.resources[nomenclator.Root.String()]
}

// Errors returns the list of errors that are part of this version.
func (v *Version) Errors() []*Error {
	count := len(v.errors)
	errors := make([]*Error, count)
	index := 0
	for _, err := range v.errors {
		errors[index] = err
		index++
	}
	return errors
}

// FindError returns the error with the given name, or nil if there is no such error.
func (v *Version) FindError(name *names.Name) *Error {
	if name == nil {
		return nil
	}
	return v.errors[name.String()]
}

// AddError adds the given error to the version.
func (v *Version) AddError(err *Error) {
	if err != nil {
		v.errors[err.Name().String()] = err
		err.SetOwner(v)
	}
}

// AddErrors adds the given errors to the version.
func (v *Version) AddErrors(errors []*Error) {
	for _, err := range errors {
		v.AddError(err)
	}
}

func (v *Version) addScalarType(name *names.Name) {
	// Add the scalar type:
	scalarType := NewType()
	scalarType.SetKind(ScalarType)
	scalarType.SetName(name)
	v.AddType(scalarType)

	// Add the list type:
	listType := NewType()
	listType.SetKind(ListType)
	listType.SetName(names.Cat(name, nomenclator.List))
	listType.SetElement(scalarType)
	v.AddType(listType)
}
