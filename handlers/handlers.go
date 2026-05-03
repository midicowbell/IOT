package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"iot/models"
	"iot/service"
	"net/http"
	"time"
)

type WeightDTO struct {
	ID     int     `json:"id"`
	Weight float64 `json:"weight"`
}

type AnsDTO struct {
	Answer string    `json:"answer"`
	Time   time.Time `json:"time"`
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

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to HTTP response")
		return
	}
}

func (h *HTTPHandlers) HandleUpdateWeight(w http.ResponseWriter, r *http.Request) {
	var data WeightDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	allert, err := h.service.UpdateProductWeight(data.ID, data.Weight)
	if err != nil {
		if errors.Is(err, models.ShelfIsNotFound) {
			http.Error(w, models.ErrorToString(err, ""), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	}
}
