package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"blog-server/internal/apperr"
	"blog-server/internal/http/response"
)

const userIDKey = "userID"

type TokenVerifier func(ctx context.Context, token string) (uint64, error)

// Authenticate reads a Bearer token. With required=false, requests without a token pass through
// anonymously, but a token that is present and invalid is still rejected.
func Authenticate(verify TokenVerifier, required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			if required {
				response.Fail(c, apperr.ErrUnauthorized)
				return
			}
			c.Next()
			return
		}
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			response.Fail(c, apperr.ErrUnauthorized)
			return
		}
		id, err := verify(c.Request.Context(), token)
		if err != nil {
			response.Fail(c, err)
			return
		}
		c.Set(userIDKey, id)
		c.Next()
	}
}

func UserID(c *gin.Context) (uint64, bool) {
	v, ok := c.Get(userIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}
