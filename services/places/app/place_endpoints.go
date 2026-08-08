package app

import (
	"AirPolMonitor/services/places/domain"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type PlaceEndpoints struct {
	repo domain.PlaceRepositoryPort
}

type placeUpsertRequest struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Coordinates []float64 `json:"coordinates"`
}

type placeUpsertResponse struct {
	ID string `json:"id"`
}

func NewPlaceEndpoints(repo domain.PlaceRepositoryPort) *PlaceEndpoints {
	return &PlaceEndpoints{repo: repo}
}

func (h *PlaceEndpoints) Register(mux *http.ServeMux) {
	mux.HandleFunc("/places", h.handleCreatePlace)
	mux.HandleFunc("/places/", h.handleEditPlace)
}

func (h *PlaceEndpoints) handleCreatePlace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req placeUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	agg, err := domain.GetPlaceAggregate("", req.Name, h.repo)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := agg.HandleEditPlaceCommand(domain.EditPlaceCommand{
		ID:          req.ID,
		Name:        req.Name,
		Coordinates: req.Coordinates,
	}); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Save(r.Context(), agg); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, placeUpsertResponse{ID: req.ID})
}

func (h *PlaceEndpoints) handleEditPlace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	placeID := strings.TrimPrefix(r.URL.Path, "/places/")
	if placeID == "" || strings.Contains(placeID, "/") {
		writeJSONError(w, http.StatusBadRequest, "invalid place id in path")
		return
	}

	var req placeUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	agg, err := domain.GetPlaceAggregate(placeID, "", h.repo)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := agg.HandleEditPlaceCommand(domain.EditPlaceCommand{
		ID:          placeID,
		Name:        req.Name,
		Coordinates: req.Coordinates,
	}); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Save(r.Context(), agg); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, placeUpsertResponse{ID: placeID})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	if strings.TrimSpace(message) == "" {
		message = errors.New(http.StatusText(status)).Error()
	}
	writeJSON(w, status, map[string]string{"error": message})
}
