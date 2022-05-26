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

	"github.com/pulumi/pulumi-xyz/provider/pkg/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

func MakeRandom(ctx context.Context, length int) string {
	done := make(chan string)
	defer close(done)

	// TODO: pull this logic into a separate function and generalize the cancellation wrapper
	// 		 1. need channel as arg
	//		 2. need to return partial state on failure
	go func() {
		logging.Log(ctx, diag.Info, "beginning random generation")
		seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
		charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

		result := make([]rune, length)
		for i := range result {
			result[i] = charset[seededRand.Intn(len(charset))]
		}
		for i := 0; i <= 5; i++ {
			logging.Log(ctx, diag.Info, fmt.Sprintf("creation in progress %d/5", i))
			time.Sleep(1 * time.Second)
		}
		logging.ClearStatus(ctx)
		done <- string(result)
	}()

	select {
	case <-ctx.Done():
		return "CANCELLED"
	case r := <-done:
		return r
	}
}
