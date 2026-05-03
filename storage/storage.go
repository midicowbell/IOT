package storage

import "iot/models"

type Storage interface {
	GetShelves() []models.Shelf
	GetShelfByID(id int) (*models.Shelf, error)
	UpdateWeight(shelfID int, weight float64) error
	AddShelf(shelf *models.Shelf) error
	DeleteShelf(id int) error
}
