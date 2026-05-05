package models

import (
	"encoding/json"
	"errors"
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
