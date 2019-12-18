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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	amv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/clustersmgmt/v1"
	"github.com/openshift-online/ocm-api-metamodel/tests/api/errors"
)

var _ = Describe("Writer", func() {
	It("Can write empty map of strings", func() {
		object, err := cmv1.NewCluster().
			Properties(map[string]string{}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = cmv1.MarshalCluster(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"kind": "Cluster"
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

	It("Can write empty map of objects", func() {
		object, err := amv1.NewAccessToken().
			Auths(map[string]*amv1.AuthBuilder{}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = amv1.MarshalAccessToken(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{}`))
	})

	It("Can write map of objects with one value", func() {
		object, err := amv1.NewAccessToken().
			Auths(map[string]*amv1.AuthBuilder{
				"my.com": amv1.NewAuth().Username("myuser").Email("mymail"),
			}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = amv1.MarshalAccessToken(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"auths": {
				"my.com": {
					"username": "myuser",
					"email": "mymail"
				}
			}
		}`))
	})

	It("Can write map of objects with two values", func() {
		object, err := amv1.NewAccessToken().
			Auths(map[string]*amv1.AuthBuilder{
				"my.com":   amv1.NewAuth().Username("myuser").Email("mymail"),
				"your.com": amv1.NewAuth().Username("youruser").Email("yourmail"),
			}).
			Build()
		Expect(err).ToNot(HaveOccurred())
		buffer := new(bytes.Buffer)
		err = amv1.MarshalAccessToken(object, buffer)
		Expect(err).ToNot(HaveOccurred())
		Expect(buffer).To(MatchJSON(`{
			"auths": {
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
})
