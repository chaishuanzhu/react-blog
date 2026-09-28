package model

import "time"

type CategoryRef struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type TagRef struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type ArticleSummary struct {
	ID          uint64       `json:"id"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	Category    *CategoryRef `json:"category"`
	Tags        []TagRef     `json:"tags"`
	PublishedAt time.Time    `json:"publishedAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type ArticleDetail struct {
	ArticleSummary
	Content string `json:"content"`
}

type Category struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	ArticleCount int64     `json:"articleCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

type CategoryList struct {
	Items              []Category `json:"items"`
	UncategorizedCount int64      `json:"uncategorizedCount"`
}

type Tag struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	ArticleCount int64     `json:"articleCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Moment struct {
	ID        uint64    `json:"id"`
	Content   string    `json:"content"`
	Images    []string  `json:"images"`
	CreatedAt time.Time `json:"createdAt"`
}

type FriendLink struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Avatar      string    `json:"avatar"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Changelog struct {
	ID       uint64    `json:"id"`
	Items    []string  `json:"items"`
	LoggedAt time.Time `json:"loggedAt"`
}

type Page struct {
	Key       string    `json:"key"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Site struct {
	Notice        string `json:"notice"`
	ViewCount     uint64 `json:"viewCount"`
	ArticleCount  int64  `json:"articleCount"`
	CategoryCount int64  `json:"categoryCount"`
	TagCount      int64  `json:"tagCount"`
}

type ArticleFilter struct {
	Keyword  string
	Category string
	Tag      string
	Page     int
	PageSize int
}
