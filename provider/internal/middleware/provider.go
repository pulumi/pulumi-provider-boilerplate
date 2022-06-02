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
	"time"

	"github.com/jpillora/backoff"
	"github.com/pulumi/pulumi-xyz/provider/internal/errors"
	"github.com/pulumi/pulumi-xyz/provider/internal/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func Create(ctx context.Context, fn CreateFunction, inputs resource.PropertyMap) (resource.PropertyMap, error) {
	// TODO: getter -- need a way to lookup current state for a URN for checkpointing

	retryPolicy := backoff.Backoff{
		Min:    1 * time.Second,
		Max:    30 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	resultCh := make(chan resource.PropertyMap)
	errCh := make(chan error)

	go fn(ctx, resultCh, errCh, inputs)

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
				fn(ctx, resultCh, errCh, inputs)
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

// TODO: Read
func Read(ctx context.Context, fn ReadFunction) (map[string]any, error) {
	return nil, nil
}

// TODO: Update
func Update(ctx context.Context, fn UpdateFunction, inputs, olds resource.PropertyMap) (map[string]any, error) {
	return nil, nil
}

// TODO: Delete
func Delete(ctx context.Context, fn DeleteFunction, inputs resource.PropertyMap) (map[string]any, error) {
	return nil, nil
}
