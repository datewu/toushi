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

var NotFountResponse = errResponse(http.StatusNotFound,
	"the requested resource could not be found",
)
var EditConflictResponse = errResponse(http.StatusConflict,
	"unable to update the record due to an edit conflict, please try later",
)

var RateLimitExceededResponse = errResponse(http.StatusTooManyRequests,
	"rate limit exceeded",
)

var InvalidCredentialsResponse = errResponse(http.StatusBadRequest,
	"invalid authentication credentials",
)

var InvalidAuthenticationTokenResponse = errResponse(http.StatusBadRequest,
	"invalid or missing authentication token",
)

var AuthenticationRequireResponse = errResponse(http.StatusUnauthorized,
	"you must be authenticated to access this resource",
)

var InactiveAccountResponse = errResponse(http.StatusForbidden,
	"your user account must be activated to access this resource",
)

var NotPermittedResponse = errResponse(http.StatusForbidden,
	"your user account doesn't have the necessary permissions to access this resource",
)

var MethodNotAllowResponse http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s mehtod is not supported for this resource", r.Method)
	errResponse(http.StatusMethodNotAllowed, msg)(w, r)
}

func BadRequestResponse(err error) http.HandlerFunc {
	return errResponse(http.StatusBadRequest, err)
}

func FailedValidationResponse(errs map[string]string) http.HandlerFunc {
	return errResponse(http.StatusBadRequest, errs)
}

func ServerErrResponse(err interface{}) http.HandlerFunc {
	errs := map[string]interface{}{
		"error":  "the server encountered a problem and could not process your request",
		"detail": err,
	}
	return errResponse(http.StatusInternalServerError, errs)
}
