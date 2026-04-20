package dto

type CreateShortURLRequest struct {
	URL              string `json:"url"`
	ExpiresInSeconds *int64 `json:"expires_in_seconds,omitempty"`
}

type CreateShortURLResponse struct {
	ShortURL string `json:"short_url"`
	ID       string `json:"id"`
}
