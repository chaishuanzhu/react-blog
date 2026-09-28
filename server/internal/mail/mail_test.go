package mail

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"blog-server/internal/config"
)

func TestBuildMessage(t *testing.T) {
	cfg := config.SMTP{From: "no_reply@example.com", FromName: "飞鸟小站"}
	raw := string(buildMessage(cfg, Message{To: "a@example.com", Subject: "新评论", HTML: "<p>你好</p>"}, time.Unix(0, 0)))

	head, body, ok := strings.Cut(raw, "\r\n\r\n")
	if !ok {
		t.Fatal("message has no header/body separator")
	}
	for _, want := range []string{"To: a@example.com", "Subject: =?UTF-8?b?", "Message-ID: <", "@example.com>", "Content-Transfer-Encoding: base64"} {
		if !strings.Contains(head, want) {
			t.Errorf("header missing %q in:\n%s", want, head)
		}
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(body, "\r\n", ""))
	if err != nil || string(decoded) != "<p>你好</p>" {
		t.Errorf("body = %q, err = %v", decoded, err)
	}
}

func TestRenderCommentEscapesHTML(t *testing.T) {
	out, err := RenderComment(CommentMail{Author: "x", Content: `<script>alert(1)</script>`, Link: "http://localhost/post?title=a"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "<script>") {
		t.Error("comment content was not escaped")
	}
}
