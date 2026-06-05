package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kaanchinar/url-shortener/dto"
	"github.com/kaanchinar/url-shortener/service"
)

type URLHandler struct {
	svc     *service.URLService
	baseURL string
}

func NewURLHandler(s *service.URLService, baseURL string) *URLHandler {
	return &URLHandler{svc: s, baseURL: strings.TrimRight(baseURL, "/")}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateShortURLRequest

	if r.Body == nil {
		http.Error(w, `{"error":"request body is required"}`, http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"invalid JSON: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, `{"error":"url field is required"}`, http.StatusBadRequest)
		return
	}

	shortID, err := h.svc.ShortenUrl(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Build the short URL from configurable base URL, or fall back to request headers
	var fullShortURL string
	if h.baseURL != "" {
		fullShortURL = fmt.Sprintf("%s/%s", h.baseURL, shortID)
	} else {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		if fwdProto := r.Header.Get("X-Forwarded-Proto"); fwdProto != "" {
			scheme = fwdProto
		}
		host := r.Host
		if fwdHost := r.Header.Get("X-Forwarded-Host"); fwdHost != "" {
			host = fwdHost
		}
		fullShortURL = fmt.Sprintf("%s://%s/%s", scheme, host, shortID)
	}

	res := dto.CreateShortURLResponse{
		ID:       shortID,
		ShortURL: fullShortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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