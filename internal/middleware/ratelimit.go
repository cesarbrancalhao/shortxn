package middleware

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mtx      sync.RWMutex
	r        rate.Limit
	b        int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

func (rl *RateLimiter) RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
		}

		return next(c)
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mtx.Lock()
	defer rl.mtx.Unlock()

	limiter, exists := rl.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.visitors[ip] = limiter
	}

	return limiter
}
