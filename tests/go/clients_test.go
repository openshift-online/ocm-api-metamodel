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

// This file contains tests for the methods and send requests and receive responses.

package tests

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	amv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/clustersmgmt/v1"
	"github.com/openshift-online/ocm-api-metamodel/tests/go/generated/errors"
	sbv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/statusboard/v1"
)

var _ = Describe("Client", func() {
	var server *Server
	var transport http.RoundTripper

	BeforeEach(func() {
		server = NewServer()
		transport = NewTransport(server)
	})

	AfterEach(func() {
		server.Close()
	})

	It("Can read response with empty map of objects", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"auths": {}
		}`))
		client := amv1.NewClient(transport, "/api/accounts_mgmt/v1")
		response, err := client.AccessToken().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		auths := response.Map()
		Expect(auths).To(BeEmpty())
	})

	It("Can read response with map of objects with one value", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"map": {
				"my.com": {
					"username": "myuser",
					"email": "mymail"
				}
			}
		}`))
		client := amv1.NewClient(transport, "/api/accounts_mgmt/v1")
		response, err := client.AccessToken().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		auths := response.Map()
		Expect(auths).To(HaveLen(1))
		auth := auths["my.com"]
		Expect(auth).ToNot(BeNil())
		Expect(auth.Username()).To(Equal("myuser"))
		Expect(auth.Email()).To(Equal("mymail"))
	})

	It("Can read response with map of objects with two values", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
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
		client := amv1.NewClient(transport, "/api/accounts_mgmt/v1")
		response, err := client.AccessToken().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		auths := response.Map()
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

	It("Can retrieve version metadata", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"server_version": "123"
		}`))
		client := cmv1.NewClient(transport, "/api/clusters_mgmt/v1")
		response, err := client.Get().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		body := response.Body()
		Expect(body).ToNot(BeNil())
		Expect(body.ServerVersion()).To(Equal("123"))
	})

	It("Can execute action with one input parameter", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodPost,
					"/api/clusters_mgmt/v1/register_disconnected",
				),
				VerifyJSON(`{
					"cluster": {
						"kind": "Cluster",
						"name": "mycluster",
						"external_id": "456"
					}
				}`),
				RespondWith(
					http.StatusOK,
					`{
						"cluster": {
							"kind": "Cluster",
							"id": "123",
							"name": "mycluster",
							"external_id": "456"
						}
					}`,
				),
			),
		)

		// Prepare the description of the cluster:
		cluster, err := cmv1.NewCluster().
			Name("mycluster").
			ExternalID("456").
			Build()
		Expect(err).ToNot(HaveOccurred())

		// Send the request:
		client := cmv1.NewClient(transport, "/api/clusters_mgmt/v1")
		response, err := client.RegisterDisconnected().
			Cluster(cluster).
			Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())

		// Verify the response:
		cluster = response.Cluster()
		Expect(cluster).ToNot(BeNil())
		Expect(cluster.ID()).To(Equal("123"))
		Expect(cluster.Name()).To(Equal("mycluster"))
		Expect(cluster.ExternalID()).To(Equal("456"))
	})

	It("Can execute action with multiple input parameters", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodPost,
					"/api/clusters_mgmt/v1/register_cluster",
				),
				VerifyJSON(`{
					"subscription_id": "123",
					"external_id": "456"
				}`),
				RespondWith(
					http.StatusOK,
					`{
						"cluster": {
							"id": "123",
							"name": "mycluster",
							"external_id": "456"
						}
					}`,
				),
			),
		)

		// Send the request:
		client := cmv1.NewClient(transport, "/api/clusters_mgmt/v1")
		response, err := client.RegisterCluster().
			SubscriptionID("123").
			ExternalID("456").
			Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())

		// Verify the response:
		cluster := response.Cluster()
		Expect(cluster).ToNot(BeNil())
		Expect(cluster.ID()).To(Equal("123"))
		Expect(cluster.Name()).To(Equal("mycluster"))
		Expect(cluster.ExternalID()).To(Equal("456"))
	})

	It("Can retrieve nil list", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodGet,
					"/api/clusters_mgmt/v1/clusters",
				),
				RespondWith(http.StatusOK, `{}`),
			),
		)

		// Send the request:
		client := cmv1.NewClustersClient(transport, "/api/clusters_mgmt/v1/clusters")
		response, err := client.List().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())

		// Verify the result:
		items := response.Items()
		Expect(items).To(BeNil())
	})

	It("Can retrieve empty list", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodGet,
					"/api/clusters_mgmt/v1/clusters",
				),
				RespondWith(
					http.StatusOK,
					`{
						"items": []
					}`,
				),
			),
		)

		// Send the request:
		client := cmv1.NewClustersClient(transport, "/api/clusters_mgmt/v1/clusters")
		response, err := client.List().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())

		// Verify the response:
		items := response.Items()
		Expect(items).ToNot(BeNil())
		Expect(items.Len()).To(BeZero())
	})

	It("Can retrieve list with one element", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodGet,
					"/api/clusters_mgmt/v1/clusters",
				),
				RespondWith(
					http.StatusOK,
					`{
						"items": [
							{
								"id": "123",
								"name": "mycluster"
							}
						]
					}`,
				),
			),
		)

		// Send the request:
		client := cmv1.NewClustersClient(transport, "/api/clusters_mgmt/v1/clusters")
		response, err := client.List().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())

		// Verify the response:
		items := response.Items()
		Expect(items).ToNot(BeNil())
		Expect(items.Len()).To(Equal(1))
		slice := items.Slice()
		Expect(slice).ToNot(BeNil())
		Expect(slice).To(HaveLen(1))
		item := slice[0]
		Expect(item).ToNot(BeNil())
		Expect(item.ID()).To(Equal("123"))
		Expect(item.Name()).To(Equal("mycluster"))
	})

	It("Can retrieve list with two elements", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodGet,
					"/api/clusters_mgmt/v1/clusters",
				),
				RespondWith(
					http.StatusOK,
					`{
						"items": [
							{
								"id": "123",
								"name": "mycluster"
							},
							{
								"id": "456",
								"name": "yourcluster"
							}
						]
					}`,
				),
			),
		)

		// Send the request:
		client := cmv1.NewClustersClient(transport, "/api/clusters_mgmt/v1/clusters")
		response, err := client.List().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())

		// Verify the response:
		items := response.Items()
		Expect(items).ToNot(BeNil())
		Expect(items.Len()).To(Equal(2))
		slice := items.Slice()
		Expect(slice).ToNot(BeNil())
		Expect(slice).To(HaveLen(2))
		first := slice[0]
		Expect(first).ToNot(BeNil())
		Expect(first.ID()).To(Equal("123"))
		Expect(first.Name()).To(Equal("mycluster"))
		second := slice[1]
		Expect(second).ToNot(BeNil())
		Expect(second.ID()).To(Equal("456"))
		Expect(second.Name()).To(Equal("yourcluster"))
	})

	It("Sends paging parameters", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodGet,
					"/api/clusters_mgmt/v1/clusters",
				),
				VerifyFormKV("page", "123"),
				VerifyFormKV("size", "456"),
				RespondWith(http.StatusOK, `{}`),
			),
		)

		// Send the request:
		client := cmv1.NewClustersClient(transport, "/api/clusters_mgmt/v1/clusters")
		response, err := client.List().
			Page(123).
			Size(456).
			Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
	})

	It("Can retrieve paging parameters", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodGet,
					"/api/clusters_mgmt/v1/clusters",
				),
				RespondWith(
					http.StatusOK,
					`{
						"page": 123,
						"size": 456,
						"total": 789,
						"items": []
					}`,
				),
			),
		)

		// Send the request:
		client := cmv1.NewClustersClient(transport, "/api/clusters_mgmt/v1/clusters")
		response, err := client.List().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())

		// Verify the response:
		Expect(response.Page()).To(Equal(123))
		Expect(response.Size()).To(Equal(456))
		Expect(response.Total()).To(Equal(789))
	})

	DescribeTable(
		"Custom query parameters",
		func(value interface{}, expected string) {
			// Prepare the server:
			server.AppendHandlers(
				CombineHandlers(
					VerifyFormKV("my", expected),
					RespondWith(http.StatusOK, `{}`),
				),
			)

			// Send the request:
			client := cmv1.NewClusterClient(transport, "")
			_, err := client.Get().Parameter("my", value).Send()
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("True", true, "true"),
		Entry("False", false, "false"),
		Entry("Zero integer", 0, "0"),
		Entry("Positive integer", 123, "123"),
		Entry("Negative integer", -123, "-123"),
		Entry("Positive integer", 123, "123"),
		Entry("Negative integer", -123, "-123"),
		Entry("Plain string", "myvalue", "myvalue"),
		Entry("Zero float", 0.0, "0"),
		Entry("Positive float", 123.1, "123.1"),
		Entry("Negative float", -123.1, "-123.1"),
		Entry("String that requires encoding", "Áá", "Áá"),
		Entry("String with slash", "my/value", "my/value"),
		Entry("String with space", "my value", "my value"),
		Entry("String with leading space", " myvalue", " myvalue"),
		Entry("String with trailing space", "myvalue ", "myvalue "),
	)

	It("Returns error with HTTP status code", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodDelete,
					"/api/clusters_mgmt/v1/clusters/123",
				),
				RespondWith(
					http.StatusNotFound,
					`{
						"kind": "Error",
						"id": "404",
						"href": "/api/clusters_mgmt/v1/errors/404",
						"code": "CLUSTERS-MGMT-404",
						"reason": "Cluster '123' doesn't exist"
					}`,
				),
			),
		)

		// Send the request:
		client := cmv1.NewClusterClient(transport, "/api/clusters_mgmt/v1/clusters/123")
		_, err := client.Delete().Send()
		Expect(err).To(HaveOccurred())

		// Verify the error:
		var sdkErr *errors.Error
		Expect(err).To(BeAssignableToTypeOf(sdkErr))
		sdkErr = err.(*errors.Error)
		Expect(sdkErr.Status()).To(Equal(http.StatusNotFound))
	})

	It("Sends date parameter in RFC3339 format", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodGet,
					"/api/status_board/v1/statuses",
				),
				VerifyFormKV("created_after", "2022-01-25T14:57:02Z"),
				RespondWith(http.StatusOK, `{}`),
			),
		)

		// Send the request:
		client := sbv1.NewStatusesClient(transport, "/api/status_board/v1/statuses")
		date, err := time.Parse(time.RFC3339, "2022-01-25T15:57:02+01:00")
		Expect(err).ToNot(HaveOccurred())
		response, err := client.List().
			CreatedAfter(date).
			Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
	})

	It("Accepts `204 No Content` with empty response body", func() {
		// Prepare the server:
		server.AppendHandlers(
			RespondWith(http.StatusNoContent, ""),
		)

		// Create a transport that replaces the response body with an empty reader:
		transport = &EmptyResponseBodyTransport{
			wrapped: transport,
		}

		// Send the request:
		client := cmv1.NewClustersClient(transport, "/api/clusters_mgmt/v1/clusters")
		body, err := cmv1.NewCluster().
			Name("my-cluster").
			Build()
		Expect(err).ToNot(HaveOccurred())
		response, err := client.Add().
			Body(body).
			Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response.Status()).To(Equal(http.StatusNoContent))
		Expect(response.Body()).To(BeNil())
	})
})
