package service

import (
	"fmt"
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
	return "Статус: Вес в норме", nil
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
