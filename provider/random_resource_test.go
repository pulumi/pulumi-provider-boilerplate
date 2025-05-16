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

package provider

import (
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	integration "github.com/pulumi/pulumi-go-provider/integration"
	presource "github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
)

func TestRandomResource(t *testing.T) {
	t.Parallel()

	server, err := integration.NewServer(t.Context(),
		"xyz",
		semver.Version{Minor: 1},
		integration.WithProvider(Provider()),
	)
	require.NoError(t, err)

	integration.LifeCycleTest{
		Resource: "xyz:index:Random",
		Create: integration.Operation{
			Inputs: presource.FromResourcePropertyMap(presource.NewPropertyMapFromMap(map[string]interface{}{
				"length": 24,
			})),

			Hook: func(inputs, output property.Map) {
				t.Logf("Outputs: %v", output)
				result := output.Get("result").AsString()
				assert.Len(t, result, 24)
			},
		},
		Updates: []integration.Operation{
			{
				Inputs: presource.FromResourcePropertyMap(presource.NewPropertyMapFromMap(map[string]interface{}{
					"length": 10,
				})),

				Hook: func(inputs, output property.Map) {
					result := output.Get("result").AsString()
					assert.Len(t, result, 10)
				},
			},
		},
	}.Run(t, server)
}
