package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(log *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			log.WithFields(logrus.Fields{
				"method":     c.Request().Method,
				"path":       c.Request().URL.Path,
				"status":     c.Response().Status,
				"duration":   time.Since(start),
				"ip":         c.RealIP(),
				"user_agent": c.Request().UserAgent(),
			}).Info("Request handled")

			return err
		}
	}
}
