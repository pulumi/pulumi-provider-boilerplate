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

	"github.com/pulumi/pulumi-xyz/provider/internal/errors"
	"github.com/pulumi/pulumi-xyz/provider/internal/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

/*
CRUD method
- idempotent
- optionally handle context cancellation
- transparently handle credential refresh here
*/

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
	// This simulates a retry loop. A real example would include the cancellation context on any network requests.
	for i := range result {
		select {
		case <-ctx.Done():
			errCh <- errors.CancellationError[string]{Result: string(result), Err: fmt.Errorf("CANCELLED")}
			return
		default:
			logging.Info(ctx, fmt.Sprintf("creation in progress %d/%d", i, length))
			time.Sleep(500 * time.Millisecond)
		}
		result[i] = charset[seededRand.Intn(len(charset))]
	}

	logging.ClearStatus(ctx)
	resultCh <- string(result)
}

/*
Middleware params
- context (includes URN, host)
- arg bag (resource.PropertyMap)
- optional retry config -- TODO: maybe not needed
*/

// TODO: retry args needed? maybe just default to something sensible and set overall timeout using context
//type RetryArgs struct {
//	Count int
//}

func getState(ctx context.Context) {

}

func Middleware(ctx context.Context, input resource.PropertyMap) (map[string]any, error) {
	return nil, nil
}

func MakeRandom(ctx context.Context, length int) (map[string]any, error) {
	resultCh := make(chan string)
	errCh := make(chan errors.ResourceError[string])

	// TODO: getter -- need a way to lookup current state for a URN for checkpointing

	// TODO: configure with retry args -- default to something sensible if not provided (no retry?)
	for i := 0; i < 3; i++ {

	}
	// TODO: this function should be looked up from a provider-specific mapping. probably use a function to abstract the lookup by URN
	go genRandom(ctx, resultCh, errCh, length)

	// TODO: readiness check. same as above, look up from function by URN

	// (context, result, err) could be common args for CRUD + readiness functions

	select {
	case <-ctx.Done():
		// TODO: This should get the current state before exiting
		if ctx.Err().Error() == context.DeadlineExceeded.Error() {
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
