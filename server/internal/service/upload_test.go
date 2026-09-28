package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"blog-server/internal/apperr"
	"blog-server/internal/config"
)

func testUploads(t *testing.T) *Uploads {
	t.Helper()
	u, err := NewUploads(config.OSS{
		Endpoint: "https://oss.example.com", Region: "us-east-1",
		AccessKeyID: "id", SecretAccessKey: "secret", Bucket: "blog",
		PublicBaseURL: "https://oss.example.com/blog", Dir: "images",
	})
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func TestUploadTicket(t *testing.T) {
	now := time.Date(2026, 9, 28, 8, 0, 0, 0, time.UTC)
	tk, err := testUploads(t).Ticket(context.Background(), "Photo.PNG", 1234, now)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(tk.Key, "images/2026/09/28080000-") || !strings.HasSuffix(tk.Key, ".png") {
		t.Errorf("key = %q", tk.Key)
	}
	if tk.URL != "https://oss.example.com/blog/"+tk.Key || tk.Method != "PUT" || tk.Headers["Content-Type"] != "image/png" {
		t.Errorf("unexpected ticket: %+v", tk)
	}

	u, err := url.Parse(tk.UploadURL)
	if err != nil {
		t.Fatal(err)
	}
	if u.Host != "oss.example.com" || u.Path != "/blog/"+tk.Key {
		t.Errorf("upload URL should be path-style, got %s", tk.UploadURL)
	}
	q := u.Query()
	if got := q.Get("X-Amz-SignedHeaders"); got != "content-length;content-type;host" {
		t.Errorf("signed headers = %q", got)
	}
	if q.Get("X-Amz-Expires") != "600" || q.Get("X-Amz-Signature") == "" {
		t.Errorf("missing expiry or signature: %s", u.RawQuery)
	}
}

func TestUploadTicketRejects(t *testing.T) {
	disabled, err := NewUploads(config.OSS{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := disabled.Ticket(context.Background(), "a.png", 10, time.Now()); err == nil {
		t.Error("expected error when storage is not configured")
	}

	u := testUploads(t)
	cases := []struct {
		name string
		size int64
	}{{"evil.svg", 10}, {"page.html", 10}, {"noext", 10}, {"a.png", 0}, {"a.png", 11 << 20}}
	for _, c := range cases {
		var apiErr *apperr.Error
		if _, err := u.Ticket(context.Background(), c.name, c.size, time.Now()); !errors.As(err, &apiErr) || apiErr.Status != 400 {
			t.Errorf("%s (%d bytes): expected 400, got %v", c.name, c.size, err)
		}
	}
}
