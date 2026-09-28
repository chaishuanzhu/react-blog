package model

import "time"

type Comment struct {
	ID        uint64    `json:"id"`
	ParentID  *uint64   `json:"parentId"`
	Nickname  string    `json:"nickname"`
	Website   string    `json:"website"`
	Avatar    string    `json:"avatar"`
	Content   string    `json:"content"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentThread struct {
	Comment
	Replies []Comment `json:"replies"`
}

type CommentInput struct {
	ArticleID uint64 `json:"articleId"`
	ParentID  uint64 `json:"parentId"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Website   string `json:"website"`
	Content   string `json:"content"`
}

type CommentAuthor struct {
	IP        string
	UserAgent string
	IsAdmin   bool
	Nickname  string
	Email     string
	Website   string
	Avatar    string
}
