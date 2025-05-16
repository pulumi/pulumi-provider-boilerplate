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
	"context"
	"math/rand/v2"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// Random is the controller for the resource.
//
// Each resource has a controlling struct.
// Resource behavior is determined by implementing methods on the controlling struct.
// The `Create` method is mandatory, but other methods are optional.
// - Check: Remap inputs before they are typed.
// - Diff: Change how instances of a resource are compared.
// - Update: Mutate a resource in place.
// - Read: Get the state of a resource from the backing provider.
// - Delete: Custom logic when the resource is deleted.
// - Annotate: Describe fields and set defaults for a resource.
// - WireDependencies: Control how outputs and secrets flows through values.
type Random struct{}

// RandomArgs are the inputs to the random resource's constructor.
//
// Each resource has an input struct, defining what arguments it accepts.
type RandomArgs struct {
	// Fields projected into Pulumi must be public and hava a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but it's generally a
	// good idea.
	Length int `pulumi:"length"`
}

// RandomState is what's persisted in state.
//
// Each resource has a state, describing the fields that exist on the resource.
type RandomState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	RandomArgs
	// Here we define a required output called result.
	Result string `pulumi:"result"`
}

// Create creates a new instance of the random resource.
//
// All resources must implement Create at a minimum.
func (Random) Create(
	ctx context.Context,
	req infer.CreateRequest[RandomArgs],
) (infer.CreateResponse[RandomState], error) {
	name := req.Name
	input := req.Inputs
	preview := req.DryRun
	state := RandomState{RandomArgs: input}
	if preview {
		return infer.CreateResponse[RandomState]{ID: name, Output: state}, nil
	}
	state.Result = makeRandom(input.Length)
	return infer.CreateResponse[RandomState]{ID: name, Output: state}, nil
}

func makeRandom(length int) string {
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") // SED_SKIP

	result := make([]rune, length)
	for i := range result {
		result[i] = charset[rand.IntN(len(charset))] //nolint:gosec // Intentionally weak random.
	}
	return string(result)
}
