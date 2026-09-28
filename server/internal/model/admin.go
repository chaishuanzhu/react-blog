package model

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Website   string    `json:"website"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user"`
}

type ProfileInput struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Website  string `json:"website"`
}

type AdminArticle struct {
	ID          uint64       `json:"id"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	Content     string       `json:"content,omitempty"`
	Status      string       `json:"status"`
	Category    *CategoryRef `json:"category"`
	Tags        []TagRef     `json:"tags"`
	PublishedAt time.Time    `json:"publishedAt"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type ArticleInput struct {
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	CategoryID  *uint64    `json:"categoryId"`
	TagIDs      []uint64   `json:"tagIds"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"publishedAt"`
}

type AdminArticleFilter struct {
	Status     string
	Keyword    string
	CategoryID uint64
	TagID      uint64
	Page       int
	PageSize   int
}

type NameInput struct {
	Name string `json:"name"`
}

type AdminComment struct {
	ID        uint64       `json:"id"`
	ParentID  *uint64      `json:"parentId"`
	Article   *ArticleLink `json:"article"`
	Nickname  string       `json:"nickname"`
	Email     string       `json:"email"`
	Website   string       `json:"website"`
	Avatar    string       `json:"avatar"`
	Content   string       `json:"content"`
	IsAdmin   bool         `json:"isAdmin"`
	IP        string       `json:"ip"`
	CreatedAt time.Time    `json:"createdAt"`
}

type ArticleLink struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type MomentInput struct {
	Content   string     `json:"content"`
	Images    []string   `json:"images"`
	CreatedAt *time.Time `json:"createdAt"`
}

type FriendLinkInput struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
}

type ChangelogInput struct {
	Items    []string   `json:"items"`
	LoggedAt *time.Time `json:"loggedAt"`
}

type ContentInput struct {
	Content string `json:"content"`
}

type NoticeInput struct {
	Notice string `json:"notice"`
}

type AdminStats struct {
	PublishedCount  int64  `json:"publishedCount"`
	DraftCount      int64  `json:"draftCount"`
	CategoryCount   int64  `json:"categoryCount"`
	TagCount        int64  `json:"tagCount"`
	CommentCount    int64  `json:"commentCount"`
	MomentCount     int64  `json:"momentCount"`
	FriendLinkCount int64  `json:"friendLinkCount"`
	ViewCount       uint64 `json:"viewCount"`
}

type Created struct {
	ID uint64 `json:"id"`
}
