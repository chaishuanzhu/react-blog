package http

import (
	"context"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.yaml.in/yaml/v3"

	"blog-server/api"
	"blog-server/internal/config"
)

var pathParam = regexp.MustCompile(`:(\w+)`)

// TestSpecCoversRoutes keeps api/openapi.yaml in sync with the routes registered in NewRouter.
func TestSpecCoversRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{APIDocs: true, AllowedOrigins: []string{"http://localhost:3000"}}
	r, err := NewRouter(context.Background(), cfg, nil, nil)
	if err != nil {
		t.Fatalf("new router: %v", err)
	}
	var routes []string
	for _, rt := range r.Routes() {
		if strings.HasPrefix(rt.Path, "/api/docs") || rt.Path == "/api/openapi.yaml" {
			continue
		}
		routes = append(routes, rt.Method+" "+pathParam.ReplaceAllString(rt.Path, "{$1}"))
	}

	var spec struct {
		Paths map[string]map[string]any `yaml:"paths"`
	}
	if err := yaml.Unmarshal(api.OpenAPI, &spec); err != nil {
		t.Fatalf("parse openapi.yaml: %v", err)
	}
	var documented []string
	for path, item := range spec.Paths {
		for method := range item {
			switch method {
			case "get", "post", "put", "patch", "delete":
				documented = append(documented, strings.ToUpper(method)+" "+path)
			}
		}
	}

	if len(routes) < 40 {
		t.Fatalf("only %d routes collected, expected the full API", len(routes))
	}
	for _, rt := range routes {
		if !slices.Contains(documented, rt) {
			t.Errorf("route %s is missing from api/openapi.yaml", rt)
		}
	}
	for _, d := range documented {
		if !slices.Contains(routes, d) {
			t.Errorf("api/openapi.yaml documents %s, which is not registered", d)
		}
	}
}
