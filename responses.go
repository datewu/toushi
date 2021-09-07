package toushi

import (
	"fmt"
	"net/http"
)

// OKJSON handle 200 respose
func OKJSON(w http.ResponseWriter, r *http.Request, data Envelope) {
	WriteJSON(w, http.StatusOK, data, nil)
}

// ErrResponse handle err respose
func ErrResponse(w http.ResponseWriter, r *http.Request, status int, msg interface{}) {
	data := Envelope{"error": msg}
	WriteJSON(w, status, data, nil)
}

func errResponse(code int, msg interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ErrResponse(w, r, code, msg)
	}
}

// NotFoundResponse handle 404 respose
var NotFountResponse = errResponse(http.StatusNotFound,
	"the requested resource could not be found",
)

// EditConflictResponse handle 409 respose
var EditConflictResponse = errResponse(http.StatusConflict,
	"unable to update the record due to an edit conflict, please try later",
)

// RateLimitExceededResponse handle 429 respose
var RateLimitExceededResponse = errResponse(http.StatusTooManyRequests,
	"rate limit exceeded",
)

// InvalidCredentialsResponse handle 400 respose
var InvalidCredentialsResponse = errResponse(http.StatusBadRequest,
	"invalid authentication credentials",
)

// InvalidAuthenticationTokenResponse handle 401 respose
var InvalidAuthenticationTokenResponse = errResponse(http.StatusUnauthorized,
	"invalid or missing authentication token",
)

// AuthenicationRequiredResponse handle 401 respose
var AuthenticationRequireResponse = errResponse(http.StatusUnauthorized,
	"you must be authenticated to access this resource",
)

// InactiveUserResponse handle 403 respose
var InactiveAccountResponse = errResponse(http.StatusForbidden,
	"your user account must be activated to access this resource",
)

// NotPermittedResponse handle 403 respose
var NotPermittedResponse = errResponse(http.StatusForbidden,
	"your user account doesn't have the necessary permissions to access this resource",
)

// MethodNotAllowedResponse handle 405 respose
var MethodNotAllowResponse http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s mehtod is not supported for this resource", r.Method)
	errResponse(http.StatusMethodNotAllowed, msg)(w, r)
}

// BadRequestResponse handle 400 respose
func BadRequestResponse(err error) http.HandlerFunc {
	return errResponse(http.StatusBadRequest, err.Error())
}

// FailedValidationResponse handle 400 respose
func FailedValidationResponse(errs map[string]string) http.HandlerFunc {
	return errResponse(http.StatusBadRequest, errs)
}

// ServerErrorResponse handle 500 respose
func ServerErrResponse(err error) http.HandlerFunc {
	errs := map[string]interface{}{
		"error":  "the server encountered a problem and could not process your request",
		"detail": err.Error(),
	}
	return errResponse(http.StatusInternalServerError, errs)
}
