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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	amv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/clustersmgmt/v1"
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
		client := amv1.NewAccessTokenClient(transport, "", "")
		response, err := client.Post().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		token := response.Body()
		Expect(token).ToNot(BeNil())
		auths := token.Auths()
		Expect(auths).To(BeEmpty())
	})

	It("Can read response with map of objects with one value", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"auths": {
				"my.com": {
					"username": "myuser",
					"email": "mymail"
				}
			}
		}`))
		client := amv1.NewAccessTokenClient(transport, "", "")
		response, err := client.Post().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		token := response.Body()
		Expect(token).ToNot(BeNil())
		auths := token.Auths()
		Expect(auths).To(HaveLen(1))
		auth := auths["my.com"]
		Expect(auth).ToNot(BeNil())
		Expect(auth.Username()).To(Equal("myuser"))
		Expect(auth.Email()).To(Equal("mymail"))
	})

	It("Can read response with map of objects with two values", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
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
		client := amv1.NewAccessTokenClient(transport, "", "")
		response, err := client.Post().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		token := response.Body()
		Expect(token).ToNot(BeNil())
		auths := token.Auths()
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
		client := cmv1.NewClient(transport, "", "")
		response, err := client.Get().Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		body := response.Body()
		Expect(body).ToNot(BeNil())
		Expect(body.ServerVersion()).To(Equal("123"))
	})

	It("Can execute non REST method", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodPost,
					"/api/clusters_mgmt/v1/register_cluster",
				),
				RespondWith(
					http.StatusOK,
					`{
						"id": "123",
						"name": "mycluster",
						"external_id": "456"
					}`,
				),
			),
		)

		// Prepare the request body:
		body, err := cmv1.NewCluster().
			ID("123").
			Name("mycluster").
			ExternalID("456").
			Build()
		Expect(err).ToNot(HaveOccurred())

		// Send the request:
		client := cmv1.NewClient(transport, "/api/clusters_mgmt/v1", "")
		response, err := client.RegisterCluster().
			Body(body).
			Send()
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())

		// Verify the response:
		body = response.Body()
		Expect(body).ToNot(BeNil())
		Expect(body.ID()).To(Equal("123"))
		Expect(body.Name()).To(Equal("mycluster"))
		Expect(body.ExternalID()).To(Equal("456"))
	})
})
