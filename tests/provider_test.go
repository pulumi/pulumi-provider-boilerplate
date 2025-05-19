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

package tests

import (
	"context"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	xyz "github.com/pulumi/pulumi-provider-boilerplate/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
)

func TestRandomCreate(t *testing.T) {
	t.Parallel()

	prov := provider(t)

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("Random"),
		Properties: property.NewMap(map[string]property.Value{
			"length": property.New(12.0),
		}),

		DryRun: false,
	})

	require.NoError(t, err)
	result := response.Properties.Get("result").AsString()
	assert.Len(t, result, 12)
}

// urn is a helper function to build an urn for running integration tests.
func urn(typ string) resource.URN {
	return resource.NewURN("stack", "proj", "",
		tokens.Type("test:index:"+typ), "name")
}

// Create a test server.
func provider(t *testing.T) integration.Server {
	s, err := integration.NewServer(
		context.Background(),
		xyz.Name,
		semver.MustParse("1.0.0"),
		integration.WithProvider(xyz.Provider()),
	)
	require.NoError(t, err)
	return s
}
