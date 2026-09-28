package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"blog-server/internal/apperr"
	"blog-server/internal/model"
	"blog-server/internal/store"
)

const summaryLength = 150

var pageKeys = map[string]store.PagesPageKey{
	"about-site": store.PagesPageKeyAboutSite,
	"about-me":   store.PagesPageKeyAboutMe,
}

type Public struct {
	q store.Querier
}

func NewPublic(q store.Querier) *Public {
	return &Public{q: q}
}

func (s *Public) ListArticles(ctx context.Context, f model.ArticleFilter) ([]model.ArticleSummary, int64, error) {
	keyword, category, tag := nullString(f.Keyword), nullString(f.Category), nullString(f.Tag)
	if keyword.Valid {
		keyword.String = likeContains(keyword.String)
	}

	total, err := s.q.CountPublishedArticles(ctx, store.CountPublishedArticlesParams{
		Keyword: keyword, Category: category, Tag: tag,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("count articles: %w", err)
	}
	if total == 0 {
		return []model.ArticleSummary{}, 0, nil
	}

	rows, err := s.q.ListPublishedArticles(ctx, store.ListPublishedArticlesParams{
		Keyword: keyword, Category: category, Tag: tag,
		Limit:  int32(f.PageSize),
		Offset: int32((f.Page - 1) * f.PageSize),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list articles: %w", err)
	}

	ids := make([]uint64, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	tags, err := s.tagsByArticle(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	items := make([]model.ArticleSummary, len(rows))
	for i, r := range rows {
		items[i] = model.ArticleSummary{
			ID:          r.ID,
			Title:       r.Title,
			Summary:     Summarize(r.Content, summaryLength),
			Category:    categoryRef(r.CategoryID, r.CategoryName),
			Tags:        tagsOrEmpty(tags[r.ID]),
			PublishedAt: r.PublishedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	}
	return items, total, nil
}

func (s *Public) GetArticle(ctx context.Context, id uint64) (*model.ArticleDetail, error) {
	r, err := s.q.GetPublishedArticle(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get article: %w", err)
	}
	tags, err := s.tagsByArticle(ctx, []uint64{r.ID})
	if err != nil {
		return nil, err
	}
	return &model.ArticleDetail{
		ArticleSummary: model.ArticleSummary{
			ID:          r.ID,
			Title:       r.Title,
			Summary:     Summarize(r.Content, summaryLength),
			Category:    categoryRef(r.CategoryID, r.CategoryName),
			Tags:        tagsOrEmpty(tags[r.ID]),
			PublishedAt: r.PublishedAt,
			UpdatedAt:   r.UpdatedAt,
		},
		Content: r.Content,
	}, nil
}

func (s *Public) ListCategories(ctx context.Context) (*model.CategoryList, error) {
	rows, err := s.q.ListCategoriesWithCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	uncategorized, err := s.q.CountUncategorizedPublished(ctx)
	if err != nil {
		return nil, fmt.Errorf("count uncategorized: %w", err)
	}
	items := make([]model.Category, len(rows))
	for i, r := range rows {
		items[i] = model.Category{ID: r.ID, Name: r.Name, ArticleCount: r.ArticleCount, CreatedAt: r.CreatedAt}
	}
	return &model.CategoryList{Items: items, UncategorizedCount: uncategorized}, nil
}

func (s *Public) ListTags(ctx context.Context) ([]model.Tag, error) {
	rows, err := s.q.ListTagsWithCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	items := make([]model.Tag, len(rows))
	for i, r := range rows {
		items[i] = model.Tag{ID: r.ID, Name: r.Name, ArticleCount: r.ArticleCount, CreatedAt: r.CreatedAt}
	}
	return items, nil
}

func (s *Public) ListMoments(ctx context.Context, page, pageSize int) ([]model.Moment, int64, error) {
	total, err := s.q.CountMoments(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count moments: %w", err)
	}
	rows, err := s.q.ListMoments(ctx, store.ListMomentsParams{
		Limit:  int32(pageSize),
		Offset: int32((page - 1) * pageSize),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list moments: %w", err)
	}
	items := make([]model.Moment, len(rows))
	for i, r := range rows {
		images, err := decodeStrings(r.Images)
		if err != nil {
			return nil, 0, fmt.Errorf("decode moment %d images: %w", r.ID, err)
		}
		items[i] = model.Moment{ID: r.ID, Content: r.Content, Images: images, CreatedAt: r.CreatedAt}
	}
	return items, total, nil
}

func (s *Public) ListFriendLinks(ctx context.Context) ([]model.FriendLink, error) {
	rows, err := s.q.ListFriendLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list friend links: %w", err)
	}
	items := make([]model.FriendLink, len(rows))
	for i, r := range rows {
		items[i] = model.FriendLink{
			ID: r.ID, Name: r.Name, URL: r.Url, Avatar: r.Avatar,
			Description: r.Description, CreatedAt: r.CreatedAt,
		}
	}
	return items, nil
}

func (s *Public) ListChangelogs(ctx context.Context) ([]model.Changelog, error) {
	rows, err := s.q.ListChangelogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("list changelogs: %w", err)
	}
	items := make([]model.Changelog, len(rows))
	for i, r := range rows {
		entries, err := decodeStrings(r.Items)
		if err != nil {
			return nil, fmt.Errorf("decode changelog %d items: %w", r.ID, err)
		}
		items[i] = model.Changelog{ID: r.ID, Items: entries, LoggedAt: r.LoggedAt}
	}
	return items, nil
}

func (s *Public) ListProjects(ctx context.Context) ([]model.Project, error) {
	rows, err := s.q.ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	items := make([]model.Project, len(rows))
	for i, r := range rows {
		items[i] = model.Project{
			ID: r.ID, Name: r.Name, Description: r.Description,
			Cover: r.Cover, URL: r.Url, SortOrder: r.SortOrder,
		}
	}
	return items, nil
}

func (s *Public) GetPage(ctx context.Context, key string) (*model.Page, error) {
	dbKey, ok := pageKeys[key]
	if !ok {
		return nil, apperr.ErrNotFound
	}
	p, err := s.q.GetPage(ctx, dbKey)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get page: %w", err)
	}
	return &model.Page{Key: key, Content: p.Content, UpdatedAt: p.UpdatedAt}, nil
}

func (s *Public) GetSite(ctx context.Context) (*model.Site, error) {
	settings, err := s.q.GetSiteSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get site settings: %w", err)
	}
	stats, err := s.q.GetSiteStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get site stats: %w", err)
	}
	return &model.Site{
		Notice:        settings.Notice,
		ViewCount:     settings.ViewCount,
		ArticleCount:  stats.ArticleCount,
		CategoryCount: stats.CategoryCount,
		TagCount:      stats.TagCount,
	}, nil
}

func (s *Public) RecordView(ctx context.Context) (uint64, error) {
	res, err := s.q.IncrementViewCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("increment view count: %w", err)
	}
	n, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("read view count: %w", err)
	}
	return uint64(n), nil
}

func (s *Public) tagsByArticle(ctx context.Context, ids []uint64) (map[uint64][]model.TagRef, error) {
	out := make(map[uint64][]model.TagRef, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	rows, err := s.q.ListTagsForArticles(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list article tags: %w", err)
	}
	for _, r := range rows {
		out[r.ArticleID] = append(out[r.ArticleID], model.TagRef{ID: r.ID, Name: r.Name})
	}
	return out, nil
}

func categoryRef(id sql.NullInt64, name sql.NullString) *model.CategoryRef {
	if !id.Valid {
		return nil
	}
	return &model.CategoryRef{ID: uint64(id.Int64), Name: name.String}
}

func tagsOrEmpty(tags []model.TagRef) []model.TagRef {
	if tags == nil {
		return []model.TagRef{}
	}
	return tags
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func decodeStrings(raw json.RawMessage) ([]string, error) {
	out := []string{}
	if len(raw) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []string{}
	}
	return out, nil
}
