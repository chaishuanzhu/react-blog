package service

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"blog-server/internal/apperr"
	"blog-server/internal/config"
	appmail "blog-server/internal/mail"
	"blog-server/internal/model"
	"blog-server/internal/store"
)

const (
	maxCommentRunes  = 1000
	guestbookTitle   = "留言板"
	commentMailRunes = 300
)

var (
	reNickname = regexp.MustCompile(`^[\p{Han}A-Za-z0-9_\- ]{2,16}$`)
	reQQEmail  = regexp.MustCompile(`^([1-9][0-9]{4,11})@qq\.com$`)
)

type Comments struct {
	q        store.Querier
	notifier appmail.Notifier
	cfg      *config.Config
}

func NewComments(q store.Querier, notifier appmail.Notifier, cfg *config.Config) *Comments {
	return &Comments{q: q, notifier: notifier, cfg: cfg}
}

// commentRow matches the column list shared by the comment read queries, so each generated row type converts to it directly.
type commentRow struct {
	ID        uint64
	ArticleID sql.NullInt64
	ParentID  sql.NullInt64
	Nickname  string
	Website   string
	Avatar    string
	Content   string
	IsAdmin   bool
	CreatedAt time.Time
}

func (r commentRow) toModel() model.Comment {
	c := model.Comment{
		ID: r.ID, Nickname: r.Nickname, Website: r.Website, Avatar: r.Avatar,
		Content: r.Content, IsAdmin: r.IsAdmin, CreatedAt: r.CreatedAt,
	}
	if r.ParentID.Valid {
		pid := uint64(r.ParentID.Int64)
		c.ParentID = &pid
	}
	return c
}

type articleRef struct {
	ID    uint64
	Title string
}

func (s *Comments) List(ctx context.Context, articleID uint64, page, pageSize int) ([]model.CommentThread, int64, error) {
	nullArticleID := sql.NullInt64{}
	if articleID != 0 {
		ref, err := s.q.GetPublishedArticleRef(ctx, articleID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, apperr.ErrNotFound
		}
		if err != nil {
			return nil, 0, fmt.Errorf("get article: %w", err)
		}
		nullArticleID = sql.NullInt64{Int64: int64(ref.ID), Valid: true}
	}

	total, err := s.q.CountTopLevelComments(ctx, nullArticleID)
	if err != nil {
		return nil, 0, fmt.Errorf("count comments: %w", err)
	}
	if total == 0 {
		return []model.CommentThread{}, 0, nil
	}
	rows, err := s.q.ListTopLevelComments(ctx, store.ListTopLevelCommentsParams{
		ArticleID: nullArticleID,
		Limit:     int32(pageSize),
		Offset:    int32((page - 1) * pageSize),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list comments: %w", err)
	}

	threads := make([]model.CommentThread, len(rows))
	index := make(map[uint64]int, len(rows))
	ids := make([]sql.NullInt64, len(rows))
	for i, r := range rows {
		threads[i] = model.CommentThread{Comment: commentRow(r).toModel(), Replies: []model.Comment{}}
		index[r.ID] = i
		ids[i] = sql.NullInt64{Int64: int64(r.ID), Valid: true}
	}
	if len(ids) > 0 {
		replies, err := s.q.ListRepliesForParents(ctx, ids)
		if err != nil {
			return nil, 0, fmt.Errorf("list replies: %w", err)
		}
		for _, r := range replies {
			if i, ok := index[uint64(r.ParentID.Int64)]; ok {
				threads[i].Replies = append(threads[i].Replies, commentRow(r).toModel())
			}
		}
	}
	return threads, total, nil
}

func (s *Comments) Create(ctx context.Context, in model.CommentInput, author model.CommentAuthor) (*model.Comment, error) {
	in.Content = strings.TrimSpace(in.Content)
	if in.Content == "" {
		return nil, apperr.BadRequest("content is required")
	}
	if utf8.RuneCountInString(in.Content) > maxCommentRunes {
		return nil, apperr.BadRequest(fmt.Sprintf("content must be at most %d characters", maxCommentRunes))
	}

	nickname, email, website, avatar := author.Nickname, author.Email, author.Website, author.Avatar
	if !author.IsAdmin {
		var err error
		if nickname, email, website, err = s.validateVisitor(in); err != nil {
			return nil, err
		}
		avatar = avatarFor(email)
	}

	var article *articleRef
	if in.ArticleID != 0 {
		ref, err := s.q.GetPublishedArticleRef(ctx, in.ArticleID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.BadRequest("article not found")
		}
		if err != nil {
			return nil, fmt.Errorf("get article: %w", err)
		}
		article = &articleRef{ID: ref.ID, Title: ref.Title}
	}

	var repliedTo *store.GetCommentRow
	threadID := sql.NullInt64{}
	if in.ParentID != 0 {
		parent, err := s.q.GetComment(ctx, in.ParentID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.BadRequest("parent comment not found")
		}
		if err != nil {
			return nil, fmt.Errorf("get parent comment: %w", err)
		}
		if !sameArticle(parent.ArticleID, article) {
			return nil, apperr.BadRequest("parent comment belongs to a different page")
		}
		repliedTo = &parent
		// Threads are one level deep: a reply to a reply joins the top-level thread.
		threadID = sql.NullInt64{Int64: int64(parent.ID), Valid: true}
		if parent.ParentID.Valid {
			threadID = parent.ParentID
		}
	}

	articleID := sql.NullInt64{}
	if article != nil {
		articleID = sql.NullInt64{Int64: int64(article.ID), Valid: true}
	}
	res, err := s.q.CreateComment(ctx, store.CreateCommentParams{
		ArticleID: articleID,
		ParentID:  threadID,
		Nickname:  nickname,
		Email:     email,
		Website:   website,
		Avatar:    avatar,
		Content:   in.Content,
		IsAdmin:   author.IsAdmin,
		Ip:        author.IP,
		UserAgent: truncateString(author.UserAgent, 255),
	})
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("read comment id: %w", err)
	}
	created, err := s.q.GetCommentByID(ctx, uint64(id))
	if err != nil {
		return nil, fmt.Errorf("load created comment: %w", err)
	}
	comment := commentRow(created).toModel()

	s.notify(comment, email, article, repliedTo)
	return &comment, nil
}

