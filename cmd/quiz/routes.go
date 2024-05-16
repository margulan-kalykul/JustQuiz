package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	r.HandleFunc("/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	players := r.PathPrefix("/v1").Subrouter()

	// Players list
	players.HandleFunc("/players", app.getPlayersList).Methods("GET")
	// Create a new player
	players.HandleFunc("/players", app.createPlayerHandler).Methods("POST")
	// Get a player by id
	players.HandleFunc("/players/{id:[0-9]+}", app.getPlayerHandler).Methods("GET")
	// Update player data with id
	players.HandleFunc("/players/{id:[0-9]+}", app.updatePlayerHandler).Methods("PUT")
	// Delete player by id
	players.HandleFunc("/players/{id:[0-9]+}", app.requirePermissions("player:write", app.deletePlayerHandler)).Methods("DELETE")

	users := r.PathPrefix("/v1").Subrouter()
	// User handlers with Authentication
	users.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
