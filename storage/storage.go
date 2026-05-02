package storage

import "iot/models"

type Storage interface {
	GetShelves() []models.Shelf
	UpdateWeight(shelfID int, weight float64) error
}
