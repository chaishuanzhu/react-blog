package service

import "testing"

func TestSummarize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		max  int
		want string
	}{
		{"plain", "hello world", 20, "hello world"},
		{"heading and emphasis", "# Title\n\nSome **bold** and _italic_ text", 50, "Title Some bold and italic text"},
		{"code fence removed", "before\n```go\nfmt.Println(1)\n```\nafter", 50, "before after"},
		{"link keeps text", "see [docs](https://x.y) and ![img](a.png)", 50, "see docs and"},
		{"inline code and html", "use `go test` <br/> now", 50, "use go test now"},
		{"list and quote", "- one\n- two\n> quoted", 50, "one two quoted"},
		{"truncates runes", "你好世界，这是一段测试", 4, "你好世界…"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Summarize(tt.in, tt.max); got != tt.want {
				t.Errorf("Summarize(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestLikeContains(t *testing.T) {
	if got, want := likeContains(`50%_off\`), `%50\%\_off\\%`; got != want {
		t.Errorf("likeContains = %q, want %q", got, want)
	}
}
