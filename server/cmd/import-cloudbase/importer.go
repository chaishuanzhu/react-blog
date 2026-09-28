package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type importer struct {
	tx         *sql.Tx
	adminEmail string

	categoryIDs map[string]int64
	tagIDs      map[string]int64
	articleIDs  map[string]int64 // keyed by original titleEng
	commentIDs  map[string]int64 // keyed by original _id

	counts   map[string]int
	warnings []string
}

func (im *importer) run(ctx context.Context, e *export) error {
	im.categoryIDs = map[string]int64{}
	im.tagIDs = map[string]int64{}
	im.articleIDs = map[string]int64{}
	im.commentIDs = map[string]int64{}
	im.counts = map[string]int{}

	steps := []struct {
		name string
		fn   func(context.Context, *export) error
	}{
		{"categories", im.importCategories},
		{"tags", im.importTags},
		{"articles", im.importArticles},
		{"comments", im.importComments},
		{"moments", im.importMoments},
		{"friend_links", im.importLinks},
		{"changelogs", im.importLogs},
		{"pages", im.importPages},
		{"site_settings", im.importSiteSettings},
	}
	for _, s := range steps {
		if err := s.fn(ctx, e); err != nil {
			return fmt.Errorf("%s: %w", s.name, err)
		}
	}
	return nil
}

func (im *importer) warn(format string, args ...any) {
	im.warnings = append(im.warnings, fmt.Sprintf(format, args...))
}

func (im *importer) report() {
	keys := make([]string, 0, len(im.counts))
	for k := range im.counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println("imported rows:")
	for _, k := range keys {
		fmt.Printf("  %-14s %d\n", k, im.counts[k])
	}
	if len(im.warnings) > 0 {
		fmt.Printf("warnings (%d):\n", len(im.warnings))
		for _, w := range im.warnings {
			fmt.Println("  -", w)
		}
	}
}

func (im *importer) insert(ctx context.Context, table, query string, args ...any) (int64, error) {
	res, err := im.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	im.counts[table]++
	return res.LastInsertId()
}

func (im *importer) ensureCategory(ctx context.Context, name string, created time.Time) (int64, error) {
	if id, ok := im.categoryIDs[name]; ok {
		return id, nil
	}
	id, err := im.insert(ctx, "categories", "INSERT INTO categories (name, created_at) VALUES (?, ?)", name, created)
	if err != nil {
		return 0, err
	}
	im.categoryIDs[name] = id
	return id, nil
}

func (im *importer) ensureTag(ctx context.Context, name string, created time.Time) (int64, error) {
	if id, ok := im.tagIDs[name]; ok {
		return id, nil
	}
	id, err := im.insert(ctx, "tags", "INSERT INTO tags (name, created_at) VALUES (?, ?)", name, created)
	if err != nil {
		return 0, err
	}
	im.tagIDs[name] = id
	return id, nil
}

func (im *importer) importCategories(ctx context.Context, e *export) error {
	for _, c := range e.Categories {
		name := strings.TrimSpace(c.Name)
		if name == "" {
			im.warn("category %s has an empty name, skipped", c.ID)
			continue
		}
		if _, ok := im.categoryIDs[name]; ok {
			im.warn("duplicate category %q merged", name)
			continue
		}
		if _, err := im.ensureCategory(ctx, name, c.Date.Or(time.Now())); err != nil {
			return err
		}
	}
	return nil
}

func (im *importer) importTags(ctx context.Context, e *export) error {
	for _, t := range e.Tags {
		name := strings.TrimSpace(t.Name)
		if name == "" {
			im.warn("tag %s has an empty name, skipped", t.ID)
			continue
		}
		if _, ok := im.tagIDs[name]; ok {
			im.warn("duplicate tag %q merged", name)
			continue
		}
		if _, err := im.ensureTag(ctx, name, t.Date.Or(time.Now())); err != nil {
			return err
		}
	}
	return nil
}

