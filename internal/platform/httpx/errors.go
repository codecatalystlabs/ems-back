package httpx

import "net/http"

var (
	ErrUnauthorized = http.StatusUnauthorized
	ErrForbidden    = http.StatusForbidden
	ErrNotFound     = http.StatusNotFound
)
