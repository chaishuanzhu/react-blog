package mail

import (
	"bytes"
	"html/template"
)

type CommentMail struct {
	SiteName    string
	Heading     string
	Target      string
	Author      string
	Content     string
	ParentShown bool
	ParentName  string
	Parent      string
	Link        string
}

var commentTmpl = template.Must(template.New("comment").Parse(`<!DOCTYPE html>
<html><body style="margin:0;padding:24px;background:#f5f6f7;font-family:-apple-system,'PingFang SC','Microsoft YaHei',sans-serif;color:#333">
<div style="max-width:560px;margin:0 auto;background:#fff;border-radius:8px;padding:24px 28px">
  <h2 style="margin:0 0 16px;font-size:18px">{{.Heading}}</h2>
  <p style="margin:0 0 12px;color:#666">来自「{{.Target}}」</p>
  {{if .ParentShown}}
  <div style="margin:0 0 12px;padding:10px 14px;background:#f7f7f7;border-left:3px solid #ccc;color:#777;white-space:pre-wrap">{{.ParentName}}：{{.Parent}}</div>
  {{end}}
  <div style="margin:0 0 20px;padding:12px 14px;background:#f0f7ff;border-left:3px solid #4a90e2;white-space:pre-wrap"><strong>{{.Author}}</strong>：{{.Content}}</div>
  <a href="{{.Link}}" style="display:inline-block;padding:8px 18px;background:#4a90e2;color:#fff;text-decoration:none;border-radius:4px">查看完整内容</a>
  <p style="margin:24px 0 0;font-size:12px;color:#aaa">此邮件由 {{.SiteName}} 自动发送，请勿直接回复。</p>
</div>
</body></html>`))

func RenderComment(data CommentMail) (string, error) {
	var buf bytes.Buffer
	if err := commentTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
