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

	"github.com/openshift-online/ocm-api-metamodel/tests/api"
	am "github.com/openshift-online/ocm-api-metamodel/tests/api/accountsmgmt"
	az "github.com/openshift-online/ocm-api-metamodel/tests/api/authorizations"
	cm "github.com/openshift-online/ocm-api-metamodel/tests/api/clustersmgmt"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/clustersmgmt/v1"
)

var _ = Describe("Server", func() {
	var adapter *api.Adapter
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		adapter = api.NewAdapter(&MyServer{})
		recorder = httptest.NewRecorder()
	})

	It("Can receive a request and return response", func() {
		request := httptest.NewRequest(http.MethodGet, "/clusters_mgmt/v1/clusters", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
	})

	It("Returns the list of clusters with a trailing slash", func() {
		request := httptest.NewRequest(http.MethodGet, "/clusters_mgmt/v1/clusters/", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
	})

	It("Returns a 404 for an unknown resource", func() {
		request := httptest.NewRequest(http.MethodGet, "/foo", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Can get a list of clusters", func() {
		request := httptest.NewRequest(http.MethodGet, "/clusters_mgmt/v1/clusters", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"page": 0,
			"size": 0,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "test-list-clusters"
				}
			]
		}`))
	})

	It("Can get a list of clusters by page", func() {
		request := httptest.NewRequest(
			http.MethodGet,
			"/clusters_mgmt/v1/clusters?page=2",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"page": 2,
			"size": 0,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "test-list-clusters"
				}
			]
		}`))
	})

	It("Can get a list of clusters by size", func() {
		request := httptest.NewRequest(
			http.MethodGet,
			"/clusters_mgmt/v1/clusters?size=2",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"page": 0,
			"size": 2,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "test-list-clusters"
				}
			]
		}`))
	})

	It("Can get a list of clusters by size and page", func() {
		request := httptest.NewRequest(
			http.MethodGet,
			"/clusters_mgmt/v1/clusters?size=2&page=1",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"page": 1,
			"size": 2,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "test-list-clusters"
				}
			]
		}`))
	})

	It("Can get a cluster by identifier", func() {
		request := httptest.NewRequest(
			http.MethodGet,
			"/clusters_mgmt/v1/clusters/123",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"kind": "Cluster",
			"name": "test-get-cluster-by-id"
		}`))
	})

	It("Can get a cluster sub resource by id", func() {
		request := httptest.NewRequest(
			http.MethodGet,
			"/clusters_mgmt/v1/clusters/123/identity_providers",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
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
		request := httptest.NewRequest(http.MethodGet, "/foo_mgmt", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 404 for an unknown version", func() {
		request := httptest.NewRequest(http.MethodGet, "/clusters_mgmt/v100", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 404 for an unknown resource", func() {
		request := httptest.NewRequest(http.MethodGet, "/clusters_mgmt/flusters", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 404 for an unknown sub-resource", func() {
		request := httptest.NewRequest(
			http.MethodGet,
			"/clusters_mgmt/v1/clusters/123/foo",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 404 if the server returns nil for a locator", func() {
		request := httptest.NewRequest(http.MethodGet, "/clusters_mgmt/v1/nil", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("Returns 405 for unsupported service method", func() {
		request := httptest.NewRequest(http.MethodPost, "/clusters_mgmt", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("Returns 405 for unsupported version method", func() {
		request := httptest.NewRequest(http.MethodPost, "/clusters_mgmt/v1", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("Returns 405 for unsupported resource method", func() {
		request := httptest.NewRequest(http.MethodPut, "/clusters_mgmt/v1/clusters", nil)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("Returns 405 for unsupported sub-resource method", func() {
		request := httptest.NewRequest(
			http.MethodPut,
			"/clusters_mgmt/v1/clusters/123",
			nil,
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("Supports non REST method", func() {
		request := httptest.NewRequest(
			http.MethodPost,
			"/clusters_mgmt/v1/register_cluster",
			strings.NewReader(`{
				"id": "123",
				"name": "mycluster",
				"external_id": "456"
			}`),
		)
		adapter.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"kind": "Cluster",
			"id": "123",
			"name": "mycluster",
			"external_id": "456"
		}`))
	})
})

