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

package nomenclator

import (
	"gitlab.cee.redhat.com/service/ocm-api-metamodel/pkg/names"
)

var (
	// A:
	Add = names.ParseUsingCase("Add")

	// B:
	Boolean = names.ParseUsingCase("Boolean")
	Builder = names.ParseUsingCase("Builder")

	// C:
	Client  = names.ParseUsingCase("Client")
	Clients = names.ParseUsingCase("Clients")

	// D:
	Data   = names.ParseUsingCase("Data")
	Date   = names.ParseUsingCase("Date")
	Delete = names.ParseUsingCase("Delete")

	// E:
	Error  = names.ParseUsingCase("Error")
	Errors = names.ParseUsingCase("Errors")

	// F:
	Float = names.ParseUsingCase("Float")

	// G:
	Get = names.ParseUsingCase("Get")

	// H:
	HREF    = names.ParseUsingCase("HREF")
	Handler = names.ParseUsingCase("Handler")
	Helpers = names.ParseUsingCase("Helpers")

	// I:
	ID      = names.ParseUsingCase("ID")
	Index   = names.ParseUsingCase("Index")
	Integer = names.ParseUsingCase("Integer")

	// J:
	Kind = names.ParseUsingCase("Kind")

	// L:
	Link = names.ParseUsingCase("Link")
	List = names.ParseUsingCase("List")
	Long = names.ParseUsingCase("Long")

	// M:
	Map     = names.ParseUsingCase("Map")
	Marshal = names.ParseUsingCase("Marshal")

	// N:
	New = names.ParseUsingCase("New")

	// P:
	Page  = names.ParseUsingCase("Page")
	Parse = names.ParseUsingCase("Parse")
	Post  = names.ParseUsingCase("Post")

	// R:
	Read     = names.ParseUsingCase("Read")
	Reader   = names.ParseUsingCase("Reader")
	Readers  = names.ParseUsingCase("Readers")
	Request  = names.ParseUsingCase("Request")
	Resource = names.ParseUsingCase("Resource")
	Response = names.ParseUsingCase("Response")
	Root     = names.ParseUsingCase("Root")
	Router   = names.ParseUsingCase("Router")

	// S:
	Service = names.ParseUsingCase("Service")
	Set     = names.ParseUsingCase("Set")
	Spec    = names.ParseUsingCase("Spec")
	String  = names.ParseUsingCase("String")
	Server  = names.ParseUsingCase("Server")

	// T:
	Type = names.ParseUsingCase("Type")

	// U:
	Unmarshal = names.ParseUsingCase("Unmarshal")
	Unwrap    = names.ParseUsingCase("Unwrap")
	Update    = names.ParseUsingCase("Update")

	// W:
	Wrap = names.ParseUsingCase("Wrap")
)
