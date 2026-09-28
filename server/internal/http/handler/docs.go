package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files/v2"

	"blog-server/api"
)

// Relative to /api/docs/, so the page also works behind the webpack dev-server proxy.
const swaggerInitializer = `window.onload = function () {
  window.ui = SwaggerUIBundle({
    url: '../openapi.yaml',
    dom_id: '#swagger-ui',
    deepLinking: true,
    persistAuthorization: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: 'StandaloneLayout'
  });
};
`

// Docs serves the OpenAPI spec and a bundled Swagger UI.
type Docs struct{}

func (Docs) Spec(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", api.OpenAPI)
}

var swaggerAssets = http.StripPrefix("/api/docs", http.FileServerFS(swaggerFiles.FS))

// UI must be mounted at /api/docs/*file.
func (Docs) UI(c *gin.Context) {
	// Swagger UI sets inline styles; everything else is same-origin.
	c.Header("Content-Security-Policy",
		"default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; frame-ancestors 'none'")
	switch file := c.Param("file"); {
	case file == "/swagger-initializer.js":
		c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte(swaggerInitializer))
	case strings.HasSuffix(file, ".map"):
		c.Status(http.StatusNotFound)
	default:
		swaggerAssets.ServeHTTP(c.Writer, c.Request)
	}
}
