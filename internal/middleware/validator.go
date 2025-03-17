package middleware

import (
	"Shortxn/internal/domain"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type URLValidator struct {
	MaxLength int
}

func NewURLValidator(maxLength int) *URLValidator {
	return &URLValidator{MaxLength: maxLength}
}

func (v *URLValidator) ValidateURL(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input domain.URLRequest

		decoder := json.NewDecoder(c.Request().Body)
		if err := decoder.Decode(&input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format. Expected: {\"url\": \"https://example.com\"}")
		}

		if input.URL == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "URL cannot be empty")
		}

		if len(input.URL) > v.MaxLength {
			return echo.NewHTTPError(http.StatusBadRequest, "URL too long")
		}

		_, err := url.ParseRequestURI(input.URL)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL format")
		}
		c.Set("requestBody", input)

		return next(c)
	}
}
