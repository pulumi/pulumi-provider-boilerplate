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
	"time"
)

// SessionContext contains a context to manage the session for a provider.
type SessionContext struct {
	session context.Context
	cancel  context.CancelFunc
}

// NewSessionContext creates a context to manage the session for a provider. This session is created when a provider
// starts and can be used to signal for early cancellation of any ongoing operations.
func NewSessionContext() *SessionContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &SessionContext{session: ctx, cancel: cancel}
}

// Cancel signals the termination of the session context.
func (s *SessionContext) Cancel() {
	s.cancel()
}

// Join merges the cancellation context of a request and the session. This method should be called at the beginning of
// each gRPC method in a provider to ensure that cancellations are handled from either context. If timeoutSeconds
// argument is greater than 0, then the cancellation will be triggered automatically once that time has elapsed.
//
// Note that the join does not copy any values from the session context, so be sure to set values on the request
// context instead.
func (s *SessionContext) Join(ctx context.Context, timeoutSeconds int) context.Context {
	var joined context.Context
	var cancel context.CancelFunc

	if timeoutSeconds > 0 {
		joined, cancel = context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	} else {
		joined, cancel = context.WithCancel(ctx)
	}

	// Start a goroutine that blocks on the session context. If the session is cancelled, then cancel the joined
	// context.
	go func() {
		<-s.session.Done()
		cancel()
	}()

	return joined
}
