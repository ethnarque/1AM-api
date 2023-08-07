package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	p := payload{"error": message}

	err := app.writeJSON(w, status, p, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	msg := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "the request resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflitResponse(w http.ResponseWriter, r *http.Request) {
	msg := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, msg)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	msg := "rate limiter exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, msg)
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	msg := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, msg)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	msg := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, msg)
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	msg := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, msg)
}

func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	msg := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, msg)
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	msg := "your user account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, msg)
}
