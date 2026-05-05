package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var ErrorNegativeWeight = errors.New("weight cannot be negative")
var ShelfIsEmpty = errors.New("shelf is empty")
var ShelfIsNotFound = errors.New("shelf not found")

func ErrorToString(err error, anyinformation string) string {
	return err.Error() + anyinformation
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func NewErrorDTO(err error) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}
}

func HTTPError(w http.ResponseWriter, errDTO ErrorDTO, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, ShelfIsNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, ShelfIsEmpty), errors.Is(err, ErrorNegativeWeight):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(errDTO.ToString()))
}
