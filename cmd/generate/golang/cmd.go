/*
Copyright (c) 2020 Red Hat, Inc.

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

package golang

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/openshift-online/ocm-api-metamodel/pkg/generators"
	"github.com/openshift-online/ocm-api-metamodel/pkg/generators/golang"
	"github.com/openshift-online/ocm-api-metamodel/pkg/generators/openapi"
	"github.com/openshift-online/ocm-api-metamodel/pkg/http"
	"github.com/openshift-online/ocm-api-metamodel/pkg/language"
	"github.com/openshift-online/ocm-api-metamodel/pkg/reporter"
)

// Cmd is the definition of the command:
var Cmd = &cobra.Command{
	Use:   "go",
	Short: "Generate Go code",
	Long:  "Generate Go code.",
	Run:   run,
}

// Values of the command line arguments:
var args struct {
	paths  []string
	base   string
	output string
}

func init() {
	flags := Cmd.Flags()
	flags.StringSliceVar(
		&args.paths,
		"model",
		[]string{},
		"File or directory containing the model. If it is a directory then all .model"+
			"files inside it and its sub directories will be loaded. If used "+
			"multiple times then all the specified files and directories will be "+
			"loaded, in the same order that they appear in the command line.",
	)
	flags.StringVar(
		&args.base,
		"base",
		"",
		"Import path of the base package for the generated code.",
	)
	flags.StringVar(
		&args.output,
		"output",
		"",
		"Directory where the source code will be generated.",
	)
}

func run(cmd *cobra.Command, argv []string) {
	// Create the reporter:
	reporter := reporter.NewReporter()

	// Check command line options:
	ok := true
	if len(args.paths) == 0 {
		reporter.Errorf("Option '--model' is mandatory")
		ok = false
	}
	if args.base == "" {
		reporter.Errorf("Option '--base' is mandatory")
		ok = false
	}
	if args.output == "" {
		reporter.Errorf("Option '--output' is mandatory")
		ok = false
	}
	if !ok {
		os.Exit(1)
	}

	// Read the model:
	model, err := language.NewReader().
		Reporter(reporter).
		Inputs(args.paths).
		Read()
	if err != nil {
		reporter.Errorf("Can't read model: %v", err)
		os.Exit(1)
	}

	// Create the calculators:
	goPackagesCalculator, err := golang.NewPackagesCalculator().
		Reporter(reporter).
		Base(args.base).
		Build()
	if err != nil {
		reporter.Errorf("Can't create Go packages calculator: %v", err)
		os.Exit(1)
	}
	goNamesCalculator, err := golang.NewNamesCalculator().
		Reporter(reporter).
		Build()
	if err != nil {
		reporter.Errorf("Can't create Go names calculator: %v", err)
		os.Exit(1)
	}
	goTypesCalculator, err := golang.NewTypesCalculator().
		Reporter(reporter).
		Packages(goPackagesCalculator).
		Names(goNamesCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create Go types calculator: %v", err)
		os.Exit(1)
	}
	bindingCalculator, err := http.NewBindingCalculator().
		Reporter(reporter).
		Build()
	if err != nil {
		reporter.Errorf("Can't create HTTP binding calculator: %v", err)
		os.Exit(1)
	}
	openapiNamesCalculator, err := openapi.NewNamesCalculator().
		Reporter(reporter).
		Build()
	if err != nil {
		reporter.Errorf("Can't create OpenAPI names calculator: %v", err)
		os.Exit(1)
	}

	// We will store here all the code generators that we will later run:
	var gens []generators.Generator
	var gen generators.Generator

	// Create the errors generator:
	gen, err = golang.NewErrorsGenerator().
		Reporter(reporter).
		Model(model).
		Output(args.output).
		Packages(goPackagesCalculator).
		Names(goNamesCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create errors generator: %v", err)
		os.Exit(1)
	}
	gens = append(gens, gen)

	// Create the helpers generator:
	gen, err = golang.NewHelpersGenerator().
		Reporter(reporter).
		Model(model).
		Output(args.output).
		Packages(goPackagesCalculator).
		Names(goNamesCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create helpers generator: %v", err)
		os.Exit(1)
	}
	gens = append(gens, gen)

	// Create the types generator:
	gen, err = golang.NewTypesGenerator().
		Reporter(reporter).
		Model(model).
		Output(args.output).
		Packages(goPackagesCalculator).
		Names(goNamesCalculator).
		Types(goTypesCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create types generator: %v", err)
		os.Exit(1)
	}
	gens = append(gens, gen)

	// Create the builders generator:
	gen, err = golang.NewBuildersGenerator().
		Reporter(reporter).
		Model(model).
		Output(args.output).
		Packages(goPackagesCalculator).
		Names(goNamesCalculator).
		Types(goTypesCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create builders generator: %v", err)
		os.Exit(1)
	}
	gens = append(gens, gen)

	// Create the clients generator:
	gen, err = golang.NewClientsGenerator().
		Reporter(reporter).
		Model(model).
		Output(args.output).
		Packages(goPackagesCalculator).
		Names(goNamesCalculator).
		Types(goTypesCalculator).
		Binding(bindingCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create clients generator: %v", err)
		os.Exit(1)
	}
	gens = append(gens, gen)

	// Create the resource generator:
	gen, err = golang.NewServersGenerator().
		Reporter(reporter).
		Model(model).
		Output(args.output).
		Packages(goPackagesCalculator).
		Names(goNamesCalculator).
		Types(goTypesCalculator).
		Binding(bindingCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create servers generator: %v", err)
		os.Exit(1)
	}
	gens = append(gens, gen)

	// Create the JSON readers generator:
	gen, err = golang.NewJSONSupportGenerator().
		Reporter(reporter).
		Model(model).
		Output(args.output).
		Packages(goPackagesCalculator).
		Names(goNamesCalculator).
		Types(goTypesCalculator).
		Binding(bindingCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create JSON readers generator: %v", err)
		os.Exit(1)
	}
	gens = append(gens, gen)

	// Create the OpenAPI specifications generator:
	gen, err = golang.NewOpenAPIGenerator().
		Reporter(reporter).
		Model(model).
		Output(args.output).
		Packages(goPackagesCalculator).
		Names(openapiNamesCalculator).
		Binding(bindingCalculator).
		Build()
	if err != nil {
		reporter.Errorf("Can't create OpenAPI specifications generator: %v", err)
		os.Exit(1)
	}
	gens = append(gens, gen)

	// Run the generators:
	for _, gen := range gens {
		err = gen.Run()
		if err != nil {
			reporter.Errorf("Generation failed: %v", err)
			os.Exit(1)
		}
	}

	// Run the formatter:
	reporter.Infof("Formating generated files")
	sources := []string{}
	err = filepath.Walk(args.output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() && filepath.Ext(path) == ".go" {
			sources = append(sources, path)
		}
		return nil
	})
	if err != nil {
		reporter.Errorf("Can't collect generated sources: %v", err)
		os.Exit(1)
	}
	sort.Strings(sources)
	args := []string{
		"-w",
	}
	args = append(args, sources...)
	formatter := exec.Command("goimports", args...)
	formatter.Stdout = os.Stdout
	formatter.Stderr = os.Stderr
	err = formatter.Run()
	if err != nil {
		reporter.Errorf("Can't format generated sources: %v", err)
		os.Exit(1)
	}

	// Bye:
	os.Exit(0)
}
