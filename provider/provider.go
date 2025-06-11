// Copyright 2025, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package provider implements a simple random resource and component.
package provider

import (
	"fmt"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

// Name controls how this provider is referenced in package names and elsewhere.
const Name string = "provider-boilerplate"

// Provider creates a new instance of the provider.
func Provider() p.Provider {
	p, err := infer.NewProviderBuilder().
		WithDisplayName("pulumi-provider-boilerplate").
		WithDescription("An example built with pulumi-go-provider.").
		WithHomepage("https://www.pulumi.com").
		WithNamespace("pulumi").
		WithResources(infer.Resource(Random{})).
		WithComponents(infer.ComponentF(NewRandomComponent)).
		WithConfig(infer.Config(&Config{})).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		}).Build()
	if err != nil {
		panic(fmt.Errorf("unable to build provider: %w", err))
	}
	return p
}

// Config defines provider-level configuration
type Config struct {
	Scream *bool `pulumi:"itsasecret,optional"`
}
