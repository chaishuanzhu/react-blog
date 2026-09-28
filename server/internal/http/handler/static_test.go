package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func writeFile(t *testing.T, name, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestSPAServe(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "index.html"), "INDEX")
	writeFile(t, filepath.Join(root, "js", "app.123.js"), "JS")
	writeFile(t, filepath.Join(root, "assets", "logo.svg"), "SVG")
	writeFile(t, filepath.Join(filepath.Dir(root), "secret.txt"), "SECRET")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.NoRoute(SPA{Root: root, Prefix: "/admin"}.Serve)

	cases := []struct {
		path, wantBody, wantCache string
		wantCode                  int
	}{
		{"/admin/", "INDEX", "no-cache", 200},
		{"/admin/home", "INDEX", "no-cache", 200},
		{"/admin/js/app.123.js", "JS", "immutable", 200},
		{"/admin/assets/logo.svg", "SVG", "max-age=3600", 200},
		{"/admin/js/missing.js", "", "", 404},
		{"/admin/js", "INDEX", "no-cache", 200},
		{"/admin/../secret.txt", "", "", 404},
		{"/admin/..%2fsecret.txt", "", "", 404},
	}
	for _, tc := range cases {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if w.Code != tc.wantCode {
			t.Errorf("%s: code %d, want %d", tc.path, w.Code, tc.wantCode)
			continue
		}
		if tc.wantCode != 200 {
			continue
		}
		if w.Body.String() != tc.wantBody {
			t.Errorf("%s: body %q, want %q", tc.path, w.Body.String(), tc.wantBody)
		}
		if !strings.Contains(w.Header().Get("Cache-Control"), tc.wantCache) {
			t.Errorf("%s: Cache-Control %q, want %q", tc.path, w.Header().Get("Cache-Control"), tc.wantCache)
		}
	}
}
