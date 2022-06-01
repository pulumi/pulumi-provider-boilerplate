// Copyright 2016-2022, Pulumi Corporation.
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

package middleware

import (
	"context"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type CreateFunction func(
	ctx context.Context,
	resultCh chan<- resource.PropertyMap,
	errCh chan<- error,
	inputs resource.PropertyMap,
)

type ReadFunction func(
	ctx context.Context,
	resultCh chan<- resource.PropertyMap,
	errCh chan<- error,
)

type UpdateFunction func(
	ctx context.Context,
	resultCh chan<- resource.PropertyMap,
	errCh chan<- error,
	inputs resource.PropertyMap,
	old resource.PropertyMap,
)

type DeleteFunction func(
	ctx context.Context,
	resultCh chan<- resource.PropertyMap,
	errCh chan<- error,
)
