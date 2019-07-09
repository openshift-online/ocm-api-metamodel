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

// Error is the representation of a catagery of errors.
type Error struct {
	owner *Version
	doc   string
	name  *names.Name
	code  int
}

// NewError creates a new error.
func NewError() *Error {
	return new(Error)
}

// Owner returns the version that owns this error.
func (e *Error) Owner() *Version {
	return e.owner
}

// SetOwner sets the version that owns this error.
func (e *Error) SetOwner(version *Version) {
	e.owner = version
}

// Doc returns the documentation of this error.
func (e *Error) Doc() string {
	return e.doc
}

// SetDoc sets the documentation of this error.
func (e *Error) SetDoc(value string) {
	e.doc = value
}

// Name returns the name of this error.
func (e *Error) Name() *names.Name {
	return e.name
}

// SetName sets the name of this error.
func (e *Error) SetName(value *names.Name) {
	e.name = value
}

// Code returns the numeric code of this error.
func (e *Error) Code() int {
	return e.code
}

// SetCode sets the numeric code of this error.
func (e *Error) SetCode(value int) {
	e.code = value
}
