package main

import (
	"Shortxn/internal/config"
	"Shortxn/internal/domain"
	"Shortxn/internal/infra/metrics"
	"Shortxn/internal/infra/postgres"
	"Shortxn/internal/infra/rabbitmq"
	"Shortxn/internal/infra/redis"
	"Shortxn/internal/middleware"
	"Shortxn/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Config struct {
	PostgresURL string
	RedisAddr   string
	ServerPort  string
	RabbitMQURL string
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	rateLimiter := middleware.NewRateLimiter(rate.Limit(cfg.RateLimit), cfg.RateBurst)
	urlValidator := middleware.NewURLValidator(cfg.MaxURLLength)

	urlRepo, err := postgres.NewURLRepository(cfg.PostgresURL)
	if err != nil {
		panic(err)
	}

	cacheRepo := redis.NewCacheRepository(cfg.RedisAddr)
	urlService := service.NewURLService(urlRepo, cacheRepo)

	publisher, err := rabbitmq.NewPublisher(cfg.RabbitMQURL)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	e := echo.New()

	e.HTTPErrorHandler = customErrorHandler

	e.Use(middleware.LoggingMiddleware(log))
	e.Use(rateLimiter.RateLimit)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.POST("/shorten", func(c echo.Context) error {
		start := time.Now()
		defer func() {
			metrics.ResponseTime.WithLabelValues("shorten", "POST").Observe(time.Since(start).Seconds())
		}()

		req := c.Get("requestBody").(domain.URLRequest)
		log.WithField("url", req.URL).Info("Shortening URL")

		url, err := urlService.ShortenURL(req.URL)
		if err != nil {
			log.WithError(err).Error("Failed to shorten URL")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to shorten URL")
		}

		metrics.URLShortened.Inc()
		return c.JSON(http.StatusOK, ShortenResponse{ShortURL: url.ShortURL})
	}, urlValidator.ValidateURL)

	e.GET("/:id", func(c echo.Context) error {
		start := time.Now()
		defer func() {
			metrics.ResponseTime.WithLabelValues("redirect", "GET").Observe(time.Since(start).Seconds())
		}()

		id := c.Param("id")
		metrics.URLRedirects.WithLabelValues(id).Inc()

		event := &domain.ClickEvent{
			URLId:     id,
			Timestamp: time.Now(),
			UserAgent: c.Request().UserAgent(),
			IPAddress: c.RealIP(),
		}
		go publisher.PublishClickEvent(event)

		if longURL, err := cacheRepo.Get(id); err == nil {
			return c.Redirect(http.StatusMovedPermanently, longURL)
		}

		url, err := urlRepo.GetByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		cacheRepo.Set(id, url.LongURL, 24*time.Hour)

		go urlRepo.IncrementClicks(id) // Increment asynchronously

		return c.Redirect(http.StatusMovedPermanently, url.LongURL)
	})

	e.Logger.Fatal(e.Start(cfg.ServerPort))
}

func customErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		msg = err.Error()
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, map[string]interface{}{
				"message": fmt.Sprintf("%v", msg),
			})
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
