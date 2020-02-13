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

// This file contains tests for servers.

package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/openshift-online/ocm-api-metamodel/tests/go/generated"
	am "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/accountsmgmt"
	az "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/authorizations"
	cm "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/clustersmgmt"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/go/generated/clustersmgmt/v1"
)

var _ = Describe("Server", func() {
	var (
		server   *MyServer
		adapter  *generated.Adapter
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		// Create the server:
		server = &MyServer{
			clustersMgmt: &MyCMServer{
				v1: &MyCMV1Server{
					clusters: &MyClustersServer{
						cluster: &MyClusterServer{},
					},
				},
			},
		}

		// Create the adapter:
		adapter = generated.NewAdapter(server)

		// Create the recorder:
		recorder = httptest.NewRecorder()
	})

	Describe("Parameter default values", func() {
		It("Assigns integer default", func() {
			// Prepare the server:
			server.clustersMgmt.v1.clusters.list = func(
				ctx context.Context,
				request *cmv1.ClustersListServerRequest,
				response *cmv1.ClustersListServerResponse,
			) error {
				Expect(request.Size()).To(Equal(100))
				return nil
			}

			// Send the request:
			request := httptest.NewRequest(
				http.MethodGet,
				"/api/clusters_mgmt/v1/clusters",
				nil,
			)
			adapter.ServeHTTP(recorder, request)
		})

		It("Overrides integer default", func() {
			// Prepare the server:
			server.clustersMgmt.v1.clusters.list = func(
				ctx context.Context,
				request *cmv1.ClustersListServerRequest,
				response *cmv1.ClustersListServerResponse,
			) error {
				Expect(request.Size()).To(Equal(10))
				return nil
			}

			// Send the request:
			request := httptest.NewRequest(
				http.MethodGet,
				"/api/clusters_mgmt/v1/clusters?size=10",
				nil,
			)
			adapter.ServeHTTP(recorder, request)
		})

		It("Assigns boolean default", func() {
			// Prepare the server:
			server.clustersMgmt.v1.clusters.cluster.del = func(
				ctx context.Context,
				request *cmv1.ClusterDeleteServerRequest,
				response *cmv1.ClusterDeleteServerResponse,
			) error {
				Expect(request.Deprovision()).To(BeTrue())
				return nil
			}

			// Send the request:
			request := httptest.NewRequest(
				http.MethodDelete,
				"/api/clusters_mgmt/v1/clusters/123",
				nil,
			)
			adapter.ServeHTTP(recorder, request)
		})

		It("Overrides boolean default", func() {
			// Prepare the server:
			server.clustersMgmt.v1.clusters.cluster.del = func(
				ctx context.Context,
				request *cmv1.ClusterDeleteServerRequest,
				response *cmv1.ClusterDeleteServerResponse,
			) error {
				Expect(request.Deprovision()).To(BeFalse())
				return nil
			}

			// Send the request:
			request := httptest.NewRequest(
				http.MethodDelete,
				"/api/clusters_mgmt/v1/clusters/123?deprovision=false",
				nil,
			)
			adapter.ServeHTTP(recorder, request)
		})

		It("Assigns string default", func() {
			// Prepare the server:
			server.clustersMgmt.v1.clusters.cluster.del = func(
				ctx context.Context,
				request *cmv1.ClusterDeleteServerRequest,
				response *cmv1.ClusterDeleteServerResponse,
			) error {
				Expect(request.Reason()).To(Equal("myreason"))
				return nil
			}

			// Send the request:
			request := httptest.NewRequest(
				http.MethodDelete,
				"/api/clusters_mgmt/v1/clusters/123",
				nil,
			)
			adapter.ServeHTTP(recorder, request)
		})

		It("Overrides string default", func() {
			// Prepare the server:
			server.clustersMgmt.v1.clusters.cluster.del = func(
				ctx context.Context,
				request *cmv1.ClusterDeleteServerRequest,
				response *cmv1.ClusterDeleteServerResponse,
			) error {
				Expect(request.Reason()).To(Equal("yourreason"))
				return nil
			}

			// Send the request:
			request := httptest.NewRequest(
				http.MethodDelete,
				"/api/clusters_mgmt/v1/clusters/123?reason=yourreason",
				nil,
			)
			adapter.ServeHTTP(recorder, request)
		})
	})

	It("Returns the list of clusters with a trailing slash", func() {
		// Prepare the server:
		server.clustersMgmt.v1.clusters.list = func(
			ctx context.Context,
			request *cmv1.ClustersListServerRequest,
			response *cmv1.ClustersListServerResponse,
		) error {
			items, err := cmv1.NewClusterList().Build()
			if err != nil {
				return err
			}
			response.Items(items)
			response.Page(1)
			response.Size(0)
			response.Total(0)
			return nil
		}

		// Send the request:
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/clusters_mgmt/v1/clusters/",
			nil,
		)
		adapter.ServeHTTP(recorder, request)

		// Verify the response:
		Expect(recorder.Code).To(Equal(http.StatusOK))
	})

	It("Returns a 404 for an unknown resource", func() {
		request := httptest.NewRequest(http.MethodGet, "/foo", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Can get a list of clusters", func() {
		// Prepare the server:
		server.clustersMgmt.v1.clusters.list = func(
			ctx context.Context,
			request *cmv1.ClustersListServerRequest,
			response *cmv1.ClustersListServerResponse,
		) error {
			items, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().
						Name("mycluster"),
				).
				Build()
			if err != nil {
				return err
			}
			response.Items(items)
			response.Page(1)
			response.Size(1)
			response.Total(1)

			return nil
		}

		// Send the request:
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/clusters_mgmt/v1/clusters",
			nil,
		)
		adapter.ServeHTTP(recorder, request)

		// Verify the response:
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"kind": "ClusterList",
			"page": 1,
			"size": 1,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "mycluster"
				}
			]
		}`))
	})

	It("Can get a list of clusters by page", func() {
		// Prepare the server:
		server.clustersMgmt.v1.clusters.list = func(
			ctx context.Context,
			request *cmv1.ClustersListServerRequest,
			response *cmv1.ClustersListServerResponse,
		) error {
			// Verify the request:
			Expect(request.Page()).To(Equal(2))

			// Send the response:
			items, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().
						Name("mycluster"),
				).
				Build()
			if err != nil {
				return err
			}
			response.Items(items)
			response.Page(2)
			response.Size(1)
			response.Total(1)

			return nil
		}

		// Send the request:
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/clusters_mgmt/v1/clusters?page=2",
			nil,
		)
		adapter.ServeHTTP(recorder, request)

		// Verify the response:
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"kind": "ClusterList",
			"page": 2,
			"size": 1,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "mycluster"
				}
			]
		}`))
	})

	It("Can get a list of clusters by size", func() {
		// Prepare the server:
		server.clustersMgmt.v1.clusters.list = func(
			ctx context.Context,
			request *cmv1.ClustersListServerRequest,
			response *cmv1.ClustersListServerResponse,
		) error {
			// Verify the request:
			Expect(request.Size()).To(Equal(2))

			// Send the response:
			items, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().
						Name("mycluster"),
					cmv1.NewCluster().
						Name("yourcluster"),
				).
				Build()
			if err != nil {
				return err
			}
			response.Items(items)
			response.Page(1)
			response.Size(2)
			response.Total(2)

			return nil
		}

		// Send the request:
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/clusters_mgmt/v1/clusters?size=2",
			nil,
		)
		adapter.ServeHTTP(recorder, request)

		// Verify the response:
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"kind": "ClusterList",
			"page": 1,
			"size": 2,
			"total": 2,
			"items": [
				{
					"kind": "Cluster",
					"name": "mycluster"
				},
				{
					"kind": "Cluster",
					"name": "yourcluster"
				}
			]
		}`))
	})

	It("Can get a list of clusters by size and page", func() {
		// Prepare the server:
		server.clustersMgmt.v1.clusters.list = func(
			ctx context.Context,
			request *cmv1.ClustersListServerRequest,
			response *cmv1.ClustersListServerResponse,
		) error {
			// Verify the request:
			Expect(request.Page()).To(Equal(2))
			Expect(request.Size()).To(Equal(2))

			// Send the response:
			items, err := cmv1.NewClusterList().
				Items(
					cmv1.NewCluster().
						Name("mycluster"),
					cmv1.NewCluster().
						Name("yourcluster"),
				).
				Build()
			if err != nil {
				return err
			}
			response.Items(items)
			response.Page(2)
			response.Size(2)
			response.Total(4)

			return nil
		}

		// Send the request:
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/clusters_mgmt/v1/clusters?size=2&page=2",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"kind": "ClusterList",
			"page": 2,
			"size": 2,
			"total": 4,
			"items": [
				{
					"kind": "Cluster",
					"name": "mycluster"
				},
				{
					"kind": "Cluster",
					"name": "yourcluster"
				}
			]
		}`))
	})

	It("Can get a cluster by identifier", func() {
		// Prepare the server:
		server.clustersMgmt.v1.clusters.cluster.get = func(
			ctx context.Context,
			request *cmv1.ClusterGetServerRequest,
			response *cmv1.ClusterGetServerResponse,
		) error {
			cluster, err := cmv1.NewCluster().
				Name("mycluster").
				Build()
			if err != nil {
				return err
			}
			response.Body(cluster)
			return nil
		}

		// Send the request:
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/clusters_mgmt/v1/clusters/123",
			nil,
		)
		adapter.ServeHTTP(recorder, request)

		// Verify the response:
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"kind": "Cluster",
			"name": "mycluster"
		}`))
	})

	It("Can get a cluster sub resource by id", func() {
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/clusters_mgmt/v1/clusters/123/identity_providers",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"kind": "IdentityProviderList",
			"page": 1,
			"size": 1,
			"total": 1,
			"items": [
				{
					"kind": "IdentityProvider",
					"name": "test-list-identity-providers"
				}
			]
		}`))
	})

	It("Returns 404 for an unknown service", func() {
		request := httptest.NewRequest(http.MethodGet, "/api/foo_mgmt", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 404 for an unknown version", func() {
		request := httptest.NewRequest(http.MethodGet, "/api/clusters_mgmt/v100", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 404 for an unknown resource", func() {
		request := httptest.NewRequest(http.MethodGet, "/api/clusters_mgmt/flusters", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 404 for an unknown sub-resource", func() {
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/clusters_mgmt/v1/clusters/123/foo",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 404 if the server returns nil for a locator", func() {
		request := httptest.NewRequest(http.MethodGet, "/api/clusters_mgmt/v1/nil", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 405 for unsupported service method", func() {
		request := httptest.NewRequest(http.MethodPost, "/api/clusters_mgmt", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("Returns 405 for unsupported version method", func() {
		request := httptest.NewRequest(http.MethodPost, "/api/clusters_mgmt/v1", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("Returns 405 for unsupported resource method", func() {
		request := httptest.NewRequest(
			http.MethodPut,
			"/api/clusters_mgmt/v1/clusters",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("Returns 405 for unsupported sub-resource method", func() {
		request := httptest.NewRequest(
			http.MethodPut,
			"/api/clusters_mgmt/v1/clusters/123",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("Supports non REST method", func() {
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/clusters_mgmt/v1/register_cluster",
			strings.NewReader(`{
				"subscription_id": "123",
				"external_id": "456"
			}`),
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"cluster": {
				"kind": "Cluster",
				"id": "123",
				"external_id": "456"
			}
		}`))
	})
})

// MyServer is the implementation of the top level server.
type MyServer struct {
	clustersMgmt *MyCMServer
}

// Make sure that we implement the interface:
var _ generated.Server = &MyServer{}

func (s *MyServer) AccountsMgmt() am.Server {
	return nil
}

func (s *MyServer) Authorizations() az.Server {
	return nil
}

func (s *MyServer) ClustersMgmt() cm.Server {
	return s.clustersMgmt
}

// MyCMServer is the implementation of the clusters management server.
type MyCMServer struct {
	v1 *MyCMV1Server
}

// Make sure that we implement the interface:
var _ cm.Server = &MyCMServer{}

func (s *MyCMServer) V1() cmv1.Server {
	return s.v1
}

// MyCMV1Server is the implementation of version 1 of the clusters management server.
type MyCMV1Server struct {
	clusters *MyClustersServer
}

// Make sure that we implement the interface:
var _ cmv1.Server = &MyCMV1Server{}

func (s *MyCMV1Server) RegisterCluster(ctx context.Context,
	request *cmv1.RegisterClusterServerRequest,
	response *cmv1.RegisterClusterServerResponse) error {
	// Create the cluster:
	cluster, err := cmv1.NewCluster().
		ID(request.SubscriptionID()).
		ExternalID(request.ExternalID()).
		Build()
	Expect(err).ToNot(HaveOccurred())

	// Return the cluster:
	response.Cluster(cluster)

	return nil
}

func (s *MyCMV1Server) RegisterDisconnected(ctx context.Context,
	request *cmv1.RegisterDisconnectedServerRequest,
	response *cmv1.RegisterDisconnectedServerResponse) error {
	response.Cluster(request.Cluster())
	return nil
}

func (s *MyCMV1Server) Clusters() cmv1.ClustersServer {
	return s.clusters
}

func (s *MyCMV1Server) Nil() cmv1.NilServer {
	// This should always return nil, as it is used in the tests to check what happens when
	// a locator returns nil.
	return nil
}

// MyClustersServer is the implementation of the server that manages the collection of clusters.
type MyClustersServer struct {
	// Methods:
	list func(
		ctx context.Context,
		request *cmv1.ClustersListServerRequest,
		response *cmv1.ClustersListServerResponse,
	) error

	// Locators:
	cluster *MyClusterServer
}

func (s *MyClustersServer) List(ctx context.Context, request *cmv1.ClustersListServerRequest,
	response *cmv1.ClustersListServerResponse) error {
	return s.list(ctx, request, response)
}

func (s *MyClustersServer) Add(ctx context.Context, request *cmv1.ClustersAddServerRequest,
	response *cmv1.ClustersAddServerResponse) error {
	return nil
}

func (s *MyClustersServer) Cluster(id string) cmv1.ClusterServer {
	s.cluster.id = id
	return s.cluster
}

// MyCluster servier is the implementation of the server that manages a specific cluster.
type MyClusterServer struct {
	// Identifier:
	id string

	// Get method:
	get func(
		ctx context.Context,
		request *cmv1.ClusterGetServerRequest,
		response *cmv1.ClusterGetServerResponse,
	) error

	// Delete method:
	del func(
		ctx context.Context,
		request *cmv1.ClusterDeleteServerRequest,
		response *cmv1.ClusterDeleteServerResponse,
	) error
}

// Make sure we implement the interface:
var _ cmv1.ClusterServer = &MyClusterServer{}

func (s *MyClusterServer) Get(ctx context.Context, request *cmv1.ClusterGetServerRequest,
	response *cmv1.ClusterGetServerResponse) error {
	return s.get(ctx, request, response)
}

func (s *MyClusterServer) Update(ctx context.Context, request *cmv1.ClusterUpdateServerRequest,
	response *cmv1.ClusterUpdateServerResponse) error {
	return nil
}

func (s *MyClusterServer) Delete(ctx context.Context, request *cmv1.ClusterDeleteServerRequest,
	response *cmv1.ClusterDeleteServerResponse) error {
	return s.del(ctx, request, response)
}

func (s *MyClusterServer) Groups() cmv1.GroupsServer {
	return nil
}

func (s *MyClusterServer) IdentityProviders() cmv1.IdentityProvidersServer {
	return &MyIdentityProvidersServer{}
}

// MyIdentityProvidersServer is the implementation of the server that manages the collection of
// identity providers of a cluster.
type MyIdentityProvidersServer struct {
	// Empty on purpose.
}

// Make sure we implement the interface:
var _ cmv1.IdentityProvidersServer = &MyIdentityProvidersServer{}

func (s *MyIdentityProvidersServer) List(ctx context.Context,
	request *cmv1.IdentityProvidersListServerRequest,
	response *cmv1.IdentityProvidersListServerResponse) error {
	items, err := cmv1.NewIdentityProviderList().
		Items(cmv1.NewIdentityProvider().Name("test-list-identity-providers")).
		Build()
	if err != nil {
		return err
	}
	response.Items(items)
	response.Page(1)
	response.Size(1)
	response.Total(1)
	return nil
}

func (s *MyIdentityProvidersServer) Add(ctx context.Context,
	request *cmv1.IdentityProvidersAddServerRequest,
	response *cmv1.IdentityProvidersAddServerResponse) error {
	return nil
}

func (s *MyIdentityProvidersServer) IdentityProvider(id string) cmv1.IdentityProviderServer {
	return nil
}
