// SPDX-License-Identifier: Apache-2.0

package chi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ajay-ghs/luraJailbreak/v2/config"
	"github.com/ajay-ghs/luraJailbreak/v2/proxy"
	"github.com/go-chi/chi/v5"
)

func BenchmarkEndpointHandler_ko(b *testing.B) {
	p := func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		return nil, fmt.Errorf("This is %s", "a dummy error")
	}
	endpoint := &config.EndpointConfig{
		Timeout:     time.Second,
		CacheTTL:    6 * time.Hour,
		QueryString: []string{"b"},
	}

	router := chi.NewRouter()
	router.Handle("/_chi_endpoint/", NewEndpointHandler(endpoint, p))

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/_chi_endpoint/a?b=1", http.NoBody)
	req.Header.Set("Content-Type", "application/json")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkEndpointHandler_ok(b *testing.B) {
	pResp := proxy.Response{
		Data:       map[string]interface{}{},
		Io:         io.NopCloser(&bytes.Buffer{}),
		IsComplete: true,
		Metadata:   proxy.Metadata{},
	}
	p := func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		return &pResp, nil
	}
	endpoint := &config.EndpointConfig{
		Timeout:     time.Second,
		CacheTTL:    6 * time.Hour,
		QueryString: []string{"b"},
	}

	router := chi.NewRouter()
	router.Handle("/_chi_endpoint/", NewEndpointHandler(endpoint, p))

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/_chi_endpoint/a?b=1", http.NoBody)
	req.Header.Set("Content-Type", "application/json")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkEndpointHandler_ko_Parallel(b *testing.B) {
	p := func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		return nil, fmt.Errorf("This is %s", "a dummy error")
	}
	endpoint := &config.EndpointConfig{
		Timeout:     time.Second,
		CacheTTL:    6 * time.Hour,
		QueryString: []string{"b"},
	}

	router := chi.NewRouter()
	router.Handle("/_chi_endpoint/", NewEndpointHandler(endpoint, p))

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/_chi_endpoint/a?b=1", http.NoBody)
	req.Header.Set("Content-Type", "application/json")

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

func BenchmarkEndpointHandler_ok_Parallel(b *testing.B) {
	pResp := proxy.Response{
		Data:       map[string]interface{}{},
		Io:         io.NopCloser(&bytes.Buffer{}),
		IsComplete: true,
		Metadata:   proxy.Metadata{},
	}
	p := func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		return &pResp, nil
	}
	endpoint := &config.EndpointConfig{
		Timeout:     time.Second,
		CacheTTL:    6 * time.Hour,
		QueryString: []string{"b"},
	}

	router := chi.NewRouter()
	router.Handle("/_chi_endpoint/", NewEndpointHandler(endpoint, p))

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/_chi_endpoint/a?b=1", http.NoBody)
	req.Header.Set("Content-Type", "application/json")

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}
