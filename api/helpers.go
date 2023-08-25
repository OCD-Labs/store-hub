package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type envelop map[string]interface{}

// ValidationError represents a custom validation error that
// contains information about the violated fields and their messages.
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// shouldBindBody reads/parses & validates request body.
func (s *StoreHub) shouldBindBody(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	// Restrict r.Body to 1MB
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)

	// Disallow unknown fields.
	decoder.DisallowUnknownFields()
	err := decoder.Decode(obj)

	if err != nil {
		// expected error types
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("request body contains badly-formatted JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("request body contains badly-formatted JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("request body contains badly-formatted JSON for the field: %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("request body contains badly-formatted JSON ((at character %d))", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("empty request body")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown field: %s", field)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("request body must not be larger than %d bytes", maxBytes)
		default:
			return err
		}
	}

	// ensure r.Body is one json object
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("request body must contain only a single JSON")
	}

	err = s.bindValidation(w, r, obj)

	return err
}

// writeJSON writes and sends JSON response.
func (s *StoreHub) writeJSON(
	w http.ResponseWriter,
	statusCode int,
	data envelop,
	header http.Header,
) error {
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp = append(resp, '\n')

	for key, value := range header {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)

	return nil
}

// bindValidation is a helper function that binds the request
// (body, path variable and query string) to the given interface
// and validates it with the specified validator.
func (s *StoreHub) bindValidation(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) error {
	validate := validator.New()

	if err := validate.Struct(data); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		validationErrors := make([]ValidationError, 0, len(errs))

		for _, err := range errs {
			validationErrors = append(validationErrors, ValidationError{
				Field: err.Field(),
				Error: fmt.Sprintf("%s validation failed on '%s'", err.Tag(), err.Param()),
			})
		}

		s.logAndRespond(w, r, http.StatusBadRequest, validationErrors, err)

		return err
	}

	return nil
}

// shouldBindQuery extracts and binds query string to a struct using reflection:
func (s *StoreHub) shouldBindQuery(w http.ResponseWriter, r *http.Request, obj interface{}) (err error) {
	values := r.URL.Query()

	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fType := t.Field(i)

		tag := fType.Tag.Get("querystr")
		if tag == "-" {
			continue
		}
		value := values.Get(tag)

		if fType.Type == reflect.TypeOf(time.Time{}) {
			parsedTime, err := time.Parse("2006-01-02", value)
			if err != nil { // TODO: consider if ignoring the error is appropriate
				continue
			}
			f.Set(reflect.ValueOf(parsedTime))
		} else {
			setValue(f, value)
		}
	}

	err = s.bindValidation(w, r, obj)

	return
}

// ShouldBindPathVars extracts and binds path variables to a struct using reflection:
func (s *StoreHub) ShouldBindPathVars(w http.ResponseWriter, r *http.Request, obj interface{}) (err error) {
	params := httprouter.ParamsFromContext(r.Context())
	if params == nil {
		err = fmt.Errorf("URL parameters not found")
		s.errorResponse(w, r, http.StatusBadRequest, err)
		return err
	}

	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fType := t.Field(i)

		tag := fType.Tag.Get("path")
		if tag == "-" {
			continue
		}
		val := params.ByName(tag)

		setValue(f, val)
	}

	err = s.bindValidation(w, r, obj)

	return
}

// setValue leverage reflection to set values.
func setValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(uintValue)
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(floatValue)
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			field.SetBool(boolValue)
		}
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			if value == "" {
				field.Set(reflect.Zero(field.Type()))
				return
			}

			items := strings.Split(value, ",")
			if len(items) > 0 {
				field.Set(reflect.ValueOf(items))
				return
			}

			field.Set(reflect.Zero(field.Type()))
		}
	}
}

// errorResponse writes error response.
func (s *StoreHub) errorResponse(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	message interface{},
) {
	resp := envelop{
		"status": "error",
		"error": envelop{
			"message": message,
		},
	}
	err := s.writeJSON(w, statusCode, resp, nil)
	if err != nil {
		log.Error().Err(err).
			Str("request_method", r.Method).
			Str("request_url", r.URL.String()).
			Msg("failed to write response body")
		w.WriteHeader(500)
	}
}

// logAndRespond logs the error and sends the response.
func (s *StoreHub) logAndRespond(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	message interface{},
	err error,
) {
	// Log the error
	if err != nil {
		log.Error().Err(err).
			Str("request_method", r.Method).
			Str("request_url", r.URL.String()).
			Msg("Error processing request")
	}

	// Send the response
	s.errorResponse(w, r, statusCode, message)
}
