package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blog-server/internal/apperr"
	"blog-server/internal/http/response"
)

// SecurityHeaders sets headers suited to a JSON-only API: nothing it returns should be rendered,
// framed, or sniffed as another content type.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		c.Next()
	}
}

// SetPageSecurityHeaders replaces the API headers with a looser set for the frontend pages. Images
// come from OSS, QQ avatars and arbitrary friend-link sites, so the CSP only locks down framing,
// plugins and <base>.
func SetPageSecurityHeaders(h http.Header) {
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "SAMEORIGIN")
	h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	h.Set("Content-Security-Policy", "frame-ancestors 'self'; object-src 'none'; base-uri 'self'")
}

// BodyLimit caps request bodies at def bytes, or at the override for matched route paths
// (as registered, e.g. "/api/v1/admin/articles/:id").
func BodyLimit(def int64, overrides map[string]int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil || c.Request.Body == http.NoBody {
			c.Next()
			return
		}
		limit := def
		if n, ok := overrides[c.FullPath()]; ok {
			limit = n
		}
		if c.Request.ContentLength > limit {
			response.Fail(c, apperr.ErrTooLarge)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}
