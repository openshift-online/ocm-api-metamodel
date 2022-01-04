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

// This file contains the definition of the test suite.

package tests

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main")
}

// NewTransport creates an instance of the transport that will be used by the tests to talk to the
// tests server.
func NewTransport(server *Server) http.RoundTripper {
	return &Transport{
		server:  server,
		wrapped: &http.Transport{},
	}
}

// Transport is the transport that will be used by the tests to talk to the tests server. It takes
// care of basic things like adding the server address to the path calculated by the client.
type Transport struct {
	server  *Server
	wrapped *http.Transport
}

// RoundTrip implements the RoundTripper interface.
func (t *Transport) RoundTrip(request *http.Request) (response *http.Response, err error) {
	request.URL.Scheme = "http"
	request.URL.Host = t.server.Addr()
	request.Header.Set("Content-type", "application/json")
	response, err = t.wrapped.RoundTrip(request)
	return
}
