package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaanchinar/url-shortener/dto"
	"github.com/kaanchinar/url-shortener/service"
)

type URLHandler struct {
	svc *service.URLService
}

func NewURLHandler(s *service.URLService) *URLHandler {
	return &URLHandler{svc: s}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateShortURLRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusBadRequest)
		return
	}

	shortID, err := h.svc.ShortenUrl(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusBadRequest)
		return
	}

	fullShortURL := fmt.Sprintf("http://localhost:3000/%s", shortID)

	res := dto.CreateShortURLResponse{
		ID:       shortID,
		ShortURL: fullShortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return
	}
}

func (h *URLHandler) GetLongURL(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")
	if shortID == "" {
		http.Error(w, "Invalid short ID", http.StatusBadRequest)
		return
	}

	urlModel, err := h.svc.GetUrlById(r.Context(), shortID)

	if errors.Is(err, service.ErrURLExpired) {
		http.Error(w, "URL has expired", http.StatusGone)
		return
	}
	if err != nil || urlModel == nil {
		http.Error(w, "Failed to get url", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, urlModel.OriginalURL, http.StatusFound)
}
