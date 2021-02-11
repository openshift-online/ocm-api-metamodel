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

// This file contains tests for the Poll method.

package tests

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/clustersmgmt/v1"
)

var _ = Describe("Poll", func() {
	var server *Server
	var transport http.RoundTripper

	BeforeEach(func() {
		server = NewServer()
		transport = NewTransport(server)
	})

	AfterEach(func() {
		server.Close()
	})

	It("Refuses to start without a timeout or deadline", func() {
		client := cmv1.NewClusterClient(transport, "")
		ctx := context.Background()
		_, err := client.Poll().
			Interval(1 * time.Millisecond).
			StartContext(ctx)
		Expect(err).To(HaveOccurred())
	})

	It("Refuses to start without an interval", func() {
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond)
		_, err := client.Poll().
			StartContext(ctx)
		Expect(err).To(HaveOccurred())
	})

	It("Refuses to start without zero interval", func() {
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond)
		_, err := client.Poll().
			Interval(0).
			StartContext(ctx)
		Expect(err).To(HaveOccurred())
	})

	It("Refuses to start without negative interval", func() {
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond)
		_, err := client.Poll().
			Interval(-1 * time.Millisecond).
			StartContext(ctx)
		Expect(err).To(HaveOccurred())
	})

	It("Finishes if the first request returns the default status", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"kind": "Cluster",
			"id": "123",
			"href": "/api/clusters_mgmt/v1/clusters/123",
			"name": "mycluster"
		}`))
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
		response, err := client.Poll().
			Interval(1 * time.Millisecond).
			StartContext(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		status := response.Status()
		Expect(status).To(Equal(http.StatusOK))
		result := response.Body()
		Expect(result).ToNot(BeNil())
		Expect(result.ID()).To(Equal("123"))
		Expect(result.Name()).To(Equal("mycluster"))
	})

	It("Finishes if the first request returns a non default status", func() {
		server.AppendHandlers(RespondWith(http.StatusNotFound, `{
			"kind": "Error",
			"id": "404",
			"href": "/api/clusters_mgmt/v1/errors/404",
			"code": "CLUSTERS-MGMT-404",
			"reason": "Not found"
		}`))
		client := cmv1.NewClusterClient(transport, "/api/clusters_mgmt/v1/clusters/123")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
		response, err := client.Poll().
			Interval(1 * time.Millisecond).
			Status(http.StatusNotFound).
			StartContext(ctx)
		Expect(err).To(HaveOccurred())
		status := response.Status()
		Expect(status).To(Equal(http.StatusNotFound))
		result := response.Error()
		Expect(result).ToNot(BeNil())
		Expect(result.ID()).To(Equal("404"))
		Expect(result.Code()).To(Equal("CLUSTERS-MGMT-404"))
		Expect(result.Reason()).To(Equal("Not found"))
	})

	It("Finishes if the second request returns the default status", func() {
		server.AppendHandlers(RespondWith(http.StatusNotFound, `{
			"kind": "Error",
			"id": "404",
			"href": "/api/clusters_mgmt/v1/errors/404",
			"code": "CLUSTERS-MGMT-404",
			"reason": "Not found"
		}`))
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"kind": "Cluster",
			"id": "123",
			"href": "/api/clusters_mgmt/v1/clusters/123",
			"name": "mycluster"
		}`))
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
		response, err := client.Poll().
			Interval(1 * time.Millisecond).
			StartContext(ctx)
		Expect(err).ToNot(HaveOccurred())
		status := response.Status()
		Expect(status).To(Equal(http.StatusOK))
		result := response.Body()
		Expect(result).ToNot(BeNil())
		Expect(result.ID()).To(Equal("123"))
		Expect(result.Name()).To(Equal("mycluster"))
	})

	It("Finishes if the second request returns the default status", func() {
		server.AppendHandlers(RespondWith(http.StatusNotFound, `{
			"kind": "Error",
			"id": "404",
			"href": "/api/clusters_mgmt/v1/errors/404",
			"code": "CLUSTERS-MGMT-404",
			"reason": "Not found"
		}`))
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"kind": "Cluster",
			"id": "123",
			"href": "/api/clusters_mgmt/v1/clusters/123",
			"name": "mycluster"
		}`))
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
		response, err := client.Poll().
			Interval(1 * time.Millisecond).
			StartContext(ctx)
		Expect(err).ToNot(HaveOccurred())
		status := response.Status()
		Expect(status).To(Equal(http.StatusOK))
		result := response.Body()
		Expect(result).ToNot(BeNil())
		Expect(result.ID()).To(Equal("123"))
		Expect(result.Name()).To(Equal("mycluster"))
	})

	It("Finishes if the timeout expires", func() {
		server.AppendHandlers(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(20 * time.Millisecond)
			},
			RespondWith(http.StatusOK, `{
				"kind": "Cluster",
				"id": "123",
				"href": "/api/clusters_mgmt/v1/clusters/123",
				"name": "mycluster"
			}`),
		)
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
		response, err := client.Poll().
			Interval(1 * time.Millisecond).
			StartContext(ctx)
		Expect(err).To(HaveOccurred())
		Expect(ctx.Err()).ToNot(BeNil())
		Expect(response).To(BeNil())
	})

	It("Finishes if the first request satisfies the conditions", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"kind": "Cluster",
			"id": "123",
			"href": "/api/clusters_mgmt/v1/clusters/123",
			"name": "mycluster",
			"state": "ready"
		}`))
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
		response, err := client.Poll().
			Interval(1 * time.Millisecond).
			Predicate(func(response *cmv1.ClusterGetResponse) bool {
				return response.Body().State() == cmv1.ClusterStateReady
			}).
			StartContext(ctx)
		Expect(err).ToNot(HaveOccurred())
		status := response.Status()
		Expect(status).To(Equal(http.StatusOK))
		result := response.Body()
		Expect(result).ToNot(BeNil())
		Expect(result.ID()).To(Equal("123"))
		Expect(result.Name()).To(Equal("mycluster"))
		Expect(result.State()).To(Equal(cmv1.ClusterStateReady))
	})

	It("Finishes if the second request satisfies the conditions", func() {
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"kind": "Cluster",
			"id": "123",
			"href": "/api/clusters_mgmt/v1/clusters/123",
			"name": "mycluster",
			"state": "installing"
		}`))
		server.AppendHandlers(RespondWith(http.StatusOK, `{
			"kind": "Cluster",
			"id": "123",
			"href": "/api/clusters_mgmt/v1/clusters/123",
			"name": "mycluster",
			"state": "ready"
		}`))
		client := cmv1.NewClusterClient(transport, "")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
		response, err := client.Poll().
			Interval(1 * time.Millisecond).
			Predicate(func(response *cmv1.ClusterGetResponse) bool {
				return response.Body().State() == cmv1.ClusterStateReady
			}).
			StartContext(ctx)
		Expect(err).ToNot(HaveOccurred())
		status := response.Status()
		Expect(status).To(Equal(http.StatusOK))
		result := response.Body()
		Expect(result).ToNot(BeNil())
		Expect(result.ID()).To(Equal("123"))
		Expect(result.Name()).To(Equal("mycluster"))
		Expect(result.State()).To(Equal(cmv1.ClusterStateReady))
	})
})
