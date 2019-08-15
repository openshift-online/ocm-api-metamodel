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
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"

	v1 "gitlab.cee.redhat.com/service/ocm-api-metamodel/tests/api/clustersmgmt/v1"
)

type MyTestClustersServer struct{}

func (s *MyTestClustersServer) List(request *v1.ClustersListServerRequest,
	response *v1.ClustersListServerResponse) error {
	// Set a status code 200. Return empty response.
	response.SetStatusCode(200)
	return nil
}

func (s *MyTestClustersServer) Add(request *v1.ClustersAddServerRequest, response *v1.ClustersAddServerResponse) error {
	// Set a status code 200. Return empty response.
	response.SetStatusCode(200)
	return nil
}

func (s *MyTestClustersServer) Cluster(id string) v1.ClusterServer {
	return nil
}

var _ = Describe("Server", func() {
	It("Can recieve a request and return response", func() {
		myTestClustersServer := new(MyTestClustersServer)
		clustersAdapter := v1.NewClustersServerAdapter(myTestClustersServer)

		request := httptest.NewRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters", nil)
		recorder := httptest.NewRecorder()
		clustersAdapter.ServeHTTP(recorder, request)
	})
})
