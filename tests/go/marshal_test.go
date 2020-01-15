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

// This file contains tests for the methods that write objects to streams.

package tests

import (
	"bytes"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	amv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/clustersmgmt/v1"
	"github.com/openshift-online/ocm-api-metamodel/tests/go/generated/errors"
)

var _ = Describe("Marshal", func() {
	It("Can write nil map of strings", func() {
		object, err := cmv1.NewCluster().
			Properties(nil).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = cmv1.MarshalCluster(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"kind": "Cluster"
		}`))
	})

	It("Can write empty map of strings", func() {
		object, err := cmv1.NewCluster().
			Properties(map[string]string{}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = cmv1.MarshalCluster(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"kind": "Cluster",
			"properties": {}
		}`))
	})

	It("Can write map of strings with one value", func() {
		object, err := cmv1.NewCluster().
			Properties(map[string]string{
				"mykey": "myvalue",
			}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = cmv1.MarshalCluster(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"kind": "Cluster",
			"properties": {
				"mykey": "myvalue"
			}
		}`))
	})

	It("Can write map of strings with two value", func() {
		object, err := cmv1.NewCluster().
			Properties(map[string]string{
				"mykey":   "myvalue",
				"yourkey": "yourvalue",
			}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = cmv1.MarshalCluster(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"kind": "Cluster",
			"properties": {
				"mykey": "myvalue",
				"yourkey": "yourvalue"
			}
		}`))
	})

	It("Can write nil map of objects", func() {
		object, err := amv1.NewRegistryAuths().
			Map(nil).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = amv1.MarshalRegistryAuths(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{}`))
	})

	It("Can write empty map of objects", func() {
		object, err := amv1.NewRegistryAuths().
			Map(map[string]*amv1.RegistryAuthBuilder{}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = amv1.MarshalRegistryAuths(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"map": {}
		}`))
	})

	It("Can write map of objects with one value", func() {
		object, err := amv1.NewRegistryAuths().
			Map(map[string]*amv1.RegistryAuthBuilder{
				"my.com": amv1.NewRegistryAuth().
					Username("myuser").
					Email("mymail"),
			}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = amv1.MarshalRegistryAuths(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"map": {
				"my.com": {
					"username": "myuser",
					"email": "mymail"
				}
			}
		}`))
	})

	It("Can write map of objects with two values", func() {
		object, err := amv1.NewRegistryAuths().
			Map(map[string]*amv1.RegistryAuthBuilder{
				"my.com": amv1.NewRegistryAuth().
					Username("myuser").
					Email("mymail"),
				"your.com": amv1.NewRegistryAuth().
					Username("youruser").
					Email("yourmail"),
			}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = amv1.MarshalRegistryAuths(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"map": {
				"my.com": {
					"username": "myuser",
					"email": "mymail"
				},
				"your.com": {
					"username": "youruser",
					"email": "yourmail"
				}
			}
		}`))
	})

	It("Can write an error", func() {
		object, err := errors.NewError().
			ID("401").
			HREF("/api/clusters_mgmt/v1/errors/401").
			Code("CLUSTERS-MGMT-401").
			Reason("My reason").
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = errors.MarshalError(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"kind": "Error",
			"id": "401",
			"href": "/api/clusters_mgmt/v1/errors/401",
			"code": "CLUSTERS-MGMT-401",
			"reason": "My reason"
		}`))
	})

	It("Can write empty list of booleans", func() {
		buffer := &bytes.Buffer{}
		object := []bool{}
		err := cmv1.MarshalBooleanList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[]`))
	})

	It("Can write list with false", func() {
		buffer := &bytes.Buffer{}
		object := []bool{false}
		err := cmv1.MarshalBooleanList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[false]`))
	})

	It("Can write list with true", func() {
		buffer := &bytes.Buffer{}
		object := []bool{true}
		err := cmv1.MarshalBooleanList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[true]`))
	})

	It("Can write list with multiple booleans", func() {
		buffer := &bytes.Buffer{}
		object := []bool{true, false}
		err := cmv1.MarshalBooleanList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[true, false]`))
	})

	It("Can write list of integers", func() {
		buffer := &bytes.Buffer{}
		object := []int{-1, 0, 1}
		err := cmv1.MarshalIntegerList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[-1, 0, 1]`))
	})

	It("Can write empty list of integers", func() {
		buffer := &bytes.Buffer{}
		object := []int{}
		err := cmv1.MarshalIntegerList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[]`))
	})

	It("Can write list with zero", func() {
		buffer := &bytes.Buffer{}
		object := []int{0}
		err := cmv1.MarshalIntegerList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[0]`))
	})

	It("Can write list with positive integer", func() {
		buffer := &bytes.Buffer{}
		object := []int{123}
		err := cmv1.MarshalIntegerList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[123]`))
	})

	It("Can write list with negative integer", func() {
		buffer := &bytes.Buffer{}
		object := []int{-123}
		err := cmv1.MarshalIntegerList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[-123]`))
	})

	It("Can write list with multiple integers", func() {
		buffer := &bytes.Buffer{}
		object := []int{-123, 0, 123}
		err := cmv1.MarshalIntegerList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[-123, 0, 123]`))
	})

	It("Can write empty list of floats", func() {
		buffer := &bytes.Buffer{}
		object := []float64{}
		err := cmv1.MarshalFloatList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[]`))
	})

	It("Can write list with float zero", func() {
		buffer := &bytes.Buffer{}
		object := []float64{0.0}
		err := cmv1.MarshalFloatList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[0.0]`))
	})

	It("Can write list with positive float", func() {
		buffer := &bytes.Buffer{}
		object := []float64{123.456}
		err := cmv1.MarshalFloatList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[123.456]`))
	})

	It("Can write list with positive float", func() {
		buffer := &bytes.Buffer{}
		object := []float64{-123.456}
		err := cmv1.MarshalFloatList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[-123.456]`))
	})

	It("Can write list with multiple floats", func() {
		buffer := &bytes.Buffer{}
		object := []float64{-123.456, 0, 123.456}
		err := cmv1.MarshalFloatList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[-123.456, 0, 123.456]`))
	})

	It("Can write empty list of strings", func() {
		buffer := &bytes.Buffer{}
		object := []string{}
		err := cmv1.MarshalStringList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[]`))
	})

	It("Can write list with empty string", func() {
		buffer := &bytes.Buffer{}
		object := []string{""}
		err := cmv1.MarshalStringList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[""]`))
	})

	It("Can write list with one string", func() {
		buffer := &bytes.Buffer{}
		object := []string{"mystring"}
		err := cmv1.MarshalStringList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`["mystring"]`))
	})

	It("Can write list with multiple strings", func() {
		buffer := &bytes.Buffer{}
		object := []string{"mystring", "yourstring"}
		err := cmv1.MarshalStringList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`["mystring", "yourstring"]`))
	})

	It("Can write empty list of dates", func() {
		buffer := &bytes.Buffer{}
		object := []time.Time{}
		err := cmv1.MarshalDateList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[]`))
	})

	It("Can write list with one date", func() {
		buffer := &bytes.Buffer{}
		object := []time.Time{
			time.Date(2019, time.July, 14, 15, 16, 17, 0, time.UTC),
		}
		err := cmv1.MarshalDateList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[
			"2019-07-14T15:16:17Z"
		]`))
	})

	It("Can write list with multiple dates", func() {
		buffer := &bytes.Buffer{}
		object := []time.Time{
			time.Date(2019, time.July, 14, 15, 16, 17, 0, time.UTC),
			time.Date(2019, time.August, 15, 16, 17, 18, 0, time.UTC),
		}
		err := cmv1.MarshalDateList(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`[
			"2019-07-14T15:16:17Z",
			"2019-08-15T16:17:18Z"
		]`))
	})
})
