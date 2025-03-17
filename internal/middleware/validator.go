package middleware

import (
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
		var input struct {
			URL string `json:"url"`
		}

		if err := c.Bind(&input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid input format")
		}

		if len(input.URL) > v.MaxLength {
			return echo.NewHTTPError(http.StatusBadRequest, "URL too long")
		}

		_, err := url.ParseRequestURI(input.URL)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid URL format")
		}

		return next(c)
	}
}
