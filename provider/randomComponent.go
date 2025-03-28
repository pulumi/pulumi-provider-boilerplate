// Copyright 2016-2023, Pulumi Corporation.
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
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Similar to resources, components have a controlling struct.
// The NewRandomComponent function is responsible for creating
// the component by composing together other resources.
type RandomComponent struct {
	pulumi.ResourceState                     // Component state needs this for tracking nested resource states.
	RandomComponentArgs                      // Include all the input fields in the state.
	Password             pulumi.StringOutput `pulumi:"password"`
}

// Similar to resources, components have an input struct, defining what arguments it accepts.
type RandomComponentArgs struct {
	Length pulumi.IntInput `pulumi:"length"`
}

func NewRandomComponent(ctx *pulumi.Context, name string, args RandomComponentArgs, opts ...pulumi.ResourceOption) (*RandomComponent, error) {
	// Initialize the component state.
	comp := &RandomComponent{
		RandomComponentArgs: args,
	}
	// Register the component resource to which we will attach all other resources.
	err := ctx.RegisterComponentResource(Name+":index:RandomComponent", name, comp, opts...)
	if err != nil {
		return nil, err
	}

	// Construct the arguments for the sub-resource.
	pArgs := &random.RandomPasswordArgs{
		Length: args.Length,
	}

	// We can access provider configuration too if needed.
	config := infer.GetConfig[Config](ctx.Context())
	if config.Scream != nil {
		pArgs.Lower = pulumi.BoolPtr(*config.Scream)
	}

	// Create the sub-resource. Ensure that the sub-resource is parented to the component resource.
	password, err := random.NewRandomPassword(ctx, name+"-password", pArgs, pulumi.Parent(comp))
	if err != nil {
		return nil, err
	}

	// Update the state of the component with output from the sub-resource.
	comp.Password = password.Result
	return comp, nil
}
