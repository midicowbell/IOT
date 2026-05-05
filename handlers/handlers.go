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

func NewAnsDTO(str string) AnsDTO {
	return AnsDTO{
		Answer: str,
		Time:   time.Now(),
	}
}

type HTTPHandlers struct {
	service *service.StockService
}

func NewHTTPHandler(service *service.StockService) *HTTPHandlers {
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
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	allert, err := h.service.UpdateProductWeight(data.ID, data.Weight)
	if err != nil {
		errDTO := models.NewErrorDTO(err)
		if errors.Is(err, models.ShelfIsNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
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
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	shelf, err := models.NewShelf(shelfDTO.ID, 0)
	if err != nil {
		errDTO := models.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}
	if err := h.service.AddShelf(shelf); err != nil {
		errDTO := models.NewErrorDTO(err)
		if errors.Is(err, models.ShelfIsEmpty) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPHandlers) DeleteShelf(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}
	if err := h.service.DeleteShelf(id); err != nil {
		errDTO := models.NewErrorDTO(err)
		if errors.Is(err, models.ShelfIsNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
