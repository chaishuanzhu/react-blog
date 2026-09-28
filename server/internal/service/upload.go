package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"blog-server/internal/apperr"
	"blog-server/internal/config"
)

const (
	uploadURLTTL   = 10 * time.Minute
	maxUploadBytes = 10 << 20
)

// SVG is excluded because it can carry scripts when served from a public bucket.
var allowedImageTypes = map[string]string{
	".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
	".gif": "image/gif", ".webp": "image/webp", ".avif": "image/avif",
}

type UploadTicket struct {
	Method    string            `json:"method"`
	UploadURL string            `json:"uploadUrl"`
	Headers   map[string]string `json:"headers"`
	Key       string            `json:"key"`
	URL       string            `json:"url"`
	ExpiresAt time.Time         `json:"expiresAt"`
	MaxBytes  int64             `json:"maxBytes"`
}

type Uploads struct {
	cfg    config.OSS
	client *minio.Client
}

func NewUploads(cfg config.OSS) (*Uploads, error) {
	s := &Uploads{cfg: cfg}
	if !cfg.Enabled() {
		return s, nil
	}
	endpoint, err := url.Parse(cfg.Endpoint)
	if err != nil || endpoint.Host == "" {
		return nil, fmt.Errorf("OSS_ENDPOINT must be a URL such as https://oss.example.com")
	}
	s.client, err = minio.New(endpoint.Host, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure:       endpoint.Scheme == "https",
		Region:       cfg.Region,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("init object storage client: %w", err)
	}
	return s, nil
}

// Ticket returns a presigned PUT URL. Content-Type and Content-Length are part of the signature,
// so the browser can only upload a file of the declared type and size to this exact key.
func (s *Uploads) Ticket(ctx context.Context, filename string, size int64, now time.Time) (*UploadTicket, error) {
	if s.client == nil {
		return nil, &apperr.Error{Status: 503, Code: "UPLOAD_DISABLED", Message: "image upload is not configured"}
	}
	ext := strings.ToLower(path.Ext(strings.TrimSpace(filename)))
	contentType, ok := allowedImageTypes[ext]
	if !ok {
		return nil, apperr.BadRequest("only jpg, png, gif, webp and avif images are allowed")
	}
	if size <= 0 || size > maxUploadBytes {
		return nil, apperr.BadRequest(fmt.Sprintf("size must be between 1 and %d bytes", maxUploadBytes))
	}

	id := make([]byte, 8)
	_, _ = rand.Read(id)
	utc := now.UTC()
	key := path.Join(s.cfg.Dir, utc.Format("2006/01"), utc.Format("02150405")+"-"+hex.EncodeToString(id)+ext)

	signed := http.Header{}
	signed.Set("Content-Type", contentType)
	signed.Set("Content-Length", strconv.FormatInt(size, 10))
	u, err := s.client.PresignHeader(ctx, http.MethodPut, s.cfg.Bucket, key, uploadURLTTL, url.Values{}, signed)
	if err != nil {
		return nil, fmt.Errorf("presign upload: %w", err)
	}

	return &UploadTicket{
		Method:    http.MethodPut,
		UploadURL: u.String(),
		Headers:   map[string]string{"Content-Type": contentType},
		Key:       key,
		URL:       s.cfg.PublicBaseURL + "/" + key,
		ExpiresAt: utc.Add(uploadURLTTL),
		MaxBytes:  maxUploadBytes,
	}, nil
}
