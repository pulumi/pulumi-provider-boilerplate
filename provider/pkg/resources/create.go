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

package resources

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pulumi/pulumi-xyz/provider/internal/errors"
	"github.com/pulumi/pulumi-xyz/provider/internal/logging"
	"github.com/pulumi/pulumi-xyz/provider/internal/middleware"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

func Create(ctx context.Context) middleware.CreateFunction {
	urn, ok := ctx.Value("urn").(resource.URN)
	contract.Assertf(ok, "context missing required value: urn")

	logging.Info(ctx, fmt.Sprintf("urn type: %s", urn.Type().String()))
	switch urn.Type().String() {
	case "xyz:index:Random":
		return createRandom
	default:
		return func(ctx context.Context, resultCh chan<- resource.PropertyMap, errCh chan<- error, inputs resource.PropertyMap) {
		}
	}
}

func createRandom(ctx context.Context, resultCh chan<- resource.PropertyMap, errCh chan<- error, inputs resource.PropertyMap) {
	length := int(inputs["length"].NumberValue()) // TODO: validate inputs

	logging.Info(ctx, "beginning random generation")
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	result := make([]rune, length)
	// This simulates a retry loop. A real example would include the cancellation context on any network requests.
	for i := range result {
		select {
		case <-ctx.Done():
			errCh <- errors.CancellationError[string]{Result: string(result), Err: fmt.Errorf("CANCELLED")}
			return
		default:
			logging.Info(ctx, fmt.Sprintf("creation in progress %d/%d", i, length))
			//if i == 2 {
			//	errCh <- errors.OperationError[string]{Result: string(result), Err: fmt.Errorf("Simulated")}
			//	return
			//}
			time.Sleep(500 * time.Millisecond)
		}
		result[i] = charset[seededRand.Intn(len(charset))]
	}

	logging.ClearStatus(ctx)
	resultMap := resource.NewPropertyMapFromMap(map[string]any{
		"result": string(result),
	})
	resultCh <- resultMap
}
