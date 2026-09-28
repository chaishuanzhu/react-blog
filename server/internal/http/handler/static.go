package handler

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SPA serves a built single-page app mounted at Prefix ("" or e.g. "/admin"): existing files are
// served as-is and any other extensionless path falls back to index.html for client-side routing.
type SPA struct {
	Root   string
	Prefix string
}

func (s SPA) Serve(c *gin.Context) {
	rel := path.Clean("/" + strings.TrimPrefix(c.Request.URL.Path, s.Prefix))
	if rel != "/" {
		if serveFile(c, filepath.Join(s.Root, filepath.FromSlash(rel)), assetCacheControl(rel)) {
			return
		}
		// A missing script or image must not be answered with HTML.
		if path.Ext(rel) != "" {
			c.Status(http.StatusNotFound)
			return
		}
	}
	if !serveFile(c, filepath.Join(s.Root, "index.html"), "no-cache") {
		c.Status(http.StatusNotFound)
	}
}

func serveFile(c *gin.Context, name, cacheControl string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		return false
	}
	c.Header("Cache-Control", cacheControl)
	http.ServeContent(c.Writer, c.Request, info.Name(), info.ModTime(), f)
	return true
}

// Bundles under js/ and css/ carry a content hash in their file name.
func assetCacheControl(rel string) string {
	if strings.HasPrefix(rel, "/js/") || strings.HasPrefix(rel, "/css/") {
		return "public, max-age=31536000, immutable"
	}
	return "public, max-age=3600"
}
