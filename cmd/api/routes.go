package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	r := httprouter.New()

	r.NotFound = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	r.HandlerFunc(http.MethodGet, "/v1/tracks", app.requireActivatedUser(app.listTracksHandler))
	r.HandlerFunc(http.MethodPost, "/v1/tracks", app.requireActivatedUser(app.createTrackHandler))
	r.HandlerFunc(http.MethodGet, "/v1/tracks/:id", app.requireActivatedUser(app.showTrackHandler))
	r.HandlerFunc(http.MethodPatch, "/v1/tracks/:id", app.requireActivatedUser(app.updateTrackHandler))
	r.HandlerFunc(http.MethodDelete, "/v1/tracks/:id", app.requireActivatedUser(app.deleteTrackHandler))

	r.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	r.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	r.HandlerFunc(http.MethodPost, "/v1/tokens/auth", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(r)))

}
