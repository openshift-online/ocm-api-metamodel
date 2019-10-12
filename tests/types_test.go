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

// This file contains tests for types.

package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	amv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/clustersmgmt/v1"
)

var _ = Describe("Type", func() {
	Describe("Kind", func() {
		It("Returns nil nil", func() {
			var object *cmv1.Cluster
			Expect(object.Kind()).To(Equal(cmv1.ClusterNilKind))
		})
	})

	Describe("Link", func() {
		It("Returns false on nil", func() {
			var object *cmv1.Cluster
			Expect(object.Link()).To(BeFalse())
		})
	})

	Describe("ID", func() {
		It("Can get value of nil", func() {
			var object *cmv1.Cluster
			Expect(object.ID()).To(BeEmpty())
		})

		It("Can check value of nil", func() {
			var object *cmv1.Cluster
			value, ok := object.GetID()
			Expect(ok).To(BeFalse())
			Expect(value).To(BeEmpty())
		})
	})

	Describe("HREF", func() {
		It("Can get value of nil", func() {
			var object *cmv1.Cluster
			Expect(object.HREF()).To(BeEmpty())
		})

		It("Can check value of nil", func() {
			var object *cmv1.Cluster
			value, ok := object.GetHREF()
			Expect(ok).To(BeFalse())
			Expect(value).To(BeEmpty())
		})
	})

	Describe("String attribute", func() {
		It("Can get value of nil", func() {
			var object *cmv1.Cluster
			Expect(object.Name()).To(BeEmpty())
		})

		It("Can check value of nil", func() {
			var object *cmv1.Cluster
			value, ok := object.GetName()
			Expect(ok).To(BeFalse())
			Expect(value).To(BeEmpty())
		})
	})

	Describe("Get", func() {
		It("Returns nil for nil list ", func() {
			var list *cmv1.ClusterList
			Expect(list.Get(0)).To(BeNil())
		})

		It("Returns nil for empty list ", func() {
			list, err := cmv1.NewClusterList().Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Get(0)).To(BeNil())
		})

		It("Returns nil for negative index", func() {
			list, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("123"),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Get(-1)).To(BeNil())
		})

		It("Returns nil for positive index out of range", func() {
			list, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("123"),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Get(1)).To(BeNil())
		})

		It("Returns first item for zero", func() {
			list, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("0"),
					cmv1.NewCluster().ID("1"),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			item := list.Get(0)
			Expect(item).ToNot(BeNil())
			Expect(item.ID()).To(Equal("0"))
		})

		It("Returns second item for one", func() {
			list, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("0"),
					cmv1.NewCluster().ID("1"),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			item := list.Get(1)
			Expect(item).ToNot(BeNil())
			Expect(item.ID()).To(Equal("1"))
		})
	})

	Describe("Len", func() {
		It("Returns zero for nil list ", func() {
			var list *cmv1.ClusterList
			Expect(list.Len()).To(BeZero())
		})

		It("Returns zero for empty list ", func() {
			list, err := cmv1.NewClusterList().Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Len()).To(BeZero())
		})

		It("Returns one for list with one element", func() {
			list, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("123"),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Len()).To(Equal(1))
		})

		It("Returns two for list with two elements", func() {
			list, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("123"),
					cmv1.NewCluster().ID("456"),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Len()).To(Equal(2))
		})
	})

	Describe("Empty", func() {
		It("Returns `true` for nil object ", func() {
			var object *cmv1.Cluster
			Expect(object.Empty()).To(BeTrue())
		})

		It("Returns `true` for an empty object", func() {
			object, err := cmv1.NewCluster().Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object.Empty()).To(BeTrue())
		})

		It("Returns `false` for an object with identifier", func() {
			object, err := cmv1.NewCluster().
				ID("123").
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object.Empty()).To(BeFalse())
		})

		It("Returns `false` for an object with an string attribute", func() {
			object, err := cmv1.NewCluster().
				Name("mycluster").
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(object.Empty()).To(BeFalse())
		})

		It("Returns `true` for nil list ", func() {
			var list *cmv1.ClusterList
			Expect(list.Empty()).To(BeTrue())
		})

		It("Returns `true` for empty list ", func() {
			list, err := cmv1.NewClusterList().Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeTrue())
		})

		It("Returns `false` for list with one element", func() {
			list, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("123"),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for list with two elements", func() {
			list, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().ID("123"),
					cmv1.NewCluster().ID("456"),
				).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `true` for empty map of strings", func() {
			list, err := cmv1.NewCluster().
				Properties(map[string]string{}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeTrue())
		})

		It("Returns `false` for map of strings with one value", func() {
			list, err := cmv1.NewCluster().
				Properties(map[string]string{
					"mykey": "myvalue",
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for map of strings with two values", func() {
			list, err := cmv1.NewCluster().
				Properties(map[string]string{
					"mykey":   "myvalue",
					"yourkey": "yourvalue",
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `true` for empty map of objects", func() {
			list, err := amv1.NewAccessToken().
				Auths(map[string]*amv1.AuthBuilder{}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeTrue())
		})

		It("Returns `false` for map of objects with one value", func() {
			list, err := amv1.NewAccessToken().
				Auths(map[string]*amv1.AuthBuilder{
					"my.com": amv1.NewAuth().Username("myuser"),
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeFalse())
		})

		It("Returns `false` for map of objects with one value", func() {
			list, err := amv1.NewAccessToken().
				Auths(map[string]*amv1.AuthBuilder{
					"my.com":   amv1.NewAuth().Username("myuser"),
					"your.com": amv1.NewAuth().Username("youruser"),
				}).
				Build()
			Expect(err).ToNot(HaveOccurred())
			Expect(list.Empty()).To(BeFalse())
		})
	})
})
