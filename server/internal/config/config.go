package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env            string
	HTTPAddr       string
	DatabaseDSN    string
	JWTSecret      string
	JWTTTL         time.Duration
	AllowedOrigins []string
	TrustedProxies []string
	AdminEmail     string
	AdminPassword  string
	AdminNickname  string

	SiteName string
	SiteURL  string

	// Built frontends served by this process; empty disables them (development uses the dev servers).
	WebDir   string
	AdminDir string

	// Serve the OpenAPI spec and Swagger UI under /api/docs/.
	APIDocs bool

	SMTP        SMTP
	NotifyEmail string

	OSS OSS

	CommentRatePerMinute int
	CommentBurst         int
}

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

func (s SMTP) Enabled() bool {
	return s.Host != "" && s.Username != "" && s.Password != ""
}

// OSS configures an S3-compatible object store (RustFS, MinIO, AWS S3, R2, ...).
type OSS struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	PublicBaseURL   string
	Dir             string
}

func (o OSS) Enabled() bool {
	return o.Endpoint != "" && o.AccessKeyID != "" && o.SecretAccessKey != "" && o.Bucket != ""
}

func Load() (*Config, error) {
	ttl, err := time.ParseDuration(getenv("JWT_TTL", "168h"))
	if err != nil {
		return nil, fmt.Errorf("parse JWT_TTL: %w", err)
	}
	smtpPort, err := getInt("SMTP_PORT", 465)
	if err != nil {
		return nil, err
	}
	ratePerMinute, err := getInt("COMMENT_RATE_PER_MINUTE", 3)
	if err != nil {
		return nil, err
	}
	burst, err := getInt("COMMENT_BURST", 3)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Env:            getenv("APP_ENV", "development"),
		HTTPAddr:       getenv("HTTP_ADDR", ":8080"),
		DatabaseDSN:    os.Getenv("DATABASE_DSN"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTTTL:         ttl,
		AllowedOrigins: splitList(getenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001")),
		TrustedProxies: splitList(getenv("TRUSTED_PROXIES", "127.0.0.1,::1")),
		AdminEmail:     strings.TrimSpace(os.Getenv("ADMIN_EMAIL")),
		AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
		AdminNickname:  strings.TrimSpace(os.Getenv("ADMIN_NICKNAME")),

		SiteName: getenv("SITE_NAME", "Blog"),
		SiteURL:  strings.TrimRight(getenv("SITE_URL", "http://localhost:3000"), "/"),

		WebDir:   os.Getenv("WEB_DIR"),
		AdminDir: os.Getenv("ADMIN_DIR"),

		SMTP: SMTP{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     smtpPort,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     getenv("SMTP_FROM", os.Getenv("SMTP_USERNAME")),
			FromName: os.Getenv("SMTP_FROM_NAME"),
		},
		NotifyEmail: strings.TrimSpace(os.Getenv("NOTIFY_EMAIL")),

		OSS: OSS{
			Endpoint:        strings.TrimRight(os.Getenv("OSS_ENDPOINT"), "/"),
			Region:          getenv("OSS_REGION", "us-east-1"),
			AccessKeyID:     os.Getenv("OSS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("OSS_SECRET_ACCESS_KEY"),
			Bucket:          os.Getenv("OSS_BUCKET"),
			Dir:             strings.Trim(getenv("OSS_DIR", "images"), "/"),
		},

		CommentRatePerMinute: ratePerMinute,
		CommentBurst:         burst,
	}

	if cfg.APIDocs, err = getBool("API_DOCS", !cfg.IsProduction()); err != nil {
		return nil, err
	}

	cfg.OSS.PublicBaseURL = strings.TrimRight(getenv("OSS_PUBLIC_BASE_URL", cfg.OSS.Endpoint+"/"+cfg.OSS.Bucket), "/")

	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN is required")
	}
	if cfg.IsProduction() && len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
	}
	return cfg, nil
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) (int, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", key)
	}
	return n, nil
}

func getBool(key string, fallback bool) (bool, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback, nil
	}
	b, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s must be true or false", key)
	}
	return b, nil
}

func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
