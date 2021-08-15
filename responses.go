package toushi

import (
	"fmt"
	"net/http"
)

// ErrResponse handle err respose
func ErrResponse(w http.ResponseWriter, r *http.Request, status int, msg interface{}) {
	data := envelope{"error": msg}
	WriteJSON(w, status, data, nil)
}

type errResponse struct {
	msg  interface{}
	code int
}

// ServeHTTP implement http.Hanlder
func (res *errResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ErrResponse(w, r, res.code, res.msg)
}

var NotFountResponse = errResponse{
	msg:  "the requested resource could not be found",
	code: http.StatusNotFound,
}
var EditConflictResponse = errResponse{
	msg:  "unable to update the record due to an edit conflict, please try later",
	code: http.StatusConflict,
}

var RateLimitExceededResponse = errResponse{
	msg:  "rate limit exceeded",
	code: http.StatusTooManyRequests,
}

var InvalidCredentialsResponse = errResponse{
	msg:  "invalid authentication credentials",
	code: http.StatusBadRequest,
}

var InvalidAuthenticationTokenResponse = errResponse{
	msg:  "invalid or missing authentication token",
	code: http.StatusBadRequest,
}

var AuthenticationRequireResponse = errResponse{
	msg:  "you must be authenticated to access this resource",
	code: http.StatusUnauthorized,
}

var InactiveAccountResponse = errResponse{
	msg:  "your user account must be activated to access this resource",
	code: http.StatusForbidden,
}

var NotPermittedResponse = errResponse{
	msg:  "your user account doesn't have the necessary permissions to access this resource",
	code: http.StatusForbidden,
}

func MethodNotAllowResponse(method string) http.Handler {
	return &errResponse{
		msg:  fmt.Sprintf("the %s mehtod is not supported for this resource", method),
		code: http.StatusMethodNotAllowed,
	}
}

func BadRequestResponse(err error) http.Handler {
	return &errResponse{
		msg:  err,
		code: http.StatusBadRequest,
	}
}

func FailedValidationResponse(errs map[string]string) http.Handler {
	return &errResponse{
		msg:  errs,
		code: http.StatusUnprocessableEntity,
	}
}

func ServerErrResponse(err interface{}) http.Handler {
	errs := map[string]interface{}{
		"error":  "the server encountered a problem and could not process your request",
		"detail": err,
	}
	return &errResponse{
		msg:  errs,
		code: http.StatusInternalServerError,
	}
}
