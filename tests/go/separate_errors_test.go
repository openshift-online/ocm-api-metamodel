/*
Copyright (c) 2020 Red Hat, Inc.

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

// This file contains tests for errors.

package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-online/ocm-api-metamodel/tests/go/separately_generated/sdk-go/errors"
)

var _ = Describe("Separately Generated Errors", func() {
	DescribeTable(
		"Conversion to string",
		func(input, expected string) {
			object, err := errors.UnmarshalError(input)
			Expect(err).ToNot(HaveOccurred())
			Expect(object).ToNot(BeNil())

			// Check the `Error` method:
			actual := object.Error()
			Expect(actual).To(Equal(expected))

			// Check the `String` method:
			actual = object.String()
			Expect(actual).To(Equal(expected))
		},
		Entry(
			"Empty",
			`{}`,
			"unknown error",
		),
		Entry(
			"Only identifier",
			`{
				"id": "401"
			}`,
			"identifier is '401'",
		),
		Entry(
			"Only code",
			`{
				"code": "CLUSTERS-MGMT-401"
			}`,
			"code is 'CLUSTERS-MGMT-401'",
		),
		Entry(
			"Only operation identifier",
			`{
				"operation_id": "456"
			}`,
			"operation identifier is '456'",
		),
		Entry(
			"Only reason",
			`{
				"reason": "My reason"
			}`,
			`My reason`,
		),
		Entry(
			"Identifier and code",
			`{
				"id": "401",
				"code": "CLUSTERS-MGMT-401"
			}`,
			"identifier is '401' and code is 'CLUSTERS-MGMT-401'",
		),
		Entry(
			"Identifier, code and operation identifier",
			`{
				"id": "401",
				"code": "CLUSTERS-MGMT-401",
				"operation_id": "456"
			}`,
			"identifier is '401', code is 'CLUSTERS-MGMT-401' and operation "+
				"identifier is '456'",
		),
		Entry(
			"Identifier, code, operation identifier and reason",
			`{
				"id": "401",
				"code": "CLUSTERS-MGMT-401",
				"reason": "My reason",
				"operation_id": "456"
			}`,
			"identifier is '401', code is 'CLUSTERS-MGMT-401' and operation "+
				"identifier is '456': My reason",
		),
		Entry(
			"Status, identifier, code, operation identifier and reason",
			`{
				"id": "401",
				"status": 401,
				"code": "CLUSTERS-MGMT-401",
				"reason": "My reason",
				"operation_id": "456"
			}`,
			"status is 401, identifier is '401', code is 'CLUSTERS-MGMT-401' and "+
				"operation identifier is '456': My reason",
		),
	)
})
