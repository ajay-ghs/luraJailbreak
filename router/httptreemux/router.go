// SPDX-License-Identifier: Apache-2.0

/*
Package httptreemux provides some basic implementations for building routers based on dimfeld/httptreemux
*/
package httptreemux

import (
	"net/http"

	"github.com/ajay-ghs/luraJailbreak/v2/logging"
	"github.com/ajay-ghs/luraJailbreak/v2/proxy"
	"github.com/ajay-ghs/luraJailbreak/v2/router"
	"github.com/ajay-ghs/luraJailbreak/v2/router/mux"
	"github.com/ajay-ghs/luraJailbreak/v2/transport/http/server"
	"github.com/dimfeld/httptreemux/v5"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DefaultFactory returns a net/http mux router factory with the injected proxy factory and logger
func DefaultFactory(pf proxy.Factory, logger logging.Logger) router.Factory {
	return mux.NewFactory(DefaultConfig(pf, logger))
}

// DefaultConfig returns the struct that collects the parts the router should be built from
func DefaultConfig(pf proxy.Factory, logger logging.Logger) mux.Config {
	return mux.Config{
		Engine:         NewEngine(httptreemux.NewContextMux()),
		Middlewares:    []mux.HandlerMiddleware{},
		HandlerFactory: mux.CustomEndpointHandler(mux.NewRequestBuilder(ParamsExtractor)),
		ProxyFactory:   pf,
		Logger:         logger,
		DebugPattern:   "/__debug/{params}",
		RunServer:      server.RunServer,
	}
}

func ParamsExtractor(r *http.Request) map[string]string {
	params := map[string]string{}
	title := cases.Title(language.Und)
	for key, value := range httptreemux.ContextParams(r.Context()) {
		params[title.String(key)] = value
	}
	return params
}

func NewEngine(m *httptreemux.ContextMux) Engine {
	return Engine{m}
}

type Engine struct {
	r *httptreemux.ContextMux
}

// Handle implements the mux.Engine interface from the lura router package
func (g Engine) Handle(pattern, method string, handler http.Handler) {
	g.r.Handle(method, pattern, handler.(http.HandlerFunc))
}

// ServeHTTP implements the http:Handler interface from the stdlib
func (g Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.r.ServeHTTP(mux.NewHTTPErrorInterceptor(w), r)
}
