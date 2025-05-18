/*
Copyright 2023 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package diff

import (
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	intcmp "github.com/fluxcd/tanka-controller/internal/cmp"
)

// CompareOption is a function that modifies the unstructured object before
// comparing.
type CompareOption func(u *unstructured.Unstructured)

// WithoutStatus removes the status field from the unstructured object
// before comparing.
func WithoutStatus() CompareOption {
	return func(u *unstructured.Unstructured) {
		delete(u.Object, "status")
	}
}

// Unstructured compares two unstructured objects and returns a diff and
// a bool indicating whether the objects are equal.
func Unstructured(x, y *unstructured.Unstructured, opts ...CompareOption) (string, bool) {
	if len(opts) > 0 {
		x = x.DeepCopy()
		y = y.DeepCopy()
	}

	for _, opt := range opts {
		opt(x)
		opt(y)
	}

	r := intcmp.SimpleUnstructuredReporter{}
	_ = cmp.Diff(x.UnstructuredContent(), y.UnstructuredContent(), cmp.Reporter(&r))
	return r.String(), r.String() == ""
}
