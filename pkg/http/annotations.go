/*
Copyright (c) 2022 Red Hat, Inc.

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

package http

import (
	"fmt"

	"github.com/openshift-online/ocm-api-metamodel/v2/pkg/concepts"
)

// httpName checks if the given concept has a `http` annotation. If it does then it returns the
// value of the `name` parameter. If it doesn't, it returns an empty string.
func httpName(concept concepts.Annotated) string {
	annotation := concept.GetAnnotation("http")
	if annotation == nil {
		return ""
	}
	name := annotation.FindParameter("name")
	if name == nil {
		return ""
	}
	return fmt.Sprintf("%s", name)
}

// jsonName checks if the given concept has a `json` annotation. If it does then it returns the
// value of the `name` parameter. If it doesn't, it returns an empty string.
func jsonName(concept concepts.Annotated) string {
	annotation := concept.GetAnnotation("json")
	if annotation == nil {
		return ""
	}
	name := annotation.FindParameter("name")
	if name == nil {
		return ""
	}
	return fmt.Sprintf("%s", name)
}
