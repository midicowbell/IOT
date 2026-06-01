package service

import (
	"context"
	"fmt"
	"iot/models"
	"iot/storage"
)

type StockService struct {
	repo storage.Storage
}

func NewStockService(s storage.Storage) *StockService {
	return &StockService{
		repo: s,
	}
}

func (s *StockService) UpdateProductWeight(ctx context.Context, shelfID int, newWeight float64) (string, error) {
	shelf, err := s.repo.GetShelfByID(ctx, shelfID)
	if err != nil {
		return "", err
	}
	status := "OK"
	if shelf.Product == nil {
		status = "EMPTY"
	} else {
		status = "REFILL"
	}
	if err := s.repo.UpdateWeight(ctx, shelfID, newWeight, status); err != nil {
		return "", err
	}
	if shelf.NeedsRefill() {
		return fmt.Sprintf("ALERT: На полке %d (%s) заканчивается товар! Осталось: %d шт.",
			shelf.ID, shelf.Product.Name, shelf.GetQuantity()), nil
	}
	return "STATUS: Вес в норме", nil
}

func (s *StockService) GetFullStatus(ctx context.Context) ([]models.Shelf, error) {
	shelves, err := s.repo.GetShelves(ctx)
	if err != nil {
		return nil, err
	}
	var report []models.Shelf
	for _, shelf := range shelves {
		shelf.Quantity = shelf.GetQuantity()
		report = append(report, shelf)
	}
	return report, nil
}

func (s *StockService) AddShelf(ctx context.Context, shelf *models.Shelf) error {
	return s.repo.AddShelf(ctx, shelf)
}

func (s *StockService) DeleteShelf(ctx context.Context, shelfId int) error {
	return s.repo.DeleteShelf(ctx, shelfId)
}

func (s *StockService) GetShelfByID(ctx context.Context, shelfId int) (*models.Shelf, error) {
	return s.repo.GetShelfByID(ctx, shelfId)
}

func (s *StockService) FillShelf(ctx context.Context, shelfID int, product *models.Product) error {
	_, err := s.repo.GetShelfByID(ctx, shelfID)
	if err != nil {
		return err
	}
	return s.repo.UpdateShelfProduct(ctx, shelfID, &product.ID)
}

func (s *StockService) DeleteProduct(ctx context.Context, shelfID int) error {
	shelf, err := s.repo.GetShelfByID(ctx, shelfID)
	if err != nil {
		return err
	}
	if shelf.Product == nil {
		return models.ShelfIsEmpty
	}
	return s.repo.UpdateShelfProduct(ctx, shelfID, nil)
}

func (s *StockService) GetProduct(ctx context.Context, shelfID int) (*models.Product, error) {
	shelf, err := s.repo.GetShelfByID(ctx, shelfID)
	if err != nil {
		return nil, err
	}
	return shelf.GetProduct()
}
