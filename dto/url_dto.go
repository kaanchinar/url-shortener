package dto

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type CreateShortURLResponse struct {
	ShortURL string `json:"short_url"`
	ID       string `json:"id"`
}
