package storage

import (
	"errors"
	"iot/models"
	"sync"
)

type MemStorage struct {
	mtx     sync.RWMutex
	shelves map[int]*models.Shelf
}

func NewMemStorage() *MemStorage {
	return &MemStorage{shelves: make(map[int]*models.Shelf)}
}

func (m *MemStorage) GetShelves() []models.Shelf {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	ans := make([]models.Shelf, 0, len(m.shelves))
	for _, val := range m.shelves {
		ans = append(ans, *val)
	}
	return ans
}

func (m *MemStorage) UpdateWeight(shelfID int, weight float64) error {
	if weight < 0 {
		return errors.New("weight should not be less than 0")
	}
	m.mtx.Lock()
	defer m.mtx.Unlock()
	shelf, ok := m.shelves[shelfID]
	if !ok {
		return errors.New("shelf with this ID not found")
	}
	shelf.CurrWeight = weight
	return nil
}
