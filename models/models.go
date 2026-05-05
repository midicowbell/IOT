package models

type Product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Weight    float64 `json:"weight"`
	MinWeight float64 `json:"min_weight"`
	Unit      string  `json:"unit"`
}

type Shelf struct {
	ID         int      `json:"id"`
	Product    *Product `json:"product"`
	CurrWeight float64  `json:"curr_weight"`
}

func NewProduct(id int, name string, weight float64, minWeight float64, unit string) (*Product, error) {
	if weight < 0 || minWeight < 0 {
		return nil, ErrorNegativeWeight
	}
	return &Product{
		ID:        id,
		Name:      name,
		Weight:    weight,
		MinWeight: minWeight,
		Unit:      unit,
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
	return int(s.CurrWeight) / int(s.Product.Weight)
}

func (p *Product) UpdateMinWeight(minWeight float64) error {
	if minWeight < 0 {
		return ErrorNegativeWeight
	}
	p.MinWeight = minWeight
	return nil
}
