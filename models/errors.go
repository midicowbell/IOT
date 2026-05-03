package models

import "errors"

var ErrorNegativeWeight = errors.New("weight cannot be negative")
var ShelfIsEmpty = errors.New("shelf is empty")
var ShelfIsNotFound = errors.New("shelf not found")

func ErrorToString(err error, anyinformation string) string {
	return err.Error() + anyinformation
}
