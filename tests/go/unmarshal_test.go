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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	amv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/clustersmgmt/v1"
	"github.com/openshift-online/ocm-api-metamodel/tests/go/generated/errors"
)

var _ = Describe("Unmarshal", func() {
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

	It("Ignores wrong struct kind", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"kind": "Junk"
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		Expect(object.Kind()).To(Equal(cmv1.ClusterKind))
		Expect(object.Link()).To(BeFalse())
	})

	It("Ignores wrong list link kind", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"groups": {
				"kind": "Junk"
			}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		list := object.Groups()
		Expect(list).ToNot(BeNil())
		Expect(list.Kind()).To(Equal(cmv1.GroupListKind))
		Expect(list.Link()).To(BeFalse())
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

	It("Can floating point zero", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"factor": 0.0
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Factor()).To(Equal(0.0))
	})

	It("Can read floating point one", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"factor": 1.0
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Factor()).To(Equal(1.0))
	})

	It("Can read floating point minus one", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"factor": -1.0
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.Factor()).To(Equal(-1.0))
	})

	It("Can read date", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"created_at": "2019-07-12T17:12:57.019504Z"
		}`)
		Expect(err).ToNot(HaveOccurred())
		date := object.CreatedAt()
		Expect(date.Year()).To(Equal(2019))
		Expect(date.Month()).To(Equal(time.July))
		Expect(date.Day()).To(Equal(12))
		Expect(date.Hour()).To(Equal(17))
		Expect(date.Minute()).To(Equal(12))
		Expect(date.Second()).To(Equal(57))
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
		Expect(object).To(BeEmpty())
	})

	It("Can read list with one element", func() {
		list, err := cmv1.UnmarshalClusterList(`[
			{
				"name": "mycluster"
			}
		]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(list).ToNot(BeNil())
		Expect(list).To(HaveLen(1))
		Expect(list[0]).ToNot(BeNil())
		Expect(list[0].Name()).To(Equal("mycluster"))
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
		Expect(list).To(HaveLen(2))
		Expect(list[0]).ToNot(BeNil())
		Expect(list[0].Name()).To(Equal("mycluster"))
		Expect(list[1]).ToNot(BeNil())
		Expect(list[1].Name()).To(Equal("yourcluster"))
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

	It("Can read empty map of strings", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"properties": {}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		properties := object.Properties()
		Expect(properties).ToNot(BeNil())
		Expect(properties).To(BeEmpty())
	})

	It("Can read map of strings with one value", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"properties": {
				"mykey": "myvalue"
			}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		properties := object.Properties()
		Expect(properties).ToNot(BeNil())
		Expect(properties).To(HaveLen(1))
		Expect(properties).To(HaveKeyWithValue("mykey", "myvalue"))
	})

	It("Can read map of strings with two values", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"properties": {
				"mykey": "myvalue",
				"yourkey": "yourvalue"
			}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		properties := object.Properties()
		Expect(properties).ToNot(BeNil())
		Expect(properties).To(HaveLen(2))
		Expect(properties).To(HaveKeyWithValue("mykey", "myvalue"))
		Expect(properties).To(HaveKeyWithValue("yourkey", "yourvalue"))
	})

	It("Can read empty map of objects", func() {
		object, err := amv1.UnmarshalRegistryAuths(`{}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
	})

	It("Can read map of objects with one value", func() {
		object, err := amv1.UnmarshalRegistryAuths(`{
			"map": {
				"my.com": {
					"username": "myuser",
					"email": "mymail"
				}
			}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		auths := object.Map()
		Expect(auths).ToNot(BeNil())
		Expect(auths).To(HaveLen(1))
		auth := auths["my.com"]
		Expect(auth).ToNot(BeNil())
		Expect(auth.Username()).To(Equal("myuser"))
		Expect(auth.Email()).To(Equal("mymail"))
	})

	It("Can read map of objects with two values", func() {
		object, err := amv1.UnmarshalRegistryAuths(`{
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
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		auths := object.Map()
		Expect(auths).ToNot(BeNil())
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

	It("Can read an error", func() {
		object, err := errors.UnmarshalError(`{
			"kind": "Error",
			"id": "401",
			"href": "/api/clusters_mgmt/v1/errors/401",
			"code": "CLUSTERS-MGMT-401",
			"reason": "My reason",
			"operation_id": "456",
			"details": {"kind" : "cluster error"}
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).ToNot(BeNil())
		Expect(object.Kind()).To(Equal(errors.ErrorKind))
		Expect(object.ID()).To(Equal("401"))
		Expect(object.HREF()).To(Equal("/api/clusters_mgmt/v1/errors/401"))
		Expect(object.Code()).To(Equal("CLUSTERS-MGMT-401"))
		Expect(object.Reason()).To(Equal("My reason"))
		Expect(object.OperationID()).To(Equal("456"))
		Expect(object.Details()).To(Equal(map[string]interface{}{
			"kind": "cluster error"}))
	})

	It("Can read mixed known and unknown attributes", func() {
		object, err := cmv1.UnmarshalCluster(`{
			"junk": "rubbish",
			"id": "123",
			"first": {
				"second": {
					"third": [
						{
							"fourth": 4
						}
					]
				}
			},
			"name": "mycluster",
			"debris": "litter"
		}`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object.ID()).To(Equal("123"))
		Expect(object.Name()).To(Equal("mycluster"))
	})

	It("Can read empty list of booleans", func() {
		object, err := cmv1.UnmarshalBooleanList(`[]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(BeEmpty())
	})

	It("Can read list with boolean true", func() {
		object, err := cmv1.UnmarshalBooleanList(`[true]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]bool{true}))
	})

	It("Can read list with boolean false", func() {
		object, err := cmv1.UnmarshalBooleanList(`[false]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]bool{false}))
	})

	It("Can read list with multiple booleans", func() {
		object, err := cmv1.UnmarshalBooleanList(`[true, false]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]bool{true, false}))
	})

	It("Can read empty list of integers", func() {
		object, err := cmv1.UnmarshalIntegerList(`[]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(BeEmpty())
	})

	It("Can read list with integer zero", func() {
		object, err := cmv1.UnmarshalIntegerList(`[0]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]int{0}))
	})

	It("Can read list with positive integer", func() {
		object, err := cmv1.UnmarshalIntegerList(`[123]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]int{123}))
	})

	It("Can read list with negative integer", func() {
		object, err := cmv1.UnmarshalIntegerList(`[-123]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]int{-123}))
	})

	It("Can read list with multiple integers", func() {
		object, err := cmv1.UnmarshalIntegerList(`[-123, 0, 123]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]int{-123, 0, 123}))
	})

	It("Can read empty list of floats", func() {
		object, err := cmv1.UnmarshalFloatList(`[]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(BeEmpty())
	})

	It("Can read list with float zero", func() {
		object, err := cmv1.UnmarshalFloatList(`[0.0]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]float64{0.0}))
	})

	It("Can read list with positive float", func() {
		object, err := cmv1.UnmarshalFloatList(`[123.456]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]float64{123.456}))
	})

	It("Can read list with negative float", func() {
		object, err := cmv1.UnmarshalFloatList(`[-123.456]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]float64{-123.456}))
	})

	It("Can read list with multiple floats", func() {
		object, err := cmv1.UnmarshalFloatList(`[-123.456, 0, 123.456]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]float64{-123.456, 0.0, 123.456}))
	})

	It("Can read empty list of strings", func() {
		object, err := cmv1.UnmarshalStringList(`[]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(BeEmpty())
	})

	It("Can read list with empty string", func() {
		object, err := cmv1.UnmarshalStringList(`[""]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]string{""}))
	})

	It("Can read list with one string", func() {
		object, err := cmv1.UnmarshalStringList(`["mystring"]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]string{"mystring"}))
	})

	It("Can read list with multiple strings", func() {
		object, err := cmv1.UnmarshalStringList(`["mystring", "", "yourstring"]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]string{"mystring", "", "yourstring"}))
	})

	It("Can read empty list of dates", func() {
		object, err := cmv1.UnmarshalDateList(`[]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(BeEmpty())
	})

	It("Can read list with one date", func() {
		object, err := cmv1.UnmarshalDateList(`[
			"2019-07-14T15:16:17Z"
		]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]time.Time{
			time.Date(2019, time.July, 14, 15, 16, 17, 0, time.UTC),
		}))
	})

	It("Can read list with multiple dates", func() {
		object, err := cmv1.UnmarshalDateList(`[
			"2019-07-14T15:16:17Z",
			"2019-08-15T16:17:18Z"
		]`)
		Expect(err).ToNot(HaveOccurred())
		Expect(object).To(Equal([]time.Time{
			time.Date(2019, time.July, 14, 15, 16, 17, 0, time.UTC),
			time.Date(2019, time.August, 15, 16, 17, 18, 0, time.UTC),
		}))
	})
})
