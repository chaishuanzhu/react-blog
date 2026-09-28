package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// msTime accepts the date shapes found in CloudBase exports:
// epoch milliseconds (number or numeric string), {"$date": <ms>} and {"$date": "<RFC 3339>"}.
type msTime struct {
	time.Time
	Set bool
}

func (t *msTime) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	if b[0] == '{' {
		var wrapped struct {
			Date json.RawMessage `json:"$date"`
		}
		if err := json.Unmarshal(b, &wrapped); err != nil {
			return err
		}
		return t.UnmarshalJSON(wrapped.Date)
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if ms, err := strconv.ParseInt(s, 10, 64); err == nil {
			t.Time, t.Set = time.UnixMilli(ms).UTC(), true
			return nil
		}
		for _, layout := range []string{time.RFC3339Nano, "2006-01-02 15:04:05", "2006-01-02"} {
			if parsed, err := time.ParseInLocation(layout, s, time.Local); err == nil {
				t.Time, t.Set = parsed.UTC(), true
				return nil
			}
		}
		return fmt.Errorf("unrecognized date %q", s)
	}
	var ms float64
	if err := json.Unmarshal(b, &ms); err != nil {
		return err
	}
	t.Time, t.Set = time.UnixMilli(int64(ms)).UTC(), true
	return nil
}

func (t msTime) Or(fallback time.Time) time.Time {
	if t.Set {
		return t.Time
	}
	return fallback
}

type cbCategory struct {
	ID   string `json:"_id"`
	Name string `json:"class"`
	Date msTime `json:"date"`
}

type cbTag struct {
	ID   string `json:"_id"`
	Name string `json:"tag"`
	Date msTime `json:"date"`
}

type cbArticle struct {
	ID       string   `json:"_id"`
	Title    string   `json:"title"`
	TitleEng string   `json:"titleEng"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	Classes  string   `json:"classes"`
	Date     msTime   `json:"date"`
	Post     bool     `json:"post"`
}

type cbComment struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Link      string `json:"link"`
	Content   string `json:"content"`
	Date      msTime `json:"date"`
	Avatar    string `json:"avatar"`
	PostTitle string `json:"postTitle"`
	ReplyID   string `json:"replyId"`
}

type cbSay struct {
	ID      string   `json:"_id"`
	Content string   `json:"content"`
	Imgs    []string `json:"imgs"`
	Date    msTime   `json:"date"`
}

type cbLink struct {
	ID     string `json:"_id"`
	Name   string `json:"name"`
	Link   string `json:"link"`
	Avatar string `json:"avatar"`
	Descr  string `json:"descr"`
	Date   msTime `json:"date"`
}

type cbLog struct {
	ID         string   `json:"_id"`
	Date       msTime   `json:"date"`
	LogContent []string `json:"logContent"`
}

type cbAbout struct {
	ID      string `json:"_id"`
	Content string `json:"content"`
}

type cbNotice struct {
	ID     string `json:"_id"`
	Notice string `json:"notice"`
}

type cbSiteCount struct {
	ID    string `json:"_id"`
	Count int64  `json:"count"`
}

type export struct {
	Categories []cbCategory
	Tags       []cbTag
	Articles   []cbArticle
	Comments   []cbComment
	Says       []cbSay
	Links      []cbLink
	Logs       []cbLog
	About      []cbAbout
	Notice     []cbNotice
	SiteCount  []cbSiteCount
}

func loadExport(dir string) (*export, error) {
	e := &export{}
	loaders := []struct {
		name string
		load func(path string) error
	}{
		{"classes", func(p string) error { return readDocuments(p, &e.Categories) }},
		{"tags", func(p string) error { return readDocuments(p, &e.Tags) }},
		{"articles", func(p string) error { return readDocuments(p, &e.Articles) }},
		{"allComments", func(p string) error { return readDocuments(p, &e.Comments) }},
		{"says", func(p string) error { return readDocuments(p, &e.Says) }},
		{"links", func(p string) error { return readDocuments(p, &e.Links) }},
		{"logs", func(p string) error { return readDocuments(p, &e.Logs) }},
		{"about", func(p string) error { return readDocuments(p, &e.About) }},
		{"notice", func(p string) error { return readDocuments(p, &e.Notice) }},
		{"siteCount", func(p string) error { return readDocuments(p, &e.SiteCount) }},
	}
	for _, l := range loaders {
		path, err := findCollectionFile(dir, l.name)
		if err != nil {
			return nil, err
		}
		if path == "" {
			fmt.Printf("  %-12s (file not found, skipped)\n", l.name)
			continue
		}
		if err := l.load(path); err != nil {
			return nil, fmt.Errorf("%s: %w", filepath.Base(path), err)
		}
	}
	return e, nil
}

func findCollectionFile(dir, name string) (string, error) {
	for _, ext := range []string{".json", ".jsonl"} {
		p := filepath.Join(dir, name+ext)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
	}
	return "", nil
}

// readDocuments accepts either a JSON array or JSON Lines (one document per line),
// which is what the CloudBase console exports.
func readDocuments[T any](path string, dst *[]T) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil
	}
	if trimmed[0] == '[' {
		return json.Unmarshal(trimmed, dst)
	}

	sc := bufio.NewScanner(bytes.NewReader(trimmed))
	sc.Buffer(make([]byte, 0, 1024*1024), 64*1024*1024)
	line := 0
	for sc.Scan() {
		line++
		text := strings.TrimSpace(sc.Text())
		if text == "" {
			continue
		}
		var doc T
		if err := json.Unmarshal([]byte(text), &doc); err != nil {
			return fmt.Errorf("line %d: %w", line, err)
		}
		*dst = append(*dst, doc)
	}
	return sc.Err()
}
