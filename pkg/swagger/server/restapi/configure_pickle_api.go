// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"pickle-api/pkg/swagger/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name PickleAPI --spec ../../swagger.yml --principal interface{} --exclude-main

func configureFlags(api *operations.PickleAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PickleAPIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.BinProducer = runtime.ByteStreamProducer()
	api.JSONProducer = runtime.JSONProducer()
	api.TxtProducer = runtime.TextProducer()

	if api.GetHelloUserHandler == nil {
		api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(func(params operations.GetHelloUserParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetHelloUser has not yet been implemented")
		})
	}
	if api.GetPickleNameHandler == nil {
		api.GetPickleNameHandler = operations.GetPickleNameHandlerFunc(func(params operations.GetPickleNameParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPickleName has not yet been implemented")
		})
	}
	if api.GetPickleRandomHandler == nil {
		api.GetPickleRandomHandler = operations.GetPickleRandomHandlerFunc(func(params operations.GetPickleRandomParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPickleRandom has not yet been implemented")
		})
	}
	if api.GetPicklesHandler == nil {
		api.GetPicklesHandler = operations.GetPicklesHandlerFunc(func(params operations.GetPicklesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPickles has not yet been implemented")
		})
	}
	if api.CheckHealthHandler == nil {
		api.CheckHealthHandler = operations.CheckHealthHandlerFunc(func(params operations.CheckHealthParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.CheckHealth has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
