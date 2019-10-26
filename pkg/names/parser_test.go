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

// This file contains tests for the name parser.

package names

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	DescribeTable("Parse using case",
		func(text, expected string) {
			Expect(ParseUsingCase(text).String()).To(Equal(expected))
		},
		Entry("One character", "X", "x"),
		Entry("Two characters", "Pi", "pi"),
		Entry("Three characters", "Max", "max"),
		Entry("One word", "Cluster", "cluster"),
		Entry("Two words", "CloudProvider", "cloud_provider"),
		Entry("Three words", "CloudProviderRegion", "cloud_provider_region"),
		Entry("Two character initialism", "IP", "IP"),
		Entry("Two character initialism after word", "MainIP", "main_IP"),
		Entry("Two character initialism before word", "IPAddress", "IP_address"),
		Entry("Two character initialism between words", "MainIPAddress", "main_IP_address"),
		Entry("Three character initialism", "CPU", "CPU"),
		Entry("Three character initialism after word", "TotalCPU", "total_CPU"),
		Entry("Three character initialism before word", "CPUUsage", "CPU_usage"),
		Entry("Three character initialism between words", "TotalCPUUsage", "total_CPU_usage"),
		Entry("Two very short words", "PiX", "pi_x"),
		Entry("Two chars initialism plural", "IDs", "IDs"),
		Entry("Two chars initialism after word", "ClusterIDs", "cluster_IDs"),
		Entry("Two chars initialism before word", "IDsList", "IDs_list"),
		Entry("Three chars initialism plural", "CPUs", "CPUs"),
		Entry("Three chars initialism plural after word", "ClusterCPUs", "cluster_CPUs"),
		Entry("Three chars initialism plural before word", "CPUsList", "CPUs_list"),
	)
})
