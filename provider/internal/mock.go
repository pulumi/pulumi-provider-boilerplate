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

package internal

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pulumi/pulumi-xyz/provider/pkg/errors"
	"github.com/pulumi/pulumi-xyz/provider/pkg/logging"
)

func genRandom(
	ctx context.Context,
	resultCh chan<- string,
	errCh chan<- errors.ResourceError[string],
	length int,
) {
	logging.Info(ctx, "beginning random generation")
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	result := make([]rune, length)
	for i := range result {
		select {
		case <-ctx.Done():
			errCh <- errors.CancellationError[string]{Result: string(result), Err: fmt.Errorf("CANCELLED")}
			return
		default:
			logging.Info(ctx, fmt.Sprintf("creation in progress %d/5", i))
			time.Sleep(1 * time.Second)
		}
		result[i] = charset[seededRand.Intn(len(charset))]
	}

	logging.ClearStatus(ctx)
	resultCh <- string(result)
}

/*
Want:

1. Consistent cancellation handling
2. Retries
3. Consistent partial state checkpointing
4. Optional dispatch to awaiter
5. Consistent error handling

TODO:

1. Allow author to separate the operation from awaiting ready
2. Where to put retry logic? Should this be specific to await logic?
 1. Want transport to automatically refresh / handle transient network error
3. Get current state on cancel

*/

func MakeRandom(ctx context.Context, length int) (map[string]any, errors.ResourceError[string]) {
	resultCh := make(chan string)
	errCh := make(chan errors.ResourceError[string])

	// TODO: retry logic here?
	go genRandom(ctx, resultCh, errCh, length)

	select {
	case <-ctx.Done():
		// TODO: This should get the current state before exiting
		if ctx.Err().Error() == "context deadline exceeded" {
			return map[string]any{}, errors.TimeoutError[string]{
				Result: "TODO",
				Err:    ctx.Err(),
			}
		} else {
			return map[string]any{}, errors.CancellationError[string]{
				Result: "TODO",
				Err:    ctx.Err(),
			}
		}
	case r := <-resultCh:
		return map[string]any{"result": r}, nil
	case err := <-errCh:
		// TODO: This should get the current state before exiting
		return map[string]any{}, err
	}
}
