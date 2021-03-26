// Copyright 2016-2021 The Libsacloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sacloud

import (
	"log"
	"net/http/httputil"

	"go.uber.org/ratelimit"

	"net/http"
	"sync"
)

// RateLimitRoundTripper 秒間アクセス数を制限するためのhttp.RoundTripper実装
type RateLimitRoundTripper struct {
	// Transport 親となるhttp.RoundTripper、nilの場合http.DefaultTransportが利用される
	Transport http.RoundTripper
	// RateLimitPerSec 秒あたりのリクエスト数
	RateLimitPerSec int

	once      sync.Once
	rateLimit ratelimit.Limiter
}

// RoundTrip http.RoundTripperの実装
func (r *RateLimitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.once.Do(func() {
		r.rateLimit = ratelimit.New(r.RateLimitPerSec)
	})
	if r.Transport == nil {
		r.Transport = http.DefaultTransport
	}

	r.rateLimit.Take()
	return r.Transport.RoundTrip(req)
}

// TracingRoundTripper リクエスト/レスポンスのトレースログを出力するためのhttp.RoundTripper実装
type TracingRoundTripper struct {
	// Transport 親となるhttp.RoundTripper、nilの場合http.DefaultTransportが利用される
	Transport http.RoundTripper
}

// RoundTrip http.RoundTripperの実装
func (r *TracingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.Transport == nil {
		r.Transport = http.DefaultTransport
	}

	data, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] \trequest: %s %s\n==============================\n%s\n============================\n", req.Method, req.URL.String(), string(data))

	res, err := r.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	data, err = httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] \tresponse: %s %s\n==============================\n%s\n============================\n", req.Method, req.URL.String(), string(data))

	return res, err
}
