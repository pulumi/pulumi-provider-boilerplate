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

package session

import (
	"context"
	"time"
)

// Context contains a session context for a provider.
type Context struct {
	session context.Context
	cancel  context.CancelFunc
}

// NewContext creates a session context for a provider. This session is created when a provider starts and can be used
// to signal for early cancellation of any ongoing operations.
func NewContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{session: ctx, cancel: cancel}
}

// Cancel signals the termination of the session context.
func (s *Context) Cancel() {
	s.cancel()
}

// Join merges the cancellation context of a request and the session. This method should be called at the beginning of
// each gRPC method in a provider to ensure that cancellations are handled from either context. If timeoutSeconds
// argument is greater than 0, then the cancellation will be triggered automatically once that time has elapsed.
func (s *Context) Join(ctx context.Context, timeoutSeconds int) (context.Context, context.CancelFunc) {
	var joined context.Context
	var cancel context.CancelFunc

	if timeoutSeconds > 0 {
		joined, cancel = context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	} else { // TODO: set a sensible default timeout (30 mins?)
		joined, cancel = context.WithCancel(ctx)
	}

	// Start a goroutine that blocks on the session context. If the session is cancelled, then cancel the joined
	// context. The goroutine exits once the joined context completes.
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-s.session.Done():
			cancel()
		}
	}()

	return joined, cancel
}
