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

// This file contains tests for builders.

package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	amv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/separately_generated/sdk-go/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/separately_generated/sdk-go/clustersmgmt/v1"
	"github.com/openshift-online/ocm-api-metamodel/tests/go/separately_generated/sdk-go/errors"
)

var _ = Describe("Separately Generated Builder", func() {
	Describe("Build", func() {
		It("Can set empty string list attribute", func() {
			object, err := cmv1.NewGithubIdentityProvider().
				Teams().
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object).ToNot(BeNil())
			slice := object.Teams()
			Expect(slice).ToNot(BeNil())
			Expect(slice).To(BeEmpty())
		})

		It("Can set string list attribute with one value", func() {
			object, err := cmv1.NewGithubIdentityProvider().
				Teams("a-team").
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object).ToNot(BeNil())
			slice := object.Teams()
			Expect(slice).ToNot(BeNil())
			Expect(slice).To(HaveLen(1))
			Expect(slice[0]).To(Equal("a-team"))
		})

		It("Can set string list attribute with two values", func() {
			object, err := cmv1.NewGithubIdentityProvider().
				Teams("a-team", "b-team").
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object).ToNot(BeNil())
			slice := object.Teams()
			Expect(slice).ToNot(BeNil())
			Expect(slice).To(HaveLen(2))
			Expect(slice[0]).To(Equal("a-team"))
			Expect(slice[1]).To(Equal("b-team"))
		})

		It("Can set empty struct list attribute", func() {
			object, err := cmv1.NewCluster().
				Groups(
					cmv1.NewGroupList(),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object).ToNot(BeNil())
			list := object.Groups()
			Expect(list).ToNot(BeNil())
			Expect(list.Empty()).To(BeTrue())
		})

		It("Can set struct list attribute with one value", func() {
			items := cmv1.NewGroupList().Items(
				cmv1.NewGroup().ID("a-group"),
			)
			object, err := cmv1.NewCluster().
				Groups(items).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object).ToNot(BeNil())
			list := object.Groups()
			Expect(list).ToNot(BeNil())
			Expect(list.Len()).To(Equal(1))
			slice := list.Slice()
			Expect(slice).ToNot(BeEmpty())
			Expect(slice).To(HaveLen(1))
			Expect(slice[0]).ToNot(BeNil())
			Expect(slice[0].ID()).To(Equal("a-group"))
		})

		It("Can set struct list attribute with two values", func() {
			items := cmv1.NewGroupList().Items(
				cmv1.NewGroup().ID("a-group"),
				cmv1.NewGroup().ID("b-group"),
			)
			object, err := cmv1.NewCluster().
				Groups(items).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object).ToNot(BeNil())
			list := object.Groups()
			Expect(list).ToNot(BeNil())
			Expect(list.Len()).To(Equal(2))
			slice := list.Slice()
			Expect(slice).ToNot(BeEmpty())
			Expect(slice).To(HaveLen(2))
			Expect(slice[0]).ToNot(BeNil())
			Expect(slice[0].ID()).To(Equal("a-group"))
			Expect(slice[1]).ToNot(BeNil())
			Expect(slice[1].ID()).To(Equal("b-group"))
		})

		It("Can set empty map of strings", func() {
			object, err := cmv1.NewCluster().
				Properties(map[string]string{}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			properties := object.Properties()
			Expect(properties).To(BeEmpty())
		})

		It("Can set map of strings with one value", func() {
			object, err := cmv1.NewCluster().
				Properties(map[string]string{
					"mykey": "myvalue",
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			properties := object.Properties()
			Expect(properties).To(HaveLen(1))
			Expect(properties).To(HaveKeyWithValue("mykey", "myvalue"))
		})

		It("Can set map of strings with two values", func() {
			object, err := cmv1.NewCluster().
				Properties(map[string]string{
					"mykey":   "myvalue",
					"yourkey": "yourvalue",
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			properties := object.Properties()
			Expect(properties).To(HaveLen(2))
			Expect(properties).To(HaveKeyWithValue("mykey", "myvalue"))
			Expect(properties).To(HaveKeyWithValue("yourkey", "yourvalue"))
		})

		It("Can set empty map of objects", func() {
			object, err := amv1.NewRegistryAuths().
				Map(map[string]*amv1.RegistryAuthBuilder{}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			auths := object.Map()
			Expect(auths).To(BeEmpty())
		})

		It("Can set map of objects with one value", func() {
			object, err := amv1.NewRegistryAuths().
				Map(map[string]*amv1.RegistryAuthBuilder{
					"my.com": amv1.NewRegistryAuth().
						Username("myuser").
						Email("mymail"),
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			auths := object.Map()
			Expect(auths).To(HaveLen(1))
			auth := auths["my.com"]
			Expect(auth).ToNot(BeNil())
			Expect(auth.Username()).To(Equal("myuser"))
			Expect(auth.Email()).To(Equal("mymail"))
		})

		It("Can set map of objects with two values", func() {
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
			auths := object.Map()
			Expect(auths).To(HaveLen(2))
			first := auths["my.com"]
			Expect(first).ToNot(BeNil())
			Expect(first.Username()).To(Equal("myuser"))
			Expect(first.Email()).To(Equal("mymail"))
			second := auths["your.com"]
			Expect(second).ToNot(BeNil())
			Expect(second.Username()).To(Equal("youruser"))
			Expect(second.Email()).To(Equal("yourmail"))
		})

		It("Can build an error", func() {
			object, err := errors.NewError().
				ID("401").
				HREF("/api/clusters_mgmt/v1/errors/401").
				Code("CLUSTERS-MGMT-401").
				Reason("My reason").
				OperationID("456").
				Details(map[string]interface{}{
					"kind": "cluster error"}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object.ID()).To(Equal("401"))
			Expect(object.HREF()).To(Equal("/api/clusters_mgmt/v1/errors/401"))
			Expect(object.Code()).To(Equal("CLUSTERS-MGMT-401"))
			Expect(object.Reason()).To(Equal("My reason"))
			Expect(object.OperationID()).To(Equal("456"))
			Expect(object.Details()).To(Equal(map[string]interface{}{
				"kind": "cluster error"}))
		})
	})

	Describe("Copy", func() {
		It("Copies simple attribute", func() {
			original, err := cmv1.NewCluster().
				ID("123").
				Name("my").
				Build()
			Expect(err).ToNot(HaveOccurred())
			replica, err := cmv1.NewCluster().
				Copy(original).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(replica.ID()).To(Equal("123"))
			Expect(replica.Name()).To(Equal("my"))
		})

		It("Discards existing values", func() {
			original, err := cmv1.NewCluster().
				ID("123").
				Name("my").
				Build()
			Expect(err).ToNot(HaveOccurred())
			replica, err := cmv1.NewCluster().
				ID("456").
				Name("your").
				Copy(original).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(replica.ID()).To(Equal("123"))
			Expect(replica.Name()).To(Equal("my"))
		})

		It("Copies empty map of strings", func() {
			original, err := cmv1.NewCluster().
				Properties(map[string]string{}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			replica, err := cmv1.NewCluster().
				Copy(original).
				Build()
			Expect(err).ToNot(HaveOccurred())
			properties := replica.Properties()
			Expect(properties).To(BeEmpty())
		})

		It("Copies map of strings with one value", func() {
			original, err := cmv1.NewCluster().
				Properties(map[string]string{
					"mykey": "myvalue",
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			replica, err := cmv1.NewCluster().
				Copy(original).
				Build()
			Expect(err).ToNot(HaveOccurred())
			properties := replica.Properties()
			Expect(properties).To(HaveLen(1))
			Expect(properties).To(HaveKeyWithValue("mykey", "myvalue"))
		})

		It("Copies map of strings with two values", func() {
			original, err := cmv1.NewCluster().
				Properties(map[string]string{
					"mykey":   "myvalue",
					"yourkey": "yourvalue",
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			replica, err := cmv1.NewCluster().
				Copy(original).
				Build()
			Expect(err).ToNot(HaveOccurred())
			properties := replica.Properties()
			Expect(properties).To(HaveLen(2))
			Expect(properties).To(HaveKeyWithValue("mykey", "myvalue"))
			Expect(properties).To(HaveKeyWithValue("yourkey", "yourvalue"))
		})

		It("Copies empty map of objects", func() {
			original, err := amv1.NewRegistryAuths().
				Map(map[string]*amv1.RegistryAuthBuilder{}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			replica, err := amv1.NewRegistryAuths().
				Copy(original).
				Build()
			Expect(err).ToNot(HaveOccurred())
			auths := replica.Map()
			Expect(auths).To(BeEmpty())
		})

		It("Copies map of objects with one value", func() {
			original, err := amv1.NewRegistryAuths().
				Map(map[string]*amv1.RegistryAuthBuilder{
					"my.com": amv1.NewRegistryAuth().
						Username("myuser").
						Email("mymail"),
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			replica, err := amv1.NewRegistryAuths().
				Copy(original).
				Build()
			Expect(err).ToNot(HaveOccurred())
			auths := replica.Map()
			Expect(auths).To(HaveLen(1))
			auth := auths["my.com"]
			Expect(auth).ToNot(BeNil())
			Expect(auth.Username()).To(Equal("myuser"))
			Expect(auth.Email()).To(Equal("mymail"))
		})

		It("Copies map of objects with two values", func() {
			original, err := amv1.NewRegistryAuths().
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
			replica, err := amv1.NewRegistryAuths().
				Copy(original).
				Build()
			Expect(err).ToNot(HaveOccurred())
			auths := replica.Map()
			Expect(auths).To(HaveLen(2))
			first := auths["my.com"]
			Expect(first).ToNot(BeNil())
			Expect(first.Username()).To(Equal("myuser"))
			Expect(first.Email()).To(Equal("mymail"))
			second := auths["your.com"]
			Expect(second).ToNot(BeNil())
			Expect(second.Username()).To(Equal("youruser"))
			Expect(second.Email()).To(Equal("yourmail"))
		})
	})

	Describe("Empty", func() {
		It("Returns `true` for nil builder ", func() {
			var builder *cmv1.ClusterBuilder
			Expect(builder.Empty()).To(BeTrue())
		})

		It("Returns `true` for an empty builder", func() {
			builder := cmv1.NewCluster()
			Expect(builder.Empty()).To(BeTrue())
		})

		It("Returns `false` for an builder with identifier", func() {
			builder := cmv1.NewCluster().
				ID("123")
			Expect(builder.Empty()).To(BeFalse())
		})

		It("Returns `false` for an builder with an string attribute", func() {
			builder := cmv1.NewCluster().
				Name("mycluster")
			Expect(builder.Empty()).To(BeFalse())
		})

		It("Returns `true` for nil list ", func() {
			var list *cmv1.ClusterListBuilder
			Expect(list.Empty()).To(BeTrue())
		})

		It("Returns `true` for empty list ", func() {
			list := cmv1.NewClusterList()
			Expect(list.Empty()).To(BeTrue())
		})

		It("Returns `false` for list with one element", func() {
			list := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("123"),
				)
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for list with two elements", func() {
			list := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("123"),
					cmv1.NewCluster().ID("456"),
				)
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for empty map of strings", func() {
			list := cmv1.NewCluster().
				Properties(map[string]string{})
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for map of strings with one value", func() {
			list := cmv1.NewCluster().
				Properties(map[string]string{
					"mykey": "myvalue",
				})
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for map of strings with two values", func() {
			list := cmv1.NewCluster().
				Properties(map[string]string{
					"mykey":   "myvalue",
					"yourkey": "yourvalue",
				})
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for empty map of builders", func() {
			list := amv1.NewRegistryAuths().
				Map(map[string]*amv1.RegistryAuthBuilder{})
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for map of builders with one value", func() {
			list := amv1.NewRegistryAuths().
				Map(map[string]*amv1.RegistryAuthBuilder{
					"my.com": amv1.NewRegistryAuth().Username("myuser"),
				})
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for map of builders with one value", func() {
			list := amv1.NewRegistryAuths().
				Map(map[string]*amv1.RegistryAuthBuilder{
					"my.com":   amv1.NewRegistryAuth().Username("myuser"),
					"your.com": amv1.NewRegistryAuth().Username("youruser"),
				})
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `true` empty link", func() {
			builder := cmv1.NewCluster().
				Link(true)
			Expect(builder.Empty()).To(BeTrue())
		})
	})
})
