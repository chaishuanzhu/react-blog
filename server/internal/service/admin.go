package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"blog-server/internal/apperr"
	"blog-server/internal/model"
	"blog-server/internal/store"
)

const (
	maxTagsPerArticle = 20
	maxMomentImages   = 9
	maxContentBytes   = 4 << 20
)

type Admin struct {
	db *sql.DB
	q  *store.Queries
}

func NewAdmin(db *sql.DB) *Admin {
	return &Admin{db: db, q: store.New(db)}
}

// ---- articles ----

func (s *Admin) ListArticles(ctx context.Context, f model.AdminArticleFilter) ([]model.AdminArticle, int64, error) {
	status := store.NullArticlesStatus{}
	if f.Status != "" {
		if f.Status != string(store.ArticlesStatusDraft) && f.Status != string(store.ArticlesStatusPublished) {
			return nil, 0, apperr.BadRequest("status must be draft or published")
		}
		status = store.NullArticlesStatus{ArticlesStatus: store.ArticlesStatus(f.Status), Valid: true}
	}
	keyword := nullString(f.Keyword)
	if keyword.Valid {
		keyword.String = likeContains(keyword.String)
	}
	categoryID := sql.NullInt64{Int64: int64(f.CategoryID), Valid: f.CategoryID != 0}
	tagID := sql.NullInt64{Int64: int64(f.TagID), Valid: f.TagID != 0}

	total, err := s.q.AdminCountArticles(ctx, store.AdminCountArticlesParams{
		Status: status, Keyword: keyword, CategoryID: categoryID, TagID: tagID,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("count articles: %w", err)
	}
	rows, err := s.q.AdminListArticles(ctx, store.AdminListArticlesParams{
		Status: status, Keyword: keyword, CategoryID: categoryID, TagID: tagID,
		Limit: int32(f.PageSize), Offset: int32((f.Page - 1) * f.PageSize),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list articles: %w", err)
	}

	ids := make([]uint64, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	tags, err := (&Public{q: s.q}).tagsByArticle(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	items := make([]model.AdminArticle, len(rows))
	for i, r := range rows {
		items[i] = adminArticle(store.AdminGetArticleRow(r), tags[r.ID], false)
	}
	return items, total, nil
}

func (s *Admin) GetArticle(ctx context.Context, id uint64) (*model.AdminArticle, error) {
	r, err := s.q.AdminGetArticle(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get article: %w", err)
	}
	tags, err := (&Public{q: s.q}).tagsByArticle(ctx, []uint64{id})
	if err != nil {
		return nil, err
	}
	a := adminArticle(r, tags[id], true)
	return &a, nil
}

func (s *Admin) CreateArticle(ctx context.Context, in model.ArticleInput) (uint64, error) {
	p, tagIDs, err := validateArticle(in)
	if err != nil {
		return 0, err
	}
	var id uint64
	err = s.inTx(ctx, func(q *store.Queries) error {
		res, err := q.CreateArticle(ctx, store.CreateArticleParams{
			Title: p.Title, Content: p.Content, CategoryID: p.CategoryID,
			Status: p.Status, PublishedAt: p.PublishedAt,
		})
		if err != nil {
			return err
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		id = uint64(lastID)
		return setArticleTags(ctx, q, id, tagIDs)
	})
	if err != nil {
		return 0, articleWriteError(err)
	}
	return id, nil
}

func (s *Admin) UpdateArticle(ctx context.Context, id uint64, in model.ArticleInput) error {
	p, tagIDs, err := validateArticle(in)
	if err != nil {
		return err
	}
	err = s.inTx(ctx, func(q *store.Queries) error {
		n, err := q.UpdateArticle(ctx, store.UpdateArticleParams{
			Title: p.Title, Content: p.Content, CategoryID: p.CategoryID,
			Status: p.Status, PublishedAt: p.PublishedAt, ID: id,
		})
		if err != nil {
			return err
		}
		if n == 0 {
			return apperr.ErrNotFound
		}
		if err := q.DeleteArticleTags(ctx, id); err != nil {
			return err
		}
		return setArticleTags(ctx, q, id, tagIDs)
	})
	return articleWriteError(err)
}

func (s *Admin) DeleteArticle(ctx context.Context, id uint64) error {
	return deleted(s.q.DeleteArticle(ctx, id))
}

func validateArticle(in model.ArticleInput) (store.CreateArticleParams, []uint64, error) {
	var p store.CreateArticleParams
	var err error
	if p.Title, err = requiredText("title", in.Title, 200); err != nil {
		return p, nil, err
	}

	switch in.Status {
	case "", string(store.ArticlesStatusDraft):
		p.Status = store.ArticlesStatusDraft
	case string(store.ArticlesStatusPublished):
		p.Status = store.ArticlesStatusPublished
	default:
		return p, nil, apperr.BadRequest("status must be draft or published")
	}
	p.Content = in.Content
	if len(p.Content) > maxContentBytes {
		return p, nil, apperr.BadRequest("content is too large")
	}
	if p.Status == store.ArticlesStatusPublished && strings.TrimSpace(p.Content) == "" {
		return p, nil, apperr.BadRequest("content is required to publish")
	}

	p.PublishedAt = time.Now().UTC()
	if in.PublishedAt != nil {
		p.PublishedAt = in.PublishedAt.UTC()
	}
	if in.CategoryID != nil && *in.CategoryID != 0 {
		p.CategoryID = sql.NullInt64{Int64: int64(*in.CategoryID), Valid: true}
	}

	seen := map[uint64]bool{}
	var tagIDs []uint64
	for _, id := range in.TagIDs {
		if id != 0 && !seen[id] {
			seen[id] = true
			tagIDs = append(tagIDs, id)
		}
	}
	if len(tagIDs) > maxTagsPerArticle {
		return p, nil, apperr.BadRequest(fmt.Sprintf("an article can have at most %d tags", maxTagsPerArticle))
	}
	return p, tagIDs, nil
}

func setArticleTags(ctx context.Context, q *store.Queries, articleID uint64, tagIDs []uint64) error {
	for _, tagID := range tagIDs {
		if err := q.AddArticleTag(ctx, store.AddArticleTagParams{ArticleID: articleID, TagID: tagID}); err != nil {
			return err
		}
	}
	return nil
}

func articleWriteError(err error) error {
	switch {
	case err == nil:
		return nil
	case isMySQLError(err, mysqlFKViolation):
		return apperr.BadRequest("category or tag does not exist")
	}
	return err
}

func adminArticle(r store.AdminGetArticleRow, tags []model.TagRef, withContent bool) model.AdminArticle {
	a := model.AdminArticle{
		ID: r.ID, Title: r.Title, Summary: Summarize(r.Content, summaryLength),
		Status: string(r.Status), Category: categoryRef(r.CategoryID, r.CategoryName),
		Tags: tagsOrEmpty(tags), PublishedAt: r.PublishedAt, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
	if withContent {
		a.Content = r.Content
	}
	return a
}

// ---- categories & tags ----

func (s *Admin) CreateCategory(ctx context.Context, in model.NameInput) (uint64, error) {
	name, err := requiredText("name", in.Name, 64)
	if err != nil {
		return 0, err
	}
	return created(s.q.CreateCategory(ctx, name))
}

func (s *Admin) RenameCategory(ctx context.Context, id uint64, in model.NameInput) error {
	name, err := requiredText("name", in.Name, 64)
	if err != nil {
		return err
	}
	return updated(s.q.RenameCategory(ctx, store.RenameCategoryParams{Name: name, ID: id}))
}

func (s *Admin) DeleteCategory(ctx context.Context, id uint64) error {
	return deleted(s.q.DeleteCategory(ctx, id))
}

func (s *Admin) CreateTag(ctx context.Context, in model.NameInput) (uint64, error) {
	name, err := requiredText("name", in.Name, 64)
	if err != nil {
		return 0, err
	}
	return created(s.q.CreateTag(ctx, name))
}

func (s *Admin) RenameTag(ctx context.Context, id uint64, in model.NameInput) error {
	name, err := requiredText("name", in.Name, 64)
	if err != nil {
		return err
	}
	return updated(s.q.RenameTag(ctx, store.RenameTagParams{Name: name, ID: id}))
}

func (s *Admin) DeleteTag(ctx context.Context, id uint64) error {
	return deleted(s.q.DeleteTag(ctx, id))
}

// ---- comments ----

func (s *Admin) ListComments(ctx context.Context, page, pageSize int) ([]model.AdminComment, int64, error) {
	total, err := s.q.AdminCountComments(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count comments: %w", err)
	}
	rows, err := s.q.AdminListComments(ctx, store.AdminListCommentsParams{
		Limit: int32(pageSize), Offset: int32((page - 1) * pageSize),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list comments: %w", err)
	}
	items := make([]model.AdminComment, len(rows))
	for i, r := range rows {
		c := model.AdminComment{
			ID: r.ID, Nickname: r.Nickname, Email: r.Email, Website: r.Website, Avatar: r.Avatar,
			Content: r.Content, IsAdmin: r.IsAdmin, IP: r.Ip, CreatedAt: r.CreatedAt,
		}
		if r.ParentID.Valid {
			pid := uint64(r.ParentID.Int64)
			c.ParentID = &pid
		}
		if r.ArticleID.Valid {
			c.Article = &model.ArticleLink{ID: uint64(r.ArticleID.Int64), Title: r.ArticleTitle.String}
		}
		items[i] = c
	}
	return items, total, nil
}

func (s *Admin) DeleteComment(ctx context.Context, id uint64) error {
	return deleted(s.q.DeleteComment(ctx, id))
}

// ---- moments ----

func validateMoment(in model.MomentInput) (string, json.RawMessage, time.Time, error) {
	content, err := requiredText("content", in.Content, 2000)
	if err != nil {
		return "", nil, time.Time{}, err
	}
	images := []string{}
	for _, img := range in.Images {
		if strings.TrimSpace(img) == "" {
			continue
		}
		u, err := requiredURL("image", img)
		if err != nil {
			return "", nil, time.Time{}, err
		}
		images = append(images, u)
	}
	if len(images) > maxMomentImages {
		return "", nil, time.Time{}, apperr.BadRequest(fmt.Sprintf("at most %d images", maxMomentImages))
	}
	raw, _ := json.Marshal(images)
	at := time.Now().UTC()
	if in.CreatedAt != nil {
		at = in.CreatedAt.UTC()
	}
	return content, raw, at, nil
}

func (s *Admin) CreateMoment(ctx context.Context, in model.MomentInput) (uint64, error) {
	content, images, at, err := validateMoment(in)
	if err != nil {
		return 0, err
	}
	return created(s.q.CreateMoment(ctx, store.CreateMomentParams{Content: content, Images: images, CreatedAt: at}))
}

func (s *Admin) UpdateMoment(ctx context.Context, id uint64, in model.MomentInput) error {
	content, images, at, err := validateMoment(in)
	if err != nil {
		return err
	}
	return updated(s.q.UpdateMoment(ctx, store.UpdateMomentParams{Content: content, Images: images, CreatedAt: at, ID: id}))
}

func (s *Admin) DeleteMoment(ctx context.Context, id uint64) error {
	return deleted(s.q.DeleteMoment(ctx, id))
}

// ---- friend links ----

func validateFriendLink(in model.FriendLinkInput) (store.CreateFriendLinkParams, error) {
	var p store.CreateFriendLinkParams
	var err error
	if p.Name, err = requiredText("name", in.Name, 64); err != nil {
		return p, err
	}
	if p.Url, err = requiredURL("url", in.URL); err != nil {
		return p, err
	}
	if p.Avatar, err = optionalURL("avatar", in.Avatar); err != nil {
		return p, err
	}
	if p.Description, err = optionalText("description", in.Description, 255); err != nil {
		return p, err
	}
	return p, nil
}

func (s *Admin) CreateFriendLink(ctx context.Context, in model.FriendLinkInput) (uint64, error) {
	p, err := validateFriendLink(in)
	if err != nil {
		return 0, err
	}
	return created(s.q.CreateFriendLink(ctx, p))
}

func (s *Admin) UpdateFriendLink(ctx context.Context, id uint64, in model.FriendLinkInput) error {
	p, err := validateFriendLink(in)
	if err != nil {
		return err
	}
	return updated(s.q.UpdateFriendLink(ctx, store.UpdateFriendLinkParams{
		Name: p.Name, Url: p.Url, Avatar: p.Avatar, Description: p.Description, ID: id,
	}))
}

func (s *Admin) DeleteFriendLink(ctx context.Context, id uint64) error {
	return deleted(s.q.DeleteFriendLink(ctx, id))
}

// ---- changelogs ----

func validateChangelog(in model.ChangelogInput) (json.RawMessage, time.Time, error) {
	items := []string{}
	for _, item := range in.Items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if utf8.RuneCountInString(item) > 500 {
			return nil, time.Time{}, apperr.BadRequest("each item must be at most 500 characters")
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, time.Time{}, apperr.BadRequest("items must contain at least one entry")
	}
	raw, _ := json.Marshal(items)
	at := time.Now().UTC()
	if in.LoggedAt != nil {
		at = in.LoggedAt.UTC()
	}
	return raw, at, nil
}

func (s *Admin) CreateChangelog(ctx context.Context, in model.ChangelogInput) (uint64, error) {
	items, at, err := validateChangelog(in)
	if err != nil {
		return 0, err
	}
	return created(s.q.CreateChangelog(ctx, store.CreateChangelogParams{Items: items, LoggedAt: at}))
}

func (s *Admin) UpdateChangelog(ctx context.Context, id uint64, in model.ChangelogInput) error {
	items, at, err := validateChangelog(in)
	if err != nil {
		return err
	}
	return updated(s.q.UpdateChangelog(ctx, store.UpdateChangelogParams{Items: items, LoggedAt: at, ID: id}))
}

func (s *Admin) DeleteChangelog(ctx context.Context, id uint64) error {
	return deleted(s.q.DeleteChangelog(ctx, id))
}

// ---- projects ----

func validateProject(in model.ProjectInput) (store.CreateProjectParams, error) {
	var p store.CreateProjectParams
	var err error
	if p.Name, err = requiredText("name", in.Name, 64); err != nil {
		return p, err
	}
	if p.Description, err = optionalText("description", in.Description, 255); err != nil {
		return p, err
	}
	if p.Cover, err = optionalURL("cover", in.Cover); err != nil {
		return p, err
	}
	if p.Url, err = optionalURL("url", in.URL); err != nil {
		return p, err
	}
	p.SortOrder = in.SortOrder
	return p, nil
}

func (s *Admin) CreateProject(ctx context.Context, in model.ProjectInput) (uint64, error) {
	p, err := validateProject(in)
	if err != nil {
		return 0, err
	}
	return created(s.q.CreateProject(ctx, p))
}

func (s *Admin) UpdateProject(ctx context.Context, id uint64, in model.ProjectInput) error {
	p, err := validateProject(in)
	if err != nil {
		return err
	}
	return updated(s.q.UpdateProject(ctx, store.UpdateProjectParams{
		Name: p.Name, Description: p.Description, Cover: p.Cover, Url: p.Url, SortOrder: p.SortOrder, ID: id,
	}))
}

func (s *Admin) DeleteProject(ctx context.Context, id uint64) error {
	return deleted(s.q.DeleteProject(ctx, id))
}

// ---- pages, notice, stats ----

func (s *Admin) UpdatePage(ctx context.Context, key string, in model.ContentInput) error {
	dbKey, ok := pageKeys[key]
	if !ok {
		return apperr.ErrNotFound
	}
	if len(in.Content) > maxContentBytes {
		return apperr.BadRequest("content is too large")
	}
	return updated(s.q.UpdatePage(ctx, store.UpdatePageParams{Content: in.Content, PageKey: dbKey}))
}

func (s *Admin) UpdateNotice(ctx context.Context, in model.NoticeInput) error {
	notice, err := optionalText("notice", in.Notice, 255)
	if err != nil {
		return err
	}
	return s.q.UpdateNotice(ctx, notice)
}

func (s *Admin) Stats(ctx context.Context) (*model.AdminStats, error) {
	r, err := s.q.AdminStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("admin stats: %w", err)
	}
	return &model.AdminStats{
		PublishedCount: r.PublishedCount, DraftCount: r.DraftCount, CategoryCount: r.CategoryCount,
		TagCount: r.TagCount, CommentCount: r.CommentCount, MomentCount: r.MomentCount,
		FriendLinkCount: r.FriendLinkCount, ProjectCount: r.ProjectCount, ViewCount: r.ViewCount,
	}, nil
}

// ---- helpers ----

func (s *Admin) inTx(ctx context.Context, fn func(q *store.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(s.q.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func created(res sql.Result, err error) (uint64, error) {
	if isMySQLError(err, mysqlDuplicateEntry) {
		return 0, apperr.Conflict("name already exists")
	}
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func updated(n int64, err error) error {
	if isMySQLError(err, mysqlDuplicateEntry) {
		return apperr.Conflict("name already exists")
	}
	if err != nil {
		return err
	}
	if n == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func deleted(n int64, err error) error {
	if err != nil {
		return err
	}
	if n == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
