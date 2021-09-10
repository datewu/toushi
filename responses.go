package toushi

import (
	"fmt"
	"net/http"
)

// OKJSON handle 200 respose
func OKJSON(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, data, nil)
}

// OKText handle 200 respose
func OKText(w http.ResponseWriter, text string) {
	WriteStr(w, http.StatusOK, text, nil)
}

func errResponse(code int, msg interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := Envelope{"error": msg}
		WriteJSON(w, code, data, nil)
	}
}

// HandleNotFound handle 404 response
var HandleNotFound = errResponse(http.StatusNotFound,
	"the requested resource could not be found",
)

// HandleEditConflict handle 409 response
var HandleEditConflict = errResponse(http.StatusConflict,
	"unable to update the record due to an edit conflict, please try later",
)

// HandleRateLimitExceede handle 429 response
var HandleRateLimitExceede = errResponse(http.StatusTooManyRequests,
	"rate limit exceeded",
)

// HandleInvalidCredentials handle 400 response
var HandleInvalidCredentials = errResponse(http.StatusBadRequest,
	"invalid authentication credentials",
)

// HandleInvalidAuthenticationToken handle 401 response
var HandleInvalidAuthenticationToken = errResponse(http.StatusUnauthorized,
	"invalid or missing authentication token",
)

// HandleAuthenticationRequire handle 401 response
var HandleAuthenticationRequire = errResponse(http.StatusUnauthorized,
	"you must be authenticated to access this resource",
)

// HandleInactiveAccount handle 403 response
var HandleInactiveAccount = errResponse(http.StatusForbidden,
	"your user account must be activated to access this resource",
)

// HandleNotPermitted handle 403 response
var HandleNotPermitted = errResponse(http.StatusForbidden,
	"your user account doesn't have the necessary permissions to access this resource",
)

// HandleMethodNotAllow handle 405 response
var HandleMethodNotAllow http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s mehtod is not supported for this resource", r.Method)
	errResponse(http.StatusMethodNotAllowed, msg)(w, r)
}

// HandleBadRequest handle 400 response with custom message
func HandleBadRequestMsg(msg string) http.HandlerFunc {
	return errResponse(http.StatusBadRequest, msg)
}

// HandleBadRequestErr handle 400 response with a error
func HandleBadRequestErr(err error) http.HandlerFunc {
	return HandleBadRequestMsg(err.Error())
}

// HandleFailedValidation handle 400 response
func HandleFailedValidation(errs map[string]string) http.HandlerFunc {
	return errResponse(http.StatusBadRequest, errs)
}

// HandleServerErr handle 500 response
func HandleServerErr(err error) http.HandlerFunc {
	errs := map[string]interface{}{
		"error":  "the server encountered a problem and could not process your request",
		"detail": err.Error(),
	}
	return errResponse(http.StatusInternalServerError, errs)
}
