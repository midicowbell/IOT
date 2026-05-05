package service

import (
	"fmt"
	"iot/models"
	"iot/storage"
)

type ShelfStatus struct {
	ID            int     `json:"id"`
	ProductName   string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	CurrentWeight float64 `json:"current_weight"`
	Status        string  `json:"status"`
}
type StockService struct {
	repo storage.Storage
}

func NewStockService(s storage.Storage) *StockService {
	return &StockService{
		repo: s,
	}
}

func (s *StockService) UpdateProductWeight(shelfID int, newWeight float64) (string, error) {
	err := s.repo.UpdateWeight(shelfID, newWeight)
	if err != nil {
		return "", err
	}
	shelf, err := s.repo.GetShelfByID(shelfID)
	if err != nil {
		return "", err
	}
	if shelf.NeedsRefill() {
		return fmt.Sprintf("ALERT: На полке %d (%s) заканчивается товар! Осталось: %d шт.",
			shelf.ID, shelf.Product.Name, shelf.GetQuantity()), nil
	}
	return "STATUS: Вес в норме", nil
}

func (s *StockService) GetFullStatus() []ShelfStatus {
	shelves := s.repo.GetShelves()
	var report []ShelfStatus

	for _, shelf := range shelves {
		status := "OK"
		if shelf.NeedsRefill() {
			status = "REFILL"
		}
		report = append(report, ShelfStatus{
			ID:            shelf.ID,
			ProductName:   shelf.Product.Name,
			Quantity:      shelf.GetQuantity(),
			CurrentWeight: shelf.CurrWeight,
			Status:        status,
		})
	}
	return report
}

func (s *StockService) AddShelf(shelf *models.Shelf) error {
	return s.repo.AddShelf(shelf)
}

func (s *StockService) DeleteShelf(id int) error {
	return s.repo.DeleteShelf(id)
}

func (s *StockService) GetShelfByID(id int) (*models.Shelf, error) {
	return s.repo.GetShelfByID(id)
}

func (s *StockService) FillShelf(shelfID int, product *models.Product) error {
	shelf, err := s.repo.GetShelfByID(shelfID)
	if err != nil {
		return err
	}
	shelf.SetProduct(product)
	return nil
}
