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
	"fmt"
	"math/rand"
	"time"

	"github.com/jpillora/backoff"
	"github.com/pulumi/pulumi-xyz/provider/internal/errors"
	"github.com/pulumi/pulumi-xyz/provider/internal/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

type createFunction func(context.Context, chan<- resource.PropertyMap, chan<- error, resource.PropertyMap)

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

func create(ctx context.Context) createFunction {
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

func Create(ctx context.Context, inputs resource.PropertyMap) (resource.PropertyMap, error) {
	// TODO: getter -- need a way to lookup current state for a URN for checkpointing

	retryPolicy := backoff.Backoff{
		Min:    1 * time.Second,
		Max:    30 * time.Second,
		Factor: 1.5,
		Jitter: true,
	}

	createFunc := create(ctx)

	resultCh := make(chan resource.PropertyMap)
	errCh := make(chan error)

	go createFunc(ctx, resultCh, errCh, inputs)

	var result resource.PropertyMap

CREATE:
	for {
		select {
		case <-ctx.Done():
			// TODO: This should get the current state before exiting
			if ctx.Err().Error() == context.DeadlineExceeded.Error() {
				return resource.PropertyMap{}, errors.TimeoutError[string]{
					Result: "TODO",
					Err:    ctx.Err(),
				}
			} else {
				return resource.PropertyMap{}, errors.CancellationError[string]{
					Result: "TODO",
					Err:    ctx.Err(),
				}
			}
		case result = <-resultCh:
			break CREATE
		case err := <-errCh:
			logging.Info(ctx, fmt.Sprintf("create error: %s", err))

			// TODO: check for retryable error
			duration := retryPolicy.Duration()
			attempt := int(retryPolicy.Attempt())
			if attempt > 5 {
				logging.Info(ctx, "create failed: max retries exceeded")
				return resource.PropertyMap{}, err
			}
			go func() { // Launch sleep in a goroutine so we don't block
				logging.Info(ctx, fmt.Sprintf("waiting for %s before retrying create", duration))
				time.Sleep(duration)

				logging.Info(ctx, fmt.Sprintf("retrying create, attempt %d", attempt))
				go createFunc(ctx, resultCh, errCh, inputs)
			}()

			// TODO: return now if non-retryable
		}
	}

	// TODO: readiness check - if no await logic, return here
READY:
	for {
		select {
		case <-ctx.Done():
			// TODO: This should get the current state before exiting
			if ctx.Err().Error() == context.DeadlineExceeded.Error() {
				return resource.PropertyMap{}, errors.TimeoutError[string]{
					Result: "TODO",
					Err:    ctx.Err(),
				}
			} else {
				return resource.PropertyMap{}, errors.CancellationError[string]{
					Result: "TODO",
					Err:    ctx.Err(),
				}
			}
		default:
			break READY
		}
	}

	return result, nil
}

func Read(ctx context.Context, input resource.PropertyMap) (map[string]any, error) {
	return nil, nil
}

func Update(ctx context.Context, input resource.PropertyMap) (map[string]any, error) {
	return nil, nil
}

func Delete(ctx context.Context, input resource.PropertyMap) (map[string]any, error) {
	return nil, nil
}