func (s *Comments) validateVisitor(in model.CommentInput) (nickname, email, website string, err error) {
	nickname = strings.TrimSpace(in.Nickname)
	if !reNickname.MatchString(nickname) || utf8.RuneCountInString(nickname) > 16 {
		return "", "", "", apperr.BadRequest("nickname must be 2-16 characters of letters, digits, Chinese, _ or -")
	}

	email = strings.ToLower(strings.TrimSpace(in.Email))
	addr, perr := mail.ParseAddress(email)
	if perr != nil || addr.Address != email || len(email) > 128 {
		return "", "", "", apperr.BadRequest("email is invalid")
	}

	website = strings.TrimSpace(in.Website)
	if website != "" {
		u, perr := url.Parse(website)
		if perr != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || len(website) > 255 {
			return "", "", "", apperr.BadRequest("website must be an http(s) URL")
		}
	}

	if s.isReservedIdentity(nickname, email) {
		return "", "", "", apperr.ErrForbidden
	}
	return nickname, email, website, nil
}

func (s *Comments) isReservedIdentity(nickname, email string) bool {
	if s.cfg.AdminNickname != "" && strings.EqualFold(nickname, s.cfg.AdminNickname) {
		return true
	}
	for _, reserved := range []string{s.cfg.AdminEmail, s.cfg.NotifyEmail, s.cfg.SMTP.From} {
		if reserved != "" && strings.EqualFold(email, reserved) {
			return true
		}
	}
	return false
}

func (s *Comments) notify(c model.Comment, authorEmail string, article *articleRef, repliedTo *store.GetCommentRow) {
	target, link := guestbookTitle, s.cfg.SiteURL+"/msg"
	if article != nil {
		target = article.Title
		link = fmt.Sprintf("%s/post/%d", s.cfg.SiteURL, article.ID)
	}
	link += fmt.Sprintf("#comment-%d", c.ID)

	base := appmail.CommentMail{
		SiteName: s.cfg.SiteName,
		Target:   target,
		Author:   c.Nickname,
		Content:  truncateString(c.Content, commentMailRunes),
		Link:     link,
	}
	if repliedTo != nil {
		base.ParentShown = true
		base.ParentName = repliedTo.Nickname
		base.Parent = truncateString(repliedTo.Content, commentMailRunes)
	}

	if !c.IsAdmin && s.cfg.NotifyEmail != "" {
		m := base
		m.Heading = fmt.Sprintf("「%s」有新评论", target)
		if repliedTo != nil {
			m.Heading = fmt.Sprintf("%s 回复了 %s", c.Nickname, repliedTo.Nickname)
		}
		s.send(s.cfg.NotifyEmail, fmt.Sprintf("[%s] %s", s.cfg.SiteName, m.Heading), m)
	}

	if repliedTo == nil || repliedTo.IsAdmin || repliedTo.Email == "" ||
		strings.EqualFold(repliedTo.Email, authorEmail) || strings.EqualFold(repliedTo.Email, s.cfg.NotifyEmail) {
		return
	}
	m := base
	m.Heading = fmt.Sprintf("%s 回复了你的评论", c.Nickname)
	s.send(repliedTo.Email, fmt.Sprintf("[%s] 你的评论有了新回复", s.cfg.SiteName), m)
}

func (s *Comments) send(to, subject string, data appmail.CommentMail) {
	html, err := appmail.RenderComment(data)
	if err != nil {
		slog.Error("render comment mail", "err", err)
		return
	}
	s.notifier.Enqueue(appmail.Message{To: to, Subject: subject, HTML: html})
}

func sameArticle(parentArticle sql.NullInt64, article *articleRef) bool {
	if article == nil {
		return !parentArticle.Valid
	}
	return parentArticle.Valid && uint64(parentArticle.Int64) == article.ID
}

func avatarFor(email string) string {
	if m := reQQEmail.FindStringSubmatch(email); m != nil {
		return "https://q1.qlogo.cn/g?b=qq&nk=" + m[1] + "&s=100"
	}
	sum := md5.Sum([]byte(email))
	return "https://cravatar.cn/avatar/" + hex.EncodeToString(sum[:]) + "?d=identicon&s=100"
}

func truncateString(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	return string([]rune(s)[:maxRunes]) + "…"
}
