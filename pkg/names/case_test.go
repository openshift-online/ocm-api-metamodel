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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable(
	"Convert to lower camel case",
	func(text, expected string) {
		Expect(ToLowerCamel(text)).To(Equal(expected))
	},
	Entry("One character", "X", "x"),
	Entry("Two characters", "Pi", "pi"),
	Entry("Three characters", "Max", "max"),
	Entry("One word", "Cluster", "cluster"),
	Entry("Two words", "CloudProvider", "cloudProvider"),
	Entry("Three words", "CloudProviderRegion", "cloudProviderRegion"),
	Entry("Two character initialism", "IP", "ip"),
	Entry("Two character initialism after word", "MainIP", "mainIP"),
	Entry("Two character initialism before word", "IPAddress", "ipAddress"),
	Entry("Two character initialism between words", "MainIPAddress", "mainIPAddress"),
	Entry("Three character initialism", "CPU", "cpu"),
	Entry("Three character initialism after word", "TotalCPU", "totalCPU"),
	Entry("Three character initialism before word", "CPUUsage", "cpuUsage"),
	Entry("Three character initialism between words", "TotalCPUUsage", "totalCPUUsage"),
	Entry("Two very short words", "PiX", "piX"),
	Entry("Two chars initialism plural", "IDs", "ids"),
	Entry("Two chars initialism after word", "ClusterIDs", "clusterIDs"),
	Entry("Two chars initialism before word", "IDsList", "idsList"),
	Entry("Three chars initialism plural", "CPUs", "cpus"),
	Entry("Three chars initialism plural after word", "ClusterCPUs", "clusterCPUs"),
	Entry("Three chars initialism plural before word", "CPUsList", "cpusList"),
	Entry("Words with numbers", "ClientX509CertURL", "clientX509CertURL"),
	Entry("Version number", "v1", "v1"),
)

var _ = DescribeTable(
	"Convert to upper camel case",
	func(text, expected string) {
		Expect(ToUpperCamel(text)).To(Equal(expected))
	},
	Entry("One character", "X", "X"),
	Entry("Two characters", "Pi", "Pi"),
	Entry("Three characters", "Max", "Max"),
	Entry("One word", "Cluster", "Cluster"),
	Entry("Two words", "CloudProvider", "CloudProvider"),
	Entry("Three words", "CloudProviderRegion", "CloudProviderRegion"),
	Entry("Two character initialism", "IP", "IP"),
	Entry("Two character initialism after word", "MainIP", "MainIP"),
	Entry("Two character initialism before word", "IPAddress", "IPAddress"),
	Entry("Two character initialism between words", "MainIPAddress", "MainIPAddress"),
	Entry("Three character initialism", "CPU", "CPU"),
	Entry("Three character initialism after word", "TotalCPU", "TotalCPU"),
	Entry("Three character initialism before word", "CPUUsage", "CPUUsage"),
	Entry("Three character initialism between words", "TotalCPUUsage", "TotalCPUUsage"),
	Entry("Two very short words", "PiX", "PiX"),
	Entry("Two chars initialism plural", "IDs", "IDs"),
	Entry("Two chars initialism after word", "ClusterIDs", "ClusterIDs"),
	Entry("Two chars initialism before word", "IDsList", "IDsList"),
	Entry("Three chars initialism plural", "CPUs", "CPUs"),
	Entry("Three chars initialism plural after word", "ClusterCPUs", "ClusterCPUs"),
	Entry("Three chars initialism plural before word", "CPUsList", "CPUsList"),
	Entry("Words with numbers", "ClientX509CertURL", "ClientX509CertURL"),
	Entry("Version number", "v1", "V1"),
)

var _ = DescribeTable(
	"Convert to snake case",
	func(text, expected string) {
		Expect(ToSnake(text)).To(Equal(expected))
	},
	Entry("One character", "X", "x"),
	Entry("Two characters", "Pi", "pi"),
	Entry("Three characters", "Max", "max"),
	Entry("One word", "Cluster", "cluster"),
	Entry("Two words", "CloudProvider", "cloud_provider"),
	Entry("Three words", "CloudProviderRegion", "cloud_provider_region"),
	Entry("Two character initialism", "IP", "ip"),
	Entry("Two character initialism after word", "MainIP", "main_ip"),
	Entry("Two character initialism before word", "IPAddress", "ip_address"),
	Entry("Two character initialism between words", "MainIPAddress", "main_ip_address"),
	Entry("Three character initialism", "CPU", "cpu"),
	Entry("Three character initialism after word", "TotalCPU", "total_cpu"),
	Entry("Three character initialism before word", "CPUUsage", "cpu_usage"),
	Entry("Three character initialism between words", "TotalCPUUsage", "total_cpu_usage"),
	Entry("Two very short words", "PiX", "pi_x"),
	Entry("Two chars initialism plural", "IDs", "ids"),
	Entry("Two chars initialism after word", "ClusterIDs", "cluster_ids"),
	Entry("Two chars initialism before word", "IDsList", "ids_list"),
	Entry("Three chars initialism plural", "CPUs", "cpus"),
	Entry("Three chars initialism plural after word", "ClusterCPUs", "cluster_cpus"),
	Entry("Three chars initialism plural before word", "CPUsList", "cpus_list"),
	Entry("Words with numbers", "ClientX509CertURL", "client_x509_cert_url"),
	Entry("Version number", "v1", "v1"),
)
