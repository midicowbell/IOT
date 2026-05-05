package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"iot/models"
	"iot/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type WeightDTO struct {
	ID     int     `json:"id"`
	Weight float64 `json:"weight"`
}

type AnsDTO struct {
	Answer string    `json:"answer"`
	Time   time.Time `json:"time"`
}

type ShelfDTO struct {
	ID int `json:"id"`
}
type ProductDTO struct {
	ProductID  int     `json:"product_id"`
	ShelfID    int     `json:"shelf_id"`
	Name       string  `json:"name"`
	UnitWeight float64 `json:"unit_weight"`
	MinWeight  float64 `json:"min_weight"`
}

func NewAnsDTO(str string) AnsDTO {
	return AnsDTO{
		Answer: str,
		Time:   time.Now(),
	}
}

type HTTPHandlers struct {
	service *service.StockService
}

func NewHTTPHandlers(service *service.StockService) *HTTPHandlers {
	return &HTTPHandlers{
		service: service,
	}
}

func (h *HTTPHandlers) HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	status := h.service.GetFullStatus()

	b, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to HTTP response")
		return
	}
}

func (h *HTTPHandlers) HandleUpdateWeight(w http.ResponseWriter, r *http.Request) {
	var data WeightDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	allert, err := h.service.UpdateProductWeight(data.ID, data.Weight)
	if err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	ans := NewAnsDTO(allert)
	b, err := json.MarshalIndent(ans, "", "    ")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to HTTP response")
		return
	}
}

func (h *HTTPHandlers) HandleAddShelf(w http.ResponseWriter, r *http.Request) {
	var shelfDTO ShelfDTO
	if err := json.NewDecoder(r.Body).Decode(&shelfDTO); err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	shelf, err := models.NewShelf(shelfDTO.ID, 0)
	if err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	if err := h.service.AddShelf(shelf); err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPHandlers) DeleteShelf(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errDTO := models.NewErrorDTO(errors.New("invalid shelf ID format"))
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteShelf(id); err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandlers) HandleAddProduct(w http.ResponseWriter, r *http.Request) {
	var recProduct ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&recProduct); err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	product, err := models.NewProduct(recProduct.ProductID, recProduct.Name, recProduct.UnitWeight, recProduct.MinWeight)
	if err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	if err := h.service.FillShelf(recProduct.ShelfID, product); err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	b, err := json.MarshalIndent(product, "", "    ")
	if err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to http response")
		return
	}
}

func (h *HTTPHandlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errDTO := models.NewErrorDTO(errors.New("invalid shelf ID format"))
		models.HTTPError(w, errDTO, err)
		return
	}
	if err := h.service.DeleteProduct(id); err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandlers) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errDTO := models.NewErrorDTO(errors.New("invalid shelf ID format"))
		models.HTTPError(w, errDTO, err)
		return
	}
	product, err := h.service.GetProduct(id)
	if err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}

	b, err := json.MarshalIndent(product, "", "    ")
	if err != nil {
		errDTO := models.NewErrorDTO(err)
		models.HTTPError(w, errDTO, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to http response")
		return
	}
}
