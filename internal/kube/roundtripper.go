/*
Copyright The Helm Authors.

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

// TODO: this file was taken from https://github.com/helm/helm/blob/main/pkg/cli/roundtripper.go
// We should be able to get rid of it once https://github.com/helm/helm/issues/13052 is addressed
// A PR (https://github.com/helm/helm/pull/13383) was open and merged but not yet picked into any release
package kube

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type retryingRoundTripper struct {
	wrapped http.RoundTripper
}

func (rt *retryingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.roundTrip(req, 1, nil)
}

func (rt *retryingRoundTripper) roundTrip(req *http.Request, retry int, prevResp *http.Response) (*http.Response, error) {
	if retry < 0 {
		return prevResp, nil
	}
	resp, rtErr := rt.wrapped.RoundTrip(req)
	if rtErr != nil {
		return resp, rtErr
	}
	if resp.StatusCode < 500 {
		return resp, nil
	}
	if resp.Header.Get("content-type") != "application/json" {
		return resp, nil
	}
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return resp, err
	}

	var ke kubernetesError
	r := bytes.NewReader(b)
	err = json.NewDecoder(r).Decode(&ke)
	r.Seek(0, io.SeekStart)
	resp.Body = io.NopCloser(r)
	if err != nil {
		return resp, err
	}
	if ke.Code < 500 {
		return resp, nil
	}
	// Matches messages like "etcdserver: leader changed"
	if strings.HasSuffix(ke.Message, "etcdserver: leader changed") {
		return rt.roundTrip(req, retry-1, resp)
	}
	// Matches messages like "rpc error: code = Unknown desc = raft proposal dropped"
	if strings.HasSuffix(ke.Message, "raft proposal dropped") {
		return rt.roundTrip(req, retry-1, resp)
	}
	return resp, nil
}

type kubernetesError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
