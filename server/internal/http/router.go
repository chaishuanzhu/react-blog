package http

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"blog-server/internal/apperr"
	"blog-server/internal/config"
	"blog-server/internal/http/handler"
	"blog-server/internal/http/middleware"
	"blog-server/internal/http/response"
	"blog-server/internal/mail"
	"blog-server/internal/service"
	"blog-server/internal/store"
)

func NewRouter(ctx context.Context, cfg *config.Config, db *sql.DB, notifier mail.Notifier) (*gin.Engine, error) {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		slog.Warn("invalid TRUSTED_PROXIES, trusting none", "err", err)
		_ = r.SetTrustedProxies(nil)
	}
	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.SecurityHeaders(),
		middleware.BodyLimit(1<<20, map[string]int64{
			"/api/v1/admin/articles":     5 << 20,
			"/api/v1/admin/articles/:id": 5 << 20,
			"/api/v1/admin/pages/:key":   5 << 20,
		}),
		cors.New(cors.Config{
			AllowOrigins:     cfg.AllowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
			ExposeHeaders:    []string{"X-Request-ID"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}),
	)

	q := store.New(db)
	authSvc := service.NewAuth(q, cfg)
	if err := authSvc.EnsureAdmin(ctx); err != nil {
		return nil, err
	}
	adminSvc := service.NewAdmin(db)
	uploads, err := service.NewUploads(cfg.OSS)
	if err != nil {
		return nil, err
	}

	health := &handler.Health{DB: db}
	pub := &handler.Public{Svc: service.NewPublic(q)}
	comments := &handler.Comment{Svc: service.NewComments(q, notifier, cfg), Auth: authSvc}
	auth := &handler.Auth{Svc: authSvc}
	admin := &handler.Admin{Svc: adminSvc, Uploads: uploads}

	requireAdmin := middleware.Authenticate(authSvc.VerifyToken, true)
	optionalAdmin := middleware.Authenticate(authSvc.VerifyToken, false)

	r.GET("/healthz", health.Live)
	r.GET("/readyz", health.Ready)

	if cfg.APIDocs {
		docs := handler.Docs{}
		r.GET("/api/openapi.yaml", docs.Spec)
		// /api/docs is redirected to /api/docs/ by gin's RedirectTrailingSlash.
		r.GET("/api/docs/*file", gzip.Gzip(gzip.DefaultCompression), docs.UI)
	}

	api := r.Group("/api/v1")
	api.GET("/articles", pub.ListArticles)
	api.GET("/articles/:id", pub.GetArticle)
	api.GET("/categories", pub.ListCategories)
	api.GET("/tags", pub.ListTags)
	api.GET("/moments", pub.ListMoments)
	api.GET("/friend-links", pub.ListFriendLinks)
	api.GET("/changelogs", pub.ListChangelogs)
	api.GET("/projects", pub.ListProjects)
	api.GET("/pages/:key", pub.GetPage)
	api.GET("/site", pub.GetSite)
	api.POST("/site/views", pub.RecordView)
	api.GET("/comments", comments.List)
	api.POST("/comments", optionalAdmin, middleware.RateLimitByIP(cfg.CommentRatePerMinute, cfg.CommentBurst), comments.Create)

	api.POST("/auth/login", middleware.RateLimitByIP(5, 5), auth.Login)
	me := api.Group("/auth", requireAdmin)
	me.GET("/me", auth.Me)
	me.PUT("/profile", auth.UpdateProfile)
	me.PUT("/password", auth.ChangePassword)

	adm := api.Group("/admin", requireAdmin)
	adm.GET("/stats", admin.Stats)

	adm.GET("/articles", admin.ListArticles)
	adm.GET("/articles/:id", admin.GetArticle)
	adm.POST("/articles", handler.CreateHandler(adminSvc.CreateArticle))
	adm.PUT("/articles/:id", handler.UpdateHandler(adminSvc.UpdateArticle))
	adm.DELETE("/articles/:id", handler.DeleteHandler(adminSvc.DeleteArticle))

	adm.GET("/categories", pub.ListCategories)
	adm.POST("/categories", handler.CreateHandler(adminSvc.CreateCategory))
	adm.PUT("/categories/:id", handler.UpdateHandler(adminSvc.RenameCategory))
	adm.DELETE("/categories/:id", handler.DeleteHandler(adminSvc.DeleteCategory))

	adm.GET("/tags", pub.ListTags)
	adm.POST("/tags", handler.CreateHandler(adminSvc.CreateTag))
	adm.PUT("/tags/:id", handler.UpdateHandler(adminSvc.RenameTag))
	adm.DELETE("/tags/:id", handler.DeleteHandler(adminSvc.DeleteTag))

	adm.GET("/comments", admin.ListComments)
	adm.DELETE("/comments/:id", handler.DeleteHandler(adminSvc.DeleteComment))

	adm.GET("/moments", pub.ListMoments)
	adm.POST("/moments", handler.CreateHandler(adminSvc.CreateMoment))
	adm.PUT("/moments/:id", handler.UpdateHandler(adminSvc.UpdateMoment))
	adm.DELETE("/moments/:id", handler.DeleteHandler(adminSvc.DeleteMoment))

	adm.GET("/friend-links", pub.ListFriendLinks)
	adm.POST("/friend-links", handler.CreateHandler(adminSvc.CreateFriendLink))
	adm.PUT("/friend-links/:id", handler.UpdateHandler(adminSvc.UpdateFriendLink))
	adm.DELETE("/friend-links/:id", handler.DeleteHandler(adminSvc.DeleteFriendLink))

	adm.GET("/changelogs", pub.ListChangelogs)
	adm.POST("/changelogs", handler.CreateHandler(adminSvc.CreateChangelog))
	adm.PUT("/changelogs/:id", handler.UpdateHandler(adminSvc.UpdateChangelog))
	adm.DELETE("/changelogs/:id", handler.DeleteHandler(adminSvc.DeleteChangelog))

	adm.GET("/projects", pub.ListProjects)
	adm.POST("/projects", handler.CreateHandler(adminSvc.CreateProject))
	adm.PUT("/projects/:id", handler.UpdateHandler(adminSvc.UpdateProject))
	adm.DELETE("/projects/:id", handler.DeleteHandler(adminSvc.DeleteProject))

	adm.PUT("/pages/:key", admin.UpdatePage)
	adm.PUT("/site/notice", admin.UpdateNotice)
	adm.POST("/uploads", admin.UploadTicket)

	r.NoRoute(gzip.Gzip(gzip.DefaultCompression), frontends(cfg))

	return r, nil
}

const adminPrefix = "/admin"

// frontends serves the admin SPA under /admin/ and the blog SPA everywhere else. Unknown /api paths
// and non-GET requests keep the JSON 404.
func frontends(cfg *config.Config) gin.HandlerFunc {
	web := handler.SPA{Root: cfg.WebDir}
	admin := handler.SPA{Root: cfg.AdminDir, Prefix: adminPrefix}
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if (c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead) || strings.HasPrefix(p, "/api/") {
			response.Fail(c, apperr.ErrNotFound)
			return
		}
		switch {
		case cfg.AdminDir != "" && p == adminPrefix:
			c.Redirect(http.StatusMovedPermanently, adminPrefix+"/")
		case cfg.AdminDir != "" && strings.HasPrefix(p, adminPrefix+"/"):
			middleware.SetPageSecurityHeaders(c.Writer.Header())
			admin.Serve(c)
		case cfg.WebDir != "":
			middleware.SetPageSecurityHeaders(c.Writer.Header())
			web.Serve(c)
		default:
			response.Fail(c, apperr.ErrNotFound)
		}
	}
}
