// SPDX-License-Identifier: Apache-2.0

/*
Package gorilla provides some basic implementations for building routers based on gorilla/mux
*/
package gorilla

import (
	"net/http"

	gorilla "github.com/gorilla/mux"

	"github.com/ajay-ghs/luraJailbreak/v2/logging"
	"github.com/ajay-ghs/luraJailbreak/v2/proxy"
	"github.com/ajay-ghs/luraJailbreak/v2/router"
	"github.com/ajay-ghs/luraJailbreak/v2/router/mux"
	"github.com/ajay-ghs/luraJailbreak/v2/transport/http/server"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DefaultFactory returns a net/http mux router factory with the injected proxy factory and logger
func DefaultFactory(pf proxy.Factory, logger logging.Logger) router.Factory {
	return mux.NewFactory(DefaultConfig(pf, logger))
}

// DefaultConfig returns the struct that collects the parts the router should be builded from
func DefaultConfig(pf proxy.Factory, logger logging.Logger) mux.Config {
	return mux.Config{
		Engine:         gorillaEngine{gorilla.NewRouter()},
		Middlewares:    []mux.HandlerMiddleware{},
		HandlerFactory: mux.CustomEndpointHandler(mux.NewRequestBuilder(gorillaParamsExtractor)),
		ProxyFactory:   pf,
		Logger:         logger,
		DebugPattern:   "/__debug/{params}",
		EchoPattern:    "/__echo/{params}",
		RunServer:      server.RunServer,
	}
}

func gorillaParamsExtractor(r *http.Request) map[string]string {
	params := map[string]string{}
	title := cases.Title(language.Und)
	for key, value := range gorilla.Vars(r) {
		params[title.String(key)] = value
	}
	return params
}

type gorillaEngine struct {
	r *gorilla.Router
}

// Handle implements the mux.Engine interface from the lura router package
func (g gorillaEngine) Handle(pattern, method string, handler http.Handler) {
	g.r.Handle(pattern, handler).Methods(method)
}

// ServeHTTP implements the http:Handler interface from the stdlib
func (g gorillaEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.r.ServeHTTP(mux.NewHTTPErrorInterceptor(w), r)
}