// MyServer is the implementation of the top level server.
type MyServer struct {
	// Empty on purpose.
}

// Make sure that we implement the interface:
var _ api.Server = &MyServer{}

func (s *MyServer) AccountsMgmt() am.Server {
	return nil
}

func (s *MyServer) Authorizations() az.Server {
	return nil
}

func (s *MyServer) ClustersMgmt() cm.Server {
	return &MyCMServer{}
}

// MyCMServer is the implementation of the clusters management server.
type MyCMServer struct {
	// Empty on purpose.
}

// Make sure that we implement the interface:
var _ cm.Server = &MyCMServer{}

func (s *MyCMServer) V1() cmv1.Server {
	return &MyCMV1Server{}
}

// MyCMV1Server is the implementation of version 1 of the clusters management server.
type MyCMV1Server struct {
	// Empty on purpose.
}

// Make sure that we implement the interface:
var _ cmv1.Server = &MyCMV1Server{}

func (s *MyCMV1Server) RegisterCluster(ctx context.Context,
	request *cmv1.RootRegisterClusterServerRequest,
	response *cmv1.RootRegisterClusterServerResponse) error {
	response.Body(request.Body())
	return nil
}

func (s *MyCMV1Server) RegisterDisconnected(ctx context.Context,
	request *cmv1.RootRegisterDisconnectedServerRequest,
	response *cmv1.RootRegisterDisconnectedServerResponse) error {
	response.Body(request.Body())
	return nil
}

func (s *MyCMV1Server) Clusters() cmv1.ClustersServer {
	return &MyClustersServer{}
}

func (s *MyCMV1Server) Nil() cmv1.NilServer {
	// This should always return nil, as it is used in the tests to check what happens when
	// a locator returns nil.
	return nil
}

// MyClustersServer is the implementation of the server that manages the collection of clusters.
type MyClustersServer struct {
	// Empty on purpose.
}

func (s *MyClustersServer) List(ctx context.Context, request *cmv1.ClustersListServerRequest,
	response *cmv1.ClustersListServerResponse) error {
	items, err := cmv1.NewClusterList().
		Items(cmv1.NewCluster().Name("test-list-clusters")).
		Build()
	if err != nil {
		return err
	}
	response.Items(items)
	response.Page(request.Page())
	response.Size(request.Size())
	response.Total(items.Len())
	return nil
}

func (s *MyClustersServer) Add(ctx context.Context, request *cmv1.ClustersAddServerRequest,
	response *cmv1.ClustersAddServerResponse) error {
	return nil
}

func (s *MyClustersServer) Cluster(id string) cmv1.ClusterServer {
	return &MyClusterServer{}
}

// MyCluster servier is the implementation of the server that manages a specific cluster.
type MyClusterServer struct {
	// Empty on purpose.
}

// Make sure we implement the interface:
var _ cmv1.ClusterServer = &MyClusterServer{}

func (s *MyClusterServer) Get(ctx context.Context, request *cmv1.ClusterGetServerRequest,
	response *cmv1.ClusterGetServerResponse) error {
	cluster, err := cmv1.NewCluster().Name("test-get-cluster-by-id").Build()
	if err != nil {
		return err
	}
	response.Body(cluster)
	return nil
}

func (s *MyClusterServer) Update(ctx context.Context, request *cmv1.ClusterUpdateServerRequest,
	response *cmv1.ClusterUpdateServerResponse) error {
	return nil
}

func (s *MyClusterServer) Delete(ctx context.Context, request *cmv1.ClusterDeleteServerRequest,
	response *cmv1.ClusterDeleteServerResponse) error {
	return nil
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