func (im *importer) importArticles(ctx context.Context, e *export) error {
	articles := append([]cbArticle(nil), e.Articles...)
	sort.SliceStable(articles, func(i, j int) bool { return articles[i].Date.Time.Before(articles[j].Date.Time) })

	for _, a := range articles {
		titleEng := strings.TrimSpace(a.TitleEng)
		published := a.Date.Or(time.Now())
		var categoryID sql.NullInt64
		if name := strings.TrimSpace(a.Classes); name != "" {
			if _, known := im.categoryIDs[name]; !known {
				im.warn("article %q: category %q not in classes collection, created", a.Title, name)
			}
			id, err := im.ensureCategory(ctx, name, published)
			if err != nil {
				return err
			}
			categoryID = sql.NullInt64{Int64: id, Valid: true}
		}

		status := "draft"
		if a.Post {
			status = "published"
		}
		id, err := im.insert(ctx, "articles",
			`INSERT INTO articles (title, content, category_id, status, published_at, created_at)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			a.Title, a.Content, categoryID, status, published, published)
		if err != nil {
			return fmt.Errorf("article %q: %w", a.Title, err)
		}
		// CloudBase comments reference articles by titleEng; the first article wins on duplicates.
		switch _, dup := im.articleIDs[titleEng]; {
		case titleEng == "":
			im.warn("article %s (%q) has no titleEng, its comments cannot be matched", a.ID, a.Title)
		case dup:
			im.warn("article %q: duplicate titleEng %q, its comments go to the earlier article", a.Title, titleEng)
		default:
			im.articleIDs[titleEng] = id
		}

		seen := map[string]bool{}
		for _, raw := range a.Tags {
			name := strings.TrimSpace(raw)
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			if _, known := im.tagIDs[name]; !known {
				im.warn("article %q: tag %q not in tags collection, created", a.Title, name)
			}
			tagID, err := im.ensureTag(ctx, name, published)
			if err != nil {
				return err
			}
			if _, err := im.insert(ctx, "article_tags",
				"INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", id, tagID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (im *importer) importComments(ctx context.Context, e *export) error {
	var top, replies []cbComment
	for _, c := range e.Comments {
		if c.ReplyID == "" {
			top = append(top, c)
		} else {
			replies = append(replies, c)
		}
	}
	byDate := func(list []cbComment) {
		sort.SliceStable(list, func(i, j int) bool { return list[i].Date.Time.Before(list[j].Date.Time) })
	}
	byDate(top)
	byDate(replies)

	for _, c := range append(top, replies...) {
		var articleID sql.NullInt64
		if c.PostTitle != "" {
			id, ok := im.articleIDs[c.PostTitle]
			if !ok {
				im.warn("comment %s: article %q not found, skipped", c.ID, c.PostTitle)
				continue
			}
			articleID = sql.NullInt64{Int64: id, Valid: true}
		}
		var parentID sql.NullInt64
		if c.ReplyID != "" {
			id, ok := im.commentIDs[c.ReplyID]
			if !ok {
				im.warn("comment %s: parent %s not found, skipped", c.ID, c.ReplyID)
				continue
			}
			parentID = sql.NullInt64{Int64: id, Valid: true}
		}
		isAdmin := im.adminEmail != "" && strings.EqualFold(strings.TrimSpace(c.Email), im.adminEmail)

		id, err := im.insert(ctx, "comments",
			`INSERT INTO comments (article_id, parent_id, nickname, email, website, avatar, content, is_admin, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			articleID, parentID, truncateRunes(c.Name, 32), truncateRunes(c.Email, 128),
			truncateRunes(c.Link, 255), truncateRunes(c.Avatar, 255), c.Content, isAdmin, c.Date.Or(time.Now()))
		if err != nil {
			return fmt.Errorf("comment %s: %w", c.ID, err)
		}
		im.commentIDs[c.ID] = id
	}
	return nil
}

func (im *importer) importMoments(ctx context.Context, e *export) error {
	for _, s := range e.Says {
		images := []string{}
		for _, img := range s.Imgs {
			if img = strings.TrimSpace(img); img != "" {
				images = append(images, img)
			}
		}
		raw, _ := json.Marshal(images)
		if _, err := im.insert(ctx, "moments",
			"INSERT INTO moments (content, images, created_at) VALUES (?, ?, ?)",
			s.Content, raw, s.Date.Or(time.Now())); err != nil {
			return err
		}
	}
	return nil
}

func (im *importer) importLinks(ctx context.Context, e *export) error {
	for _, l := range e.Links {
		if _, err := im.insert(ctx, "friend_links",
			"INSERT INTO friend_links (name, url, avatar, description, created_at) VALUES (?, ?, ?, ?, ?)",
			l.Name, l.Link, l.Avatar, truncateRunes(l.Descr, 255), l.Date.Or(time.Now())); err != nil {
			return err
		}
	}
	return nil
}

func (im *importer) importLogs(ctx context.Context, e *export) error {
	for _, l := range e.Logs {
		items := l.LogContent
		if items == nil {
			items = []string{}
		}
		raw, _ := json.Marshal(items)
		if _, err := im.insert(ctx, "changelogs",
			"INSERT INTO changelogs (items, logged_at) VALUES (?, ?)", raw, l.Date.Or(time.Now())); err != nil {
			return err
		}
	}
	return nil
}

// The old blog picks the "about site" document as data[0] after ordering by _id descending,
// and "about me" as data[1].
func (im *importer) importPages(ctx context.Context, e *export) error {
	docs := append([]cbAbout(nil), e.About...)
	sort.SliceStable(docs, func(i, j int) bool { return docs[i].ID > docs[j].ID })
	keys := []string{"about_site", "about_me"}
	for i, key := range keys {
		if i >= len(docs) {
			im.warn("about collection has no document for %s", key)
			continue
		}
		if _, err := im.tx.ExecContext(ctx, "UPDATE pages SET content = ? WHERE page_key = ?", docs[i].Content, key); err != nil {
			return err
		}
		im.counts["pages"]++
	}
	if len(docs) > len(keys) {
		im.warn("about collection has %d documents, only the first %d were used", len(docs), len(keys))
	}
	return nil
}

func (im *importer) importSiteSettings(ctx context.Context, e *export) error {
	if len(e.Notice) > 0 {
		if _, err := im.tx.ExecContext(ctx, "UPDATE site_settings SET notice = ? WHERE id = 1",
			truncateRunes(e.Notice[0].Notice, 255)); err != nil {
			return err
		}
		im.counts["site_settings"]++
	}
	if len(e.SiteCount) > 0 {
		if _, err := im.tx.ExecContext(ctx, "UPDATE site_settings SET view_count = ? WHERE id = 1",
			e.SiteCount[0].Count); err != nil {
			return err
		}
		im.counts["site_settings"]++
	}
	return nil
}

func truncateRunes(s string, n int) string {
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
