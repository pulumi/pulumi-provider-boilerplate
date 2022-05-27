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

package await

import (
	"fmt"
)

type foo interface {
	string | map[string]any
}

type ResourceError[T foo] interface {
	Error() string
	Object() T
}

type CancellationError[T foo] struct {
	Result T
	Err    error
}

func (e CancellationError[T]) Error() string {
	return fmt.Sprint("resource operation was cancelled")
}

func (e CancellationError[T]) Object() T {
	return e.Result
}

type TimeoutError[T foo] struct {
	Result T
	Err    error
}

func (e TimeoutError[T]) Error() string {
	return fmt.Sprint("resource operation timed out")
}

func (e TimeoutError[T]) Object() T {
	return e.Result
}

//type PartialStringError struct {
//	Result string
//	Err    error
//}
//
//func (e PartialStringError) Error() string {
//	return e.Err.Error()
//}
//
//func (e PartialStringError) Object() any {
//	return e.Result
//}
