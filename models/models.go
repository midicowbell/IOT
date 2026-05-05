package models

type Product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Weight    float64 `json:"weight"`     // weight is given in grams only
	MinWeight float64 `json:"min_weight"` // weight is given in grams only
}

type Shelf struct {
	ID         int      `json:"id"`
	Product    *Product `json:"product"`
	CurrWeight float64  `json:"curr_weight"`
}

func NewProduct(id int, name string, weight float64, minWeight float64) (*Product, error) {
	if weight < 0 || minWeight < 0 {
		return nil, ErrorNegativeWeight
	}
	return &Product{
		ID:        id,
		Name:      name,
		Weight:    weight,
		MinWeight: minWeight,
	}, nil
}
func NewShelf(id int, weight float64) (*Shelf, error) {
	if weight < 0 {
		return nil, ErrorNegativeWeight
	}
	return &Shelf{
		ID:         id,
		CurrWeight: weight,
	}, nil
}
func (s *Shelf) SetProduct(p *Product) {
	s.Product = p
}

func (s *Shelf) DeleteProduct() error {
	if s.Product == nil {
		return ShelfIsEmpty
	}
	s.Product = nil
	return nil
}

func (s *Shelf) GetProduct() (*Product, error) {
	if s.Product == nil {
		return nil, ShelfIsEmpty
	}
	return s.Product, nil
}

func (s *Shelf) UpdateWeight(weight float64) error {
	if weight < 0 {
		return ErrorNegativeWeight
	}
	s.CurrWeight = weight
	return nil
}

func (s *Shelf) NeedsRefill() bool {
	if s.Product == nil {
		return false
	}
	return s.CurrWeight < s.Product.MinWeight
}

func (s *Shelf) GetQuantity() int {
	if s.Product == nil || s.Product.Weight == 0 {
		return 0
	}
	return int(s.CurrWeight / s.Product.Weight)
}

func (p *Product) UpdateMinWeight(minWeight float64) error {
	if minWeight < 0 {
		return ErrorNegativeWeight
	}
	p.MinWeight = minWeight
	return nil
}
