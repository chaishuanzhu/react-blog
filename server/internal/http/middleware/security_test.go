package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newLimitedRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders(), BodyLimit(16, map[string]int64{"/big": 64}))
	read := func(c *gin.Context) {
		if _, err := io.ReadAll(c.Request.Body); err != nil {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}
		c.Status(http.StatusNoContent)
	}
	r.POST("/small", read)
	r.POST("/big", read)
	return r
}

func TestBodyLimit(t *testing.T) {
	r := newLimitedRouter()
	cases := []struct {
		path    string
		size    int
		chunked bool
		want    int
	}{
		{"/small", 16, false, http.StatusNoContent},
		{"/small", 17, false, http.StatusRequestEntityTooLarge},
		{"/small", 17, true, http.StatusRequestEntityTooLarge},
		{"/big", 64, false, http.StatusNoContent},
		{"/big", 65, true, http.StatusRequestEntityTooLarge},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(strings.Repeat("a", tc.size)))
		if tc.chunked {
			req.ContentLength = -1
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.want {
			t.Errorf("%s size=%d chunked=%v: got %d, want %d", tc.path, tc.size, tc.chunked, w.Code, tc.want)
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	newLimitedRouter().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/small", nil))
	for _, h := range []string{"X-Content-Type-Options", "X-Frame-Options", "Referrer-Policy", "Content-Security-Policy"} {
		if w.Header().Get(h) == "" {
			t.Errorf("missing %s", h)
		}
	}
}
