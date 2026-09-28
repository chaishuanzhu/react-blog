package service

import (
	"database/sql"
	"errors"
	"sort"
	"testing"

	"blog-server/internal/apperr"
	"blog-server/internal/config"
	appmail "blog-server/internal/mail"
	"blog-server/internal/model"
	"blog-server/internal/store"
)

type recordingNotifier struct{ to []string }

func (r *recordingNotifier) Enqueue(m appmail.Message) { r.to = append(r.to, m.To) }

func testComments(n appmail.Notifier) *Comments {
	return NewComments(nil, n, &config.Config{
		SiteName:      "飞鸟小站",
		SiteURL:       "https://blog.example.com",
		AdminNickname: "飞鸟",
		AdminEmail:    "admin@example.com",
		NotifyEmail:   "notify@example.com",
		SMTP:          config.SMTP{From: "no_reply@example.com"},
	})
}

func TestValidateVisitor(t *testing.T) {
	s := testComments(nil)
	valid := model.CommentInput{Nickname: "访客_1", Email: " Guest@Example.com ", Website: "https://g.example.com", Content: "hi"}

	nick, email, site, err := s.validateVisitor(valid)
	if err != nil || nick != "访客_1" || email != "guest@example.com" || site != "https://g.example.com" {
		t.Fatalf("valid input rejected or not normalized: %q %q %q %v", nick, email, site, err)
	}

	cases := map[string]func(*model.CommentInput){
		"short nickname":     func(in *model.CommentInput) { in.Nickname = "a" },
		"symbol nickname":    func(in *model.CommentInput) { in.Nickname = "<b>x</b>" },
		"bad email":          func(in *model.CommentInput) { in.Email = "not-an-email" },
		"email with name":    func(in *model.CommentInput) { in.Email = "Bob <bob@example.com>" },
		"javascript website": func(in *model.CommentInput) { in.Website = "javascript:alert(1)" },
	}
	for name, mutate := range cases {
		in := valid
		mutate(&in)
		if _, _, _, err := s.validateVisitor(in); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}

	reserved := map[string]func(*model.CommentInput){
		"admin nickname": func(in *model.CommentInput) { in.Nickname = "飞鸟" },
		"admin email":    func(in *model.CommentInput) { in.Email = "ADMIN@example.com" },
		"notify email":   func(in *model.CommentInput) { in.Email = "notify@example.com" },
	}
	for name, mutate := range reserved {
		in := valid
		mutate(&in)
		if _, _, _, err := s.validateVisitor(in); !errors.Is(err, apperr.ErrForbidden) {
			t.Errorf("%s: expected ErrForbidden, got %v", name, err)
		}
	}
}

func TestAvatarFor(t *testing.T) {
	if got := avatarFor("12345678@qq.com"); got != "https://q1.qlogo.cn/g?b=qq&nk=12345678&s=100" {
		t.Errorf("qq avatar = %q", got)
	}
	if got := avatarFor("a@example.com"); got != "https://cravatar.cn/avatar/b418773a2c51fb9777a1648346fa7394?d=identicon&s=100" {
		t.Errorf("cravatar avatar = %q", got)
	}
}

func TestNotifyRecipients(t *testing.T) {
	article := &articleRef{ID: 1, Title: "Hello"}
	visitorParent := &store.GetCommentRow{ID: 7, Nickname: "甲", Email: "jia@example.com", Content: "first"}
	adminParent := &store.GetCommentRow{ID: 8, Nickname: "飞鸟", Email: "admin@example.com", IsAdmin: true}

	tests := []struct {
		name        string
		isAdmin     bool
		authorEmail string
		parent      *store.GetCommentRow
		want        []string
	}{
		{"new visitor comment notifies admin", false, "yi@example.com", nil, []string{"notify@example.com"}},
		{"visitor reply notifies admin and parent", false, "yi@example.com", visitorParent, []string{"jia@example.com", "notify@example.com"}},
		{"replying to yourself skips parent", false, "jia@example.com", visitorParent, []string{"notify@example.com"}},
		{"visitor reply to admin notifies admin once", false, "yi@example.com", adminParent, []string{"notify@example.com"}},
		{"admin reply notifies only the visitor", true, "admin@example.com", visitorParent, []string{"jia@example.com"}},
		{"admin top-level comment sends nothing", true, "admin@example.com", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &recordingNotifier{}
			s := testComments(rec)
			c := model.Comment{ID: 99, Nickname: "乙", Content: "hi", IsAdmin: tt.isAdmin}
			s.notify(c, tt.authorEmail, article, tt.parent)

			sort.Strings(rec.to)
			if len(rec.to) != len(tt.want) {
				t.Fatalf("recipients = %v, want %v", rec.to, tt.want)
			}
			for i := range tt.want {
				if rec.to[i] != tt.want[i] {
					t.Fatalf("recipients = %v, want %v", rec.to, tt.want)
				}
			}
		})
	}
}

func TestSameArticle(t *testing.T) {
	a := &articleRef{ID: 3}
	if !sameArticle(sql.NullInt64{}, nil) || sameArticle(sql.NullInt64{}, a) ||
		!sameArticle(sql.NullInt64{Int64: 3, Valid: true}, a) || sameArticle(sql.NullInt64{Int64: 4, Valid: true}, a) {
		t.Error("sameArticle mismatch")
	}
}
