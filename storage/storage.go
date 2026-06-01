package storage

import (
	"context"
	"iot/models"
)

type Storage interface {
	GetShelves(ctx context.Context) ([]models.Shelf, error)
	GetShelfByID(ctx context.Context, shelfID int) (*models.Shelf, error)
	UpdateWeight(ctx context.Context, shelfID int, weight float64, status string) error
	UpdateShelfProduct(ctx context.Context, shelfID int, productID *int) error
	AddShelf(ctx context.Context, shelf *models.Shelf) error
	DeleteShelf(ctx context.Context, shelfID int) error
}
