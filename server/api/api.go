// Package api holds the OpenAPI description of the HTTP API. Keep openapi.yaml in sync with
// internal/http/router.go; TestSpecCoversRoutes fails when they diverge.
package api

import _ "embed"

//go:embed openapi.yaml
var OpenAPI []byte
