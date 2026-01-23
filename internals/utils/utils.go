package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func ReadParam(r *http.Request, paramName string) (int64, error) {
	idParam := chi.URLParam(r, paramName)
	if idParam == "" {
		return 0, errors.New("Invalid id parameter")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter type")
	}

	return id, nil
}

func ReadQueryParamString(r *http.Request, paramName string) (string, error) {
	val := r.URL.Query().Get(paramName)
	if val == "" {
		return "", errors.New("invalid query parameter")
	}
	return val, nil
}

func ReadQueryParamInt64(r *http.Request, paramName string) (int64, error) {
	val := r.URL.Query().Get(paramName)
	if val == "" {
		return 0, errors.New("invalid query parameter")
	}

	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.New("invalid query parameter type")
	}

	return id, nil
}


func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}