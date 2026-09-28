package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"blog-server/internal/apperr"
	"blog-server/internal/http/response"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimitByIP allows perMinute requests per client IP with the given burst.
func RateLimitByIP(perMinute, burst int) gin.HandlerFunc {
	if perMinute <= 0 {
		return func(c *gin.Context) { c.Next() }
	}
	var (
		mu       sync.Mutex
		visitors = map[string]*visitor{}
		every    = rate.Every(time.Minute / time.Duration(perMinute))
	)

	go func() {
		for range time.Tick(5 * time.Minute) {
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 10*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		v, ok := visitors[ip]
		if !ok {
			v = &visitor{limiter: rate.NewLimiter(every, burst)}
			visitors[ip] = v
		}
		v.lastSeen = time.Now()
		allowed := v.limiter.Allow()
		mu.Unlock()

		if !allowed {
			response.Fail(c, apperr.ErrTooMany)
			return
		}
		c.Next()
	}
}
