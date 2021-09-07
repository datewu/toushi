package toushi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// ErrNoToken is returned when a token is not found in the request
var ErrNoToken = errors.New("no token")

func GetBeareToken(r *http.Request, name string) (string, error) {
	head, err := GetToken(r, name)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(head, "Bearer ") {
		return "", errors.New("token must be a Bearer token")
	}
	return strings.TrimPrefix(head, "Bearer "), nil
}

func GetToken(r *http.Request, name string) (string, error) {
	if name == "" {
		name = "token"
	}
	q := ReadQuery(r, name, "") // for ws query
	if q != "" {
		return q, nil
	}
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", ErrNoToken
	}
	return token, nil
}

func ReadQuery(r *http.Request, key string, defaultValue string) string {
	qs := r.URL.Query()
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func ReadCSVQuery(r *http.Request, key string, defaultValue []string) []string {
	qs := r.URL.Query()
	cs := qs.Get(key)
	if cs == "" {
		return defaultValue
	}
	return strings.Split(cs, ",")
}

func ReadInt64Query(r *http.Request, key string, defaultValue int64) int64 {
	qs := r.URL.Query()
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

func ReadParams(r *http.Request, name string) string {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName(name)
}

func ReadInt64Param(r *http.Request, name string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName(name), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	const maxBodySize = 8 * 1_048_576 // 8MB for max readJSON body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBodySize))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalErr *json.UnmarshalTypeError
		var invalidUnmarshalErr *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalErr):
			if unmarshalErr.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalErr.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", unmarshalErr.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

			// an open issue at https://github.com/golang/go/issues/29035
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

			// an open issue at https://github.com/golang/go/issues/30715
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBodySize)

		case errors.As(err, &invalidUnmarshalErr):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single json value")
	}
	return nil
}
