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

package errors

import (
	"fmt"
)

type resource interface {
	string | int | map[string]any
}

// ResourceError includes information about an error that occurred during a Pulumi update, and an object representing
// the most recent state of the corresponding resource. This allows the provider to save partial state for operations
// that failed to complete so that they can be resumed in another update.
type ResourceError[T resource] interface {
	Error() string
	Object() T
}

// CancellationError indicates that the resource operation was cancelled before it completed.
type CancellationError[T resource] struct {
	Result T
	Err    error
}

func (e CancellationError[T]) Error() string {
	return fmt.Sprint("resource operation was cancelled")
}

func (e CancellationError[T]) Object() T {
	return e.Result
}

// TimeoutError indicates that the resource operation timed out.
type TimeoutError[T resource] struct {
	Result T
	Err    error
}

func (e TimeoutError[T]) Error() string {
	return fmt.Sprint("resource operation timed out")
}

func (e TimeoutError[T]) Object() T {
	return e.Result
}

// OperationError indicates that the resource operation failed.
type OperationError[T resource] struct {
	Result T
	Err    error
}

func (e OperationError[T]) Error() string {
	return fmt.Sprintf("resource operation failed: %s", e.Err)
}

func (e OperationError[T]) Object() T {
	return e.Result
}

// ReadinessError indicates that a resource was created, but failed to become ready.
type ReadinessError[T resource] struct {
	Result T
	Err    error
}

func (e ReadinessError[T]) Error() string {
	return fmt.Sprintf("resource was created but failed to become ready: %s", e.Err)
}

func (e ReadinessError[T]) Object() T {
	return e.Result
}

// Statically verify that the error types implement the expected interfaces.
var _ error = (*CancellationError[string])(nil)
var _ ResourceError[string] = (*CancellationError[string])(nil)
var _ error = (*TimeoutError[string])(nil)
var _ ResourceError[string] = (*TimeoutError[string])(nil)
var _ error = (*OperationError[string])(nil)
var _ ResourceError[string] = (*OperationError[string])(nil)
var _ error = (*ReadinessError[string])(nil)
var _ ResourceError[string] = (*ReadinessError[string])(nil)
