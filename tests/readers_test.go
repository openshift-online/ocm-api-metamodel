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

// This file contains tests for the methods that read objects from streams.

package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/clustersmgmt/v1"
)

var _ = Describe("Reader", func() {
	It("Can read empty object", func() {
		object, err := cmv1.UnmarshalCluster(`{}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		Expect(object.Kind()).To(Equal(cmv1.ClusterKind))
		Expect(object.Link()).To(BeFalse())
		Expect(object.ID()).To(BeEmpty())
		Expect(object.HREF()).To(BeEmpty())
		Expect(object.Name()).To(BeEmpty())
	})

	It("Can read empty link", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"kind": "ClusterLink"
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		Expect(object.Kind()).To(Equal(cmv1.ClusterLinkKind))
		Expect(object.Link()).To(BeTrue())
	})

	It("Can read basic object", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"kind": "Cluster",
			"id": "123",
			"href": "/api/clusters_mgmt/v1/clusters/123",
			"name": "mycluster"
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		Expect(object.Kind()).To(Equal(cmv1.ClusterKind))
		Expect(object.Link()).To(BeFalse())
		Expect(object.ID()).To(Equal("123"))
		Expect(object.HREF()).To(Equal("/api/clusters_mgmt/v1/clusters/123"))
		Expect(object.Name()).To(Equal("mycluster"))
	})

	It("Can read true", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"managed": true
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Managed()).To(BeTrue())
	})

	It("Can read false", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"managed": false
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Managed()).To(BeFalse())
	})

	It("Can read integer zero", func() {
		object, err := cmv1.UnmarshalClusterNodes(`{
			"compute": 0
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Compute()).To(Equal(0))
	})

	It("Can read integer one", func() {
		object, err := cmv1.UnmarshalClusterNodes(`{
			"compute": 1
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Compute()).To(Equal(1))
	})

	It("Can read integer minus one", func() {
		object, err := cmv1.UnmarshalClusterNodes(`{
			"compute": -1
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Compute()).To(Equal(-1))
	})

	It("Can read false", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"managed": false
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Managed()).To(BeFalse())
	})

	It("Can read object with one unknown attribute", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"myname": "myvalue"
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
	})

	It("Can read object with two unknown attributes", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"myname": "myvalue",
			"yourname": "yourvalue"
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
	})

	It("Can read an empty list of objects", func() {
		object, err := cmv1.UnmarshalClusterList(`[]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		Expect(object.Len()).To(BeZero())
	})

	It("Can read list with one element", func() {
		list, err := cmv1.UnmarshalClusterList(`[
			{
				"name": "mycluster"
			}
		]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(list).ToNot(BeNil())
		Expect(list.Len()).To(Equal(1))
		slice := list.Slice()
		Expect(slice).ToNot(BeNil())
		Expect(len(slice)).To(Equal(1))
		Expect(slice[0]).ToNot(BeNil())
		Expect(slice[0].Name()).To(Equal("mycluster"))
	})

	It("Can read list with two elements", func() {
		list, err := cmv1.UnmarshalClusterList(`[
			{
				"name": "mycluster"
			},
			{
				"name": "yourcluster"
			}
		]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(list).ToNot(BeNil())
		Expect(list.Len()).To(Equal(2))
		slice := list.Slice()
		Expect(slice).ToNot(BeNil())
		Expect(len(slice)).To(Equal(2))
		Expect(slice[0]).ToNot(BeNil())
		Expect(slice[0].Name()).To(Equal("mycluster"))
		Expect(slice[1]).ToNot(BeNil())
		Expect(slice[1].Name()).To(Equal("yourcluster"))
	})

	It("Can read empty string list attribute", func() {
		object, err := cmv1.UnmarshalGithubIdentityProvider(`{
			"teams": []
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		slice := object.Teams()
		Expect(slice).ToNot(BeNil())
		Expect(slice).To(BeEmpty())
	})

	It("Can read string list attribute with one element", func() {
		object, err := cmv1.UnmarshalGithubIdentityProvider(`{
			"teams": [
				"a-team"
			]
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		slice := object.Teams()
		Expect(slice).ToNot(BeNil())
		Expect(slice).To(HaveLen(1))
		Expect(slice[0]).To(Equal("a-team"))
	})

	It("Can read string list attribute with two elements", func() {
		object, err := cmv1.UnmarshalGithubIdentityProvider(`{
			"teams": [
				"a-team",
				"b-team"
			]
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		slice := object.Teams()
		Expect(slice).ToNot(BeNil())
		Expect(slice).To(HaveLen(2))
		Expect(slice[0]).To(Equal("a-team"))
		Expect(slice[1]).To(Equal("b-team"))
	})

	It("Can read attribute that is link to a list of instances of a class", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"groups": {
				"kind": "GroupListLink",
				"href": "/api/clusters_mgmt/v1/clusters/123/groups"
			}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		list := object.Groups()
		Expect(list).ToNot(BeNil())
		Expect(list.Kind()).To(Equal(cmv1.GroupListLinkKind))
		Expect(list.Link()).To(BeTrue())
		Expect(list.HREF()).To(Equal("/api/clusters_mgmt/v1/clusters/123/groups"))
		Expect(list.Len()).To(BeZero())
	})

	It("Can read attribute that is a list of zero instances of a class", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"groups": {
				"kind": "GroupList",
				"href": "/api/clusters_mgmt/v1/clusters/123/groups",
				"items": []
			}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		list := object.Groups()
		Expect(list).ToNot(BeNil())
		Expect(list.Kind()).To(Equal(cmv1.GroupListKind))
		Expect(list.Link()).To(BeFalse())
		Expect(list.HREF()).To(Equal("/api/clusters_mgmt/v1/clusters/123/groups"))
		Expect(list.Len()).To(BeZero())
	})

	It("Can read attribute that is a list of one instance of a class", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"groups": {
				"kind": "GroupList",
				"href": "/api/clusters_mgmt/v1/clusters/123/groups",
				"items": [
					{
						"kind": "Group",
						"href": "/api/clusters_mgmt/v1/clusters/123/groups/456",
						"id": "456"
					}
				]
			}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		list := object.Groups()
		Expect(list).ToNot(BeNil())
		Expect(list.Kind()).To(Equal(cmv1.GroupListKind))
		Expect(list.Link()).To(BeFalse())
		Expect(list.HREF()).To(Equal("/api/clusters_mgmt/v1/clusters/123/groups"))
		Expect(list.Len()).To(Equal(1))
		slice := list.Slice()
		Expect(slice).ToNot(BeNil())
		Expect(slice).To(HaveLen(1))
		Expect(slice[0].Kind()).To(Equal(cmv1.GroupKind))
		Expect(slice[0].HREF()).To(Equal("/api/clusters_mgmt/v1/clusters/123/groups/456"))
		Expect(slice[0].ID()).To(Equal("456"))
	})

	It("Can read attribute that is a list of two instances of a class", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"groups": {
				"kind": "GroupList",
				"href": "/api/clusters_mgmt/v1/clusters/123/groups",
				"items": [
					{
						"kind": "Group",
						"href": "/api/clusters_mgmt/v1/clusters/123/groups/456",
						"id": "456"
					},
					{
						"kind": "Group",
						"href": "/api/clusters_mgmt/v1/clusters/123/groups/789",
						"id": "789"
					}
				]
			}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		list := object.Groups()
		Expect(list).ToNot(BeNil())
		Expect(list.Kind()).To(Equal(cmv1.GroupListKind))
		Expect(list.Link()).To(BeFalse())
		Expect(list.HREF()).To(Equal("/api/clusters_mgmt/v1/clusters/123/groups"))
		Expect(list.Len()).To(Equal(2))
		slice := list.Slice()
		Expect(slice).ToNot(BeNil())
		Expect(slice).To(HaveLen(2))
		Expect(slice[0].Kind()).To(Equal(cmv1.GroupKind))
		Expect(slice[0].HREF()).To(Equal("/api/clusters_mgmt/v1/clusters/123/groups/456"))
		Expect(slice[0].ID()).To(Equal("456"))
		Expect(slice[1].Kind()).To(Equal(cmv1.GroupKind))
		Expect(slice[1].HREF()).To(Equal("/api/clusters_mgmt/v1/clusters/123/groups/789"))
		Expect(slice[1].ID()).To(Equal("789"))
	})
})
