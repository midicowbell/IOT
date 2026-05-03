package models

import "errors"

var ErorrNegativeWeight = errors.New("weight cannot be negative")
var ShelfIsEmpty = errors.New("shelf is empty")
var ShelfIsNotFound = errors.New("shelf not found")
